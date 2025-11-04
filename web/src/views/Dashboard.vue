<template>
  <!--begin::Row-->
  <div class="row g-5 g-xl-10 mb-5 mb-xl-10">
    <!--begin::Col-->
    <div class="col-xl-3">
      <!--begin::Card-->
      <div class="card card-xl-stretch mb-xl-10">
        <div class="card-body d-flex flex-column">
          <div class="d-flex flex-stack mb-5">
            <span class="text-gray-500 fw-semibold fs-6">总文章数</span>
            <span class="text-gray-500 fw-bold fs-6">
              <i class="bi bi-file-text text-primary fs-5"></i>
            </span>
          </div>
          <div class="d-flex align-items-center">
            <span class="text-gray-800 fw-bold fs-2qx me-2 lh-1">{{ stats.articleCount }}</span>
            <span class="text-gray-500 fw-semibold fs-7 align-self-end pb-1">篇</span>
          </div>
          <router-link to="/admin/article/index" class="text-primary fw-semibold fs-7 mt-5">
            查看文章管理 <i class="bi bi-arrow-right ms-1"></i>
          </router-link>
        </div>
      </div>
      <!--end::Card-->
    </div>
    <!--end::Col-->

    <!--begin::Col-->
    <div class="col-xl-3">
      <!--begin::Card-->
      <div class="card card-xl-stretch mb-xl-10">
        <div class="card-body d-flex flex-column">
          <div class="d-flex flex-stack mb-5">
            <span class="text-gray-500 fw-semibold fs-6">分类数量</span>
            <span class="text-gray-500 fw-bold fs-6">
              <i class="bi bi-folder text-success fs-5"></i>
            </span>
          </div>
          <div class="d-flex align-items-center">
            <span class="text-gray-800 fw-bold fs-2qx me-2 lh-1">{{ stats.categoryCount }}</span>
            <span class="text-gray-500 fw-semibold fs-7 align-self-end pb-1">个</span>
          </div>
          <router-link to="/admin/category" class="text-success fw-semibold fs-7 mt-5">
            查看分类管理 <i class="bi bi-arrow-right ms-1"></i>
          </router-link>
        </div>
      </div>
      <!--end::Card-->
    </div>
    <!--end::Col-->

    <!--begin::Col-->
    <div class="col-xl-3">
      <!--begin::Card-->
      <div class="card card-xl-stretch mb-xl-10">
        <div class="card-body d-flex flex-column">
          <div class="d-flex flex-stack mb-5">
            <span class="text-gray-500 fw-semibold fs-6">评论数量</span>
            <span class="text-gray-500 fw-bold fs-6">
              <i class="bi bi-chat-dots text-info fs-5"></i>
            </span>
          </div>
          <div class="d-flex align-items-center">
            <span class="text-gray-800 fw-bold fs-2qx me-2 lh-1">{{ stats.commentCount }}</span>
            <span class="text-gray-500 fw-semibold fs-7 align-self-end pb-1">条</span>
          </div>
          <router-link to="/admin/comment" class="text-info fw-semibold fs-7 mt-5">
            查看评论管理 <i class="bi bi-arrow-right ms-1"></i>
          </router-link>
        </div>
      </div>
      <!--end::Card-->
    </div>
    <!--end::Col-->

    <!--begin::Col-->
    <div class="col-xl-3">
      <!--begin::Card-->
      <div class="card card-xl-stretch mb-xl-10">
        <div class="card-body d-flex flex-column">
          <div class="d-flex flex-stack mb-5">
            <span class="text-gray-500 fw-semibold fs-6">友情链接</span>
            <span class="text-gray-500 fw-bold fs-6">
              <i class="bi bi-link-45deg text-warning fs-5"></i>
            </span>
          </div>
          <div class="d-flex align-items-center">
            <span class="text-gray-800 fw-bold fs-2qx me-2 lh-1">{{ stats.friendlinkCount }}</span>
            <span class="text-gray-500 fw-semibold fs-7 align-self-end pb-1">个</span>
          </div>
          <router-link to="/admin/friendlink" class="text-warning fw-semibold fs-7 mt-5">
            查看友链管理 <i class="bi bi-arrow-right ms-1"></i>
          </router-link>
        </div>
      </div>
      <!--end::Card-->
    </div>
    <!--end::Col-->
  </div>
  <!--end::Row-->

  <!--begin::Row-->
  <div class="row g-5 g-xl-10">
    <!--begin::Col-->
    <div class="col-xl-12">
      <!--begin::Card-->
      <div class="card card-xl-stretch">
        <div class="card-header border-0 pt-5">
          <h3 class="card-title align-items-start flex-column">
            <span class="card-label fw-bold fs-3 mb-1">快捷操作</span>
            <span class="text-muted mt-1 fw-semibold fs-7">快速访问常用功能</span>
          </h3>
        </div>
        <div class="card-body pt-0">
          <div class="d-flex flex-wrap gap-5">
            <router-link to="/admin/article/index" class="btn btn-flex btn-light-primary flex-column flex-center text-center w-200px h-150px">
              <span class="svg-icon svg-icon-3x mb-5">
                <i class="bi bi-file-text fs-2x text-primary"></i>
              </span>
              <span class="fs-5 fw-bold mb-2">文章管理</span>
              <span class="fs-7 text-gray-600">查看和编辑文章</span>
            </router-link>

            <router-link to="/admin/category" class="btn btn-flex btn-light-success flex-column flex-center text-center w-200px h-150px">
              <span class="svg-icon svg-icon-3x mb-5">
                <i class="bi bi-folder fs-2x text-success"></i>
              </span>
              <span class="fs-5 fw-bold mb-2">分类管理</span>
              <span class="fs-7 text-gray-600">管理文章分类</span>
            </router-link>

            <router-link to="/admin/article/content" class="btn btn-flex btn-light-warning flex-column flex-center text-center w-200px h-150px">
              <span class="svg-icon svg-icon-3x mb-5">
                <i class="bi bi-plus-circle fs-2x text-warning"></i>
              </span>
              <span class="fs-5 fw-bold mb-2">创建文章</span>
              <span class="fs-7 text-gray-600">撰写新的博客文章</span>
            </router-link>

            <router-link to="/admin/comment" class="btn btn-flex btn-light-info flex-column flex-center text-center w-200px h-150px">
              <span class="svg-icon svg-icon-3x mb-5">
                <i class="bi bi-chat-dots fs-2x text-info"></i>
              </span>
              <span class="fs-5 fw-bold mb-2">评论管理</span>
              <span class="fs-7 text-gray-600">管理用户评论</span>
            </router-link>

            <router-link to="/admin/friendlink" class="btn btn-flex btn-light-secondary flex-column flex-center text-center w-200px h-150px">
              <span class="svg-icon svg-icon-3x mb-5">
                <i class="bi bi-link-45deg fs-2x text-secondary"></i>
              </span>
              <span class="fs-5 fw-bold mb-2">友链管理</span>
              <span class="fs-7 text-gray-600">管理友情链接</span>
            </router-link>

            <router-link to="/index" class="btn btn-flex btn-light-dark flex-column flex-center text-center w-200px h-150px">
              <span class="svg-icon svg-icon-3x mb-5">
                <i class="bi bi-house fs-2x text-dark"></i>
              </span>
              <span class="fs-5 fw-bold mb-2">博客首页</span>
              <span class="fs-7 text-gray-600">查看博客前台</span>
            </router-link>
          </div>
        </div>
      </div>
      <!--end::Card-->
    </div>
    <!--end::Col-->
  </div>
  <!--end::Row-->
</template>

<script lang="ts">
import { defineComponent, onMounted, ref } from "vue";
import service from "@/utils/request";

export default defineComponent({
  name: "main-dashboard",
  setup() {
    const stats = ref({
      articleCount: 0,
      categoryCount: 0,
      commentCount: 0,
      friendlinkCount: 0
    });

    const loading = ref(false);

    const loadStats = async () => {
      loading.value = true;
      try {
        // 获取分类数量
        const categoryRes: any = await service.get("/base/getCategoryList");
        if (categoryRes && categoryRes.data && Array.isArray(categoryRes.data)) {
          stats.value.categoryCount = categoryRes.data.length;
        }

        // 获取文章数量
        const articleRes: any = await service.get("/base/getArticleList", {
          params: {
            page: 1,
            pageSize: 1
          }
        });
        if (articleRes && articleRes.data && articleRes.data.total !== undefined) {
          stats.value.articleCount = articleRes.data.total || 0;
        }

        // 获取评论数量
        const commentRes: any = await service.get("/admin/comment/getCommentList", {
          params: {
            page: 1,
            pageSize: 1
          }
        });
        if (commentRes && commentRes.data && commentRes.data.total !== undefined) {
          stats.value.commentCount = commentRes.data.total || 0;
        }

        // 获取友链数量
        const friendlinkRes: any = await service.get("/admin/friendlink/getFriendlinkList", {
          params: {
            page: 1,
            pageSize: 1
          }
        });
        if (friendlinkRes && friendlinkRes.data && friendlinkRes.data.total !== undefined) {
          stats.value.friendlinkCount = friendlinkRes.data.total || 0;
        }
      } catch (error) {
        console.error('获取统计数据失败:', error);
      } finally {
        loading.value = false;
      }
    };

    onMounted(() => {
      loadStats();
    });

    return {
      stats,
      loading
    };
  },
});
</script>

<style scoped>
.w-200px {
  width: 200px !important;
}

.h-150px {
  height: 150px !important;
}

.fs-2qx {
  font-size: 2.25rem !important;
}
</style>
