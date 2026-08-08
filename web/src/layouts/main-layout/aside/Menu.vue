<template>
  <!--begin::Menu wrapper-->
  <div id="kt_aside_menu_wrapper" ref="scrollElRef" class="w-100 h-100 hover-scroll-overlay-y d-flex pe-2" data-kt-scroll="true"
    data-kt-scroll-activate="{default: false, lg: true}" data-kt-scroll-height="auto"
    data-kt-scroll-dependencies="#kt_aside_logo, #kt_aside_footer"
    data-kt-scroll-wrappers="#kt_aside, #kt_aside_menu, #kt_aside_menu_wrapper" data-kt-scroll-offset="100">
    <!--begin::Menu-->
    <div class="menu menu-column menu-rounded menu-sub-indention menu-active-bg fw-semibold my-auto w-100" id="kt_aside_menu"
      data-kt-menu="true">

      <div class="menu-item">
        <router-link :to="{ name: 'blog-home' }" class="menu-link" active-class="active">
          <span class="menu-icon">
            <span class="svg-icon svg-icon-5">
              <inline-svg :src="getAssetPath('/media/icons/duotune/arrows/arr001.svg')" />
            </span>
          </span>
          <span class="menu-title">首页</span>
        </router-link>
      </div>

      <template v-for="(item, i) in category" :key="i">
        <div v-if="item.fid == 0" :class="{ show: hasActiveChildren(item.ID) }" class="menu-item menu-accordion"
          data-kt-menu-sub="accordion" data-kt-menu-trigger="click">
          <span class="menu-link">
            <span class="menu-icon">
              <span class="svg-icon svg-icon-5">
                <inline-svg :src="getAssetPath('/media/icons/duotune/arrows/arr001.svg')" />
              </span>
            </span>
            <span class="menu-title">{{ item.name }}</span>
            <span class="menu-arrow"></span>
          </span>

          <div class="menu-sub menu-sub-accordion">
            <template v-for="(menuItem, j) in category" :key="j">
              <div v-if="menuItem.fid === item.ID" class="menu-item">
                <router-link :to="'/category/' + menuItem.ID" class="menu-link" active-class="active"
                  :class="{ active: hasActiveSubChildren(menuItem.ID) }">
                  <span class="menu-bullet">
                    <span class="bullet bullet-dot"></span>
                  </span>
                  <span class="menu-title">{{ menuItem.name }}</span>
                </router-link>
              </div>
            </template>
          </div>
        </div>
      </template>

    </div>
    <!--end::Menu-->
  </div>
  <!--end::Menu wrapper-->
</template>

<script lang="ts">
import { getAssetPath } from "@/core/helpers/assets";
import { defineComponent, nextTick, onMounted, ref, watch } from "vue";
import { useRoute } from "vue-router";
import { asideMenuIcons } from "@/core/helpers/config";
import service from "@/utils/request";
import type { CategoryList } from "@/core/blog/CategoryTypes"
import { MenuComponent } from "@/assets/ts/components/MenuComponent";


export default defineComponent({
  name: "kt-menu",
  components: {},
  setup() {
    const route = useRoute();
    const scrollElRef = ref<null | HTMLElement>(null);
    const category = ref<Array<CategoryList>>([])

    onMounted(() => {
      if (scrollElRef.value) {
        scrollElRef.value.scrollTop = 0;
      }
      getMenuList();
    });

    const getMenuList = () => {
      service.get("/base/getCategoryList")
        .then(res => {
          category.value = res.data
        })
        .catch(err => {
          console.log(err)
        })
    }

    //重写是否拥有激活的子元素判断条件
    const hasActiveChildren = (match: number) => {
      let result = false;
      for (let index = 0; index < category.value.length; index++) {
        const item = category.value[index];
        if (item.fid === match && Number(route.path.match(/\d+/g)?.[0]) === item.ID) {
          result = true;
          break;
        }
      }
      return result;
    }

    //判断子元素是否激活
    const hasActiveSubChildren = (match: number) => {
      if (match === Number(route.path.match(/\d+/g)?.[0])) {
        return true;
      }
      return false;
    }

    const collapseAccordions = () => {
      const menuElement = document.getElementById("kt_aside_menu");
      if (!menuElement) return;
      const menu = MenuComponent.getInstance(menuElement);
      menuElement.querySelectorAll<HTMLElement>(".menu-item.show[data-kt-menu-trigger]")
        .forEach((item) => menu?.hide(item));
    };

    watch(() => route.name, async (routeName) => {
      if (routeName === "blog-home") {
        collapseAccordions();
        return;
      }
      if (routeName === "category") {
        await nextTick();
        document.querySelectorAll<HTMLElement>("#kt_aside_menu .menu-item.menu-accordion.show > .menu-sub")
          .forEach((subMenu) => subMenu.removeAttribute("style"));
      }
    }, { flush: "sync" });

    return {
      category,
      hasActiveChildren,
      hasActiveSubChildren,
      asideMenuIcons,
      getAssetPath,
    };
  },
});
</script>
