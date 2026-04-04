'use client';

import { ThemeProvider } from 'next-themes';
import { UserAuthProvider } from '@/context/user-auth';
import { Toaster } from '@/components/ui/sonner';

export function Providers({ children }: { children: React.ReactNode }) {
  return (
    <ThemeProvider
      attribute="class"
      defaultTheme="dark"
      enableSystem
      disableTransitionOnChange
    >
      <UserAuthProvider>
        {children}
      </UserAuthProvider>
      <Toaster richColors position="top-right" />
    </ThemeProvider>
  );
}
