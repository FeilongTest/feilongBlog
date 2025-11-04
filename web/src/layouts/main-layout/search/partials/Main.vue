<template>
  <div class="mb-4">
    <!--begin::Heading-->
    <div class="d-flex flex-stack fw-semobold mb-4" v-if="searchHistory.length > 0">
      <!--begin::Label-->
      <span class="text-muted fs-6 me-2">最近搜索:</span>
      <!--end::Label-->
      <button 
        @click="clearHistory"
        class="btn btn-sm btn-link text-muted p-0"
        type="button"
      >
        清除
      </button>
    </div>
    <!--end::Heading-->
    <!--begin::Items-->
    <div class="scroll-y mh-200px mh-lg-325px" v-if="searchHistory.length > 0">
      <!--begin::Item-->
      <div 
        v-for="(item, index) in searchHistory" 
        :key="index"
        class="d-flex align-items-center mb-5 cursor-pointer"
        @click="handleHistoryClick(item)"
      >
        <!--begin::Symbol-->
        <div class="symbol symbol-40px symbol-circle me-4">
          <span class="symbol-label bg-light-primary">
            <i class="bi bi-clock-history fs-5 text-primary"></i>
          </span>
        </div>
        <!--end::Symbol-->
        <!--begin::Title-->
        <div class="d-flex flex-column flex-grow-1">
          <span class="fs-6 text-gray-800 fw-semobold">{{ item }}</span>
          <span class="fs-7 text-muted">点击搜索</span>
        </div>
        <!--end::Title-->
        <!--begin::Delete-->
        <button 
          @click.stop="removeHistoryItem(index)"
          class="btn btn-sm btn-icon btn-light-danger"
          type="button"
        >
          <i class="bi bi-x-lg fs-6"></i>
        </button>
        <!--end::Delete-->
      </div>
      <!--end::Item-->
    </div>
    <!--end::Items-->
    
    <!--begin::Empty-->
    <div v-else class="text-center py-10">
      <i class="bi bi-search fs-3x text-muted mb-3"></i>
      <p class="text-muted fs-6 mb-0">暂无搜索历史</p>
      <p class="text-muted fs-7 mt-2">开始搜索文章标题吧</p>
    </div>
    <!--end::Empty-->
  </div>
</template>

<script lang="ts">
import { getAssetPath } from "@/core/helpers/assets";
import { defineComponent, ref, onMounted, onUnmounted } from "vue";
import { useRouter } from "vue-router";

export default defineComponent({
  name: "kt-main",
  components: {},
  setup() {
    const router = useRouter();
    const searchHistory = ref<string[]>([]);

    // 加载搜索历史
    const loadSearchHistory = () => {
      try {
        const history = JSON.parse(localStorage.getItem("searchHistory") || "[]");
        searchHistory.value = history;
      } catch (error) {
        console.error("加载搜索历史失败:", error);
        searchHistory.value = [];
      }
    };

    // 点击历史记录
    const handleHistoryClick = (keyword: string) => {
      router.push({
        name: "blog-home",
        query: { keyword: keyword },
      });
      // 关闭搜索菜单
      const menuElement = document.getElementById("kt-search-menu");
      if (menuElement) {
        const menuInstance = (window as any).KTMenu?.getInstance(menuElement);
        if (menuInstance) {
          menuInstance.hide();
        }
      }
    };

    // 删除单条历史记录
    const removeHistoryItem = (index: number) => {
      searchHistory.value.splice(index, 1);
      localStorage.setItem("searchHistory", JSON.stringify(searchHistory.value));
    };

    // 清除所有历史记录
    const clearHistory = () => {
      searchHistory.value = [];
      localStorage.removeItem("searchHistory");
    };

    onMounted(() => {
      loadSearchHistory();
      // 监听存储事件，当搜索历史更新时自动刷新
      window.addEventListener('storage', loadSearchHistory);
      // 监听自定义事件（同页面内更新）
      window.addEventListener('searchHistoryUpdated', loadSearchHistory);
    });

    // 组件卸载时移除事件监听
    onUnmounted(() => {
      window.removeEventListener('storage', loadSearchHistory);
      window.removeEventListener('searchHistoryUpdated', loadSearchHistory);
    });

    return {
      searchHistory,
      handleHistoryClick,
      removeHistoryItem,
      clearHistory,
      getAssetPath,
    };
  },
});
</script>

<style scoped>
.cursor-pointer {
  cursor: pointer;
  transition: background-color 0.2s;
}

.cursor-pointer:hover {
  background-color: rgba(0, 0, 0, 0.02);
  border-radius: 0.5rem;
}
</style>
