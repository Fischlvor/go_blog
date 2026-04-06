import 'server-only';

import { serverRequest } from './shared';
import { FALLBACK_SITE } from './fallback';
import type { Website } from '@/lib/client-api/public/website';

export async function getWebsiteInfoServer(): Promise<Website> {
  try {
    return await serverRequest<Website>('/website/info', {
      next: { revalidate: 300 },
    });
  } catch (error) {
    console.error('[SSR] getWebsiteInfoServer failed:', error);
    return FALLBACK_SITE;
  }
}
