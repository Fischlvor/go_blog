'use client';

export const dynamic = 'force-dynamic';

import { useEffect, useState } from 'react';
import { toast } from 'sonner';
import { getUserInfo, resetPassword, updateUserInfo } from '@/lib/client-api/user/user';
import type { User } from '@/lib/client-api/types';
import { Button } from '@/components/ui/button';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Dialog, DialogContent, DialogFooter, DialogHeader, DialogTitle } from '@/components/ui/dialog';
import { Input } from '@/components/ui/input';
import { Textarea } from '@/components/ui/textarea';

export default function UserInfoPage() {
  const [user, setUser] = useState<User | null>(null);
  const [nickname, setNickname] = useState('');
  const [avatar, setAvatar] = useState('');
  const [address, setAddress] = useState('');
  const [signature, setSignature] = useState('');
  const [pwdOpen, setPwdOpen] = useState(false);
  const [password, setPassword] = useState('');
  const [newPassword, setNewPassword] = useState('');

  const load = async () => {
    try {
      const res = await getUserInfo();
      setUser(res);
      setNickname(res.nickname || '');
      setAvatar(res.avatar || '');
      setAddress(res.address || '');
      setSignature(res.signature || '');
    } catch {
      toast.error('获取用户信息失败');
    }
  };

  useEffect(() => {
    getUserInfo()
      .then((res) => {
        setUser(res);
        setNickname(res.nickname || '');
        setAvatar(res.avatar || '');
        setAddress(res.address || '');
        setSignature(res.signature || '');
      })
      .catch(() => {
        toast.error('获取用户信息失败');
      });
  }, []);

  const onSave = async () => {
    try {
      await updateUserInfo({ nickname, avatar, address, signature });
      toast.success('更新成功');
      await load();
    } catch {
      toast.error('更新失败');
    }
  };

  const onResetPassword = async () => {
    if (!password || !newPassword) {
      toast.error('请输入原密码和新密码');
      return;
    }
    try {
      await resetPassword(password, newPassword);
      toast.success('修改密码成功');
      setPwdOpen(false);
      setPassword('');
      setNewPassword('');
    } catch {
      toast.error('修改密码失败');
    }
  };

  return (
    <div className="grid gap-4 lg:grid-cols-2">
      <Card>
        <CardHeader><CardTitle>用户信息</CardTitle></CardHeader>
        <CardContent className="space-y-3">
          <Input value={avatar} onChange={(e) => setAvatar(e.target.value)} placeholder="头像URL" />
          <Input value={nickname} onChange={(e) => setNickname(e.target.value)} placeholder="昵称" />
          <Input value={address} onChange={(e) => setAddress(e.target.value)} placeholder="地址" />
          <Textarea value={signature} onChange={(e) => setSignature(e.target.value)} placeholder="签名" rows={3} />
          <div className="text-sm text-muted-foreground">UUID：{user?.uuid || '-'}</div>
          <Button onClick={onSave}>保存信息</Button>
        </CardContent>
      </Card>

      <Card>
        <CardHeader><CardTitle>操作</CardTitle></CardHeader>
        <CardContent className="space-y-3">
          <Button variant="outline" onClick={() => setPwdOpen(true)}>修改密码</Button>
        </CardContent>
      </Card>

      <Dialog open={pwdOpen} onOpenChange={setPwdOpen}>
        <DialogContent>
          <DialogHeader><DialogTitle>修改密码</DialogTitle></DialogHeader>
          <div className="space-y-3">
            <Input value={password} onChange={(e) => setPassword(e.target.value)} type="password" placeholder="原密码" />
            <Input value={newPassword} onChange={(e) => setNewPassword(e.target.value)} type="password" placeholder="新密码" />
          </div>
          <DialogFooter>
            <Button variant="outline" onClick={() => setPwdOpen(false)}>取消</Button>
            <Button onClick={onResetPassword}>确定</Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>
    </div>
  );
}
