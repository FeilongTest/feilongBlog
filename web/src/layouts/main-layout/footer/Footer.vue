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
          <a href="/rss.xml" class="menu-link px-2">RSS</a>
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
import { computed, defineComponent, onMounted } from "vue";
import { footerWidthFluid } from "@/core/helpers/config";
import { useSiteStore } from "@/stores/site";

export default defineComponent({
  name: "KTFooter",
  setup() {
    const siteStore = useSiteStore();
    const icp = computed(() => siteStore.setting.icp);
    onMounted(() => {
      siteStore.loadSiteSetting().catch(() => undefined);
    });
    return {
      footerWidthFluid,
      icp,
    };
  },
});
</script>
