/**
 * 管理员反馈接口
 */

import { adminRequest } from '../http';
import type { Feedback, PageResult } from '../types';

export interface FeedbackReplyRequest {
  id: number;
  reply: string;
}

export interface FeedbackDeleteRequest {
  ids: number[];
}

export interface AdminFeedbackListQuery {
  page?: number;
  page_size?: number;
  status?: string | null;
}

export async function adminListFeedbacks(query?: AdminFeedbackListQuery): Promise<PageResult<Feedback>> {
  const params = new URLSearchParams();
  if (query?.page) params.set('page', String(query.page));
  if (query?.page_size) params.set('page_size', String(query.page_size));
  if (query?.status) params.set('filter.status', query.status);
  const qs = params.toString();
  return adminRequest<PageResult<Feedback>>(`/feedback/list${qs ? `?${qs}` : ''}`);
}

export async function adminReplyFeedback(data: FeedbackReplyRequest): Promise<void> {
  return adminRequest<void>('/feedback/reply', {
    method: 'PUT',
    body: JSON.stringify(data),
  });
}

export async function adminDeleteFeedbacks(data: FeedbackDeleteRequest): Promise<void> {
  return adminRequest<void>('/feedback/delete', {
    method: 'DELETE',
    body: JSON.stringify(data),
  });
}
