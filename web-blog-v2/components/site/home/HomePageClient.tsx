'use client';

import Link from 'next/link';
import { useEffect, useRef, useState } from 'react';
import { motion } from 'framer-motion';
import type { Variants } from 'framer-motion';
import { ArrowRight, ArrowDown } from 'lucide-react';
import { buttonVariants } from '@/components/ui/button';
import { ArticleCard } from '@/components/site/article/ArticleCard';
import type { Article } from '@/lib/client-api/types';
import type { Website } from '@/lib/client-api/public/website';
import { cn } from '@/lib/utils';

const fadeUp: Variants = {
  hidden: { opacity: 0, y: 16 },
  show: { opacity: 1, y: 0, transition: { duration: 0.45 } },
};
const stagger: Variants = {
  hidden: {},
  show: { transition: { staggerChildren: 0.09 } },
};

function useTypewriter(text: string, speed = 55, startDelay = 400) {
  const [displayed, setDisplayed] = useState('');
  const [done, setDone] = useState(false);

  useEffect(() => {
    setDisplayed('');
    setDone(false);
    if (!text) {
      setDone(true);
      return;
    }
    const timeout = setTimeout(() => {
      let i = 0;
      const interval = setInterval(() => {
        i++;
        setDisplayed(text.slice(0, i));
        if (i >= text.length) {
          clearInterval(interval);
          setDone(true);
        }
      }, speed);
      return () => clearInterval(interval);
    }, startDelay);
    return () => clearTimeout(timeout);
  }, [text, speed, startDelay]);

  return { displayed, done };
}

interface HomePageClientProps {
  site: Partial<Website>;
  articles: Article[];
}

export function HomePageClient({ site, articles }: HomePageClientProps) {
  const postsRef = useRef<HTMLElement>(null);
  const heroRef = useRef<HTMLDivElement>(null);
  const { displayed: typedName, done: typeDone } = useTypewriter(site.name || 'developer', 60, 300);

  const scrollToPosts = () => {
    postsRef.current?.scrollIntoView({ behavior: 'smooth', block: 'start' });
  };

  return (
    <div className="min-h-screen">
      <section ref={heroRef} className="relative flex flex-col justify-center items-center overflow-hidden min-h-[calc(100vh-4rem)]">
        <div
          className="pointer-events-none absolute inset-0 -z-10"
          style={{ backgroundImage: 'radial-gradient(ellipse 70% 50% at 50% 0%, hsl(var(--primary)/0.10), transparent)' }}
        />
        <div
          className="pointer-events-none absolute inset-0 -z-10 opacity-[0.025] dark:opacity-[0.05]"
          style={{
            backgroundImage: 'radial-gradient(circle, hsl(var(--foreground)) 1px, transparent 1px)',
            backgroundSize: '28px 28px',
          }}
        />

        <div className="max-w-2xl mx-auto px-4 w-full relative z-10 text-center">
          <motion.div className="space-y-6" initial="hidden" animate="show" variants={stagger}>
            <motion.div variants={fadeUp} className="font-mono text-left inline-block">
              <div className="px-6 py-4 text-sm min-w-[280px]">
                <div className="flex items-center gap-2">
                  <span className="text-green-500">~</span>
                  <span className="text-primary">$</span>
                  <span className="text-foreground">whoami</span>
                </div>
                <div className="mt-1 min-h-[4.5rem] md:min-h-[5.5rem] flex items-center">
                  <h1 className="text-4xl md:text-6xl font-bold tracking-tight leading-[1.05]">
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
                </div>
              </div>
            </motion.div>

            <motion.div variants={fadeUp} className="flex flex-wrap gap-3 justify-center pt-4">
              <Link
                href="/articles"
                className={cn(
                  buttonVariants({ size: 'lg' }),
                  'rounded-full bg-gradient-to-r from-violet-600 to-cyan-500 text-white border-0 hover:opacity-90 shadow-lg shadow-violet-500/20 font-mono'
                )}
              >
                ./read<ArrowRight className="ml-1.5 h-4 w-4" />
              </Link>
              <Link href="/about" className={cn(buttonVariants({ variant: 'outline', size: 'lg' }), 'rounded-full font-mono')}>
                ./about
              </Link>
            </motion.div>
          </motion.div>
        </div>

        <motion.button
          type="button"
          onClick={scrollToPosts}
          className="absolute bottom-10 left-1/2 -translate-x-1/2 text-muted-foreground/40 hover:text-muted-foreground transition-colors cursor-pointer z-20"
          initial={{ opacity: 0, y: -6 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ delay: 1.8, duration: 0.5 }}
        >
          <ArrowDown className="h-5 w-5 animate-bounce" />
        </motion.button>
      </section>

      <section ref={postsRef} className="py-20 max-w-6xl mx-auto px-4" id="posts">
        <motion.div initial="hidden" animate="show" variants={stagger}>
          <motion.div className="flex items-end justify-between mb-12" variants={fadeUp}>
            <div>
              <p className="text-xs font-mono text-primary tracking-widest uppercase mb-2">// latest_articles</p>
              <h2 className="text-3xl font-bold tracking-tight">最新文章</h2>
              <p className="text-muted-foreground mt-1.5 text-sm">分享技术思考、实践与探索</p>
            </div>
            <Link href="/articles" className={cn(buttonVariants({ variant: 'ghost', size: 'sm' }), 'group cursor-pointer font-mono')}>
              ls -a
              <ArrowRight className="ml-1 h-3.5 w-3.5 transition-transform group-hover:translate-x-1" />
            </Link>
          </motion.div>

          {articles.length > 0 ? (
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
