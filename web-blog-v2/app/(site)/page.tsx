'use client';

import Link from 'next/link';
import { useEffect, useState } from 'react';
import { motion } from 'framer-motion';
import type { Variants } from 'framer-motion';
import { ArrowRight, ArrowDown } from 'lucide-react';
import { GithubIcon, BilibiliIcon, SteamIcon } from '@/components/common/icons';
import { buttonVariants } from '@/components/ui/button';
import { ArticleCard } from '@/components/site/article/ArticleCard';
import { useSite } from '@/context/site';
import { listArticles } from '@/lib/api/public/article';
import { Skeleton } from '@/components/ui/skeleton';
import type { Article } from '@/lib/api/types';
import { cn } from '@/lib/utils';

const fadeUp: Variants = {
  hidden: { opacity: 0, y: 16 },
  show: { opacity: 1, y: 0, transition: { duration: 0.45 } },
};
const stagger: Variants = {
  hidden: {},
  show: { transition: { staggerChildren: 0.09 } },
};

// ─── 打字机 Hook ──────────────────────────────────────────
function useTypewriter(text: string, speed = 55, startDelay = 400) {
  const [displayed, setDisplayed] = useState('');
  const [done, setDone] = useState(false);
  useEffect(() => {
    setDisplayed('');
    setDone(false);
    if (!text) { setDone(true); return; }
    const timeout = setTimeout(() => {
      let i = 0;
      const interval = setInterval(() => {
        i++;
        setDisplayed(text.slice(0, i));
        if (i >= text.length) { clearInterval(interval); setDone(true); }
      }, speed);
      return () => clearInterval(interval);
    }, startDelay);
    return () => clearTimeout(timeout);
  }, [text, speed, startDelay]);
  return { displayed, done };
}

export default function HomePage() {
  const { site, isLoading: siteLoading } = useSite();
  const [articles, setArticles] = useState<Article[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    listArticles({ page: 1, page_size: 6, order: 'desc' })
      .then(data => setArticles(data.list))
      .catch(() => {})
      .finally(() => setLoading(false));
  }, []);

  const typeTarget = siteLoading ? '' : (site.name || 'developer');
  const { displayed: typedName, done: typeDone } = useTypewriter(typeTarget, 60, 300);

  const socialLinks = [
    site.github_url   && { href: site.github_url,   Icon: GithubIcon,   label: 'GitHub',   hoverClass: 'hover:text-foreground' },
    site.bilibili_url && { href: site.bilibili_url, Icon: BilibiliIcon, label: 'Bilibili', hoverClass: 'hover:text-[#00a1d6]' },
    site.steam_url    && { href: site.steam_url,    Icon: SteamIcon,    label: 'Steam',    hoverClass: 'hover:text-[#66c0f4]' },
  ].filter(Boolean) as { href: string; Icon: React.FC<{className?: string}>; label: string; hoverClass: string }[];

  return (
    <div className="min-h-screen">
      {/* ─── Hero ─── */}
      <section className="relative flex flex-col justify-center overflow-hidden min-h-[80vh] md:min-h-[75vh]">
        {/* bg decorations */}
        <div className="pointer-events-none absolute inset-0 -z-10"
          style={{ backgroundImage: 'radial-gradient(ellipse 70% 50% at 50% 0%, hsl(var(--primary)/0.10), transparent)' }}
        />
        <div className="pointer-events-none absolute inset-0 -z-10 opacity-[0.025] dark:opacity-[0.05]"
          style={{
            backgroundImage: 'radial-gradient(circle, hsl(var(--foreground)) 1px, transparent 1px)',
            backgroundSize: '28px 28px',
          }}
        />

        <div className="max-w-6xl mx-auto px-4 py-24 md:py-32 w-full">
          <motion.div
            className="space-y-8 max-w-2xl"
            initial="hidden" animate="show" variants={stagger}
          >
            {/* prompt line */}
            <motion.div variants={fadeUp} className="font-mono text-sm text-muted-foreground flex items-center gap-2">
              <span className="text-primary">$</span>
              <span>whoami</span>
            </motion.div>

            {/* name line */}
            <motion.div variants={fadeUp}>
              <h1 className="text-4xl md:text-6xl font-bold tracking-tight leading-[1.05] font-mono">
                <span className="bg-gradient-to-r from-violet-500 via-purple-400 to-cyan-400 bg-clip-text text-transparent">
                  {typedName}
                </span>
                <span
                  className={cn(
                    'ml-0.5 inline-block w-[3px] h-[0.9em] bg-primary align-middle',
                    typeDone ? 'animate-[blink_1s_step-end_infinite]' : 'opacity-100'
                  )}
                />
              </h1>
            </motion.div>

            {/* job + slogan */}
            <motion.div variants={fadeUp} className="space-y-2">
              {siteLoading ? (
                <Skeleton className="h-4 w-32" />
              ) : site.job ? (
                <p className="font-mono text-sm text-primary/80 tracking-widest uppercase">
                  // {site.job}
                </p>
              ) : null}
            </motion.div>

            {/* CTA buttons */}
            <motion.div variants={fadeUp} className="flex flex-wrap gap-3">
              <Link href="/articles"
                className={cn(buttonVariants({ size: 'lg' }),
                  'rounded-full bg-gradient-to-r from-violet-600 to-cyan-500 text-white border-0 hover:opacity-90 shadow-lg shadow-violet-500/20 font-mono'
                )}>
                ./read<ArrowRight className="ml-1.5 h-4 w-4" />
              </Link>
              <Link href="/about"
                className={cn(buttonVariants({ variant: 'outline', size: 'lg' }), 'rounded-full font-mono')}>
                ./about
              </Link>
            </motion.div>
          </motion.div>
        </div>

        {/* scroll hint */}
        <motion.div
          className="absolute bottom-8 left-1/2 -translate-x-1/2 text-muted-foreground/40"
          initial={{ opacity: 0, y: -6 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ delay: 1.8, duration: 0.5 }}
        >
          <ArrowDown className="h-5 w-5 animate-bounce" />
        </motion.div>
      </section>

      {/* ─── Latest Articles ─── */}
      <section className="py-20 max-w-6xl mx-auto px-4" id="posts">
        <motion.div
          initial="hidden" animate="show" variants={stagger}
        >
          <motion.div className="flex items-end justify-between mb-12" variants={fadeUp}>
            <div>
              <p className="text-xs font-mono text-primary tracking-widest uppercase mb-2">// latest_articles</p>
              <h2 className="text-3xl font-bold tracking-tight">最新文章</h2>
              <p className="text-muted-foreground mt-1.5 text-sm">分享技术思考、实践与探索</p>
            </div>
            <Link href="/articles"
              className={cn(buttonVariants({ variant: 'ghost', size: 'sm' }), 'group cursor-pointer font-mono')}>
              ls -a
              <ArrowRight className="ml-1 h-3.5 w-3.5 transition-transform group-hover:translate-x-1" />
            </Link>
          </motion.div>

          {loading ? (
            <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
              {[0,1,2].map(i => <div key={i} className="h-80 rounded-2xl bg-muted/60 animate-pulse" />)}
            </div>
          ) : articles.length > 0 ? (
            <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
              {articles.map((article, i) => (
                <motion.div key={article.id} variants={fadeUp} custom={i}>
                  <ArticleCard article={article} />
                </motion.div>
              ))}
            </div>
          ) : (
            <div className="text-center py-20">
              <p className="text-muted-foreground font-mono text-sm">// no posts yet</p>
            </div>
          )}
        </motion.div>
      </section>
    </div>
  );
}
