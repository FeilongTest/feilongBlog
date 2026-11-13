<template>
  <div class="d-flex gap-2 mb-5" role="group" aria-label="预览尺寸">
    <button
      type="button"
      class="btn btn-sm"
      :class="mobile ? 'btn-light' : 'btn-primary'"
      @click="mobile = false"
    >
      桌面
    </button>
    <button
      type="button"
      class="btn btn-sm"
      :class="mobile ? 'btn-primary' : 'btn-light'"
      @click="mobile = true"
    >
      手机
    </button>
    <span class="text-muted fs-7 my-auto">仅供预览，未保存的内容不会公开</span>
  </div>
  <article
    class="preview-body border rounded p-5 mx-auto"
    :class="{ 'preview-mobile': mobile }"
  >
    <h1 class="fs-2 fw-bold text-gray-900">
      {{ article.title || "未命名草稿" }}
    </h1>
    <p v-if="article.excerpt" class="text-muted fs-5">{{ article.excerpt }}</p>
    <div v-if="article.tagNames?.length" class="d-flex flex-wrap gap-2 mb-5">
      <span
        v-for="tag in article.tagNames"
        :key="tag"
        class="badge badge-light-primary"
        >{{ tag }}</span
      >
    </div>
    <img
      v-if="article.pic"
      :src="resolveMediaUrl(article.pic)"
      alt="文章封面"
      class="mw-100 rounded mb-5"
    />
    <div class="preview-content" v-html="html"></div>
  </article>
</template>
<script setup lang="ts">
import { computed, ref } from "vue";
import DOMPurify from "dompurify";
import MarkdownIt from "markdown-it";
import { resolveMediaUrl } from "@/utils/media";
const props = defineProps<{
  article: {
    title?: string;
    content?: string;
    contentFormat?: string;
    excerpt?: string;
    pic?: string;
    tagNames?: string[];
  };
}>();
const mobile = ref(false);
const markdown = new MarkdownIt({ html: false, linkify: true, breaks: true });
const html = computed(() =>
  DOMPurify.sanitize(
    props.article.contentFormat === "markdown"
      ? markdown.render(props.article.content || "")
      : props.article.content || ""
  )
);
</script>
<style scoped>
.preview-body {
  width: 100%;
  transition: max-width 0.15s;
  background: var(--bs-body-bg);
}
.preview-mobile {
  max-width: 390px;
}
.preview-content {
  line-height: 1.9;
  overflow-wrap: anywhere;
}
.preview-content :deep(img) {
  max-width: 100%;
  height: auto;
}
.preview-content :deep(pre) {
  overflow: auto;
  padding: 1rem;
  border-radius: 0.5rem;
  background: var(--bs-gray-100);
}
.preview-content :deep(table) {
  display: block;
  overflow-x: auto;
  max-width: 100%;
}
</style>
