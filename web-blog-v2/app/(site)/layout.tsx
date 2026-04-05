import { Navbar } from '@/components/site/layout/Navbar';
import { Footer } from '@/components/site/layout/Footer';
import { getWebsiteInfoServer } from '@/lib/server-api/website';

export const dynamic = 'force-dynamic';

async function loadSiteInfo() {
  try {
    return await getWebsiteInfoServer();
  } catch {
    return {
      title: '博客',
      name: '博客',
      description: '',
      github_url: '',
      bilibili_url: '',
      steam_url: '',
      icp_filing: '',
    };
  }
}

export default async function SiteLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  const site = await loadSiteInfo();

  return (
    <div className="min-h-screen flex flex-col">
      <Navbar title={site.title || '博客'} />
      <main className="flex-1">{children}</main>
      <Footer site={site} />
    </div>
  );
}
