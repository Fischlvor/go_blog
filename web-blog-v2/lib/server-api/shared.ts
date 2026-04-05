import 'server-only';

interface ApiEnvelope<T> {
  code: string;
  message: string;
  data?: T;
}

interface ServerRequestOptions extends RequestInit {
  next?: NextFetchRequestConfig;
}

function normalizeBaseUrl(value: string): string {
  return value.replace(/\/$/, '');
}

function normalizePath(value: string): string {
  return value.startsWith('/') ? value : `/${value}`;
}

export function getSiteUrl(): string {
  const siteUrl = process.env.SITE_URL || process.env.NEXT_PUBLIC_SITE_URL;
  if (!siteUrl) {
    throw new Error('Missing SITE_URL or NEXT_PUBLIC_SITE_URL');
  }
  return normalizeBaseUrl(siteUrl);
}

function getApiBase(): string {
  const versionPath = normalizePath(process.env.NEXT_PUBLIC_API_V1 || '/api/v1');
  const internalBase = process.env.SERVER_API_BASE;
  if (internalBase) {
    return `${normalizeBaseUrl(internalBase)}${versionPath}`;
  }

  const siteUrl = process.env.SITE_URL || process.env.NEXT_PUBLIC_SITE_URL;

  if (!siteUrl) {
    throw new Error('Missing SERVER_API_BASE and SITE_URL/NEXT_PUBLIC_SITE_URL');
  }

  return `${normalizeBaseUrl(siteUrl)}${versionPath}`;
}

export async function serverRequest<T>(path: string, options: ServerRequestOptions = {}): Promise<T> {
  const { headers, ...restOptions } = options;
  const response = await fetch(`${getApiBase()}${normalizePath(path)}`, {
    ...restOptions,
    headers: {
      Accept: 'application/json',
      ...headers,
    },
  });

  const json = (await response.json()) as ApiEnvelope<T>;

  if (!response.ok) {
    throw new Error(json.message || `HTTP ${response.status}`);
  }

  if (json.code !== '0000') {
    throw new Error(json.message || 'Unknown error');
  }

  return json.data as T;
}
