import { delHtmlTag } from "@/utils/html";

/** 中文按字、拉丁文字按词估算；仅在正文渲染后统计。 */
export function readingStats(html: string) {
  const text = delHtmlTag(html);
  const chinese = (text.match(/[\u3400-\u9fff]/g) || []).length;
  const words = (text.replace(/[\u3400-\u9fff]/g, " ").match(/[a-zA-Z0-9]+(?:['’-][a-zA-Z0-9]+)*/g) || []).length;
  return { count: chinese + words, minutes: Math.max(1, Math.ceil(chinese / 300 + words / 200)) };
}
