'use client';

import { useEffect, useRef, useState } from 'react';
import ReactMarkdown from 'react-markdown';
import remarkGfm from 'remark-gfm';
import { cn } from '@/lib/utils';
import type { Components } from 'react-markdown';

// ── Mermaid 图表 ──────────────────────────────────────────────
function MermaidBlock({ code }: { code: string }) {
  const ref = useRef<HTMLDivElement>(null);
  const [svg, setSvg] = useState<string | null>(null);
  const [error, setError] = useState(false);

  useEffect(() => {
    let cancelled = false;
    (async () => {
      try {
        const mermaid = (await import('mermaid')).default;
        const isDark = document.documentElement.classList.contains('dark');
        mermaid.initialize({
          startOnLoad: false,
          theme: isDark ? 'dark' : 'neutral',
          securityLevel: 'loose',
        });
        const id = `mermaid-${Math.random().toString(36).slice(2)}`;
        const { svg: rendered } = await mermaid.render(id, code);
        if (!cancelled) setSvg(rendered);
      } catch {
        if (!cancelled) setError(true);
      }
    })();
    return () => { cancelled = true; };
  }, [code]);

  if (error) return (
    <pre className="rounded-xl border bg-muted/50 p-4 text-xs text-muted-foreground overflow-x-auto">{code}</pre>
  );
  if (!svg) return (
    <div className="flex items-center justify-center h-24 rounded-xl border bg-muted/30 text-xs text-muted-foreground">渲染中...</div>
  );
  return (
    <div ref={ref} className="flex justify-center my-4 overflow-x-auto"
      dangerouslySetInnerHTML={{ __html: svg }} />
  );
}

// ── 代码块（含复制按钮）──────────────────────────────────────
function CodeBlock({ className, children }: { className?: string; children?: React.ReactNode }) {
  const lang = /language-(\S+)/.exec(className || '')?.[1] ?? '';
  const code = String(children).replace(/\n$/, '');
  const [copied, setCopied] = useState(false);

  if (lang === 'mermaid') return <MermaidBlock code={code} />;

  const handleCopy = async () => {
    try {
      await navigator.clipboard.writeText(code);
      setCopied(true);
      setTimeout(() => setCopied(false), 1500);
    } catch { /* ignore */ }
  };

  return (
    <div className="relative group">
      <button
        type="button"
        onClick={handleCopy}
        className="absolute right-3 top-3 z-10 opacity-0 group-hover:opacity-100 transition-opacity rounded-md border border-border/60 bg-background/80 backdrop-blur-sm px-2 py-1 text-xs text-muted-foreground hover:text-foreground cursor-pointer"
      >
        {copied ? '已复制' : '复制'}
      </button>
      <pre className={cn('overflow-x-auto rounded-xl border bg-zinc-100 dark:bg-muted/50 p-4 text-sm text-zinc-800 dark:text-zinc-200', className)}>
        <code className={className}>{code}</code>
      </pre>
    </div>
  );
}

// ── 标题（自动加 id）─────────────────────────────────────────
function makeHeading(Tag: 'h1' | 'h2' | 'h3' | 'h4' | 'h5' | 'h6') {
  return function Heading({ children, ...props }: React.HTMLAttributes<HTMLHeadingElement>) {
    const text = String(children ?? '');
    const id = text.toLowerCase().replace(/\s+/g, '-').replace(/[^\w\u4e00-\u9fa5-]/g, '');
    return <Tag id={id} {...props}>{children}</Tag>;
  };
}

const COMPONENTS: Components = {
  h1: makeHeading('h1'),
  h2: makeHeading('h2'),
  h3: makeHeading('h3'),
  h4: makeHeading('h4'),
  h5: makeHeading('h5'),
  h6: makeHeading('h6'),
  // pre 包裹代码块（含无语言标识的纯文本块）
  pre({ children }) {
    // 提取 code 子元素
    const child = children as React.ReactElement<{ className?: string; children?: React.ReactNode }>;
    const className = child?.props?.className ?? '';
    const code = String(child?.props?.children ?? '').replace(/\n$/, '');
    return <CodeBlock className={className}>{code}</CodeBlock>;
  },
  // 行内代码（父节点不是 pre，react-markdown 直接调用）
  code({ className, children, ...props }) {
    if (className?.startsWith('language-')) {
      // fallback：直接渲染（正常情况由 pre 处理）
      return <CodeBlock className={className}>{String(children)}</CodeBlock>;
    }
    return (
      <code
        className="font-mono text-sm bg-muted px-1.5 py-0.5 rounded text-primary/90 before:content-none after:content-none"
        {...props}
      >
        {children}
      </code>
    );
  },
  // 图片
  img({ src, alt }) {
    return (
      <img
        src={src} alt={alt ?? ''}
        className="rounded-xl border max-w-full h-auto my-4"
        loading="lazy"
      />
    );
  },
  // 表格
  table({ children }) {
    return (
      <div className="overflow-x-auto my-4">
        <table className="min-w-full divide-y divide-border">{children}</table>
      </div>
    );
  },
  // 引用
  blockquote({ children }) {
    return (
      <blockquote className="border-l-4 border-primary/40 pl-4 italic text-muted-foreground my-4">
        {children}
      </blockquote>
    );
  },
  // 链接
  a({ href, children }) {
    const external = href?.startsWith('http');
    return (
      <a
        href={href}
        target={external ? '_blank' : undefined}
        rel={external ? 'noopener noreferrer' : undefined}
        className="text-primary underline-offset-4 hover:underline"
      >
        {children}
      </a>
    );
  },
};

// ── 主组件 ───────────────────────────────────────────────────
interface MarkdownContentProps {
  content: string;
  className?: string;
}

export function MarkdownContent({ content, className }: MarkdownContentProps) {
  return (
    <div
      className={cn(
        // prose 基础
        'prose prose-neutral dark:prose-invert max-w-none',
        // 标题间距
        'prose-headings:scroll-mt-24 prose-headings:font-bold prose-headings:tracking-tight',
        // h1
        'prose-h1:text-3xl prose-h1:border-b prose-h1:border-border prose-h1:pb-3',
        // h2
        'prose-h2:text-2xl prose-h2:border-b prose-h2:border-border/50 prose-h2:pb-2',
        // h3
        'prose-h3:text-xl',
        // 正文
        'prose-p:leading-7 prose-p:text-foreground/90',
        // 列表
        'prose-li:text-foreground/90',
        // 链接（由自定义组件处理，这里只做 fallback）
        'prose-a:text-primary prose-a:no-underline hover:prose-a:underline',
        // 行内代码（由自定义组件处理）
        'prose-code:before:content-none prose-code:after:content-none',
        // hr
        'prose-hr:border-border',
        // 表格
        'prose-th:bg-muted/60 prose-td:border-border prose-th:border-border',
        // strong
        'prose-strong:text-foreground prose-strong:font-semibold',
        className
      )}
    >
      <ReactMarkdown remarkPlugins={[remarkGfm]} components={COMPONENTS}>
        {content}
      </ReactMarkdown>
    </div>
  );
}
