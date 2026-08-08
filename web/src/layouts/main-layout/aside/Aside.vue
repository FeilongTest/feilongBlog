<template>
  <!--begin::Aside-->
  <div
    id="kt_aside"
    class="aside py-9"
    data-kt-drawer="true"
    data-kt-drawer-name="aside"
    data-kt-drawer-activate="{default: true, lg: false}"
    data-kt-drawer-overlay="true"
    data-kt-drawer-width="{default:'200px', '300px': '250px'}"
    data-kt-drawer-direction="start"
    data-kt-drawer-toggle="#kt_aside_toggle"
  >
    <!--begin::Aside logo-->
    <div class="aside-logo flex-column-auto px-9 mb-9" id="kt_aside_logo">
      <!--begin::Logo-->
      <router-link :to="{ name: 'blog-home' }" class="studio-logo" aria-label="返回博客首页">
        studio.
      </router-link>
      <!--end::Logo-->
    </div>
    <!--end::Aside logo-->

    <!--begin::Aside menu-->
    <div class="aside-menu flex-column-fluid ps-5 pe-3 overflow-hidden">
      <KTMenu v-if="!authStore.isAuthenticated"></KTMenu>
      <KTAdminMenu v-if="authStore.isAuthenticated"></KTAdminMenu>
    </div>
    <!--end::Aside menu-->

    <!--begin::Footer-->
    <div class="aside-footer flex-column-auto px-9" id="kt_aside_footer">
      <!--begin::User panel-->
      <div class="d-flex flex-stack">
        <!--begin::Wrapper-->
        <div class="d-flex align-items-center">
          <!--begin::Avatar-->
          <div class="symbol symbol-circle symbol-40px">
            <img class="account-avatar account-avatar-40" :src="avatarUrl" alt="" @error="useDefaultAvatar" />
          </div>
          <!--end::Avatar-->

          <!--begin::User info-->
          <div class="ms-2">
            <!--begin::Name-->
            <a
              href="#"
              class="text-gray-800 text-hover-primary fs-6 fw-bold lh-1"
              >{{ displayName }}</a
            >
            <!--end::Name-->

            <!--begin::Major-->
            <span v-if="displayMeta" class="aside-user-meta text-muted fw-semobold d-block fs-7 mt-1"
              >{{ displayMeta }}</span
            >
            <!--end::Major-->
          </div>
          <!--end::User info-->
        </div>
        <!--end::Wrapper-->

        <!--begin::User menu-->
        <div class="ms-1">
          <div
            class="btn btn-sm btn-icon btn-active-color-primary position-relative me-n2"
            data-kt-menu-trigger="click"
            data-kt-menu-overflow="true"
            data-kt-menu-placement="top-end"
          >
            <span class="svg-icon svg-icon-1">
              <inline-svg
                :src="getAssetPath('/media/icons/duotune/coding/cod001.svg')"
              />
            </span>
          </div>
          <UserMenu />
        </div>
        <!--end::User menu-->
      </div>
      <!--end::User panel-->
    </div>
    <!--end::Footer-->
  </div>
  <!--end::Aside-->
</template>

<script lang="ts">
import { getAssetPath } from "@/core/helpers/assets";
import { computed, defineComponent, onMounted } from "vue";
import { useI18n } from "vue-i18n";
import { useAuthStore } from "@/stores/auth";
import KTAdminMenu from "@/layouts/main-layout/aside/AdminMenu.vue";
import KTMenu from "@/layouts/main-layout/aside/Menu.vue";
import UserMenu from "@/layouts/main-layout/menus/UserAccountMenu.vue";
import { asideTheme } from "@/core/helpers/config";

export default defineComponent({
  name: "KTAside",
  components: {
    KTMenu,
    KTAdminMenu,
    UserMenu,
  },
  props: {
    lightLogo: String,
    darkLogo: String,
  },
  setup() {
    const { t } = useI18n();
    const authStore = useAuthStore();
    const profile = computed(() => authStore.user?.user);
    const publicProfile = computed(() => authStore.publicProfile);
    const defaultAvatar = getAssetPath("/media/avatars/300-1.jpg");
    const avatarUrl = computed(() => (authStore.isAuthenticated ? profile.value?.pic : publicProfile.value?.pic) || defaultAvatar);
    const displayName = computed(() => authStore.isAuthenticated
      ? (profile.value?.truename || profile.value?.userName || "飞龙Test")
      : (publicProfile.value?.trueName || "飞龙Test"));
    const displayMeta = computed(() => authStore.isAuthenticated
      ? (profile.value?.bio || "Golang Dev")
      : (publicProfile.value?.bio || "Golang Dev"));
    const useDefaultAvatar = (event: Event) => {
      const image = event.currentTarget as HTMLImageElement;
      if (image.src.endsWith(defaultAvatar)) return;
      image.src = defaultAvatar;
    };
    onMounted(() => {
      if (authStore.isAuthenticated && !profile.value) authStore.refreshProfile().catch(() => undefined);
      if (!authStore.isAuthenticated && !publicProfile.value?.bio) authStore.refreshPublicProfile().catch(() => undefined);
    });

    return {
      authStore,
      asideTheme,
      t,
      getAssetPath,
      avatarUrl,
      displayName,
      displayMeta,
      useDefaultAvatar,
    };
  },
});
</script>

<style scoped>
.studio-logo {
  color: #181c32;
  font-size: 1.65rem;
  font-weight: 800;
  letter-spacing: -.055em;
}

[data-theme="dark"] .studio-logo { color: #fff; }
.aside-user-meta { line-height: 1.35; padding-bottom: 1px; }
.account-avatar-40 { display: block; width: 40px !important; height: 40px !important; max-width: none; object-fit: cover; object-position: center; border-radius: 50%; }
</style>
