import 'server-only';

import { serverRequest } from './shared';
import { FALLBACK_ARTICLE_DETAIL, FALLBACK_ARTICLE_PAGE, FALLBACK_CATEGORIES } from './fallback';
import type { Article, ArticleCategory, ArticleDetail, ArticleListQuery, PageResult } from '@/lib/client-api/types';

function buildArticleQuery(query?: ArticleListQuery): string {
  const params = new URLSearchParams();
  if (query?.page) params.set('page', String(query.page));
  if (query?.page_size) params.set('page_size', String(query.page_size));
  if (query?.keyword) params.set('keyword', query.keyword);
  if (query?.['filter.category_id']) params.set('filter.category_id', String(query['filter.category_id']));
  if (query?.['filter.tag_id']) params.set('filter.tag_id', String(query['filter.tag_id']));
  if (query?.sort_by) params.set('sort_by', query.sort_by);
  if (query?.order) params.set('order', query.order);
  const qs = params.toString();
  return qs ? `?${qs}` : '';
}

export async function listArticlesServer(query?: ArticleListQuery): Promise<PageResult<Article>> {
  try {
    return await serverRequest<PageResult<Article>>(`/article/search${buildArticleQuery(query)}`, {
      next: { revalidate: 60 },
    });
  } catch {
    return {
      ...FALLBACK_ARTICLE_PAGE,
      current_page: query?.page || 1,
      page_size: query?.page_size || 0,
    };
  }
}

export async function getArticleServer(slug: string): Promise<ArticleDetail> {
  try {
    return await serverRequest<ArticleDetail>(`/article/${slug}`, {
      next: { revalidate: 60 },
    });
  } catch {
    return {
      ...FALLBACK_ARTICLE_DETAIL,
      slug,
    };
  }
}

export async function listCategoriesServer(): Promise<ArticleCategory[]> {
  try {
    return await serverRequest<ArticleCategory[]>('/article/category', {
      next: { revalidate: 300 },
    });
  } catch {
    return FALLBACK_CATEGORIES;
  }
}
