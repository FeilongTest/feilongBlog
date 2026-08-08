<template>
  <article class="card feed-card">
    <div class="card-body pb-5">
      <div class="d-flex align-items-center mb-5">
        <div class="symbol symbol-45px me-5">
          <img :src="getAssetPath('/media/avatars/300-1.jpg')" alt="飞龙" />
        </div>
        <div class="d-flex flex-column flex-grow-1 min-w-0">
          <router-link :to="articleUrl" class="feed-title text-gray-900 text-hover-primary fs-6 fw-bold">
            {{ article.title }}
          </router-link>
          <span class="text-gray-500 fw-semibold fs-7 mt-1">
            {{ dayjs.unix(Number(article.ctime)).format('YYYY-MM-DD HH:mm') }}
          </span>
        </div>
        <span v-if="article.isTop === 1" class="badge badge-light-success ms-3">置顶</span>
      </div>

      <router-link v-if="article.type !== 0 && article.pic" :to="articleUrl" class="feed-media d-block rounded overflow-hidden mb-5">
        <img v-if="!imageFailed" :src="imageUrl" class="feed-image" :alt="article.title" loading="lazy"
          @error="imageFailed = true" />
        <div v-else class="feed-image-fallback d-flex flex-column align-items-center justify-content-center">
          <i class="bi bi-image fs-2x text-primary mb-2"></i>
          <span class="text-gray-500 fs-7 fw-semibold">图片暂不可用</span>
        </div>
      </router-link>

      <div v-if="article.type === 0" class="feed-content-full text-gray-800 fs-6 fw-normal mb-5" v-html="sanitizedContent"></div>
      <p v-else class="feed-content text-gray-800 fs-6 fw-normal mb-5">{{ getExcerpt(article.content, 180) }}</p>

      <div class="d-flex align-items-center">
        <router-link :to="articleUrl" class="btn btn-sm btn-light btn-color-muted btn-active-light-primary px-4 py-2 me-4">
          <span class="svg-icon svg-icon-3 me-1">
            <inline-svg :src="getAssetPath('/media/icons/duotune/coding/cod006.svg')" />
          </span>
          {{ article.view }}
        </router-link>
        <router-link :to="articleUrl" class="btn btn-sm btn-light btn-color-muted btn-active-light-success px-4 py-2 me-4">
          <span class="svg-icon svg-icon-3 me-1">
            <inline-svg :src="getAssetPath('/media/icons/duotune/communication/com012.svg')" />
          </span>
          {{ article.commentCount || 0 }}
        </router-link>
        <button type="button" class="btn btn-sm btn-light btn-color-muted btn-active-light-danger px-4 py-2 like-button" :class="{ 'is-liked': liked, 'is-animating': liking }" :disabled="liking" @click="toggleLike">
          <span class="svg-icon svg-icon-3 like-icon">
            <inline-svg :src="getAssetPath('/media/icons/duotune/general/gen030.svg')" />
          </span>
          {{ likeCount }}
        </button>
      </div>
    </div>
  </article>
</template>

<script setup lang="ts">
import { computed, ref } from "vue";
import dayjs from "dayjs";
import type { ArticleList } from "@/core/blog/ArticleTypes";
import { getAssetPath } from "@/core/helpers/assets";
import { getExcerpt } from "@/utils/html";
import service from "@/utils/request";
import { getLikedArticles, getVisitorId, saveArticleLike } from "@/utils/visitor";
import DOMPurify from "dompurify";

const props = defineProps<{ article: ArticleList }>();
const imageFailed = ref(false);
const liked = ref(getLikedArticles().has(props.article.ID));
const liking = ref(false);
const likeCount = ref(props.article.likeCount || 0);
const sanitizedContent = computed(() => DOMPurify.sanitize(props.article.content || ""));
const articleUrl = computed(() => `/content/${props.article.fid}/${props.article.ID}`);
const imageUrl = computed(() => {
  const value = props.article.pic?.trim();
  if (!value || /^(https?:)?\/\//i.test(value)) return value;
  return value.startsWith("/") ? value : `/${value}`;
});

const toggleLike = async () => {
  if (liking.value) return;
  liking.value = true;
  try {
    const { data } = await service.post("/base/likeArticle", { aid: props.article.ID, visitorId: getVisitorId() });
    liked.value = Boolean(data?.liked);
    likeCount.value = Number(data?.count || 0);
    saveArticleLike(props.article.ID, liked.value);
  } finally {
    liking.value = false;
  }
};
</script>

<style scoped>
.feed-card { transition: box-shadow .2s ease; }
.feed-card:hover { box-shadow: var(--bs-card-box-shadow); }
.feed-title { display: -webkit-box; overflow: hidden; -webkit-line-clamp: 1; -webkit-box-orient: vertical; }
.feed-content { display: -webkit-box; overflow: hidden; -webkit-line-clamp: 2; -webkit-box-orient: vertical; line-height: 1.75; min-height: 3.5em; }
.feed-content-full { line-height: 1.75; overflow-wrap: anywhere; }
.feed-content-full :deep(p:last-child) { margin-bottom: 0; }
.feed-content-full :deep(img) { display: block; max-width: 100%; height: auto; border-radius: .625rem; margin: 1rem 0; }
.like-button .like-icon { transition: transform .18s ease; }
.like-button.is-liked { color: var(--kt-danger) !important; background-color: var(--kt-danger-light) !important; }
.like-button.is-liked .like-icon { color: var(--kt-danger) !important; }
.like-button.is-animating .like-icon { animation: heart-pop .32s ease; }
@keyframes heart-pop { 0% { transform: scale(1); } 45% { transform: scale(1.35) rotate(-8deg); } 100% { transform: scale(1); } }
.feed-media { background: #f5f8fa; }
.feed-image { display: block; width: 100%; max-height: 420px; min-height: 220px; object-fit: cover; }
.feed-image-fallback { width: 100%; min-height: 220px; background: linear-gradient(135deg, #f1f6ff, #f7f4ff); }
@media (prefers-reduced-motion: reduce) { .like-button .like-icon { transition: none; } .like-button.is-animating .like-icon { animation: none; } }
</style>
