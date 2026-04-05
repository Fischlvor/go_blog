/**
 * 管理员资源接口
 */

import { adminRequest } from '../http';
import type { PageResult } from '../types';

export interface AdminResource {
  id: number;
  file_key: string;
  file_name: string;
  file_url: string;
  file_size: number;
  mime_type: string;
  transcode_status: number;
  transcode_url?: string;
  thumbnail_url?: string;
  created_at: string;
  updated_at: string;
}

export interface AdminResourceListQuery {
  page?: number;
  page_size?: number;
  file_name?: string | null;
  mime_type?: string | null;
}

export interface ResourceCheckRequest {
  file_hash: string;
  file_size: number;
  file_name: string;
}

export interface ResourceCheckResponse {
  exists: boolean;
  file_url?: string;
  task_id?: string;
  total_chunks?: number;
  uploaded_chunks?: number[];
  missing_chunks?: number[];
}

export interface ResourceInitRequest {
  file_hash: string;
  file_size: number;
  file_name: string;
  mime_type: string;
}

export interface ResourceInitResponse {
  task_id: string;
  total_chunks: number;
  chunk_size: number;
}

export interface ResourceUploadChunkResponse {
  success: boolean;
  chunk_number: number;
}

export interface ResourceCompleteResponse {
  file_url: string;
  file_key: string;
}

export async function adminListResources(query?: AdminResourceListQuery): Promise<PageResult<AdminResource>> {
  const params = new URLSearchParams();
  if (query?.page) params.set('page', String(query.page));
  if (query?.page_size) params.set('page_size', String(query.page_size));
  if (query?.file_name) params.set('filter.file_name', query.file_name);
  if (query?.mime_type) params.set('filter.mime_type', query.mime_type);
  const qs = params.toString();
  return adminRequest<PageResult<AdminResource>>(`/resources/list${qs ? `?${qs}` : ''}`);
}

export async function adminDeleteResources(ids: number[]): Promise<void> {
  return adminRequest<void>('/resources/delete', {
    method: 'POST',
    body: JSON.stringify({ ids }),
  });
}

export async function adminGetResourceMaxSize(): Promise<{ max_size: number }> {
  return adminRequest<{ max_size: number }>('/resources/max-size');
}

export async function adminCheckResource(data: ResourceCheckRequest): Promise<ResourceCheckResponse> {
  return adminRequest<ResourceCheckResponse>('/resources/check', {
    method: 'POST',
    body: JSON.stringify(data),
  });
}

export async function adminInitResource(data: ResourceInitRequest): Promise<ResourceInitResponse> {
  return adminRequest<ResourceInitResponse>('/resources/init', {
    method: 'POST',
    body: JSON.stringify(data),
  });
}

export async function adminUploadResourceChunk(
  taskId: string,
  chunkNumber: number,
  chunkData: Blob,
  signal?: AbortSignal
): Promise<ResourceUploadChunkResponse> {
  const formData = new FormData();
  formData.append('task_id', taskId);
  formData.append('chunk_number', String(chunkNumber));
  formData.append('chunk_data', chunkData);

  return adminRequest<ResourceUploadChunkResponse>('/resources/upload-chunk', {
    method: 'POST',
    body: formData,
    signal,
  });
}

export async function adminCompleteResource(taskId: string): Promise<ResourceCompleteResponse> {
  return adminRequest<ResourceCompleteResponse>('/resources/complete', {
    method: 'POST',
    body: JSON.stringify({ task_id: taskId }),
  });
}

export async function adminCancelResource(taskId: string): Promise<void> {
  return adminRequest<void>('/resources/cancel', {
    method: 'POST',
    body: JSON.stringify({ task_id: taskId }),
  });
}
