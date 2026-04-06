import Link from 'next/link';
import { Separator } from '@/components/ui/separator';
import { GithubIcon, BilibiliIcon, SteamIcon } from '@/components/common/icons';
import { getFieldValue, type Website } from '@/lib/client-api/public/website';

const NAV = [
  ['/', '首页'],
  ['/articles', '文章'],
  ['/archive', '归档'],
  ['/links', '友链'],
  ['/about', '关于'],
] as const;

interface FooterProps {
  site: Partial<Website>;
}

function getStartYear(createdAt?: string) {
  if (!createdAt) return new Date().getFullYear();
  const date = new Date(createdAt);
  if (Number.isNaN(date.getTime())) return new Date().getFullYear();
  return date.getFullYear();
}

function getRunningDuration(createdAt?: string) {
  if (!createdAt) return '';

  const start = new Date(createdAt);
  if (Number.isNaN(start.getTime())) return '';

  const now = new Date();
  if (start > now) return '已运行 0天';

  let years = now.getFullYear() - start.getFullYear();
  let months = now.getMonth() - start.getMonth();
  let days = now.getDate() - start.getDate();

  if (days < 0) {
    months -= 1;
    const prevMonthDays = new Date(now.getFullYear(), now.getMonth(), 0).getDate();
    days += prevMonthDays;
  }

  if (months < 0) {
    years -= 1;
    months += 12;
  }

  const parts = [];
  if (years > 0) parts.push(`${years}年`);
  if (months > 0) parts.push(`${months}月`);
  if (days > 0 || parts.length === 0) parts.push(`${days}天`);

  return `已运行 ${parts.join('')}`;
}

export function Footer({ site }: FooterProps) {
  const socials = [
    getFieldValue(site.github_url) && { href: getFieldValue(site.github_url), Icon: GithubIcon, label: 'GitHub', hover: 'hover:text-foreground' },
    getFieldValue(site.bilibili_url) && { href: getFieldValue(site.bilibili_url), Icon: BilibiliIcon, label: 'Bilibili', hover: 'hover:text-[#00a1d6]' },
    getFieldValue(site.steam_url) && { href: getFieldValue(site.steam_url), Icon: SteamIcon, label: 'Steam', hover: 'hover:text-[#66c0f4]' },
  ].filter(Boolean) as { href: string; Icon: React.FC<{ className?: string }>; label: string; hover: string }[];

  const createdAt = getFieldValue(site.created_at);
  const startYear = getStartYear(createdAt);
  const runningDuration = getRunningDuration(createdAt);
  const title = getFieldValue(site.title) || '博客';
  const name = getFieldValue(site.name) || title;
  const description = getFieldValue(site.description);
  const version = getFieldValue(site.version);
  const icpFiling = getFieldValue(site.icp_filing);

  return (
    <footer className="border-t border-border/50 bg-background/80 mt-auto">
      <div className="max-w-6xl mx-auto px-4 py-10">
        <div className="grid grid-cols-1 md:grid-cols-3 gap-8">
          <div className="space-y-3">
            <h3 className="font-bold text-lg bg-gradient-to-r from-violet-600 to-cyan-500 bg-clip-text text-transparent">
              {title}
            </h3>
            {description && (
              <p className="text-sm text-muted-foreground max-w-xs leading-relaxed">{description}</p>
            )}
          </div>

          <div className="space-y-3">
            <h4 className="text-sm font-semibold font-mono text-muted-foreground tracking-widest uppercase">{'// nav'}</h4>
            <ul className="space-y-2">
              {NAV.map(([href, label]) => (
                <li key={href}>
                  <Link href={href} className="text-sm text-muted-foreground hover:text-foreground transition-colors">
                    {label}
                  </Link>
                </li>
              ))}
            </ul>
          </div>

          <div className="space-y-3">
            <h4 className="text-sm font-semibold font-mono text-muted-foreground tracking-widest uppercase">{'// social'}</h4>
            {socials.length > 0 ? (
              <ul className="space-y-2.5">
                {socials.map(({ href, Icon, label, hover }) => (
                  <li key={label}>
                    <a
                      href={href}
                      target="_blank"
                      rel="noopener noreferrer"
                      className={`inline-flex items-center gap-2 text-sm text-muted-foreground transition-colors ${hover}`}
                    >
                      <Icon className="h-4 w-4" />
                      {label}
                    </a>
                  </li>
                ))}
              </ul>
            ) : (
              <p className="text-xs font-mono text-muted-foreground/40">{'// empty'}</p>
            )}
          </div>
        </div>

        <Separator className="my-6" />

        <div className="flex flex-col sm:flex-row justify-between items-center gap-2 text-xs text-foreground/60 font-mono">
          <div className="flex items-center gap-3 flex-wrap justify-center sm:justify-start">
            <p>{startYear} {name}. All rights reserved.</p>
            {runningDuration && <p>{runningDuration}</p>}
            {version && <p>v{version}</p>}
          </div>
          <div className="flex items-center gap-3">
            {icpFiling && (
              <a
                href="https://beian.miit.gov.cn"
                target="_blank"
                rel="noopener noreferrer"
                className="underline underline-offset-2 hover:text-muted-foreground transition-colors"
              >
                {icpFiling}
              </a>
            )}
            <p>Built with Next.js</p>
          </div>
        </div>
      </div>
    </footer>
  );
}
