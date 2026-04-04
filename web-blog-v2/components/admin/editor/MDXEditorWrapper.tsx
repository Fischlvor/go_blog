'use client';

import {
  forwardRef,
  useEffect,
  useImperativeHandle,
  useRef,
  type ForwardedRef,
} from 'react';
import Vditor from 'vditor';
import 'vditor/dist/index.css';

export interface VditorEditorMethods {
  setMarkdown: (markdown: string) => void;
  getMarkdown: () => string;
  focus: () => void;
}

interface MDXEditorWrapperProps {
  markdown: string;
  onChange: (markdown: string) => void;
  onImageUpload?: (file: File) => Promise<string>;
  placeholder?: string;
}

const TOOLBAR: Array<string> = [
  'undo',
  'redo',
  '|',
  'headings',
  'bold',
  'italic',
  'strike',
  'quote',
  'line',
  '|',
  'list',
  'ordered-list',
  'check',
  'outdent',
  'indent',
  '|',
  'code',
  'inline-code',
  'insert-after',
  'insert-before',
  'table',
  'link',
  'upload',
  '|',
  'fullscreen',
  'outline',
  'help',
];

export const MDXEditorWrapper = forwardRef<VditorEditorMethods, MDXEditorWrapperProps>(
  ({ markdown, onChange, onImageUpload, placeholder }, ref: ForwardedRef<VditorEditorMethods>) => {
    const containerRef = useRef<HTMLDivElement | null>(null);
    const vditorRef = useRef<Vditor | null>(null);
    const currentValueRef = useRef(markdown);
    const onChangeRef = useRef(onChange);
    const isReadyRef = useRef(false);
    const hasInitializedRef = useRef(false);

    useEffect(() => {
      onChangeRef.current = onChange;
    }, [onChange]);

    useImperativeHandle(ref, () => ({
      setMarkdown: (value: string) => {
        currentValueRef.current = value;
        vditorRef.current?.setValue(value);
      },
      getMarkdown: () => vditorRef.current?.getValue() ?? currentValueRef.current,
      focus: () => {
        vditorRef.current?.focus();
      },
    }), []);

    useEffect(() => {
      if (!containerRef.current || vditorRef.current || hasInitializedRef.current) return;

      const isDark = document.documentElement.classList.contains('dark');
      hasInitializedRef.current = true;

      const instance = new Vditor(containerRef.current, {
        mode: 'sv',
        height: '100%',
        minHeight: 640,
        cache: { enable: false },
        toolbarConfig: { pin: true },
        toolbar: TOOLBAR,
        placeholder: placeholder || '开始撰写文章内容...',
        value: currentValueRef.current,
        theme: isDark ? 'dark' : 'classic',
        icon: 'material',
        preview: {
          mode: 'both',
          maxWidth: Number.MAX_SAFE_INTEGER,
          theme: { current: isDark ? 'dark' : 'light' },
          hljs: { style: isDark ? 'native' : 'github' },
          markdown: {
            toc: true,
            footnotes: true,
            codeBlockPreview: true,
            mathBlockPreview: true,
          },
        },
        counter: {
          enable: true,
          type: 'markdown',
        },
        upload: onImageUpload
          ? {
              accept: 'image/*',
              multiple: false,
              handler: async (files) => {
                try {
                  const succMap: Record<string, string> = {};

                  for (const file of files) {
                    const url = await onImageUpload(file);
                    succMap[file.name] = url;
                  }

                  return JSON.stringify({
                    msg: '',
                    code: 0,
                    data: {
                      errFiles: [],
                      succMap,
                    },
                  });
                } catch (error) {
                  return (error as Error)?.message || '上传失败';
                }
              },
            }
          : undefined,
        after: () => {
          isReadyRef.current = true;
          vditorRef.current?.setValue(currentValueRef.current);
        },
        input: (value) => {
          currentValueRef.current = value;
          onChangeRef.current(value);
        },
      });

      vditorRef.current = instance;

      const observer = new MutationObserver(() => {
        const dark = document.documentElement.classList.contains('dark');
        instance.setTheme(dark ? 'dark' : 'classic', dark ? 'dark' : 'light', dark ? 'native' : 'github');
      });

      observer.observe(document.documentElement, {
        attributes: true,
        attributeFilter: ['class'],
      });

      return () => {
        observer.disconnect();
        isReadyRef.current = false;
        hasInitializedRef.current = false;
        vditorRef.current?.destroy();
        vditorRef.current = null;
      };
    }, [onImageUpload, placeholder]);

    useEffect(() => {
      currentValueRef.current = markdown;

      if (!isReadyRef.current || !vditorRef.current) return;

      if (vditorRef.current.getValue() !== markdown) {
        vditorRef.current.setValue(markdown);
      }
    }, [markdown]);

    return <div ref={containerRef} className="h-full [&_.vditor]:h-full [&_.vditor-content]:max-w-none" />;
  }
);

MDXEditorWrapper.displayName = 'MDXEditorWrapper';

