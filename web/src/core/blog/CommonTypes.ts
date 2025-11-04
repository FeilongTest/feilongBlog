export interface CommonResponse<T> {
    code: number;
    data: T;
    msg: string;
}