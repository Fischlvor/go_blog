'use client';

export const dynamic = 'force-dynamic';

import { useEffect, useState } from 'react';
import { toast } from 'sonner';
import {
  adminGetEmailConfig,
  adminGetGaodeConfig,
  adminGetJwtConfig,
  adminGetQiniuConfig,
  adminGetQQConfig,
  adminGetSystemConfig,
  adminGetWebsiteConfig,
  adminUpdateEmailConfig,
  adminUpdateGaodeConfig,
  adminUpdateJwtConfig,
  adminUpdateQiniuConfig,
  adminUpdateQQConfig,
  adminUpdateSystemConfig,
  adminUpdateWebsiteConfig,
  type EmailConfig,
  type GaodeConfig,
  type JwtConfig,
  type QiniuConfig,
  type QQConfig,
  type SystemConfig,
} from '@/lib/client-api/admin/config';
import type { Website } from '@/lib/client-api/public/website';
import { Button } from '@/components/ui/button';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Input } from '@/components/ui/input';

export default function AdminAppConfigPage() {
  const [website, setWebsite] = useState<Website | null>(null);
  const [system, setSystem] = useState<SystemConfig | null>(null);
  const [email, setEmail] = useState<EmailConfig | null>(null);
  const [qq, setQQ] = useState<QQConfig | null>(null);
  const [qiniu, setQiniu] = useState<QiniuConfig | null>(null);
  const [jwt, setJwt] = useState<JwtConfig | null>(null);
  const [gaode, setGaode] = useState<GaodeConfig | null>(null);

  useEffect(() => {
    adminGetWebsiteConfig().then(setWebsite).catch(() => setWebsite(null));
    adminGetSystemConfig().then(setSystem).catch(() => setSystem(null));
    adminGetEmailConfig().then(setEmail).catch(() => setEmail(null));
    adminGetQQConfig().then(setQQ).catch(() => setQQ(null));
    adminGetQiniuConfig().then(setQiniu).catch(() => setQiniu(null));
    adminGetJwtConfig().then(setJwt).catch(() => setJwt(null));
    adminGetGaodeConfig().then(setGaode).catch(() => setGaode(null));
  }, []);

  const save = async (fn: () => Promise<void>) => {
    try { await fn(); toast.success('保存成功'); } catch { toast.error('保存失败'); }
  };

  return (
    <div className="space-y-4">
      <Card><CardHeader><CardTitle>网站配置</CardTitle></CardHeader><CardContent className="space-y-2">{website ? <><Input value={website.title || ''} onChange={(e) => setWebsite({ ...website, title: e.target.value })} placeholder="标题" /><Input value={website.name || ''} onChange={(e) => setWebsite({ ...website, name: e.target.value })} placeholder="名称" /><Input value={website.job || ''} onChange={(e) => setWebsite({ ...website, job: e.target.value })} placeholder="职业" /><Input value={website.icp_filing || ''} onChange={(e) => setWebsite({ ...website, icp_filing: e.target.value })} placeholder="ICP备案" /><Button onClick={() => save(() => adminUpdateWebsiteConfig(website))}>保存网站配置</Button></> : <p className="text-sm text-muted-foreground">加载失败</p>}</CardContent></Card>

      <Card><CardHeader><CardTitle>系统配置</CardTitle></CardHeader><CardContent className="space-y-2">{system ? <><label className="flex items-center gap-2 text-sm"><input type="checkbox" checked={!!system.use_multipoint} onChange={(e) => setSystem({ ...system, use_multipoint: e.target.checked })} /> 多点登录</label><Input value={system.sessions_secret || ''} onChange={(e) => setSystem({ ...system, sessions_secret: e.target.value })} placeholder="sessions_secret" /><Input value={system.oss_type || ''} onChange={(e) => setSystem({ ...system, oss_type: e.target.value })} placeholder="oss_type" /><Button onClick={() => save(() => adminUpdateSystemConfig(system))}>保存系统配置</Button></> : <p className="text-sm text-muted-foreground">加载失败</p>}</CardContent></Card>

      <Card><CardHeader><CardTitle>邮箱配置</CardTitle></CardHeader><CardContent className="space-y-2">{email ? <><Input value={email.host || ''} onChange={(e) => setEmail({ ...email, host: e.target.value })} placeholder="host" /><Input value={String(email.port || '')} onChange={(e) => setEmail({ ...email, port: Number(e.target.value) || 0 })} placeholder="port" /><Input value={email.from || ''} onChange={(e) => setEmail({ ...email, from: e.target.value })} placeholder="from" /><Button onClick={() => save(() => adminUpdateEmailConfig(email))}>保存邮箱配置</Button></> : <p className="text-sm text-muted-foreground">加载失败</p>}</CardContent></Card>

      <Card><CardHeader><CardTitle>QQ登录配置</CardTitle></CardHeader><CardContent className="space-y-2">{qq ? <><label className="flex items-center gap-2 text-sm"><input type="checkbox" checked={!!qq.enable} onChange={(e) => setQQ({ ...qq, enable: e.target.checked })} /> 启用</label><Input value={qq.app_id || ''} onChange={(e) => setQQ({ ...qq, app_id: e.target.value })} placeholder="app_id" /><Input value={qq.app_key || ''} onChange={(e) => setQQ({ ...qq, app_key: e.target.value })} placeholder="app_key" /><Input value={qq.redirect_uri || ''} onChange={(e) => setQQ({ ...qq, redirect_uri: e.target.value })} placeholder="redirect_uri" /><Button onClick={() => save(() => adminUpdateQQConfig(qq))}>保存QQ配置</Button></> : <p className="text-sm text-muted-foreground">加载失败</p>}</CardContent></Card>

      <Card><CardHeader><CardTitle>七牛云配置</CardTitle></CardHeader><CardContent className="space-y-2">{qiniu ? <><Input value={qiniu.zone || ''} onChange={(e) => setQiniu({ ...qiniu, zone: e.target.value })} placeholder="zone" /><Input value={qiniu.bucket || ''} onChange={(e) => setQiniu({ ...qiniu, bucket: e.target.value })} placeholder="bucket" /><Input value={qiniu.img_path || ''} onChange={(e) => setQiniu({ ...qiniu, img_path: e.target.value })} placeholder="img_path" /><Button onClick={() => save(() => adminUpdateQiniuConfig(qiniu))}>保存七牛云配置</Button></> : <p className="text-sm text-muted-foreground">加载失败</p>}</CardContent></Card>

      <Card><CardHeader><CardTitle>JWT配置</CardTitle></CardHeader><CardContent className="space-y-2">{jwt ? <><Input value={jwt.access_token_secret || ''} onChange={(e) => setJwt({ ...jwt, access_token_secret: e.target.value })} placeholder="access_token_secret" /><Input value={jwt.refresh_token_secret || ''} onChange={(e) => setJwt({ ...jwt, refresh_token_secret: e.target.value })} placeholder="refresh_token_secret" /><Input value={jwt.access_token_expiry_time || ''} onChange={(e) => setJwt({ ...jwt, access_token_expiry_time: e.target.value })} placeholder="access_token_expiry_time" /><Button onClick={() => save(() => adminUpdateJwtConfig(jwt))}>保存JWT配置</Button></> : <p className="text-sm text-muted-foreground">加载失败</p>}</CardContent></Card>

      <Card><CardHeader><CardTitle>高德配置</CardTitle></CardHeader><CardContent className="space-y-2">{gaode ? <><label className="flex items-center gap-2 text-sm"><input type="checkbox" checked={!!gaode.enable} onChange={(e) => setGaode({ ...gaode, enable: e.target.checked })} /> 启用</label><Input value={gaode.key || ''} onChange={(e) => setGaode({ ...gaode, key: e.target.value })} placeholder="key" /><Button onClick={() => save(() => adminUpdateGaodeConfig(gaode))}>保存高德配置</Button></> : <p className="text-sm text-muted-foreground">加载失败</p>}</CardContent></Card>
    </div>
  );
}
