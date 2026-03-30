'use client';

import { Suspense } from 'react';
import { useEffect } from 'react';
import { useRouter, useSearchParams } from 'next/navigation';
import { toast } from 'sonner';
import { useUserAuth } from '@/context/user-auth';
import { handleSSOCallback } from '@/lib/api/public/auth';

function SSOCallbackContent() {
  const router = useRouter();
  const searchParams = useSearchParams();
  const { setToken, refreshUser } = useUserAuth();

  useEffect(() => {
    const code = searchParams.get('code');
    const state = searchParams.get('state') || '';
    const redirectUri = `${window.location.origin}/sso-callback`;

    // 解析 state 中的 return_url
    let returnUrl = '/';
    try {
      const parsed = JSON.parse(decodeURIComponent(state));
      if (parsed.return_url) returnUrl = parsed.return_url;
    } catch {
      // state 不是 JSON，忽略
    }

    if (!code) {
      toast.error('登录失败，缺少授权码');
      router.replace('/login');
      return;
    }

    handleSSOCallback(code, state, redirectUri)
      .then(res => {
        setToken(res.access_token);
        return refreshUser();
      })
      .then(() => {
        toast.success('登录成功');
        router.replace(returnUrl);
      })
      .catch(() => {
        toast.error('登录失败，请重试');
        router.replace('/login');
      });
  // eslint-disable-next-line react-hooks/exhaustive-deps
  }, []);

  return (
    <div className="min-h-screen flex items-center justify-center">
      <div className="text-center space-y-3">
        <div className="h-8 w-8 border-2 border-primary border-t-transparent rounded-full animate-spin mx-auto" />
        <p className="text-muted-foreground text-sm">正在处理登录...</p>
      </div>
    </div>
  );
}

export default function SSOCallbackPage() {
  return (
    <Suspense fallback={
      <div className="min-h-screen flex items-center justify-center">
        <div className="text-center space-y-3">
          <div className="h-8 w-8 border-2 border-primary border-t-transparent rounded-full animate-spin mx-auto" />
          <p className="text-muted-foreground text-sm">正在处理登录...</p>
        </div>
      </div>
    }>
      <SSOCallbackContent />
    </Suspense>
  );
}
