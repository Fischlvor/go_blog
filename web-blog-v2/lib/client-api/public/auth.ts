/**
 * 公开认证接口
 */

import { publicRequest } from '../http';
import type { SSOLoginUrlResponse } from '../types';

export interface SSOCallbackResponse {
  access_token: string;
  token_type: string;
  expires_in: number;
  user_info: {
    uuid: string;
    nickname: string;
    avatar: string;
    email?: string;
  };
  state: string;
}

/**
 * 获取 SSO 登录 URL
 * GET /api/v1/auth/sso_login_url
 */
export async function getSSOLoginUrl(redirectUri: string, returnUrl?: string): Promise<SSOLoginUrlResponse> {
  const params = new URLSearchParams();
  params.set('redirect_uri', redirectUri);
  if (returnUrl) params.set('return_url', returnUrl);
  return publicRequest<SSOLoginUrlResponse>(`/auth/sso_login_url?${params.toString()}`);
}

/**
 * SSO 回调处理
 * GET /api/v1/auth/callback?code=...&state=...&redirect_uri=...
 */
export async function handleSSOCallback(
  code: string,
  state?: string,
  redirectUri?: string,
): Promise<SSOCallbackResponse> {
  const params = new URLSearchParams();
  params.set('code', code);
  if (state) params.set('state', state);
  if (redirectUri) params.set('redirect_uri', redirectUri);
  return publicRequest<SSOCallbackResponse>(`/auth/callback?${params.toString()}`);
}
