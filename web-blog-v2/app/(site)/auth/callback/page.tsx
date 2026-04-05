'use client';

import { Suspense } from 'react';
import { useEffect, useRef } from 'react';
import { useRouter, useSearchParams } from 'next/navigation';
import { toast } from 'sonner';
import { handleSSOCallback } from '@/lib/client-api/public/auth';
import { useUserAuth } from '@/context/user-auth';
import { getClientCallbackUrl } from '@/lib/utils/site-url';

function SSOCallbackContent() {
  const router = useRouter();
  const searchParams = useSearchParams();
  const { setToken, refreshUser } = useUserAuth();
  const handled = useRef(false);

  useEffect(() => {
    if (handled.current) return;
    handled.current = true;

    const code = searchParams.get('code');
    const state = searchParams.get('state');
    const redirectUri = getClientCallbackUrl('/auth/callback');

    if (!code) {
      toast.error('登录失败：缺少授权码');
      router.replace('/');
      return;
    }

    (async () => {
      try {
        const res = await handleSSOCallback(code, state ?? undefined, redirectUri);
        setToken(res.access_token);
        await refreshUser();
        toast.success('登录成功');
        const returnUrl = state ? decodeURIComponent(state) : '/';
        router.replace(returnUrl.startsWith('/') ? returnUrl : '/');
      } catch (e: unknown) {
        toast.error(e instanceof Error ? e.message : '登录失败，请重试');
        router.replace('/');
      }
    })();
  }, [searchParams, router, setToken, refreshUser]);

  return (
    <div className="min-h-screen flex items-center justify-center">
      <div className="text-center space-y-3">
        <div className="inline-block h-8 w-8 animate-spin rounded-full border-4 border-primary border-r-transparent" />
        <p className="text-sm text-muted-foreground">正在完成登录...</p>
      </div>
    </div>
  );
}

export default function SSOCallbackPage() {
  return (
    <Suspense
      fallback={
        <div className="min-h-screen flex items-center justify-center">
          <div className="text-center space-y-3">
            <div className="inline-block h-8 w-8 animate-spin rounded-full border-4 border-primary border-r-transparent" />
            <p className="text-sm text-muted-foreground">正在完成登录...</p>
          </div>
        </div>
      }
    >
      <SSOCallbackContent />
    </Suspense>
  );
}
