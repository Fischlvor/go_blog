/**
 * 管理员配置接口
 */

import { adminRequest } from '../http';
import type { Website } from '../public/website';

export interface SystemConfig {
  use_multipoint: boolean;
  sessions_secret: string;
  oss_type: string;
}

export interface EmailConfig {
  host: string;
  port: number;
  from: string;
  nickname: string;
  secret: string;
  is_ssl: boolean;
}

export interface QQConfig {
  enable: boolean;
  app_id: string;
  app_key: string;
  redirect_uri: string;
}

export interface QiniuConfig {
  zone: string;
  bucket: string;
  img_path: string;
  access_key: string;
  secret_key: string;
  use_https: boolean;
  use_cdn_domains: boolean;
}

export interface JwtConfig {
  access_token_secret: string;
  refresh_token_secret: string;
  access_token_expiry_time: string;
  refresh_token_expiry_time: string;
  issuer: string;
}

export interface GaodeConfig {
  enable: boolean;
  key: string;
}

export function adminGetWebsiteConfig(): Promise<Website> {
  return adminRequest<Website>('/config/website');
}

export function adminUpdateWebsiteConfig(data: Website): Promise<void> {
  return adminRequest<void>('/config/website', {
    method: 'PUT',
    body: JSON.stringify(data),
  });
}

export function adminGetSystemConfig(): Promise<SystemConfig> {
  return adminRequest<SystemConfig>('/config/system');
}

export function adminUpdateSystemConfig(data: SystemConfig): Promise<void> {
  return adminRequest<void>('/config/system', {
    method: 'PUT',
    body: JSON.stringify(data),
  });
}

export function adminGetEmailConfig(): Promise<EmailConfig> {
  return adminRequest<EmailConfig>('/config/email');
}

export function adminUpdateEmailConfig(data: EmailConfig): Promise<void> {
  return adminRequest<void>('/config/email', {
    method: 'PUT',
    body: JSON.stringify(data),
  });
}

export function adminGetQQConfig(): Promise<QQConfig> {
  return adminRequest<QQConfig>('/config/qq');
}

export function adminUpdateQQConfig(data: QQConfig): Promise<void> {
  return adminRequest<void>('/config/qq', {
    method: 'PUT',
    body: JSON.stringify(data),
  });
}

export function adminGetQiniuConfig(): Promise<QiniuConfig> {
  return adminRequest<QiniuConfig>('/config/qiniu');
}

export function adminUpdateQiniuConfig(data: QiniuConfig): Promise<void> {
  return adminRequest<void>('/config/qiniu', {
    method: 'PUT',
    body: JSON.stringify(data),
  });
}

export function adminGetJwtConfig(): Promise<JwtConfig> {
  return adminRequest<JwtConfig>('/config/jwt');
}

export function adminUpdateJwtConfig(data: JwtConfig): Promise<void> {
  return adminRequest<void>('/config/jwt', {
    method: 'PUT',
    body: JSON.stringify(data),
  });
}

export function adminGetGaodeConfig(): Promise<GaodeConfig> {
  return adminRequest<GaodeConfig>('/config/gaode');
}

export function adminUpdateGaodeConfig(data: GaodeConfig): Promise<void> {
  return adminRequest<void>('/config/gaode', {
    method: 'PUT',
    body: JSON.stringify(data),
  });
}
