<!-- eslint-disable vue/multi-word-component-names -->
<template>
  <div v-if="keyword" class="card mb-7">
    <div class="card-body d-flex flex-wrap align-items-center justify-content-between gap-3 py-5">
      <div class="d-flex align-items-center">
        <span class="svg-icon svg-icon-2x svg-icon-primary me-3">
          <inline-svg :src="getAssetPath('/media/icons/duotune/general/gen021.svg')" />
        </span>
        <div>
          <div class="fw-bold text-gray-900">“{{ keyword }}”的搜索结果</div>
          <div class="text-muted fs-7">{{ loading ? "正在搜索…" : `找到 ${total} 篇文章` }}</div>
        </div>
      </div>
      <button type="button" class="btn btn-sm btn-light-primary" @click="clearSearch">清除搜索</button>
    </div>
  </div>

  <div v-if="loading" class="row g-5 g-xl-8" aria-label="正在加载文章">
    <div v-for="column in 2" :key="column" class="col-xl-6">
      <BlogFeedSkeleton v-for="item in 3" :key="item" class="mb-5 mb-xl-8" />
    </div>
  </div>

  <div v-else-if="articleList.length" class="row g-5 g-xl-8 feed-results">
    <div class="col-xl-6">
      <BlogFeeds v-for="article in leftColumn" :key="article.ID" class="mb-5 mb-xl-8" :article="article" />
    </div>
    <div class="col-xl-6">
      <BlogFeeds v-for="article in rightColumn" :key="article.ID" class="mb-5 mb-xl-8" :article="article" />
    </div>
  </div>

  <div v-else class="card py-15 text-center feed-results">
    <div class="card-body">
      <i class="bi bi-search fs-3x text-gray-300"></i>
      <h3 class="mt-5">没有找到相关文章</h3>
      <p class="text-muted mb-0">换一个关键词试试看。</p>
    </div>
  </div>

  <div v-if="!loading" class="d-flex flex-stack flex-wrap gap-4 pt-5 feed-results">
    <div class="fs-6 fw-semibold text-gray-600">共 {{ total }} 篇文章</div>
    <TablePagination v-if="total > pageSize" :total-pages="Math.ceil(total / pageSize)" :total="total"
      :per-page="pageSize" :current-page="page" @page-change="pageChange" />
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, ref, watch } from "vue";
import { useRouter } from "vue-router";
import BlogFeeds from "@/components/blog/BlogFeeds.vue";
import BlogFeedSkeleton from "@/components/blog/BlogFeedSkeleton.vue";
import TablePagination from "@/components/kt-datatable/table-partials/table-content/table-footer/TablePagination.vue";
import type { ArticleList } from "@/core/blog/ArticleTypes";
import { absoluteUrl, applySeo } from "@/core/seo";
import { getAssetPath } from "@/core/helpers/assets";
import { useSiteStore } from "@/stores/site";
import service from "@/utils/request";

const page = ref(1);
const total = ref(0);
const pageSize = ref(10);
const articleList = ref<ArticleList[]>([]);
const fid = ref(0);
const keyword = ref("");
const loading = ref(true);
let requestSequence = 0;
const router = useRouter();
const siteStore = useSiteStore();

const leftColumn = computed(() => articleList.value.filter((_, index) => index % 2 === 0));
const rightColumn = computed(() => articleList.value.filter((_, index) => index % 2 === 1));

const syncRoute = () => {
  const route = router.currentRoute.value;
  fid.value = route.name === "category" ? Number(route.params.id) : 0;
  keyword.value = route.query.keyword ? String(route.query.keyword) : "";
};

// 首页、分类页与搜索结果页各自的 SEO 信息
const updateSeo = async () => {
  const site = await siteStore.loadSiteSetting();

  // 站内搜索结果不参与收录，避免低质量页面被索引
  if (keyword.value) {
    applySeo({
      title: `“${keyword.value}”的搜索结果 - ${site.siteName}`,
      description: `在 ${site.siteName} 中搜索“${keyword.value}”的相关文章。`,
      canonical: absoluteUrl(router.currentRoute.value.fullPath),
      robots: "noindex, follow",
      siteName: site.siteName,
      type: "website",
    });
    return;
  }

  if (fid.value) {
    await siteStore.loadCategories();
    const category = siteStore.categoryName(fid.value);
    const title = `${category || "文章"} 分类 - ${site.siteName}`;
    applySeo({
      title,
      description: `${site.siteName} 中「${category || "该"}」分类下的全部文章列表。`,
      keywords: [category, site.siteName, site.keywords].filter(Boolean).join(","),
      canonical: absoluteUrl(`/category/${fid.value}`),
      siteName: site.siteName,
      type: "website",
      section: category || undefined,
      jsonLd: [
        {
          "@context": "https://schema.org",
          "@type": "CollectionPage",
          name: title,
          url: absoluteUrl(`/category/${fid.value}`),
        },
      ],
    });
    return;
  }

  applySeo({
    title: `${site.siteName} - ${site.siteSlogan}`,
    description: site.description,
    keywords: site.keywords,
    canonical: absoluteUrl("/"),
    siteName: site.siteName,
    type: "website",
    jsonLd: [
      {
        "@context": "https://schema.org",
        "@type": "WebSite",
        name: site.siteName,
        url: absoluteUrl("/"),
        description: site.description,
      },
    ],
  });
};

const getArticleList = async () => {
  const sequence = ++requestSequence;
  loading.value = true;
  try {
    const { data } = await service.get("/base/getArticleList", {
      params: { page: page.value, pageSize: pageSize.value, fid: fid.value, keyword: keyword.value || undefined },
    });
    if (sequence !== requestSequence) return;
    total.value = data?.total || 0;
    articleList.value = data?.list || [];
  } finally {
    if (sequence === requestSequence) {
      loading.value = false;
      void updateSeo();
    }
  }
};

const clearSearch = () => router.push({ name: "blog-home" });

const pageChange = (newPage: number) => {
  page.value = newPage;
  getArticleList();
  window.scrollTo({ top: 0, behavior: "smooth" });
};

onMounted(() => {
  syncRoute();
  getArticleList();
});

watch(() => router.currentRoute.value.fullPath, () => {
  page.value = 1;
  syncRoute();
  getArticleList();
});
</script>

<style scoped>
.feed-results { animation: content-enter .28s ease-out; }
@keyframes content-enter {
  from { opacity: 0; transform: translateY(6px); }
  to { opacity: 1; transform: translateY(0); }
}
@media (prefers-reduced-motion: reduce) {
  .feed-results { animation: none; }
}
</style>
