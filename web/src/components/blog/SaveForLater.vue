<template>
  <span class="d-inline-flex flex-column align-items-start gap-1">
    <button
      type="button"
      class="btn btn-sm"
      :class="saved ? 'btn-light-primary' : 'btn-light'"
      :aria-pressed="saved"
      @click="toggle"
    >
      <i
        aria-hidden="true"
        class="bi"
        :class="saved ? 'bi-bookmark-fill' : 'bi-bookmark'"
      ></i
      >{{ saved ? "已收藏" : "稍后阅读" }}
    </button>
    <span
      v-if="message"
      class="fs-8 text-muted"
      role="status"
      aria-live="polite"
      >{{ message }}</span
    >
  </span>
</template>
<script setup lang="ts">
import { computed, ref } from "vue";
import { useReadingListStore } from "@/stores/readingList";
const props = defineProps<{
  article: { ID: number; fid: number; title: string };
}>();
const list = useReadingListStore();
const message = ref("");
const saved = computed(() => list.has(props.article.ID));
function toggle() {
  message.value = list.toggle(props.article);
}
</script>
