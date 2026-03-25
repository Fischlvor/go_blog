'use client';

import { useEffect, useState } from 'react';
import { useParams, useRouter } from 'next/navigation';
import { ArrowLeft, Calendar, Eye, Hash } from 'lucide-react';
import { HeartIcon } from '@/components/common/icons';
import { Button } from '@/components/ui/button';
import { Badge } from '@/components/ui/badge';
import { Separator } from '@/components/ui/separator';
import { Skeleton } from '@/components/ui/skeleton';
import { Avatar, AvatarFallback, AvatarImage } from '@/components/ui/avatar';
import { toast } from 'sonner';
import { getArticle } from '@/lib/api/public/article';
import { toggleArticleLike, removeArticleLike } from '@/lib/api/user/article';
import { getArticleComments, createComment } from '@/lib/api/public/comment';
import { useUserAuth } from '@/context/user-auth';
import { MarkdownContent } from '@/components/site/article/MarkdownContent';
import type { ArticleDetail, Comment } from '@/lib/api/types';
import { cn } from '@/lib/utils';

function formatDate(d: string | null | undefined) {
  if (!d) return '';
  return new Date(d).toLocaleDateString('zh-CN', { year: 'numeric', month: 'long', day: 'numeric' });
}

interface TocItem { id: string; text: string; level: number; }

function extractToc(markdown: string): TocItem[] {
  const items: TocItem[] = [];
  const lines = markdown.split('\n');
  for (const line of lines) {
    const m = line.match(/^(#{1,4})\s+(.+)/);
    if (!m) continue;
    const level = m[1].length;
    const text = m[2].trim();
    const id = text.toLowerCase().replace(/\s+/g, '-').replace(/[^\w\u4e00-\u9fa5-]/g, '');
    items.push({ id, text, level });
  }
  return items;
}

function TableOfContents({ toc }: { toc: TocItem[] }) {
  const [active, setActive] = useState('');
  useEffect(() => {
    const obs = new IntersectionObserver(
      (entries) => { entries.forEach(e => { if (e.isIntersecting) setActive(e.target.id); }); },
      { rootMargin: '-80px 0px -60% 0px' }
    );
    toc.forEach(({ id }) => { const el = document.getElementById(id); if (el) obs.observe(el); });
    return () => obs.disconnect();
  }, [toc]);
  if (!toc.length) return null;
  return (
    <div className="sticky top-24 w-56 xl:w-64 flex-shrink-0 hidden lg:block">
      <div className="rounded-xl border bg-card/50 backdrop-blur-sm p-4">
        <p className="text-xs font-mono text-primary tracking-widest uppercase mb-3">// 目录</p>
        <div className="space-y-0.5">
          {toc.map(item => (
            <a key={item.id} href={`#${item.id}`}
              onClick={e => { e.preventDefault(); document.getElementById(item.id)?.scrollIntoView({ behavior: 'smooth' }); }}
              className={cn('block text-xs py-1 truncate transition-colors cursor-pointer',
                item.level === 2 ? 'pl-3' : item.level === 3 ? 'pl-6' : item.level >= 4 ? 'pl-9' : 'pl-0',
                active === item.id ? 'text-primary font-medium' : 'text-muted-foreground hover:text-foreground'
              )}
            >{item.text}</a>
          ))}
        </div>
      </div>
    </div>
  );
}
export default function ArticleDetailPage() {
  const { slug } = useParams<{ slug: string }>();
  const router = useRouter();
  const { isLoggedIn } = useUserAuth();
  const [article, setArticle] = useState<ArticleDetail | null>(null);
  const [comments, setComments] = useState<Comment[]>([]);
  const [toc, setToc] = useState<TocItem[]>([]);
  const [loading, setLoading] = useState(true);
  const [liked, setLiked] = useState(false);
  const [likes, setLikes] = useState(0);
  const [liking, setLiking] = useState(false);
  const [commentText, setCommentText] = useState('');
  const [submitting, setSubmitting] = useState(false);
  const [replyTo, setReplyTo] = useState<number | null>(null);

  useEffect(() => {
    if (!slug) return;
    Promise.all([getArticle(slug), getArticleComments(slug)])
      .then(([art, cmts]) => {
        setArticle(art); setComments(cmts);
        setLikes(art.like?.likes ?? 0); setLiked(art.like?.liked === true);
        setToc(extractToc(art.content || ''));
      })
      .catch(() => toast.error('文章加载失败'))
      .finally(() => setLoading(false));
  }, [slug]);

  const handleLike = async () => {
    if (!isLoggedIn) { toast.info('请先登录'); return; }
    if (liking || !slug) return;
    setLiking(true);
    try {
      const res = liked ? await removeArticleLike(slug) : await toggleArticleLike(slug);
      setLiked(res.liked === true); setLikes(res.likes);
    } catch { toast.error('操作失败'); } finally { setLiking(false); }
  };

  const handleComment = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!isLoggedIn) { toast.info('请先登录'); return; }
    if (!commentText.trim() || !slug) return;
    setSubmitting(true);
    try {
      await createComment({ article_slug: slug, parent_id: replyTo, content: commentText.trim() });
      toast.success('评论已提交，等待审核');
      setCommentText(''); setReplyTo(null);
      setComments(await getArticleComments(slug));
    } catch { toast.error('评论失败'); } finally { setSubmitting(false); }
  };

  if (loading) return (
    <div className="max-w-5xl mx-auto px-4 py-12 space-y-6">
      <Skeleton className="h-6 w-20" />
      <Skeleton className="h-12 w-4/5" />
      <Skeleton className="h-4 w-48" />
      <Skeleton className="aspect-video w-full rounded-2xl" />
      {[1,2,3,4].map(i => <Skeleton key={i} className="h-4 w-full" />)}
    </div>
  );

  if (!article) return (
    <div className="max-w-5xl mx-auto px-4 py-24 text-center space-y-4">
      <p className="font-mono text-muted-foreground">// 404 · 文章不存在</p>
      <Button onClick={() => router.back()}>返回</Button>
    </div>
  );
  return (
    <div className="max-w-6xl mx-auto px-4 py-10">
      <button type="button" onClick={() => router.back()}
        className="mb-8 inline-flex items-center gap-1.5 text-sm text-muted-foreground hover:text-foreground transition-colors cursor-pointer">
        <ArrowLeft className="h-4 w-4" /> 返回
      </button>

      <header className="mb-10">
        <div className="flex flex-wrap gap-2 mb-4">
          {article.category && <Badge className="rounded-full bg-primary/10 text-primary border-primary/20">{article.category.name}</Badge>}
          {article.is_featured && <Badge variant="secondary" className="rounded-full">精选</Badge>}
        </div>
        <h1 className="text-3xl md:text-5xl font-bold tracking-tight leading-tight mb-6">{article.title}</h1>
        <div className="flex flex-wrap items-center gap-x-6 gap-y-2 text-sm text-muted-foreground mb-5">
          <div className="flex items-center gap-2">
            <Avatar className="h-7 w-7">
              <AvatarImage src={article.author?.avatar} />
              <AvatarFallback className="text-xs">{article.author?.nickname?.[0]}</AvatarFallback>
            </Avatar>
            <span className="font-medium text-foreground">{article.author?.nickname}</span>
          </div>
          <span className="flex items-center gap-1"><Calendar className="h-3.5 w-3.5" />{formatDate(article.published_at)}</span>
          <span className="flex items-center gap-1"><Eye className="h-3.5 w-3.5" />{article.views} 次阅读</span>
        </div>
        {article.tags?.length > 0 && (
          <div className="flex flex-wrap gap-2">
            {article.tags.map(tag => (
              <span key={tag.id} className="inline-flex items-center gap-1 rounded-full border px-2.5 py-0.5 text-xs font-mono text-muted-foreground cursor-default">
                <Hash className="h-2.5 w-2.5" />{tag.name}
              </span>
            ))}
          </div>
        )}
      </header>

      {article.featured_image && (
        <div className="aspect-video rounded-2xl overflow-hidden mb-10 ring-1 ring-border">
          <img src={article.featured_image} alt={article.title} className="w-full h-full object-cover" />
        </div>
      )}

      <div className="flex gap-12 items-start">
        <main className="flex-1 min-w-0">
          <MarkdownContent content={article.content || ''} />

          <div className="flex justify-center my-16">
            <button type="button" onClick={handleLike} disabled={liking}
              className={cn('group flex items-center gap-3 rounded-full px-8 py-3.5 text-sm font-medium border transition-all cursor-pointer',
                liked ? 'bg-rose-500/10 border-rose-500/30 text-rose-500' : 'bg-card text-muted-foreground hover:text-foreground border-border'
              )}>
              <HeartIcon className={cn('h-5 w-5 transition-transform group-hover:scale-110', liked && 'fill-rose-500 text-rose-500')} />
              {likes} 人觉得不错
            </button>
          </div>

          <Separator className="mb-10" />

          <section id="comments">
            <h2 className="text-xl font-bold mb-6">评论 <span className="text-base font-normal text-muted-foreground">({comments.length})</span></h2>
            {isLoggedIn ? (
              <form onSubmit={handleComment} className="mb-10 space-y-3">
                {replyTo && (
                  <div className="flex items-center gap-2 text-xs text-muted-foreground bg-muted/50 rounded-lg px-3 py-2">
                    回复 #{replyTo}
                    <button type="button" className="underline cursor-pointer" onClick={() => setReplyTo(null)}>取消</button>
                  </div>
                )}
                <textarea
                  className="w-full rounded-xl border border-input bg-background px-4 py-3 text-sm resize-none focus:outline-none focus:ring-2 focus:ring-ring"
                  rows={4} placeholder="写下你的评论..." maxLength={2000}
                  value={commentText} onChange={e => setCommentText(e.target.value)}
                />
                <Button type="submit" disabled={submitting || !commentText.trim()} className="rounded-full">
                  {submitting ? '提交中...' : '发布评论'}
                </Button>
              </form>
            ) : (
              <p className="text-sm text-muted-foreground mb-8 p-4 rounded-xl border bg-muted/30">
                <button type="button" className="text-primary hover:underline font-medium cursor-pointer" onClick={() => {
                  const redirectUri = `${window.location.origin}/auth/callback`;
                  import('@/lib/api/public/auth').then(({ getSSOLoginUrl }) =>
                    getSSOLoginUrl(redirectUri, window.location.pathname).then(r => { window.location.href = r.sso_login_url; })
                  );
                }}>登录</button> 后参与评论
              </p>
            )}
            <div className="space-y-6">
              {comments.map(comment => (
                <div key={comment.id} className="flex gap-3 group">
                  <Avatar className="h-8 w-8 flex-shrink-0 mt-0.5">
                    <AvatarImage src={comment.user?.avatar} />
                    <AvatarFallback className="text-xs">{comment.user?.nickname?.[0]}</AvatarFallback>
                  </Avatar>
                  <div className="flex-1 rounded-xl bg-muted/30 px-4 py-3 space-y-1">
                    <div className="flex items-center justify-between">
                      <span className="text-sm font-medium">{comment.user?.nickname}</span>
                      <span className="text-xs text-muted-foreground">{formatDate(comment.created_at)}</span>
                    </div>
                    <p className="text-sm leading-relaxed">{comment.content}</p>
                    {isLoggedIn && (
                      <button type="button" className="text-xs text-muted-foreground hover:text-primary cursor-pointer transition-colors" onClick={() => setReplyTo(comment.id)}>回复</button>
                    )}
                  </div>
                </div>
              ))}
              {comments.length === 0 && (
                <p className="text-center text-muted-foreground py-10 font-mono text-sm">// 暂无评论，来说点什么吧</p>
              )}
            </div>
          </section>
        </main>
        <TableOfContents toc={toc} />
      </div>
    </div>
  );
}
