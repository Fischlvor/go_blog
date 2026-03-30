import Link from 'next/link';
import { redirect } from 'next/navigation';
import { cookies } from 'next/headers';
import { AdminSidebar } from '@/components/admin/layout/Sidebar';
import { AdminHeader } from '@/components/admin/layout/Header';

async function checkAdminAuth(): Promise<boolean> {
  // Admin 通过 JWT token 鉴权，token 存在 localStorage（客户端）
  // 服务端 layout 只做基本骨架，具体鉴权由客户端组件完成
  return true;
}

export default async function DashboardLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <div className="flex h-screen overflow-hidden bg-background">
      <AdminSidebar />
      <div className="flex flex-col flex-1 overflow-hidden">
        <AdminHeader />
        <main className="flex-1 overflow-y-auto p-6">
          {children}
        </main>
      </div>
    </div>
  );
}
