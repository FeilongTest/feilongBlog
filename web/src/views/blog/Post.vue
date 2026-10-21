<template>
  <div class="reading-progress" aria-hidden="true">
    <span class="reading-progress-bar" :style="{ transform: `scaleX(${readingProgress})` }"></span>
  </div>

  <PostSkeleton v-if="loading" />
  <div v-else-if="loadError" class="card py-15 text-center post-loaded">
    <div class="card-body">
      <i class="bi bi-cloud-slash fs-3x text-gray-300"></i>
      <h3 class="mt-5">文章加载失败</h3>
      <p class="text-muted">网络可能暂时不可用，请稍后重试。</p>
      <button type="button" class="btn btn-sm btn-light-primary me-3" @click="retryArticle">重新加载</button>
      <router-link :to="{ name: 'blog-home' }" class="btn btn-sm btn-light">返回首页</router-link>
    </div>
  </div>
  <div v-else class="card post-loaded">
    <div class="card-body p-lg-20 pb-lg-0">
      <div class="d-flex flex-column flex-xl-row">
        <div class="flex-lg-row-fluid me-xl-15">
          <div ref="articleRef" class="mb-17">
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

              <button
                v-if="articleInfo.type !== 0 && articleInfo.pic"
                type="button"
                class="article-cover"
                aria-label="查看缩略图大图"
                @click="openLightbox(articleInfo.pic, articleInfo.title)"
              >
                <img :src="articleInfo.pic" :alt="articleInfo.title" />
              </button>
            </div>

            <div ref="contentRef" class="article-content mb-8" @click="onContentClick" v-html="sanitizedContent"></div>

            <div class="d-flex align-items-center border border-dashed card-rounded p-5 p-lg-10 mb-14">
              <div class="text-center flex-shrink-0 me-7 me-lg-13">
                <div class="symbol symbol-70px symbol-circle mb-2">
                  <img :src="getAssetPath('/media/avatars/300-1.jpg')" alt="飞龙" />
                </div>
                <div>
                  <router-link to="/about" class="text-gray-700 text-hover-primary fw-bold">飞龙</router-link>
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
          <div v-if="toc.length" class="mb-16">
            <h4 class="text-dark mb-7">文章目录</h4>
            <nav class="toc-nav" aria-label="文章目录">
              <a
                v-for="item in toc"
                :key="item.id"
                class="toc-link"
                :class="{ 'toc-level-3': item.level === 3, 'toc-active': activeHeading === item.id }"
                :href="`#${item.id}`"
                @click.prevent="scrollToHeading(item.id)"
              >
                {{ item.text }}
              </a>
            </nav>
          </div>

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
            <div v-if="summaryLoading" class="placeholder-glow" aria-label="正在加载相关分类">
              <span v-for="item in 3" :key="item" class="placeholder col-12 d-block mb-5"></span>
            </div>
            <template v-else>
              <div v-for="item in articleSummary" :key="item.fid" class="d-flex flex-stack fw-semibold fs-5 text-muted mb-4 sidebar-results">
                <router-link class="text-muted text-hover-primary pe-2" :to="{ name: 'category', params: { id: item.fid } }">
                  <span class="badge badge-light-warning fw-bold my-2 me-3">推荐板块</span>{{ item.name }}
                </router-link>
                <div>共 <span class="badge badge-light-info fw-bold my-2">{{ item.total }}</span> 篇</div>
              </div>
            </template>
          </div>

          <div v-if="asideLoading || aside.related.length" class="mb-16">
            <h4 class="text-dark mb-7">相关文章</h4>
            <div v-if="asideLoading" class="placeholder-glow" aria-label="正在加载相关文章">
              <div v-for="item in 3" :key="item" class="d-flex align-items-center mb-7">
                <span class="placeholder rounded latest-image-placeholder me-4"></span>
                <div class="flex-grow-1">
                  <span class="placeholder col-12 d-block mb-3"></span>
                  <span class="placeholder col-7 d-block"></span>
                </div>
              </div>
            </div>
            <div v-else v-for="item in aside.related" :key="item.ID" class="d-flex mb-7 sidebar-results">
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

          <div v-if="asideLoading || aside.hot.length" class="mb-16">
            <h4 class="text-dark mb-7">热门文章</h4>
            <div v-if="asideLoading" class="placeholder-glow" aria-label="正在加载热门文章">
              <div v-for="item in 3" :key="item" class="d-flex align-items-center mb-7">
                <span class="placeholder rounded latest-image-placeholder me-4"></span>
                <div class="flex-grow-1">
                  <span class="placeholder col-12 d-block mb-3"></span>
                  <span class="placeholder col-7 d-block"></span>
                </div>
              </div>
            </div>
            <div v-else v-for="item in aside.hot" :key="item.ID" class="d-flex mb-7 sidebar-results">
              <div class="symbol symbol-60px symbol-2by3 me-4">
                <div class="symbol-label" :style="{ backgroundImage: `url(${item.pic})` }"></div>
              </div>
              <div class="align-self-center">
                <span class="text-dark fw-bold text-hover-primary fs-6 pe-4 article-link" @click="gotoContent(item)">
                  {{ item.title }}
                </span>
                <span class="text-muted fs-8 d-block mt-1">{{ item.view || 0 }} 次阅读</span>
              </div>
            </div>
          </div>

          <div>
            <h4 class="text-dark mb-7">最新发布</h4>
            <div v-if="latestLoading" class="placeholder-glow" aria-label="正在加载最新文章">
              <div v-for="item in 4" :key="item" class="d-flex align-items-center mb-7">
                <span class="placeholder rounded latest-image-placeholder me-4"></span>
                <div class="flex-grow-1">
                  <span class="placeholder col-12 d-block mb-3"></span>
                  <span class="placeholder col-7 d-block"></span>
                </div>
              </div>
            </div>
            <template v-else>
              <div v-for="item in articleList" :key="item.ID" class="d-flex mb-7 sidebar-results">
                <div class="symbol symbol-60px symbol-2by3 me-4">
                  <div class="symbol-label" :style="{ backgroundImage: `url(${item.pic})` }"></div>
                </div>
                <div class="align-self-center">
                  <span class="text-dark fw-bold text-hover-primary fs-6 pe-4 article-link" @click="gotoContent(item)">
                    {{ item.title }}
                  </span>
                </div>
              </div>
            </template>
          </div>
        </div>
      </div>
    </div>
  </div>

  <Teleport to="body">
    <div v-if="lightboxSrc" class="image-lightbox" @click.self="closeLightbox">
      <button type="button" class="image-lightbox-close" aria-label="关闭大图" @click="closeLightbox">×</button>
      <img :src="lightboxSrc" :alt="lightboxAlt" />
    </div>
  </Teleport>
</template>

<script setup lang="ts">
import { computed, nextTick, onBeforeUnmount, onMounted, ref, watch } from "vue";
import { useRoute, useRouter } from "vue-router";
import dayjs from "dayjs";
import DOMPurify from "dompurify";
import MarkdownIt from "markdown-it";
import Prism from "prismjs";
import CommentSection from "@/components/blog/CommentSection.vue";
import PostSkeleton from "@/components/blog/PostSkeleton.vue";
import type { ArticleBrief, ArticleList, ArticelSummary, PostAside } from "@/core/blog/ArticleTypes";
import { absoluteUrl, applySeo, stripHtml } from "@/core/seo";
import { getAssetPath } from "@/core/helpers/assets";
import { useSiteStore } from "@/stores/site";
import service from "@/utils/request";

const route = useRoute();
const router = useRouter();
const siteStore = useSiteStore();
const articleInfo = ref<ArticleList>({} as ArticleList);
const articleSummary = ref<ArticelSummary[]>([]);
const articleList = ref<ArticleList[]>([]);
const searchKeyword = ref("");
const contentRef = ref<HTMLElement>();
const articleRef = ref<HTMLElement>();
const loading = ref(true);
const loadError = ref(false);
const summaryLoading = ref(true);
const latestLoading = ref(true);
let requestSequence = 0;

// 文章目录与阅读进度
interface TocItem {
  id: string;
  text: string;
  level: number;
}
const activeHeading = ref("");
const readingProgress = ref(0);
// 相关文章与热门文章
const aside = ref<PostAside>({ related: [], hot: [] });
const asideLoading = ref(true);
// 图片灯箱
const lightboxSrc = ref("");
const lightboxAlt = ref("");
const markdown = new MarkdownIt({ html: false, linkify: true, breaks: true });

const decodeHeadingText = (html: string) => html
  .replace(/<[^>]+>/g, "")
  .replace(/&nbsp;/g, " ")
  .replace(/&amp;/g, "&")
  .replace(/&lt;/g, "<")
  .replace(/&gt;/g, ">")
  .replace(/&quot;/g, "\"")
  .replace(/&#39;/g, "'")
  .replace(/\s+/g, " ")
  .trim();

// 把标题锚点写进 HTML。只改 DOM 的话，正文重新渲染后锚点会丢，目录就跳不到。
const headingPattern = () => /<h([1-4])([^>]*)>([\s\S]*?)<\/h\1>/gi;

const injectHeadingIds = (html: string) => {
  const matches = Array.from(html.matchAll(headingPattern()));
  const hasSection = matches.some((match) => match[1] !== "1");
  const items: TocItem[] = [];
  let index = 0;
  const output = html.replace(headingPattern(), (full, level: string, attrs: string, inner: string) => {
    if (hasSection && level === "1") return full;
    const text = decodeHeadingText(inner);
    if (!text) return full;
    index += 1;
    const existing = attrs.match(/\sid=["']([^"']+)["']/i);
    const id = existing?.[1] || `section-${index}`;
    items.push({ id, text, level: Number(level) });
    if (existing) return full;
    return `<h${level}${attrs} id="${id}">${inner}</h${level}>`;
  });
  return { html: output, toc: items };
};

const renderedArticle = computed(() => {
  const source = articleInfo.value.content || "";
  const html = articleInfo.value.contentFormat === "markdown" ? markdown.render(source) : source;
  return injectHeadingIds(DOMPurify.sanitize(html, { ADD_ATTR: ["id"] }));
});
const sanitizedContent = computed(() => renderedArticle.value.html);
const toc = computed(() => renderedArticle.value.toc);

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
  summaryLoading.value = true;
  try {
    const res: any = await service.post("/base/getSummary", { id: fid });
    articleSummary.value = res.data || [];
  } catch {
    articleSummary.value = [];
  } finally {
    summaryLoading.value = false;
  }
};

const getArticleList = async () => {
  latestLoading.value = true;
  try {
    const res: any = await service.get("/base/getArticleList", { params: { page: 1, pageSize: 5 } });
    articleList.value = res.data?.list || [];
  } catch {
    articleList.value = [];
  } finally {
    latestLoading.value = false;
  }
};

// 站点信息与文章内容共同决定详情页的 SEO 元信息
const applyArticleSeo = async (article: ArticleList) => {
  const site = await siteStore.loadSiteSetting();
  const canonical = absoluteUrl(`/content/${article.fid}/${article.ID}`);
  const publishedTime = article.ctime ? dayjs.unix(article.ctime).toISOString() : undefined;
  const modifiedTime = article.editTime ? dayjs.unix(article.editTime).toISOString() : undefined;
  const description = stripHtml(sanitizedContent.value) || `${site.siteName} 的文章：${article.title}`;

  applySeo({
    title: `${article.title} - ${site.siteName}`,
    description,
    keywords: [article.title, site.siteName].filter(Boolean).join(","),
    canonical,
    type: "article",
    image: article.pic || undefined,
    siteName: site.siteName,
    publishedTime,
    modifiedTime,
    jsonLd: [
      {
        "@context": "https://schema.org",
        "@type": "BlogPosting",
        headline: article.title,
        description,
        url: canonical,
        mainEntityOfPage: { "@type": "WebPage", "@id": canonical },
        author: { "@type": "Person", name: site.author || site.siteName, url: absoluteUrl("/about") },
        publisher: { "@type": "Organization", name: site.siteName },
        datePublished: publishedTime,
        dateModified: modifiedTime,
        image: article.pic ? [article.pic] : undefined,
      },
    ],
  });
};

const getArticle = async (id: number) => {
  const sequence = ++requestSequence;
  loading.value = true;
  loadError.value = false;
  activeHeading.value = "";
  readingProgress.value = 0;
  aside.value = { related: [], hot: [] };
  try {
    const res: any = await service.post("/base/getArticle", { id });
    if (sequence !== requestSequence) return;
    articleInfo.value = res.data;
    await applyArticleSeo(articleInfo.value);
    loading.value = false;
    await Promise.all([getSummary(articleInfo.value.fid), enhanceCodeBlocks(), getPostAside()]);
    updateReadingState();
  } catch {
    if (sequence === requestSequence) loadError.value = true;
  } finally {
    if (sequence === requestSequence) loading.value = false;
  }
};

const retryArticle = () => {
  const articleId = Number(route.params.id);
  if (articleId) getArticle(articleId);
};

const gotoContent = (article: ArticleList | ArticleBrief) => {
  router.replace({ name: "blog-content", params: { fid: article.fid, id: article.ID } });
};

const readingOffset = () => {
  const header = document.getElementById("kt_header");
  return (header?.offsetHeight || 72) + 16;
};

// 点击目录跳转到对应标题。用 scrollIntoView，避免页面实际滚动容器不是 window 时跳不动。
const scrollToHeading = (id: string) => {
  const target = document.getElementById(id);
  if (!target) return;
  const reduceMotion = window.matchMedia("(prefers-reduced-motion: reduce)").matches;
  target.scrollIntoView({ behavior: reduceMotion ? "auto" : "smooth", block: "start" });
  activeHeading.value = id;
};

const openLightbox = (src: string, alt = "") => {
  const url = src.trim();
  if (!url) return;
  lightboxSrc.value = url;
  lightboxAlt.value = alt;
};

const closeLightbox = () => {
  lightboxSrc.value = "";
  lightboxAlt.value = "";
};

// 正文图片和被链接包住的缩略图都打开大图，而不是跟着外链跳走
const onContentClick = (event: MouseEvent) => {
  const target = event.target;
  if (!(target instanceof Element)) return;
  const image = target.closest("img");
  if (!image || !contentRef.value?.contains(image)) return;
  const link = image.closest("a");
  if (link) {
    const href = link.getAttribute("href") || "";
    const sameImage = href === image.getAttribute("src") || href === image.currentSrc;
    if (!sameImage && href && !href.startsWith("#")) return;
    event.preventDefault();
  }
  openLightbox(image.currentSrc || image.src, image.alt || articleInfo.value.title || "");
};

// 阅读进度按正文区域计算，顶部细条用缩放绘制，避免宽度动画拖在滚动后面
const updateReadingState = () => {
  const article = articleRef.value;
  const offset = readingOffset();
  if (article) {
    const rect = article.getBoundingClientRect();
    const distance = rect.height - (window.innerHeight - offset);
    if (distance <= 1) {
      readingProgress.value = rect.top <= offset ? 1 : 0;
    } else {
      const scrolled = offset - rect.top;
      readingProgress.value = Math.min(1, Math.max(0, scrolled / distance));
    }
  }

  if (!toc.value.length) return;
  let current = "";
  for (const item of toc.value) {
    const element = document.getElementById(item.id);
    if (!element) continue;
    if (element.getBoundingClientRect().top <= offset + 4) current = item.id;
    else break;
  }
  activeHeading.value = current;
};

let scrollFrame = 0;
const handleScroll = () => {
  if (scrollFrame) return;
  scrollFrame = window.requestAnimationFrame(() => {
    scrollFrame = 0;
    updateReadingState();
  });
};

const handleKeydown = (event: KeyboardEvent) => {
  if (event.key === "Escape") closeLightbox();
};

// 相关文章与热门文章
const getPostAside = async () => {
  asideLoading.value = true;
  try {
    const res: any = await service.get("/base/getPostAside", {
      params: { fid: articleInfo.value.fid, exclude: articleInfo.value.ID, limit: 5 },
    });
    aside.value = {
      related: res.data?.related || [],
      hot: res.data?.hot || [],
    };
  } catch {
    aside.value = { related: [], hot: [] };
  } finally {
    asideLoading.value = false;
  }
};

onMounted(() => {
  window.addEventListener("scroll", handleScroll, { passive: true });
  window.addEventListener("keydown", handleKeydown);
});

onBeforeUnmount(() => {
  window.removeEventListener("scroll", handleScroll);
  window.removeEventListener("keydown", handleKeydown);
  if (scrollFrame) window.cancelAnimationFrame(scrollFrame);
});

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
.post-loaded { animation: content-enter .28s ease-out; }
.placeholder { background-color: var(--kt-gray-300); }
.latest-image-placeholder { width: 60px; height: 80px; flex: 0 0 60px; }
.sidebar-results { animation: content-enter .24s ease-out; }

/* 顶部阅读进度。参考 PaperMod / NexT：固定在视口最上沿，用缩放而不是改宽度 */
.reading-progress {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  z-index: 1200;
  height: 3px;
  pointer-events: none;
}

.reading-progress-bar {
  display: block;
  width: 100%;
  height: 100%;
  transform: scaleX(0);
  transform-origin: left center;
  background: var(--kt-primary);
}

.article-cover {
  display: block;
  width: 100%;
  margin: 2rem 0 0;
  padding: 0;
  border: 0;
  appearance: none;
  background: transparent;
  font: inherit;
  text-align: left;
  cursor: pointer;
}

.article-cover img {
  display: block;
  width: 100%;
  height: 350px;
  object-fit: cover;
  border-radius: 0.75rem;
}

/* 文章目录 */
.toc-nav {
  display: flex;
  flex-direction: column;
  max-height: 60vh;
  overflow-y: auto;
  border-left: 2px solid var(--kt-gray-300);
}

.toc-link {
  padding: 0.35rem 0.85rem;
  margin-left: -2px;
  border-left: 2px solid transparent;
  color: var(--kt-text-gray-600);
  font-size: 0.875rem;
  line-height: 1.5;
  text-decoration: none;
  transition: color 0.15s ease, border-color 0.15s ease;
}

.toc-link:hover { color: var(--kt-primary); }
.toc-level-3 { padding-left: 1.6rem; font-size: 0.82rem; }
.toc-active { border-left-color: var(--kt-primary); color: var(--kt-primary); font-weight: 600; }

.article-content :deep(img) { cursor: pointer; }

.image-lightbox {
  position: fixed;
  inset: 0;
  z-index: 2000;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 2rem;
  background: rgba(15, 23, 42, 0.88);
  cursor: default;
  animation: lightbox-enter 0.18s ease-out;
}

.image-lightbox img {
  max-width: min(1100px, 100%);
  max-height: 100%;
  object-fit: contain;
  border-radius: 0.5rem;
  box-shadow: 0 24px 60px rgba(0, 0, 0, 0.45);
  cursor: default;
}

.image-lightbox-close {
  position: absolute;
  top: 1rem;
  right: 1rem;
  width: 40px;
  height: 40px;
  border: 0;
  border-radius: 999px;
  background: rgba(255, 255, 255, 0.16);
  color: #fff;
  font-size: 28px;
  line-height: 1;
  cursor: pointer;
}

@keyframes lightbox-enter {
  from { opacity: 0; }
  to { opacity: 1; }
}
@keyframes content-enter {
  from { opacity: 0; transform: translateY(6px); }
  to { opacity: 1; transform: translateY(0); }
}

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
  scroll-margin-top: 104px;
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

@media (prefers-reduced-motion: reduce) {
  .post-loaded, .sidebar-results { animation: none; }
}
</style>
