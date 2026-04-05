export function getPublicSiteUrl(): string {
  return process.env.NEXT_PUBLIC_SITE_URL || '';
}

export function getClientCallbackUrl(path: string): string {
  const normalizedPath = path.startsWith('/') ? path : `/${path}`;
  const siteUrl = getPublicSiteUrl();

  if (siteUrl) {
    return `${siteUrl.replace(/\/$/, '')}${normalizedPath}`;
  }

  if (typeof window !== 'undefined') {
    return `${window.location.origin}${normalizedPath}`;
  }

  return normalizedPath;
}
