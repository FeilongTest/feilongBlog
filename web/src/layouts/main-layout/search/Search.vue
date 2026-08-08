<template>
  <MenuComponent menu-selector="#kt-search-menu">
    <template v-slot:toggle>
      <!--begin::Search-->
      <div
        id="kt_header_search"
        class="d-flex align-items-stretch"
        data-kt-menu-target="#kt-search-menu"
        data-kt-menu-trigger="click"
        data-kt-menu-attach="parent"
        data-kt-menu-placement="bottom-end"
        data-kt-menu-flip="bottom"
      >
        <!--begin::Search toggle-->
        <div class="d-flex align-items-center" id="kt_header_search_toggle">
          <button
            type="button"
            class="btn btn-icon btn-color-gray-700 btn-active-color-primary btn-outline w-40px h-40px"
            aria-label="搜索文章"
            title="搜索文章"
          >
            <span class="svg-icon svg-icon-1">
              <inline-svg
                :src="getAssetPath('/media/icons/duotune/general/gen021.svg')"
              />
            </span>
          </button>
        </div>
        <!--end::Search toggle-->
      </div>
      <!--end::Search-->
    </template>
    <template v-slot:content>
      <!--begin::Menu-->
      <div
        class="menu menu-sub menu-sub-dropdown menu-column p-7 w-325px w-md-375px"
        data-kt-menu="true"
        id="kt-search-menu"
      >
        <!--begin::Wrapper-->
        <div>
          <!--begin::Form-->
          <form class="w-100 position-relative mb-3" autocomplete="off">
            <!--begin::Icon-->
            <span
              class="svg-icon svg-icon-2 svg-icon-lg-1 svg-icon-gray-500 position-absolute top-50 translate-middle-y ms-0"
            >
              <inline-svg
                :src="getAssetPath('/media/icons/duotune/general/gen021.svg')"
              />
            </span>
            <!--end::Icon-->

            <!--begin::Input-->
            <input
              ref="inputRef"
              v-model="search"
              @input="searching"
              @keyup.enter="handleSearch"
              type="text"
              class="form-control form-control-flush ps-10"
              name="search"
              placeholder="搜索文章标题..."
            />
            <!--end::Input-->

            <!--begin::Spinner-->
            <span
              v-if="loading"
              class="position-absolute top-50 end-0 translate-middle-y lh-0 me-1"
            >
              <span
                class="spinner-border h-15px w-15px align-middle text-gray-400"
              ></span>
            </span>
            <!--end::Spinner-->

            <!--begin::Reset-->
            <span
              v-show="search.length && !loading"
              @click="reset()"
              class="btn btn-flush btn-active-color-primary position-absolute top-50 end-0 translate-middle-y lh-0"
            >
              <span class="svg-icon svg-icon-2 svg-icon-lg-1 me-0">
                <inline-svg
                  :src="getAssetPath('/media/icons/duotune/arrows/arr061.svg')"
                />
              </span>
            </span>
            <!--end::Reset-->

          </form>
          <!--end::Form-->

          <!--begin::Separator-->
          <div class="separator border-gray-200 mb-6" v-if="state !== 'main'"></div>
          <!--end::Separator-->
          
          <!--begin::Search Results-->
          <div v-if="state === 'results'" class="mb-4">
            <div class="d-flex flex-column align-items-center py-10">
              <div class="text-center mb-5">
                <i class="bi bi-search fs-2x text-primary mb-3"></i>
                <p class="text-gray-600 fs-6 mb-4">
                  搜索 "<strong>{{ search }}</strong>"
                </p>
                <button 
                  @click="handleSearch"
                  class="btn btn-primary btn-sm"
                  type="button"
                >
                  <i class="bi bi-arrow-right me-2"></i>
                  查看搜索结果
                </button>
              </div>
            </div>
          </div>
          <!--end::Search Results-->
          
          <PartialMain v-else-if="state === 'main'"></PartialMain>
          <Empty v-else-if="state === 'empty'"></Empty>
        </div>
        <!--end::Wrapper-->

      </div>
      <!--end::Menu-->
    </template>
  </MenuComponent>
</template>

<script lang="ts">
import { getAssetPath } from "@/core/helpers/assets";
import { defineComponent, ref } from "vue";
import { useRouter } from "vue-router";
import PartialMain from "@/layouts/main-layout/search/partials/Main.vue";
import Empty from "@/layouts/main-layout/search/partials/Empty.vue";
import MenuComponent from "@/components/menu/MenuComponent.vue";

export default defineComponent({
  name: "kt-search",
  components: {
    PartialMain,
    Empty,
    MenuComponent,
  },
  setup() {
    const router = useRouter();
    const search = ref<string>("");
    const state = ref<
      "main" | "empty" | "advanced-options" | "preferences" | "results"
    >("main");
    const loading = ref<boolean>(false);
    const inputRef = ref<HTMLInputElement | null>(null);

    const searching = (e: Event) => {
      const target = e.target as HTMLInputElement;
      const keyword = target.value.trim();
      
      if (keyword.length <= 1) {
        state.value = "main";
        loading.value = false;
      } else {
        // 实时搜索，显示结果
        loadResults(keyword);
      }
    };

    const loadResults = async (keyword: string) => {
      loading.value = true;
      try {
        // 这里可以调用 API 获取搜索结果预览
        // 暂时先显示结果状态，让用户点击搜索按钮跳转
        state.value = "results";
      } catch (error) {
        console.error("搜索失败:", error);
        state.value = "empty";
      } finally {
        loading.value = false;
      }
    };

    const handleSearch = () => {
      const keyword = search.value.trim();
      if (keyword.length > 0) {
        // 保存搜索历史到 localStorage
        saveSearchHistory(keyword);
        
        // 触发自定义事件，通知 Main.vue 刷新
        window.dispatchEvent(new Event('searchHistoryUpdated'));
        
        // 跳转到搜索结果页面
        router.push({
          name: "blog-home",
          query: { keyword: keyword },
        });
        // 关闭搜索菜单（通过触发菜单关闭事件）
        const menuElement = document.getElementById("kt-search-menu");
        if (menuElement) {
          const menuInstance = (window as any).KTMenu?.getInstance(menuElement);
          if (menuInstance) {
            menuInstance.hide();
          }
        }
        // 重置搜索框
        search.value = "";
        state.value = "main";
      }
    };

    const saveSearchHistory = (keyword: string) => {
      try {
        const history = JSON.parse(localStorage.getItem("searchHistory") || "[]");
        // 移除重复项
        const filtered = history.filter((item: string) => item !== keyword);
        // 添加到开头
        filtered.unshift(keyword);
        // 最多保存 10 条
        const limited = filtered.slice(0, 10);
        localStorage.setItem("searchHistory", JSON.stringify(limited));
      } catch (error) {
        console.error("保存搜索历史失败:", error);
      }
    };

    const reset = () => {
      search.value = "";
      state.value = "main";
    };

    const setState = (
      curr: "main" | "empty" | "advanced-options" | "preferences" | "results"
    ) => {
      state.value = curr;
    };

    return {
      search,
      state,
      loading,
      searching,
      reset,
      inputRef,
      setState,
      handleSearch,
      saveSearchHistory,
      getAssetPath,
    };
  },
});
</script>
