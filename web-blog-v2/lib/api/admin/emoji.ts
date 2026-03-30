/**
 * 管理员表情接口
 */

import { adminRequest, authRequest, publicRequest } from '../http';
import type { PageResult } from '../types';

export interface AdminEmoji {
  id: number;
  key: string;
  filename: string;
  group_name: string;
  group_key?: string;
  sprite_group: number;
  sprite_position_x: number;
  sprite_position_y: number;
  file_size: number;
  cdn_url: string;
  upload_time: string;
  status: number;
  created_by: string;
  created_at: string;
  updated_at: string;
}

export interface AdminEmojiListQuery {
  page?: number;
  page_size?: number;
  keyword?: string;
  group_key?: string;
  sprite_group?: number;
}

export async function adminListEmojis(query?: AdminEmojiListQuery): Promise<PageResult<AdminEmoji>> {
  const params = new URLSearchParams();
  if (query?.page) params.set('page', String(query.page));
  if (query?.page_size) params.set('page_size', String(query.page_size));
  if (query?.keyword) params.set('keyword', query.keyword);
  if (query?.group_key) params.set('group_key', query.group_key);
  if (query?.sprite_group !== undefined) params.set('sprite_group', String(query.sprite_group));
  const qs = params.toString();
  return adminRequest<PageResult<AdminEmoji>>(`/emoji/list${qs ? `?${qs}` : ''}`);
}

export interface AdminEmojiGroup {
  id: number;
  group_name: string;
  group_key: string;
  description: string;
  sort_order: number;
  emoji_count: number;
  status: number;
  created_by: string;
  created_at: string;
  updated_at: string;
}

export interface AdminEmojiGroupCreateRequest {
  group_name: string;
  description?: string;
  sort_order?: number;
}

export interface AdminEmojiGroupUpdateRequest extends AdminEmojiGroupCreateRequest {
  id: number;
}

export function adminListEmojiGroups(): Promise<AdminEmojiGroup[]> {
  return publicRequest<AdminEmojiGroup[]>('/emoji/groups');
}

export function adminCreateEmojiGroup(data: AdminEmojiGroupCreateRequest): Promise<void> {
  return adminRequest<void>('/emoji/groups', {
    method: 'POST',
    body: JSON.stringify(data),
  });
}

export function adminUpdateEmojiGroup(id: number, data: AdminEmojiGroupCreateRequest): Promise<void> {
  return adminRequest<void>(`/emoji/groups/${id}`, {
    method: 'PUT',
    body: JSON.stringify(data),
  });
}

export function adminDeleteEmojiGroup(id: number): Promise<void> {
  return adminRequest<void>(`/emoji/groups/${id}`, {
    method: 'DELETE',
  });
}

export interface AdminEmojiSprite {
  id: number;
  sprite_group: number;
  filename: string;
  cdn_url: string;
  width: number;
  height: number;
  emoji_count: number;
  file_size: number;
  status: number;
  created_at: string;
  updated_at: string;
}

export function adminListEmojiSprites(): Promise<AdminEmojiSprite[]> {
  return adminRequest<AdminEmojiSprite[]>('/emoji/sprites');
}

export function adminDeleteEmoji(id: number): Promise<void> {
  return adminRequest<void>(`/emoji/${id}`, { method: 'DELETE' });
}

export function adminRestoreEmoji(id: number): Promise<void> {
  return adminRequest<void>(`/emoji/${id}/restore`, { method: 'PUT' });
}

export function adminUploadEmojis(formData: FormData): Promise<void> {
  return authRequest<void>('/emoji/upload', {
    method: 'POST',
    body: formData,
    apiVersion: 'admin',
  });
}

export function adminRegenerateEmojiSprites(groupKeys?: string[]): Promise<void> {
  return adminRequest<void>('/emoji/regenerate', {
    method: 'POST',
    body: JSON.stringify({ group_keys: groupKeys || [] }),
  });
}
