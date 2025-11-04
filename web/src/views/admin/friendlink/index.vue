<template>
  <div class="card">
    <!--begin::Header-->
    <div class="card-header border-0 pt-5">
      <h3 class="card-title align-items-start flex-column">
        <span class="card-label fw-bold fs-3 mb-1">友链管理</span>
        <span class="text-muted mt-1 fw-semobold fs-7">管理友情链接</span>
      </h3>
      <div class="card-toolbar">
        <button
          type="button"
          class="btn btn-primary"
          @click="handleAdd"
        >
          <span class="svg-icon svg-icon-2">
            <inline-svg
              :src="getAssetPath('/media/icons/duotune/arrows/arr075.svg')"
            />
          </span>
          新增友链
        </button>
      </div>
    </div>
    <!--end::Header-->

    <!--begin::Body-->
    <div class="card-body py-3">
      <!--begin::Table-->
      <el-table :data="friendlinkList" style="width: 100%" v-loading="loading">
        <el-table-column prop="ID" label="ID" width="80" />
        <el-table-column label="Logo" width="100">
          <template #default="{ row }">
            <el-avatar v-if="row.logo" :src="row.logo" :size="40" />
            <el-avatar v-else :size="40">
              <i class="bi bi-link-45deg fs-2"></i>
            </el-avatar>
          </template>
        </el-table-column>
        <el-table-column prop="name" label="友链名称" width="150" />
        <el-table-column prop="url" label="友链地址" min-width="250" show-overflow-tooltip />
        <el-table-column prop="desc" label="描述" min-width="200" show-overflow-tooltip />
        <el-table-column label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="row.status === 0 ? 'success' : 'danger'">
              {{ row.status === 0 ? '显示' : '隐藏' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="150" fixed="right">
          <template #default="{ row }">
            <button
              type="button"
              class="btn btn-icon btn-bg-light btn-active-color-primary btn-sm me-1"
              @click="handleEdit(row)"
              title="编辑"
            >
              <span class="svg-icon svg-icon-3">
                <inline-svg :src="getAssetPath('/media/icons/duotune/art/art005.svg')" />
              </span>
            </button>
            <button
              type="button"
              class="btn btn-icon btn-bg-light btn-active-color-primary btn-sm"
              @click="handleDelete(row)"
              title="删除"
            >
              <span class="svg-icon svg-icon-3">
                <inline-svg :src="getAssetPath('/media/icons/duotune/general/gen027.svg')" />
              </span>
            </button>
          </template>
        </el-table-column>
      </el-table>
      <!--end::Table-->

      <!--begin::Pagination-->
      <div class="d-flex flex-stack flex-wrap pt-10">
        <div class="fs-6 fw-semobold text-gray-700">
          共 {{ total }} 条友链
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

    <!--begin::Dialog-->
    <el-dialog
      v-model="dialogVisible"
      :title="isEdit ? '编辑友链' : '新增友链'"
      width="600px"
      @close="resetForm"
    >
      <el-form :model="formData" :rules="rules" ref="formRef" label-width="100px">
        <el-form-item label="友链名称" prop="name">
          <el-input v-model="formData.name" placeholder="请输入友链名称" />
        </el-form-item>
        <el-form-item label="友链地址" prop="url">
          <el-input v-model="formData.url" placeholder="请输入友链地址" />
        </el-form-item>
        <el-form-item label="Logo地址" prop="logo">
          <el-input v-model="formData.logo" placeholder="请输入Logo图片地址" />
        </el-form-item>
        <el-form-item label="描述">
          <el-input
            v-model="formData.desc"
            type="textarea"
            :rows="3"
            placeholder="请输入友链描述"
          />
        </el-form-item>
        <el-form-item label="状态">
          <el-radio-group v-model="formData.status">
            <el-radio :label="0">显示</el-radio>
            <el-radio :label="1">隐藏</el-radio>
          </el-radio-group>
        </el-form-item>
      </el-form>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="dialogVisible = false">取消</el-button>
          <el-button type="primary" @click="handleSubmit">确定</el-button>
        </span>
      </template>
    </el-dialog>
    <!--end::Dialog-->
  </div>
</template>

<script lang="ts">
import { defineComponent, ref, reactive, onMounted } from "vue";
import service from "@/utils/request";
import type { Friendlink, FriendlinkFormData } from "@/core/blog/FriendlinkTypes";
import { ElMessage, ElMessageBox, type FormInstance, type FormRules } from "element-plus";
import { getAssetPath } from "@/core/helpers/assets";
import Swal from "sweetalert2";

export default defineComponent({
  name: "admin-friendlink",
  setup() {
    const friendlinkList = ref<Friendlink[]>([]);
    const loading = ref(false);
    const page = ref(1);
    const pageSize = ref(10);
    const total = ref(0);
    const dialogVisible = ref(false);
    const isEdit = ref(false);
    const formRef = ref<FormInstance>();

    const formData = reactive<FriendlinkFormData>({
      name: "",
      url: "",
      logo: "",
      desc: "",
      status: 0,
    });

    const rules: FormRules = {
      name: [{ required: true, message: "请输入友链名称", trigger: "blur" }],
      url: [{ required: true, message: "请输入友链地址", trigger: "blur" }],
    };

    // 获取友链列表
    const getFriendlinkList = async () => {
      loading.value = true;
      try {
        const res: any = await service.get("/admin/friendlink/getFriendlinkList", {
          params: {
            page: page.value,
            pageSize: pageSize.value,
          },
        });
        if (res && res.data) {
          friendlinkList.value = res.data.list || [];
          total.value = res.data.total || 0;
        }
      } catch (error) {
        console.error("获取友链列表失败:", error);
        ElMessage.error("获取友链列表失败");
      } finally {
        loading.value = false;
      }
    };

    // 新增友链
    const handleAdd = () => {
      isEdit.value = false;
      resetForm();
      dialogVisible.value = true;
    };

    // 编辑友链
    const handleEdit = (row: Friendlink) => {
      isEdit.value = true;
      formData.ID = row.ID;
      formData.name = row.name;
      formData.url = row.url;
      formData.logo = row.logo;
      formData.desc = row.desc;
      formData.status = row.status;
      dialogVisible.value = true;
    };

    // 提交表单
    const handleSubmit = async () => {
      if (!formRef.value) return;
      await formRef.value.validate(async (valid) => {
        if (valid) {
          try {
            if (isEdit.value) {
              await service.put("/admin/friendlink/updateFriendlink", formData);
              ElMessage.success("更新成功");
            } else {
              await service.post("/admin/friendlink/createFriendlink", formData);
              ElMessage.success("创建成功");
            }
            dialogVisible.value = false;
            getFriendlinkList();
          } catch (error) {
            console.error("操作失败:", error);
            ElMessage.error("操作失败");
          }
        }
      });
    };

    // 删除友链
    const handleDelete = (row: Friendlink) => {
      Swal.fire({
        title: `确定要删除友链"${row.name}"吗？`,
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
            await service.delete("/admin/friendlink/deleteFriendlink", {
              data: { ID: row.ID },
            });
            ElMessage.success("删除成功");
            getFriendlinkList();
          } catch (error) {
            console.error("删除失败:", error);
            ElMessage.error("删除失败");
          }
        }
      });
    };

    // 重置表单
    const resetForm = () => {
      formData.ID = undefined;
      formData.name = "";
      formData.url = "";
      formData.logo = "";
      formData.desc = "";
      formData.status = 0;
      formRef.value?.clearValidate();
    };

    // 分页变化
    const pageChange = (newPage: number) => {
      page.value = newPage;
      getFriendlinkList();
    };

    const itemPageChange = (newSize: number) => {
      pageSize.value = newSize;
      page.value = 1;
      getFriendlinkList();
    };

    onMounted(() => {
      getFriendlinkList();
    });

    return {
      friendlinkList,
      loading,
      page,
      pageSize,
      total,
      dialogVisible,
      isEdit,
      formRef,
      formData,
      rules,
      getFriendlinkList,
      handleAdd,
      handleEdit,
      handleSubmit,
      handleDelete,
      resetForm,
      pageChange,
      itemPageChange,
      getAssetPath,
      Math,
    };
  },
});
</script>

<style scoped>
:deep(.el-table) {
  font-size: 13px;
}
</style>

