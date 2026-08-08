<template>
  <!--begin::Toolbar wrapper-->
  <div class="d-flex align-items-stretch flex-shrink-0">
    <div class="topbar d-flex align-items-stretch flex-shrink-0">
      <!--begin::Search-->
      <div class="d-flex align-items-center ms-1 ms-lg-3">
        <KTSearch></KTSearch>
      </div>
      <!--end::Search-->

      <!--begin::Theme mode-->
      <div class="d-flex align-items-center ms-3 ms-lg-4">
        <!--begin::Menu toggle-->
        <a
          href="#"
          class="btn btn-icon btn-color-gray-700 btn-active-color-primary btn-outline w-40px h-40px"
          data-kt-menu-trigger="{default:'click', lg: 'hover'}"
          data-kt-menu-attach="parent"
          data-kt-menu-placement="bottom-end"
        >
          <span class="svg-icon theme-light-show svg-icon-1">
            <inline-svg
              :src="getAssetPath('/media/icons/duotune/general/gen060.svg')"
            />
          </span>
          <span class="svg-icon theme-dark-show svg-icon-1">
            <inline-svg
              :src="getAssetPath('/media/icons/duotune/general/gen061.svg')"
            />
          </span>
        </a>
        <!--begin::Menu toggle-->
        <KTThemeModeSwitcher></KTThemeModeSwitcher>
      </div>
      <!--end::Theme mode-->

      <!-- <button
        class="btn btn-icon btn-active-icon-primary w-40px h-40px d-xxl-none ms-2 me-n2"
        id="kt_sidebar_toggler"
      >
        <span class="svg-icon svg-icon-2x">
          <inline-svg
            :src="getAssetPath('/media/icons/duotune/coding/cod001.svg')"
          />
        </span>
      </button>
      删除多余按钮 -->
    </div>
  </div>
  <!--end::Toolbar wrapper-->
</template>

<script lang="ts">
import { getAssetPath } from "@/core/helpers/assets";
import { defineComponent } from "vue";
import KTSearch from "@/layouts/main-layout/search/Search.vue";
import KTThemeModeSwitcher from "@/layouts/main-layout/theme-mode/ThemeModeSwitcher.vue";

export default defineComponent({
  name: "layout-topbar",
  components: {
    KTSearch,
    KTThemeModeSwitcher,
  },
  setup() {
    return {
      getAssetPath,
    };
  },
});
</script>
