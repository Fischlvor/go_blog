import type {Model, PageInfo, PageResult} from "@/api/common";
import type {ApiResponse} from "@/utils/request";
import service, { adminService } from "@/utils/request";

export interface Image extends Model {
    name: string;
    url: string;
    size: number;
    mime_type: string;
}

export interface ImageUploadResponse {
    url: string;  // 完整 URL（包含 CDN 域名）
}

export interface ImageDeleteRequest {
    ids: number[];
}

export const imageDelete = (data: ImageDeleteRequest): Promise<ApiResponse<undefined>> => {
    return adminService({
        url: '/image/delete',
        method: 'delete',
        data: data,
    });
}

export interface ImageListRequest extends PageInfo {
    name: string | null;
    mime_type: string | null;
}

export const imageList = (data: ImageListRequest): Promise<ApiResponse<PageResult<Image>>> => {
    return adminService({
        url: '/image/list',
        method: 'get',
        params: data,
    });
}
