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
    cover: string;
    title: string;
    category: string;
    tags: string[];
    abstract: string;
    content: string;
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
    id: string;
    cover: string;
    title: string;
    category: string;
    tags: string[];
    abstract: string;
    content: string;
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
