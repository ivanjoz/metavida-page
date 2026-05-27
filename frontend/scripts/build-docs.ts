#!/usr/bin/env bun
/**
 * Post-build step for GitHub Pages.
 *
 * Takes the SvelteKit adapter-static output (frontend/build) and publishes it
 * into the repo-root `docs/` folder that GitHub Pages serves from, adding the
 * artifacts GitHub Pages needs:
 *   - .nojekyll  -> stops Jekyll from stripping the `_app/` asset folder
 *   - 404.html   -> SPA fallback for non-prerendered routes (created by the
 *                   adapter; also used as index fallback if `/` wasn't rendered)
 *   - CNAME      -> copied from static/CNAME or the CNAME env var, if present
 *
 * Run after `vite build`:  bun run scripts/build-docs.ts
 */
import { rm, mkdir, cp, readdir, writeFile, stat } from 'node:fs/promises';
import { existsSync } from 'node:fs';
import path from 'node:path';

const frontendDir = path.resolve(import.meta.dir, '..');
const repoRoot = path.resolve(frontendDir, '..');
const buildDir = path.join(frontendDir, 'build');
const docsDir = path.join(repoRoot, 'docs');

async function main() {
  if (!existsSync(buildDir)) {
    console.error(`No build output found at ${buildDir}. Run \`vite build\` first.`);
    process.exit(1);
  }

  console.log(`Publishing ${path.relative(repoRoot, buildDir)} -> ${path.relative(repoRoot, docsDir)}`);

  // Replace docs/ wholesale so removed files don't linger.
  await rm(docsDir, { recursive: true, force: true });
  await mkdir(docsDir, { recursive: true });

  for (const entry of await readdir(buildDir)) {
    await cp(path.join(buildDir, entry), path.join(docsDir, entry), { recursive: true });
  }

  // Disable Jekyll so files/folders starting with `_` (e.g. _app) are served.
  await writeFile(path.join(docsDir, '.nojekyll'), '');

  // Ensure a 404.html exists for SPA fallback routing on GitHub Pages.
  if (!existsSync(path.join(docsDir, '404.html'))) {
    const indexHtml = path.join(docsDir, 'index.html');
    if (existsSync(indexHtml)) {
      await cp(indexHtml, path.join(docsDir, '404.html'));
    } else {
      console.warn('No 404.html or index.html in build output — SPA fallback may not work.');
    }
  }

  // Optional custom domain. Prefer static/CNAME (auto-copied by the adapter),
  // otherwise honor a CNAME env var.
  const cnameFromEnv = process.env.CNAME?.trim();
  if (cnameFromEnv && !existsSync(path.join(docsDir, 'CNAME'))) {
    await writeFile(path.join(docsDir, 'CNAME'), cnameFromEnv + '\n');
    console.log(`Wrote CNAME -> ${cnameFromEnv}`);
  }

  const fileCount = await countFiles(docsDir);
  console.log(`Done. ${fileCount} files ready in docs/ for GitHub Pages.`);
}

async function countFiles(dir: string): Promise<number> {
  let total = 0;
  for (const entry of await readdir(dir)) {
    const full = path.join(dir, entry);
    const info = await stat(full);
    total += info.isDirectory() ? await countFiles(full) : 1;
  }
  return total;
}

main().catch((err) => {
  console.error(err);
  process.exit(1);
});
