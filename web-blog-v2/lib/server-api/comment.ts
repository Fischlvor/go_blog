import 'server-only';

import { serverRequest } from './shared';
import { FALLBACK_COMMENTS } from './fallback';
import type { Comment, PageResult } from '@/lib/client-api/types';

export async function getArticleCommentsServer(articleSlug: string, page = 1, pageSize = 20): Promise<Comment[]> {
  const params = new URLSearchParams();
  params.set('page', String(page));
  params.set('page_size', String(pageSize));

  try {
    const result = await serverRequest<PageResult<Comment>>(`/comment/${articleSlug}?${params.toString()}`, {
      next: { revalidate: 60 },
    });

    return result.list ?? [];
  } catch {
    return FALLBACK_COMMENTS;
  }
}
