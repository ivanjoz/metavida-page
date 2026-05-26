import { dev } from '$app/environment';

// Keep production SSR intact, but skip dev SSR to avoid duplicate dev rendering work.
export const ssr = !dev;
