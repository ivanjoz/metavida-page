package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/klauspost/compress/zstd"
)

const metavidaServiceName = "metavida"

type ServerCredentials struct {
	Host             string `json:"host"`
	User             string `json:"user"`
	Key              string `json:"key"`
	Arch             string `json:"arch"`
	RemoteBinaryPath string `json:"bin"`
}

type Credentials struct {
	Server ServerCredentials `json:"SERVER"`
}

func DeployVPS() {
	fmt.Println("Starting Metavida VPS deployment...")

	// Read deploy targets from credentials.json so host-specific rules stay outside source control.
	credentialsFilePath := "../credentials.json"
	credentialsContent, readCredentialsError := os.ReadFile(credentialsFilePath)
	if readCredentialsError != nil {
		fmt.Printf("Error reading credentials.json: %v\n", readCredentialsError)
		return
	}

	var credentials Credentials
	if parseCredentialsError := json.Unmarshal(credentialsContent, &credentials); parseCredentialsError != nil {
		fmt.Printf("Error parsing credentials.json: %v\n", parseCredentialsError)
		return
	}

	server := credentials.Server
	if server.Host == "" {
		fmt.Println("Error: SERVER is empty in credentials.json")
		return
	}

	// Store build artifacts outside backend/ so deploys do not dirty source folders.
	localTempDirectoryPath := "../tmp"
	if createTempDirectoryError := os.MkdirAll(localTempDirectoryPath, 0o755); createTempDirectoryError != nil {
		fmt.Printf("Error creating tmp directory: %v\n", createTempDirectoryError)
		return
	}

	targetArchitecture := resolveTargetArchitecture(server.Arch)
	loginUser := server.User
	if loginUser == "" {
		loginUser = "root"
	}

	if server.RemoteBinaryPath == "" {
		fmt.Printf("Error: bin path is empty for host %s\n", server.Host)
		return
	}

	serverTarget := fmt.Sprintf("%s@%s", loginUser, server.Host)
	resolvedKeyPath := expandHomePath(server.Key)
	localCompressedBinaryPath, buildArtifactError := buildCompressedBinaryForArchitecture(localTempDirectoryPath, targetArchitecture)
	if buildArtifactError != nil {
		fmt.Printf("Error preparing %s artifact for host %s: %v\n", targetArchitecture, server.Host, buildArtifactError)
		return
	}

	remoteCompressedBinaryPath := server.RemoteBinaryPath + ".zst"
	remoteBinaryDirectory := filepath.Dir(server.RemoteBinaryPath)
	remoteCredentialsPath := filepath.Join(remoteBinaryDirectory, "credentials.json")

	fmt.Printf("Deploying to %s\n", serverTarget)
	fmt.Printf("Debug: host=%s user=%s arch=%s key=%q bin=%s\n", server.Host, loginUser, targetArchitecture, resolvedKeyPath, server.RemoteBinaryPath)

	fmt.Printf("Ensuring remote directory exists: %s\n", remoteBinaryDirectory)
	ensureRemoteDirectoryCommand := buildSSHCommand(resolvedKeyPath, serverTarget, fmt.Sprintf("mkdir -p %s", remoteBinaryDirectory))
	ensureRemoteDirectoryCommand.Stdout = os.Stdout
	ensureRemoteDirectoryCommand.Stderr = os.Stderr

	if ensureRemoteDirectoryError := ensureRemoteDirectoryCommand.Run(); ensureRemoteDirectoryError != nil {
		fmt.Printf("Error creating remote directory on %s: %v\n", serverTarget, ensureRemoteDirectoryError)
		return
	}

	fmt.Printf("Uploading compressed binary to %s:%s...\n", serverTarget, remoteCompressedBinaryPath)
	uploadBinaryCommand := buildRsyncCommand(resolvedKeyPath, localCompressedBinaryPath, serverTarget, remoteCompressedBinaryPath)
	uploadBinaryCommand.Stdout = os.Stdout
	uploadBinaryCommand.Stderr = os.Stderr

	if uploadError := uploadBinaryCommand.Run(); uploadError != nil {
		fmt.Printf("Error uploading binary to %s: %v\n", serverTarget, uploadError)
		return
	}

	fmt.Printf("Decompressing binary on %s...\n", serverTarget)
	decompressRemoteBinaryCommand := buildSSHCommand(
		resolvedKeyPath,
		serverTarget,
		fmt.Sprintf(
			"zstd -d --force %s -o %s && rm %s && chmod +x %s",
			remoteCompressedBinaryPath,
			server.RemoteBinaryPath,
			remoteCompressedBinaryPath,
			server.RemoteBinaryPath,
		),
	)
	decompressRemoteBinaryCommand.Stdout = os.Stdout
	decompressRemoteBinaryCommand.Stderr = os.Stderr

	if decompressError := decompressRemoteBinaryCommand.Run(); decompressError != nil {
		fmt.Printf("Error decompressing on %s: %v\n", serverTarget, decompressError)
		return
	}

	fmt.Printf("Uploading credentials to %s:%s...\n", serverTarget, remoteCredentialsPath)
	uploadCredentialsCommand := buildRsyncCommand(resolvedKeyPath, credentialsFilePath, serverTarget, remoteCredentialsPath)
	uploadCredentialsCommand.Stdout = os.Stdout
	uploadCredentialsCommand.Stderr = os.Stderr

	if uploadCredentialsError := uploadCredentialsCommand.Run(); uploadCredentialsError != nil {
		fmt.Printf("Error uploading credentials to %s: %v\n", serverTarget, uploadCredentialsError)
		return
	}

	isServiceRunning, serviceCheckError := isRemoteServiceRunning(resolvedKeyPath, serverTarget, metavidaServiceName)
	if serviceCheckError != nil {
		fmt.Printf("Error checking %s service on %s: %v\n", metavidaServiceName, serverTarget, serviceCheckError)
		return
	}

	if !isServiceRunning {
		// First-time server setup may not have the service installed yet; deploy the binary and stop here.
		fmt.Printf("%s service is not running on %s. Binary deployed; skipping restart.\n", metavidaServiceName, serverTarget)
		fmt.Printf("Deployment completed for %s.\n", serverTarget)
		fmt.Println("Deployment complete!")
		return
	}

	fmt.Printf("Restarting running %s service on %s...\n", metavidaServiceName, serverTarget)
	restartServiceCommand := buildSSHCommand(resolvedKeyPath, serverTarget, fmt.Sprintf("systemctl restart %s", metavidaServiceName))
	restartServiceCommand.Stdout = os.Stdout
	restartServiceCommand.Stderr = os.Stderr

	if restartServiceError := restartServiceCommand.Run(); restartServiceError != nil {
		fmt.Printf("Error restarting service on %s: %v\n", serverTarget, restartServiceError)
		return
	}

	fmt.Printf("Deployment completed for %s.\n", serverTarget)
	fmt.Println("Deployment complete!")
}

func resolveTargetArchitecture(configuredArchitecture string) string {
	normalizedArchitecture := strings.TrimSpace(strings.ToLower(configuredArchitecture))
	if normalizedArchitecture == "arm64" {
		return "arm64"
	}

	return "amd64"
}

func buildCompressedBinaryForArchitecture(
	localTempDirectoryPath string,
	targetArchitecture string,
) (string, error) {
	// Build the single configured architecture for the configured VPS target.
	localBinaryName := fmt.Sprintf("metavida_app_linux_%s", targetArchitecture)
	localBinaryPath := filepath.Join(localTempDirectoryPath, localBinaryName)
	localCompressedBinaryPath := localBinaryPath + ".zst"

	fmt.Printf("Compiling backend for linux/%s...\n", targetArchitecture)
	buildBackendCommand := exec.Command("go", "build", "-o", localBinaryPath, ".")
	buildBackendCommand.Dir = "../backend"
	buildBackendCommand.Env = append(os.Environ(), "GOOS=linux", "GOARCH="+targetArchitecture)
	buildBackendCommand.Stdout = os.Stdout
	buildBackendCommand.Stderr = os.Stderr

	if buildBackendError := buildBackendCommand.Run(); buildBackendError != nil {
		return "", buildBackendError
	}
	fmt.Printf("Compilation successful for linux/%s.\n", targetArchitecture)

	fmt.Printf("Compressing linux/%s binary with Zstd...\n", targetArchitecture)
	if compressBinaryError := compressBinary(localBinaryPath, localCompressedBinaryPath); compressBinaryError != nil {
		return "", compressBinaryError
	}
	fmt.Printf("Compression successful for linux/%s.\n", targetArchitecture)

	return localCompressedBinaryPath, nil
}

func compressBinary(localBinaryPath string, localCompressedBinaryPath string) error {
	inputBinaryFile, openInputError := os.Open(localBinaryPath)
	if openInputError != nil {
		return openInputError
	}
	defer inputBinaryFile.Close()

	outputCompressedFile, createOutputError := os.Create(localCompressedBinaryPath)
	if createOutputError != nil {
		return createOutputError
	}
	defer outputCompressedFile.Close()

	compressionWriter, createWriterError := zstd.NewWriter(outputCompressedFile)
	if createWriterError != nil {
		return createWriterError
	}
	defer compressionWriter.Close()

	_, copyError := io.Copy(compressionWriter, inputBinaryFile)
	return copyError
}

func isRemoteServiceRunning(resolvedKeyPath string, serverTarget string, serviceName string) (bool, error) {
	// Missing or inactive services return exit code 3; only unexpected SSH/systemctl failures are fatal.
	checkServiceCommand := buildSSHCommand(resolvedKeyPath, serverTarget, fmt.Sprintf("systemctl is-active --quiet %s", serviceName))
	serviceCheckError := checkServiceCommand.Run()
	if serviceCheckError == nil {
		return true, nil
	}

	if exitError, ok := serviceCheckError.(*exec.ExitError); ok && exitError.ExitCode() == 3 {
		return false, nil
	}

	return false, serviceCheckError
}

func expandHomePath(originalPath string) string {
	if originalPath == "" || !strings.HasPrefix(originalPath, "~") {
		return originalPath
	}

	homeDirectory, homeDirectoryError := os.UserHomeDir()
	if homeDirectoryError != nil {
		return originalPath
	}

	return filepath.Join(homeDirectory, strings.TrimPrefix(originalPath, "~"))
}

func buildSSHCommand(resolvedKeyPath string, serverTarget string, remoteCommand string) *exec.Cmd {
	sshArguments := []string{}
	if resolvedKeyPath != "" {
		sshArguments = append(sshArguments, "-i", resolvedKeyPath)
	}

	sshArguments = append(sshArguments, serverTarget, remoteCommand)
	return exec.Command("ssh", sshArguments...)
}

func buildRsyncCommand(resolvedKeyPath string, localCompressedBinaryPath string, serverTarget string, remoteCompressedBinaryPath string) *exec.Cmd {
	sshTransportCommand := "ssh"
	if resolvedKeyPath != "" {
		sshTransportCommand = fmt.Sprintf("ssh -i %s", resolvedKeyPath)
	}

	return exec.Command(
		"rsync",
		"-ahP",
		"-e", sshTransportCommand,
		localCompressedBinaryPath,
		fmt.Sprintf("%s:%s", serverTarget, remoteCompressedBinaryPath),
	)
}
