/**
 * 管理员用户接口
 */

import { adminRequest } from '../http';
import type { User, PageResult } from '../types';

export interface AdminUserListQuery {
  page?: number;
  page_size?: number;
  keyword?: string | null;
}

export interface LoginRecord {
  id: number;
  user_uuid: string;
  user: User;
  login_method: string;
  ip: string;
  address: string;
  os: string;
  device_info: string;
  browser_info: string;
  status: string;
  created_at: string;
  updated_at: string;
}

export async function adminListUsers(query?: AdminUserListQuery): Promise<PageResult<User>> {
  const params = new URLSearchParams();
  if (query?.page) params.set('page', String(query.page));
  if (query?.page_size) params.set('page_size', String(query.page_size));
  if (query?.keyword) params.set('keyword', query.keyword);
  const qs = params.toString();
  return adminRequest<PageResult<User>>(`/user/list${qs ? `?${qs}` : ''}`);
}

export async function adminFreezeUser(id: number): Promise<void> {
  return adminRequest<void>('/user/freeze', {
    method: 'PUT',
    body: JSON.stringify({ id }),
  });
}

export async function adminUnfreezeUser(id: number): Promise<void> {
  return adminRequest<void>('/user/unfreeze', {
    method: 'PUT',
    body: JSON.stringify({ id }),
  });
}

export async function adminListLoginRecords(query?: { page?: number; page_size?: number; uuid?: string | null }): Promise<PageResult<LoginRecord>> {
  const params = new URLSearchParams();
  if (query?.page) params.set('page', String(query.page));
  if (query?.page_size) params.set('page_size', String(query.page_size));
  if (query?.uuid) params.set('uuid', query.uuid);
  const qs = params.toString();
  return adminRequest<PageResult<LoginRecord>>(`/user/loginList${qs ? `?${qs}` : ''}`);
}
