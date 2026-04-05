import { HomePageClient } from '@/components/site/home/HomePageClient';
import { listArticlesServer } from '@/lib/server-api/article';
import { getWebsiteInfoServer } from '@/lib/server-api/website';

export default async function HomePage() {
  const [site, articleResult] = await Promise.all([
    getWebsiteInfoServer(),
    listArticlesServer({ page: 1, page_size: 6, order: 'desc' }),
  ]);

  return <HomePageClient site={site} articles={articleResult.list} />;
}
