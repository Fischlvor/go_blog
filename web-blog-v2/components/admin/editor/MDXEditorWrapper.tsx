'use client';

import { forwardRef, type ForwardedRef } from 'react';
import {
  MDXEditor,
  type MDXEditorMethods,
  type MDXEditorProps,
  headingsPlugin,
  listsPlugin,
  quotePlugin,
  thematicBreakPlugin,
  markdownShortcutPlugin,
  linkPlugin,
  linkDialogPlugin,
  imagePlugin,
  tablePlugin,
  codeBlockPlugin,
  codeMirrorPlugin,
  diffSourcePlugin,
  toolbarPlugin,
  UndoRedo,
  BoldItalicUnderlineToggles,
  BlockTypeSelect,
  CodeToggle,
  CreateLink,
  InsertImage,
  InsertTable,
  InsertThematicBreak,
  InsertCodeBlock,
  ListsToggle,
  Separator,
  DiffSourceToggleWrapper,
} from '@mdxeditor/editor';
import '@mdxeditor/editor/style.css';
import '@/styles/mdx-editor.css';

function EditorToolbar() {
  return (
    <DiffSourceToggleWrapper>
      <UndoRedo />
      <Separator />
      <BoldItalicUnderlineToggles />
      <CodeToggle />
      <Separator />
      <BlockTypeSelect />
      <Separator />
      <ListsToggle />
      <Separator />
      <CreateLink />
      <InsertImage />
      <InsertTable />
      <InsertThematicBreak />
      <InsertCodeBlock />
    </DiffSourceToggleWrapper>
  );
}

interface MDXEditorWrapperProps extends Omit<MDXEditorProps, 'ref'> {
  onImageUpload?: (file: File) => Promise<string>;
}

export const MDXEditorWrapper = forwardRef<MDXEditorMethods, MDXEditorWrapperProps>(
  ({ onImageUpload, ...props }, ref: ForwardedRef<MDXEditorMethods>) => {
    return (
      <MDXEditor
        plugins={[
          headingsPlugin(),
          listsPlugin(),
          quotePlugin(),
          thematicBreakPlugin(),
          markdownShortcutPlugin(),
          linkPlugin(),
          linkDialogPlugin(),
          imagePlugin({
            imageUploadHandler: onImageUpload
              ? async (image: File) => {
                  return await onImageUpload(image);
                }
              : undefined,
          }),
          tablePlugin(),
          codeBlockPlugin({ defaultCodeBlockLanguage: 'txt' }),
          codeMirrorPlugin({
            codeBlockLanguages: {
              txt: 'Text',
              js: 'JavaScript',
              ts: 'TypeScript',
              tsx: 'TSX',
              jsx: 'JSX',
              css: 'CSS',
              html: 'HTML',
              json: 'JSON',
              go: 'Go',
              python: 'Python',
              bash: 'Bash',
              sql: 'SQL',
              yaml: 'YAML',
              markdown: 'Markdown',
            },
          }),
          diffSourcePlugin({ viewMode: 'rich-text' }),
          toolbarPlugin({ toolbarContents: () => <EditorToolbar /> }),
        ]}
        {...props}
        ref={ref}
      />
    );
  }
);

MDXEditorWrapper.displayName = 'MDXEditorWrapper';
