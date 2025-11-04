<template>
  <div
    class="modal fade"
    id="kt_modal_category"
    ref="categoryModalRef"
    tabindex="-1"
    aria-hidden="true"
  >
    <!--begin::Modal dialog-->
    <div class="modal-dialog modal-dialog-centered mw-650px">
      <!--begin::Modal content-->
      <div class="modal-content">
        <!--begin::Modal header-->
        <div class="modal-header" id="kt_modal_category_header">
          <!--begin::Modal title-->
          <h2 class="fw-bold">{{ getModalTitle() }}分类</h2>
          <!--end::Modal title-->

          <!--begin::Close-->
          <div
            id="kt_modal_category_close"
            data-bs-dismiss="modal"
            class="btn btn-icon btn-sm btn-active-icon-primary"
          >
            <span class="svg-icon svg-icon-1">
              <inline-svg
                :src="getAssetPath('/media/icons/duotune/arrows/arr061.svg')"
              />
            </span>
          </div>
          <!--end::Close-->
        </div>
        <!--end::Modal header-->
        <!--begin::Form-->
        <el-form
          @submit.prevent="submit()"
          :model="formData"
          :rules="rules"
          ref="formRef"
        >
          <!--begin::Modal body-->
          <div class="modal-body py-10 px-lg-17">
            <!--begin::Scroll-->
            <div
              class="scroll-y me-n7 pe-7"
              id="kt_modal_category_scroll"
              data-kt-scroll="true"
              data-kt-scroll-activate="{default: false, lg: true}"
              data-kt-scroll-max-height="auto"
              data-kt-scroll-dependencies="#kt_modal_category_header"
              data-kt-scroll-wrappers="#kt_modal_category_scroll"
              data-kt-scroll-offset="300px"
            >
              <!--begin::Input group-->
              <div class="fv-row mb-7">
                <!--begin::Label-->
                <label class="required fs-6 fw-semobold mb-2">分类名称</label>
                <!--end::Label-->

                <!--begin::Input-->
                <el-form-item prop="name">
                  <el-input
                    v-model="formData.name"
                    type="text"
                    placeholder="请输入分类"
                  />
                </el-form-item>
                <!--end::Input-->
              </div>
              <!--end::Input group-->

              <!--begin::Input group-->
              <div class="fv-row mb-7">
                <!--begin::Label-->
                <label class="fs-6 fw-semobold mb-2">
                  <span class="required">上级分类</span>

                  <i
                    class="fas fa-exclamation-circle ms-1 fs-7"
                    data-bs-toggle="tooltip"
                    title="Email address must be active"
                  ></i>
                </label>
                <!--end::Label-->

                <!--begin::Input-->
                <el-form-item prop="fid">
                  <el-select style="width: 100%;" v-model="formData.ID">
                    <el-option :value="0" label="根目录" />
                    <template v-for="item in data">
                      <el-option v-if="item.fid == 0" :value="item.ID" :label="item.name" />
                    </template>
                    

                  </el-select>
                </el-form-item>
                <!--end::Input-->
              </div>
              <!--end::Input group-->

              <!--begin::Input group-->
              <div class="fv-row mb-15">
                <!--begin::Label-->
                <label class="fs-6 fw-semobold mb-2">分类权重</label>
                <!--end::Label-->

                <!--begin::Input-->
                <el-form-item prop="description">
                  <el-input v-model="formData.sort" type="text" />
                </el-form-item>
                <!--end::Input-->
              </div>
              <!--end::Input group-->

            </div>
            <!--end::Scroll-->
          </div>
          <!--end::Modal body-->

          <!--begin::Modal footer-->
          <div class="modal-footer flex-center">
            <!--begin::Button-->
            <button
              type="reset"
              id="kt_modal_category_cancel"
              class="btn btn-light me-3"
            >
              重置
            </button>
            <!--end::Button-->

            <!--begin::Button-->
            <button
              :data-kt-indicator="loading ? 'on' : null"
              class="btn btn-lg btn-primary"
              type="submit"
            >
              <span v-if="!loading" class="indicator-label">
                提交
                <span class="svg-icon svg-icon-3 ms-2 me-0">
                  <inline-svg
                    :src="
                      getAssetPath('/media/icons/duotune/arrows/arr064.svg')
                    "
                  />
                </span>
              </span>
              <span v-if="loading" class="indicator-progress">
                Please wait...
                <span
                  class="spinner-border spinner-border-sm align-middle ms-2"
                ></span>
              </span>
            </button>
            <!--end::Button-->
          </div>
          <!--end::Modal footer-->
        </el-form>
        <!--end::Form-->
      </div>
    </div>
  </div>
</template>

<script lang="ts">
import { getAssetPath } from "@/core/helpers/assets";
import { defineComponent, ref } from "vue";
import { hideModal } from "@/core/helpers/dom";
import Swal from "sweetalert2";
import type { CategoryList } from "@/core/blog/CategoryTypes";
import { watch } from "vue"; 

export default defineComponent({
  name: "category-modal",
  components: {},
  props: {
    id:Number,
    mode:String,
    data:Array<CategoryList>,
  },
  setup(props) {
    const formRef = ref<null | HTMLFormElement>(null);
    const categoryModalRef = ref<null | HTMLElement>(null);
    const loading = ref<boolean>(false);
    const formData = ref({
      ID: 0,
      name: "",
      sort: 0,
    });
  
    watch(
      () => props.id,
      (newValue: any) => {
        props.data?.forEach(element => {
          if(props.id == element.ID){
            formData.value.ID = element.fid;
            formData.value.name = element.name;
            formData.value.sort = element.sort;
          }else if(props.id == 0){
            formData.value.ID = 0;
            formData.value.name = "";
            formData.value.sort = 0;
          }
        });
      },
    )


    const rules = ref({
      name: [
        {
          required: true,
          message: "必须输入分类名称",
          trigger: "change",
        },
      ],
    });

    const getModalTitle = () => {
      return props.mode == "add" ? "新增" : "修改"
    }

    const submit = () => {
      if (!formRef.value) {
        return;
      }

      formRef.value.validate((valid: boolean) => {
        if (valid) {
          loading.value = true;

          setTimeout(() => {
            loading.value = false;

            Swal.fire({
              text: "Form has been successfully submitted!",
              icon: "success",
              buttonsStyling: false,
              confirmButtonText: "Ok, got it!",
              heightAuto: false,
              customClass: {
                confirmButton: "btn btn-primary",
              },
            }).then(() => {
              hideModal(categoryModalRef.value);
            });
          }, 2000);
        } else {
          Swal.fire({
            text: "Sorry, looks like there are some errors detected, please try again.",
            icon: "error",
            buttonsStyling: false,
            confirmButtonText: "Ok, got it!",
            heightAuto: false,
            customClass: {
              confirmButton: "btn btn-primary",
            },
          });
          return false;
        }
      });
    };

    return {
      getModalTitle,
      formData,
      rules,
      submit,
      formRef,
      loading,
      categoryModalRef,
      getAssetPath,
    };
  },
});
</script>
