export const delHtmlTag = (value?: string): string => {
  if (!value) return "";
  const document = new DOMParser().parseFromString(value, "text/html");
  return (document.body.textContent || "").replace(/\s+/g, " ").trim();
};

export const getExcerpt = (value?: string, maxLength = 220, format = "html"): string => {
  // 旧版 API 没有 summary 时的轻量兼容；完整 Markdown 渲染仅在文章页加载。
  const source = format === "markdown" ? (value || "")
    .replace(/```[\s\S]*?```/g, " ")
    .replace(/!\[[^\]]*\]\([^)]*\)/g, " ")
    .replace(/\[([^\]]*)\]\([^)]*\)/g, "$1")
    .replace(/^\s{0,3}(?:#{1,6}\s*|>\s*|[-*+]\s+)/gm, "")
    .replace(/[*_`~]+/g, "") : value;
  const text = delHtmlTag(source);
  return text.length > maxLength ? `${text.slice(0, maxLength)}…` : text;
};
