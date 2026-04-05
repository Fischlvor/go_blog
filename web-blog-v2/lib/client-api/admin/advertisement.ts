/**
 * 管理员广告接口
 */

import { adminRequest } from '../http';

export interface AdminAdvertisementCreateRequest {
  ad_image: string;
  link: string;
  title: string;
  content: string;
}

export interface AdminAdvertisementUpdateRequest extends AdminAdvertisementCreateRequest {
  id: number;
}

export interface AdminAdvertisementDeleteRequest {
  ids: number[];
}

export interface Advertisement {
  id: number;
  ad_image: string;
  link: string;
  title: string;
  content: string;
  created_at: string;
  updated_at: string;
}

export interface AdvertisementListResult {
  list: Advertisement[];
  current_page: number;
  page_size: number;
  total_items: number;
  total_pages: number;
}

export interface AdminAdvertisementListQuery {
  page?: number;
  page_size?: number;
  title?: string | null;
  content?: string | null;
}

export async function adminListAdvertisements(query?: AdminAdvertisementListQuery): Promise<AdvertisementListResult> {
  const params = new URLSearchParams();
  if (query?.page) params.set('page', String(query.page));
  if (query?.page_size) params.set('page_size', String(query.page_size));
  if (query?.title) params.set('title', query.title);
  if (query?.content) params.set('content', query.content);
  const qs = params.toString();
  return adminRequest<AdvertisementListResult>(`/advertisement/list${qs ? `?${qs}` : ''}`);
}

export async function adminCreateAdvertisement(data: AdminAdvertisementCreateRequest): Promise<void> {
  return adminRequest<void>('/advertisement/create', {
    method: 'POST',
    body: JSON.stringify(data),
  });
}

export async function adminUpdateAdvertisement(data: AdminAdvertisementUpdateRequest): Promise<void> {
  return adminRequest<void>('/advertisement/update', {
    method: 'PUT',
    body: JSON.stringify(data),
  });
}

export async function adminDeleteAdvertisements(data: AdminAdvertisementDeleteRequest): Promise<void> {
  return adminRequest<void>('/advertisement/delete', {
    method: 'DELETE',
    body: JSON.stringify(data),
  });
}
