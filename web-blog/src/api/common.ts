export interface Model {
    id: number;
    created_at: Date;
    updated_at: Date;
}

export interface PageInfo {
    page: number;
    page_size: number;
}

export interface PageResult<T> {
    list: T[];
    current_page: number;
    page_size: number;
    total_items: number;
    total_pages: number;
}

export interface Hit<T> {
    _id: string;
    _source: T;
}
