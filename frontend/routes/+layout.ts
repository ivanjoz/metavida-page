import { dev } from '$app/environment';

// Keep production SSR intact, but skip dev SSR to avoid duplicate dev rendering work.
export const ssr = !dev;

// Statically render public pages into /docs for GitHub Pages. The dynamic app
// sections opt out via their own +layout.ts (admin/client) and fall back to the
// 404.html SPA shell.
export const prerender = true;
