'use client';

import { createContext, useContext, useEffect, useState } from 'react';
import { getWebsiteInfo } from '@/lib/api/public/website';
import type { Website } from '@/lib/api/public/website';

interface SiteContextValue {
  site: Partial<Website>;
  isLoading: boolean;
}

const SiteContext = createContext<SiteContextValue>({ site: {}, isLoading: true });

export function SiteProvider({ children }: { children: React.ReactNode }) {
  const [site, setSite] = useState<Partial<Website>>({});
  const [isLoading, setIsLoading] = useState(true);

  useEffect(() => {
    getWebsiteInfo()
      .then(setSite)
      .catch(() => {})
      .finally(() => setIsLoading(false));
  }, []);

  return (
    <SiteContext.Provider value={{ site, isLoading }}>
      {children}
    </SiteContext.Provider>
  );
}

export function useSite(): SiteContextValue {
  return useContext(SiteContext);
}
