/**
 * 管理员 AI 管理接口
 */

import { adminRequest } from '../http';
import type { PageResult } from '../types';

export interface AIModel {
  id: number;
  name: string;
  display_name?: string;
  provider?: string;
  endpoint?: string;
  api_key?: string;
  model?: string;
  max_tokens?: number;
  temperature?: number;
  status?: string;
  is_active?: boolean;
  created_at: string;
  updated_at: string;
}

export interface AISession {
  id: number;
  title?: string;
  user_uuid?: string;
  user?: {
    uuid: string;
    username: string;
    avatar: string;
  };
  model?: string;
  message_count?: number;
  created_at: string;
  updated_at: string;
}

export interface AIMessage {
  id: number;
  session_id: number;
  role: string;
  content: string;
  tokens?: number;
  created_at: string;
  updated_at: string;
}

export interface AdminAIModelQuery {
  page?: number;
  page_size?: number;
  name?: string;
  provider?: string;
}

export function adminListAIModels(query?: AdminAIModelQuery): Promise<PageResult<AIModel>> {
  const params = new URLSearchParams();
  if (query?.page) params.set('page', String(query.page));
  if (query?.page_size) params.set('page_size', String(query.page_size));
  if (query?.name) params.set('filter.name', query.name);
  if (query?.provider) params.set('filter.provider', query.provider);
  const qs = params.toString();
  return adminRequest<PageResult<AIModel>>(`/ai-management/model${qs ? `?${qs}` : ''}`);
}

export function adminGetAIModel(id: number): Promise<AIModel> {
  return adminRequest<AIModel>(`/ai-management/model/${id}`);
}

export function adminCreateAIModel(data: Omit<AIModel, 'id' | 'created_at' | 'updated_at'>): Promise<void> {
  return adminRequest<void>('/ai-management/model', {
    method: 'POST',
    body: JSON.stringify(data),
  });
}

export function adminUpdateAIModel(data: Omit<AIModel, 'created_at' | 'updated_at'>): Promise<void> {
  return adminRequest<void>('/ai-management/model', {
    method: 'PUT',
    body: JSON.stringify(data),
  });
}

export function adminDeleteAIModel(id: number): Promise<void> {
  return adminRequest<void>(`/ai-management/model/${id}`, {
    method: 'DELETE',
  });
}

export interface AdminAISessionQuery {
  page?: number;
  page_size?: number;
  keyword?: string;
  user_uuid?: string;
}

export function adminListAISessions(query?: AdminAISessionQuery): Promise<PageResult<AISession>> {
  const params = new URLSearchParams();
  if (query?.page) params.set('page', String(query.page));
  if (query?.page_size) params.set('page_size', String(query.page_size));
  if (query?.keyword) params.set('keyword', query.keyword);
  if (query?.user_uuid) params.set('filter.user_uuid', query.user_uuid);
  const qs = params.toString();
  return adminRequest<PageResult<AISession>>(`/ai-management/session${qs ? `?${qs}` : ''}`);
}

export function adminGetAISession(id: number): Promise<AISession & { messages?: AIMessage[] }> {
  return adminRequest<AISession & { messages?: AIMessage[] }>(`/ai-management/session/${id}`);
}

export function adminDeleteAISession(id: number): Promise<void> {
  return adminRequest<void>(`/ai-management/session/${id}`, {
    method: 'DELETE',
  });
}

export interface AdminAIMessageQuery {
  page?: number;
  page_size?: number;
  session_id?: number;
  role?: string;
}

export function adminListAIMessages(query?: AdminAIMessageQuery): Promise<PageResult<AIMessage>> {
  const params = new URLSearchParams();
  if (query?.page) params.set('page', String(query.page));
  if (query?.page_size) params.set('page_size', String(query.page_size));
  if (query?.session_id) params.set('filter.session_id', String(query.session_id));
  if (query?.role) params.set('filter.role', query.role);
  const qs = params.toString();
  return adminRequest<PageResult<AIMessage>>(`/ai-management/message${qs ? `?${qs}` : ''}`);
}

export function adminGetAIMessage(id: number): Promise<AIMessage> {
  return adminRequest<AIMessage>(`/ai-management/message/${id}`);
}

export function adminDeleteAIMessage(id: number): Promise<void> {
  return adminRequest<void>(`/ai-management/message/${id}`, {
    method: 'DELETE',
  });
}
