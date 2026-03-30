'use client';

import { useRef, useState } from 'react';
import { createPortal } from 'react-dom';
import { Avatar, AvatarFallback, AvatarImage } from '@/components/ui/avatar';

export interface UserCellPopoverInfo {
  uuid?: string;
  nickname?: string;
  username?: string;
  avatar?: string;
  email?: string;
  address?: string;
  signature?: string;
}

interface UserCellPopoverProps {
  user?: UserCellPopoverInfo | null;
  fallbackText?: string;
}

function getDisplayName(user?: UserCellPopoverInfo | null) {
  if (!user) return '';
  return user.nickname || user.username || user.uuid || '未知用户';
}

export function UserCellPopover({ user, fallbackText = '-' }: UserCellPopoverProps) {
  const [open, setOpen] = useState(false);
  const [position, setPosition] = useState<{ left: number; top: number } | null>(null);
  const anchorRef = useRef<HTMLDivElement | null>(null);

  if (!user) {
    return <span className="text-xs text-muted-foreground">{fallbackText}</span>;
  }

  const displayName = getDisplayName(user);

  const openPopover = () => {
    if (!anchorRef.current || typeof window === 'undefined') return;
    const rect = anchorRef.current.getBoundingClientRect();
    const panelWidth = 320;
    const panelHeight = 180;
    const viewportWidth = window.innerWidth;
    const viewportHeight = window.innerHeight;

    let left = rect.left;
    if (left + panelWidth > viewportWidth - 12) {
      left = Math.max(12, viewportWidth - panelWidth - 12);
    }

    let top = rect.bottom + 8;
    if (top + panelHeight > viewportHeight - 12) {
      top = Math.max(12, rect.top - panelHeight - 8);
    }

    setPosition({ left, top });
    setOpen(true);
  };

  const closePopover = () => setOpen(false);

  return (
    <>
      <div ref={anchorRef} className="inline-flex items-center gap-2" onMouseEnter={openPopover} onMouseLeave={closePopover}>
        <Avatar size="sm">
          {user.avatar ? <AvatarImage src={user.avatar} alt={displayName} /> : null}
          <AvatarFallback>{displayName.slice(0, 1)}</AvatarFallback>
        </Avatar>
        <span className="max-w-[120px] truncate text-sm">{displayName}</span>
      </div>

      {open && position && typeof document !== 'undefined'
        ? createPortal(
          <div className="fixed z-[9999] w-80 rounded-xl border border-border bg-background p-3 text-xs shadow-2xl" style={{ left: position.left, top: position.top }} onMouseEnter={() => setOpen(true)} onMouseLeave={closePopover}>
            <div className="flex items-center gap-3">
              <Avatar>
                {user.avatar ? <AvatarImage src={user.avatar} alt={displayName} /> : null}
                <AvatarFallback>{displayName.slice(0, 1)}</AvatarFallback>
              </Avatar>
              <div className="min-w-0">
                <div className="truncate text-sm font-medium">{displayName}</div>
                <div className="truncate text-muted-foreground">{user.uuid || '-'}</div>
              </div>
            </div>

            <div className="mt-3 space-y-1 text-muted-foreground">
              <div className="truncate">邮箱：{user.email || '-'}</div>
              <div className="truncate">地址：{user.address || '-'}</div>
              <div className="line-clamp-2">签名：{user.signature || '-'}</div>
            </div>
          </div>,
          document.body
        )
        : null}
    </>
  );
}
