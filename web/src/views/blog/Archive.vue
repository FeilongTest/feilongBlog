<template>
  <div class="card">
    <div class="card-body p-lg-17">
      <div class="d-flex flex-wrap align-items-end justify-content-between gap-3 mb-10">
        <div>
          <span class="text-primary fw-bold fs-7 text-uppercase">Archive</span>
          <h1 class="text-gray-900 fw-bold fs-2x mt-3 mb-2">文章归档</h1>
          <p class="text-muted mb-0">共 {{ entries.length }} 篇已发布文章</p>
        </div>
        <router-link to="/" class="btn btn-sm btn-light">返回首页</router-link>
      </div>

      <div v-if="loading" class="text-muted">正在整理归档…</div>
      <div v-else-if="loadError" class="text-center py-10">
        <h3>归档加载失败</h3>
        <button type="button" class="btn btn-sm btn-light-primary mt-4" @click="loadArchive">重新加载</button>
      </div>
      <div v-else-if="!groups.length" class="text-muted">暂时还没有文章。</div>
      <section v-for="group in groups" :key="group.year" class="mb-10">
        <h2 class="fs-3 text-gray-900 mb-4">{{ group.year }}</h2>
        <ul class="list-unstyled mb-0">
          <li v-for="item in group.items" :key="item.id" class="d-flex flex-wrap gap-3 py-3 border-bottom">
            <time class="text-muted fw-semibold" style="min-width: 7rem">{{ formatDate(item.ctime) }}</time>
            <router-link :to="`/content/${item.fid}/${item.id}`" class="text-gray-900 text-hover-primary fw-bold">
              {{ item.title }}
            </router-link>
          </li>
        </ul>
      </section>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, ref } from "vue";
import dayjs from "dayjs";
import { absoluteUrl, applySeo } from "@/core/seo";
import { useSiteStore } from "@/stores/site";
import service from "@/utils/request";

interface ArchiveEntry {
  id: number;
  fid: number;
  title: string;
  ctime: number;
}

const siteStore = useSiteStore();
const entries = ref<ArchiveEntry[]>([]);
const loading = ref(true);
const loadError = ref(false);

const groups = computed(() => {
  const grouped = new Map<string, ArchiveEntry[]>();
  entries.value.forEach((item) => {
    const year = item.ctime ? dayjs.unix(Number(item.ctime)).format("YYYY") : "更早";
    const list = grouped.get(year) || [];
    list.push(item);
    grouped.set(year, list);
  });
  return Array.from(grouped, ([year, items]) => ({ year, items }));
});

const formatDate = (timestamp: number) => (timestamp ? dayjs.unix(Number(timestamp)).format("YYYY-MM-DD") : "");

const loadArchive = async () => {
  loading.value = true;
  loadError.value = false;
  try {
    const { data } = await service.get("/base/getArchive");
    entries.value = data?.list || [];
  } catch {
    entries.value = [];
    loadError.value = true;
  } finally {
    loading.value = false;
  }
};

onMounted(async () => {
  const site = await siteStore.loadSiteSetting();
  applySeo({
    title: `文章归档 - ${site.siteName}`,
    description: `${site.siteName} 的全部文章归档。`,
    keywords: ["归档", site.siteName, site.keywords].filter(Boolean).join(","),
    canonical: absoluteUrl("/archive"),
    siteName: site.siteName,
    type: "website",
    jsonLd: [
      {
        "@context": "https://schema.org",
        "@type": "CollectionPage",
        name: `文章归档 - ${site.siteName}`,
        url: absoluteUrl("/archive"),
      },
    ],
  });
  await loadArchive();
});
</script>
