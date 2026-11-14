<template>
  <div class="article-workspace">
    <div v-if="loadingEditor" class="card card-body mb-5" role="status">
      正在加载编辑器…
    </div>
    <div v-else-if="editorError" class="alert alert-danger" role="alert">
      文章加载失败，编辑已暂停，避免覆盖原文。<button
        class="btn btn-sm btn-light ms-3"
        @click="loadEditor"
      >
        重新加载
      </button>
    </div>
    <div
      v-if="!loadingEditor && saveConflict"
      class="alert alert-warning"
      role="alert"
    >
      服务器文章已在其他页面更新。本机草稿仍然保留，重新载入后可比较并选择是否恢复。<button
        type="button"
        class="btn btn-sm btn-light ms-3"
        @click="reloadLatest"
      >
        载入服务器最新版本
      </button>
    </div>
    <div
      v-else-if="!loadingEditor && pendingDraft"
      class="alert alert-warning d-flex flex-wrap align-items-center gap-3"
      role="status"
    >
      <span>发现未发布的本机草稿，恢复前请注意服务器文章可能已更新。</span>
      <button class="btn btn-sm btn-primary" @click="draft.restore">
        恢复草稿
      </button>
      <button class="btn btn-sm btn-light" @click="draft.discard">
        丢弃草稿
      </button>
    </div>
    <div
      v-if="!loadingEditor && !editorError"
      :inert="submitting"
      :aria-busy="submitting"
      class="d-flex flex-column flex-lg-row"
    >
      <!-- 正文开始 -->
      <div class="flex-lg-row-fluid me-lg-15 order-1 order-lg-1 mb-10 mb-lg-0">
        <div class="card card-flush py-5 mb-5 mb-lg-10">
          <!--begin::Card body-->
          <div class="card-body pt-0">
            <VForm
              id="kt_modal_create_api_key_form"
              class="form"
              @submit="submitBtn"
            >
              <!--begin::Input group-->
              <div class="mb-5 fv-row">
                <!--begin::Label-->
                <label
                  for="article-title"
                  class="required fs-5 fw-semobold mb-3"
                  >文章标题</label
                >
                <!--end::Label-->

                <!--begin::Input-->
                <Field
                  type="text"
                  class="form-control form-control-lg form-control-solid"
                  placeholder="请输入文章标题"
                  maxlength="200"
                  id="article-title"
                  name="title"
                  v-model="articleInfo.title"
                />
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
                      <label
                        for="article-category"
                        class="required fs-5 fw-semobold mb-3"
                        >文章分类</label
                      >

                      <Field
                        class="form-select form-select-lg form-select-solid select2-hidden-accessible"
                        placeholder="请选择文章分类"
                        id="article-category"
                        name="fid"
                        v-model="articleInfo.fid"
                        as="select"
                      >
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
                      <label
                        for="article-type"
                        class="required fs-5 fw-semobold mb-3"
                        >文章样式</label
                      >

                      <Field
                        as="select"
                        class="form-select form-select-lg form-select-solid"
                        placeholder="请选择文章样式"
                        id="article-type"
                        name="type"
                        v-model="articleInfo.type"
                      >
                        <option label="普通文章" :value="1"></option>
                        <option label="纯文本动态" :value="0"></option>
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

              <div
                class="d-flex align-items-center justify-content-between mb-3"
              >
                <label class="fs-5 fw-semobold mb-0">内容格式</label>
                <div class="btn-group" role="group" aria-label="内容格式">
                  <button
                    type="button"
                    :class="[
                      'btn btn-sm',
                      articleInfo.contentFormat === 'html'
                        ? 'btn-primary'
                        : 'btn-light-primary',
                    ]"
                    :disabled="uploading"
                    @click="articleInfo.contentFormat = 'html'"
                  >
                    富文本
                  </button>
                  <button
                    type="button"
                    :class="[
                      'btn btn-sm',
                      articleInfo.contentFormat === 'markdown'
                        ? 'btn-primary'
                        : 'btn-light-primary',
                    ]"
                    :disabled="uploading"
                    @click="articleInfo.contentFormat = 'markdown'"
                  >
                    Markdown
                  </button>
                </div>
              </div>

              <div
                v-if="articleInfo.contentFormat !== 'markdown'"
                class="editor—wrapper"
              >
                <Toolbar
                  v-if="editorRef"
                  :key="editorGeneration"
                  style="border-bottom: 1px solid #ccc"
                  :editor="editorRef"
                  :defaultConfig="toolbarConfig"
                  :mode="mode"
                />
                <Editor
                  :key="editorGeneration"
                  class="article-editor"
                  v-model="valueHtml"
                  :defaultConfig="editorConfig"
                  :mode="mode"
                  @onCreated="handleCreated"
                />
              </div>
              <MarkdownEditor
                v-else
                ref="markdownEditorRef"
                v-model="markdownValue"
                @uploading="markdownUploading = $event"
              />
            </VForm>
          </div>
          <!--end::Card body-->
        </div>
      </div>
      <!-- 正文结束 -->
      <!-- 侧边栏开始 -->
      <div
        class="flex-column flex-lg-row-auto w-100 w-lg-250px w-xl-300px mb-10 order-2 order-lg-2"
      >
        <div
          class="card card-flush mb-0"
          data-kt-sticky="true"
          data-kt-sticky-name="view-article-upload"
          data-kt-sticky-offset="{default: false, lg: '200px'}"
          data-kt-sticky-width="{lg: '250px', xl: '300px'}"
          data-kt-sticky-left="auto"
          data-kt-sticky-top="20px"
          data-kt-sticky-animation="false"
          data-kt-sticky-zindex="95"
        >
          <div class="card-body pt-0 fs-6">
            <section class="mt-7 mb-7" aria-label="发布设置">
              <label for="article-publish-state" class="form-label fw-bold"
                >发布状态</label
              >
              <select
                id="article-publish-state"
                v-model.number="publishStatus"
                class="form-select form-select-solid mb-3"
              >
                <option :value="2">草稿（不公开）</option>
                <option :value="0">公开发布</option>
                <option :value="3">定时发布</option>
                <option :value="1">隐藏（不公开）</option>
              </select>
              <template v-if="publishStatus === 3">
                <label for="article-publish-at" class="form-label"
                  >发布时间</label
                >
                <input
                  id="article-publish-at"
                  v-model="publishTime"
                  type="datetime-local"
                  class="form-control form-control-solid"
                />
                <p class="form-text">
                  时区：{{ timezone }}。服务器会在到期后 30 秒内检查发布。
                </p>
              </template>
              <p v-if="publishStatus === 2" class="text-muted fs-7 mb-0">
                保存到服务器后，可在其他设备继续编辑。
              </p>
            </section>
            <section class="mb-7">
              <label for="article-tags" class="form-label fw-bold"
                >文章标签</label
              >
              <el-select
                id="article-tags"
                v-model="tagNames"
                multiple
                filterable
                allow-create
                default-first-option
                :multiple-limit="10"
                placeholder="选择或输入标签，最多 10 个"
                class="w-100"
              >
                <el-option
                  v-for="tag in availableTags"
                  :key="tag.id"
                  :label="tag.name"
                  :value="tag.name"
                />
              </el-select>
            </section>
            <section class="mb-7">
              <label for="article-excerpt" class="form-label fw-bold"
                >列表摘要</label
              >
              <textarea
                id="article-excerpt"
                v-model="excerpt"
                maxlength="300"
                rows="3"
                class="form-control form-control-solid"
                placeholder="留空时自动从正文生成"
              ></textarea>
              <span class="text-muted fs-8">{{ excerpt.length }}/300</span>
            </section>
            <details class="mb-7">
              <summary class="fw-bold cursor-pointer mb-3">
                搜索引擎展示
              </summary>
              <label for="article-seo-title" class="form-label">SEO 标题</label>
              <input
                id="article-seo-title"
                v-model="seoTitle"
                maxlength="200"
                class="form-control form-control-solid mb-3"
                placeholder="留空时使用文章标题"
              />
              <label for="article-seo-description" class="form-label"
                >SEO 描述</label
              >
              <textarea
                id="article-seo-description"
                v-model="seoDescription"
                maxlength="300"
                rows="3"
                class="form-control form-control-solid"
                placeholder="留空时使用摘要"
              ></textarea>
              <div class="border rounded p-3 mt-3">
                <div class="text-primary fw-semibold text-break">
                  {{ seoTitle || articleInfo.title || "文章标题" }}
                </div>
                <div class="text-muted fs-7 mt-2 text-break">
                  {{ seoDescription || excerpt || "发布后将显示文章摘要。" }}
                </div>
              </div>
            </details>
            <div class="d-flex gap-2 mb-5">
              <button
                type="button"
                class="btn btn-sm btn-light-primary"
                @click="openPreview"
              >
                预览文章</button
              ><button
                v-if="articleInfo.ID"
                type="button"
                class="btn btn-sm btn-light"
                @click="openHistory"
              >
                历史版本
              </button>
            </div>
            <!-- 缩略图开始 -->
            <div v-if="Number(articleInfo.type) !== 0" class="mt-7 mb-5 row">
              <!--begin::Label-->
              <label class="fs-5 mb-7 fw-semobold">缩略图</label>
              <!--end::Label-->

              <!--begin::Image input-->
              <div
                class="img-input img-input-outline"
                data-kt-image-input="true"
              >
                <!--begin::Preview existing avatar-->
                <div
                  class="img-input-wrapper rounded bgi-no-repeat bgi-size-contain bgi-position-x-center min-h-150px p-8"
                  :style="`background-image: url(${getImg(data.img)})`"
                ></div>
                <!--end::Preview existing avatar-->

                <!--begin::Label-->
                <label
                  class="btn btn-icon btn-circle btn-active-color-primary w-25px h-25px bg-body shadow"
                  data-kt-image-input-action="change"
                  data-bs-toggle="tooltip"
                  title="编辑图片"
                  @change="changeImage"
                >
                  <i class="bi bi-pencil-fill fs-7"></i>

                  <!--begin::Inputs-->
                  <input
                    type="file"
                    name="avatar"
                    aria-label="上传缩略图"
                    :disabled="uploading"
                    accept=".png, .jpg, .jpeg"
                  />
                  <input type="hidden" name="avatar_remove" />
                  <!--end::Inputs-->
                </label>
                <!--end::Label-->

                <!--begin::Remove-->
                <span
                  class="btn btn-icon btn-circle btn-active-color-primary w-25px h-25px bg-body shadow"
                  data-kt-image-input-action="remove"
                  data-bs-toggle="tooltip"
                  @click="removeImage"
                  title="移除图片"
                >
                  <i class="bi bi-x fs-2"></i>
                </span>
                <!--end::Remove-->
              </div>
              <!--end::Image input-->

              <!--begin::Input group-->
              <div class="mt-5">
                <Field
                  type="text"
                  name="fname"
                  class="form-control form-control-lg form-control-solid mb-3 mb-lg-0"
                  placeholder="上传或者填写图片外链"
                  v-model="data.imgStr"
                />
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
                <span v-if="uploading">图片正在上传，完成后即可保存。</span>
                <span v-else-if="storageError" class="text-danger"
                  >本机草稿保存失败，请及时发布或备份内容。</span
                >
                <span v-else-if="savedAt"
                  >本机草稿已保存 {{ draftTime }}，尚未发布</span
                >
                <span v-else>{{
                  dirty
                    ? "正在保存本机草稿…"
                    : "本机草稿保留 7 天，仅此浏览器可见"
                }}</span>
              </p>
              <button
                type="button"
                :disabled="
                  submitting || uploading || !!pendingDraft || saveConflict
                "
                @click="submitBtn"
                class="d-none d-lg-block form-control from-control-lg from-control-solid bg-info mb-3 mb-lg-0"
                style="color: aliceblue"
              >
                {{
                  submitting
                    ? "正在保存…"
                    : uploading
                    ? "图片上传中…"
                    : articleMode
                }}
              </button>
            </div>
            <!--end::Input group-->
          </div>
        </div>
      </div>
      <!-- 侧边栏结束 -->
    </div>
  </div>
  <div
    v-if="!loadingEditor && !editorError"
    class="mobile-publish-bar d-lg-none"
  >
    <span class="text-muted fs-7">{{
      !articleInfo.ID
        ? "尚未保存到服务器"
        : dirty
        ? "有未保存的编辑"
        : "已与服务器同步"
    }}</span
    ><button
      type="button"
      class="btn btn-primary btn-sm"
      :disabled="submitting || uploading || !!pendingDraft || saveConflict"
      @click="submitBtn"
    >
      {{ submitting ? "正在保存…" : uploading ? "图片上传中…" : articleMode }}
    </button>
  </div>
  <el-dialog
    v-model="showPreview"
    title="文章预览"
    width="min(960px, 95vw)"
    destroy-on-close
    ><ArticlePreview :article="previewArticle"
  /></el-dialog>
  <el-dialog
    v-model="showHistory"
    title="文章历史版本"
    width="min(1000px, 95vw)"
    destroy-on-close
  >
    <p class="text-muted">
      最近 50 个保存版本。载入后先进入编辑器，需再次保存才会影响服务器文章。
    </p>
    <div v-if="historyLoading" role="status">正在读取历史版本…</div>
    <div v-else-if="historyError" class="alert alert-danger">
      历史版本加载失败。<button
        class="btn btn-sm btn-light ms-3"
        @click="openHistory"
      >
        重试
      </button>
    </div>
    <div v-else-if="!revisions.length" class="text-muted">
      暂无历史版本。下次保存后将开始记录。
    </div>
    <div v-else class="row g-5">
      <div class="col-md-4">
        <button
          v-for="revision in revisions"
          :key="revision.id"
          type="button"
          class="btn btn-sm w-100 text-start mb-2"
          :class="
            selectedRevision?.id === revision.id
              ? 'btn-light-primary'
              : 'btn-light'
          "
          @click="selectRevision(revision)"
        >
          <div class="fw-bold">{{ revision.title }}</div>
          <div class="fs-8 text-muted">
            版本 {{ revision.version }} ·
            {{ new Date(revision.savedAt * 1000).toLocaleString() }}
          </div>
        </button>
      </div>
      <div class="col-md-8">
        <template v-if="revisionSnapshot"
          ><button
            type="button"
            class="btn btn-sm btn-primary mb-4"
            @click="loadRevision"
          >
            载入此版本到编辑器</button
          ><ArticlePreview :article="revisionSnapshot"
        /></template>
        <p v-else class="text-muted">选择一个版本查看内容</p>
      </div>
    </div>
  </el-dialog>
</template>
<script setup lang="ts">
import "@wangeditor/editor/dist/css/style.css";
import {
  defineAsyncComponent,
  computed,
  onBeforeUnmount,
  onMounted,
  ref,
  shallowRef,
  watch,
} from "vue";
import {
  onBeforeRouteLeave,
  onBeforeRouteUpdate,
  useRoute,
  useRouter,
} from "vue-router";
import { Editor, Toolbar } from "@wangeditor/editor-for-vue";
import type { IDomEditor } from "@wangeditor/editor";
const MarkdownEditor = defineAsyncComponent(
  () => import("@/components/admin/MarkdownEditor.vue")
);
const ArticlePreview = defineAsyncComponent(
  () => import("@/components/admin/ArticlePreview.vue")
);
import type {
  ArticleRevision,
  ArticleSnapshot,
  Tag,
  ArticleStatus,
} from "@/core/blog/PublishingTypes";
import { localDateTime, unixDateTime } from "@/core/blog/PublishingTypes";
import { ErrorMessage, Field, Form as VForm } from "vee-validate";
import { ElNotification } from "element-plus";
import Swal from "sweetalert2";
import DOMPurify from "dompurify";
import { resolveMediaUrl } from "@/utils/media";
import { getAssetPath } from "@/core/helpers/assets";
import service from "@/utils/request";
import { useAuthStore } from "@/stores/auth";
import { useLocalDraft } from "@/composables/useLocalDraft";
import type { CategoryList } from "@/core/blog/CategoryTypes";
import type { ArticleList } from "@/core/blog/ArticleTypes";

interface ArticleDraft {
  title: string;
  fid: number;
  type: number;
  contentFormat: "html" | "markdown";
  html: string;
  markdown: string;
  pic: string;
  status?: ArticleStatus;
  publishAt?: number;
  tagNames?: string[];
  excerpt?: string;
  seoTitle?: string;
  seoDescription?: string;
}
const route = useRoute();
const router = useRouter();
const auth = useAuthStore();
const category = ref<CategoryList[]>([]);
const publishStatus = ref<ArticleStatus>(2);
const publishTime = ref("");
const excerpt = ref("");
const seoTitle = ref("");
const seoDescription = ref("");
const tagNames = ref<string[]>([]);
const availableTags = ref<Tag[]>([]);
const timezone = Intl.DateTimeFormat().resolvedOptions().timeZone;
const showPreview = ref(false);
const showHistory = ref(false);
const historyLoading = ref(false);
const historyError = ref(false);
const revisions = ref<ArticleRevision[]>([]);
const selectedRevision = ref<ArticleRevision>();
const revisionSnapshot = ref<ArticleSnapshot>();
const previewArticle = computed(() => ({
  title: articleInfo.value.title,
  content:
    articleInfo.value.contentFormat === "markdown"
      ? markdownValue.value
      : valueHtml.value,
  contentFormat: articleInfo.value.contentFormat,
  excerpt: excerpt.value,
  tagNames: tagNames.value,
  pic: data.value.imgStr,
}));
const articleInfo = ref<Partial<ArticleList>>({
  title: "",
  type: 1,
  contentFormat: "html",
});
const valueHtml = ref("");
const markdownValue = ref("");
const markdownEditorRef = ref<{ flush: () => void }>();
function flushEditor() {
  if (articleInfo.value.contentFormat === "markdown")
    markdownEditorRef.value?.flush();
}
function openPreview() {
  flushEditor();
  showPreview.value = true;
}
const data = ref({ img: "/media/svg/files/blank-image.svg", imgStr: "" });
const editorRef = shallowRef<IDomEditor>();
const editorGeneration = ref(0);
const mode = "default";
const toolbarConfig = window.matchMedia("(max-width: 991px)").matches
  ? {
      toolbarKeys: [
        "headerSelect",
        "bold",
        "italic",
        "bulletedList",
        "insertLink",
        "uploadImage",
        "undo",
        "redo",
        {
          key: "group-more",
          title: "更多格式",
          menuKeys: [
            "blockquote",
            "codeBlock",
            "insertTable",
            "numberedList",
            "fullScreen",
          ],
        },
      ],
    }
  : {};
const loadingEditor = ref(true);
const editorError = ref(false);
const submitting = ref(false);
const saveConflict = ref(false);
const coverUploading = ref(false);
const markdownUploading = ref(false);
const richUploads = ref(0);
const uploading = computed(
  () => coverUploading.value || markdownUploading.value || richUploads.value > 0
);
let loadSequence = 0;
const articleMode = computed(() =>
  publishStatus.value === 2
    ? "保存草稿"
    : publishStatus.value === 3
    ? "保存定时发布"
    : publishStatus.value === 1
    ? "保存为隐藏"
    : articleInfo.value.ID && articleInfo.value.status === 0
    ? "保存修改"
    : "发布文章"
);
const snapshot = (): ArticleDraft => ({
  title: articleInfo.value.title || "",
  fid: Number(articleInfo.value.fid || 0),
  type: Number(articleInfo.value.type ?? 1),
  status: publishStatus.value,
  publishAt: publishStatus.value === 3 ? unixDateTime(publishTime.value) : 0,
  tagNames: [...tagNames.value],
  excerpt: excerpt.value,
  seoTitle: seoTitle.value,
  seoDescription: seoDescription.value,
  contentFormat: articleInfo.value.contentFormat || "html",
  html: valueHtml.value === "<p><br></p>" ? "" : valueHtml.value,
  markdown: markdownValue.value,
  pic: data.value.imgStr,
});
const applyDraft = (draft: ArticleDraft) => {
  Object.assign(articleInfo.value, {
    title: draft.title,
    fid: draft.fid,
    type: draft.type,
    contentFormat: draft.contentFormat,
  });
  publishStatus.value =
    draft.status ?? (articleInfo.value.status as ArticleStatus) ?? 2;
  publishTime.value = localDateTime(draft.publishAt || 0);
  tagNames.value = draft.tagNames || articleInfo.value.tagNames || [];
  excerpt.value = draft.excerpt ?? articleInfo.value.excerpt ?? "";
  seoTitle.value = draft.seoTitle ?? articleInfo.value.seoTitle ?? "";
  seoDescription.value =
    draft.seoDescription ?? articleInfo.value.seoDescription ?? "";
  editorRef.value = undefined;
  editorGeneration.value++;
  valueHtml.value = DOMPurify.sanitize(draft.html);
  markdownValue.value = draft.markdown;
  data.value = {
    img: draft.pic || "/media/svg/files/blank-image.svg",
    imgStr: draft.pic,
  };
};
const validateDraft = (value: unknown): value is ArticleDraft => {
  const d = value as ArticleDraft;
  return (
    !!d &&
    typeof d.title === "string" &&
    Number.isInteger(d.fid) &&
    [0, 1].includes(d.type) &&
    ["html", "markdown"].includes(d.contentFormat) &&
    typeof d.html === "string" &&
    typeof d.markdown === "string" &&
    typeof d.pic === "string" &&
    (d.status === undefined || [0, 1, 2, 3].includes(d.status)) &&
    (d.publishAt === undefined || Number.isFinite(d.publishAt)) &&
    ["excerpt", "seoTitle", "seoDescription"].every(
      (key) => d[key] === undefined || typeof d[key] === "string"
    ) &&
    (d.tagNames === undefined ||
      (Array.isArray(d.tagNames) &&
        d.tagNames.length <= 10 &&
        d.tagNames.every((name) => typeof name === "string")))
  );
};
const draft = useLocalDraft(snapshot, applyDraft, validateDraft);
const { dirty, savedAt, storageError, pending: pendingDraft } = draft;
const draftTime = computed(() =>
  savedAt.value ? new Date(savedAt.value).toLocaleTimeString("zh-CN") : ""
);

async function loadEditor() {
  const sequence = ++loadSequence;
  draft.pause();
  loadingEditor.value = true;
  editorError.value = false;
  showHistory.value = false;
  showPreview.value = false;
  selectedRevision.value = undefined;
  revisionSnapshot.value = undefined;
  revisions.value = [];
  try {
    const id = route.params.id ? Number(route.params.id) : 0;
    if (!Number.isSafeInteger(id) || id < 0 || (!!route.params.id && id === 0))
      throw new Error("文章编号不合法");
    const [profile, categories, article, tags] = await Promise.all([
      auth.user?.user?.ID
        ? Promise.resolve(auth.user.user)
        : auth.refreshProfile(),
      service.get("/base/getCategoryList"),
      id
        ? service.post("/admin/article/getArticle", { id })
        : Promise.resolve({ data: undefined }),
      service.get("/admin/article/tags").catch(() => ({ data: [] })),
    ]);
    if (sequence !== loadSequence) return;
    if (!profile?.ID) throw new Error("无法获取账户资料");
    category.value = categories.data || [];
    const first =
      category.value.find((item) => item.fid !== 0) || category.value[0];
    articleInfo.value = article.data || {
      title: "",
      type: 1,
      contentFormat: "html",
      fid: first?.ID || 0,
      status: 2,
      version: 1,
    };
    availableTags.value = Array.isArray(tags.data) ? tags.data : [];
    publishStatus.value = (articleInfo.value.status ?? 0) as ArticleStatus;
    publishTime.value = localDateTime(
      articleInfo.value.publishAt || Math.floor(Date.now() / 1000) + 3600
    );
    tagNames.value =
      articleInfo.value.tagNames ||
      articleInfo.value.tags?.map((tag) => tag.name) ||
      [];
    excerpt.value = articleInfo.value.excerpt || "";
    seoTitle.value = articleInfo.value.seoTitle || "";
    seoDescription.value = articleInfo.value.seoDescription || "";
    articleInfo.value.contentFormat ||= "html";
    articleInfo.value.status ??= 0;
    articleInfo.value.version ??= 1;
    saveConflict.value = false;
    valueHtml.value =
      articleInfo.value.contentFormat === "html"
        ? DOMPurify.sanitize(articleInfo.value.content || "")
        : "";
    markdownValue.value =
      articleInfo.value.contentFormat === "markdown"
        ? articleInfo.value.content || ""
        : "";
    const pic = resolveMediaUrl(articleInfo.value.pic);
    data.value = {
      img: pic || "/media/svg/files/blank-image.svg",
      imgStr: pic,
    };
    draft.initialize(`${profile.ID}:${id || "new"}`);
  } catch {
    if (sequence === loadSequence) editorError.value = true;
  } finally {
    if (sequence === loadSequence) loadingEditor.value = false;
  }
}

async function mayLeave() {
  if (submitting.value || uploading.value) return false;
  flushEditor();
  draft.flush();
  if (!dirty.value || !auth.isAuthenticated) return true;
  const { isConfirmed } = await Swal.fire({
    title: "离开编辑页面？",
    text: storageError.value
      ? "本机草稿未能保存，离开可能丢失修改。"
      : "修改尚未发布，已保存在本机草稿中。",
    icon: "warning",
    showCancelButton: true,
    confirmButtonText: "离开页面",
    cancelButtonText: "继续编辑",
  });
  return isConfirmed;
}
onBeforeRouteLeave(mayLeave);
onBeforeRouteUpdate(
  (to, from) => to.params.id === from.params.id || mayLeave()
);
const beforeUnload = (event: BeforeUnloadEvent) => {
  flushEditor();
  if (!dirty.value && !submitting.value && !uploading.value) return;
  draft.flush();
  event.preventDefault();
  event.returnValue = "";
};
onMounted(() => {
  window.addEventListener("beforeunload", beforeUnload);
  void loadEditor();
});
watch(
  () => route.params.id,
  () => {
    void loadEditor();
  }
);
onBeforeUnmount(() => {
  ++loadSequence;
  window.removeEventListener("beforeunload", beforeUnload);
  editorRef.value?.destroy();
});
const handleCreated = (editor: IDomEditor) => {
  editorRef.value = editor;
};
const getImg = (path: string) =>
  path.includes("/media/") ? getAssetPath(path) : resolveMediaUrl(path);
const removeImage = () => {
  data.value = { img: "/media/svg/files/blank-image.svg", imgStr: "" };
};
async function changeImage(event: Event) {
  const input = event.target as HTMLInputElement;
  const file = input.files?.[0];
  if (!file || uploading.value) return;
  coverUploading.value = true;
  try {
    const form = new FormData();
    form.append("file", file);
    const response = await service.post("/admin/file/upload", form);
    const pic = resolveMediaUrl(response.data.url);
    data.value = { img: pic, imgStr: pic };
  } catch {
    /* 请求拦截器已显示上传失败原因。 */
  } finally {
    coverUploading.value = false;
    input.value = "";
  }
}
const editorConfig = {
  placeholder: "请输入内容...",
  MENU_CONF: {
    uploadImage: {
      maxFileSize: 10 * 1024 * 1024,
      allowedFileTypes: ["image/*"],
      base64LimitSize: 0,
      async customUpload(
        file: File,
        insert: (url: string, alt: string, href: string) => void
      ) {
        if (file.size > 10 * 1024 * 1024) {
          ElNotification({
            title: "图片过大",
            message: "单张图片不能超过 10 MB",
            type: "warning",
          });
          return;
        }
        richUploads.value++;
        try {
          const form = new FormData();
          form.append("file", file);
          const response = await service.post("/admin/file/upload", form);
          insert(
            resolveMediaUrl(response.data.url),
            response.data.alt || "",
            ""
          );
        } catch {
          // Keep editing after the request interceptor reports the upload failure.
        } finally {
          richUploads.value--;
        }
      },
    },
  },
};
async function reloadLatest() {
  draft.flush();
  if (storageError.value) {
    const answer = await Swal.fire({
      title: "本机草稿尚未保存",
      text: "重新载入会替换当前编辑器内容，请先复制需要保留的内容。",
      icon: "warning",
      showCancelButton: true,
      confirmButtonText: "仍然载入",
      cancelButtonText: "继续编辑",
    });
    if (!answer.isConfirmed) return;
  }
  await loadEditor();
}
async function openHistory() {
  const articleID = articleInfo.value.ID;
  const sequence = loadSequence;
  showHistory.value = true;
  historyLoading.value = true;
  historyError.value = false;
  selectedRevision.value = undefined;
  revisionSnapshot.value = undefined;
  try {
    const response = await service.get("/admin/article/revisions", {
      params: { id: articleID },
    });
    if (sequence !== loadSequence) return;
    revisions.value = Array.isArray(response.data) ? response.data : [];
  } catch {
    if (sequence === loadSequence) historyError.value = true;
  } finally {
    if (sequence === loadSequence) historyLoading.value = false;
  }
}
async function selectRevision(revision: ArticleRevision) {
  const sequence = loadSequence;
  selectedRevision.value = revision;
  revisionSnapshot.value = undefined;
  try {
    const response = await service.get("/admin/article/revision", {
      params: { id: articleInfo.value.ID, revisionId: revision.id },
    });
    if (sequence === loadSequence && selectedRevision.value?.id === revision.id)
      revisionSnapshot.value = response.data;
  } catch {
    ElNotification({
      title: "读取版本失败",
      message: "请重新选择或稍后再试",
      type: "error",
    });
  }
}
async function loadRevision() {
  const saved = revisionSnapshot.value;
  if (!saved) return;
  if (dirty.value) {
    const answer = await Swal.fire({
      title: "替换当前未保存的编辑？",
      text: "载入历史版本将替换编辑器内容，服务器文章暂不改变。",
      showCancelButton: true,
      confirmButtonText: "载入版本",
      cancelButtonText: "继续编辑",
    });
    if (!answer.isConfirmed) return;
  }
  applyDraft({
    ...snapshot(),
    title: saved.title,
    fid: saved.fid,
    type: saved.type,
    contentFormat: saved.contentFormat,
    html: saved.contentFormat === "html" ? saved.content : "",
    markdown: saved.contentFormat === "markdown" ? saved.content : "",
    pic: saved.pic,
    excerpt: saved.excerpt,
    seoTitle: saved.seoTitle,
    seoDescription: saved.seoDescription,
    tagNames: saved.tagNames,
  });
  showHistory.value = false;
  draft.flush();
  ElNotification({
    title: "已载入历史内容",
    message: "内容尚未保存或发布，请确认后再保存",
    type: "success",
  });
}
async function submitBtn() {
  if (
    submitting.value ||
    uploading.value ||
    loadingEditor.value ||
    editorError.value ||
    pendingDraft.value ||
    saveConflict.value
  )
    return;
  flushEditor();
  const current = snapshot();
  if (current.status !== 2 && (!current.title.trim() || !current.fid)) {
    ElNotification({
      title: "请完善文章",
      message: "请填写标题并选择分类",
      type: "warning",
    });
    return;
  }
  if (
    current.status === 3 &&
    (!current.publishAt || current.publishAt <= Math.floor(Date.now() / 1000))
  ) {
    ElNotification({ title: "请选择未来的发布时间", type: "warning" });
    return;
  }
  if (
    articleInfo.value.ID &&
    articleInfo.value.status === 0 &&
    current.status !== 0
  ) {
    const answer = await Swal.fire({
      title: "暂时撤下已发布文章？",
      text: "保存后文章将不再公开，草稿、隐藏或定时发布状态均不会向访客展示。",
      icon: "warning",
      showCancelButton: true,
      confirmButtonText: "确认调整",
      cancelButtonText: "继续编辑",
    });
    if (!answer.isConfirmed) return;
  }
  submitting.value = true;
  try {
    const id = articleInfo.value.ID;
    const response = await service({
      url: id ? "/admin/article/updateArticle" : "/admin/article/createArticle",
      method: id ? "put" : "post",
      data: {
        ...articleInfo.value,
        title: current.title.trim(),
        fid: current.fid,
        type: current.type,
        version: articleInfo.value.version || 1,
        status: current.status,
        publishAt: current.publishAt,
        tagNames: current.tagNames,
        excerpt: current.excerpt,
        seoTitle: current.seoTitle,
        seoDescription: current.seoDescription,
        content:
          current.contentFormat === "markdown"
            ? current.markdown
            : current.html,
        contentFormat: current.contentFormat,
        pic: current.type === 0 ? "" : resolveMediaUrl(current.pic),
      },
    });
    articleInfo.value.version =
      response.data?.version || (articleInfo.value.version || 1) + 1;
    articleInfo.value.status = current.status ?? 2;
    if (!id && response.data?.id) articleInfo.value.ID = response.data.id;
    draft.markSaved(current);
    ElNotification({
      title: current.status === 0 ? "发布内容已保存" : "保存成功",
      message:
        current.status === 2
          ? "草稿已保存到服务器，访客不可见"
          : current.status === 3
          ? "定时发布已设置"
          : "文章已保存到服务器",
      type: "success",
    });
    submitting.value = false;
    if (!id && response.data?.id)
      await router.replace({
        name: "admin-article-content",
        params: { id: response.data.id },
      });
    else if (!id) await router.replace({ name: "admin-article-index" });
  } catch (error: any) {
    if (error?.response?.status === 409) saveConflict.value = true;
    draft.flush();
  } finally {
    submitting.value = false;
  }
}
</script>

<style lang="scss" scoped>
.article-editor {
  height: 500px;
  overflow-y: hidden;
}
.mobile-publish-bar {
  position: fixed;
  left: 0;
  right: 0;
  bottom: 0;
  z-index: 105;
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 0.75rem 1rem calc(0.75rem + env(safe-area-inset-bottom));
  background: var(--bs-body-bg);
  border-top: 1px solid var(--bs-border-color);
  box-shadow: 0 -2px 12px #00000008;
}
@media (max-width: 991px) {
  .article-workspace {
    padding-bottom: 90px;
  }
  .article-editor {
    height: 380px;
  }
}

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
