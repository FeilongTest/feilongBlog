<template>
  <!--begin::Demos drawer-->
  <div
    id="kt_engage_demos_label"
    class="bg-body"
    data-kt-drawer="true"
    data-kt-drawer-name="explore"
    data-kt-drawer-activate="true"
    data-kt-drawer-overlay="true"
    data-kt-drawer-width="{default:'calc(100vw - 20px)', 'sm': '350px', 'lg': '475px'}"
    data-kt-drawer-direction="end"
    data-kt-drawer-toggle="#kt_engage_demos_toggle"
    data-kt-drawer-close="#kt_engage_demos_close"
  >
    <!--begin::Card-->
    <div class="card shadow-none w-100">
      <!--begin::Header-->
      <div class="card-header" id="kt_explore_header">
        <h3 class="card-title fw-bold text-gray-700">
          <i class="bi bi-link-45deg fs-2 me-2"></i>友情链接
        </h3>

        <div class="card-toolbar">
          <button
            type="button"
            class="btn btn-sm btn-icon btn-active-light-primary me-n5"
            id="kt_engage_demos_close"
          >
            <span class="svg-icon svg-icon-2">
              <inline-svg
                :src="getAssetPath('/media/icons/duotune/arrows/arr061.svg')"
              />
            </span>
          </button>
        </div>
      </div>
      <!--end::Header-->

      <!--begin::Body-->
      <div class="card-body" id="kt_explore_body">
        <!--begin::Content-->
        <div
          id="kt_explore_scroll"
          class="friend-links-scroll scroll-y me-n5 pe-5"
          data-kt-scroll="true"
          data-kt-scroll-height="auto"
          data-kt-scroll-wrappers="#kt_explore_body"
          data-kt-scroll-dependencies="#kt_explore_header"
          data-kt-scroll-offset="5px"
        >
          <!--begin::Wrapper-->
          <div class="mb-0">
            <!--begin::Description-->
            <div class="text-center mb-6">
              <p class="text-gray-600 fs-6 mb-3">
                欢迎交换友情链接！如需添加，请联系我。
              </p>
              <div class="d-flex flex-column gap-2">
                <a v-if="contactEmail" :href="`mailto:${contactEmail}`" class="text-primary fs-7 text-break">
                  <i class="bi bi-envelope fs-6 me-1"></i>{{ contactEmail }}
                </a>
                <span v-else class="text-muted fs-7"><i class="bi bi-envelope fs-6 me-1"></i>暂未设置联系邮箱</span>
                <a href="https://github.com/FeilongTest" target="_blank" class="text-primary fs-7">
                  <i class="bi bi-github fs-6 me-1"></i>GitHub
                </a>
              </div>
            </div>
            <!--end::Description-->

            <!--begin::友情链接列表-->
            <div class="mb-0">
              <h5 class="fw-bold mb-4 text-gray-800 fs-6">
                <i class="bi bi-star-fill text-warning me-2 fs-7"></i>推荐网站
              </h5>

              <!--begin::Loading-->
              <div v-if="loading" class="text-center py-10">
                <div class="spinner-border text-primary" role="status">
                  <span class="visually-hidden">加载中...</span>
                </div>
              </div>
              <!--end::Loading-->

              <!--begin::友情链接-->
              <template v-else-if="friendLinks.length > 0">
                <template v-for="link in friendLinks" :key="link.ID">
                  <div class="friend-link-card mb-3 p-3 rounded border border-gray-300 hover-elevate-up">
                    <!--begin::Icon-->
                    <div class="symbol symbol-35px me-3 flex-shrink-0">
                      <img v-if="link.logo" :src="link.logo" :alt="link.name" class="w-100 h-100 rounded" />
                      <div v-else class="symbol-label fs-4 fw-semobold bg-light-primary text-primary">
                        {{ link.name.charAt(0) }}
                      </div>
                    </div>
                    <!--end::Icon-->
                    
                    <!--begin::Info-->
                    <div class="friend-link-info">
                      <a :href="link.url" target="_blank" class="text-gray-800 text-hover-primary fw-bold fs-7 mb-1 d-block text-truncate">
                        {{ link.name }}
                      </a>
                      <div class="text-muted fs-8 text-truncate" :title="link.desc">{{ link.desc || '暂无描述' }}</div>
                    </div>
                    <!--end::Info-->

                    <!--begin::Link-->
                    <a :href="link.url" target="_blank" class="friend-link-action btn btn-xs btn-light-primary">
                      访问
                    </a>
                    <!--end::Link-->
                  </div>
                </template>
              </template>
              <!--end::友情链接-->

              <!--begin::Empty-->
              <div v-else class="text-center py-10">
                <i class="bi bi-link-45deg fs-3x text-muted mb-3"></i>
                <p class="text-muted mb-0">暂无友情链接</p>
              </div>
              <!--end::Empty-->

              <!--begin::添加提示-->
              <div class="text-center mt-5 p-4 bg-light-info rounded">
                <i class="bi bi-info-circle fs-3 text-info mb-2"></i>
                <p class="text-gray-700 mb-0 fs-7">
                  想要添加您的网站？<br>
                  <a href="https://github.com/FeilongTest" target="_blank" class="text-primary fw-bold">联系我</a>
                </p>
              </div>
              <!--end::添加提示-->
            </div>
            <!--end::友情链接列表-->
          </div>
          <!--end::Wrapper-->
        </div>
        <!--end::Content-->
      </div>
      <!--end::Body-->
    </div>
    <!--end::Card-->
  </div>
  <!--end::Demos drawer-->
</template>

<script lang="ts">
import { getAssetPath } from "@/core/helpers/assets";
import { computed, defineComponent, onBeforeUnmount, onMounted, ref } from "vue";
import service from "@/utils/request";
import type { Friendlink } from "@/core/blog/FriendlinkTypes";
import { useAuthStore } from "@/stores/auth";

export default defineComponent({
  name: "kt-demos-drawer",
  setup() {
    const friendLinks = ref<Friendlink[]>([]);
    const loading = ref(false);
    const loaded = ref(false);
    const authStore = useAuthStore();
    const contactEmail = computed(() => authStore.publicProfile?.email || "");

    // 获取友情链接列表
    const getFriendlinks = async () => {
      if (loaded.value || loading.value) return;
      loading.value = true;
      try {
        const [res] = await Promise.all([
          service.get("/base/getAllFriendlinks"),
          authStore.refreshPublicProfile().catch(() => undefined),
        ]) as [any, unknown];
        if (res && res.data && Array.isArray(res.data)) {
          friendLinks.value = res.data;
        }
        loaded.value = true;
      } catch (error) {
        console.error("获取友链失败:", error);
      } finally {
        loading.value = false;
      }
    };

    const handleDrawerOpen = () => getFriendlinks();

    onMounted(() => {
      document.getElementById("kt_engage_demos_toggle")?.addEventListener("click", handleDrawerOpen);
    });

    onBeforeUnmount(() => {
      document.getElementById("kt_engage_demos_toggle")?.removeEventListener("click", handleDrawerOpen);
    });

    return {
      friendLinks,
      loading,
      contactEmail,
      getAssetPath,
    };
  },
});
</script>

<style scoped>
.friend-links-scroll {
  max-width: 100%;
  overflow-x: hidden;
}

.friend-link-card {
  display: grid;
  grid-template-columns: 35px minmax(0, 1fr) auto;
  align-items: center;
  column-gap: 0.75rem;
  width: 100%;
  min-width: 0;
  overflow: hidden;
}

.friend-link-info {
  min-width: 0;
  overflow: hidden;
}

.friend-link-action {
  flex-shrink: 0;
  margin-left: 0.25rem;
}

.hover-elevate-up {
  transition: all 0.2s ease;
}

.hover-elevate-up:hover {
  transform: translateY(-1px);
  box-shadow: 0 0.25rem 0.75rem rgba(0, 0, 0, 0.1);
}

.text-truncate {
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

/* 确保友链卡片更紧凑 */
.symbol-35px {
  width: 35px;
  height: 35px;
  min-width: 35px;
}

@media (max-width: 575.98px) {
  #kt_explore_body {
    padding: 1.25rem;
  }

  .friend-links-scroll {
    margin-right: 0 !important;
    padding-right: 0 !important;
  }

  .friend-link-card {
    column-gap: 0.625rem;
  }

  .friend-link-action {
    padding-right: 0.75rem;
    padding-left: 0.75rem;
  }
}
</style>
