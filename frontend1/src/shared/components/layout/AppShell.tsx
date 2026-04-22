import { PropsWithChildren, ReactNode } from 'react';
import { Header } from '@/shared/components/layout/Header';

type AppShellProps = PropsWithChildren<{
  hero?: ReactNode;
}>;

export function AppShell({ hero, children }: AppShellProps) {
  return (
    <div className="min-h-screen">
      <Header />
      {hero}
      <main className="mx-auto max-w-6xl px-4 py-8 md:px-6">{children}</main>
    </div>
  );
}
