export const Env = {
  apiBase: import.meta.env.VITE_API_BASE || 'http://localhost:3591',
  makeRoute(route: string) {
    if (route.startsWith('http')) return route;
    return `${this.apiBase}${route.startsWith('/') ? route : `/${route}`}`;
  }
};

