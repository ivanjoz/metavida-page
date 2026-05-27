import { browser } from '$app/environment';

export { browser };

// Injected by vite (define) from credentials.json "ENPOINT" at build time.
declare const __PROD_API_BASE__: string;

const TOKEN_KEY = 'metavida-token';
const TOKEN_EXP_KEY = 'metavida-token-exp';
const USER_KEY = 'metavida-user';

const LOCAL_API_BASE = 'http://localhost:14447';
const PROD_API_BASE = (typeof __PROD_API_BASE__ !== 'undefined' && __PROD_API_BASE__)
  || import.meta.env.VITE_API_BASE

const isLocalHost = browser && /^(localhost|127\.0\.0\.1|\[?::1\]?)$/.test(window.location.hostname);

// Use the local backend only when actually served from localhost; everything
// else (deployed site, LAN-accessed dev server) hits the production endpoint.
// Trailing slashes are trimmed so makeAbsoluteRoute doesn't produce `//`.
const resolveApiBase = () => (isLocalHost ? LOCAL_API_BASE : PROD_API_BASE).replace(/\/+$/, '');

const makeAbsoluteRoute = (baseRoute: string, route: string) => {
  // Normalize app-relative API paths so shared Genix helpers can call one route builder.
  if (route.startsWith('http')) return route;
  return `${baseRoute}${route.startsWith('/') ? route : `/${route}`}`;
};

const makeCDNRoute = (...segments: string[]) => {
  // Shared image components pass path fragments; empty CDN means local/static paths.
  return segments.filter(Boolean).join('/').replaceAll('//', '/');
};

// Returns 0 (SSR), 2 (logged in), or 3 (not logged in) — matches Genix convention.
export const checkIsLogin = (): number => {
  if (!browser) return 0;
  const token = localStorage.getItem(TOKEN_KEY);
  const exp = parseInt(localStorage.getItem(TOKEN_EXP_KEY) || '0');
  if (!token || !exp || Math.floor(Date.now() / 1000) > exp) return 3;
  return 2;
};

export const Env = {
  apiBase: resolveApiBase(),
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
  clearAccesos: (() => {
    if (!browser) return;
    localStorage.removeItem(TOKEN_KEY);
    localStorage.removeItem(TOKEN_EXP_KEY);
    localStorage.removeItem(USER_KEY);
    window.location.href = '/login';
  }) as (() => void) | null,
  getToken: () => browser ? localStorage.getItem(TOKEN_KEY) || '' : '',
  setSession: (token: string, expTime: number, userInfo: string) => {
    if (!browser) return;
    localStorage.setItem(TOKEN_KEY, token);
    localStorage.setItem(TOKEN_EXP_KEY, String(expTime));
    localStorage.setItem(USER_KEY, userInfo);
  },
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
