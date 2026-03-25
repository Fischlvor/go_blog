/**
 * 公开文章接口（无需认证）
 * 路由：/api/v1/article/*
 */

import { publicRequest } from '../http';
import type { Article, ArticleDetail, ArticleCategory, ArticleTag, ArticleListQuery, PageResult } from '../types';

/**
 * 获取文章列表
 * GET /api/v1/article/search
 */
export async function listArticles(query?: ArticleListQuery): Promise<PageResult<Article>> {
  const params = new URLSearchParams();
  if (query?.page) params.set('page', String(query.page));
  if (query?.page_size) params.set('page_size', String(query.page_size));
  if (query?.keyword) params.set('keyword', query.keyword);
  if (query?.['filter.category_id']) params.set('filter.category_id', String(query['filter.category_id']));
  if (query?.['filter.tag_id']) params.set('filter.tag_id', String(query['filter.tag_id']));
  if (query?.sort_by) params.set('sort_by', query.sort_by);
  if (query?.order) params.set('order', query.order);
  const qs = params.toString();
  return publicRequest<PageResult<Article>>(`/article/search${qs ? `?${qs}` : ''}`);
}

/**
 * 获取文章详情（按 slug）
 * GET /api/v1/article/:slug
 */
export async function getArticle(slug: string): Promise<ArticleDetail> {
  return publicRequest<ArticleDetail>(`/article/${slug}`);
}

/**
 * 获取所有分类
 * GET /api/v1/article/category
 */
export async function listCategories(): Promise<ArticleCategory[]> {
  return publicRequest<ArticleCategory[]>('/article/category');
}

/**
 * 获取所有标签
 * GET /api/v1/article/tags
 */
export async function listTags(): Promise<ArticleTag[]> {
  return publicRequest<ArticleTag[]>('/article/tags');
}
