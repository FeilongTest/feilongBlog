<template>
  <div class="card">
    <div class="card-header align-items-center gap-3">
      <div>
        <h1 class="fs-3 mb-1">标签管理</h1>
        <p class="text-muted fs-7 mb-0">分类用于内容层级，标签连接相同主题。</p>
      </div>
      <button type="button" class="btn btn-sm btn-primary" @click="edit()">
        新建标签
      </button>
    </div>
    <div class="card-body">
      <input
        v-model="query"
        aria-label="筛选标签"
        placeholder="搜索标签名称"
        class="form-control form-control-solid mb-6"
        style="max-width: 400px"
      />
      <p v-if="loading" role="status">正在加载标签…</p>
      <div v-else-if="error" class="alert alert-danger">
        加载失败。<button class="btn btn-sm btn-light" @click="load">
          重试
        </button>
      </div>
      <div v-else class="table-responsive">
        <table class="table align-middle">
          <thead>
            <tr>
              <th>标签</th>
              <th>描述</th>
              <th>已发布文章</th>
              <th class="text-end">操作</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="tag in filtered" :key="tag.id">
              <td class="fw-semibold">{{ tag.name }}</td>
              <td class="text-muted">{{ tag.description || "—" }}</td>
              <td>{{ tag.count || 0 }}</td>
              <td class="text-end text-nowrap">
                <button
                  class="btn btn-sm btn-light-primary me-2"
                  @click="edit(tag)"
                >
                  编辑</button
                ><button
                  class="btn btn-sm btn-light-danger"
                  :disabled="deleting === tag.id"
                  @click="remove(tag)"
                >
                  移除
                </button>
              </td>
            </tr>
          </tbody>
        </table>
        <p v-if="!filtered.length" class="text-center text-muted py-5">
          暂无标签，可在此创建，也可在文章编辑器中输入。
        </p>
      </div>
    </div>
    <el-dialog
      v-model="showEditor"
      :title="form.id ? '编辑标签' : '新建标签'"
      width="min(500px, 95vw)"
      ><form @submit.prevent="save">
        <label for="tag-name" class="form-label">标签名称</label
        ><input
          id="tag-name"
          v-model="form.name"
          required
          maxlength="32"
          class="form-control form-control-solid mb-4"
        /><label for="tag-description" class="form-label">标签描述</label
        ><textarea
          id="tag-description"
          v-model="form.description"
          maxlength="300"
          rows="4"
          class="form-control form-control-solid mb-5"
        ></textarea
        ><button type="submit" :disabled="saving" class="btn btn-primary">
          {{ saving ? "正在保存…" : "保存标签" }}
        </button>
      </form></el-dialog
    >
  </div>
</template>
<script setup lang="ts">
import { computed, onMounted, reactive, ref } from "vue";
import { ElNotification } from "element-plus";
import Swal from "sweetalert2";
import type { Tag } from "@/core/blog/PublishingTypes";
import service from "@/utils/request";
const tags = ref<Tag[]>([]);
const query = ref("");
const loading = ref(true);
const error = ref(false);
const saving = ref(false);
const deleting = ref(0);
const showEditor = ref(false);
const form = reactive<Tag>({ id: 0, name: "", description: "" });
const filtered = computed(() =>
  tags.value.filter((tag) =>
    tag.name
      .toLocaleLowerCase()
      .includes(query.value.trim().toLocaleLowerCase())
  )
);
async function load() {
  loading.value = true;
  error.value = false;
  try {
    const response = await service.get("/admin/article/tags");
    tags.value = Array.isArray(response.data) ? response.data : [];
  } catch {
    error.value = true;
  } finally {
    loading.value = false;
  }
}
function edit(tag?: Tag) {
  Object.assign(form, tag || { id: 0, name: "", description: "" });
  showEditor.value = true;
}
async function save() {
  if (saving.value) return;
  saving.value = true;
  try {
    await service.post("/admin/article/tags", {
      id: form.id,
      name: form.name.trim(),
      description: form.description.trim(),
    });
    showEditor.value = false;
    await load();
    ElNotification({ title: "标签已保存", type: "success" });
  } catch {
  } finally {
    saving.value = false;
  }
}
async function remove(tag: Tag) {
  const answer = await Swal.fire({
    title: `移除标签“${tag.name}”？`,
    text: "标签与文章的关联会移除，文章本身保留。",
    icon: "warning",
    showCancelButton: true,
    confirmButtonText: "移除标签",
    cancelButtonText: "取消",
  });
  if (!answer.isConfirmed) return;
  deleting.value = tag.id;
  try {
    await service.delete("/admin/article/tags", { data: { id: tag.id } });
    await load();
  } catch {
  } finally {
    deleting.value = 0;
  }
}
onMounted(load);
</script>
