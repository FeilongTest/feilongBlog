<template>
    <!--begin::Feeds Widget 5-->
    <div class="card" :class="widgetClasses">
      <!--begin::Body-->
      <div class="card-body pb-0">
        <!--begin::Header-->
        <div class="d-flex align-items-center mb-5">
          <!--begin::User-->
          <div class="d-flex align-items-center flex-grow-1">
            <!--begin::Avatar-->
            <div class="symbol symbol-45px me-5">
              <img :src="getAssetPath('/media/avatars/300-1.jpg')" alt="" />
            </div>
            <!--end::Avatar-->
  
            <!--begin::Info-->
            <div class="d-flex flex-column">

              <router-link :to="'/content/'+ article?.fid + '/'+article?.ID" class="title-text text-gray-800 text-hover-primary fs-6 fw-bold"
                >
                <span class="badge badge-light-success fw-bold my-2 me-3">{{ article?.isTop === 1?"置顶":""}}</span>
                {{article?.title}}
              </router-link>
  
              <span class="text-gray-400 fw-semobold">{{moment.unix(Number(article?.ctime)).format("YYYY-MM-DD HH:mm:ss")}}</span>
            </div>
            <!--end::Info-->
          </div>
          <!--end::User-->
  
          <!--begin::Menu-->
          <div class="my-0">
            <button
              type="button"
              class="btn btn-sm btn-icon btn-color-primary btn-active-light-primary"
              data-kt-menu-trigger="click"
              data-kt-menu-placement="bottom-end"
              data-kt-menu-flip="top-end"
            >
              <span class="svg-icon svg-icon-2">
                <inline-svg
                  :src="getAssetPath('/media/icons/duotune/general/gen024.svg')"
                />
              </span>
            </button>
  
            <Dropdown2></Dropdown2>
          </div>
          <!--end::Menu-->
        </div>
        <!--end::Header-->
  
        <!--begin::Post-->
        <div class="mb-5">
          <!--begin::Image-->
          
          <div v-show="article?.pic?true:false"
            class="bgi-no-repeat bgi-size-cover rounded min-h-250px mb-5"
            :style="{
              backgroundImage: `url(${article?.pic})`,
            }"
          ></div>
          <!--end::Image-->
  
          <!--begin::Text-->
          <div class="content-text text-gray-800 mb-5 max-line5">
            {{ delHtmlTag(article?.content) }}
          </div>
          <!--end::Text-->
  
          <!--begin::Toolbar-->
          <div class="d-flex align-items-center mb-5">
            <a
              class="btn btn-sm btn-light btn-color-muted btn-active-light-info px-4 py-2 me-4"
            >
              <span class="svg-icon svg-icon-2">
                <inline-svg
                  :src="getAssetPath('/media/icons/duotune/coding/cod006.svg')"
                />
              </span>
              {{article?.view}}
            </a>

            <a
              class="btn btn-sm btn-light btn-color-muted btn-active-light-success px-4 py-2 me-4"
            >
              <span class="svg-icon svg-icon-3">
                <inline-svg
                  :src="
                    getAssetPath('/media/icons/duotune/communication/com012.svg')
                  "
                />
              </span>
              {{article?.commentCount || 0}}
            </a>
  
            <a
              class="btn btn-sm btn-light btn-color-muted btn-active-light-danger px-4 py-2"
            >
              <span class="svg-icon svg-icon-2">
                <inline-svg
                  :src="getAssetPath('/media/icons/duotune/general/gen030.svg')"
                />
              </span>
              暂不支持
            </a>
          </div>
          <!--end::Toolbar-->
        </div>
        <!--end::Post-->
  
      </div>
      <!--end::Body-->
    </div>
    <!--end::Feeds Widget 5-->
  </template>
  
  <script lang="ts">
  import { getAssetPath } from "@/core/helpers/assets";
  import { defineComponent,type PropType } from "vue";
  import Dropdown2 from "@/components/dropdown/Dropdown2.vue";
  import type { ArticleList } from "@/core/blog/ArticleTypes";
  import { delHtmlTag } from "@/utils/html"
  import moment from "moment";
  
  export default defineComponent({
    name: "Blog-Feeds",
    props: {
      widgetClasses: String,
      article:Object as PropType<ArticleList>,
    },
    components: {
      Dropdown2,
    },
    setup() {

      return {
        getAssetPath,
        delHtmlTag,
        moment,
      };
    },
  });
  </script>
  
  <style scoped>

.title-text {
     display: -webkit-box;
    overflow: hidden; /*超出宽度部分隐藏*/
    text-overflow: ellipsis; /*超出部分以点号代替*/
    -webkit-line-clamp: 1; /*行数*/
    -webkit-box-orient: vertical;
    overflow: hidden;
}
.content-text {
    display: -webkit-box;
    overflow: hidden; /*超出宽度部分隐藏*/
    text-overflow: ellipsis; /*超出部分以点号代替*/
    -webkit-line-clamp: 1; /*行数*/
    -webkit-box-orient: vertical;
    overflow: hidden;
}
</style>