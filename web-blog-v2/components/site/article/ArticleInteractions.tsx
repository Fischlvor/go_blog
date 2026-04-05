'use client';

import { useEffect, useState } from 'react';
import { ArrowDown } from 'lucide-react';
import { HeartIcon } from '@/components/common/icons';
import { Button } from '@/components/ui/button';
import { Separator } from '@/components/ui/separator';
import { Avatar, AvatarFallback, AvatarImage } from '@/components/ui/avatar';
import { toast } from 'sonner';
import { toggleArticleLike, removeArticleLike } from '@/lib/client-api/user/article';
import { getArticleComments, createComment } from '@/lib/client-api/public/comment';
import { useUserAuth } from '@/context/user-auth';
import { getClientCallbackUrl } from '@/lib/utils/site-url';
import type { Comment } from '@/lib/client-api/types';
import { cn } from '@/lib/utils';

function formatDate(d: string | null | undefined) {
  if (!d) return '';
  return new Date(d).toLocaleDateString('zh-CN', { year: 'numeric', month: 'long', day: 'numeric' });
}

interface ArticleInteractionsProps {
  slug: string;
  initialComments: Comment[];
  initialLiked: boolean;
  initialLikes: number;
}

export function ArticleInteractions({
  slug,
  initialComments,
  initialLiked,
  initialLikes,
}: ArticleInteractionsProps) {
  const { isLoggedIn } = useUserAuth();
  const [comments, setComments] = useState<Comment[]>(initialComments);
  const [liked, setLiked] = useState(initialLiked);
  const [likes, setLikes] = useState(initialLikes);
  const [liking, setLiking] = useState(false);
  const [commentText, setCommentText] = useState('');
  const [submitting, setSubmitting] = useState(false);
  const [replyTo, setReplyTo] = useState<number | null>(null);
  const [showScrollDown, setShowScrollDown] = useState(true);

  useEffect(() => {
    const onScroll = () => {
      const scrolled = window.scrollY + window.innerHeight;
      const total = document.documentElement.scrollHeight;
      setShowScrollDown(scrolled < total - 120);
    };
    window.addEventListener('scroll', onScroll, { passive: true });
    onScroll();
    return () => window.removeEventListener('scroll', onScroll);
  }, []);

  useEffect(() => {
    setComments(initialComments);
  }, [initialComments]);

  const handleLike = async () => {
    if (!isLoggedIn) {
      toast.info('请先登录');
      return;
    }
    if (liking) return;
    setLiking(true);
    try {
      const res = liked ? await removeArticleLike(slug) : await toggleArticleLike(slug);
      setLiked(res.liked === true);
      setLikes(res.likes);
    } catch {
      toast.error('操作失败');
    } finally {
      setLiking(false);
    }
  };

  const handleComment = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!isLoggedIn) {
      toast.info('请先登录');
      return;
    }
    if (!commentText.trim()) return;
    setSubmitting(true);
    try {
      await createComment({ article_slug: slug, parent_id: replyTo, content: commentText.trim() });
      toast.success('评论已提交，等待审核');
      setCommentText('');
      setReplyTo(null);
      setComments(await getArticleComments(slug));
    } catch {
      toast.error('评论失败');
    } finally {
      setSubmitting(false);
    }
  };

  const handleLogin = async () => {
    const redirectUri = getClientCallbackUrl('/auth/callback');
    const { getSSOLoginUrl } = await import('@/lib/client-api/public/auth');
    const result = await getSSOLoginUrl(redirectUri, window.location.pathname);
    window.location.href = result.sso_login_url;
  };

  return (
    <>
      <div className="flex justify-center my-12">
        <button
          type="button"
          onClick={handleLike}
          disabled={liking}
          className={cn(
            'group flex items-center gap-3 rounded-full px-8 py-3.5 text-sm font-medium border transition-all cursor-pointer',
            liked ? 'bg-rose-500/10 border-rose-500/30 text-rose-500' : 'bg-card text-muted-foreground hover:text-foreground border-border'
          )}
        >
          <HeartIcon className={cn('h-5 w-5 transition-transform group-hover:scale-110', liked && 'fill-rose-500 text-rose-500')} />
          {likes} 人觉得不错
        </button>
      </div>

      <Separator className="mb-10" />

      <section id="comments">
        <h2 className="text-xl font-bold mb-6">
          评论 <span className="text-base font-normal text-muted-foreground">({comments.length})</span>
        </h2>
        {isLoggedIn ? (
          <form onSubmit={handleComment} className="mb-10 space-y-3">
            {replyTo && (
              <div className="flex items-center gap-2 text-xs text-muted-foreground bg-muted/50 rounded-lg px-3 py-2">
                回复 #{replyTo}
                <button type="button" className="underline cursor-pointer" onClick={() => setReplyTo(null)}>
                  取消
                </button>
              </div>
            )}
            <textarea
              className="w-full rounded-xl border border-input bg-background px-4 py-3 text-sm resize-none focus:outline-none focus:ring-2 focus:ring-ring"
              rows={4}
              placeholder="写下你的评论..."
              maxLength={2000}
              value={commentText}
              onChange={(e) => setCommentText(e.target.value)}
            />
            <Button type="submit" disabled={submitting || !commentText.trim()} className="rounded-full">
              {submitting ? '提交中...' : '发布评论'}
            </Button>
          </form>
        ) : (
          <p className="text-sm text-muted-foreground mb-8 p-4 rounded-xl border bg-muted/30">
            <button type="button" className="text-primary hover:underline font-medium cursor-pointer" onClick={handleLogin}>
              登录
            </button>{' '}
            后参与评论
          </p>
        )}
        <div className="space-y-6">
          {comments.map((comment) => (
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
                  <button
                    type="button"
                    className="text-xs text-muted-foreground hover:text-primary cursor-pointer transition-colors"
                    onClick={() => setReplyTo(comment.id)}
                  >
                    回复
                  </button>
                )}
              </div>
            </div>
          ))}
          {comments.length === 0 && (
            <p className="text-center text-muted-foreground py-10 font-mono text-sm">// 暂无评论，来说点什么吧</p>
          )}
        </div>
      </section>

      {showScrollDown && (
        <button
          type="button"
          onClick={() => document.getElementById('comments')?.scrollIntoView({ behavior: 'smooth' })}
          className="fixed bottom-8 right-8 z-50 flex items-center justify-center w-10 h-10 rounded-full border border-border bg-background/80 backdrop-blur-sm shadow-md text-muted-foreground hover:text-foreground hover:border-primary/50 transition-all cursor-pointer"
          title="直达评论区"
        >
          <ArrowDown className="h-4 w-4" />
        </button>
      )}
    </>
  );
}
