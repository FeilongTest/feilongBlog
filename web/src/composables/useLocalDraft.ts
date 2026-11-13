import { computed, onBeforeUnmount, ref, shallowRef, watch } from "vue";

interface Draft<T> { version: 1; savedAt: number; value: T }
const MAX_AGE = 7 * 24 * 60 * 60 * 1000;

/** 按账号及文章隔离，只在本机保存；服务器保存成功后清除。 */
export function useLocalDraft<T>(snapshot: () => T, apply: (value: T) => void, validate: (value: unknown) => value is T) {
  const ready = ref(false);
  const baseline = ref("");
  const savedAt = ref(0);
  const storageError = ref(false);
  const pending = shallowRef<Draft<T> | null>(null);
  let key = "";
  let timer: ReturnType<typeof setTimeout> | undefined;
  const serialized = computed(() => JSON.stringify(snapshot()));
  const dirty = computed(() => ready.value && baseline.value !== serialized.value);

  function flush() {
    clearTimeout(timer);
    if (!ready.value || !key || pending.value) return;
    try {
      if (!dirty.value) { localStorage.removeItem(key); savedAt.value = 0; return; }
      const draft: Draft<T> = { version: 1, savedAt: Date.now(), value: JSON.parse(serialized.value) };
      localStorage.setItem(key, JSON.stringify(draft));
      savedAt.value = draft.savedAt;
      storageError.value = false;
    } catch { storageError.value = true; }
  }

  function pause() { flush(); ready.value = false; clearTimeout(timer); }

  function initialize(nextKey: string) {
    ready.value = false;
    key = `feilong:article-draft:${nextKey}`;
    baseline.value = serialized.value;
    pending.value = null;
    savedAt.value = 0;
    storageError.value = false;
    try {
      const raw = localStorage.getItem(key);
      if (raw) {
        const stored = JSON.parse(raw);
        if (stored?.version === 1 && Number.isFinite(stored.savedAt) && stored.savedAt <= Date.now() && Date.now() - stored.savedAt < MAX_AGE && validate(stored.value)) {
          if (JSON.stringify(stored.value) !== baseline.value) pending.value = stored;
          else localStorage.removeItem(key);
        } else localStorage.removeItem(key);
      }
    } catch { storageError.value = true; }
    ready.value = true;
  }

  function restore() {
    if (!pending.value) return;
    apply(pending.value.value);
    pending.value = null;
    flush();
  }

  function discard() {
    pending.value = null;
    try { localStorage.removeItem(key); storageError.value = false; }
    catch { storageError.value = true; }
    if (dirty.value) flush();
  }

  function markSaved(value: T = snapshot()) {
    baseline.value = JSON.stringify(value);
    pending.value = null;
    flush();
  }

  watch(serialized, () => { clearTimeout(timer); timer = setTimeout(flush, 750); });
  onBeforeUnmount(flush);
  return { dirty, savedAt, storageError, pending, initialize, pause, restore, discard, markSaved, flush };
}
