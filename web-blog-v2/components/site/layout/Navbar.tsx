'use client';

import Link from 'next/link';
import { usePathname } from 'next/navigation';
import { useEffect, useState } from 'react';
import { Menu, X, Sun, Moon } from 'lucide-react';
import { useTheme } from 'next-themes';
import { Button, buttonVariants } from '@/components/ui/button';
import { Sheet, SheetContent, SheetTrigger } from '@/components/ui/sheet';
import { cn } from '@/lib/utils';
import { useUserAuth } from '@/context/user-auth';
import { UserMenu } from '@/components/site/auth/UserMenu';
import { getSSOLoginUrl } from '@/lib/client-api/public/auth';
import { getClientCallbackUrl } from '@/lib/utils/site-url';

const NAV_ITEMS = [
  { href: '/', label: '首页' },
  { href: '/articles', label: '文章' },
  { href: '/archive', label: '归档' },
  { href: '/links', label: '友链' },
  { href: '/about', label: '关于' },
];

function ThemeToggle() {
  const { theme, setTheme } = useTheme();
  const [mounted, setMounted] = useState(false);
  useEffect(() => setMounted(true), []);
  if (!mounted) return <div className="w-8 h-8" />;

  const isDark = theme === 'dark';

  return (
    <Button variant="ghost" size="icon" onClick={() => setTheme(isDark ? 'light' : 'dark')} aria-label="切换主题">
      {isDark ? <Sun className="h-4 w-4" /> : <Moon className="h-4 w-4" />}
    </Button>
  );
}

function NavbarActions() {
  const pathname = usePathname();
  const { isLoggedIn } = useUserAuth();
  const [open, setOpen] = useState(false);
  const [ssoLoading, setSsoLoading] = useState(false);

  const handleSSOLogin = async () => {
    setSsoLoading(true);
    try {
      const redirectUri = getClientCallbackUrl('/auth/callback');
      const returnUrl = pathname;
      const res = await getSSOLoginUrl(redirectUri, returnUrl);
      window.location.href = res.sso_login_url;
    } finally {
      setSsoLoading(false);
    }
  };

  return (
    <>
      <nav className="hidden md:flex items-center gap-1">
        {NAV_ITEMS.map((item) => (
          <Link
            key={item.href}
            href={item.href}
            className={cn(
              'px-3 py-1.5 rounded-md text-sm font-medium transition-colors',
              pathname === item.href
                ? 'bg-primary/10 text-primary'
                : 'text-muted-foreground hover:text-foreground hover:bg-accent'
            )}
          >
            {item.label}
          </Link>
        ))}
      </nav>

      <div className="flex items-center gap-1">
        <ThemeToggle />
        <div className="hidden md:block">
          {isLoggedIn ? (
            <UserMenu />
          ) : (
            <Button size="sm" variant="outline" onClick={handleSSOLogin} disabled={ssoLoading}>
              {ssoLoading ? '跳转中...' : '登录'}
            </Button>
          )}
        </div>

        <Sheet open={open} onOpenChange={setOpen}>
          <SheetTrigger
            render={
              <button
                type="button"
                className={cn(buttonVariants({ variant: 'ghost', size: 'icon' }), 'md:hidden')}
                aria-label="菜单"
              />
            }
          >
            {open ? <X className="h-5 w-5" /> : <Menu className="h-5 w-5" />}
          </SheetTrigger>
          <SheetContent side="right" className="w-64 pt-10">
            <nav className="flex flex-col gap-1">
              {NAV_ITEMS.map((item) => (
                <Link
                  key={item.href}
                  href={item.href}
                  onClick={() => setOpen(false)}
                  className={cn(
                    'px-4 py-2.5 rounded-md text-sm font-medium transition-colors',
                    pathname === item.href
                      ? 'bg-primary/10 text-primary'
                      : 'text-muted-foreground hover:text-foreground hover:bg-accent'
                  )}
                >
                  {item.label}
                </Link>
              ))}
              <div className="pt-4 border-t mt-4">
                {isLoggedIn ? (
                  <UserMenu mobile />
                ) : (
                  <Button className="w-full" onClick={() => { setOpen(false); handleSSOLogin(); }} disabled={ssoLoading}>
                    {ssoLoading ? '跳转中...' : '登录'}
                  </Button>
                )}
              </div>
            </nav>
          </SheetContent>
        </Sheet>
      </div>
    </>
  );
}

interface NavbarProps {
  title: string;
}

export function Navbar({ title }: NavbarProps) {
  return (
    <header className="sticky top-0 z-50 w-full border-b border-border/50 bg-background/80 backdrop-blur-md">
      <div className="max-w-6xl mx-auto px-4 h-16 flex items-center justify-between">
        <Link href="/" className="flex items-center gap-2 font-bold text-lg tracking-tight">
          <span className="bg-gradient-to-r from-violet-600 to-cyan-500 bg-clip-text text-transparent">
            {title}
          </span>
        </Link>
        <NavbarActions />
      </div>
    </header>
  );
}
