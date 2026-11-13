/** API 返回绝对 R2 地址或 /blog/...；兼容历史本地上传路径。 */
export function resolveMediaUrl(raw?: string): string {
  const value = (raw || "").trim();
  if (!value) return "";
  if (/^(?:https?:)?\/\//i.test(value)) return value;
  if (/^[a-z][a-z0-9+.-]*:/i.test(value)) return "";
  const path = value.replace(/^\.?\//, "");
  if (/^(?:public\/)?uploads\/file\//i.test(path)) {
    return `/blog/${path.replace(/^public\//i, "")}`;
  }
  if (/^public\/uploads\//i.test(path) && import.meta.env.VITE_MEDIA_BASE_URL) {
    return `${String(import.meta.env.VITE_MEDIA_BASE_URL).replace(/\/$/, "")}/legacy/Public/Uploads/${path.slice("public/uploads/".length)}`;
  }
  return `/${path}`;
}
