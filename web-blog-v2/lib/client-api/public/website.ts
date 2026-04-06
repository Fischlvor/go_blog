/**
 * 网站公开接口
 */

import { publicRequest } from '../http';

export interface SettingField {
  value: string;
  setting_key: string;
}

export function getFieldValue(field?: SettingField | null): string {
  return field?.value ?? '';
}

export interface Website {
  avatar: SettingField;
  title: SettingField;
  description: SettingField;
  profile_intro: SettingField;
  tech_stack: SettingField;
  work_experiences: SettingField;
  version: SettingField;
  created_at: SettingField;
  icp_filing: SettingField;
  bilibili_url: SettingField;
  github_url: SettingField;
  steam_url: SettingField;
  name: SettingField;
  job: SettingField;
  address: SettingField;
  email: SettingField;
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
