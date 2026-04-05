import 'server-only';

import type { Website } from '@/lib/client-api/public/website';
import type { Article, ArticleCategory, ArticleDetail, Comment, FriendLink, PageResult } from '@/lib/client-api/types';

export const FALLBACK_SITE: Website = {
  avatar: '',
  logo: '',
  full_logo: '',
  title: '博客',
  slogan: '',
  slogan_en: '',
  description: '',
  version: '',
  created_at: '',
  icp_filing: '',
  public_security_filing: '',
  bilibili_url: '',
  github_url: '',
  steam_url: '',
  name: '博客',
  job: '',
  address: '',
  email: '',
  qq_image: '',
  wechat_image: '',
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
