/**
 * 管理员评论接口
 */

import { adminRequest } from '../http';
import type { Comment, PageResult } from '../types';

export interface AdminCommentListQuery {
  page?: number;
  page_size?: number;
  keyword?: string | null;
  article_slug?: string | null;
  user_uuid?: string | null;
}

export async function adminListComments(query?: AdminCommentListQuery): Promise<PageResult<Comment>> {
  const params = new URLSearchParams();
  if (query?.page) params.set('page', String(query.page));
  if (query?.page_size) params.set('page_size', String(query.page_size));
  if (query?.keyword) params.set('keyword', query.keyword);
  if (query?.article_slug) params.set('filter.article_slug', query.article_slug);
  if (query?.user_uuid) params.set('filter.user_uuid', query.user_uuid);
  const qs = params.toString();
  return adminRequest<PageResult<Comment>>(`/comment/list${qs ? `?${qs}` : ''}`);
}

export async function adminDeleteComment(id: number): Promise<void> {
  return adminRequest<void>(`/comment/delete/${id}`, {
    method: 'DELETE',
  });
}
