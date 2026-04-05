/**
 * 网站公开接口
 */

import { publicRequest } from '../http';

export interface Website {
  avatar: string;
  logo: string;
  full_logo: string;
  title: string;
  slogan: string;
  slogan_en: string;
  description: string;
  version: string;
  created_at: string;
  icp_filing: string;
  public_security_filing: string;
  bilibili_url: string;
  github_url: string;
  steam_url: string;
  name: string;
  job: string;
  address: string;
  email: string;
  qq_image: string;
  wechat_image: string;
}

export async function getWebsiteInfo(): Promise<Website> {
  return publicRequest<Website>('/website/info');
}

export async function getWebsiteCarousel(): Promise<string[]> {
  return publicRequest<string[]>('/website/carousel');
}

export interface HotItem {
  index: number;
  title: string;
  description: string;
  image: string;
  popularity: string;
  url: string;
}

export interface HotSearchData {
  source: string;
  update_time: string;
  hot_list: HotItem[];
}

export async function getWebsiteNews(source: string): Promise<HotSearchData> {
  return publicRequest<HotSearchData>(`/website/news?source=${source}`);
}

export interface FooterLink {
  title: string;
  link: string;
}

export async function getFooterLinks(): Promise<FooterLink[]> {
  return publicRequest<FooterLink[]>('/website/footerLink');
}
