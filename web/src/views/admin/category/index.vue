<template>
  <div class="category-page">
    <div class="category-hero card border-0 mb-6 overflow-hidden">
      <div class="card-body p-7 p-lg-9 position-relative">
        <div class="hero-glow hero-glow-one"></div>
        <div class="hero-glow hero-glow-two"></div>
        <div class="d-flex flex-column flex-lg-row align-items-lg-center justify-content-between position-relative gap-6">
          <div>
            <div class="d-flex align-items-center gap-3 mb-3">
              <span class="hero-icon"><i class="bi bi-diagram-3-fill"></i></span>
              <span class="text-uppercase fw-bold fs-8 letter-spacing text-primary">Content taxonomy</span>
            </div>
            <h1 class="fs-2x fw-bolder text-gray-900 mb-2">分类管理</h1>
            <p class="text-gray-600 fs-6 mb-0">整理内容层级，让读者更快找到感兴趣的文章。</p>
          </div>
          <button type="button" class="btn btn-primary btn-lg px-6 shadow-sm" @click="add()">
            <i class="bi bi-plus-lg fs-4 me-2"></i>新建分类
          </button>
        </div>
      </div>
    </div>

    <div class="row g-5 mb-6">
      <div v-for="stat in stats" :key="stat.label" class="col-12 col-sm-4">
        <div class="card h-100 border-0 stat-card">
          <div class="card-body d-flex align-items-center p-5">
            <span class="stat-icon" :class="`bg-light-${stat.color} text-${stat.color}`">
              <i :class="[stat.icon, `text-${stat.color}`]"></i>
            </span>
            <div class="ms-4">
              <div class="fs-2 fw-bolder text-gray-900 lh-1 mb-2">{{ stat.value }}</div>
              <div class="text-muted fw-semibold fs-7">{{ stat.label }}</div>
            </div>
          </div>
        </div>
      </div>
    </div>

    <div class="card border-0 category-list-card">
      <div class="card-header border-0 py-6 min-h-auto">
        <div class="card-title flex-column align-items-start m-0">
          <h3 class="fw-bold text-gray-900 mb-1">分类结构</h3>
          <span class="text-muted fs-7">拖拽排序暂未开放，可通过权重控制展示顺序</span>
        </div>
        <div class="card-toolbar d-flex gap-3">
          <div class="position-relative search-box">
            <i class="bi bi-search position-absolute top-50 translate-middle-y ms-4 text-gray-500"></i>
            <input v-model.trim="keyword" class="form-control form-control-solid ps-11" placeholder="搜索分类名称" />
          </div>
          <button class="btn btn-light-primary btn-icon" title="刷新" :disabled="loading" @click="getCategoryList">
            <i class="bi bi-arrow-clockwise fs-4" :class="{ 'spin': loading }"></i>
          </button>
        </div>
      </div>

      <div class="card-body pt-0">
        <div v-if="loading && category.length === 0" class="empty-state">
          <span class="spinner-border text-primary mb-4"></span>
          <span class="text-muted">正在加载分类...</span>
        </div>

        <div v-else-if="filteredGroups.length === 0" class="empty-state">
          <span class="empty-icon mb-4"><i class="bi bi-folder2-open"></i></span>
          <h4 class="fw-bold text-gray-800 mb-2">没有找到分类</h4>
          <p class="text-muted mb-5">{{ keyword ? '换个关键词试试看' : '创建第一个分类，开始整理博客内容' }}</p>
          <button v-if="!keyword" class="btn btn-sm btn-primary" @click="add()">新建分类</button>
        </div>

        <div v-else class="category-tree">
          <section v-for="group in filteredGroups" :key="group.root.ID" class="category-group">
            <div class="category-row root-row">
              <button class="expand-button" :class="{ collapsed: !isExpanded(group.root.ID) }" @click="toggle(group.root.ID)">
                <i class="bi bi-chevron-down"></i>
              </button>
              <span class="folder-icon"><i class="bi bi-folder-fill"></i></span>
              <div class="category-main">
                <div class="d-flex align-items-center flex-wrap gap-2">
                  <span class="category-name">{{ group.root.name }}</span>
                  <span class="badge badge-light-primary">一级分类</span>
                </div>
                <span class="category-meta">ID {{ group.root.ID }} · {{ group.children.length }} 个子分类</span>
              </div>
              <div class="weight-block">
                <span class="weight-value">{{ group.root.sort }}</span>
                <span class="weight-label">权重</span>
              </div>
              <div class="row-actions">
                <button class="btn btn-sm btn-light-primary" @click="modify(group.root.ID)"><i class="bi bi-pencil me-2"></i>编辑</button>
                <button class="btn btn-sm btn-icon btn-light-danger" title="删除" @click="remove(group.root)"><i class="bi bi-trash3"></i></button>
              </div>
            </div>

            <div v-show="isExpanded(group.root.ID)" class="children-wrap">
              <div v-for="child in group.children" :key="child.ID" class="category-row child-row">
                <span class="tree-line"><i class="bi bi-arrow-return-right"></i></span>
                <span class="child-icon"><i class="bi bi-tag-fill"></i></span>
                <div class="category-main">
                  <span class="category-name fs-6">{{ child.name }}</span>
                  <span class="category-meta">ID {{ child.ID }} · 隶属于 {{ group.root.name }}</span>
                </div>
                <div class="weight-block">
                  <span class="weight-value">{{ child.sort }}</span>
                  <span class="weight-label">权重</span>
                </div>
                <div class="row-actions">
                  <button class="btn btn-sm btn-light-primary" @click="modify(child.ID)"><i class="bi bi-pencil me-2"></i>编辑</button>
                  <button class="btn btn-sm btn-icon btn-light-danger" title="删除" @click="remove(child)"><i class="bi bi-trash3"></i></button>
                </div>
              </div>
              <button class="add-child-row" @click="add(group.root.ID)">
                <i class="bi bi-plus-circle me-2"></i>在“{{ group.root.name }}”下添加子分类
              </button>
            </div>
          </section>
        </div>
      </div>
    </div>

    <CategoryModal :mode="modalMode" :data="category" :id="modalId" :parent-id="modalParentId" @saved="getCategoryList" />
  </div>
</template>

<script lang="ts">
import { computed, defineComponent, onMounted, ref } from "vue";
import { Modal } from "bootstrap";
import Swal from "sweetalert2";
import { ElNotification } from "element-plus";
import service from "@/utils/request";
import type { CategoryList } from "@/core/blog/CategoryTypes";
import CategoryModal from "./modal.vue";

export default defineComponent({
  name: "admin-category",
  components: { CategoryModal },
  setup() {
    const category = ref<Array<CategoryList>>([]);
    const modalMode = ref<"add" | "modify">("add");
    const modalId = ref(0);
    const modalParentId = ref(0);
    const keyword = ref("");
    const loading = ref(false);
    const expandedIds = ref<Set<number>>(new Set());

    const groups = computed(() => category.value
      .filter(item => item.fid === 0)
      .map(root => ({
        root,
        children: category.value.filter(item => item.fid === root.ID),
      })));

    const filteredGroups = computed(() => {
      const query = keyword.value.toLowerCase();
      if (!query) return groups.value;
      return groups.value
        .map(group => ({
          root: group.root,
          children: group.root.name.toLowerCase().includes(query)
            ? group.children
            : group.children.filter(item => item.name.toLowerCase().includes(query)),
        }))
        .filter(group => group.root.name.toLowerCase().includes(query) || group.children.length > 0);
    });

    const stats = computed(() => [
      { label: "全部分类", value: category.value.length, icon: "bi bi-collection-fill", color: "primary" },
      { label: "一级分类", value: groups.value.length, icon: "bi bi-folder-fill", color: "success" },
      { label: "子分类", value: category.value.filter(item => item.fid !== 0).length, icon: "bi bi-tags-fill", color: "info" },
    ]);

    const getCategoryList = async () => {
      loading.value = true;
      try {
        const { data } = await service.get("/admin/category/getCategoryList");
        category.value = data || [];
        expandedIds.value = new Set(category.value.filter(item => item.fid === 0).map(item => item.ID));
      } finally {
        loading.value = false;
      }
    };

    const showModal = () => new Modal("#kt_modal_category").show();
    const modify = (id: number) => {
      modalMode.value = "modify";
      modalId.value = id;
      modalParentId.value = 0;
      showModal();
    };
    const add = (parentId = 0) => {
      modalMode.value = "add";
      modalId.value = 0;
      modalParentId.value = parentId;
      showModal();
    };

    const remove = async (item: CategoryList) => {
      const childCount = category.value.filter(categoryItem => categoryItem.fid === item.ID).length;
      if (childCount > 0) {
        await Swal.fire({ title: "暂时不能删除", text: `“${item.name}”下还有 ${childCount} 个子分类，请先移动或删除它们。`, icon: "info", confirmButtonText: "知道了" });
        return;
      }
      const result = await Swal.fire({
        title: `删除“${item.name}”？`, text: "删除后无法恢复，请确认该分类已不再使用。", icon: "warning",
        showCancelButton: true, confirmButtonText: "确认删除", cancelButtonText: "取消",
        customClass: { confirmButton: "btn btn-danger", cancelButton: "btn btn-light" }, buttonsStyling: false,
      });
      if (!result.isConfirmed) return;
      await service.delete("/admin/category/delCategory", { data: { ID: item.ID } });
      ElNotification({ title: "操作成功", message: "分类已删除", type: "success" });
      await getCategoryList();
    };

    const toggle = (id: number) => {
      const next = new Set(expandedIds.value);
      next.has(id) ? next.delete(id) : next.add(id);
      expandedIds.value = next;
    };
    const isExpanded = (id: number) => expandedIds.value.has(id) || Boolean(keyword.value);

    onMounted(getCategoryList);
    return { category, modalMode, modalId, modalParentId, keyword, loading, stats, filteredGroups, getCategoryList, modify, add, remove, toggle, isExpanded };
  },
});
</script>

<style lang="scss" scoped>
.category-page { max-width: 1320px; margin: 0 auto; }
.category-hero { background: linear-gradient(120deg, #f2f6ff 0%, #fff 62%, #f1edff 100%); box-shadow: 0 12px 35px rgba(31, 42, 74, .06); }
.hero-glow { position: absolute; border-radius: 50%; filter: blur(2px); opacity: .65; }
.hero-glow-one { width: 180px; height: 180px; right: 11%; top: -100px; background: rgba(54, 153, 255, .16); }
.hero-glow-two { width: 120px; height: 120px; right: 2%; bottom: -70px; background: rgba(114, 57, 234, .14); }
.hero-icon, .stat-icon { display: inline-flex; align-items: center; justify-content: center; border-radius: 12px; }
.hero-icon { width: 42px; height: 42px; background: #fff; color: var(--kt-primary); box-shadow: 0 8px 20px rgba(54, 153, 255, .15); font-size: 1.25rem; }
.letter-spacing { letter-spacing: .12em; }
.stat-card, .category-list-card { box-shadow: 0 10px 32px rgba(31, 42, 74, .055); }
.stat-icon { width: 52px; height: 52px; font-size: 1.35rem; flex: 0 0 auto; }
.search-box { width: 240px; }
.category-tree { border: 1px solid var(--kt-gray-200); border-radius: 14px; overflow: hidden; }
.category-group + .category-group { border-top: 1px solid var(--kt-gray-200); }
.category-row { display: grid; grid-template-columns: 36px 46px minmax(220px, 1fr) 90px 150px; align-items: center; gap: 10px; min-height: 86px; padding: 15px 20px; transition: background-color .18s ease; }
.category-row:hover { background: var(--kt-gray-100); }
.root-row { background: var(--kt-body-bg); }
.child-row { min-height: 72px; background: rgba(248, 249, 252, .65); border-top: 1px dashed var(--kt-gray-300); padding-left: 38px; }
.expand-button { width: 32px; height: 32px; border: 0; border-radius: 9px; background: var(--kt-gray-100); color: var(--kt-gray-600); transition: transform .2s ease, background .2s ease; }
.expand-button:hover { background: var(--kt-primary-light); color: var(--kt-primary); }
.expand-button.collapsed i { display: inline-block; transform: rotate(-90deg); }
.folder-icon, .child-icon { width: 42px; height: 42px; display: inline-flex; align-items: center; justify-content: center; border-radius: 11px; font-size: 1.2rem; }
.folder-icon { color: #f6a723; background: #fff5df; }
.child-icon { width: 36px; height: 36px; color: var(--kt-info); background: var(--kt-info-light); font-size: 1rem; }
.tree-line { color: var(--kt-gray-400); text-align: center; font-size: 1.15rem; }
.category-main { display: flex; min-width: 0; flex-direction: column; gap: 5px; }
.category-name { color: var(--kt-gray-900); font-size: 1.08rem; font-weight: 700; overflow: hidden; text-overflow: ellipsis; white-space: nowrap; }
.category-meta { color: var(--kt-gray-500); font-size: .85rem; font-weight: 500; }
.weight-block { display: flex; flex-direction: column; align-items: center; }
.weight-value { color: var(--kt-gray-800); font-size: 1rem; font-weight: 700; }
.weight-label { color: var(--kt-gray-500); font-size: .75rem; }
.row-actions { display: flex; justify-content: flex-end; gap: 8px; }
.add-child-row { width: 100%; border: 0; border-top: 1px dashed var(--kt-gray-300); background: rgba(248, 249, 252, .65); color: var(--kt-primary); padding: 14px 20px 14px 92px; text-align: left; font-weight: 600; transition: background .18s ease; }
.add-child-row:hover { background: var(--kt-primary-light); }
.empty-state { min-height: 330px; display: flex; flex-direction: column; align-items: center; justify-content: center; text-align: center; }
.empty-icon { width: 72px; height: 72px; display: inline-flex; align-items: center; justify-content: center; border-radius: 18px; background: var(--kt-primary-light); color: var(--kt-primary); font-size: 2rem; }
.spin { animation: spin .8s linear infinite; }
@keyframes spin { to { transform: rotate(360deg); } }
@media (max-width: 767.98px) {
  .card-header { align-items: stretch; }
  .card-toolbar, .search-box { width: 100%; }
  .category-row { grid-template-columns: 32px 38px minmax(0, 1fr) auto; padding: 13px 12px; }
  .weight-block { display: none; }
  .row-actions { grid-column: 3 / 5; justify-content: flex-start; }
  .child-row { padding-left: 15px; }
  .add-child-row { padding-left: 28px; }
}
</style>
