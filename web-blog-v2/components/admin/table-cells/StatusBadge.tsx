'use client';

import { Badge } from '@/components/ui/badge';

type StatusTone = 'default' | 'secondary' | 'destructive' | 'outline';

interface StatusBadgeProps {
  text: string;
  tone?: StatusTone;
}

export function StatusBadge({ text, tone = 'secondary' }: StatusBadgeProps) {
  return <Badge variant={tone}>{text}</Badge>;
}

export function boolStatusBadge(value: boolean, trueText = '是', falseText = '否') {
  return <StatusBadge text={value ? trueText : falseText} tone={value ? 'destructive' : 'secondary'} />;
}
