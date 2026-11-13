<template>
  <div class="card">
    <div class="card-body p-lg-15">
      <div
        class="d-flex justify-content-between flex-wrap gap-4 align-items-end mb-7"
      >
        <div>
          <span class="text-primary fw-bold fs-7">READ LATER</span>
          <h1 class="text-gray-900 mt-3 mb-2">稍后阅读</h1>
          <p class="text-muted mb-0">
            清单仅在此浏览器保存，无需登录。清除站点数据后无法恢复。
          </p>
        </div>
        <button
          v-if="list.entries.length"
          type="button"
          class="btn btn-sm btn-light-danger"
          @click="clear"
        >
          清空清单
        </button>
      </div>
      <div v-if="list.readError" class="alert alert-warning">
        无法读取本机阅读清单。请检查浏览器存储设置，或清空损坏的清单后重试。<button
          class="btn btn-sm btn-light ms-3"
          @click="clear"
        >
          重置清单
        </button>
      </div>
      <div v-if="availabilityError" class="alert alert-warning">
        暂时无法核对文章状态，仍可尝试打开已保存的链接。<button
          class="btn btn-sm btn-light ms-3"
          @click="checkAvailable"
        >
          重新核对
        </button>
      </div>
      <template v-if="list.entries.length"
        ><div class="d-flex flex-wrap gap-3 mb-6">
          <input
            v-model="query"
            type="search"
            aria-label="搜索阅读清单"
            class="form-control form-control-solid"
            style="max-width: 420px"
            placeholder="查找已保存的文章"
          /><select
            v-model="sort"
            aria-label="阅读清单排序"
            class="form-select form-select-solid w-auto"
          >
            <option value="newest">最近收藏</option>
            <option value="oldest">最早收藏</option>
          </select>
        </div>
        <p class="text-muted fs-7" role="status">
          {{ filtered.length }} 篇文章<span v-if="checking">
            · 正在核对文章状态…</span
          >
        </p>
        <ul class="list-unstyled">
          <li
            v-for="entry in filtered"
            :key="entry.id"
            class="d-flex gap-4 justify-content-between align-items-center py-5 border-bottom"
          >
            <div>
              <router-link
                v-if="!unavailable(entry.id)"
                :to="`/content/${available[entry.id]?.fid ?? entry.fid}/${
                  entry.id
                }`"
                class="fw-bold fs-5 text-gray-900 text-hover-primary"
                >{{ available[entry.id]?.title || entry.title }}</router-link
              ><span v-else class="text-muted fw-bold"
                >{{ entry.title }}
                <span class="badge badge-light ms-2">文章已下线</span></span
              >
              <div class="text-muted fs-8 mt-2">
                收藏于 {{ new Date(entry.savedAt).toLocaleString() }}
              </div>
            </div>
            <button
              type="button"
              class="btn btn-sm btn-light"
              :aria-label="`移除 ${entry.title}`"
              @click="remove(entry.id)"
            >
              移除
            </button>
          </li>
        </ul>
        <p v-if="!filtered.length" class="text-muted">
          没有找到符合条件的收藏。
        </p></template
      >
      <div v-else-if="!list.readError" class="text-center py-12">
        <i class="bi bi-bookmark fs-3x text-primary"></i>
        <h2 class="fs-4 mt-5">给下一次阅读留个位置</h2>
        <p class="text-muted">
          在文章卡片或详情页点击“稍后阅读”，最多可保存 100 篇。
        </p>
        <router-link to="/" class="btn btn-sm btn-primary"
          >去发现文章</router-link
        >
      </div>
      <p v-if="message" class="text-muted fs-7" role="status">{{ message }}</p>
    </div>
  </div>
</template>
<script setup lang="ts">
import { computed, onMounted, onBeforeUnmount, ref, watch } from "vue";
import Swal from "sweetalert2";
import { useReadingListStore } from "@/stores/readingList";
import { absoluteUrl, applySeo } from "@/core/seo";
import { useSiteStore } from "@/stores/site";
import service from "@/utils/request";
const list = useReadingListStore();
const query = ref("");
const sort = ref("newest");
const message = ref("");
const checking = ref(false);
const checked = ref(false);
const availabilityError = ref(false);
const available = ref<
  Record<number, { ID: number; fid: number; title: string }>
>({});
let sequence = 0;
let active = true;
const filtered = computed(() =>
  list.entries
    .filter((entry) =>
      (available.value[entry.id]?.title || entry.title)
        .toLocaleLowerCase()
        .includes(query.value.trim().toLocaleLowerCase())
    )
    .sort((a, b) =>
      sort.value === "newest" ? b.savedAt - a.savedAt : a.savedAt - b.savedAt
    )
);
const unavailable = (id: number) =>
  checked.value &&
  !checking.value &&
  !availabilityError.value &&
  !available.value[id];
async function checkAvailable() {
  const request = ++sequence;
  const ids = list.entries.map((entry) => entry.id);
  if (!ids.length) {
    available.value = {};
    checked.value = true;
    checking.value = false;
    return;
  }
  checking.value = true;
  availabilityError.value = false;
  try {
    const response = await service.post("/base/getReadingArticles", { ids });
    if (request !== sequence) return;
    const items = Array.isArray(response.data) ? response.data : [];
    available.value = Object.fromEntries(items.map((item) => [item.ID, item]));
    checked.value = true;
  } catch {
    if (request === sequence) availabilityError.value = true;
  } finally {
    if (request === sequence) checking.value = false;
  }
}
function remove(id: number) {
  message.value = list.remove(id)
    ? "已移除"
    : "移除失败，请检查浏览器存储设置。";
}
async function clear() {
  const answer = await Swal.fire({
    title: "清空本机阅读清单？",
    text: "只清除此浏览器的收藏，不影响站点文章。",
    showCancelButton: true,
    confirmButtonText: "清空清单",
    cancelButtonText: "取消",
  });
  if (answer.isConfirmed)
    message.value = list.clear()
      ? "清单已清空"
      : "清空失败，请检查浏览器存储设置。";
}
watch(
  () => list.entries.map((entry) => entry.id).join(","),
  () => {
    void checkAvailable();
  }
);
onMounted(async () => {
  const site = await useSiteStore().loadSiteSetting();
  if (!active) return;
  applySeo({
    title: `稍后阅读 - ${site.siteName}`,
    canonical: absoluteUrl("/reading-list"),
    robots: "noindex, nofollow",
    siteName: site.siteName,
  });
  void checkAvailable();
});
onBeforeUnmount(() => {
  sequence++;
  active = false;
});
</script>
