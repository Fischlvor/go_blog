'use client';

export const dynamic = 'force-dynamic';

import { useEffect, useState } from 'react';
import { getUserFeedbacks } from '@/lib/api/user/feedback';
import type { Feedback } from '@/lib/api/types';
import { formatDate } from '@/lib/date';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from '@/components/ui/table';

export default function UserFeedbackPage() {
  const [items, setItems] = useState<Feedback[]>([]);

  useEffect(() => {
    getUserFeedbacks().then(setItems).catch(() => setItems([]));
  }, []);

  return (
    <Card>
      <CardHeader><CardTitle>我的反馈</CardTitle></CardHeader>
      <CardContent>
        <Table>
          <TableHeader><TableRow><TableHead>时间</TableHead><TableHead>内容</TableHead><TableHead>回复</TableHead></TableRow></TableHeader>
          <TableBody>{items.map((it) => (
            <TableRow key={it.id}><TableCell>{formatDate(it.created_at)}</TableCell><TableCell>{it.content}</TableCell><TableCell>{it.reply || '-'}</TableCell></TableRow>
          ))}</TableBody>
        </Table>
      </CardContent>
    </Card>
  );
}
