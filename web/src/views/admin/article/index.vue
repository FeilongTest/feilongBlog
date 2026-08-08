<template>
  <!--begin::Card-->
  <div class="card">
    <!--begin::Card header-->
    <div class="card-header border-0 pt-6">
      <!--begin::Card title-->
      <div class="card-title flex-wrap gap-3">
        <!--begin::Search-->
        <div class="d-flex align-items-center position-relative my-1">

          <input type="text" data-kt-subscription-table-filter="search"
            class="form-control form-control-solid w-230px ps-6" placeholder="输入标题关键字进行搜索" v-model="searchKeywords" @keyup.enter="applySearch"/>
          <span class="btn svg-icon svg-icon-1 position-absolute end-0" @click="applySearch">
            <inline-svg :src="getAssetPath('/media/icons/duotune/general/gen004.svg')" />
          </span>
        </div>
        <!--end::Search-->

        <el-select v-model="selectedCategory" class="category-filter" placeholder="全部分类" filterable @change="applyCategory">
          <el-option :value="0" label="全部分类" />
          <el-option-group v-for="group in categoryGroups" :key="group.root.ID" :label="group.root.name">
            <el-option :value="group.root.ID" :label="`${group.root.name}（一级分类）`" />
            <el-option v-for="child in group.children" :key="child.ID" :value="child.ID" :label="`　${child.name}`" />
          </el-option-group>
        </el-select>

        <button v-if="hasActiveFilters" type="button" class="btn btn-sm btn-light-primary" @click="resetFilters">
          <i class="bi bi-x-circle me-2"></i>清除筛选
        </button>
      </div>
      <!--begin::Card title-->

      <!--begin::Card toolbar-->
      <div class="card-toolbar">
        <!--begin::Toolbar-->
        <div v-if="selectedIds.length === 0" class="d-flex justify-content-end" data-kt-subscription-table-toolbar="base">
          <!--begin::Tab nav-->
          <ul class="nav nav-stretch fs-5 fw-semobold nav-line-tabs nav-line-tabs-2x border-transparent" role="tablist">
            <li class="nav-item" role="presentation">
              <a class="nav-link text-active-primary" :class="{ active: statusFilter === 'all' }" href="#" @click.prevent="setStatusFilter('all')">
                全部
              </a>
            </li>
            
            <li class="nav-item" role="presentation">
              <a class="nav-link text-active-primary" :class="{ active: statusFilter === 'top' }" href="#" @click.prevent="setStatusFilter('top')">
                仅置顶
              </a>
            </li>

            <li class="nav-item" role="presentation">
              <a class="nav-link text-active-primary ms-3" :class="{ active: statusFilter === 'hidden' }" href="#" @click.prevent="setStatusFilter('hidden')">
                仅隐藏
              </a>
            </li>

          </ul>
          <!--end::Tab nav-->
        </div>
        <!--end::Toolbar-->

        <!--begin::Group actions-->
        <div v-else class="d-flex justify-content-end align-items-center">
          <div class="fw-bold me-5">
            <span class="me-2">{{ selectedIds.length }}</span>Selected
          </div>
          <button type="button" class="btn btn-danger" @click="delArticleByIds()">
            Delete Selected
          </button>
        </div>
        <!--end::Group actions-->
      </div>
      <!--end::Card toolbar-->
    </div>
    <!--end::Card header-->

    <!--begin::Card body-->
    <div class="card-body pt-0">
      <KTDatatable @on-sort="sort" @on-items-select="onItemSelect" :data="articleList" :header="headerConfig"
        :checkbox-enabled="true" :total="total" :current-page="page" @page-change="pageChange"
        @on-items-per-page-change="itemPageChange" checkbox-label="ID">
        <template v-slot:id="{ row: article }">
          <div>{{ article.ID }}</div>
        </template>
        <template v-slot:title="{ row: article }">
          <router-link :to="`/admin/article/content/${article.ID}`" class="text-gray-800 text-hover-primary mb-1">
            {{ article.title }}
          </router-link>
        </template>
        <template v-slot:category="{ row: article }">
          <span class="badge badge-light-primary fw-semibold">{{ getCategoryName(article.fid) }}</span>
        </template>
        <template v-slot:view="{ row: article }">
          <div class="badge badge-light">{{ article.view }}</div>
        </template>
        <template v-slot:ctime="{ row: article }">
          {{ moment.unix(article.ctime).format('YYYY-MM-DD HH:mm:ss')/** 转换成13位时间戳 */ }}
        </template>
        <template v-slot:editTime="{ row: article }">
          {{ moment.unix(article.editTime).format('YYYY-MM-DD HH:mm:ss')/** 转换成13位时间戳 */ }}
        </template>
        <template v-slot:isTop="{ row: article }">
          <el-switch 
            v-model="article.isTop"
            active-text="Y"
            inactive-text="N"
            inline-prompt
            :active-value="1"
            :inactive-value="0"
            @change="()=>{updateArticle(article,'istop')}"
           />
        </template>
        <template v-slot:status="{ row: article }">
          <el-switch 

            v-model="article.status"
            style="--el-switch-on-color: #85d75b; --el-switch-off-color: #909399"
            :active-value="0"
            :inactive-value="1"
            @change="()=>{updateArticle(article,'status')}"
           />
        </template>

        <template v-slot:actions="{ row: article }">
          <div class="d-flex row">
            <router-link :to="'/admin/article/content/' + article.ID"
            class="btn btn-icon btn-bg-light btn-active-color-primary btn-sm me-1">
            <span class="svg-icon svg-icon-3">
              <inline-svg :src="getAssetPath('/media/icons/duotune/art/art005.svg')
                " />
            </span>
          </router-link>

          <a @click="delArticle(article.ID)" class="btn btn-icon btn-bg-light btn-active-color-primary btn-sm">
            <span class="svg-icon svg-icon-3">
              <inline-svg :src="getAssetPath(
                '/media/icons/duotune/general/gen027.svg'
              )
                " />
            </span>
          </a>
          </div>
        </template>
      </KTDatatable>
    </div>
    <!--end::Card body-->
  </div>
  <!--end::Card-->
</template>
  
<script lang="ts">
import { getAssetPath } from "@/core/helpers/assets";
import { computed, defineComponent, ref, onMounted } from "vue";
import KTDatatable from "@/components/kt-datatable/KTDataTable.vue";
import type { Sort } from "@/components/kt-datatable/table-partials/models";
import arraySort from "array-sort";
import service from "@/utils/request";
import type { ArticleList } from "@/core/blog/ArticleTypes"
import type { CategoryList } from "@/core/blog/CategoryTypes"
import { ElNotification } from 'element-plus'
import dayjs from "dayjs";
import Swal from "sweetalert2";



export default defineComponent({
  name: "admin-article-index",
  components: {
    KTDatatable,
  },
  setup() {
    const searchKeywords = ref("");
    const page = ref(1);
    const total = ref(0);
    const pageSize = ref(10);
    const category = ref<Array<CategoryList>>([]);
    const selectedCategory = ref(0);
    const statusFilter = ref<"all" | "top" | "hidden">("all");

    const articleList = ref<Array<ArticleList>>([]);

    const headerConfig = ref([
      {
        columnName: "ID",
        columnLabel: "id",
        sortEnabled: true,
      },
      {
        columnName: "文章名称",
        columnLabel: "title",
        sortEnabled: true,
      },
      {
        columnName: "所属分类",
        columnLabel: "category",
      },
      {
        columnName: "查看次数",
        columnLabel: "view",
        sortEnabled: true,
      },
      {
        columnName: "发布日期",
        columnLabel: "ctime",
        sortEnabled: true,
      },
      {
        columnName: "修改日期",
        columnLabel: "editTime",
        sortEnabled: true,
      },
      {
        columnName: "置顶",
        columnLabel: "isTop",
        sortEnabled: true,
      },
      {
        columnName: "显示/隐藏",
        columnLabel: "status",
        sortEnabled: true,
      },
      {
        columnName: "Actions",
        columnLabel: "actions",
      },
    ]);


    const deleteArticle = (id: number) => {
      for (let i = 0; i < articleList.value.length; i++) {
        if (articleList.value[i].ID === id) {
          articleList.value.splice(i, 1);
        }
      }
    };

    const sort = (sort: Sort) => {
      const reverse: boolean = sort.order === "asc";
      if (sort.label) {
        arraySort(articleList.value, sort.label, { reverse });
      }
    };

    const onItemSelect = (selectedItems: Array<number>) => {
      selectedIds.value = selectedItems;
    };

    const getArticleList = async () => {
      const params = {
          page: page.value,
          pageSize: pageSize.value,
          fid: selectedCategory.value || undefined,
          top: statusFilter.value === "top" || undefined,
          hide: statusFilter.value === "hidden" || undefined,
          keyword: searchKeywords.value || undefined,
      };
      service({
        url: "/admin/article/getArticleList",
        method: "get",
        params
      })
        .then(({ data }) => {
          if (data !== undefined) {
            total.value = data.total;
            articleList.value = data.list;
          }
        })
    }

    const getCategoryList = async () => {
      const { data } = await service.get("/admin/category/getCategoryList");
      category.value = data || [];
    };

    const categoryGroups = computed(() => category.value
      .filter(item => item.fid === 0)
      .map(root => ({ root, children: category.value.filter(item => item.fid === root.ID) })));

    const getCategoryName = (fid: number) => category.value.find(item => item.ID === fid)?.name || `分类 ${fid}`;
    const hasActiveFilters = computed(() => Boolean(selectedCategory.value || searchKeywords.value || statusFilter.value !== "all"));
    const refreshFromFirstPage = () => { page.value = 1; getArticleList(); };
    const applyCategory = () => refreshFromFirstPage();
    const applySearch = () => refreshFromFirstPage();
    const setStatusFilter = (value: "all" | "top" | "hidden") => { statusFilter.value = value; refreshFromFirstPage(); };
    const resetFilters = () => {
      selectedCategory.value = 0;
      searchKeywords.value = "";
      statusFilter.value = "all";
      refreshFromFirstPage();
    };

    const delArticle = async (id:number) => {
      Swal.fire(
        {
          title:"年轻人想好了吗,确定删除该文章?",
          icon:"warning",
          showCancelButton:true,
        }
      ).then(res => {
        if(res.isConfirmed){
          //从列表中删除
          deleteArticle(id);
          service({
            url: "/admin/article/delArticle",
            method: "delete",
            data: {
              id: id,
            }
          })
          .then(({ data }) => {
            if (data !== undefined) {
              ElNotification({
                duration: 2000,
                title: 'Success',
                message: "删除成功",
                type: 'success',
              })
            }
          })
        }
      })
    }

    const selectedIds = ref<Array<number>>([]);

    const delArticleByIds = async () => {
      Swal.fire(
        {
          title:"年轻人想好了吗,确定一下子删除这么多文章?",
          icon:"warning",
          buttonsStyling: false,
          confirmButtonText: "确定,我可能是疯了",
          heightAuto: false,
          showCancelButton:true,
          customClass: {
            confirmButton: "btn fw-semobold btn-light-danger",
            cancelButton: "btn fw-semobold btn-secondary",
          },
        }
      ).then(res => {
        if(res.isConfirmed){
          const ids:number[] = [];
          selectedIds.value.forEach((item) => {
            deleteArticle(item);
            ids.push(item)
          });
          selectedIds.value.length = 0;
          service({
            url: "/admin/article/delArticleByIds",
            method: "delete",
            data: {
              ids: ids,
            }
          }).then(({ data }) => {
            ElNotification({
              duration: 2000,
              title: 'Success',
              message: "批量删除成功",
              type: 'success',
            })
          })
        }
      })
    }


    const updateArticle = (data:ArticleList,mode:string) => {

      let msg = "";
      if(mode == "istop"){
        if(data.isTop == 1){
          msg = "置顶成功"
        }else{
          msg = "取消置顶"
        }
      }else{
        if(data.status == 1){
          msg = "状态修改为隐藏"
        }else{
          msg = "状态修改为显示"
        }
      }
      service({
        url:"/admin/article/updateArticle",
        method: "put",
        data:data
      }).then(({data}) => {
        if(data != undefined){
          ElNotification({ type: 'success', message: msg })
        }
      })
    }


    onMounted(async () => {
      await getCategoryList();
      getArticleList();
    })


    const pageChange = (newPage: number) => {
      page.value = newPage;
      getArticleList();
    }

    const itemPageChange = (newPage: number) => {
      pageSize.value = newPage;
      getArticleList();
    }

     



    return {
      articleList,
      searchKeywords,
      total,
      page,
      pageSize,
      pageChange,
      itemPageChange,
      headerConfig,
      sort,
      onItemSelect,
      selectedIds,
      deleteArticle,
      getAssetPath,
      moment: dayjs,
      delArticle,
      delArticleByIds,
      updateArticle,
      getArticleList,
      categoryGroups,
      selectedCategory,
      statusFilter,
      getCategoryName,
      hasActiveFilters,
      applyCategory,
      applySearch,
      setStatusFilter,
      resetFilters,
    };
  },
});
</script>

<style lang="scss" scoped>
.category-filter { width: 250px; }
:deep(.category-filter .el-select__wrapper) {
  min-height: 43px;
  border-radius: .625rem;
  background: var(--kt-gray-100);
  box-shadow: none;
}
@media (max-width: 767.98px) {
  .category-filter { width: 100%; }
}
</style>
