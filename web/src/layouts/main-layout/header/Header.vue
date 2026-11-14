<template>
  <!--begin::Header-->
  <div
    id="kt_header"
    class="header align-items-stretch"
    :data-kt-sticky="isHeaderSticky"
    data-kt-sticky-name="header"
    data-kt-sticky-offset="{default: '200px', lg: '300px'}"
  >
    <!--begin::Container-->
    <div
      id="kt_header_container"
      :class="{
        'container-fluid': headerWidthFluid,
        'container-xxl': !headerWidthFluid,
      }"
      class="d-flex flex-stack flex-wrap gap-2"
    >
      <div class="d-flex d-lg-none align-items-center ms-n2 me-2">
        <!--begin::Aside mobile toggle-->
        <button
          type="button"
          class="btn btn-icon btn-active-icon-primary"
          id="kt_aside_toggle"
          aria-label="导航菜单"
          aria-controls="kt_aside"
          :aria-expanded="asideOpen"
        >
          <span class="svg-icon svg-icon-1 mt-1">
            <inline-svg
              :src="getAssetPath('/media/icons/duotune/abstract/abs015.svg')"
            />
          </span>
        </button>
        <!--end::Aside mobile toggle-->

        <!--begin::Logo-->
        <router-link
          :to="{ name: 'blog-home' }"
          class="studio-logo studio-logo-mobile"
          aria-label="返回博客首页"
        >
          studio.
        </router-link>
        <!--end::Logo-->
      </div>

      <PageTitle />

      <!--begin::Wrapper-->
      <div
        class="d-flex align-items-stretch justify-content-between flex-lg-grow-1"
      >
        <!--begin::Navbar-->
        <div class="d-flex align-items-stretch" id="kt_header_nav">
          <KTMenu v-if="authStore.isAuthenticated"></KTMenu>
        </div>
        <!--end::Navbar-->

        <!--begin::Topbar-->
        <div class="d-flex align-items-stretch flex-shrink-0">
          <KTTopbar></KTTopbar>
        </div>
        <!--end::Topbar-->
      </div>
      <!--end::Wrapper-->
    </div>
    <!--end::Container-->
  </div>
  <!--end::Header-->
</template>

<script lang="ts">
import { getAssetPath } from "@/core/helpers/assets";
import {
  defineComponent,
  computed,
  ref,
  onMounted,
  onBeforeUnmount,
} from "vue";
import KTTopbar from "@/layouts/main-layout/header/Topbar.vue";
import PageTitle from "@/layouts/main-layout/page-title/PageTitle.vue";
import KTMenu from "@/layouts/main-layout/header/Menu.vue";
import { useAuthStore } from "@/stores/auth";

import {
  headerWidthFluid,
  headerFixed,
  headerFixedOnMobile,
  headerLeft,
  asideDisplay,
} from "@/core/helpers/config";

export default defineComponent({
  name: "KTHeader",
  components: {
    KTTopbar,
    KTMenu,
    PageTitle,
  },
  setup() {
    const authStore = useAuthStore();
    const asideOpen = ref(false);
    let observer: MutationObserver | undefined;
    onMounted(() => {
      const sync = () => {
        asideOpen.value =
          document.body.getAttribute("data-kt-drawer-aside") === "on";
      };
      sync();
      observer = new MutationObserver(sync);
      observer.observe(document.body, {
        attributes: true,
        attributeFilter: ["data-kt-drawer-aside"],
      });
    });
    onBeforeUnmount(() => observer?.disconnect());
    const isHeaderSticky = computed(() => {
      if (window.innerWidth > 768) {
        return headerFixed.value;
      } else {
        return headerFixedOnMobile.value;
      }
    });

    return {
      asideOpen,
      headerWidthFluid,
      headerLeft,
      asideDisplay,
      isHeaderSticky,
      getAssetPath,
      authStore,
    };
  },
});
</script>

<style scoped>
.studio-logo-mobile {
  color: #181c32;
  font-size: 1.35rem;
  font-weight: 800;
  letter-spacing: -0.05em;
}

[data-theme="dark"] .studio-logo-mobile {
  color: #fff;
}
</style>
