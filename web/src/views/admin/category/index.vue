<template>
  <!--begin::Tables Widget 11-->
  <div class="card">
    <!--begin::Header-->
    <div class="card-header border-0 pt-5">
      <h3 class="card-title align-items-start flex-column">
        <span class="card-label fw-bold fs-3 mb-1">分类管理</span>

        <span class="text-muted mt-1 fw-semobold fs-7">快速管理前台分类列表</span>
      </h3>
      <div class="card-toolbar">
        <button
            type="button"
            class="btn btn-primary"
            @click="add"
          >
            <span class="svg-icon svg-icon-2">
              <inline-svg
                :src="getAssetPath('/media/icons/duotune/arrows/arr075.svg')"
              />
            </span>
            新增分类
          </button>
      </div>
    </div>
    <!--end::Header-->

    <!--begin::Body-->
    <div class="card-body py-3">
      <!--begin::Table container-->
      <div class="table-responsive">
        <!--begin::Table-->
        <table class="table ">
          <!--begin::Table head-->
          <thead>
            <tr class="fw-bold text-muted bg-light">
              <th class="ps-4 min-w-200px rounded-start">分类名称</th>
              <th class="min-w-100px">权重</th>
              <th class="min-w-100px">ID</th>
              <th class="min-w-150px text-end rounded-end pe-4">操作</th>
            </tr>
          </thead>
          <!--end::Table head-->

          <!--begin::Table body-->
          <tbody>
            <template v-for="(item, index) in category">
              <tr v-if="item.fid === 0">
                <td>
                  <div class="ps-4 d-flex align-items-center">
                    <div class="d-flex justify-content-start flex-column">
                      <a href="#" class="text-dark fw-bold text-hover-primary mb-1 fs-6">{{ item.name }}</a>
                      <span class="text-muted fw-semobold text-muted d-block fs-7">一级目录</span>
                    </div>
                  </div>
                </td>

                <td>
                  <a href="#" class="text-dark fw-bold text-hover-primary d-block mb-1 fs-6">{{ item.sort }}</a>
                  <span class="text-muted fw-semobold text-muted d-block fs-7">Sort</span>
                </td>

                <td>
                  <a href="#" class="text-dark fw-bold text-hover-primary d-block mb-1 fs-6">{{ item.type }}</a>
                  <span class="text-muted fw-semobold text-muted d-block fs-7">Rejected</span>
                </td>

                <td class="text-end">
                  <button
                      type="button" 
                      class="btn btn-icon btn-bg-light btn-active-color-primary btn-sm me-1"
                      @click="modify(item.ID)">
                      <span class="svg-icon svg-icon-3">
                        <inline-svg :src="getAssetPath('/media/icons/duotune/art/art005.svg')
                          " />
                      </span>
                    </button>

                    <button 
                      type="button"
                      class="btn btn-icon btn-bg-light btn-active-color-primary btn-sm">
                      <span class="svg-icon svg-icon-3">
                        <inline-svg :src="getAssetPath(
                            '/media/icons/duotune/general/gen027.svg'
                          )
                          " />
                      </span>
                    </button>
                </td>
              </tr>
              <template v-for="sub in category">
                <tr v-if="sub.fid !== 0 && sub.fid == item.ID">
                  <td>
                    <div class="ps-4 d-flex align-items-center">
                      <div class="d-flex justify-content-start flex-column">
                        <a href="#" class="text-dark fw-bold text-hover-primary mb-1 fs-6">|—{{ sub.name }}</a>
                        <span class="text-muted fw-semobold text-muted d-block fs-7">二级目录</span>
                      </div>
                    </div>
                  </td>

                  <td>
                    <a href="#" class="text-dark fw-bold text-hover-primary d-block mb-1 fs-6">{{ sub.sort }}</a>
                    <span class="text-muted fw-semobold text-muted d-block fs-7">Sort</span>
                  </td>

                  <td>
                    <a href="#" class="text-dark fw-bold text-hover-primary d-block mb-1 fs-6">{{ sub.ID }}</a>
                    <span class="text-muted fw-semobold text-muted d-block fs-7">Rejected</span>
                  </td>

                  <td class="text-end">
                  <button
                      type="button" 
                      class="btn btn-icon btn-bg-light btn-active-color-primary btn-sm me-1"
                      @click="modify(sub.ID)">
                      <span class="svg-icon svg-icon-3">
                        <inline-svg :src="getAssetPath('/media/icons/duotune/art/art005.svg')
                          " />
                      </span>
                    </button>

                    <button 
                      type="button"
                      class="btn btn-icon btn-bg-light btn-active-color-primary btn-sm">
                      <span class="svg-icon svg-icon-3">
                        <inline-svg :src="getAssetPath(
                            '/media/icons/duotune/general/gen027.svg'
                          )
                          " />
                      </span>
                    </button>
                </td>
                </tr>

              </template>
            </template>
          </tbody>
          <!--end::Table body-->
        </table>
        <!--end::Table-->
      </div>
      <!--end::Table container-->
    </div>
    <!--begin::Body-->
  </div>
  <!--end::Tables Widget 11-->
  <!-- 绑定值 -->
  <CategoryModal :mode="modalMode" :data="category" :id="modalId"></CategoryModal>

</template>
  
<script lang="ts">
import { getAssetPath } from "@/core/helpers/assets";
import service from "@/utils/request";
import { onMounted, ref,defineComponent } from "vue";
import CategoryModal from "./modal.vue";
import { Modal } from "bootstrap";
import type { CategoryList } from "@/core/blog/CategoryTypes"

export default defineComponent({
  name: "admin-category",
  components: {
    CategoryModal,
  },
  setup() {
    const category = ref<Array<CategoryList>>([])
    const modalMode = ref("")
    const modalId = ref(0)

    onMounted(() => {
      getCategoryList()
    })

    const getCategoryList = () => {
      service.get("/admin/category/getCategoryList")
        .then(({ data }) => {
          if (data !== undefined) {
            category.value = data
          }
        })
    }
    
    const modify = (id:number) => {
      const modal = new Modal('#kt_modal_category');
      modalMode.value = "modify";
      modalId.value = id;
      modal.show();
    }

    const add = () => {
      const modal = new Modal('#kt_modal_category');
      modalMode.value = "add";
      modalId.value = 0;
      modal.show();
    }


    return {
      modalId,
      modalMode,
      category,
      add,
      modify,
      getAssetPath,
    };
  },
});
</script>
  