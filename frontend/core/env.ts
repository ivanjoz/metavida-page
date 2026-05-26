import { browser } from '$app/environment';

export { browser };

const makeAbsoluteRoute = (baseRoute: string, route: string) => {
  // Normalize app-relative API paths so shared Genix helpers can call one route builder.
  if (route.startsWith('http')) return route;
  return `${baseRoute}${route.startsWith('/') ? route : `/${route}`}`;
};

const makeCDNRoute = (...segments: string[]) => {
  // Shared image components pass path fragments; empty CDN means local/static paths.
  return segments.filter(Boolean).join('/').replaceAll('//', '/');
};

export const Env = {
  apiBase: import.meta.env.VITE_API_BASE || 'http://localhost:3591',
  CDN_URL: import.meta.env.VITE_CDN_URL || '',
  serviceWorker: '/sw.js',
  enviroment: 'metavida',
  componentIDCounter: 0,
  imageCounter: 10000,
  fetchID: 1000,
  dexieVersion: 1,
  imageWorker: null as unknown as Worker,
  ImageWorkerClass: null as unknown as new () => Worker,
  DELTA_CACHE_VERIFY_ROUTE_MEMORY: false,
  screen: browser ? window.screen : { height: -1, width: -1 },
  language: browser ? window.navigator?.language || '' : '',
  deviceMemory: browser ? (window.navigator as Navigator & { deviceMemory?: number }).deviceMemory || 0 : 0,
  clearAccesos: null as (() => void) | null,
  getToken: () => browser ? localStorage.getItem('metavida-token') || '' : '',
  canUserAccessRoute: (_routeValue?: string | null) => true,
  getPathname: () => browser ? window.location.pathname : '',
  getCompanyID: () => 0,
  getComponentID: () => {
    Env.componentIDCounter += 1;
    return Env.componentIDCounter;
  },
  makeRoute(route: string) {
    return makeAbsoluteRoute(this.apiBase, route);
  },
  makeCDNRoute
};
