export const CDN_BASE: string = (import.meta as any).env?.VITE_CDN_BASE || 'https://image.hsk423.cn';

export function cdn(path: string): string {
  const normalized = path.startsWith('/') ? path.slice(1) : path;
  return `${CDN_BASE}/${normalized}`;
}


