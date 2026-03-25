/**
 * 公开评论接口
 */

import { publicRequest, authRequest } from '../http';
import type { Comment } from '../types';

/**
 * 获取文章评论列表（分页）
 * GET /api/v1/comment/:article_slug
 */
export async function getArticleComments(articleSlug: string, page = 1, pageSize = 20): Promise<Comment[]> {
  const params = new URLSearchParams();
  params.set('page', String(page));
  params.set('page_size', String(pageSize));
  const result = await publicRequest<{ list: Comment[]; current_page: number; page_size: number; total_items: number; total_pages: number }>(`/comment/${articleSlug}?${params.toString()}`);
  return result.list ?? [];
}

/**
 * 获取最新评论
 */
export async function getNewComments(): Promise<Comment[]> {
  return publicRequest<Comment[]>('/comment/new');
}

export interface CommentCreateRequest {
  article_slug: string;
  parent_id: number | null;
  content: string;
}

/**
 * 发表评论（需要登录）
 */
export async function createComment(data: CommentCreateRequest): Promise<{ id: number }> {
  return authRequest<{ id: number }>('/comment/create', {
    method: 'POST',
    body: JSON.stringify(data),
  });
}

/**
 * 删除评论（需要登录）
 */
export async function deleteComment(ids: number[]): Promise<void> {
  return authRequest<void>('/comment/delete', {
    method: 'DELETE',
    body: JSON.stringify({ ids }),
  });
}
