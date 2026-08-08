export const delHtmlTag = (value?: string): string => {
  if (!value) return "";
  const document = new DOMParser().parseFromString(value, "text/html");
  return (document.body.textContent || "").replace(/\s+/g, " ").trim();
};

export const getExcerpt = (value?: string, maxLength = 220): string => {
  const text = delHtmlTag(value);
  return text.length > maxLength ? `${text.slice(0, maxLength)}…` : text;
};
