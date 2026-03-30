'use client';

import Image from 'next/image';
import { useEffect, useState } from 'react';
import { Mail } from 'lucide-react';
import { GithubIcon, MailIcon } from '@/components/common/icons';
import { Badge } from '@/components/ui/badge';
import { Skeleton } from '@/components/ui/skeleton';
import { Separator } from '@/components/ui/separator';
import { getWebsiteInfo } from '@/lib/api/public/website';
import type { Website } from '@/lib/api/public/website';

export default function AboutPage() {
  const [site, setSite] = useState<Partial<Website>>({});
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    getWebsiteInfo().then(setSite).catch(() => {}).finally(() => setLoading(false));
  }, []);

  if (loading) {
    return (
      <div className="max-w-3xl mx-auto px-4 py-12 space-y-6">
        <Skeleton className="h-40 w-40 rounded-full" />
        <Skeleton className="h-8 w-48" />
        <Skeleton className="h-24 w-full" />
      </div>
    );
  }

  return (
    <div className="max-w-3xl mx-auto px-4 py-12">
      {/* Profile */}
      <section className="flex flex-col sm:flex-row items-center sm:items-start gap-8 mb-12">
        <div className="flex-shrink-0">
          <Image
            src="/avatar.png" alt={site.name || '博主'}
            width={120} height={120}
            className="rounded-full ring-4 ring-border object-cover"
          />
        </div>
        <div className="space-y-3 text-center sm:text-left">
          <h1 className="text-3xl font-bold tracking-tight">{site.name || '博主'}</h1>
          {site.job && <p className="text-muted-foreground font-mono text-sm">{site.job}</p>}
          {site.address && <p className="text-sm text-muted-foreground">📍 {site.address}</p>}
          <div className="flex items-center gap-3 justify-center sm:justify-start pt-1">
            {site.github_url && (
              <a href={site.github_url} target="_blank" rel="noopener noreferrer"
                className="text-muted-foreground hover:text-foreground transition-colors">
                <GithubIcon className="h-5 w-5" />
              </a>
            )}
            {site.email && (
              <a href={`mailto:${site.email}`}
                className="text-muted-foreground hover:text-foreground transition-colors">
                <Mail className="h-5 w-5" />
              </a>
            )}
          </div>
        </div>
      </section>

      <Separator className="mb-12" />

      {/* Description */}
      {site.description && (
        <section className="mb-12">
          <h2 className="text-xl font-bold mb-4">关于我</h2>
          <p className="text-muted-foreground leading-relaxed whitespace-pre-line">{site.description}</p>
        </section>
      )}

      {/* Contact */}
      <section>
        <h2 className="text-xl font-bold mb-4">联系方式</h2>
        <div className="flex flex-wrap gap-3">
          {site.github_url && (
            <a href={site.github_url} target="_blank" rel="noopener noreferrer">
              <Badge variant="outline" className="gap-1.5 py-1.5 px-3">
                <GithubIcon className="h-3.5 w-3.5" /> GitHub
              </Badge>
            </a>
          )}
          {site.email && (
            <a href={`mailto:${site.email}`}>
              <Badge variant="outline" className="gap-1.5 py-1.5 px-3">
                <Mail className="h-3.5 w-3.5" /> {site.email}
              </Badge>
            </a>
          )}
        </div>
      </section>
    </div>
  );
}
