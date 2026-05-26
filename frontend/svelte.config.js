import adapter from '@sveltejs/adapter-static';
import { vitePreprocess } from '@sveltejs/vite-plugin-svelte';
import path from 'path';

/** @type {import('@sveltejs/kit').Config} */
const config = {
  preprocess: vitePreprocess(),
  kit: {
    adapter: adapter({
      pages: 'build',
      assets: 'build',
      fallback: 'index.html',
      precompress: false,
      strict: true
    }),
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
