/**
 * 共享 DTO 类型定义（对齐 server-blog-v2 响应结构）
 */

// ============ 分页 ============
export interface PageResult<T> {
  list: T[];
  current_page: number;
  page_size: number;
  total_items: number;
  total_pages: number;
}

// ============ 用户 ============
export interface User {
  id: number;
  uuid: string;
  nickname: string;
  username?: string;
  email: string;
  avatar: string;
  role_id: number;
  freeze: boolean;
  address?: string;
  signature?: string;
  created_at: string;
  last_login?: string | null;
}

export interface UserProfile {
  id: number;
  uuid: string;
  nickname: string;
  avatar: string;
  email?: string;
  signature?: string;
  role_id: number;
}

// ============ 文章 ============
export interface ArticleCategory {
  id: number;
  name: string;
  slug: string;
  article_count?: number;
  created_at?: string;
  updated_at?: string;
}

export interface ArticleTag {
  id: number;
  name: string;
  slug: string;
  article_count?: number;
  created_at?: string;
  updated_at?: string;
}

export interface ArticleAuthor {
  uuid: string;
  nickname: string;
  avatar: string;
}

export interface ArticleLikeInfo {
  liked?: boolean | null;
  likes: number;
}

// 文章摘要（列表用）
export interface Article {
  id: number;
  title: string;
  slug: string;
  excerpt: string;
  featured_image: string;
  author_uuid: string;
  author: ArticleAuthor;
  status: string;
  visibility: string;
  read_time: number;
  views: number;
  like: ArticleLikeInfo;
  is_featured: boolean;
  published_at: string | null;
  created_at: string;
  updated_at: string;
  category: ArticleCategory;
  tags: ArticleTag[];
}

// 文章详情（含正文）
export interface ArticleDetail extends Article {
  content: string;
  meta_title: string;
  meta_description: string;
}

export interface ArticleListQuery {
  page?: number;
  page_size?: number;
  keyword?: string;
  'filter.category_id'?: number;
  'filter.tag_id'?: number;
  sort_by?: string;
  order?: 'asc' | 'desc';
}

// ============ 评论 ============
export interface Comment {
  id: number;
  article_slug: string;
  parent_id: number | null;
  children?: Comment[];
  user: {
    uuid: string;
    nickname: string;
    avatar: string;
  };
  content: string;
  created_at: string;
  updated_at: string;
}

// ============ 友链 ============
export interface FriendLink {
  id: number;
  name: string;
  url?: string;
  link?: string;
  description: string;
  logo: string;
  status?: string;
  created_at?: string;
  updated_at?: string;
}

// ============ 反馈 ============
export interface Feedback {
  id: number;
  user_uuid?: string;
  type?: string;
  content: string;
  contact?: string;
  reply?: string | null;
  status: string;
  created_at: string;
  updated_at: string;
}

// ============ 认证 ============
export interface LoginResponse {
  access_token: string;
  user: UserProfile;
}

export interface SSOLoginUrlResponse {
  sso_login_url: string;
}
