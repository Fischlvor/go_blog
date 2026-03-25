'use client';

import { ThemeProvider } from 'next-themes';
import { UserAuthProvider } from '@/context/user-auth';
import { SiteProvider } from '@/context/site';
import { Toaster } from '@/components/ui/sonner';

export function Providers({ children }: { children: React.ReactNode }) {
  return (
    <ThemeProvider
      attribute="class"
      defaultTheme="dark"
      enableSystem
      disableTransitionOnChange
    >
      <SiteProvider>
        <UserAuthProvider>
          {children}
        </UserAuthProvider>
      </SiteProvider>
      <Toaster richColors position="top-right" />
    </ThemeProvider>
  );
}
