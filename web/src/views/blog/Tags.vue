<template>
  <div class="card">
    <div class="card-body p-lg-15">
      <div
        class="d-flex flex-wrap justify-content-between align-items-end gap-4 mb-9"
      >
        <div>
          <span class="text-primary fw-bold fs-7">TOPICS</span>
          <h1 class="text-gray-900 mt-3 mb-2">标签专题</h1>
          <p class="text-muted mb-0">从一个主题出发，发现不同分类中的文章。</p>
        </div>
        <router-link to="/archive" class="btn btn-sm btn-light"
          >查看文章归档</router-link
        >
      </div>
      <div class="d-flex flex-wrap gap-3 mb-7">
        <input
          v-model="query"
          type="search"
          aria-label="搜索标签"
          placeholder="搜索你感兴趣的标签"
          class="form-control form-control-solid flex-grow-1"
          style="max-width: 480px"
        /><select
          v-model="sort"
          aria-label="标签排序"
          class="form-select form-select-solid w-auto"
        >
          <option value="popular">文章数量</option>
          <option value="name">名称排序</option>
        </select>
      </div>
      <p v-if="loading" role="status" class="text-muted">正在整理标签…</p>
      <div v-else-if="error" class="alert alert-danger">
        标签加载失败。<button
          type="button"
          class="btn btn-sm btn-light ms-3"
          @click="load"
        >
          重新加载
        </button>
      </div>
      <template v-else
        ><p class="text-muted fs-7" role="status">
          {{ filtered.length }} 个公开主题
        </p>
        <div class="row g-5">
          <div v-for="tag in filtered" :key="tag.id" class="col-md-6 col-xl-4">
            <router-link
              :to="`/tag/${tag.id}`"
              class="topic-card border rounded p-5 d-block h-100"
              ><div
                class="d-flex justify-content-between align-items-center gap-3"
              >
                <span class="fw-bold fs-4 text-gray-900"># {{ tag.name }}</span
                ><span class="badge badge-light-primary"
                  >{{ tag.count }} 篇</span
                >
              </div>
              <p class="text-muted fs-7 mt-3 mb-0">
                {{ tag.description || `探索与 ${tag.name} 相关的文章。` }}
              </p></router-link
            >
          </div>
        </div>
        <p v-if="!filtered.length" class="text-muted py-8">
          {{
            query
              ? "没有找到这个标签，试试其他关键词。"
              : "还没有公开的标签专题。"
          }}
        </p></template
      >
    </div>
  </div>
</template>
<script setup lang="ts">
import { computed, onMounted, onBeforeUnmount, ref } from "vue";
import type { Tag } from "@/core/blog/PublishingTypes";
import service from "@/utils/request";
import { useSiteStore } from "@/stores/site";
import { absoluteUrl, applySeo } from "@/core/seo";
let active = true;
onBeforeUnmount(() => {
  active = false;
});
const tags = ref<Tag[]>([]);
const query = ref("");
const sort = ref("popular");
const loading = ref(true);
const error = ref(false);
const filtered = computed(() =>
  tags.value
    .filter((tag) =>
      (tag.name + " " + tag.description)
        .toLocaleLowerCase()
        .includes(query.value.trim().toLocaleLowerCase())
    )
    .sort((a, b) =>
      sort.value === "name"
        ? a.name.localeCompare(b.name, "zh-CN")
        : (b.count || 0) - (a.count || 0) ||
          a.name.localeCompare(b.name, "zh-CN")
    )
);
async function load() {
  loading.value = true;
  error.value = false;
  try {
    const response = await service.get("/base/getTags");
    tags.value = Array.isArray(response.data) ? response.data : [];
  } catch {
    error.value = true;
  } finally {
    loading.value = false;
  }
}
onMounted(async () => {
  const site = await useSiteStore().loadSiteSetting();
  if (!active) return;
  applySeo({
    title: `标签专题 - ${site.siteName}`,
    description: "按主题发现文章与相关内容。",
    canonical: absoluteUrl("/tags"),
    siteName: site.siteName,
  });
  void load();
});
</script>
<style scoped>
.topic-card {
  transition: transform 0.15s, border-color 0.15s;
}
.topic-card:hover {
  transform: translateY(-2px);
  border-color: var(--bs-primary) !important;
}
@media (prefers-reduced-motion: reduce) {
  .topic-card {
    transition: none;
  }
  .topic-card:hover {
    transform: none;
  }
}
</style>
