<template>
  <div class="card">
    <div class="card-body p-lg-20 pb-lg-0">
      <div class="d-flex flex-column flex-xl-row">
        <div class="flex-lg-row-fluid me-xl-15">
          <div class="mb-17">
            <div class="mb-8">
              <div class="d-flex flex-wrap mb-6">
                <div class="me-9 my-1">
                  <span class="svg-icon svg-icon-primary me-1 fs-2">
                    <inline-svg :src="getAssetPath('/media/icons/duotune/general/gen016.svg')" />
                  </span>
                  <span class="fw-bold text-gray-400">{{ formatDateTime(articleInfo.ctime) }}</span>
                </div>
                <div class="me-9 my-1">
                  <span class="svg-icon svg-icon-primary me-1 fs-2">
                    <inline-svg :src="getAssetPath('/media/icons/duotune/art/art005.svg')" />
                  </span>
                  <span class="fw-bold text-gray-400">{{ formatDateTime(articleInfo.editTime) }}</span>
                </div>
                <div class="me-9 my-1">
                  <span class="svg-icon svg-icon-primary me-1 fs-2">
                    <inline-svg :src="getAssetPath('/media/icons/duotune/coding/cod006.svg')" />
                  </span>
                  <span class="fw-bold text-gray-400">{{ articleInfo.view || 0 }}</span>
                </div>
              </div>

              <span class="text-dark text-hover-primary fs-2 fw-bold">{{ articleInfo.title }}</span>

              <div v-if="articleInfo.type !== 0 && articleInfo.pic" class="overlay mt-8">
                <div
                  class="bgi-no-repeat bgi-position-center bgi-size-cover card-rounded min-h-350px"
                  :style="{ backgroundImage: `url(${articleInfo.pic})` }"
                ></div>
              </div>
            </div>

            <div ref="contentRef" class="article-content mb-8" v-html="sanitizedContent"></div>

            <div class="d-flex align-items-center border border-dashed card-rounded p-5 p-lg-10 mb-14">
              <div class="text-center flex-shrink-0 me-7 me-lg-13">
                <div class="symbol symbol-70px symbol-circle mb-2">
                  <img :src="getAssetPath('/media/avatars/300-1.jpg')" alt="飞龙" />
                </div>
                <div>
                  <span class="text-gray-700 fw-bold">飞龙</span>
                  <span class="text-gray-400 fs-7 fw-semibold d-block mt-1">feilong</span>
                </div>
              </div>
              <div class="text-muted fw-semibold lh-lg">
                免责声明：本站提供的软件、教程和内容仅限用于学习和研究；请勿用于商业或非法用途。涉及第三方资源时请遵守相应许可协议并支持正版，如有侵权请联系我们处理。
              </div>
            </div>
          </div>

          <div class="mt-15">
            <CommentSection v-if="articleInfo.ID" :article-id="articleInfo.ID" />
          </div>
        </div>

        <div class="flex-column flex-lg-row-auto w-100 w-xl-300px mb-10">
          <div class="mb-16">
            <h4 class="text-dark mb-7">搜索文章</h4>
            <div class="position-relative">
              <span class="svg-icon svg-icon-2 svg-icon-gray-500 position-absolute top-50 translate-middle-y ms-3">
                <inline-svg :src="getAssetPath('/media/icons/duotune/general/gen021.svg')" />
              </span>
              <input
                v-model="searchKeyword"
                type="text"
                class="form-control form-control-solid ps-10"
                placeholder="输入文章标题关键字搜索"
                @keyup.enter="handleSearch"
              />
              <button
                v-if="searchKeyword"
                class="btn btn-sm btn-icon btn-active-light-primary position-absolute end-0 top-50 translate-middle-y me-2"
                type="button"
                @click="clearSearch"
              >
                <i class="bi bi-x-lg fs-6"></i>
              </button>
            </div>
          </div>

          <div class="mb-16">
            <h4 class="text-dark mb-7">相关分类</h4>
            <div v-for="item in articleSummary" :key="item.fid" class="d-flex flex-stack fw-semibold fs-5 text-muted mb-4">
              <router-link class="text-muted text-hover-primary pe-2" :to="{ name: 'category', params: { id: item.fid } }">
                <span class="badge badge-light-warning fw-bold my-2 me-3">推荐板块</span>{{ item.name }}
              </router-link>
              <div>共 <span class="badge badge-light-info fw-bold my-2">{{ item.total }}</span> 篇</div>
            </div>
          </div>

          <div>
            <h4 class="text-dark mb-7">最新发布</h4>
            <div v-for="item in articleList" :key="item.ID" class="d-flex mb-7">
              <div class="symbol symbol-60px symbol-2by3 me-4">
                <div class="symbol-label" :style="{ backgroundImage: `url(${item.pic})` }"></div>
              </div>
              <div class="align-self-center">
                <span class="text-dark fw-bold text-hover-primary fs-6 pe-4 article-link" @click="gotoContent(item)">
                  {{ item.title }}
                </span>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, nextTick, ref, watch } from "vue";
import { useRoute, useRouter } from "vue-router";
import dayjs from "dayjs";
import DOMPurify from "dompurify";
import MarkdownIt from "markdown-it";
import Prism from "prismjs";
import Swal from "sweetalert2";
import CommentSection from "@/components/blog/CommentSection.vue";
import type { ArticleList, ArticelSummary } from "@/core/blog/ArticleTypes";
import { getAssetPath } from "@/core/helpers/assets";
import service from "@/utils/request";

const route = useRoute();
const router = useRouter();
const articleInfo = ref<ArticleList>({} as ArticleList);
const articleSummary = ref<ArticelSummary[]>([]);
const articleList = ref<ArticleList[]>([]);
const searchKeyword = ref("");
const contentRef = ref<HTMLElement>();
const markdown = new MarkdownIt({ html: false, linkify: true, breaks: true });
const sanitizedContent = computed(() => {
  const source = articleInfo.value.content || "";
  const html = articleInfo.value.contentFormat === "markdown" ? markdown.render(source) : source;
  return DOMPurify.sanitize(html);
});

const formatDateTime = (timestamp?: number) =>
  timestamp ? dayjs.unix(Number(timestamp)).format("YYYY年MM月DD日 HH时mm分") : "";

const enhanceCodeBlocks = async () => {
  await nextTick();
  const container = contentRef.value;
  if (!container) return;

  container.querySelectorAll<HTMLElement>("pre").forEach((pre) => {
    const code = pre.querySelector<HTMLElement>("code");
    if (code) Prism.highlightElement(code);

    const button = document.createElement("button");
    button.type = "button";
    button.className = "article-code-copy";
    button.textContent = "复制";
    button.setAttribute("aria-label", "复制代码");
    button.addEventListener("click", async () => {
      await navigator.clipboard.writeText(code?.textContent || pre.textContent || "");
      button.textContent = "已复制";
      window.setTimeout(() => (button.textContent = "复制"), 1600);
    });
    pre.appendChild(button);
  });
};

const getSummary = async (fid: number) => {
  const res: any = await service.post("/base/getSummary", { id: fid });
  articleSummary.value = res.data || [];
};

const getArticleList = async () => {
  const res: any = await service.get("/base/getArticleList", { params: { page: 1, pageSize: 5 } });
  articleList.value = res.data?.list || [];
};

const getArticle = async (id: number) => {
  const res: any = await service.post("/base/getArticle", { id });
  if (res.code === 7) {
    await Swal.fire({
      title: "出错了",
      icon: "error",
      text: "没有找到文章信息，即将返回首页。",
      buttonsStyling: false,
      confirmButtonText: "返回首页",
      heightAuto: false,
      customClass: { confirmButton: "btn fw-semibold btn-light-primary" },
    });
    await router.push({ name: "blog-home" });
    return;
  }

  articleInfo.value = res.data;
  const appName = import.meta.env.VITE_APP_NAME || "Feilong'S Blog";
  document.title = `${articleInfo.value.title} - ${appName}`;
  await Promise.all([getSummary(articleInfo.value.fid), enhanceCodeBlocks()]);
};

const gotoContent = (article: ArticleList) => {
  router.replace({ name: "blog-content", params: { fid: article.fid, id: article.ID } });
};

const handleSearch = () => {
  const keyword = searchKeyword.value.trim();
  if (keyword) router.push({ name: "blog-home", query: { keyword } });
};

const clearSearch = () => {
  searchKeyword.value = "";
  router.push({ name: "blog-home" });
};

watch(
  () => route.params.id,
  (id) => {
    const articleId = Number(id);
    if (articleId) getArticle(articleId);
  },
  { immediate: true },
);

getArticleList();
</script>

<style scoped lang="scss">
.article-link { cursor: pointer; }

.article-content {
  color: var(--kt-text-gray-700);
  font-family: Inter, "PingFang SC", "Microsoft YaHei", "Noto Sans SC", sans-serif;
  font-size: 16px;
  font-weight: 400;
  line-height: 1.9;
  overflow-wrap: anywhere;
}

.article-content :deep(h1),
.article-content :deep(h2),
.article-content :deep(h3),
.article-content :deep(h4) {
  color: var(--kt-text-gray-900);
  font-weight: 700;
  line-height: 1.4;
}

.article-content :deep(h1) { margin: 2.8rem 0 1.2rem; font-size: 2rem; }
.article-content :deep(h2) { margin: 2.5rem 0 1rem; font-size: 1.65rem; }
.article-content :deep(h3) { margin: 2rem 0 0.85rem; font-size: 1.3rem; }
.article-content :deep(h4) { margin: 1.75rem 0 0.7rem; font-size: 1.1rem; }
.article-content :deep(p) { margin: 0 0 1.25rem; }
.article-content :deep(ul),
.article-content :deep(ol) { margin: 0 0 1.4rem; padding-left: 1.6rem; }
.article-content :deep(li) { margin: 0.35rem 0; }

.article-content :deep(a) {
  color: var(--kt-primary);
  text-decoration: underline;
  text-decoration-color: rgba(0, 158, 247, 0.3);
  text-underline-offset: 0.2em;
}

.article-content :deep(blockquote) {
  margin: 1.75rem 0;
  padding: 1rem 1.2rem;
  border-left: 4px solid var(--kt-primary);
  border-radius: 0 0.65rem 0.65rem 0;
  background: var(--kt-gray-100);
}

.article-content :deep(blockquote p:last-child) { margin-bottom: 0; }

.article-content :deep(img) {
  display: block;
  width: auto;
  max-width: 100%;
  height: auto;
  margin: 2rem auto;
  border-radius: 0.75rem;
  box-shadow: 0 12px 30px rgba(24, 24, 27, 0.1);
}

.article-content :deep(code:not([class*="language-"])) {
  padding: 0.16em 0.38em;
  border: 1px solid var(--kt-gray-300);
  border-radius: 0.35rem;
  background: var(--kt-gray-100);
  color: #be123c;
  font-size: 0.88em;
}

.article-content :deep(pre) {
  position: relative;
  max-width: 100%;
  margin: 1.75rem 0;
  padding: 1.15rem 3.75rem 1.15rem 1.15rem;
  overflow: auto;
  border: 1px solid rgba(255, 255, 255, 0.08);
  border-radius: 0.7rem;
  background: #24292e !important;
  box-shadow: 0 10px 24px rgba(24, 24, 27, 0.13);
  color: #e6edf3;
  font-size: 14px;
  line-height: 1.72;
  tab-size: 2;
  white-space: pre;
}

.article-content :deep(pre code) {
  display: block;
  padding: 0;
  border: 0;
  background: transparent !important;
  color: inherit;
  font-family: "JetBrains Mono", "SFMono-Regular", Consolas, "Liberation Mono", monospace;
  font-size: inherit;
  text-shadow: none;
  white-space: inherit;
}

.article-content :deep(.article-code-copy) {
  position: absolute;
  top: 0.75rem;
  right: 0.75rem;
  padding: 0.3rem 0.55rem;
  border: 1px solid rgba(255, 255, 255, 0.15);
  border-radius: 0.4rem;
  background: rgba(255, 255, 255, 0.08);
  color: #cbd5e1;
  font: 12px/1.3 Inter, sans-serif;
}

.article-content :deep(.article-code-copy:hover) { background: rgba(255, 255, 255, 0.16); color: #fff; }

.article-content :deep(table) {
  display: block;
  width: 100%;
  margin: 1.75rem 0;
  overflow-x: auto;
  border-collapse: collapse;
}

.article-content :deep(th),
.article-content :deep(td) { padding: 0.7rem 0.9rem; border: 1px solid var(--kt-gray-300); text-align: left; }
.article-content :deep(th) { background: var(--kt-gray-100); color: var(--kt-text-gray-900); }

@media (max-width: 767.98px) {
  .article-content { font-size: 15px; line-height: 1.85; }
  .article-content :deep(pre) { padding: 1rem 3.2rem 1rem 1rem; font-size: 13px; }
}
</style>
