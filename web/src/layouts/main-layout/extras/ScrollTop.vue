<template>
  <!-- begin::Scrolltop -->
  <div
    id="kt_scrolltop"
    ref="kt_scrolltop"
    class="scrolltop"
    data-kt-scrolltop="true"
  >
    <span class="svg-icon">
      <inline-svg
        :src="getAssetPath('/media/icons/duotune/arrows/arr066.svg')"
      />
    </span>
  </div>
  <!-- end::Scrolltop -->
</template>

<script lang="ts">
import { getAssetPath } from "@/core/helpers/assets";
import { defineComponent, onBeforeUnmount, onMounted } from "vue";
import { useRouter } from "vue-router";
import { ScrollTopComponent } from "@/assets/ts/components/_ScrollTopComponent";

const SHOW_AFTER = 240;

const pageScrollTop = () =>
  Math.max(
    window.scrollY || 0,
    document.documentElement?.scrollTop || 0,
    document.body?.scrollTop || 0,
  );

const updateScrollTopButton = () => {
  if (pageScrollTop() > SHOW_AFTER) document.body.setAttribute("data-kt-scrolltop", "on");
  else document.body.removeAttribute("data-kt-scrolltop");
};

const scrollPageToTop = () => {
  window.scrollTo({ top: 0, behavior: "smooth" });
  document.documentElement.scrollTo?.({ top: 0, behavior: "smooth" });
  document.body.scrollTo?.({ top: 0, behavior: "smooth" });
};

export default defineComponent({
  name: "KTScrollTop",
  components: {},
  setup() {
    const router = useRouter();
    const onScroll = (event: Event) => {
      const target = event.target;
      const isPage =
        target === document ||
        target === document.documentElement ||
        target === document.body;
      if (!isPage) return;
      updateScrollTopButton();
    };

    onMounted(() => {
      ScrollTopComponent.reinitialization();
      updateScrollTopButton();
      window.addEventListener("scroll", updateScrollTopButton, { passive: true });
      document.addEventListener("scroll", onScroll, { capture: true, passive: true });
      document.getElementById("kt_scrolltop")?.addEventListener("click", scrollPageToTop);
    });

    const removeAfterEach = router.afterEach(() => {
      window.setTimeout(updateScrollTopButton, 0);
    });

    onBeforeUnmount(() => {
      window.removeEventListener("scroll", updateScrollTopButton);
      document.removeEventListener("scroll", onScroll, true);
      document.getElementById("kt_scrolltop")?.removeEventListener("click", scrollPageToTop);
      removeAfterEach();
    });

    return {
      getAssetPath,
    };
  },
});
</script>

<style>
.scrolltop {
  z-index: 1090;
}
</style>
