/**
 * 管理员友链接口
 */

import { adminRequest } from '../http';
import type { FriendLink } from '../types';

export interface AdminFriendLinkCreateRequest {
  logo: string;
  link: string;
  name: string;
  description: string;
  sort?: number;
}

export interface AdminFriendLinkUpdateRequest extends AdminFriendLinkCreateRequest {
  id: number;
}

export interface AdminFriendLinkDeleteRequest {
  ids: number[];
}

export interface AdminFriendLinkListQuery {
  page?: number;
  page_size?: number;
  name?: string | null;
  description?: string | null;
}

export interface FriendLinkListResult {
  list: FriendLink[];
  current_page: number;
  page_size: number;
  total_items: number;
  total_pages: number;
}

export async function adminListFriendLinks(query?: AdminFriendLinkListQuery): Promise<FriendLinkListResult> {
  const params = new URLSearchParams();
  if (query?.page) params.set('page', String(query.page));
  if (query?.page_size) params.set('page_size', String(query.page_size));

  const list = await adminRequest<FriendLink[]>('/friendLink/list');
  const keywordName = query?.name?.trim();
  const keywordDesc = query?.description?.trim();

  const filtered = list.filter((item) => {
    const matchedName = keywordName ? item.name?.includes(keywordName) : true;
    const matchedDesc = keywordDesc ? item.description?.includes(keywordDesc) : true;
    return matchedName && matchedDesc;
  });

  const page = query?.page || 1;
  const pageSize = query?.page_size || 10;
  const start = (page - 1) * pageSize;
  const end = start + pageSize;

  return {
    list: filtered.slice(start, end),
    current_page: page,
    page_size: pageSize,
    total_items: filtered.length,
    total_pages: Math.max(1, Math.ceil(filtered.length / pageSize)),
  };
}

export async function adminCreateFriendLink(data: AdminFriendLinkCreateRequest): Promise<void> {
  return adminRequest<void>('/friendLink/create', {
    method: 'POST',
    body: JSON.stringify(data),
  });
}

export async function adminUpdateFriendLink(data: AdminFriendLinkUpdateRequest): Promise<void> {
  return adminRequest<void>('/friendLink/update', {
    method: 'PUT',
    body: JSON.stringify(data),
  });
}

export async function adminDeleteFriendLinks(data: AdminFriendLinkDeleteRequest): Promise<void> {
  return adminRequest<void>('/friendLink/delete', {
    method: 'DELETE',
    body: JSON.stringify(data),
  });
}
