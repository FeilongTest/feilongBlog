<template>
  <!--begin::Footer-->
  <div class="footer py-4 d-flex flex-lg-column" id="kt_footer">
    <!--begin::Container-->
    <div
      class="d-flex flex-column flex-md-row align-items-center justify-content-between"
      :class="{
        'container-fluid': footerWidthFluid,
        'container-xxl': !footerWidthFluid,
      }"
    >
      <!--begin::Copyright-->
      <div class="text-dark order-2 order-md-1">
        <span class="text-muted fw-semobold me-1">2026&copy;</span>
        <a
          href="https://github.com/FeilongTest"
          target="_blank"
          class="text-gray-800 text-hover-primary"
          >飞龙Test</a
        >
      </div>
      <!--end::Copyright-->

      <!--begin::Menu-->
      <ul class="menu menu-gray-600 menu-hover-primary fw-semobold order-1">
        <li class="menu-item">
          <router-link to="/about" class="menu-link px-2">关于</router-link>
        </li>
        <li class="menu-item">
          <router-link to="/archive" class="menu-link px-2">归档</router-link>
        </li>
        <li class="menu-item">
          <details
            ref="readingTools"
            class="footer-reading-tools"
            @keydown.esc.prevent.stop="closeReadingTools(true)"
          >
            <summary class="menu-link px-2">
              阅读工具<span aria-hidden="true">⌃</span>
            </summary>
            <nav class="footer-reading-menu" aria-label="阅读工具">
              <router-link to="/tags" @click="closeReadingTools()"
                >标签专题</router-link
              >
              <router-link to="/reading-list" @click="closeReadingTools()"
                >稍后阅读</router-link
              >
            </nav>
          </details>
        </li>
        <li class="menu-item">
          <a href="/rss.xml" class="menu-link px-2">RSS</a>
        </li>
        <li class="menu-item d-lg-none">
          <button
            type="button"
            class="menu-link px-2 border-0 bg-transparent"
            @click="showLinks"
          >
            友情链接
          </button>
        </li>
        <li v-if="icp" class="menu-item">
          <span class="menu-link px-2 text-muted">{{ icp }}</span>
        </li>
      </ul>
      <!--end::Menu-->
    </div>
    <!--end::Container-->
  </div>
  <!--end::Footer-->
</template>

<script lang="ts">
import {
  computed,
  defineComponent,
  onMounted,
  onBeforeUnmount,
  ref,
} from "vue";
import { footerWidthFluid } from "@/core/helpers/config";
import { useSiteStore } from "@/stores/site";
import { DrawerComponent } from "@/assets/ts/components/_DrawerComponent";

export default defineComponent({
  name: "KTFooter",
  setup() {
    const readingTools = ref<HTMLDetailsElement>();
    const closeReadingTools = (restoreFocus = false) => {
      const tools = readingTools.value;
      if (!tools) return;
      tools.open = false;
      if (restoreFocus) tools.querySelector("summary")?.focus();
    };
    const closeOnOutside = (event: PointerEvent) => {
      if (
        event.target instanceof Node &&
        !readingTools.value?.contains(event.target)
      ) {
        closeReadingTools();
      }
    };
    const siteStore = useSiteStore();
    const icp = computed(() => siteStore.setting.icp);
    onMounted(() => {
      document.addEventListener("pointerdown", closeOnOutside);
      siteStore.loadSiteSetting().catch(() => undefined);
    });
    onBeforeUnmount(() =>
      document.removeEventListener("pointerdown", closeOnOutside)
    );
    const showLinks = () => {
      DrawerComponent.getInstance("kt_engage_demos_label")?.show();
    };
    return {
      readingTools,
      closeReadingTools,
      showLinks,
      footerWidthFluid,
      icp,
    };
  },
});
</script>

<style scoped>
.footer-reading-tools {
  position: relative;
}
.footer-reading-tools summary {
  display: flex;
  align-items: center;
  gap: 0.4rem;
  cursor: pointer;
  list-style: none;
}
.footer-reading-tools summary::-webkit-details-marker {
  display: none;
}
.footer-reading-tools summary span {
  font-size: 0.8em;
  opacity: 0.6;
}
.footer-reading-tools summary:focus-visible {
  outline: 2px solid var(--kt-primary);
  outline-offset: 2px;
  border-radius: 0.35rem;
}
.footer-reading-menu {
  position: absolute;
  bottom: calc(100% + 0.5rem);
  right: 0;
  z-index: 20;
  min-width: 10rem;
  padding: 0.4rem;
  border: 1px solid var(--kt-gray-200);
  border-radius: 0.65rem;
  background: var(--kt-body-bg);
  box-shadow: 0 0.5rem 1.5rem rgba(0, 0, 0, 0.08);
}
.footer-reading-menu a {
  display: block;
  padding: 0.65rem 0.8rem;
  border-radius: 0.4rem;
  color: var(--kt-gray-700);
  white-space: nowrap;
}
.footer-reading-menu a:hover,
.footer-reading-menu a:focus-visible {
  background: var(--kt-gray-100);
  color: var(--kt-primary);
}
</style>
