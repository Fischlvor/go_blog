/**
 * 管理员文章接口
 */

import { adminRequest } from '../http';
import type { Article, PageResult } from '../types';

export interface AdminArticleListQuery {
  page?: number;
  page_size?: number;
  keyword?: string | null;
  category_id?: number | null;
  tag_id?: number | null;
  status?: string | null;
  visibility?: string | null;
}

export async function adminListArticles(query?: AdminArticleListQuery): Promise<PageResult<Article>> {
  const params = new URLSearchParams();
  if (query?.page) params.set('page', String(query.page));
  if (query?.page_size) params.set('page_size', String(query.page_size));
  if (query?.keyword) params.set('keyword', query.keyword);
  if (query?.category_id) params.set('filter.category_id', String(query.category_id));
  if (query?.tag_id) params.set('filter.tag_id', String(query.tag_id));
  if (query?.status) params.set('filter.status', query.status);
  if (query?.visibility) params.set('filter.visibility', query.visibility);
  const qs = params.toString();
  return adminRequest<PageResult<Article>>(`/article/list${qs ? `?${qs}` : ''}`);
}

export interface ArticleCreateRequest {
  title: string;
  slug?: string;
  content: string;
  excerpt?: string;
  featured_image?: string;
  category_id: number;
  tag_ids: number[];
  status: 'draft' | 'published';
  visibility?: 'public' | 'private';
  is_featured?: boolean;
}

export async function adminCreateArticle(data: ArticleCreateRequest): Promise<void> {
  return adminRequest<void>('/article/create', {
    method: 'POST',
    body: JSON.stringify(data),
  });
}

export interface ArticleUpdateRequest extends ArticleCreateRequest {
  slug: string;
}

export async function adminUpdateArticle(data: ArticleUpdateRequest): Promise<void> {
  return adminRequest<void>('/article/update', {
    method: 'PUT',
    body: JSON.stringify(data),
  });
}

export async function adminDeleteArticles(slugs: string[]): Promise<void> {
  return adminRequest<void>('/article/delete', {
    method: 'DELETE',
    body: JSON.stringify({ ids: slugs }),
  });
}

export interface ArticleSaveDraftRequest {
  title: string;
  content: string;
  excerpt?: string;
  featured_image?: string;
  category_id: number;
  tag_ids: number[];
}

export async function adminSaveDraft(data: ArticleSaveDraftRequest): Promise<{ slug: string }> {
  return adminRequest<{ slug: string }>('/article/create', {
    method: 'POST',
    body: JSON.stringify({
      ...data,
      status: 'draft',
      visibility: 'private',
    }),
  });
}

export interface CreateCategoryRequest {
  name: string;
  slug: string;
}

export async function adminCreateCategory(data: CreateCategoryRequest): Promise<{ id: number }> {
  return adminRequest<{ id: number }>('/category/create', {
    method: 'POST',
    body: JSON.stringify(data),
  });
}

export interface CreateTagRequest {
  name: string;
  slug: string;
}

export async function adminCreateTag(data: CreateTagRequest): Promise<{ id: number }> {
  return adminRequest<{ id: number }>('/tag/create', {
    method: 'POST',
    body: JSON.stringify(data),
  });
}
