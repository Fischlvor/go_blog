/**
 * 管理员文章接口
 */

import { adminRequest } from '../http';
import type { Article, PageResult } from '../types';

export interface AdminArticleListQuery {
  page?: number;
  page_size?: number;
  title?: string | null;
  category?: string | null;
  abstract?: string | null;
}

export async function adminListArticles(query?: AdminArticleListQuery): Promise<PageResult<Article>> {
  const params = new URLSearchParams();
  if (query?.page) params.set('page', String(query.page));
  if (query?.page_size) params.set('page_size', String(query.page_size));
  if (query?.title) params.set('title', query.title);
  if (query?.category) params.set('category', query.category);
  if (query?.abstract) params.set('abstract', query.abstract);
  const qs = params.toString();
  return adminRequest<PageResult<Article>>(`/admin/article/list${qs ? `?${qs}` : ''}`);
}

export interface ArticleCreateRequest {
  cover: string;
  title: string;
  category: string;
  tags: string[];
  abstract: string;
  content: string;
}

export async function adminCreateArticle(data: ArticleCreateRequest): Promise<void> {
  return adminRequest<void>('/admin/article/create', {
    method: 'POST',
    body: JSON.stringify(data),
  });
}

export interface ArticleUpdateRequest extends ArticleCreateRequest {
  id: string;
}

export async function adminUpdateArticle(data: ArticleUpdateRequest): Promise<void> {
  return adminRequest<void>('/admin/article/update', {
    method: 'PUT',
    body: JSON.stringify(data),
  });
}

export async function adminDeleteArticles(ids: string[]): Promise<void> {
  return adminRequest<void>('/admin/article/delete', {
    method: 'DELETE',
    body: JSON.stringify({ ids }),
  });
}
