/**
 * HTTP 请求封装层
 * - 自动处理 token（从 localStorage 读取）
 * - 自动解包响应（{ code, message, data }）
 * - 自动刷新 token（从响应头 X-New-Access-Token）
 * - 401 时清除登录态
 */

export interface ApiEnvelope<T> {
  code: string;
  message: string;
  data?: T;
}

export interface ApiError extends Error {
  code?: string;
  status?: number;
}

/**
 * API 版本前缀配置
 * 通过环境变量控制，支持未来扩展到 v2、v3 等
 */
export const API_PREFIXES = {
  v1: process.env.NEXT_PUBLIC_API_V1 || '/api/v1',
  v2: process.env.NEXT_PUBLIC_API_V2 || '/api/v2',
} as const;

export type ApiVersion = keyof typeof API_PREFIXES;

/** 默认版本为 v1 */

function getToken(): string | null {
  if (typeof window === 'undefined') return null;
  return localStorage.getItem('access_token');
}

function setToken(token: string): void {
  if (typeof window === 'undefined') return;
  localStorage.setItem('access_token', token);
}

function clearToken(): void {
  if (typeof window === 'undefined') return;
  localStorage.removeItem('access_token');
}

async function doFetch<T>(
  url: string,
  options: RequestInit & { withAuth?: boolean; apiVersion?: ApiVersion } = {}
): Promise<T> {
  const { withAuth = false, apiVersion = 'v1', ...fetchOptions } = options;
  const API_BASE = API_PREFIXES[apiVersion];

  const headers = new Headers(fetchOptions.headers || {});
  headers.set('Content-Type', 'application/json');

  if (withAuth) {
    const token = getToken();
    if (token) {
      headers.set('Authorization', `Bearer ${token}`);
    }
  }

  const fullUrl = url.startsWith('http') ? url : `${API_BASE}${url}`;

  const response = await fetch(fullUrl, {
    ...fetchOptions,
    headers,
  });

  // 检查响应头中的新 token
  const newToken = response.headers.get('X-New-Access-Token');
  if (newToken) {
    setToken(newToken);
  }

  const json = (await response.json()) as ApiEnvelope<T>;

  // 401 清除登录态
  if (response.status === 401) {
    clearToken();
    const error: ApiError = new Error(json.message || 'Unauthorized');
    error.code = json.code;
    error.status = 401;
    throw error;
  }

  // 非 200 系列状态码
  if (!response.ok) {
    const error: ApiError = new Error(json.message || `HTTP ${response.status}`);
    error.code = json.code;
    error.status = response.status;
    throw error;
  }

  // code !== '0000' 视为业务错误
  if (json.code !== '0000') {
    const error: ApiError = new Error(json.message || 'Unknown error');
    error.code = json.code;
    throw error;
  }

  return json.data as T;
}

/**
 * 公开接口（无需 token）
 */
export function publicRequest<T>(
  url: string,
  options?: RequestInit & { apiVersion?: ApiVersion }
): Promise<T> {
  return doFetch<T>(url, { ...options, withAuth: false });
}

/**
 * 需要用户认证的接口（带 token）
 */
export function authRequest<T>(
  url: string,
  options?: RequestInit & { apiVersion?: ApiVersion }
): Promise<T> {
  return doFetch<T>(url, { ...options, withAuth: true });
}

/**
 * 管理员接口（带 token）
 */
export function adminRequest<T>(
  url: string,
  options?: RequestInit & { apiVersion?: ApiVersion }
): Promise<T> {
  return doFetch<T>(url, { ...options, withAuth: true });
}
