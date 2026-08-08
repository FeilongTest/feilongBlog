<template>
  <div class="card">
    <!--begin::Header-->
    <div class="card-header border-0 pt-5">
      <h3 class="card-title align-items-start flex-column">
        <span class="card-label fw-bold fs-3 mb-1">评论管理</span>
        <span class="text-muted mt-1 fw-semobold fs-7">管理用户评论</span>
      </h3>
    </div>
    <!--end::Header-->

    <!--begin::Body-->
    <div class="card-body py-3">
      <!--begin::Table container-->
      <div class="table-responsive">
        <!--begin::Table-->
        <table class="table">
          <!--begin::Table head-->
          <thead>
            <tr class="fw-bold text-muted bg-light">
              <th class="ps-4 min-w-80px rounded-start">ID</th>
              <th class="min-w-150px">昵称</th>
              <th class="min-w-200px">邮箱</th>
              <th class="min-w-300px">评论内容</th>
              <th class="min-w-100px">文章ID</th>
              <th class="min-w-100px">状态</th>
              <th class="min-w-180px">评论时间</th>
              <th class="min-w-150px text-end rounded-end pe-4">操作</th>
            </tr>
          </thead>
          <!--end::Table head-->
          <!--begin::Table body-->
          <tbody>
            <tr v-for="item in commentList" :key="item.ID">
              <td class="ps-4">
                <div class="text-dark fw-bold text-hover-primary mb-1 fs-6">{{ item.ID }}</div>
              </td>
              <td>
                <div class="text-dark fw-bold text-hover-primary mb-1 fs-6">{{ item.name }}</div>
              </td>
              <td>
                <div class="text-muted fw-semibold text-muted d-block fs-7">{{ item.email }}</div>
              </td>
              <td>
                <div class="text-dark fw-bold text-hover-primary mb-1 fs-6" :title="item.content">
                  {{ item.content.length > 50 ? item.content.substring(0, 50) + '...' : item.content }}
                </div>
              </td>
              <td>
                <div class="text-dark fw-bold text-hover-primary mb-1 fs-6">{{ item.aid }}</div>
              </td>
              <td>
                <span :class="item.status === 0 ? 'badge badge-light-success' : 'badge badge-light-danger'">
                  {{ item.status === 0 ? '显示' : '隐藏' }}
                </span>
              </td>
              <td>
                <div class="text-muted fw-semibold text-muted d-block fs-7">{{ formatDate(item.ctime) }}</div>
              </td>
              <td class="text-end">
                <button
                  type="button"
                  class="btn btn-icon btn-bg-light btn-active-color-primary btn-sm me-1"
                  @click="updateStatus(item, item.status === 0 ? 1 : 0)"
                  :title="item.status === 0 ? '隐藏' : '显示'"
                >
                  <span class="svg-icon svg-icon-3">
                    <inline-svg :src="getAssetPath('/media/icons/duotune/general/gen019.svg')" />
                  </span>
                </button>
                <button
                  type="button"
                  class="btn btn-icon btn-bg-light btn-active-color-primary btn-sm"
                  @click="deleteComment(item)"
                  title="删除"
                >
                  <span class="svg-icon svg-icon-3">
                    <inline-svg :src="getAssetPath('/media/icons/duotune/general/gen027.svg')" />
                  </span>
                </button>
              </td>
            </tr>
            <tr v-if="commentList.length === 0 && !loading">
              <td colspan="8" class="text-center py-10">
                <span class="text-muted">暂无评论数据</span>
              </td>
            </tr>
          </tbody>
          <!--end::Table body-->
        </table>
        <!--end::Table-->
      </div>
      <!--end::Table container-->

      <!--begin::Pagination-->
      <div class="d-flex flex-stack flex-wrap pt-10">
        <div class="fs-6 fw-semobold text-gray-700">
          共 {{ total }} 条评论
        </div>
        <div class="d-flex align-items-center">
          <select 
            class="form-select form-select-sm form-select-solid w-70px me-3"
            v-model="pageSize"
            @change="itemPageChange(pageSize)"
          >
            <option :value="10">10</option>
            <option :value="20">20</option>
            <option :value="50">50</option>
            <option :value="100">100</option>
          </select>
          <div class="d-flex align-items-center">
            <button
              type="button"
              class="btn btn-sm btn-light-primary me-2"
              @click="pageChange(page - 1)"
              :disabled="page <= 1"
            >
              上一页
            </button>
            <span class="text-gray-700 fw-semibold me-2">
              第 {{ page }} 页 / 共 {{ Math.ceil(total / pageSize) }} 页
            </span>
            <button
              type="button"
              class="btn btn-sm btn-light-primary"
              @click="pageChange(page + 1)"
              :disabled="page >= Math.ceil(total / pageSize)"
            >
              下一页
            </button>
          </div>
        </div>
      </div>
      <!--end::Pagination-->
    </div>
    <!--end::Body-->
  </div>
</template>

<script lang="ts">
import { defineComponent, ref, onMounted } from "vue";
import service from "@/utils/request";
import type { Comment } from "@/core/blog/CommentTypes";
import { ElMessage, ElMessageBox } from "element-plus";
import dayjs from "dayjs";
import { getAssetPath } from "@/core/helpers/assets";
import Swal from "sweetalert2";

export default defineComponent({
  name: "admin-comment",
  setup() {
    const commentList = ref<Comment[]>([]);
    const loading = ref(false);
    const page = ref(1);
    const pageSize = ref(10);
    const total = ref(0);

    // 获取评论列表
    const getCommentList = async () => {
      loading.value = true;
      try {
        const res: any = await service.get("/admin/comment/getCommentList", {
          params: {
            page: page.value,
            pageSize: pageSize.value,
          },
        });
        if (res && res.data) {
          commentList.value = res.data.list || [];
          total.value = res.data.total || 0;
        }
      } catch (error) {
        console.error("获取评论列表失败:", error);
        ElMessage.error("获取评论列表失败");
      } finally {
        loading.value = false;
      }
    };

    // 更新评论状态
    const updateStatus = async (row: Comment, status: number) => {
      try {
        await service.put("/admin/comment/updateCommentStatus", {
          ID: row.ID,
          status: status,
        });
        ElMessage.success("更新成功");
        getCommentList();
      } catch (error) {
        console.error("更新状态失败:", error);
        ElMessage.error("更新状态失败");
      }
    };

    // 删除评论
    const deleteComment = async (row: Comment) => {
      Swal.fire({
        title: "确定要删除这条评论吗？",
        icon: "warning",
        showCancelButton: true,
        confirmButtonText: "确定",
        cancelButtonText: "取消",
        buttonsStyling: false,
        customClass: {
          confirmButton: "btn fw-semobold btn-light-danger",
          cancelButton: "btn fw-semobold btn-secondary",
        },
      }).then(async (result) => {
        if (result.isConfirmed) {
          try {
            await service.delete("/admin/comment/deleteComment", {
              data: { ID: row.ID },
            });
            ElMessage.success("删除成功");
            getCommentList();
          } catch (error) {
            console.error("删除失败:", error);
            ElMessage.error("删除失败");
          }
        }
      });
    };

    // 格式化日期
    const formatDate = (timestamp: number): string => {
      return dayjs.unix(timestamp).format("YYYY-MM-DD HH:mm:ss");
    };

    // 分页变化
    const pageChange = (newPage: number) => {
      page.value = newPage;
      getCommentList();
    };

    const itemPageChange = (newSize: number) => {
      pageSize.value = newSize;
      page.value = 1;
      getCommentList();
    };

    onMounted(() => {
      getCommentList();
    });

    return {
      commentList,
      loading,
      page,
      pageSize,
      total,
      getCommentList,
      updateStatus,
      deleteComment,
      formatDate,
      pageChange,
      itemPageChange,
      getAssetPath,
      Math,
    };
  },
});
</script>

