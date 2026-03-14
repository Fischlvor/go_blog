import type {Hit, PageInfo, PageResult} from "@/api/common";
import type {ApiResponse} from "@/utils/request";
import service, { adminService } from "@/utils/request";

export interface ArticleAuthor {
    uuid: string;
    nickname: string;
    avatar: string;
}

export interface ArticleLikeInfo {
    liked?: boolean;
    likes: number;
}

export interface ArticleCategory {
    id: number;
    name: string;
    slug: string;
}

export interface ArticleTagItem {
    id: number;
    name: string;
    slug: string;
}

export interface Article {
    id: number;
    title: string;
    slug: string;
    excerpt: string;
    abstract?: string;
    cover?: string;
    featured_image: string;
    author_uuid: string;
    author: ArticleAuthor;
    status: string;
    visibility: string; // 'public' | 'private'
    read_time: number;
    views: number;
    like: ArticleLikeInfo;
    is_featured: boolean;
    published_at: string | null;
    created_at: string;
    updated_at: string;
    category: ArticleCategory;
    tags: ArticleTagItem[];
    content?: string;
    meta_title?: string;
    meta_description?: string;
}

export interface ArticleLikeRequest {
    slug: string;
}

export const articleLike = (slug: string): Promise<ApiResponse<undefined>> => {
    return service({
        url: `/article/${slug}/like`,
        method: 'post',
    })
}

export const articleUnlike = (slug: string): Promise<ApiResponse<undefined>> => {
    return service({
        url: `/article/${slug}/like`,
        method: 'delete',
    })
}

export const articleLikesList = (data: PageInfo): Promise<ApiResponse<PageResult<Article>>> => {
    return service({
        url: '/article/likes',
        method: 'get',
        params: data
    })
}

export interface ArticleCreateRequest {
    title: string;        // 标题
    slug: string;         // 文章 slug
    content: string;      // 内容
    excerpt: string;      // 摘要
    featured_image: string; // 封面图
    category_id: number;  // 分类 ID
    tag_ids: number[];    // 标签 ID 列表
    status: string;       // 状态：draft, published
    visibility: string;   // 可见性：public, private
    is_featured: boolean; // 是否精选
}

export const articleCreate = (data: ArticleCreateRequest): Promise<ApiResponse<undefined>> => {
    return adminService({
        url: '/article/create',
        method: 'post',
        data: data
    });
}

export interface ArticleDeleteRequest {
    ids: string[];
}

export const articleDelete = (data: ArticleDeleteRequest): Promise<ApiResponse<undefined>> => {
    return adminService({
        url: '/article/delete',
        method: 'delete',
        data: data
    });
}

export interface ArticleUpdateRequest {
    slug: string;         // 文章 slug（作为标识）
    title: string;        // 标题
    content: string;      // 内容
    excerpt: string;      // 摘要
    featured_image: string; // 封面图
    category_id: number;  // 分类 ID
    tag_ids: number[];    // 标签 ID 列表
    status: string;       // 状态：draft, published
    visibility: string;   // 可见性：public, private
    is_featured: boolean; // 是否精选
}

export const articleUpdate = (data: ArticleUpdateRequest): Promise<ApiResponse<undefined>> => {
    return adminService({
        url: '/article/update',
        method: 'put',
        data: data
    });
}

export interface ArticleListRequest extends PageInfo {
    title: string | null;
    category: string | null;
    abstract: string | null;
}

export const articleList = (data: ArticleListRequest): Promise<ApiResponse<PageResult<Article>>> => {
    return adminService({
        url: '/article/list',
        method: 'get',
        params: data,
    });
}

export const articleInfoByID = (id: string): Promise<ApiResponse<Article>> => {
    return service({
        url: '/article/'+id,
        method: 'get',
    });
}

export interface ArticleSearchRequest extends PageInfo {
    query: string;
    category: string;
    tag: string;
    sort: string;
    order: string;
}

export const articleSearch = (data: ArticleSearchRequest): Promise<ApiResponse<PageResult<Article>>> => {
    return service({
        url: '/article/search',
        method: 'get',
        params: data,
    });
}

export interface CategoryDetail {
    id: number;
    name: string;
    slug: string;
    article_count: number;
    created_at: string;
    updated_at: string;
}

export const articleCategory = (): Promise<ApiResponse<CategoryDetail[]>> => {
    return service({
        url: '/article/category',
        method: 'get',
    });
}

export interface TagDetail {
    id: number;
    name: string;
    slug: string;
    article_count: number;
    created_at: string;
    updated_at: string;
}

export const articleTags = (): Promise<ApiResponse<TagDetail[]>> => {
    return service({
        url: '/article/tags',
        method: 'get',
    });
}

// 创建分类
export interface CreateCategoryRequest {
    name: string;
    slug: string;
}

export const createCategory = (data: CreateCategoryRequest): Promise<ApiResponse<{ id: number }>> => {
    return adminService({
        url: '/category/create',
        method: 'post',
        data,
    });
}

// 创建标签
export interface CreateTagRequest {
    name: string;
    slug: string;
}

export const createTag = (data: CreateTagRequest): Promise<ApiResponse<{ id: number }>> => {
    return adminService({
        url: '/tag/create',
        method: 'post',
        data,
    });
}
