<template>
  <RouterView />
</template>

<script lang="ts">
import { defineComponent, nextTick, onBeforeMount, onBeforeUnmount, onMounted } from "vue";
import { RouterView } from "vue-router";
import { useRouter } from "vue-router";
import { useConfigStore } from "@/stores/config";
import { useThemeStore } from "@/stores/theme";
import { useBodyStore } from "@/stores/body";
import { themeMode } from "@/core/helpers/config";
import { initializeComponents } from "@/core/plugins/keenthemes";
import { getVisitorId } from "@/utils/visitor";

export default defineComponent({
  name: "app",
  components: {
    RouterView,
  },
  setup() {
    const configStore = useConfigStore();
    const themeStore = useThemeStore();
    const bodyStore = useBodyStore();
    const router = useRouter();
    let lastRecordedPath = "";
    const recordVisit = (path: string) => {
      if (path === lastRecordedPath || path.startsWith("/admin") || path.startsWith("/sign-in") || path.startsWith("/error")) return;
      lastRecordedPath = path;
      const baseURL = import.meta.env.VITE_BLOG_API_URL || "/blog";
      fetch(`${baseURL}/base/visit`, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ visitorId: getVisitorId() }),
        keepalive: true,
      }).catch(() => undefined);
    };
    const removeAfterEach = router.afterEach((to) => recordVisit(to.path));

    onBeforeMount(() => {
      /**
       * Overrides the layout config using saved data from localStorage
       * remove this to use static config (@/core/config/DefaultLayoutConfig.ts)
       */
      configStore.overrideLayoutConfig();

      /**
       *  Sets a mode from configuration
       */
      themeStore.setThemeMode(themeMode.value);
    });

    onMounted(() => {
      recordVisit(router.currentRoute.value.path);
      nextTick(() => {
        initializeComponents();

        bodyStore.removeBodyClassName("page-loading");
      });
    });
    onBeforeUnmount(removeAfterEach);
  },
});
</script>

<style lang="scss">
html,
body,
#app {
  max-width: 100%;
  overflow-x: hidden;
}

@import "bootstrap-icons/font/bootstrap-icons.css";
@import "animate.css";
@import "sweetalert2/dist/sweetalert2.css";
@import "prism-themes/themes/prism-shades-of-purple.css";
@import "element-plus/dist/index.css";

// Main demo style scss
@import "assets/sass/element-ui.dark";
@import "assets/sass/plugins";
@import "assets/sass/style";

#app {
  display: contents;
}
</style>
