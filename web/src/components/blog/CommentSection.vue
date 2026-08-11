<template>
  <div class="comment-section">
    <!-- 评论标题 -->
    <div class="d-flex align-items-center mb-8">
      <h3 class="fw-bold text-dark mb-0">
        <i class="bi bi-chat-left-text fs-2 me-2"></i>
        评论 ({{ comments.length }})
      </h3>
    </div>

    <!-- 评论列表 -->
    <div v-if="loading" class="mb-10 placeholder-glow" aria-label="正在加载评论">
      <div v-for="item in 2" :key="item" class="card mb-5 shadow-sm">
        <div class="card-body p-6">
          <div class="d-flex align-items-center mb-4">
            <span class="placeholder rounded-circle comment-avatar me-4"></span>
            <div class="flex-grow-1">
              <span class="placeholder col-3 d-block mb-3"></span>
              <span class="placeholder col-2 d-block"></span>
            </div>
          </div>
          <span class="placeholder col-10 d-block mb-3"></span>
          <span class="placeholder col-7 d-block"></span>
        </div>
      </div>
    </div>

    <div v-else-if="comments.length > 0" class="mb-10 comment-results">
      <div
        v-for="comment in comments"
        :key="comment.ID"
        class="card mb-5 shadow-sm hover-elevate-up"
      >
        <div class="card-body p-6">
          <!-- 评论者信息 -->
          <div class="d-flex align-items-center mb-4">
            <div class="symbol symbol-45px me-4">
              <div
                class="symbol-label fs-3 fw-bold"
                :style="`background-color: ${getAvatarColor(comment.name)}`"
              >
                {{ comment.name.charAt(0).toUpperCase() }}
              </div>
            </div>
            <div class="flex-grow-1">
              <div class="fw-bold text-gray-800 fs-5">{{ comment.name }}</div>
              <div class="text-muted fs-7">
                {{ formatDate(comment.ctime) }}
              </div>
            </div>
          </div>

          <!-- 评论内容 -->
          <div class="text-gray-700 fs-6" style="white-space: pre-wrap">
            {{ comment.content }}
          </div>
        </div>
      </div>
    </div>

    <!-- 无评论提示 -->
    <div
      v-else
      class="text-center py-10 mb-10 bg-light rounded"
    >
      <i class="bi bi-chat-dots fs-2x text-muted mb-3"></i>
      <p class="text-muted fs-5 mb-0">暂无评论，快来抢沙发吧！</p>
    </div>

    <!-- 评论表单 -->
    <div class="card shadow-sm">
      <div class="card-body p-8">
        <h4 class="fw-bold text-dark mb-6">
          <i class="bi bi-pencil-square fs-3 me-2"></i>
          发表评论
        </h4>

        <form @submit.prevent="submitComment">
          <!-- 昵称和邮箱 -->
          <div class="row mb-6">
            <div class="col-md-6 mb-5 mb-md-0">
              <label class="form-label required fw-semobold">昵称</label>
              <input
                v-model="formData.name"
                type="text"
                class="form-control form-control-lg"
                placeholder="请输入您的昵称"
                required
                maxlength="50"
              />
            </div>
            <div class="col-md-6">
              <label class="form-label required fw-semobold">邮箱</label>
              <input
                v-model="formData.email"
                type="email"
                class="form-control form-control-lg"
                placeholder="请输入您的邮箱"
                required
                maxlength="100"
              />
            </div>
          </div>

          <!-- 评论内容 -->
          <div class="mb-6">
            <label class="form-label required fw-semobold">评论内容</label>
            <textarea
              v-model="formData.content"
              class="form-control form-control-lg"
              rows="5"
              placeholder="说点什么吧..."
              required
              maxlength="1000"
            ></textarea>
            <div class="text-muted fs-7 mt-2">
              已输入 {{ formData.content.length }} / 1000 字
            </div>
          </div>

          <!-- 提交按钮 -->
          <div class="d-flex justify-content-end">
            <button
              type="submit"
              class="btn btn-primary btn-lg"
              :disabled="submitting"
            >
              <span v-if="!submitting">
                <i class="bi bi-send fs-4 me-2"></i>
                发表评论
              </span>
              <span v-else>
                <span class="spinner-border spinner-border-sm me-2"></span>
                提交中...
              </span>
            </button>
          </div>
        </form>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from "vue";
import service from "@/utils/request";
import type { Comment, CommentFormData } from "@/core/blog/CommentTypes";
import { ElMessage } from "element-plus";

// Props
const props = defineProps<{
  articleId: number;
}>();

// 响应式数据
const comments = ref<Comment[]>([]);
const loading = ref(true);
const submitting = ref(false);
const formData = ref<CommentFormData>({
  name: "",
  email: "",
  content: "",
  aid: props.articleId,
});

// 获取评论列表
const getComments = async () => {
  loading.value = true;
  try {
    const res: any = await service.get("/base/getCommentList", {
      params: {
        aid: props.articleId,
      },
    });
    if (res && res.data && Array.isArray(res.data)) {
      comments.value = res.data;
    }
  } catch (error) {
    console.error("获取评论失败:", error);
  } finally {
    loading.value = false;
  }
};

// 提交评论
const submitComment = async () => {
  if (submitting.value) return;

  // 验证表单
  if (!formData.value.name.trim()) {
    ElMessage.warning("请输入昵称");
    return;
  }
  if (!formData.value.email.trim()) {
    ElMessage.warning("请输入邮箱");
    return;
  }
  if (!formData.value.content.trim()) {
    ElMessage.warning("请输入评论内容");
    return;
  }

  submitting.value = true;

  try {
    const res: any = await service.post("/base/createComment", {
      ...formData.value,
      aid: props.articleId,
      uid: 0,
    });

    if (res.code === 0) {
      ElMessage.success("评论成功！");
      // 清空表单
      formData.value.content = "";
      // 重新加载评论列表
      await getComments();
      // 滚动到评论区域
      window.scrollTo({
        top: document.querySelector(".comment-section")?.scrollTop || 0,
        behavior: "smooth",
      });
    } else {
      ElMessage.error(res.msg || "评论失败");
    }
  } catch (error: any) {
    console.error("提交评论失败:", error);
    ElMessage.error("评论失败，请稍后重试");
  } finally {
    submitting.value = false;
  }
};

// 格式化日期
const formatDate = (timestamp: number): string => {
  const date = new Date(timestamp * 1000);
  const now = new Date();
  const diff = Math.floor((now.getTime() - date.getTime()) / 1000);

  if (diff < 60) return "刚刚";
  if (diff < 3600) return `${Math.floor(diff / 60)}分钟前`;
  if (diff < 86400) return `${Math.floor(diff / 3600)}小时前`;
  if (diff < 604800) return `${Math.floor(diff / 86400)}天前`;

  const year = date.getFullYear();
  const month = String(date.getMonth() + 1).padStart(2, "0");
  const day = String(date.getDate()).padStart(2, "0");
  const hours = String(date.getHours()).padStart(2, "0");
  const minutes = String(date.getMinutes()).padStart(2, "0");

  if (year === now.getFullYear()) {
    return `${month}-${day} ${hours}:${minutes}`;
  }
  return `${year}-${month}-${day} ${hours}:${minutes}`;
};

// 获取头像颜色
const getAvatarColor = (name: string): string => {
  const colors = [
    "#667eea",
    "#764ba2",
    "#f093fb",
    "#f5576c",
    "#4facfe",
    "#00f2fe",
    "#43e97b",
    "#38f9d7",
    "#fa709a",
    "#fee140",
  ];
  const index = name.charCodeAt(0) % colors.length;
  return colors[index];
};

// 页面加载时获取评论
onMounted(() => {
  getComments();
});
</script>

<style scoped>
.hover-elevate-up {
  transition: all 0.3s ease;
}

.placeholder { background-color: var(--kt-gray-300); }
.comment-avatar { width: 45px; height: 45px; flex: 0 0 45px; }
.comment-results { animation: comment-enter .24s ease-out; }
@keyframes comment-enter {
  from { opacity: 0; transform: translateY(4px); }
  to { opacity: 1; transform: translateY(0); }
}

@media (prefers-reduced-motion: reduce) {
  .comment-results { animation: none; }
}

.hover-elevate-up:hover {
  transform: translateY(-2px);
  box-shadow: 0 0.5rem 1.5rem 0.5rem rgba(0, 0, 0, 0.075) !important;
}

.form-label.required::after {
  content: " *";
  color: #f1416c;
}
</style>

