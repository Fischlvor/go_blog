import 'server-only';

import type { Website } from '@/lib/client-api/public/website';
import type { Article, ArticleCategory, ArticleDetail, Comment, FriendLink, PageResult } from '@/lib/client-api/types';

const EMPTY_FIELD = { value: '', setting_key: '' };

export const FALLBACK_SITE: Website = {
  avatar: EMPTY_FIELD,
  title: { value: '博客', setting_key: 'website.title' },
  description: EMPTY_FIELD,
  profile_intro: EMPTY_FIELD,
  tech_stack: EMPTY_FIELD,
  work_experiences: EMPTY_FIELD,
  version: EMPTY_FIELD,
  created_at: EMPTY_FIELD,
  icp_filing: EMPTY_FIELD,
  bilibili_url: EMPTY_FIELD,
  github_url: EMPTY_FIELD,
  steam_url: { value: '', setting_key: 'profile.steam_url' },
  name: { value: '博客', setting_key: 'profile.name' },
  job: EMPTY_FIELD,
  address: EMPTY_FIELD,
  email: EMPTY_FIELD,
};

export const FALLBACK_ARTICLE_PAGE: PageResult<Article> = {
  list: [],
  total_items: 0,
  total_pages: 0,
  current_page: 1,
  page_size: 0,
};

export const FALLBACK_ARTICLE_DETAIL: ArticleDetail = {
  id: 0,
  title: '文章暂不可用',
  slug: '',
  excerpt: '',
  featured_image: '',
  author_uuid: '',
  author: {
    uuid: '',
    nickname: '',
    avatar: '',
  },
  status: '',
  visibility: '',
  read_time: 0,
  views: 0,
  like: {
    liked: false,
    likes: 0,
  },
  is_featured: false,
  published_at: '',
  created_at: '',
  updated_at: '',
  category: {
    id: 0,
    name: '',
    slug: '',
  },
  tags: [],
  content: '',
  meta_title: '',
  meta_description: '',
};

export const FALLBACK_CATEGORIES: ArticleCategory[] = [];
export const FALLBACK_FRIEND_LINKS: FriendLink[] = [];
export const FALLBACK_COMMENTS: Comment[] = [];
