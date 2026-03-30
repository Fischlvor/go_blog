/**
 * 管理员图片接口
 */

import { adminRequest } from '../http';
import type { PageResult } from '../types';

export interface AdminImage {
  id: number;
  name: string;
  url: string;
  size: number;
  mime_type: string;
  created_at: string;
  updated_at: string;
}

export interface AdminImageListQuery {
  page?: number;
  page_size?: number;
  name?: string | null;
  mime_type?: string | null;
}

export async function adminListImages(query?: AdminImageListQuery): Promise<PageResult<AdminImage>> {
  const params = new URLSearchParams();
  if (query?.page) params.set('page', String(query.page));
  if (query?.page_size) params.set('page_size', String(query.page_size));
  if (query?.name) params.set('filter.name', query.name);
  if (query?.mime_type) params.set('filter.mime_type', query.mime_type);
  const qs = params.toString();
  return adminRequest<PageResult<AdminImage>>(`/image/list${qs ? `?${qs}` : ''}`);
}

export async function adminDeleteImages(ids: number[]): Promise<void> {
  return adminRequest<void>('/image/delete', {
    method: 'DELETE',
    body: JSON.stringify({ ids }),
  });
}
