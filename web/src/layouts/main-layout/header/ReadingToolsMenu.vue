<template>
  <details
    ref="menu"
    class="reading-tools"
    @keydown.esc.prevent.stop="closeMenu(true)"
  >
    <summary
      class="btn btn-icon btn-color-gray-700 btn-active-color-primary btn-outline w-40px h-40px"
      aria-label="阅读工具"
      title="阅读工具"
    >
      <i class="bi bi-book fs-4" aria-hidden="true"></i>
    </summary>
    <nav class="reading-tools-panel" aria-label="阅读工具">
      <router-link to="/tags" @click="closeMenu()">
        <i class="bi bi-tags" aria-hidden="true"></i>
        <span>标签专题</span>
      </router-link>
      <router-link
        to="/reading-list"
        aria-label="稍后阅读清单"
        @click="closeMenu()"
      >
        <i class="bi bi-bookmark" aria-hidden="true"></i>
        <span>稍后阅读</span>
      </router-link>
    </nav>
  </details>
</template>

<script setup lang="ts">
import { onBeforeUnmount, onMounted, ref, watch } from "vue";
import { useRoute } from "vue-router";

const menu = ref<HTMLDetailsElement>();
const route = useRoute();
function closeMenu(restoreFocus = false) {
  const element = menu.value;
  if (!element) return;
  element.open = false;
  if (restoreFocus) element.querySelector("summary")?.focus();
}
function closeOnOutside(event: PointerEvent) {
  if (event.target instanceof Node && !menu.value?.contains(event.target)) {
    closeMenu();
  }
}
watch(
  () => route.fullPath,
  () => closeMenu()
);
onMounted(() => document.addEventListener("pointerdown", closeOnOutside));
onBeforeUnmount(() =>
  document.removeEventListener("pointerdown", closeOnOutside)
);
</script>

<style scoped>
.reading-tools {
  position: relative;
}
.reading-tools summary {
  display: flex;
  align-items: center;
  justify-content: center;
  list-style: none;
  cursor: pointer;
}
.reading-tools summary::-webkit-details-marker {
  display: none;
}
.reading-tools summary:focus-visible {
  outline: 2px solid var(--kt-primary);
  outline-offset: 2px;
}
.reading-tools[open] summary {
  color: var(--kt-primary);
  border-color: var(--kt-primary);
}
.reading-tools-panel {
  position: absolute;
  top: calc(100% + 0.6rem);
  right: 0;
  z-index: 1050;
  width: 11rem;
  padding: 0.45rem;
  border: 1px solid var(--kt-gray-200);
  border-radius: 0.75rem;
  background: var(--kt-body-bg);
  box-shadow: 0 0.5rem 1.5rem rgba(0, 0, 0, 0.08);
}
.reading-tools-panel a {
  display: flex;
  align-items: center;
  gap: 0.75rem;
  padding: 0.75rem;
  border-radius: 0.4rem;
  color: var(--kt-gray-700);
  white-space: nowrap;
}
.reading-tools-panel i {
  color: var(--kt-gray-500);
}
.reading-tools-panel a:hover,
.reading-tools-panel a:focus-visible {
  background: var(--kt-gray-100);
  color: var(--kt-primary);
}
</style>
