import type {Model, PageInfo, PageResult} from "@/api/common";
import type {ApiResponse} from "@/utils/request";
import service, { streamRequest, type StreamRequestConfig } from "@/utils/request";

// 表情接口
export interface Emoji extends Model {
    key: string;
    filename: string;
    group_name: string;
    sprite_group: number;
    sprite_position_x: number;
    sprite_position_y: number;
    file_size: number;
    cdn_url: string;
    upload_time: string;
    status: number;
    created_by: string;
}

// 表情组接口
export interface EmojiGroup extends Model {
    group_name: string;
    group_key: string;
    description: string;
    sort_order: number;
    emoji_count: number;
    status: number;
    created_by: string;
}

// 雪碧图接口
export interface EmojiSprite extends Model {
    sprite_group: number;
    filename: string;
    cdn_url: string;
    width: number;
    height: number;
    emoji_count: number;
    file_size: number;
    status: number;
}

// 任务接口
export interface EmojiTask extends Model {
    task_type: string;
    status: string;
    progress: number;
    message: string;
    params?: any;
    result?: any;
    created_by: string;
    started_at?: string;
    completed_at?: string;
}

// 前端配置接口
export interface EmojiConfig {
    version: string;
    total_emojis: number;
    sprites: EmojiSprite[] | null;
    mapping?: Record<string, any> | null; // 已废弃，保留用于兼容
    updated_at: string;
}

// 请求接口
export interface EmojiListRequest extends PageInfo {
    keyword?: string;
    group_key?: string;
    sprite_group?: number;
}

export interface EmojiGroupCreateRequest {
    group_name: string;
    description?: string;
    sort_order?: number;
}

export interface EmojiGroupUpdateRequest {
    group_name: string;
    description?: string;
    sort_order?: number;
}

// EmojiUploadRequest 接口已不再需要，现在直接使用FormData

// API函数

// 表情列表
export const getEmojiList = (params: EmojiListRequest): Promise<ApiResponse<PageResult<Emoji>>> => {
    return service({
        url: '/emoji/list',
        method: 'get',
        params: params,
    });
}

// 表情组列表
export const getEmojiGroups = (): Promise<ApiResponse<EmojiGroup[]>> => {
    return service({
        url: '/emoji/groups',
        method: 'get',
    });
}

// 创建表情组
export const createEmojiGroup = (data: EmojiGroupCreateRequest): Promise<ApiResponse<undefined>> => {
    return service({
        url: '/emoji/groups',
        method: 'post',
        data: data,
    });
}

// 更新表情组
export const updateEmojiGroup = (id: number, data: EmojiGroupUpdateRequest): Promise<ApiResponse<undefined>> => {
    return service({
        url: `/emoji/groups/${id}`,
        method: 'put',
        data: data,
    });
}

// 删除表情组
export const deleteEmojiGroup = (id: number): Promise<ApiResponse<undefined>> => {
    return service({
        url: `/emoji/groups/${id}`,
        method: 'delete',
    });
}

// 上传表情（multipart格式）
export const uploadEmoji = (formData: FormData): Promise<ApiResponse<{
    message: string;
    count: number;
    emojis: Emoji[];
}>> => {
    return service({
        url: '/emoji/upload',
        method: 'post',
        data: formData,
        headers: {
            'Content-Type': 'multipart/form-data'
        }
    });
}

// 流式上传表情（SSE）- 使用扩展后的通用streamRequest
export const uploadEmojiStream = (formData: FormData, onProgress?: (event: string, data: any) => void): Promise<void> => {
    const config: StreamRequestConfig = {
        url: '/emoji/upload',
        method: 'POST',
        data: formData,
        autoContentType: false, // 让浏览器自动设置multipart/form-data
        timeout: 300000 // 5分钟超时
    }
    
    return streamRequest(config, {
        onSSEEvent: (event: string, data: any) => {
            // 直接使用SSE事件回调
            onProgress?.(event, data)
        },
        onComplete: (responseData) => {
            console.log('流式上传完成:', responseData)
        },
        onError: (error) => {
            console.error('流式上传错误:', error)
            throw error
        }
    })
}

// 删除表情
export const deleteEmoji = (id: number): Promise<ApiResponse<undefined>> => {
    return service({
        url: `/emoji/${id}`,
        method: 'delete',
    });
}

// 恢复表情
export const restoreEmoji = (id: number): Promise<ApiResponse<undefined>> => {
    return service({
        url: `/emoji/${id}/restore`,
        method: 'put',
    });
}

// 重新生成雪碧图（SSE流式）
export const regenerateSpritesStream = (groupKeys: string[], onProgress?: (event: string, data: any) => void): Promise<void> => {
    const config: StreamRequestConfig = {
        url: '/emoji/regenerate',
        method: 'POST',
        data: {
            group_keys: groupKeys
        },
        timeout: 600000 // 10分钟超时
    }
    
    return streamRequest(config, {
        onSSEEvent: (event: string, data: any) => {
            // 直接使用SSE事件回调
            onProgress?.(event, data)
        },
        onComplete: (responseData) => {
            console.log('流式生成完成:', responseData)
        },
        onError: (error) => {
            console.error('流式生成错误:', error)
            throw error
        }
    })
}

// 获取任务状态
export const getTaskStatus = (id: number): Promise<ApiResponse<EmojiTask>> => {
    return service({
        url: `/emoji/task/${id}`,
        method: 'get',
    });
}

// 获取雪碧图列表
export const getSpriteList = (): Promise<ApiResponse<EmojiSprite[]>> => {
    return service({
        url: '/emoji/sprites',
        method: 'get',
    });
}

// 获取前端配置
export const getEmojiConfig = (): Promise<ApiResponse<EmojiConfig>> => {
    return service({
        url: '/emoji/config',
        method: 'get',
    });
}
