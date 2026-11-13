<template>
  <div>
  <div v-if="loadingEditor" class="card card-body mb-5" role="status">正在加载编辑器…</div>
  <div v-else-if="editorError" class="alert alert-danger" role="alert">文章加载失败，编辑已暂停，避免覆盖原文。<button class="btn btn-sm btn-light ms-3" @click="loadEditor">重新加载</button></div>
  <div v-else-if="pendingDraft" class="alert alert-warning d-flex flex-wrap align-items-center gap-3" role="status">
    <span>发现未发布的本机草稿，恢复前请注意服务器文章可能已更新。</span>
    <button class="btn btn-sm btn-primary" @click="draft.restore">恢复草稿</button>
    <button class="btn btn-sm btn-light" @click="draft.discard">丢弃草稿</button>
  </div>
  <div v-if="!loadingEditor && !editorError" :inert="submitting" :aria-busy="submitting" class="d-flex flex-column flex-lg-row">

    <!-- 正文开始 -->
    <div class="flex-lg-row-fluid me-lg-15 order-2 order-lg-1 mb-10 mb-lg-0">

      <div class="card card-flush py-5 mb-5 mb-lg-10">
        <!--begin::Card body-->
        <div class="card-body pt-0">
          <VForm id="kt_modal_create_api_key_form" class="form" @submit="submitBtn">

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
                          :value="1"
                        ></option>
                        <option
                          label="纯文本动态"
                          :value="0"
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


            <div class="d-flex align-items-center justify-content-between mb-3">
              <label class="fs-5 fw-semobold mb-0">内容格式</label>
              <div class="btn-group" role="group" aria-label="内容格式">
                <button
                  type="button"
                  :class="['btn btn-sm', articleInfo.contentFormat === 'html' ? 'btn-primary' : 'btn-light-primary']"
                  @click="articleInfo.contentFormat = 'html'"
                >
                  富文本
                </button>
                <button
                  type="button"
                  :class="['btn btn-sm', articleInfo.contentFormat === 'markdown' ? 'btn-primary' : 'btn-light-primary']"
                  @click="articleInfo.contentFormat = 'markdown'"
                >
                  Markdown
                </button>
              </div>
            </div>

            <div v-if="articleInfo.contentFormat !== 'markdown'" class="editor—wrapper">
              <Toolbar style="border-bottom: 1px solid #ccc" :editor="editorRef" :defaultConfig="toolbarConfig"
                :mode="mode" />
              <Editor style="height: 500px; overflow-y: hidden;" v-model="valueHtml" :defaultConfig="editorConfig"
                :mode="mode" @onCreated="handleCreated" />
            </div>
            <MarkdownEditor v-else v-model="markdownValue" />
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
          <div v-if="Number(articleInfo.type) !== 0" class="mt-7 mb-5 row">
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
            <p class="text-muted fs-7" role="status" aria-live="polite">
              <span v-if="storageError" class="text-danger">本机草稿保存失败，请及时发布或备份内容。</span>
              <span v-else-if="savedAt">本机草稿已保存 {{ draftTime }}，尚未发布</span>
              <span v-else>{{ dirty ? '正在保存本机草稿…' : '本机草稿保留 7 天，仅此浏览器可见' }}</span>
            </p>
            <button type="button" :disabled="submitting || uploading || !!pendingDraft" @click="submitBtn" class="form-control from-control-lg from-control-solid bg-info mb-3 mb-lg-0" style="color:aliceblue;">
              {{ submitting ? "正在保存…" : articleMode }}
            </button>
          </div>
          <!--end::Input group-->
        </div>



      </div>
    </div>
    <!-- 侧边栏结束 -->
  </div>
  </div>
</template>
<script setup lang="ts">
import '@wangeditor/editor/dist/css/style.css';
import { defineAsyncComponent, computed, onBeforeUnmount, onMounted, ref, shallowRef, watch } from 'vue';
import { onBeforeRouteLeave, onBeforeRouteUpdate, useRoute, useRouter } from 'vue-router';
import { Editor, Toolbar } from '@wangeditor/editor-for-vue';
import type { IDomEditor } from '@wangeditor/editor';
const MarkdownEditor = defineAsyncComponent(() => import('@/components/admin/MarkdownEditor.vue'));
import { ErrorMessage, Field, Form as VForm } from 'vee-validate';
import { ElNotification } from 'element-plus';
import Swal from 'sweetalert2';
import { resolveMediaUrl } from '@/utils/media';
import { getAssetPath } from '@/core/helpers/assets';
import JwtService from '@/core/services/JwtService';
import service from '@/utils/request';
import { useAuthStore } from '@/stores/auth';
import { useLocalDraft } from '@/composables/useLocalDraft';
import type { CategoryList } from '@/core/blog/CategoryTypes';
import type { ArticleList } from '@/core/blog/ArticleTypes';

interface ArticleDraft {
  title: string; fid: number; type: number; contentFormat: 'html' | 'markdown';
  html: string; markdown: string; pic: string;
}
const route = useRoute();
const router = useRouter();
const auth = useAuthStore();
const category = ref<CategoryList[]>([]);
const articleInfo = ref<Partial<ArticleList>>({ title: '', type: 1, contentFormat: 'html' });
const valueHtml = ref('');
const markdownValue = ref('');
const data = ref({ img: '/media/svg/files/blank-image.svg', imgStr: '' });
const editorRef = shallowRef<IDomEditor>();
const mode = 'default';
const toolbarConfig = {};
const loadingEditor = ref(true);
const editorError = ref(false);
const submitting = ref(false);
const uploading = ref(false);
let loadSequence = 0;
const articleMode = computed(() => articleInfo.value.ID ? '保存修改' : '发布文章');
const snapshot = (): ArticleDraft => ({
  title: articleInfo.value.title || '', fid: Number(articleInfo.value.fid || 0), type: Number(articleInfo.value.type ?? 1),
  contentFormat: articleInfo.value.contentFormat || 'html', html: valueHtml.value === '<p><br></p>' ? '' : valueHtml.value, markdown: markdownValue.value, pic: data.value.imgStr,
});
const applyDraft = (draft: ArticleDraft) => {
  Object.assign(articleInfo.value, { title: draft.title, fid: draft.fid, type: draft.type, contentFormat: draft.contentFormat });
  valueHtml.value = draft.html;
  markdownValue.value = draft.markdown;
  data.value = { img: draft.pic || '/media/svg/files/blank-image.svg', imgStr: draft.pic };
};
const validateDraft = (value: unknown): value is ArticleDraft => {
  const d = value as ArticleDraft;
  return !!d && typeof d.title === 'string' && Number.isInteger(d.fid) && [0, 1].includes(d.type)
    && ['html', 'markdown'].includes(d.contentFormat) && typeof d.html === 'string' && typeof d.markdown === 'string' && typeof d.pic === 'string';
};
const draft = useLocalDraft(snapshot, applyDraft, validateDraft);
const { dirty, savedAt, storageError, pending: pendingDraft } = draft;
const draftTime = computed(() => savedAt.value ? new Date(savedAt.value).toLocaleTimeString('zh-CN') : '');

async function loadEditor() {
  const sequence = ++loadSequence;
  draft.pause();
  loadingEditor.value = true;
  editorError.value = false;
  try {
    const id = route.params.id ? Number(route.params.id) : 0;
    if (!Number.isInteger(id) || id < 0) throw new Error('文章编号不合法');
    const [profile, categories, article] = await Promise.all([
      auth.user?.user?.ID ? Promise.resolve(auth.user.user) : auth.refreshProfile(),
      service.get('/base/getCategoryList'),
      id ? service.post('/admin/article/getArticle', { id }) : Promise.resolve({ data: undefined }),
    ]);
    if (sequence !== loadSequence) return;
    if (!profile?.ID) throw new Error('无法获取账户资料');
    category.value = categories.data || [];
    const first = category.value.find(item => item.fid !== 0) || category.value[0];
    articleInfo.value = article.data || { title: '', type: 1, contentFormat: 'html', fid: first?.ID || 0 };
    articleInfo.value.contentFormat ||= 'html';
    valueHtml.value = articleInfo.value.contentFormat === 'html' ? articleInfo.value.content || '' : '';
    markdownValue.value = articleInfo.value.contentFormat === 'markdown' ? articleInfo.value.content || '' : '';
    const pic = resolveMediaUrl(articleInfo.value.pic);
    data.value = { img: pic || '/media/svg/files/blank-image.svg', imgStr: pic };
    draft.initialize(`${profile.ID}:${id || 'new'}`);
  } catch { if (sequence === loadSequence) editorError.value = true; }
  finally { if (sequence === loadSequence) loadingEditor.value = false; }
}

async function mayLeave() {
  if (submitting.value || uploading.value) return false;
  draft.flush();
  if (!dirty.value) return true;
  const { isConfirmed } = await Swal.fire({ title: '离开编辑页面？', text: storageError.value ? '本机草稿未能保存，离开可能丢失修改。' : '修改尚未发布，已保存在本机草稿中。', icon: 'warning', showCancelButton: true, confirmButtonText: '离开页面', cancelButtonText: '继续编辑' });
  return isConfirmed;
}
onBeforeRouteLeave(mayLeave);
onBeforeRouteUpdate((to, from) => to.params.id === from.params.id || mayLeave());
const beforeUnload = (event: BeforeUnloadEvent) => {
  if (!dirty.value && !submitting.value && !uploading.value) return;
  draft.flush(); event.preventDefault(); event.returnValue = '';
};
onMounted(() => { window.addEventListener('beforeunload', beforeUnload); void loadEditor(); });
watch(() => route.params.id, () => { void loadEditor(); });
onBeforeUnmount(() => { ++loadSequence; window.removeEventListener('beforeunload', beforeUnload); editorRef.value?.destroy(); });
const handleCreated = (editor: IDomEditor) => { editorRef.value = editor; };
const getImg = (path: string) => path.includes('/media/') ? getAssetPath(path) : resolveMediaUrl(path);
const removeImage = () => { data.value = { img: '/media/svg/files/blank-image.svg', imgStr: '' }; };
async function changeImage(event: Event) {
  const input = event.target as HTMLInputElement;
  const file = input.files?.[0];
  if (!file || uploading.value) return;
  uploading.value = true;
  try {
    const form = new FormData(); form.append('file', file);
    const response = await service.post('/admin/file/upload', form);
    const pic = resolveMediaUrl(response.data.url);
    data.value = { img: pic, imgStr: pic };
  } catch { /* 请求拦截器已显示上传失败原因。 */ }
  finally { uploading.value = false; input.value = ''; }
}
const editorConfig = {
  placeholder: '请输入内容...',
  MENU_CONF: { uploadImage: {
    server: (import.meta.env.VITE_BLOG_API_URL || '/blog') + '/admin/file/upload', fieldName: 'file',
    headers: { 'x-token': JwtService.getToken() || '' }, maxFileSize: 10 * 1024 * 1024,
    allowedFileTypes: ['image/*'], base64LimitSize: 0,
    customInsert(response: any, insert: (url: string, alt: string, href: string) => void) {
      if (response.code === 0) insert(resolveMediaUrl(response.data.url), response.data.alt || '', '');
      else ElNotification({ title: '上传失败', message: response.msg || '请重试', type: 'error' });
    },
  } },
};
async function submitBtn() {
  if (submitting.value || uploading.value || loadingEditor.value || editorError.value || pendingDraft.value) return;
  const current = snapshot();
  if (!current.title.trim() || !current.fid) {
    ElNotification({ title: '请完善文章', message: '请填写标题并选择分类', type: 'warning' }); return;
  }
  submitting.value = true;
  try {
    const id = articleInfo.value.ID;
    const response = await service({ url: id ? '/admin/article/updateArticle' : '/admin/article/createArticle', method: id ? 'put' : 'post', data: {
      ...articleInfo.value, title: current.title.trim(), fid: current.fid, type: current.type,
      content: current.contentFormat === 'markdown' ? current.markdown : current.html,
      contentFormat: current.contentFormat, pic: current.type === 0 ? '' : resolveMediaUrl(current.pic),
    } });
    draft.markSaved(current);
    ElNotification({ title: id ? '保存成功' : '发布成功', message: '文章已保存到服务器', type: 'success' });
    submitting.value = false;
    if (!id && response.data?.id) await router.replace({ name: 'admin-article-content', params: { id: response.data.id } });
    else if (!id) await router.replace({ name: 'admin-article-index' });
  } catch { draft.flush(); }
  finally { submitting.value = false; }
}
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
