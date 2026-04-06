import Image from 'next/image';
import { Mail } from 'lucide-react';
import { GithubIcon } from '@/components/common/icons';
import { getWebsiteInfoServer } from '@/lib/server-api/website';
import { getFieldValue } from '@/lib/client-api/public/website';

type WorkExperience = {
  title?: string;
  company?: string;
  period?: string;
  description?: string;
};

type TechStackItem = {
  name?: string;
};

function parseWorkExperiences(raw: string): WorkExperience[] {
  try {
    const parsed = raw ? JSON.parse(raw) : [];
    return Array.isArray(parsed) ? parsed : [];
  } catch {
    return [];
  }
}

function parseTechStack(raw: string): TechStackItem[] {
  try {
    const parsed = raw ? JSON.parse(raw) : [];
    if (!Array.isArray(parsed)) {
      return [];
    }
    return parsed.map((item) => ({
      name:
        typeof item === 'string'
          ? item
          : typeof item?.name === 'string'
            ? item.name
            : '',
    }));
  } catch {
    return [];
  }
}

export default async function AboutPage() {
  const site = await getWebsiteInfoServer();
  const workExperiences = parseWorkExperiences(getFieldValue(site.work_experiences));
  const techStackItems = parseTechStack(getFieldValue(site.tech_stack)).filter((item) => item.name?.trim());

  return (
    <div className="max-w-5xl mx-auto px-4 py-12 space-y-12">
      <section className="flex flex-col sm:flex-row items-center sm:items-start gap-8">
        <div className="flex-shrink-0">
          <Image
            src={getFieldValue(site.avatar) || '/avatar.png'}
            alt={getFieldValue(site.name) || '博主'}
            width={120}
            height={120}
            className="rounded-full ring-4 ring-border object-cover"
          />
        </div>
        <div className="space-y-3 text-center sm:text-left">
          <h1 className="text-3xl font-bold tracking-tight">{getFieldValue(site.name) || '博主'}</h1>
          <p className="whitespace-pre-line text-muted-foreground leading-relaxed">{getFieldValue(site.profile_intro) || '个人介绍'}</p>
          <div className="flex items-center gap-3 justify-center sm:justify-start pt-1">
            {getFieldValue(site.github_url) && (
              <a
                href={getFieldValue(site.github_url)}
                target="_blank"
                rel="noopener noreferrer"
                className="inline-flex items-center gap-2 rounded-full border border-border/70 bg-muted/75 px-4 py-2 text-sm text-foreground/85 shadow-sm backdrop-blur-sm transition-all hover:border-foreground/20 hover:bg-accent/70 hover:text-foreground"
              >
                <GithubIcon className="h-4 w-4" />
                GitHub
              </a>
            )}
            {getFieldValue(site.email) && (
              <a
                href={`mailto:${getFieldValue(site.email)}`}
                className="inline-flex items-center gap-2 rounded-full border border-border/70 bg-muted/75 px-4 py-2 text-sm text-foreground/85 shadow-sm backdrop-blur-sm transition-all hover:border-foreground/20 hover:bg-accent/70 hover:text-foreground"
              >
                <Mail className="h-4 w-4" />
                Email
              </a>
            )}
          </div>
        </div>
      </section>

      {techStackItems.length > 0 && (
        <section className="space-y-4 rounded-3xl border border-border/70 bg-card/70 p-6 shadow-sm backdrop-blur-sm sm:p-8">
          <div className="space-y-1">
            <h2 className="text-2xl font-semibold tracking-tight">技术栈</h2>
          </div>

          <div className="flex flex-wrap gap-3">
            {techStackItems.map((item, index) => (
              <span
                key={`${item.name || 'stack'}-${index}`}
                className="inline-flex items-center rounded-full border border-blue-500/30 bg-blue-500/12 px-4 py-2 text-sm font-medium text-blue-700 shadow-sm transition-colors hover:bg-blue-500/18 dark:border-blue-400/30 dark:bg-blue-400/16 dark:text-blue-200 dark:hover:bg-blue-400/22"
              >
                {item.name}
              </span>
            ))}
          </div>
        </section>
      )}

      {workExperiences.length > 0 && (
        <section className="rounded-3xl border border-border/70 bg-card/70 p-6 shadow-sm backdrop-blur-sm sm:p-8">
          <div className="mb-6 flex items-center justify-between gap-4">
            <div>
              <h2 className="text-2xl font-semibold tracking-tight">工作经历</h2>
            </div>
          </div>

          <div className="space-y-6">
            {workExperiences.map((item, index) => (
              <div
                key={`${item.title || 'work'}-${item.company || 'company'}-${index}`}
                className="rounded-2xl border border-border/60 bg-background/70 p-5 shadow-sm transition-colors hover:bg-accent/30"
              >
                <div className="flex flex-col gap-3 sm:flex-row sm:items-start sm:justify-between">
                  <div className="space-y-1">
                    <h3 className="text-lg font-semibold text-foreground">{item.title || '未命名经历'}</h3>
                    {item.company && <p className="text-base font-medium text-primary">{item.company}</p>}
                  </div>
                  {item.period && (
                    <span className="inline-flex w-fit items-center rounded-full border border-border/70 bg-background/80 px-3 py-1 text-xs font-medium text-muted-foreground">
                      {item.period}
                    </span>
                  )}
                </div>

                {item.description && (
                  <p className="mt-4 whitespace-pre-line text-base leading-7 text-foreground">{item.description}</p>
                )}
              </div>
            ))}
          </div>
        </section>
      )}
    </div>
  );
}
