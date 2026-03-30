'use client';

export const dynamic = 'force-dynamic';

import Link from 'next/link';
import { useEffect, useState } from 'react';
import { getUserComments } from '@/lib/api/public/comment';
import type { Comment } from '@/lib/api/types';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from '@/components/ui/table';

export default function UserCommentPage() {
  const [items, setItems] = useState<Comment[]>([]);

  useEffect(() => {
    getUserComments()
      .then((res) => {
        if (Array.isArray(res)) {
          setItems(res);
          return;
        }
        const maybeList = (res as { list?: Comment[] })?.list;
        setItems(Array.isArray(maybeList) ? maybeList : []);
      })
      .catch(() => setItems([]));
  }, []);

  return (
    <Card>
      <CardHeader><CardTitle>我的评论</CardTitle></CardHeader>
      <CardContent>
        <Table>
          <TableHeader><TableRow><TableHead>文章Slug</TableHead><TableHead>内容</TableHead><TableHead>时间</TableHead></TableRow></TableHeader>
          <TableBody>{items.map((it) => (
            <TableRow key={it.id}><TableCell><Link href={`/article/${it.article_slug}`} className="text-primary hover:underline">{it.article_slug}</Link></TableCell><TableCell className="max-w-[560px] truncate">{it.content}</TableCell><TableCell>{it.created_at}</TableCell></TableRow>
          ))}</TableBody>
        </Table>
      </CardContent>
    </Card>
  );
}
