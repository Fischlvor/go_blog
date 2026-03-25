'use client';

import { useEffect, useState } from 'react';
import { ExternalLink } from 'lucide-react';
import { Card, CardContent } from '@/components/ui/card';
import { Skeleton } from '@/components/ui/skeleton';
import { getFriendLinks } from '@/lib/api/public/friendLink';
import type { FriendLink } from '@/lib/api/types';

export default function LinksPage() {
  const [links, setLinks] = useState<FriendLink[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    getFriendLinks().then(setLinks).catch(() => {}).finally(() => setLoading(false));
  }, []);

  return (
    <div className="max-w-4xl mx-auto px-4 py-12">
      <div className="mb-10">
        <h1 className="text-3xl font-bold tracking-tight">友情链接</h1>
        <p className="text-muted-foreground mt-1">互联网上志同道合的朋友们</p>
      </div>

      {loading ? (
        <div className="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 gap-4">
          {Array.from({ length: 6 }).map((_, i) => <Skeleton key={i} className="h-24 rounded-xl" />)}
        </div>
      ) : links.length > 0 ? (
        <div className="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 gap-4">
          {links.map(link => (
            <a
              key={link.id}
              href={link.url}
              target="_blank"
              rel="noopener noreferrer"
              className="group block"
            >
              <Card className="h-full transition-all duration-200 group-hover:border-primary/40 group-hover:shadow-md group-hover:-translate-y-0.5">
                <CardContent className="p-4 flex items-center gap-3">
                  {link.logo ? (
                    <img
                      src={link.logo}
                      alt={link.name}
                      className="w-10 h-10 rounded-lg object-cover flex-shrink-0"
                      onError={e => { (e.target as HTMLImageElement).style.display = 'none'; }}
                    />
                  ) : (
                    <div className="w-10 h-10 rounded-lg bg-primary/10 flex items-center justify-center flex-shrink-0">
                      <span className="text-lg font-bold text-primary">{link.name[0]}</span>
                    </div>
                  )}
                  <div className="flex-1 min-w-0">
                    <div className="flex items-center gap-1">
                      <h3 className="font-medium text-sm truncate">{link.name}</h3>
                      <ExternalLink className="h-3 w-3 text-muted-foreground flex-shrink-0 opacity-0 group-hover:opacity-100 transition-opacity" />
                    </div>
                    {link.description && (
                      <p className="text-xs text-muted-foreground line-clamp-1 mt-0.5">{link.description}</p>
                    )}
                  </div>
                </CardContent>
              </Card>
            </a>
          ))}
        </div>
      ) : (
        <p className="text-center text-muted-foreground py-20">暂无友链</p>
      )}
    </div>
  );
}
