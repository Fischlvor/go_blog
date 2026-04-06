'use client';

export const dynamic = 'force-dynamic';

import { useEffect, useMemo, useRef, useState, type ReactNode } from 'react';
import { ChevronDown, ChevronUp, PencilLine, Plus, Trash2, X } from 'lucide-react';
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
import type { SettingField, Website } from '@/lib/client-api/public/website';
import { Button } from '@/components/ui/button';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Dialog, DialogContent, DialogDescription, DialogFooter, DialogHeader, DialogTitle } from '@/components/ui/dialog';
import { Input } from '@/components/ui/input';
import { Textarea } from '@/components/ui/textarea';

type WorkExperience = {
  title: string;
  company: string;
  period: string;
  description: string;
};

type TechStackItem = {
  name: string;
};

const EMPTY_WORK_EXPERIENCE: WorkExperience = {
  title: '',
  company: '',
  period: '',
  description: '',
};

const EMPTY_TECH_STACK_ITEM: TechStackItem = {
  name: '',
};

function parseWorkExperiences(raw: string): WorkExperience[] {
  try {
    const parsed = raw ? JSON.parse(raw) : [];
    if (!Array.isArray(parsed)) {
      return [];
    }
    return parsed.map((item) => ({
      title: typeof item?.title === 'string' ? item.title : '',
      company: typeof item?.company === 'string' ? item.company : '',
      period: typeof item?.period === 'string' ? item.period : '',
      description: typeof item?.description === 'string' ? item.description : '',
    }));
  } catch {
    return [];
  }
}

function parseTechStack(raw: string): TechStackItem[] {
  try {
    const parsed = raw ? JSON.parse(raw) : [];
    if (!Array.isArray(parsed)) {
      return [];
    }
    return parsed.map((item) => ({
      name:
        typeof item === 'string'
          ? item
          : typeof item?.name === 'string'
            ? item.name
            : '',
    }));
  } catch {
    return [];
  }
}

function stringifyWorkExperiences(items: WorkExperience[]): string {
  const normalized = items.map((item) => ({
    title: item.title,
    company: item.company,
    period: item.period,
    description: item.description,
  }));

  return JSON.stringify(normalized);
}

function stringifyTechStack(items: TechStackItem[]): string {
  const normalized = items.map((item) => ({
    name: item.name,
  }));

  return JSON.stringify(normalized);
}

type WebsiteFieldProps = {
  label: string;
  field: SettingField;
  placeholder: string;
  onChange: (value: string) => void;
  multiline?: boolean;
  textareaClassName?: string;
};

type CollapsibleConfigCardProps = {
  title: string;
  children: ReactNode;
  defaultOpen?: boolean;
};

function WebsiteField({ label, field, placeholder, onChange, multiline = false, textareaClassName }: WebsiteFieldProps) {
  return (
    <div className="space-y-2">
      <div className="space-y-1">
        <label className="text-sm font-medium text-foreground">{label}</label>
        <p className="text-xs text-muted-foreground">setting_key: {field.setting_key || '未设置'}</p>
      </div>
      {multiline ? (
        <Textarea value={field.value} onChange={(e) => onChange(e.target.value)} placeholder={placeholder} className={textareaClassName ?? 'min-h-28'} />
      ) : (
        <Input value={field.value} onChange={(e) => onChange(e.target.value)} placeholder={placeholder} />
      )}
    </div>
  );
}

const EMPTY_FIELD: SettingField = {
  value: '',
  setting_key: '',
};

function ensureWebsiteFields(website: Website): Website {
  return {
    ...website,
    avatar: website.avatar ?? EMPTY_FIELD,
    title: website.title ?? EMPTY_FIELD,
    description: website.description ?? EMPTY_FIELD,
    profile_intro: website.profile_intro ?? EMPTY_FIELD,
    tech_stack: website.tech_stack ?? { value: '', setting_key: 'profile.tech_stack' },
    work_experiences: website.work_experiences ?? EMPTY_FIELD,
    version: website.version ?? EMPTY_FIELD,
    created_at: website.created_at ?? EMPTY_FIELD,
    icp_filing: website.icp_filing ?? EMPTY_FIELD,
    bilibili_url: website.bilibili_url ?? EMPTY_FIELD,
    github_url: website.github_url ?? EMPTY_FIELD,
    steam_url: website.steam_url ?? { value: '', setting_key: 'profile.steam_url' },
    name: website.name ?? EMPTY_FIELD,
    job: website.job ?? EMPTY_FIELD,
    address: website.address ?? EMPTY_FIELD,
    email: website.email ?? EMPTY_FIELD,
  };
}

function CollapsibleConfigCard({ title, children, defaultOpen = false }: CollapsibleConfigCardProps) {
  const [open, setOpen] = useState(defaultOpen);

  return (
    <Card>
      <CardHeader className="cursor-pointer select-none" onClick={() => setOpen((prev) => !prev)}>
        <div className="flex w-full items-center justify-between gap-3">
          <CardTitle>{title}</CardTitle>
          <ChevronDown className={`h-4 w-4 text-muted-foreground transition-transform ${open ? 'rotate-180' : ''}`} />
        </div>
      </CardHeader>
      {open && <CardContent className="space-y-4">{children}</CardContent>}
    </Card>
  );
}

export default function AdminAppConfigPage() {
  const [website, setWebsite] = useState<Website | null>(null);
  const [system, setSystem] = useState<SystemConfig | null>(null);
  const [email, setEmail] = useState<EmailConfig | null>(null);
  const [qq, setQQ] = useState<QQConfig | null>(null);
  const [qiniu, setQiniu] = useState<QiniuConfig | null>(null);
  const [jwt, setJwt] = useState<JwtConfig | null>(null);
  const [gaode, setGaode] = useState<GaodeConfig | null>(null);
  const [workDialogOpen, setWorkDialogOpen] = useState(false);
  const [editingWorkIndex, setEditingWorkIndex] = useState<number | null>(null);
  const [workDraft, setWorkDraft] = useState<WorkExperience>(EMPTY_WORK_EXPERIENCE);
  const [draggingTechStackIndex, setDraggingTechStackIndex] = useState<number | null>(null);
  const techStackIdRef = useRef(0);
  const [techStackKeys, setTechStackKeys] = useState<string[]>([]);
  const techStackInputRefs = useRef<Array<HTMLInputElement | null>>([]);
  const pendingFocusTechStackIndexRef = useRef<number | null>(null);

  useEffect(() => {
    adminGetWebsiteConfig().then((data) => setWebsite(ensureWebsiteFields(data))).catch(() => setWebsite(null));
    adminGetSystemConfig().then(setSystem).catch(() => setSystem(null));
    adminGetEmailConfig().then(setEmail).catch(() => setEmail(null));
    adminGetQQConfig().then(setQQ).catch(() => setQQ(null));
    adminGetQiniuConfig().then(setQiniu).catch(() => setQiniu(null));
    adminGetJwtConfig().then(setJwt).catch(() => setJwt(null));
    adminGetGaodeConfig().then(setGaode).catch(() => setGaode(null));
  }, []);

  const workExperiences = useMemo(() => parseWorkExperiences(website?.work_experiences.value || ''), [website?.work_experiences.value]);
  const techStackItems = useMemo(() => parseTechStack(website?.tech_stack.value || ''), [website?.tech_stack.value]);

  useEffect(() => {
    setTechStackKeys((prev) => {
      const next = [...prev];
      while (next.length < techStackItems.length) {
        techStackIdRef.current += 1;
        next.push(`tech-stack-${techStackIdRef.current}`);
      }
      return next.slice(0, techStackItems.length);
    });
  }, [techStackItems.length]);

  useEffect(() => {
    const targetIndex = pendingFocusTechStackIndexRef.current;
    if (targetIndex === null) {
      return;
    }

    const input = techStackInputRefs.current[targetIndex];
    if (!input) {
      return;
    }

    pendingFocusTechStackIndexRef.current = null;
    input.focus();
    input.select();
    input.scrollIntoView({ block: 'nearest', inline: 'nearest', behavior: 'smooth' });
  }, [techStackItems.length]);

  const save = async (fn: () => Promise<void>) => {
    try {
      await fn();
      toast.success('保存成功');
    } catch {
      toast.error('保存失败');
    }
  };

  const updateWebsiteField = <K extends keyof Website>(key: K, value: string) => {
    setWebsite((prev) => (prev ? { ...prev, [key]: { ...(prev[key] ?? EMPTY_FIELD), value } } : prev));
  };

  const removeWorkExperience = (index: number) => {
    updateWebsiteField(
      'work_experiences',
      stringifyWorkExperiences(workExperiences.filter((_, idx) => idx !== index)),
    );
  };

  const moveWorkExperience = (index: number, direction: 'up' | 'down') => {
    const targetIndex = direction === 'up' ? index - 1 : index + 1;
    if (targetIndex < 0 || targetIndex >= workExperiences.length) {
      return;
    }

    const next = [...workExperiences];
    [next[index], next[targetIndex]] = [next[targetIndex], next[index]];
    updateWebsiteField('work_experiences', stringifyWorkExperiences(next));
  };

  const openCreateWorkDialog = () => {
    setEditingWorkIndex(null);
    setWorkDraft(EMPTY_WORK_EXPERIENCE);
    setWorkDialogOpen(true);
  };

  const openEditWorkDialog = (index: number) => {
    setEditingWorkIndex(index);
    setWorkDraft(workExperiences[index] ?? EMPTY_WORK_EXPERIENCE);
    setWorkDialogOpen(true);
  };

  const saveWorkDraft = () => {
    const normalized = {
      title: workDraft.title.trim(),
      company: workDraft.company.trim(),
      period: workDraft.period.trim(),
      description: workDraft.description.trim(),
    };

    const next = [...workExperiences];
    if (editingWorkIndex === null) {
      next.push(normalized);
    } else {
      next[editingWorkIndex] = normalized;
    }

    updateWebsiteField('work_experiences', stringifyWorkExperiences(next));
    setWorkDialogOpen(false);
    setEditingWorkIndex(null);
    setWorkDraft(EMPTY_WORK_EXPERIENCE);
  };

  const updateTechStackItem = (index: number, value: string) => {
    const next = techStackItems.map((item, idx) => (idx === index ? { ...item, name: value } : item));
    updateWebsiteField('tech_stack', stringifyTechStack(next));
  };

  const addTechStackItem = () => {
    pendingFocusTechStackIndexRef.current = techStackItems.length;
    techStackIdRef.current += 1;
    setTechStackKeys((prev) => [...prev, `tech-stack-${techStackIdRef.current}`]);
    updateWebsiteField('tech_stack', stringifyTechStack([...techStackItems, EMPTY_TECH_STACK_ITEM]));
  };

  const removeTechStackItem = (index: number) => {
    setTechStackKeys((prev) => prev.filter((_, idx) => idx !== index));
    updateWebsiteField('tech_stack', stringifyTechStack(techStackItems.filter((_, idx) => idx !== index)));
  };

  const moveTechStackItem = (fromIndex: number, toIndex: number) => {
    if (fromIndex === toIndex || fromIndex < 0 || toIndex < 0 || fromIndex >= techStackItems.length || toIndex >= techStackItems.length) {
      return;
    }

    const nextItems = [...techStackItems];
    const [movedItem] = nextItems.splice(fromIndex, 1);
    nextItems.splice(toIndex, 0, movedItem);

    const nextKeys = [...techStackKeys];
    const [movedKey] = nextKeys.splice(fromIndex, 1);
    nextKeys.splice(toIndex, 0, movedKey);

    setTechStackKeys(nextKeys);
    updateWebsiteField('tech_stack', stringifyTechStack(nextItems));
  };

  const cleanupTechStackItem = (index: number, value: string) => {
    if (value.trim()) {
      return;
    }
    removeTechStackItem(index);
  };

  return (
    <div className="space-y-4">
      <CollapsibleConfigCard title="网站配置" defaultOpen>
        {website ? (
          <>
            <div className="grid gap-4 md:grid-cols-2">
              <WebsiteField label="网站标题" field={website.title} onChange={(value) => updateWebsiteField('title', value)} placeholder="请输入网站标题" />
              <WebsiteField label="网站描述" field={website.description} onChange={(value) => updateWebsiteField('description', value)} placeholder="请输入网站描述" />
              <WebsiteField label="版本号" field={website.version} onChange={(value) => updateWebsiteField('version', value)} placeholder="请输入版本号" />
              <WebsiteField label="创建日期" field={website.created_at} onChange={(value) => updateWebsiteField('created_at', value)} placeholder="请输入创建日期" />
              <WebsiteField label="名称" field={website.name} onChange={(value) => updateWebsiteField('name', value)} placeholder="请输入名称" />
              <WebsiteField label="职业" field={website.job} onChange={(value) => updateWebsiteField('job', value)} placeholder="请输入职业" />
              <WebsiteField label="地址" field={website.address} onChange={(value) => updateWebsiteField('address', value)} placeholder="请输入地址" />
              <WebsiteField label="联系邮箱" field={website.email} onChange={(value) => updateWebsiteField('email', value)} placeholder="请输入联系邮箱" />
              <WebsiteField label="ICP备案号" field={website.icp_filing} onChange={(value) => updateWebsiteField('icp_filing', value)} placeholder="请输入ICP备案号" />
              <WebsiteField label="头像 URL" field={website.avatar} onChange={(value) => updateWebsiteField('avatar', value)} placeholder="请输入头像 URL" />
              <WebsiteField label="GitHub URL" field={website.github_url} onChange={(value) => updateWebsiteField('github_url', value)} placeholder="请输入 GitHub URL" />
              <WebsiteField label="Bilibili URL" field={website.bilibili_url} onChange={(value) => updateWebsiteField('bilibili_url', value)} placeholder="请输入 Bilibili URL" />
              <WebsiteField label="Steam URL" field={website.steam_url} onChange={(value) => updateWebsiteField('steam_url', value)} placeholder="请输入 Steam URL" />
            </div>

            <div className="space-y-4 rounded-xl border border-border/70 bg-muted/20 p-4">
              <div className="flex items-center justify-between gap-3">
                <div>
                  <h3 className="text-sm font-semibold text-foreground">技术栈</h3>
                  <p className="mt-1 text-xs text-muted-foreground">setting_key: {website.tech_stack.setting_key || '未设置'}</p>
                  <p className="mt-1 text-xs text-muted-foreground">点击新增后会出现可直接填写的胶囊，空白胶囊失焦后会自动移除</p>
                </div>
                <Button type="button" variant="outline" size="sm" onClick={addTechStackItem}>
                  <Plus className="mr-1 h-4 w-4" />
                  新增技术栈
                </Button>
              </div>

              {techStackItems.length > 0 ? (
                <div className="h-[96px] overflow-y-auto overflow-x-hidden pr-1 [scrollbar-width:thin]">
                  <div className="flex flex-wrap content-start gap-3">
                    {techStackItems.map((item, index) => (
                      <div
                        key={techStackKeys[index] ?? `tech-stack-${index}`}
                        draggable
                        onDragStart={() => setDraggingTechStackIndex(index)}
                        onDragOver={(e) => e.preventDefault()}
                        onDrop={() => {
                          if (draggingTechStackIndex === null) {
                            return;
                          }
                          moveTechStackItem(draggingTechStackIndex, index);
                          setDraggingTechStackIndex(null);
                        }}
                        onDragEnd={() => setDraggingTechStackIndex(null)}
                        className="inline-flex h-[42px] max-w-full items-center gap-2 rounded-full border border-border/70 bg-background px-3 py-2 shadow-sm"
                      >
                        <Input
                          ref={(node) => {
                            techStackInputRefs.current[index] = node;
                          }}
                          value={item.name}
                          onChange={(e) => updateTechStackItem(index, e.target.value)}
                          onBlur={(e) => cleanupTechStackItem(index, e.target.value)}
                          onKeyDown={(e) => {
                            if (e.key === 'Enter') {
                              (e.target as HTMLInputElement).blur();
                            }
                          }}
                          placeholder="输入技术栈"
                          className="h-7 min-w-[96px] border-0 bg-transparent px-1 text-sm shadow-none focus-visible:ring-0"
                        />
                        <button
                          type="button"
                          onClick={() => removeTechStackItem(index)}
                          className="rounded-full p-1 text-muted-foreground transition-colors hover:bg-accent hover:text-foreground"
                          aria-label={`删除技术栈 ${index + 1}`}
                        >
                          <X className="h-3.5 w-3.5" />
                        </button>
                      </div>
                    ))}
                  </div>
                </div>
              ) : (
                <div className="rounded-lg border border-dashed border-border/70 bg-background/60 px-4 py-6 text-sm text-muted-foreground">
                  暂无技术栈，点击右上角“新增技术栈”开始填写。
                </div>
              )}
            </div>

            <div className="grid gap-4 lg:grid-cols-[minmax(0,0.95fr)_minmax(0,1.05fr)] lg:items-start">
              <div className="rounded-xl border border-border/70 bg-muted/20 p-4 h-[230px] overflow-hidden">
                <WebsiteField
                  label="个人介绍"
                  field={website.profile_intro}
                  onChange={(value) => updateWebsiteField('profile_intro', value)}
                  placeholder="直接输入个人介绍，回车会保留为换行"
                  multiline
                  textareaClassName="h-[180px] resize-none overflow-y-auto"
                />
              </div>

              <div className="space-y-4 rounded-xl border border-border/70 bg-muted/20 p-4 h-[230px] overflow-hidden flex flex-col">
                <div className="flex items-center justify-between gap-3">
                  <div>
                    <h3 className="text-sm font-semibold text-foreground">工作经历</h3>
                    <p className="mt-1 text-xs text-muted-foreground">setting_key: {website.work_experiences.setting_key || '未设置'}</p>
                  </div>
                  <Button type="button" variant="outline" size="sm" onClick={openCreateWorkDialog}>
                    <Plus className="mr-1 h-4 w-4" />
                    新增经历
                  </Button>
                </div>

                <div className="min-h-0 flex-1 overflow-hidden">

                {workExperiences.length > 0 ? (
                  <div className="h-full space-y-3 overflow-y-auto pr-1">
                    {workExperiences.map((item, index) => (
                      <div
                        key={`${item.title}-${item.company}-${index}`}
                        className="flex items-center justify-between gap-3 rounded-xl border border-border bg-background px-4 py-3 shadow-sm"
                      >
                        <button
                          type="button"
                          onClick={() => openEditWorkDialog(index)}
                          className="flex min-w-0 flex-1 items-center justify-between gap-3 text-left"
                        >
                          <div className="min-w-0">
                            <p className="truncate text-sm font-semibold text-foreground">{item.company || '未填写公司'}</p>
                            <p className="truncate text-xs text-muted-foreground">{item.title || '未填写职位'}</p>
                          </div>
                          <PencilLine className="h-4 w-4 shrink-0 text-muted-foreground" />
                        </button>
                        <div className="flex shrink-0 items-center gap-1">
                          <div className="flex w-[72px] items-center justify-end gap-1">
                            {index > 0 ? (
                              <Button type="button" variant="ghost" size="icon" onClick={() => moveWorkExperience(index, 'up')}>
                                <ChevronUp className="h-4 w-4" />
                              </Button>
                            ) : (
                              <span className="h-9 w-9" aria-hidden="true" />
                            )}
                            {index < workExperiences.length - 1 ? (
                              <Button type="button" variant="ghost" size="icon" onClick={() => moveWorkExperience(index, 'down')}>
                                <ChevronDown className="h-4 w-4" />
                              </Button>
                            ) : (
                              <span className="h-9 w-9" aria-hidden="true" />
                            )}
                          </div>
                          <Button type="button" variant="ghost" size="icon" onClick={() => removeWorkExperience(index)}>
                            <Trash2 className="h-4 w-4" />
                          </Button>
                        </div>
                      </div>
                    ))}
                  </div>
                ) : (
                  <div className="rounded-lg border border-dashed border-border/70 bg-background/60 px-4 py-6 text-sm text-muted-foreground">
                    暂无工作经历，点击右上角“新增经历”开始填写。
                  </div>
                )}
                </div>
              </div>
            </div>

            <Dialog open={workDialogOpen} onOpenChange={setWorkDialogOpen}>
              <DialogContent className="sm:max-w-xl">
                <DialogHeader>
                  <DialogTitle>{editingWorkIndex === null ? '新增工作经历' : `编辑工作经历 ${editingWorkIndex + 1}`}</DialogTitle>
                  <DialogDescription>在弹窗中配置单条经历的公司、职位、时间和描述，确认后写回列表。</DialogDescription>
                </DialogHeader>

                <div className="space-y-4">
                  <div className="grid gap-4 md:grid-cols-2">
                    <Input value={workDraft.company} onChange={(e) => setWorkDraft((prev) => ({ ...prev, company: e.target.value }))} placeholder="公司 / 团队" />
                    <Input value={workDraft.title} onChange={(e) => setWorkDraft((prev) => ({ ...prev, title: e.target.value }))} placeholder="职位 / 岗位名称" />
                    <Input value={workDraft.period} onChange={(e) => setWorkDraft((prev) => ({ ...prev, period: e.target.value }))} placeholder="时间段，如：2023 - 至今" className="md:col-span-2" />
                  </div>
                  <Textarea value={workDraft.description} onChange={(e) => setWorkDraft((prev) => ({ ...prev, description: e.target.value }))} placeholder="请输入这段经历的工作内容与亮点" className="min-h-32" />
                </div>

                <DialogFooter>
                  <Button type="button" variant="outline" onClick={() => setWorkDialogOpen(false)}>
                    取消
                  </Button>
                  <Button type="button" onClick={saveWorkDraft}>
                    确认
                  </Button>
                </DialogFooter>
              </DialogContent>
            </Dialog>

            <Button onClick={() => save(() => adminUpdateWebsiteConfig(website))}>保存网站配置</Button>
          </>
        ) : (
          <p className="text-sm text-muted-foreground">加载失败</p>
        )}
      </CollapsibleConfigCard>

      <CollapsibleConfigCard title="系统配置">
        {system ? (
          <>
            <label className="flex items-center gap-2 text-sm"><input type="checkbox" checked={!!system.use_multipoint} onChange={(e) => setSystem({ ...system, use_multipoint: e.target.checked })} /> 多点登录</label>
            <Input value={system.sessions_secret || ''} onChange={(e) => setSystem({ ...system, sessions_secret: e.target.value })} placeholder="sessions_secret" />
            <Input value={system.oss_type || ''} onChange={(e) => setSystem({ ...system, oss_type: e.target.value })} placeholder="oss_type" />
            <Button onClick={() => save(() => adminUpdateSystemConfig(system))}>保存系统配置</Button>
          </>
        ) : (
          <p className="text-sm text-muted-foreground">加载失败</p>
        )}
      </CollapsibleConfigCard>

      <CollapsibleConfigCard title="邮箱配置">
        {email ? (
          <>
            <Input value={email.host || ''} onChange={(e) => setEmail({ ...email, host: e.target.value })} placeholder="host" />
            <Input value={String(email.port || '')} onChange={(e) => setEmail({ ...email, port: Number(e.target.value) || 0 })} placeholder="port" />
            <Input value={email.from || ''} onChange={(e) => setEmail({ ...email, from: e.target.value })} placeholder="from" />
            <Button onClick={() => save(() => adminUpdateEmailConfig(email))}>保存邮箱配置</Button>
          </>
        ) : (
          <p className="text-sm text-muted-foreground">加载失败</p>
        )}
      </CollapsibleConfigCard>

      <CollapsibleConfigCard title="QQ登录配置">
        {qq ? (
          <>
            <label className="flex items-center gap-2 text-sm"><input type="checkbox" checked={!!qq.enable} onChange={(e) => setQQ({ ...qq, enable: e.target.checked })} /> 启用</label>
            <Input value={qq.app_id || ''} onChange={(e) => setQQ({ ...qq, app_id: e.target.value })} placeholder="app_id" />
            <Input value={qq.app_key || ''} onChange={(e) => setQQ({ ...qq, app_key: e.target.value })} placeholder="app_key" />
            <Input value={qq.redirect_uri || ''} onChange={(e) => setQQ({ ...qq, redirect_uri: e.target.value })} placeholder="redirect_uri" />
            <Button onClick={() => save(() => adminUpdateQQConfig(qq))}>保存QQ配置</Button>
          </>
        ) : (
          <p className="text-sm text-muted-foreground">加载失败</p>
        )}
      </CollapsibleConfigCard>

      <CollapsibleConfigCard title="七牛云配置">
        {qiniu ? (
          <>
            <Input value={qiniu.zone || ''} onChange={(e) => setQiniu({ ...qiniu, zone: e.target.value })} placeholder="zone" />
            <Input value={qiniu.bucket || ''} onChange={(e) => setQiniu({ ...qiniu, bucket: e.target.value })} placeholder="bucket" />
            <Input value={qiniu.img_path || ''} onChange={(e) => setQiniu({ ...qiniu, img_path: e.target.value })} placeholder="img_path" />
            <Button onClick={() => save(() => adminUpdateQiniuConfig(qiniu))}>保存七牛云配置</Button>
          </>
        ) : (
          <p className="text-sm text-muted-foreground">加载失败</p>
        )}
      </CollapsibleConfigCard>

      <CollapsibleConfigCard title="JWT配置">
        {jwt ? (
          <>
            <Input value={jwt.access_token_secret || ''} onChange={(e) => setJwt({ ...jwt, access_token_secret: e.target.value })} placeholder="access_token_secret" />
            <Input value={jwt.refresh_token_secret || ''} onChange={(e) => setJwt({ ...jwt, refresh_token_secret: e.target.value })} placeholder="refresh_token_secret" />
            <Input value={jwt.access_token_expiry_time || ''} onChange={(e) => setJwt({ ...jwt, access_token_expiry_time: e.target.value })} placeholder="access_token_expiry_time" />
            <Button onClick={() => save(() => adminUpdateJwtConfig(jwt))}>保存JWT配置</Button>
          </>
        ) : (
          <p className="text-sm text-muted-foreground">加载失败</p>
        )}
      </CollapsibleConfigCard>

      <CollapsibleConfigCard title="高德配置">
        {gaode ? (
          <>
            <label className="flex items-center gap-2 text-sm"><input type="checkbox" checked={!!gaode.enable} onChange={(e) => setGaode({ ...gaode, enable: e.target.checked })} /> 启用</label>
            <Input value={gaode.key || ''} onChange={(e) => setGaode({ ...gaode, key: e.target.value })} placeholder="key" />
            <Button onClick={() => save(() => adminUpdateGaodeConfig(gaode))}>保存高德配置</Button>
          </>
        ) : (
          <p className="text-sm text-muted-foreground">加载失败</p>
        )}
      </CollapsibleConfigCard>
    </div>
  );
}
