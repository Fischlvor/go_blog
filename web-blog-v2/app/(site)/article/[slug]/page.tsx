import type { Metadata } from 'next';
import { notFound } from 'next/navigation';
import { ArrowLeft, Calendar, Eye, Hash } from 'lucide-react';
import { Badge } from '@/components/ui/badge';
import { Avatar, AvatarFallback, AvatarImage } from '@/components/ui/avatar';
import { MarkdownContent } from '@/components/site/article/MarkdownContent';
import { ArticleInteractions } from '@/components/site/article/ArticleInteractions';
import { TableOfContents } from '@/components/site/article/TableOfContents';
import { getArticleServer } from '@/lib/server-api/article';
import { getArticleCommentsServer } from '@/lib/server-api/comment';
import { getWebsiteInfoServer } from '@/lib/server-api/website';
import { getFieldValue } from '@/lib/client-api/public/website';
import { FALLBACK_ARTICLE_DETAIL } from '@/lib/server-api/fallback';

function formatDate(d: string | null | undefined) {
  if (!d) return '';
  return new Date(d).toLocaleDateString('zh-CN', { year: 'numeric', month: 'long', day: 'numeric' });
}

interface TocItem {
  id: string;
  text: string;
  level: number;
}

function extractToc(markdown: string): TocItem[] {
  const items: TocItem[] = [];
  const lines = markdown.split('\n');
  for (const line of lines) {
    const match = line.match(/^(#{1,4})\s+(.+)/);
    if (!match) continue;
    const level = match[1].length;
    const text = match[2].trim();
    const id = text.toLowerCase().replace(/\s+/g, '-').replace(/[^\w\u4e00-\u9fa5-]/g, '');
    items.push({ id, text, level });
  }
  return items;
}

interface ArticleDetailPageProps {
  params: Promise<{ slug: string }>;
}

export async function generateMetadata({ params }: ArticleDetailPageProps): Promise<Metadata> {
  try {
    const [{ slug }, site] = await Promise.all([params, getWebsiteInfoServer()]);
    const article = await getArticleServer(slug);

    if (!article.id) {
      return {
        title: '文章不存在',
      };
    }

    const title = article.meta_title || article.title;
    const description = article.meta_description || article.excerpt || getFieldValue(site.description) || '';

    return {
      title,
      description,
      openGraph: {
        title,
        description,
        images: article.featured_image ? [article.featured_image] : undefined,
      },
    };
  } catch {
    return {
      title: '文章不存在',
    };
  }
}

export default async function ArticleDetailPage({ params }: ArticleDetailPageProps) {
  const { slug } = await params;

  try {
    const [article, comments] = await Promise.all([
      getArticleServer(slug),
      getArticleCommentsServer(slug),
    ]);

    if (!article.id || article.slug === FALLBACK_ARTICLE_DETAIL.slug) {
      notFound();
    }

    const toc = extractToc(article.content || '');

    return (
      <div className="max-w-6xl mx-auto px-4 py-10">
        <a
          href="/articles"
          className="mb-8 inline-flex items-center gap-1.5 text-sm text-muted-foreground hover:text-foreground transition-colors cursor-pointer"
        >
          <ArrowLeft className="h-4 w-4" /> 返回文章列表
        </a>

        <div className="flex gap-12 items-start mb-8">
          <header className="max-w-3xl flex-1 min-w-0">
            <div className="flex flex-wrap gap-2 mb-4">
              {article.category && (
                <Badge className="rounded-full bg-primary/10 text-primary border-primary/20">{article.category.name}</Badge>
              )}
              {article.is_featured && (
                <Badge variant="secondary" className="rounded-full">
                  精选
                </Badge>
              )}
            </div>
            <h1 className="text-3xl md:text-4xl font-bold tracking-tight leading-tight mb-4">{article.title}</h1>
            {article.excerpt && <p className="text-base text-muted-foreground leading-relaxed mb-5">{article.excerpt}</p>}
            <div className="flex items-center justify-between gap-x-5 gap-y-2 text-sm text-muted-foreground pb-5 border-b border-border flex-wrap">
              <div className="flex items-center gap-2">
                <Avatar className="h-7 w-7">
                  <AvatarImage src={article.author?.avatar} />
                  <AvatarFallback className="text-xs">{article.author?.nickname?.[0]}</AvatarFallback>
                </Avatar>
                <span className="font-medium text-foreground">{article.author?.nickname}</span>
              </div>
              <div className="flex items-center gap-4">
                <span className="flex items-center gap-1">
                  <Calendar className="h-3.5 w-3.5" />
                  {formatDate(article.published_at)}
                </span>
                <span className="flex items-center gap-1">
                  <Eye className="h-3.5 w-3.5" />
                  {article.views} 次阅读
                </span>
              </div>
            </div>
          </header>

          {article.featured_image && (
            <div className="hidden lg:block w-56 xl:w-64 flex-shrink-0">
              <img
                src={article.featured_image}
                alt={article.title}
                className="w-full h-auto max-h-64 object-cover rounded-xl ring-1 ring-border"
              />
            </div>
          )}
        </div>

        <div className="flex gap-12 items-start">
          <main className="min-w-0 max-w-3xl flex-1">
            <MarkdownContent content={article.content || ''} />

            {article.tags?.length > 0 && (
              <div className="flex flex-wrap gap-2 mt-10 pt-6 border-t border-border">
                {article.tags.map((tag) => (
                  <span
                    key={tag.id}
                    className="inline-flex items-center gap-1 rounded-full border px-2.5 py-0.5 text-xs font-mono text-muted-foreground cursor-default"
                  >
                    <Hash className="h-2.5 w-2.5" />
                    {tag.name}
                  </span>
                ))}
              </div>
            )}

            <ArticleInteractions
              slug={slug}
              initialComments={comments}
              initialLiked={false}
              initialLikes={article.like?.likes ?? 0}
            />
          </main>

          <TableOfContents toc={toc} />
        </div>
      </div>
    );
  } catch {
    notFound();
  }
}
