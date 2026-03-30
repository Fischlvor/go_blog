'use client';

import { useRef, useState } from 'react';
import { createPortal } from 'react-dom';

interface ImageThumbCellProps {
  src?: string | null;
  alt?: string;
  boxWidth?: number;
  boxHeight?: number;
  fallbackText?: string;
  previewWidth?: number;
  previewHeight?: number;
  responsive?: boolean;
}

export function ImageThumbCell({
  src,
  alt = 'thumb',
  boxWidth = 192,
  boxHeight = 108,
  fallbackText = '无图',
  previewWidth = 320,
  previewHeight = 220,
  responsive = false,
}: ImageThumbCellProps) {
  const anchorRef = useRef<HTMLDivElement | null>(null);
  const [open, setOpen] = useState(false);
  const [position, setPosition] = useState<{ left: number; top: number } | null>(null);

  const openPreview = () => {
    if (!anchorRef.current || typeof window === 'undefined') return;

    const rect = anchorRef.current.getBoundingClientRect();
    const gap = 8;

    let left = rect.right + gap;
    if (left + previewWidth > window.innerWidth - 12) {
      left = rect.left - previewWidth - gap;
    }
    left = Math.max(12, left);

    const enoughBelow = rect.top + previewHeight <= window.innerHeight - 12;
    let top = enoughBelow ? rect.top : rect.bottom - previewHeight;
    top = Math.max(12, Math.min(top, window.innerHeight - previewHeight - 12));

    setPosition({ left, top });
    setOpen(true);
  };

  return (
    <>
      <div
        ref={anchorRef}
        className="relative flex h-full w-full max-h-full max-w-full items-center justify-center"
        onMouseEnter={openPreview}
        onMouseLeave={() => setOpen(false)}
      >
        <div
          className="flex h-full w-full items-center justify-center overflow-hidden rounded-lg"
          style={responsive ? { maxWidth: '100%', maxHeight: '100%' } : { width: boxWidth, height: boxHeight, maxWidth: '100%', maxHeight: '100%' }}
        >
          {src ? (
            <img
              src={src}
              alt={alt}
              className="h-full w-full object-contain"
            />
          ) : (
            <span className="text-[10px] text-muted-foreground">{fallbackText}</span>
          )}
        </div>
      </div>

      {src && open && position && typeof document !== 'undefined'
        ? createPortal(
          <div
            className="pointer-events-none fixed z-[9999] overflow-hidden rounded-lg border border-border bg-background shadow-2xl"
            style={{ width: previewWidth, height: previewHeight, left: position.left, top: position.top }}
          >
            <img src={src} alt={alt} className="h-full w-full object-contain" />
          </div>,
          document.body
        )
        : null}
    </>
  );
}
