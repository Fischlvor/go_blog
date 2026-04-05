'use client';

import { useEffect, useState } from 'react';
import { cn } from '@/lib/utils';

interface TocItem {
  id: string;
  text: string;
  level: number;
}

interface TableOfContentsProps {
  toc: TocItem[];
}

export function TableOfContents({ toc }: TableOfContentsProps) {
  const [active, setActive] = useState('');

  useEffect(() => {
    const obs = new IntersectionObserver(
      (entries) => {
        entries.forEach((entry) => {
          if (entry.isIntersecting) setActive(entry.target.id);
        });
      },
      { rootMargin: '-80px 0px -60% 0px' }
    );

    toc.forEach(({ id }) => {
      const el = document.getElementById(id);
      if (el) obs.observe(el);
    });

    return () => obs.disconnect();
  }, [toc]);

  if (!toc.length) return null;

  return (
    <div className="sticky top-24 w-56 xl:w-64 flex-shrink-0 hidden lg:block">
      <div className="rounded-xl border bg-card/50 backdrop-blur-sm p-4">
        <p className="text-xs font-mono text-primary tracking-widest uppercase mb-3">// 目录</p>
        <div className="space-y-0.5 max-h-[calc(100vh-8rem)] overflow-y-auto pr-1">
          {toc.map((item) => (
            <a
              key={item.id}
              href={`#${item.id}`}
              onClick={(e) => {
                e.preventDefault();
                document.getElementById(item.id)?.scrollIntoView({ behavior: 'smooth' });
              }}
              className={cn(
                'block text-xs py-1 truncate transition-colors cursor-pointer',
                item.level === 2 ? 'pl-3' : item.level === 3 ? 'pl-6' : item.level >= 4 ? 'pl-9' : 'pl-0',
                active === item.id ? 'text-primary font-medium' : 'text-muted-foreground hover:text-foreground'
              )}
            >
              {item.text}
            </a>
          ))}
        </div>
      </div>
    </div>
  );
}
