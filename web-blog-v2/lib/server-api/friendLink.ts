import 'server-only';

import { serverRequest } from './shared';
import { FALLBACK_FRIEND_LINKS } from './fallback';
import type { FriendLink } from '@/lib/client-api/types';

export async function getFriendLinksServer(): Promise<FriendLink[]> {
  try {
    const data = await serverRequest<FriendLink[]>('/friendLink/info', {
      next: { revalidate: 300 },
    });
    return data;
  } catch (error) {
    console.error('[server-api] friend links request failed', error);
    return FALLBACK_FRIEND_LINKS;
  }
}
