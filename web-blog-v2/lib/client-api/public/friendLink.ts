/**
 * 友链公开接口
 */

import { publicRequest } from '../http';
import type { FriendLink } from '../types';

export async function getFriendLinks(): Promise<FriendLink[]> {
  return publicRequest<FriendLink[]>('/friendLink/info');
}
