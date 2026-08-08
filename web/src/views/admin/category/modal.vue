<template>
  <div id="kt_modal_category" ref="categoryModalRef" class="modal fade" tabindex="-1" aria-hidden="true">
    <div class="modal-dialog modal-dialog-centered mw-600px">
      <div class="modal-content border-0 overflow-hidden">
        <div class="modal-header border-0 px-8 pt-8 pb-3">
          <div class="d-flex align-items-center">
            <span class="modal-heading-icon me-4"><i :class="mode === 'add' ? 'bi bi-folder-plus' : 'bi bi-pencil-square'"></i></span>
            <div>
              <h2 class="fw-bolder text-gray-900 mb-1">{{ mode === 'add' ? '新建分类' : '编辑分类' }}</h2>
              <span class="text-muted fs-7">设置名称、所属层级与展示权重</span>
            </div>
          </div>
          <button type="button" class="btn btn-sm btn-icon btn-light" data-bs-dismiss="modal" aria-label="关闭"><i class="bi bi-x-lg"></i></button>
        </div>

        <el-form ref="formRef" :model="formData" :rules="rules" label-position="top" @submit.prevent="submit">
          <div class="modal-body px-8 py-6">
            <el-form-item label="分类名称" prop="name">
              <el-input v-model.trim="formData.name" size="large" maxlength="30" show-word-limit placeholder="例如：技术随笔" />
            </el-form-item>

            <el-form-item label="上级分类" prop="fid">
              <el-select v-model="formData.fid" size="large" class="w-100" placeholder="请选择所属层级">
                <el-option :value="0" label="设为一级分类" />
                <el-option v-for="item in parentOptions" :key="item.ID" :value="item.ID" :label="item.name" />
              </el-select>
              <div class="form-hint"><i class="bi bi-info-circle me-2"></i>选择一级分类后，当前分类会显示在它的下方。</div>
            </el-form-item>

            <el-form-item label="展示权重" prop="sort" class="mb-2">
              <el-input-number v-model="formData.sort" :min="0" :max="9999" controls-position="right" class="w-100" size="large" />
              <div class="form-hint"><i class="bi bi-sort-down me-2"></i>数字越大，分类在同级列表中的位置越靠前。</div>
            </el-form-item>
          </div>

          <div class="modal-footer border-0 bg-light px-8 py-5">
            <button type="button" class="btn btn-light" data-bs-dismiss="modal">取消</button>
            <button type="submit" class="btn btn-primary px-6" :disabled="loading">
              <span v-if="loading" class="spinner-border spinner-border-sm me-2"></span>
              {{ loading ? '正在保存' : '保存分类' }}
            </button>
          </div>
        </el-form>
      </div>
    </div>
  </div>
</template>

<script lang="ts">
import { computed, defineComponent, nextTick, ref, watch } from "vue";
import type { FormInstance, FormRules } from "element-plus";
import { ElNotification } from "element-plus";
import { hideModal } from "@/core/helpers/dom";
import type { CategoryList } from "@/core/blog/CategoryTypes";
import service from "@/utils/request";

export default defineComponent({
  name: "category-modal",
  props: {
    id: { type: Number, default: 0 },
    parentId: { type: Number, default: 0 },
    mode: { type: String, default: "add" },
    data: { type: Array<CategoryList>, default: () => [] },
  },
  emits: ["saved"],
  setup(props, { emit }) {
    const formRef = ref<FormInstance>();
    const categoryModalRef = ref<HTMLElement | null>(null);
    const loading = ref(false);
    const formData = ref({ ID: 0, name: "", fid: 0, sort: 0, type: 0 });

    const parentOptions = computed(() => props.data.filter(item => item.fid === 0 && item.ID !== props.id));
    const rules: FormRules = {
      name: [{ required: true, message: "请输入分类名称", trigger: "blur" }],
      fid: [{ required: true, message: "请选择上级分类", trigger: "change" }],
    };

    const syncForm = () => {
      const current = props.data.find(item => item.ID === props.id);
      formData.value = current && props.mode === "modify"
        ? { ID: current.ID, name: current.name, fid: current.fid, sort: current.sort, type: current.type }
        : { ID: 0, name: "", fid: props.parentId, sort: 0, type: 0 };
      nextTick(() => formRef.value?.clearValidate());
    };
    watch(() => [props.id, props.mode, props.parentId, props.data], syncForm, { deep: true, immediate: true });

    const submit = async () => {
      if (!formRef.value || !(await formRef.value.validate().catch(() => false))) return;
      loading.value = true;
      try {
        if (props.mode === "add") {
          await service.post("/admin/category/createCategory", formData.value);
        } else {
          await service.put("/admin/category/updateCategory", formData.value);
        }
        ElNotification({ title: "保存成功", message: `分类“${formData.value.name}”已更新`, type: "success" });
        hideModal(categoryModalRef.value);
        emit("saved");
      } finally {
        loading.value = false;
      }
    };

    return { formRef, categoryModalRef, loading, formData, parentOptions, rules, submit };
  },
});
</script>

<style lang="scss" scoped>
.modal-content { border-radius: 18px; box-shadow: 0 24px 70px rgba(26, 35, 62, .18); }
.modal-heading-icon { width: 48px; height: 48px; display: inline-flex; align-items: center; justify-content: center; border-radius: 13px; color: var(--kt-primary); background: var(--kt-primary-light); font-size: 1.25rem; }
.form-hint { width: 100%; margin-top: 8px; color: var(--kt-gray-500); font-size: .82rem; }
:deep(.el-form-item__label) { color: var(--kt-gray-800); font-weight: 700; font-size: .95rem; }
:deep(.el-input__wrapper), :deep(.el-select__wrapper) { border-radius: 10px; box-shadow: 0 0 0 1px var(--kt-gray-300) inset; }
</style>
