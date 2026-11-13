export type ArticleStatus = 0 | 1 | 2 | 3;
export interface Tag {
  id: number;
  name: string;
  description: string;
  count?: number;
}
export interface ArticleRevision {
  id: number;
  articleId: number;
  version: number;
  title: string;
  savedAt: number;
  reason: string;
}
export interface ArticleSnapshot {
  title: string;
  content: string;
  contentFormat: "html" | "markdown";
  pic: string;
  fid: number;
  type: number;
  excerpt: string;
  seoTitle: string;
  seoDescription: string;
  tagNames: string[];
}
export const articleStatusName: Record<ArticleStatus, string> = {
  0: "已发布",
  1: "已隐藏",
  2: "草稿",
  3: "定时发布",
};
export function localDateTime(seconds: number): string {
  if (!seconds) return "";
  const d = new Date(seconds * 1000);
  const pad = (n: number) => String(n).padStart(2, "0");
  return `${d.getFullYear()}-${pad(d.getMonth() + 1)}-${pad(d.getDate())}T${pad(
    d.getHours()
  )}:${pad(d.getMinutes())}`;
}
export function unixDateTime(value: string): number {
  const time = new Date(value).getTime();
  return Number.isFinite(time) ? Math.floor(time / 1000) : 0;
}
