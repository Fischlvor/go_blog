'use client';

import Link from 'next/link';
import { useRouter } from 'next/navigation';
import { LogOut, User } from 'lucide-react';
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu';
import { Avatar, AvatarFallback, AvatarImage } from '@/components/ui/avatar';
import { Button } from '@/components/ui/button';
import { useUserAuth } from '@/context/user-auth';
import { toast } from 'sonner';

interface UserMenuProps {
  mobile?: boolean;
}

export function UserMenu({ mobile }: UserMenuProps) {
  const { user, logout } = useUserAuth();
  const router = useRouter();

  const handleLogout = async () => {
    await logout();
    toast.success('已退出登录');
    router.push('/');
  };

  if (!user) return null;

  if (mobile) {
    return (
      <div className="space-y-1">
        <div className="flex items-center gap-3 px-4 py-2">
          <Avatar className="h-8 w-8">
            <AvatarImage src={user.avatar} alt={user.nickname} />
            <AvatarFallback>{user.nickname?.[0]?.toUpperCase()}</AvatarFallback>
          </Avatar>
          <span className="text-sm font-medium">{user.nickname}</span>
        </div>
        <Button
          variant="ghost"
          className="w-full justify-start gap-2 text-destructive"
          onClick={handleLogout}
        >
          <LogOut className="h-4 w-4" /> 退出登录
        </Button>
      </div>
    );
  }

  return (
    <DropdownMenu>
      <DropdownMenuTrigger className="rounded-full outline-none cursor-pointer">
        <Avatar className="h-8 w-8">
          <AvatarImage src={user.avatar} alt={user.nickname} />
          <AvatarFallback>{user.nickname?.[0]?.toUpperCase()}</AvatarFallback>
        </Avatar>
      </DropdownMenuTrigger>
      <DropdownMenuContent align="end" className="w-48">
        <div className="px-3 py-2">
          <p className="text-sm font-medium">{user.nickname}</p>
          <p className="text-xs text-muted-foreground truncate">{user.email}</p>
        </div>
        <DropdownMenuSeparator />
        <DropdownMenuItem
          className="flex items-center gap-2 cursor-pointer"
          onClick={() => router.push('/dashboard')}
        >
          <User className="h-4 w-4" /> 个人中心
        </DropdownMenuItem>
        <DropdownMenuSeparator />
        <DropdownMenuItem
          className="text-destructive focus:text-destructive flex items-center gap-2 cursor-pointer"
          onClick={handleLogout}
        >
          <LogOut className="h-4 w-4" /> 退出登录
        </DropdownMenuItem>
      </DropdownMenuContent>
    </DropdownMenu>
  );
}
