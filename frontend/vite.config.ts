import { readFileSync } from 'node:fs';
import path from 'node:path';
import tailwindcss from '@tailwindcss/vite';
import { sveltekit } from '@sveltejs/kit/vite';
import { defineConfig } from 'vite';

// Production API endpoint used whenever the app is not running on localhost.
// Sourced from the repo-root credentials.json ("ENPOINT"); only this single
// value is read — the rest of that file holds secrets and is never bundled.
// Overridable with VITE_API_BASE for ad-hoc builds.
function resolveProdApiBase(): string {
  if (process.env.VITE_API_BASE) return process.env.VITE_API_BASE;
  try {
    const credsPath = path.resolve(__dirname, '..', 'credentials.json');
    const creds = JSON.parse(readFileSync(credsPath, 'utf8'));
    if (typeof creds.ENPOINT === 'string' && creds.ENPOINT.trim()) {
      return creds.ENPOINT.trim();
    }
  } catch {
    // credentials.json may be absent (CI / clean checkout) — handled below.
  }
  throw new Error(
    'No production API endpoint resolved. Set VITE_API_BASE or provide credentials.json with a non-empty "ENPOINT".'
  );
}

export default defineConfig({
  plugins: [tailwindcss(), sveltekit()],
  define: {
    __PROD_API_BASE__: JSON.stringify(resolveProdApiBase())
  },
  server: {
    port: 3571,
    fs: {
      strict: false,
      allow: ['.', '..']
    }
  }
});
