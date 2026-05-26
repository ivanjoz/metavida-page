package main

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"unicode"
)

const outputPath = "backend/schemas.generated.go"
const outputPackage = "metavida/backend"

type tableEntry struct {
	TypeName    string
	PackagePath string
}

func main() {
	projectRoot, err := findProjectRootDir()
	if err != nil {
		exitWith(err)
	}

	modulePath, err := readBackendModulePath(projectRoot)
	if err != nil {
		exitWith(err)
	}

	entries, err := scanTableStructs(filepath.Join(projectRoot, "backend"), modulePath)
	if err != nil {
		exitWith(err)
	}

	aliases := assignPackageAliases(entries)
	sort.Slice(entries, func(i, j int) bool {
		return qualifiedReference(entries[i], aliases) < qualifiedReference(entries[j], aliases)
	})

	source, err := renderSchemasFile(entries, aliases, modulePath)
	if err != nil {
		exitWith(err)
	}

	fullOutputPath := filepath.Join(projectRoot, outputPath)
	if err := os.WriteFile(fullOutputPath, source, 0644); err != nil {
		exitWith(err)
	}

	fmt.Printf("generated %s (%d schemas)\n", outputPath, len(entries))
}

func findProjectRootDir() (string, error) {
	currentDir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	for {
		if pathExists(filepath.Join(currentDir, "backend")) && pathExists(filepath.Join(currentDir, "scripts")) {
			return currentDir, nil
		}
		parent := filepath.Dir(currentDir)
		if parent == currentDir {
			return "", fmt.Errorf("could not find project root containing backend/ and scripts/")
		}
		currentDir = parent
	}
}

func readBackendModulePath(projectRoot string) (string, error) {
	content, err := os.ReadFile(filepath.Join(projectRoot, "backend", "go.mod"))
	if err != nil {
		return "", err
	}
	for _, line := range strings.Split(string(content), "\n") {
		fields := strings.Fields(line)
		if len(fields) == 2 && fields[0] == "module" {
			return fields[1], nil
		}
	}
	return "", fmt.Errorf("backend/go.mod does not declare a module")
}

func scanTableStructs(backendDir string, modulePath string) ([]tableEntry, error) {
	var entries []tableEntry

	err := filepath.WalkDir(backendDir, func(path string, entry os.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if entry.IsDir() {
			switch entry.Name() {
			case ".git", "node_modules", "vendor", "db":
				return filepath.SkipDir
			}
			return nil
		}
		if filepath.Ext(path) != ".go" || strings.HasSuffix(path, "_test.go") || strings.HasSuffix(path, ".generated.go") {
			return nil
		}

		fileSet := token.NewFileSet()
		parsedFile, err := parser.ParseFile(fileSet, path, nil, parser.SkipObjectResolution)
		if err != nil {
			return fmt.Errorf("%s: %w", path, err)
		}

		packagePath, err := derivePackageImportPath(backendDir, path, modulePath)
		if err != nil {
			return err
		}

		for _, declaration := range parsedFile.Decls {
			genericDecl, ok := declaration.(*ast.GenDecl)
			if !ok || genericDecl.Tok != token.TYPE {
				continue
			}
			for _, spec := range genericDecl.Specs {
				typeSpec, ok := spec.(*ast.TypeSpec)
				if !ok || !isExported(typeSpec.Name.Name) {
					continue
				}
				structType, ok := typeSpec.Type.(*ast.StructType)
				if !ok || structType.Fields == nil || len(structType.Fields.List) == 0 {
					continue
				}
				if isBaseTableStruct(structType.Fields.List[0], typeSpec.Name.Name) {
					entries = append(entries, tableEntry{
						TypeName:    typeSpec.Name.Name,
						PackagePath: packagePath,
					})
				}
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return entries, nil
}

func isBaseTableStruct(firstField *ast.Field, ownTypeName string) bool {
	if len(firstField.Names) != 0 {
		return false
	}

	indexList, ok := firstField.Type.(*ast.IndexListExpr)
	if !ok || len(indexList.Indices) != 2 {
		return false
	}

	selector, ok := indexList.X.(*ast.SelectorExpr)
	if !ok {
		return false
	}
	packageIdent, ok := selector.X.(*ast.Ident)
	if !ok || packageIdent.Name != "db" || selector.Sel.Name != "TableStruct" {
		return false
	}

	secondArgIdent, ok := indexList.Indices[1].(*ast.Ident)
	return ok && secondArgIdent.Name == ownTypeName
}

func derivePackageImportPath(backendDir, filePath, modulePath string) (string, error) {
	relativeDir, err := filepath.Rel(backendDir, filepath.Dir(filePath))
	if err != nil {
		return "", err
	}
	if relativeDir == "." {
		return modulePath, nil
	}
	return modulePath + "/" + filepath.ToSlash(relativeDir), nil
}

func assignPackageAliases(entries []tableEntry) map[string]string {
	aliasByPackage := map[string]string{}
	for _, entry := range entries {
		if _, exists := aliasByPackage[entry.PackagePath]; exists {
			continue
		}
		aliasByPackage[entry.PackagePath] = deriveAliasForPackage(entry.PackagePath)
	}
	return aliasByPackage
}

func deriveAliasForPackage(packagePath string) string {
	if packagePath == outputPackage {
		return ""
	}
	segments := strings.Split(packagePath, "/")
	lastSegment := segments[len(segments)-1]
	if lastSegment == "types" && len(segments) >= 2 {
		return sanitizeAlias(segments[len(segments)-2] + "Types")
	}
	return sanitizeAlias(lastSegment)
}

func sanitizeAlias(alias string) string {
	alias = strings.ReplaceAll(alias, "-", "_")
	if alias == "" {
		return "pkg"
	}
	return alias
}

func qualifiedReference(entry tableEntry, aliases map[string]string) string {
	alias := aliases[entry.PackagePath]
	if alias == "" {
		return entry.TypeName
	}
	return alias + "." + entry.TypeName
}

func renderSchemasFile(entries []tableEntry, aliases map[string]string, modulePath string) ([]byte, error) {
	imports := buildImportLines(entries, aliases, modulePath)

	var buffer bytes.Buffer
	buffer.WriteString("// Code generated by scripts/schemas; DO NOT EDIT.\n")
	buffer.WriteString("package main\n\n")

	buffer.WriteString("import (\n")
	for _, line := range imports {
		buffer.WriteString("\t")
		buffer.WriteString(line)
		buffer.WriteString("\n")
	}
	buffer.WriteString(")\n\n")

	buffer.WriteString("func MakeDeploySchemas() []db.TableDeployInterface {\n")
	buffer.WriteString("\treturn []db.TableDeployInterface{\n")
	for _, entry := range entries {
		fmt.Fprintf(&buffer, "\t\tdb.Table[%s](),\n", qualifiedReference(entry, aliases))
	}
	buffer.WriteString("\t}\n")
	buffer.WriteString("}\n")

	formatted, err := format.Source(buffer.Bytes())
	if err != nil {
		return nil, fmt.Errorf("format generated source: %w\n--- source ---\n%s", err, buffer.String())
	}
	return formatted, nil
}

func buildImportLines(entries []tableEntry, aliases map[string]string, modulePath string) []string {
	requiredPackages := map[string]bool{modulePath + "/db": true}
	for _, entry := range entries {
		if entry.PackagePath != outputPackage {
			requiredPackages[entry.PackagePath] = true
		}
	}

	packagePaths := make([]string, 0, len(requiredPackages))
	for packagePath := range requiredPackages {
		packagePaths = append(packagePaths, packagePath)
	}
	sort.Strings(packagePaths)

	lines := make([]string, 0, len(packagePaths))
	for _, packagePath := range packagePaths {
		alias := aliases[packagePath]
		if alias == "" {
			alias = deriveAliasForPackage(packagePath)
		}
		segments := strings.Split(packagePath, "/")
		if alias == "" || alias == segments[len(segments)-1] {
			lines = append(lines, fmt.Sprintf("%q", packagePath))
		} else {
			lines = append(lines, fmt.Sprintf("%s %q", alias, packagePath))
		}
	}
	return lines
}

func isExported(name string) bool {
	for _, r := range name {
		return unicode.IsUpper(r)
	}
	return false
}

func pathExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func exitWith(err error) {
	fmt.Fprintln(os.Stderr, err)
	os.Exit(1)
}
