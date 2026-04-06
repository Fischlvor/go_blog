import { Navbar } from '@/components/site/layout/Navbar';
import { Footer } from '@/components/site/layout/Footer';
import { getWebsiteInfoServer } from '@/lib/server-api/website';
import { getFieldValue } from '@/lib/client-api/public/website';

export const dynamic = 'force-dynamic';

export default async function SiteLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  const site = await getWebsiteInfoServer();

  return (
    <div className="min-h-screen flex flex-col">
      <Navbar title={getFieldValue(site.title) || '博客'} />
      <main className="flex-1">{children}</main>
      <Footer site={site} />
    </div>
  );
}
