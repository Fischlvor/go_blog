/**
 * 用户文章接口（需要认证）
 */

import { authRequest } from '../http';
import type { ArticleLikeInfo, Article, PageResult } from '../types';

/**
 * 点赞/取消点赞文章
 * POST /api/v1/article/:slug/like
 */
export async function toggleArticleLike(slug: string): Promise<ArticleLikeInfo> {
  return authRequest<ArticleLikeInfo>(`/article/${slug}/like`, {
    method: 'POST',
  });
}

/**
 * 取消点赞
 * DELETE /api/v1/article/:slug/like
 */
export async function removeArticleLike(slug: string): Promise<ArticleLikeInfo> {
  return authRequest<ArticleLikeInfo>(`/article/${slug}/like`, {
    method: 'DELETE',
  });
}

/**
 * 获取用户点赞列表
 * GET /api/v1/article/likes
 */
export async function getUserLikedArticles(page = 1, pageSize = 10): Promise<PageResult<Article>> {
  const params = new URLSearchParams();
  params.set('page', String(page));
  params.set('page_size', String(pageSize));
  return authRequest<PageResult<Article>>(`/article/likes?${params.toString()}`);
}
