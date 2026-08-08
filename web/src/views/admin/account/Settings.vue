<template>
  <div class="account-settings-page">
    <div class="card mb-5 mb-xl-10">
      <div
        class="card-header border-0 cursor-pointer"
        role="button"
        data-bs-toggle="collapse"
        data-bs-target="#kt_account_profile_details"
        aria-expanded="true"
        aria-controls="kt_account_profile_details"
      >
        <div class="card-title m-0">
          <h3 class="fw-bold m-0">个人资料</h3>
        </div>
      </div>

      <div id="kt_account_profile_details" class="collapse show">
        <form id="kt_account_profile_details_form" class="form" novalidate @submit.prevent="saveProfile">
          <div class="card-body border-top p-9">
            <div class="row mb-6">
              <label class="col-lg-4 col-form-label fw-semibold fs-6">头像</label>
              <div class="col-lg-8">
                <div
                  class="image-input image-input-outline"
                  data-kt-image-input="true"
                  :style="{ backgroundImage: `url(${defaultAvatar})` }"
                >
                  <div
                    class="image-input-wrapper w-125px h-125px"
                    :style="{ backgroundImage: `url('${profileAvatar}'), url('${defaultAvatar}')` }"
                  ></div>
                  <span v-if="uploadingAvatar" class="avatar-loading">
                    <span class="spinner-border spinner-border-sm text-primary"></span>
                  </span>
                  <label
                    class="btn btn-icon btn-circle btn-active-color-primary w-25px h-25px bg-body shadow"
                    :class="{ disabled: uploadingAvatar }"
                    data-kt-image-input-action="change"
                    title="修改头像"
                  >
                    <i class="bi bi-pencil-fill fs-7"></i>
                    <input type="file" name="avatar" accept=".png,.jpg,.jpeg,.webp,.gif" :disabled="uploadingAvatar" @change="uploadAvatar" />
                    <input type="hidden" name="avatar_remove" />
                  </label>
                  <button
                    type="button"
                    class="btn btn-icon btn-circle btn-active-color-primary w-25px h-25px bg-body shadow"
                    data-kt-image-input-action="remove"
                    title="移除头像"
                    :disabled="uploadingAvatar || !profile.pic"
                    @click="profile.pic = ''"
                  >
                    <i class="bi bi-x fs-2"></i>
                  </button>
                </div>
                <div class="form-text">支持 PNG、JPG、WebP 或 GIF 格式。</div>
              </div>
            </div>

            <div class="row mb-6">
              <label class="col-lg-4 col-form-label fw-semibold fs-6">登录用户名</label>
              <div class="col-lg-8 fv-row">
                <input :value="profile.userName" type="text" class="form-control form-control-lg form-control-solid" disabled />
                <div class="form-text">登录用户名暂不支持修改</div>
              </div>
            </div>

            <div class="row mb-6">
              <label class="col-lg-4 col-form-label fw-semibold fs-6">显示名称</label>
              <div class="col-lg-8 fv-row">
                <input v-model.trim="profile.truename" type="text" maxlength="50" class="form-control form-control-lg form-control-solid" placeholder="输入管理员显示名称" />
              </div>
            </div>

            <div class="row mb-6">
              <label class="col-lg-4 col-form-label fw-semibold fs-6">个人简介</label>
              <div class="col-lg-8 fv-row">
                <input v-model.trim="profile.bio" type="text" maxlength="120" class="form-control form-control-lg form-control-solid" placeholder="例如：Golang Dev" />
                <div class="form-text">显示在侧边栏头像下方，最多 120 个字符。</div>
              </div>
            </div>

            <div class="row mb-0">
              <label class="col-lg-4 col-form-label fw-semibold fs-6">
                <span>联系邮箱</span>
                <i class="bi bi-info-circle ms-1 fs-7 text-muted" title="该邮箱会公开显示在友情链接中"></i>
              </label>
              <div class="col-lg-8 fv-row">
                <input v-model.trim="profile.email" type="email" maxlength="100" class="form-control form-control-lg form-control-solid" autocomplete="email" placeholder="name@example.com" />
                <div v-if="profileError" class="text-danger fs-7 mt-2">{{ profileError }}</div>
                <div class="form-text">该邮箱将作为博客公开联系邮箱显示。</div>
              </div>
            </div>
          </div>

          <div class="card-footer d-flex justify-content-end py-6 px-9">
            <button type="button" class="btn btn-light btn-active-light-primary me-2" :disabled="savingProfile" @click="loadProfile">放弃更改</button>
            <button type="submit" class="btn btn-primary" :disabled="savingProfile || uploadingAvatar">
              <span v-if="!savingProfile">保存更改</span>
              <span v-else>请稍候...<span class="spinner-border spinner-border-sm align-middle ms-2"></span></span>
            </button>
          </div>
        </form>
      </div>
    </div>

    <div class="card border-0 settings-card mb-6">
      <div class="card-header border-0 py-6 min-h-auto">
        <div class="card-title flex-column align-items-start m-0">
          <h3 class="fw-bold text-gray-900 mb-1">登录密码</h3>
          <span class="text-muted fs-7">密码修改成功后，需要使用新密码重新登录</span>
        </div>
      </div>

      <div class="card-body border-top p-7 p-lg-9">
        <div v-if="!editing" class="d-flex flex-column flex-md-row align-items-md-center justify-content-between gap-5">
          <div class="d-flex align-items-center gap-4">
            <span class="password-icon"><i class="bi bi-key-fill"></i></span>
            <div>
              <div class="fs-5 fw-bold text-gray-900 mb-1">Password</div>
              <div class="text-muted tracking-password">••••••••••••</div>
            </div>
          </div>
          <button type="button" class="btn btn-light-primary px-6" @click="editing = true">
            <i class="bi bi-pencil-square me-2"></i>修改密码
          </button>
        </div>

        <form v-else novalidate @submit.prevent="submit">
          <div class="notice d-flex bg-light-primary rounded border-primary border border-dashed mb-8 p-5">
            <i class="bi bi-info-circle-fill text-primary fs-2 me-4"></i>
            <div class="text-gray-700 fs-7">新密码建议同时包含大小写字母、数字和符号，长度为 8–72 位，且不能与当前密码相同。</div>
          </div>

          <div class="row g-5">
            <div class="col-12 col-lg-4">
              <label class="form-label fs-6 fw-semibold required">当前密码</label>
              <div class="position-relative">
                <input v-model="form.currentPassword" :type="showCurrent ? 'text' : 'password'" class="form-control form-control-lg form-control-solid pe-12" autocomplete="current-password" placeholder="输入当前密码" />
                <button type="button" class="password-visibility" :aria-label="showCurrent ? '隐藏当前密码' : '显示当前密码'" @click="showCurrent = !showCurrent">
                  <i :class="showCurrent ? 'bi bi-eye-slash' : 'bi bi-eye'"></i>
                </button>
              </div>
              <div v-if="errors.currentPassword" class="text-danger fs-7 mt-2">{{ errors.currentPassword }}</div>
            </div>

            <div class="col-12 col-lg-4">
              <label class="form-label fs-6 fw-semibold required">新密码</label>
              <div class="position-relative">
                <input v-model="form.newPassword" :type="showNew ? 'text' : 'password'" class="form-control form-control-lg form-control-solid pe-12" autocomplete="new-password" placeholder="输入新密码" />
                <button type="button" class="password-visibility" :aria-label="showNew ? '隐藏新密码' : '显示新密码'" @click="showNew = !showNew">
                  <i :class="showNew ? 'bi bi-eye-slash' : 'bi bi-eye'"></i>
                </button>
              </div>
              <div class="password-strength mt-3" aria-label="密码强度">
                <span v-for="index in 4" :key="index" :class="{ active: passwordStrength >= index }"></span>
              </div>
              <div v-if="errors.newPassword" class="text-danger fs-7 mt-2">{{ errors.newPassword }}</div>
            </div>

            <div class="col-12 col-lg-4">
              <label class="form-label fs-6 fw-semibold required">确认新密码</label>
              <input v-model="form.confirmPassword" type="password" class="form-control form-control-lg form-control-solid" autocomplete="new-password" placeholder="再次输入新密码" />
              <div v-if="errors.confirmPassword" class="text-danger fs-7 mt-2">{{ errors.confirmPassword }}</div>
            </div>
          </div>

          <div class="d-flex flex-wrap gap-3 mt-8">
            <button type="submit" class="btn btn-primary px-6" :disabled="submitting">
              <span v-if="submitting" class="spinner-border spinner-border-sm me-2"></span>
              {{ submitting ? "正在更新..." : "更新密码" }}
            </button>
            <button type="button" class="btn btn-light px-6" :disabled="submitting" @click="cancel">取消</button>
          </div>
        </form>
      </div>
    </div>

    <div class="notice d-flex bg-light-warning rounded border-warning border border-dashed p-6">
      <i class="bi bi-exclamation-triangle-fill text-warning fs-2x me-5"></i>
      <div>
        <h4 class="text-gray-900 fw-bold mb-1">密码安全建议</h4>
        <p class="text-gray-700 fs-7 mb-0">请勿复用其他网站的密码，也不要通过聊天、邮件或截图发送管理员密码。</p>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, reactive, ref } from "vue";
import { useRouter } from "vue-router";
import Swal from "sweetalert2";
import { ElNotification } from "element-plus";
import service from "@/utils/request";
import { useAuthStore } from "@/stores/auth";
import { getAssetPath } from "@/core/helpers/assets";
import type { User } from "@/stores/auth";

const router = useRouter();
const authStore = useAuthStore();
const editing = ref(false);
const submitting = ref(false);
const showCurrent = ref(false);
const showNew = ref(false);
const form = reactive({ currentPassword: "", newPassword: "", confirmPassword: "" });
const errors = reactive({ currentPassword: "", newPassword: "", confirmPassword: "" });
const profile = reactive<User>({} as User);
const savingProfile = ref(false);
const uploadingAvatar = ref(false);
const profileError = ref("");
const defaultAvatar = getAssetPath("/media/avatars/300-1.jpg");
const profileAvatar = computed(() => profile.pic || defaultAvatar);

const loadProfile = async () => {
  const data = await authStore.refreshProfile();
  Object.assign(profile, data || {});
};

const uploadAvatar = async (event: Event) => {
  const input = event.target as HTMLInputElement;
  const file = input.files?.[0];
  if (!file) return;
  uploadingAvatar.value = true;
  try {
    const formData = new FormData();
    formData.append("file", file);
    const { data } = await service.post("/admin/file/upload", formData);
    const uploadedUrl = String(data.url || "").trim();
    profile.pic = /^(?:https?:)?\/\//i.test(uploadedUrl) || uploadedUrl.startsWith("/")
      ? uploadedUrl
      : uploadedUrl.startsWith("public/") ? uploadedUrl : `public/${uploadedUrl}`;
    ElNotification({ title: "上传成功", message: "头像已上传，保存资料后生效", type: "success" });
  } finally {
    uploadingAvatar.value = false;
    input.value = "";
  }
};

const saveProfile = async () => {
  profileError.value = "";
  if (profile.email && !/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(profile.email)) {
    profileError.value = "请输入有效的邮箱地址";
    return;
  }
  savingProfile.value = true;
  try {
    const { data } = await service.put("/admin/user/profile", {
      trueName: profile.truename || "",
      email: profile.email || "",
      pic: profile.pic || "",
      bio: profile.bio || "",
    });
    Object.assign(profile, data);
    authStore.updateProfile(data);
    ElNotification({ title: "保存成功", message: "管理员资料已更新", type: "success" });
  } finally {
    savingProfile.value = false;
  }
};

const passwordStrength = computed(() => {
  const password = form.newPassword;
  if (!password) return 0;
  let score = password.length >= 8 ? 1 : 0;
  if (/[a-z]/.test(password) && /[A-Z]/.test(password)) score++;
  if (/\d/.test(password)) score++;
  if (/[^A-Za-z0-9]/.test(password)) score++;
  return score;
});

const resetForm = () => {
  form.currentPassword = "";
  form.newPassword = "";
  form.confirmPassword = "";
  errors.currentPassword = "";
  errors.newPassword = "";
  errors.confirmPassword = "";
  showCurrent.value = false;
  showNew.value = false;
};

const validate = () => {
  errors.currentPassword = form.currentPassword ? "" : "请输入当前密码";
  errors.newPassword = form.newPassword.length < 8
    ? "新密码至少需要 8 位"
    : form.newPassword.length > 72
      ? "新密码不能超过 72 位"
      : form.newPassword === form.currentPassword
        ? "新密码不能与当前密码相同"
        : "";
  errors.confirmPassword = form.confirmPassword !== form.newPassword ? "两次输入的新密码不一致" : "";
  return !errors.currentPassword && !errors.newPassword && !errors.confirmPassword;
};

const cancel = () => {
  resetForm();
  editing.value = false;
};

const submit = async () => {
  if (!validate() || submitting.value) return;
  submitting.value = true;
  try {
    await service.put("/admin/user/changePassword", {
      currentPassword: form.currentPassword,
      newPassword: form.newPassword,
    });
    await Swal.fire({
      title: "密码修改成功",
      text: "请使用新密码重新登录",
      icon: "success",
      confirmButtonText: "重新登录",
      buttonsStyling: false,
      heightAuto: false,
      customClass: { confirmButton: "btn btn-primary" },
    });
    authStore.logout();
    await router.replace({ name: "sign-in" });
  } finally {
    submitting.value = false;
  }
};

onMounted(() => loadProfile().catch(() => undefined));
</script>

<style scoped lang="scss">
.account-settings-page { max-width: 1320px; margin: 0 auto; }
.password-icon { display: inline-flex; align-items: center; justify-content: center; border-radius: 12px; }
.settings-card { box-shadow: 0 10px 32px rgba(31, 42, 74, .055); }
.password-icon { width: 52px; height: 52px; color: var(--kt-primary); background: var(--kt-primary-light); font-size: 1.35rem; }
.tracking-password { letter-spacing: .22em; }
.password-visibility { position: absolute; top: 50%; right: 12px; width: 34px; height: 34px; transform: translateY(-50%); border: 0; border-radius: 8px; background: transparent; color: var(--kt-gray-500); }
.password-visibility:hover { color: var(--kt-primary); background: var(--kt-primary-light); }
.password-strength { display: grid; grid-template-columns: repeat(4, 1fr); gap: 6px; }
.password-strength span { height: 4px; border-radius: 4px; background: var(--kt-gray-300); transition: background-color .2s ease; }
.password-strength span.active { background: var(--kt-success); }
.image-input { position: relative; display: inline-block; border-radius: .475rem; background-repeat: no-repeat; background-size: cover; }
.image-input-wrapper { border-radius: .475rem; background-repeat: no-repeat; background-position: center; background-size: cover; }
.image-input-outline .image-input-wrapper { border: 3px solid var(--kt-body-bg); box-shadow: var(--kt-box-shadow); }
.image-input [data-kt-image-input-action] { position: absolute; transform: translate(-50%, -50%); }
.image-input [data-kt-image-input-action="change"] { top: 0; left: 100%; cursor: pointer; }
.image-input [data-kt-image-input-action="change"] input { width: 0 !important; height: 0 !important; overflow: hidden; opacity: 0; }
.image-input [data-kt-image-input-action="remove"] { top: 100%; left: 100%; }
.avatar-loading { position: absolute; inset: 3px; display: flex; align-items: center; justify-content: center; border-radius: .475rem; background: rgba(255, 255, 255, .72); }
</style>
