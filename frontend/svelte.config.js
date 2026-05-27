import adapter from '@sveltejs/adapter-static';
import { vitePreprocess } from '@sveltejs/vite-plugin-svelte';
import path from 'path';

// Optional base path for GitHub Pages "project" hosting (e.g. /metavida-page).
// Leave empty when serving from a custom domain or user/org root page.
const base = process.env.BASE_PATH ?? '';

/** @type {import('@sveltejs/kit').Config} */
const config = {
  preprocess: vitePreprocess(),
  kit: {
    adapter: adapter({
      pages: 'build',
      assets: 'build',
      // 404.html doubles as the SPA fallback GitHub Pages serves for the
      // non-prerendered app routes (admin/client/login).
      fallback: '404.html',
      precompress: false,
      strict: false
    }),
    paths: { base },
    files: {
      assets: 'static',
      routes: 'routes',
      appTemplate: 'app.html'
    },
    alias: {
      $core: path.resolve('./core'),
      $components: path.resolve('./ui-components'),
      $libs: path.resolve('./libs'),
      $routes: path.resolve('./routes')
    }
  }
};

export default config;
