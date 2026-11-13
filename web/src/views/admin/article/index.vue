<template>
  <div class="card">
    <div class="card-header align-items-center gap-4 py-5">
      <div>
        <span class="text-primary fs-8 fw-bold">CONTENT LIBRARY</span>
        <h1 class="fs-3 mt-2 mb-1">文章管理</h1>
        <p class="text-muted fs-7 mb-0">
          整理草稿、安排发布，并回到你正在写的内容。
        </p>
      </div>
      <router-link to="/admin/article/content" class="btn btn-primary btn-sm"
        ><i class="bi bi-plus-lg"></i>写文章</router-link
      >
    </div>
    <div class="card-body">
      <nav class="d-flex flex-wrap gap-2 mb-6" aria-label="文章状态筛选">
        <button
          v-for="state in states"
          :key="state.value"
          type="button"
          class="btn btn-sm"
          :class="filters.state === state.value ? 'btn-primary' : 'btn-light'"
          @click="changeState(state.value)"
        >
          {{ state.label }}
        </button>
      </nav>
      <form
        class="d-flex flex-wrap align-items-center gap-3 mb-6"
        @submit.prevent="applyFilters"
      >
        <input
          v-model="filters.keyword"
          type="search"
          maxlength="80"
          placeholder="搜索文章标题"
          aria-label="后台搜索标题"
          class="form-control form-control-solid"
          style="max-width: 260px"
        />
        <el-select
          v-model="filters.fid"
          filterable
          aria-label="按分类筛选"
          class="library-select"
          ><el-option :value="0" label="全部分类" /><el-option
            v-for="category in categories"
            :key="category.ID"
            :value="category.ID"
            :label="category.name"
        /></el-select>
        <el-select
          v-model="filters.tag"
          filterable
          aria-label="按标签筛选"
          class="library-select"
          ><el-option :value="0" label="全部标签" /><el-option
            v-for="tag in tags"
            :key="tag.id"
            :value="tag.id"
            :label="tag.name"
        /></el-select>
        <button type="submit" class="btn btn-sm btn-light-primary">筛选</button
        ><button
          v-if="activeFilters"
          type="button"
          class="btn btn-sm btn-light"
          @click="resetFilters"
        >
          清除筛选
        </button>
      </form>
      <div
        class="d-flex flex-wrap justify-content-between align-items-center gap-3 mb-4"
      >
        <span class="text-muted fs-7" role="status">{{
          loading ? "正在加载…" : `共 ${total} 篇文章`
        }}</span>
        <div class="d-flex align-items-center gap-3">
          <button
            v-if="selected.length"
            type="button"
            class="btn btn-sm btn-light-danger"
            :disabled="busy"
            @click="removeSelected"
          >
            删除所选 {{ selected.length }} 篇</button
          ><select
            v-model.number="pageSize"
            aria-label="每页文章数"
            class="form-select form-select-sm w-auto"
            @change="changePage(1)"
          >
            <option :value="10">每页 10 篇</option>
            <option :value="25">每页 25 篇</option>
            <option :value="50">每页 50 篇</option>
          </select>
        </div>
      </div>
      <div v-if="error" class="alert alert-danger">
        文章列表加载失败。<button
          class="btn btn-sm btn-light ms-3"
          @click="load"
        >
          重试
        </button>
      </div>
      <div class="table-responsive">
        <table class="table table-row-dashed align-middle library-table">
          <thead>
            <tr class="text-muted fs-7">
              <th>
                <input
                  type="checkbox"
                  aria-label="选择当前页全部文章"
                  :checked="rows.length > 0 && selected.length === rows.length"
                  :disabled="loading || !rows.length"
                  @change="selectPage"
                />
              </th>
              <th>文章</th>
              <th>状态</th>
              <th>时间</th>
              <th>阅读</th>
              <th>置顶</th>
              <th class="text-end">操作</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="article in rows" :key="article.ID">
              <td>
                <input
                  v-model="selected"
                  type="checkbox"
                  :value="article.ID"
                  :aria-label="`选择 ${article.title}`"
                  :disabled="busy"
                />
              </td>
              <td class="library-title">
                <router-link
                  :to="`/admin/article/content/${article.ID}`"
                  class="text-gray-900 text-hover-primary fw-bold"
                  >{{ article.title }}</router-link
                >
                <div class="text-muted fs-8 mt-1">
                  {{ categoryName(article.fid) }}
                </div>
                <div
                  v-if="article.tags?.length"
                  class="d-flex flex-wrap gap-1 mt-2"
                >
                  <span
                    v-for="tag in article.tags"
                    :key="tag.id"
                    class="badge badge-light"
                    >{{ tag.name }}</span
                  >
                </div>
              </td>
              <td>
                <span class="badge" :class="statusClass(article.status)">{{
                  statusName(article.status)
                }}</span>
                <div v-if="article.status === 3" class="text-muted fs-8 mt-2">
                  {{ formatTime(article.publishAt) }}
                </div>
              </td>
              <td class="text-muted fs-8">
                <div>
                  {{ article.status === 0 ? "发布" : "创建" }}：{{
                    formatTime(article.ctime)
                  }}
                </div>
                <div class="mt-1">更新：{{ formatTime(article.editTime) }}</div>
              </td>
              <td>{{ article.view || 0 }}</td>
              <td>
                <button
                  type="button"
                  class="btn btn-icon btn-sm"
                  :class="article.isTop ? 'btn-light-warning' : 'btn-light'"
                  :aria-label="
                    article.isTop
                      ? `取消置顶 ${article.title}`
                      : `置顶 ${article.title}`
                  "
                  :disabled="busy"
                  @click="setTop(article)"
                >
                  <i class="bi bi-pin-angle"></i>
                </button>
              </td>
              <td class="text-end">
                <div class="d-flex justify-content-end flex-wrap gap-2">
                  <router-link
                    :to="`/admin/article/content/${article.ID}`"
                    class="btn btn-sm btn-light-primary"
                    >编辑</router-link
                  ><router-link
                    v-if="article.status === 0"
                    :to="`/content/${article.fid}/${article.ID}`"
                    class="btn btn-sm btn-light"
                    >查看</router-link
                  ><button
                    v-if="article.status === 0"
                    type="button"
                    class="btn btn-sm btn-light"
                    :disabled="busy"
                    @click="setVisibility(article, 1)"
                  >
                    隐藏</button
                  ><button
                    v-if="article.status === 1"
                    type="button"
                    class="btn btn-sm btn-light"
                    :disabled="busy"
                    @click="setVisibility(article, 0)"
                  >
                    重新发布</button
                  ><button
                    type="button"
                    class="btn btn-sm btn-light-danger"
                    :disabled="busy"
                    @click="removeArticle(article)"
                  >
                    删除
                  </button>
                </div>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
      <div
        v-if="!loading && !error && !rows.length"
        class="text-center text-muted py-12"
      >
        <i class="bi bi-journal fs-2x d-block mb-4"></i
        >没有符合条件的文章。<router-link
          to="/admin/article/content"
          class="d-block mt-3"
          >开始写一篇新文章</router-link
        >
      </div>
      <TablePagination
        v-if="total > pageSize"
        :total-pages="Math.ceil(total / pageSize)"
        :total="total"
        :per-page="pageSize"
        :current-page="page"
        @page-change="changePage"
      />
    </div>
  </div>
</template>
<script setup lang="ts">
import {
  computed,
  onMounted,
  onBeforeUnmount,
  reactive,
  ref,
  watch,
} from "vue";
import { useRoute, useRouter } from "vue-router";
import { ElNotification } from "element-plus";
import Swal from "sweetalert2";
import dayjs from "dayjs";
import service from "@/utils/request";
import type { ArticleList } from "@/core/blog/ArticleTypes";
import type { Tag, ArticleStatus } from "@/core/blog/PublishingTypes";
import { articleStatusName } from "@/core/blog/PublishingTypes";
import type { CategoryList } from "@/core/blog/CategoryTypes";
import TablePagination from "@/components/kt-datatable/table-partials/table-content/table-footer/TablePagination.vue";
const route = useRoute();
const router = useRouter();
const rows = ref<ArticleList[]>([]);
const categories = ref<CategoryList[]>([]);
const tags = ref<Tag[]>([]);
const total = ref(0);
const page = ref(1);
const pageSize = ref(10);
const selected = ref<number[]>([]);
const loading = ref(true);
const error = ref(false);
const busy = ref(false);
let sequence = 0;
const filters = reactive({ state: "all", keyword: "", fid: 0, tag: 0 });
const states = [
  { value: "all", label: "全部" },
  { value: "draft", label: "草稿" },
  { value: "scheduled", label: "定时发布" },
  { value: "published", label: "已发布" },
  { value: "hidden", label: "已隐藏" },
  { value: "top", label: "置顶" },
];
const activeFilters = computed(
  () =>
    filters.state !== "all" ||
    !!filters.keyword ||
    !!filters.fid ||
    !!filters.tag
);
function syncRoute() {
  filters.state = states.some((s) => s.value === route.query.state)
    ? String(route.query.state)
    : "all";
  filters.keyword = String(route.query.keyword || "").slice(0, 80);
  filters.fid = Math.max(0, Number(route.query.fid) || 0);
  filters.tag = Math.max(0, Number(route.query.tag) || 0);
  page.value = Math.min(1000000, Math.max(1, Number(route.query.page) || 1));
}
async function load() {
  const request = ++sequence;
  loading.value = true;
  error.value = false;
  selected.value = [];
  try {
    const response = await service.get("/admin/article/getArticleList", {
      params: {
        page: page.value,
        pageSize: pageSize.value,
        state: filters.state === "top" ? "all" : filters.state,
        top: filters.state === "top" || undefined,
        keyword: filters.keyword || undefined,
        fid: filters.fid || undefined,
        tagId: filters.tag || undefined,
      },
    });
    if (request !== sequence) return;
    rows.value = Array.isArray(response.data?.list) ? response.data.list : [];
    total.value = response.data?.total || 0;
    if (page.value > 1 && !rows.value.length && total.value > 0)
      changePage(Math.ceil(total.value / pageSize.value));
  } catch {
    if (request === sequence) error.value = true;
  } finally {
    if (request === sequence) loading.value = false;
  }
}
async function changePage(next: number) {
  const query = {
    state: filters.state !== "all" ? filters.state : undefined,
    keyword: filters.keyword.trim() || undefined,
    fid: filters.fid || undefined,
    tag: filters.tag || undefined,
    page: next > 1 ? next : undefined,
  };
  const before = route.fullPath;
  await router.push({ query });
  if (route.fullPath === before) void load();
}
function applyFilters() {
  void changePage(1);
}
function changeState(value: string) {
  filters.state = value;
  applyFilters();
}
function resetFilters() {
  Object.assign(filters, { state: "all", keyword: "", fid: 0, tag: 0 });
  applyFilters();
}
function selectPage(event: Event) {
  selected.value = (event.target as HTMLInputElement).checked
    ? rows.value.map((a) => a.ID)
    : [];
}
const formatTime = (seconds?: number) =>
  seconds ? dayjs.unix(seconds).format("YYYY-MM-DD HH:mm") : "—";
const categoryName = (id: number) =>
  categories.value.find((c) => c.ID === id)?.name || "未分类";
const statusName = (status: number) =>
  articleStatusName[status as ArticleStatus] || "未知状态";
const statusClass = (status: number) =>
  ({
    0: "badge-light-success",
    1: "badge-light",
    2: "badge-light-warning",
    3: "badge-light-primary",
  }[status] || "badge-light");
async function patch(
  article: ArticleList,
  fields: { status?: number; isTop?: number }
) {
  if (busy.value) return;
  busy.value = true;
  try {
    await service.patch("/admin/article/state", {
      id: article.ID,
      version: article.version,
      ...fields,
    });
    await load();
  } catch {
  } finally {
    busy.value = false;
  }
}
async function setVisibility(article: ArticleList, status: number) {
  const answer = await Swal.fire({
    title: status === 0 ? "重新公开这篇文章？" : "隐藏这篇文章？",
    text:
      status === 0
        ? "文章将向所有访客公开。"
        : "文章将从信息流、搜索、订阅和站点地图中隐藏。",
    showCancelButton: true,
    confirmButtonText: "确认",
    cancelButtonText: "取消",
  });
  if (answer.isConfirmed) await patch(article, { status });
}
const setTop = (article: ArticleList) =>
  patch(article, { isTop: article.isTop ? 0 : 1 });
async function deleteIDs(ids: number[]) {
  busy.value = true;
  try {
    await service.delete("/admin/article/delArticleByIds", { data: { ids } });
    ElNotification({ title: "文章已删除", type: "success" });
    await load();
  } catch {
  } finally {
    busy.value = false;
  }
}
async function removeArticle(article: ArticleList) {
  const answer = await Swal.fire({
    title: `删除“${article.title}”？`,
    text: "删除后无法从文章列表恢复，请确认已备份需要保留的内容。",
    icon: "warning",
    showCancelButton: true,
    confirmButtonText: "删除文章",
    cancelButtonText: "取消",
  });
  if (answer.isConfirmed) await deleteIDs([article.ID]);
}
async function removeSelected() {
  const ids = [...selected.value];
  const answer = await Swal.fire({
    title: `删除所选 ${ids.length} 篇文章？`,
    text: "此操作不可撤销。",
    icon: "warning",
    showCancelButton: true,
    confirmButtonText: "确认删除",
    cancelButtonText: "取消",
  });
  if (answer.isConfirmed) await deleteIDs(ids);
}
onMounted(async () => {
  syncRoute();
  void load();
  const [categoryResponse, tagResponse] = await Promise.allSettled([
    service.get("/admin/category/getCategoryList"),
    service.get("/admin/article/tags"),
  ]);
  if (categoryResponse.status === "fulfilled")
    categories.value = Array.isArray(categoryResponse.value.data)
      ? categoryResponse.value.data
      : [];
  if (tagResponse.status === "fulfilled")
    tags.value = Array.isArray(tagResponse.value.data)
      ? tagResponse.value.data
      : [];
});
watch(
  () => route.fullPath,
  () => {
    syncRoute();
    void load();
  }
);
onBeforeUnmount(() => {
  sequence++;
});
</script>
<style scoped>
.library-select {
  width: 170px;
}
.library-title {
  min-width: 230px;
  max-width: 380px;
}
.library-table th,
.library-table td {
  padding-top: 1rem;
  padding-bottom: 1rem;
  vertical-align: middle;
}
.library-table {
  min-width: 850px;
}
</style>
