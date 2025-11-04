// 评论类型定义
export interface Comment {
  ID?: number;
  name: string;
  email: string;
  content: string;
  uid: number;
  ctime: number;
  aid: number;
  status: number;
}

export interface CommentFormData {
  name: string;
  email: string;
  content: string;
  aid: number;
}

