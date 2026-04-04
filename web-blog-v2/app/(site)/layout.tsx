import { Navbar } from '@/components/site/layout/Navbar';
import { Footer } from '@/components/site/layout/Footer';
import { SiteProvider } from '@/context/site';

export default function SiteLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <SiteProvider>
      <div className="min-h-screen flex flex-col">
        <Navbar />
        <main className="flex-1">{children}</main>
        <Footer />
      </div>
    </SiteProvider>
  );
}
