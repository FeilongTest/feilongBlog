<template>
  <!--begin::Menu-->
  <div
      id="kt_user_menu"
      class="menu menu-sub menu-sub-dropdown menu-column menu-rounded menu-gray-600 menu-state-bg-light-primary fw-semobold py-4 fs-6 w-275px"
      data-kt-menu="true"
      style="z-index: 1050;"
  >
    <!--begin::Menu item-->
    <div class="menu-item px-3">
      <div class="menu-content d-flex align-items-center px-3">
        <!--begin::Avatar-->
        <div class="symbol symbol-50px me-5">
          <img class="account-avatar-50" alt="" :src="avatarUrl" @error="useDefaultAvatar" />
        </div>
        <!--end::Avatar-->

        <!--begin::Username-->
        <div class="d-flex flex-column">
          <div class="fw-bold d-flex align-items-center fs-5">
            {{ displayName }}
          </div>
          <span v-if="displayMeta" class="fw-semobold text-muted fs-7 text-truncate">{{ displayMeta }}</span>
        </div>
        <!--end::Username-->
      </div>
    </div>
    <!--end::Menu item-->

    <!--begin::Menu separator-->
    <div class="separator my-2"></div>
    <!--end::Menu separator-->

    <div v-if="isAuthenticated" class="menu-item px-5">
      <router-link to="/admin/account/settings" class="menu-link px-5" @click="hideMenu">
        <span class="menu-icon"><i class="bi bi-shield-lock fs-5"></i></span>
        <span class="menu-text">账户设置</span>
      </router-link>
    </div>

    <!--begin::Menu item-->
    <div class="menu-item px-5">
      <a @click="handleSignOutOrLogin" class="menu-link px-5">
        <span class="menu-icon">
          <i :class="isAuthenticated ? 'bi bi-box-arrow-right fs-5' : 'bi bi-box-arrow-in-right fs-5'"></i>
        </span>
        <span class="menu-text">{{ isAuthenticated ? '退出登录' : '进入后台' }}</span>
      </a>
    </div>
    <!--end::Menu item-->
  </div>
  <!--end::Menu-->
</template>

<script lang="ts">
import { getAssetPath } from "@/core/helpers/assets";
import { defineComponent, computed } from "vue";
import { useAuthStore } from "@/stores/auth";
import { useRouter } from "vue-router";
import { MenuComponent } from "@/assets/ts/components";

export default defineComponent({
  name: "kt-user-menu",
  components: {},
  setup() {
    const router = useRouter();
    const store = useAuthStore();
    const profile = computed(() => store.user?.user);
    const publicProfile = computed(() => store.publicProfile);
    const defaultAvatar = getAssetPath("/media/avatars/300-1.jpg");
    const avatarUrl = computed(() => (store.isAuthenticated ? profile.value?.pic : publicProfile.value?.pic) || defaultAvatar);
    const displayName = computed(() => store.isAuthenticated
      ? (profile.value?.truename || profile.value?.userName || "飞龙Test")
      : (publicProfile.value?.trueName || "飞龙Test"));
    const displayMeta = computed(() => store.isAuthenticated
      ? (profile.value?.bio || "Golang Dev")
      : (publicProfile.value?.bio || "Golang Dev"));
    const useDefaultAvatar = (event: Event) => {
      const image = event.currentTarget as HTMLImageElement;
      if (image.src.endsWith(defaultAvatar)) return;
      image.src = defaultAvatar;
    };

    const hideMenu = () => {
      const menuElement = document.getElementById("kt_user_menu");
      menuElement?.classList.remove("show");
      menuElement?.parentElement?.classList.remove("show");
      MenuComponent.hideDropdowns(undefined);
    };

    // 处理登录/退出
    const handleSignOutOrLogin = (e: Event) => {
      e.preventDefault();
      e.stopPropagation();
      
      // 隐藏菜单 - 通过触发器元素找到菜单并隐藏
      const triggerElement = document.querySelector('[data-kt-menu-trigger="click"][data-kt-menu-placement="top-end"]');
      if (triggerElement) {
        // 查找菜单元素（通过触发器查找）
        const menuElement = document.getElementById('kt_user_menu');
        if (menuElement) {
          // 尝试通过 MenuComponent 获取实例并隐藏
          const menuInstance = MenuComponent.getInstance(menuElement as HTMLElement);
          if (menuInstance) {
            // 查找触发器对应的菜单项
            const menuItems = menuElement.querySelectorAll('[data-kt-menu-trigger]');
            if (menuItems.length > 0) {
              menuInstance.hide(menuItems[0] as HTMLElement);
            }
          }
          
          // 直接移除show类
          menuElement.classList.remove('show');
          menuElement.parentElement?.classList.remove('show');
        }
        
        // 隐藏所有下拉菜单
        MenuComponent.hideDropdowns(undefined);
        
        // 移除所有相关的show类
        document.querySelectorAll('.show.menu-dropdown').forEach(el => {
          el.classList.remove('show');
        });
      }

      // 延迟执行路由跳转，确保菜单先隐藏
      setTimeout(() => {
        if (store.isAuthenticated) {
          // 已登录：退出登录
          store.logout();
          router.push({ name: "blog-home" });
        } else {
          // 未登录：跳转到登录页
          router.push({ name: "sign-in" });
        }
      }, 150);
    };

    const isAuthenticated = computed(() => {
      return store.isAuthenticated;
    });

    return {
      handleSignOutOrLogin,
      getAssetPath,
      isAuthenticated,
      hideMenu,
      avatarUrl,
      displayName,
      displayMeta,
      useDefaultAvatar,
    };
  },
});
</script>

<style scoped>
/* 确保菜单在移动端正确显示，不被分类菜单遮挡 */
.menu {
  position: relative;
  z-index: 1050 !important;
}

.menu-sub {
  position: relative;
  z-index: 1051 !important;
}

.account-avatar-50 {
  display: block;
  width: 50px !important;
  height: 50px !important;
  max-width: none;
  object-fit: cover;
  object-position: center;
  border-radius: 50%;
}

/* 移动端特殊处理 */
@media (max-width: 991.98px) {
  .menu {
    z-index: 1055 !important;
  }
  
  .menu-sub {
    z-index: 1056 !important;
  }
}
</style>
