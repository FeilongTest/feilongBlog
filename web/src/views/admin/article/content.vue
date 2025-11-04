<template>
  <div class="d-flex flex-column flex-lg-row">

    <!-- 正文开始 -->
    <div class="flex-lg-row-fluid me-lg-15 order-2 order-lg-1 mb-10 mb-lg-0">

      <div class="card card-flush py-5 mb-5 mb-lg-10">
        <!--begin::Card body-->
        <div class="card-body pt-0">
          <VForm id="kt_modal_create_api_key_form" class="form">

            <!--begin::Input group-->
            <div class="mb-5 fv-row">
              <!--begin::Label-->
              <label class="required fs-5 fw-semobold mb-3">文章标题</label>
              <!--end::Label-->

              <!--begin::Input-->
              <Field type="text" class="form-control form-control-lg form-control-solid " placeholder="请输入文章标题" name="title"
                v-model="articleInfo.title" />
              <div class="fv-plugins-message-container">
                <div class="fv-help-block">
                  <ErrorMessage name="apiName" />
                </div>
              </div>
              <!--end::Input-->
            </div>
            <!--end::Input group-->

            <!--begin::Input group-->
            <div class="mb-5 fv-row">
              <!--begin::Col-->
              <div class="col-md-12 fv-row">

                <div class="row fv-row">

                  <div class="col-6">
                    <label class="required fs-5 fw-semobold mb-3">文章分类</label>

                    <Field class="form-select form-select-lg form-select-solid select2-hidden-accessible" placeholder="请选择文章分类" name="fid"
                      v-model="articleInfo.fid" as="select">
                      <template v-for="(item, i) in category" :key="i">
                        <option
                        v-if="item.fid === 0"
                      :label="item.name"
                      :value="item.ID"
                    ></option>
                    <template v-for="(menuItem, j) in category" :key="j">
                      <option
                      v-if="menuItem.fid === item.ID"
                      :label="`|—${menuItem.name}`"
                      :value="menuItem.ID"
                    ></option>
                    </template>
                      </template>

                      </Field>
                    <div class="fv-plugins-message-container">
                      <div class="fv-help-block">
                        <ErrorMessage name="apiName" />
                      </div>
                    </div>
                  </div>

                  <div class="col-6">
                    <label class="required fs-5 fw-semobold mb-3">文章样式</label>

                    <Field as="select" class="form-select form-select-lg form-select-solid" placeholder="请选择文章样式" name="type"
                      v-model="articleInfo.type" >
                        <option
                          label="普通文章"
                          value="1"
                        ></option>
                        <option
                          label="APP文章(待开发)"
                          value="2"
                        ></option>
                      </Field>
                    <div class="fv-plugins-message-container">
                      <div class="fv-help-block">
                        <ErrorMessage name="apiName" />
                      </div>
                    </div>
                  </div>
                </div>

              </div>
              <!--end::Col-->





            </div>
            <!--end::Input group-->


            <div class="editor—wrapper">
              <Toolbar style="border-bottom: 1px solid #ccc" :editor="editorRef" :defaultConfig="toolbarConfig"
                :mode="mode" />
              <Editor style="height: 500px; overflow-y: hidden;" v-model="valueHtml" :defaultConfig="editorConfig"
                :mode="mode" @onCreated="handleCreated" />
            </div>
          </VForm>

        </div>
        <!--end::Card body-->
      </div>
    </div>
    <!-- 正文结束 -->
    <!-- 侧边栏开始 -->
    <div class="flex-column flex-lg-row-auto w-100 w-lg-250px w-xl-300px mb-10 order-1 order-lg-2">
      <div class="card card-flush mb-0" data-kt-sticky="true" data-kt-sticky-name="view-article-upload"
        data-kt-sticky-offset="{default: false, lg: '200px'}" data-kt-sticky-width="{lg: '250px', xl: '300px'}"
        data-kt-sticky-left="auto" data-kt-sticky-top="20px" data-kt-sticky-animation="false" data-kt-sticky-zindex="95">
        <div class="card-body pt-0 fs-6">

          <!-- 缩略图开始 -->
          <div class="mt-7 mb-5 row">
            <!--begin::Label-->
            <label class="fs-5 mb-7 fw-semobold">缩略图</label>
            <!--end::Label-->

            <!--begin::Image input-->
            <div class="img-input img-input-outline" data-kt-image-input="true">

              <!--begin::Preview existing avatar-->
              <div class="img-input-wrapper rounded bgi-no-repeat bgi-size-contain bgi-position-x-center min-h-150px p-8"
                :style="`background-image: url(${getImg(data.img)})`"></div>
              <!--end::Preview existing avatar-->


              <!--begin::Label-->
              <label class="btn btn-icon btn-circle btn-active-color-primary w-25px h-25px bg-body shadow"
                data-kt-image-input-action="change" data-bs-toggle="tooltip" title="编辑图片" @change="changeImage">
                <i class="bi bi-pencil-fill fs-7"></i>

                <!--begin::Inputs-->
                <input type="file" name="avatar" accept=".png, .jpg, .jpeg" />
                <input type="hidden" name="avatar_remove" />
                <!--end::Inputs-->
              </label>
              <!--end::Label-->

              <!--begin::Remove-->
              <span class="btn btn-icon btn-circle btn-active-color-primary w-25px h-25px bg-body shadow"
                data-kt-image-input-action="remove" data-bs-toggle="tooltip" @click="removeImage" title="移除图片">
                <i class="bi bi-x fs-2"></i>
              </span>
              <!--end::Remove-->
            </div>
            <!--end::Image input-->

            <!--begin::Input group-->
            <div class="mt-5">
              <Field type="text" name="fname" class="form-control form-control-lg form-control-solid mb-3 mb-lg-0"
                placeholder="上传或者填写图片外链" v-model="data.imgStr" />
            </div>
            <!--end::Input group-->

            <!--begin::Hint-->
            <!-- <div class="form-text">Allowed file types: png, jpg, jpeg.</div> -->
            <!--end::Hint-->
          </div>
          <!-- 缩略图结束 -->

          <!--begin::Seperator-->
          <div class="separator separator-dashed"></div>
          <!--end::Seperator-->

          <!--begin::Input group-->
          <div class="mt-5">
            <button @click="submitBtn" class="form-control from-control-lg from-control-solid bg-info mb-3 mb-lg-0" style="color:aliceblue;">
              {{ articleMode }}
            </button>
          </div>
          <!--end::Input group-->
        </div>



      </div>
    </div>
    <!-- 侧边栏结束 -->
  </div>
</template>
<script lang="ts">
import '@wangeditor/editor/dist/css/style.css' // 引入 css

import { defineComponent,onBeforeUnmount, ref, shallowRef, onMounted,watch } from 'vue'
import { getAssetPath } from "@/core/helpers/assets";
import { Editor, Toolbar } from '@wangeditor/editor-for-vue'
import { ErrorMessage, Field, Form as VForm } from "vee-validate";
import service from '@/utils/request';
import JwtService from "@/core/services/JwtService";
import { ElNotification } from 'element-plus'
import { useRouter } from "vue-router";
import type { CategoryList } from "@/core/blog/CategoryTypes"
import type { ArticleList } from "@/core/blog/ArticleTypes"
import Swal from "sweetalert2";


export default defineComponent({
  name: "admin-article-content",
  components: {
    Editor,
    Toolbar,
    ErrorMessage,
    Field,
    VForm,
  },
  setup() {

    //为public关键字进行后端URL重定向
    const publicApi = "public/";
    const category = ref<Array<CategoryList>>([]);
    const articleInfo = ref<ArticleList>({} as ArticleList);
    const articleMode = ref("新增");



    // 编辑器实例，必须用 shallowRef
    const editorRef = shallowRef()

    // 内容 HTML
    const valueHtml = ref('')


    const router = useRouter();

    // 模拟 ajax 异步获取内容
    onMounted(() => {
      //判断当前是否拥有id
      const currentRoute = router.currentRoute.value;
      
      if (currentRoute.params.id != "") {
          articleMode.value = "更新"
          getArticle(Number(currentRoute.params.id));
      }else{
        articleInfo.value.fid = 34
        articleInfo.value.type = 1
      }
    
      //获取分类列表
      getMenuList();      
    })

    const toolbarConfig = {}

    //自定义图片插入类型
    type InsertFnType = (url: string, alt: string, href: string) => void

    const editorConfig = { 
      placeholder: '请输入内容...',
      // 菜单配置
      MENU_CONF: {
        uploadImage: {
          server: (import.meta.env.VITE_BLOG_API_URL || '/blog') + "/admin/file/upload", // 上传图片地址
          fieldName: 'file',
          headers: {
            "x-token":`${JwtService.getToken()}`,
          },
          customInsert(res: any, insertFn: InsertFnType): void{
            if(res.code == 0){
              insertFn(publicApi+res.data.url,res.data.alt,"")
              ElNotification({
                duration:1500,
                title: 'success',
                message: "上传成功",
                type: 'success',
              })
            }else{
              ElNotification({
                duration:1500,
                title: 'error',
                message: "上传失败",
                type: 'error',
              })
            }
          },
          // 单个文件的最大体积限制，默认为 2M
          maxFileSize: 2 * 1024 * 1024, // 1M

          // 选择文件时的类型限制，默认为 ['image/*'] 。如不想限制，则设置为 []
          allowedFileTypes: ['image/*'],

          // 小于该值就插入 base64 格式（而不上传），默认为 0
          base64LimitSize: 5 * 1024 // 5kb
        
        },
      },
    }
    // 组件销毁时，也及时销毁编辑器
    onBeforeUnmount(() => {
      const editor = editorRef.value
      if (editor == null) return
      editor.destroy()
    })

    const handleCreated = (editor) => {
      editorRef.value = editor // 记录 editor 实例，重要！
    }


    const submitBtn = () => {
      let url = "/admin/article/updateArticle"
      let method = "put"
      let postData = {}
      postData = {
        ...articleInfo.value
      }
      postData = {
        ...articleInfo.value,
        pic:data.value.imgStr,
        content:valueHtml.value
      }
      if(articleMode.value == "新增"){
        url = "/admin/article/createArticle"
        method = "post"
      }
      service({
        url,
        method,
        data:postData
      }).then(({data}) => {
        if(data != undefined){
          Swal.fire({
              title:articleMode.value +"成功",
              icon:"success",
              showCancelButton:true,
            })
        }
      })
    }

    const data = ref({
      img: "/media/svg/files/blank-image.svg",
      imgStr: ""
      //初始化给一个默认图
    });

    //异步加载完成后需要重新渲染缩略图

    //移除图片方法
    const removeImage = () => {
      data.value.img = "/media/svg/files/blank-image.svg";
      data.value.imgStr = "";
    };

    //改变缩略图方法
    const changeImage = (e: any) => {
      const files = e.target.files;
      uploadImg(files[0])
    }

    const uploadImg = (img: any) => {
      let formdata = new FormData()
      formdata.append("file", img) //将每一个文件图片都加进formdata
      service.post("/admin/file/upload",formdata,{
        headers:{
          "Content-type":"multpart/form-data"
        }
      })
      .then(res => {
        data.value.img =   publicApi + res.data.url;
        data.value.imgStr = publicApi + res.data.url;
      })  
    }

    const getMenuList = () => {
      service.get("/base/getCategoryList")
        .then(res => {
          category.value = res.data
        })
    }

    const getArticle = (id:Number) => {
      service.post("/admin/article/getArticle",{id:id})
        .then(res => {
          articleInfo.value = res.data
          valueHtml.value = articleInfo.value.content
          data.value.img = articleInfo.value.pic
          data.value.imgStr = articleInfo.value.pic
        })
    }

    const getImg = (path:string) => {
      if(path.includes("/media/")){
        return getAssetPath(path);
      }
      return path;
    }

    //监视路由判断是否页面跳转
    watch(
      () => router.currentRoute.value,
      (newValue: any, oldValue: any) => {
        //更新文章页面直接进入添加文章旧数据会保留,刷新页面得以重置
        if(oldValue.name == "admin-article-content" && newValue.name == "admin-article-content"){
           if (Object.keys(oldValue.params).length !== 0) {
            //省略数据重置
            router.go(0);       
          }
        }        
      },
    )

    return {
      getAssetPath,
      editorRef,
      valueHtml,
      mode: 'default', // 或 'simple'
      toolbarConfig,
      editorConfig,
      handleCreated,
      data,
      removeImage,
      changeImage,
      category,
      articleInfo,
      submitBtn,
      articleMode,
      getImg,
    };
  }
})

</script>    

<style lang="scss" scoped>
.editor—wrapper {
  border: 1px solid #ccc;
  z-index: 100;
  /* 按需定义 */
}
// Base
.img-input {
  position: relative;
  background-repeat: no-repeat;
  background-size: cover;


  // Actions
  [data-kt-image-input-action] {
    cursor: pointer;
    position: absolute;
    transform: translate(-50%, -50%);
  }

  // Change Button
  [data-kt-image-input-action="change"] {
    left: 98%;
    top: 0;

    input {
      width: 0 !important;
      height: 0 !important;
      overflow: hidden;
      opacity: 0;
    }
  }

  // Cancel & Remove Buttons
  [data-kt-image-input-action="cancel"],
  [data-kt-image-input-action="remove"] {
    position: absolute;
    left: 98%;
    top: 100%;
  }

  // Bordered style
  &.img-input-outline {
    .img-input-wrapper {
      border: 3px solid var(--kt-body-bg);
      box-shadow: var(--kt-box-shadow);
    }
  }
}
</style>