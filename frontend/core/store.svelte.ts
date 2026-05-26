import { SvelteMap } from 'svelte/reactivity';
import type { IImageResult } from './types/common';
import { Env, browser } from './env';

export const getDeviceType = () => {
  // Match Genix breakpoints so shared responsive components keep their behavior.
  if (!browser) return 1;
  if (window.innerWidth < 740) return 3;
  if (window.innerWidth < 1140) return 2;
  return 1;
};

export interface ITopSearchLayer {
  options: any[];
  keyName: string;
  keyID: string | number;
  onSelect: (record: any) => void;
  onClear?: () => void;
  onRemove?: (record: any) => void;
}

export interface ITopDateLayer {
  selectedUnixDay: number;
  focusedUnixDay?: number;
  selectedMonthKey: number;
  label?: string;
  placeholder?: string;
  onSelect: (unixDay: number) => void;
  onClose?: () => void;
}

export const Core = $state({
  deviceType: getDeviceType(),
  popoverShowID: 0 as number | string,
  showSideLayer: 0,
  sideLayerSize: 0,
  showMobileSearchLayer: null as ITopSearchLayer | null,
  showMobileDateLayer: null as ITopDateLayer | null,
  ecommerce: { cartOption: 1 },
  openSideLayer: (layerID: number) => {
    Core.showSideLayer = layerID;
  },
  hideSideLayer: () => {
    Core.showSideLayer = 0;
  }
});

export const WeakSearchRef: WeakMap<any, {
  idToRecord: Map<string | number, any>;
  valueToRecord: Map<string, any>;
}> = new WeakMap();

export interface IFetchEvent {
  url: string;
}

export const fetchOnCourse = $state<Map<number, IFetchEvent>>(new SvelteMap());

export const fetchEvent = (fetchID: number, props: IFetchEvent | 0) => {
  // Genix HTTP helpers use ID 0 as a request for the next tracked fetch ID.
  if (fetchID === 0) {
    Env.fetchID += 1;
    return Env.fetchID;
  }

  if (props === 0) fetchOnCourse.delete(fetchID);
  else fetchOnCourse.set(fetchID, props);

  return fetchID;
};

export const openModals = $state<number[]>([]);

export const openModal = (id: number) => {
  if (!openModals.includes(id)) openModals.push(id);
};

export const closeModal = (id: number) => {
  const modalIndex = openModals.indexOf(id);
  if (modalIndex > -1) openModals.splice(modalIndex, 1);
};

export const closeAllModals = () => {
  openModals.length = 0;
};

export const imagesToUpload = new Map<number, () => Promise<IImageResult>>();
