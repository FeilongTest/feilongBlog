import { ref } from "vue";
import { defineStore } from "pinia";

export interface ReadingEntry {
  id: number;
  fid: number;
  title: string;
  savedAt: number;
}
const KEY = "feilong:reading-list:v1";
const LIMIT = 100;
function valid(value: unknown): value is ReadingEntry {
  const item = value as ReadingEntry;
  return (
    !!item &&
    Number.isSafeInteger(item.id) &&
    item.id > 0 &&
    Number.isSafeInteger(item.fid) &&
    item.fid >= 0 &&
    typeof item.title === "string" &&
    item.title.length <= 400 &&
    Number.isFinite(item.savedAt) &&
    item.savedAt > 0
  );
}

export const useReadingListStore = defineStore("reading-list", () => {
  const entries = ref<ReadingEntry[]>([]);
  const readError = ref(false);
  function load() {
    try {
      const raw = localStorage.getItem(KEY);
      const saved = raw ? JSON.parse(raw) : [];
      if (!Array.isArray(saved)) throw new Error("invalid reading list");
      const seen = new Set<number>();
      entries.value = saved
        .filter(valid)
        .filter((item) => {
          if (seen.has(item.id)) return false;
          seen.add(item.id);
          return true;
        })
        .slice(0, LIMIT);
      readError.value = false;
    } catch {
      readError.value = true;
    }
  }
  function write(next: ReadingEntry[]): boolean {
    try {
      localStorage.setItem(KEY, JSON.stringify(next));
      entries.value = next;
      readError.value = false;
      return true;
    } catch {
      return false;
    }
  }
  const has = (id: number) => entries.value.some((item) => item.id === id);
  function toggle(article: { ID: number; fid: number; title: string }): string {
    // Re-read before writing so another tab's changes are retained.
    load();
    if (readError.value) return "无法读取本机阅读清单，请检查浏览器存储设置。";
    if (has(article.ID))
      return write(entries.value.filter((item) => item.id !== article.ID))
        ? "已从稍后阅读移除"
        : "保存失败，请检查浏览器存储空间。";
    if (entries.value.length >= LIMIT)
      return "最多保存 100 篇文章，请先整理阅读清单。";
    const item: ReadingEntry = {
      id: article.ID,
      fid: article.fid,
      title: article.title.slice(0, 400),
      savedAt: Date.now(),
    };
    if (!valid(item)) return "文章信息不完整，暂时无法收藏。";
    return write([item, ...entries.value])
      ? "已加入稍后阅读"
      : "保存失败，请检查浏览器存储空间。";
  }
  function remove(id: number) {
    load();
    return (
      !readError.value && write(entries.value.filter((item) => item.id !== id))
    );
  }
  function clear() {
    return write([]);
  }
  window.addEventListener("storage", (event) => {
    if (event.key === KEY || event.key === null) load();
  });
  load();
  return { entries, readError, load, has, toggle, remove, clear };
});
