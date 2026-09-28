/**
 * 前台页面的 SEO 元信息管理。
 * 服务端已经为每个地址注入过一份完整的 meta，这里在客户端接管后保持一致，
 * 保证站内跳转时标题、描述与结构化数据同步更新。
 */

export interface SeoInput {
  title: string;
  description?: string;
  keywords?: string;
  canonical?: string;
  robots?: string;
  type?: "website" | "article";
  image?: string;
  siteName?: string;
  publishedTime?: string;
  modifiedTime?: string;
  section?: string;
  jsonLd?: Record<string, unknown>[];
}

// 由本模块写入或被服务端注入的标签统一带上该标记
const SEO_FLAG = "data-seo";

// 清理上一次写入的标签，避免站内跳转时重复
const clearSeo = () => {
  document.head.querySelectorAll(`[${SEO_FLAG}]`).forEach((element) => element.remove());
};

const upsertMeta = (attr: "name" | "property", key: string, content?: string) => {
  if (!content) return;
  let element = document.head.querySelector<HTMLMetaElement>(`meta[${attr}="${key}"]`);
  if (!element) {
    element = document.createElement("meta");
    element.setAttribute(attr, key);
    document.head.appendChild(element);
  }
  element.setAttribute("content", content);
  element.setAttribute(SEO_FLAG, "1");
};

const upsertLink = (rel: string, href?: string) => {
  if (!href) return;
  let element = document.head.querySelector<HTMLLinkElement>(`link[rel="${rel}"]`);
  if (!element) {
    element = document.createElement("link");
    element.setAttribute("rel", rel);
    document.head.appendChild(element);
  }
  element.setAttribute("href", href);
  element.setAttribute(SEO_FLAG, "1");
};

/** 把站内路径转换成绝对地址 */
export const absoluteUrl = (path: string) => {
  const normalized = path.startsWith("/") ? path : `/${path}`;
  return `${window.location.origin}${normalized}`;
};

/** 去掉 HTML 标签，用于生成页面描述 */
export const stripHtml = (html: string, limit = 150) => {
  const text = (html || "")
    .replace(/<(?:script|style)[^>]*>[\s\S]*?<\/(?:script|style)>/gi, " ")
    .replace(/<[^>]*>/g, " ")
    .replace(/\s+/g, " ")
    .trim();
  return text.length > limit ? `${text.slice(0, limit)}…` : text;
};

/** 写入当前页面的 SEO 信息 */
export const applySeo = (input: SeoInput) => {
  clearSeo();

  document.title = input.title;
  upsertMeta("name", "description", input.description);
  upsertMeta("name", "keywords", input.keywords);
  upsertMeta("name", "robots", input.robots || "index, follow");
  upsertLink("canonical", input.canonical);

  upsertMeta("property", "og:site_name", input.siteName);
  upsertMeta("property", "og:locale", "zh_CN");
  upsertMeta("property", "og:type", input.type || "website");
  upsertMeta("property", "og:title", input.title);
  upsertMeta("property", "og:description", input.description);
  upsertMeta("property", "og:url", input.canonical);
  upsertMeta("property", "og:image", input.image);
  upsertMeta("property", "article:published_time", input.publishedTime);
  upsertMeta("property", "article:modified_time", input.modifiedTime);
  upsertMeta("property", "article:section", input.section);

  upsertMeta("name", "twitter:card", input.image ? "summary_large_image" : "summary");
  if (input.image) upsertMeta("name", "twitter:image", input.image);
  upsertMeta("name", "twitter:title", input.title);
  upsertMeta("name", "twitter:description", input.description);

  (input.jsonLd || []).forEach((item) => {
    const script = document.createElement("script");
    script.type = "application/ld+json";
    script.setAttribute(SEO_FLAG, "1");
    script.textContent = JSON.stringify(item);
    document.head.appendChild(script);
  });
};
