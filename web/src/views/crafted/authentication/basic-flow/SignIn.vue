<template>
  <div class="login-panel">
    <div class="mb-10">
      <span class="login-kicker">ADMIN CONSOLE</span>
      <h2 class="login-title">欢迎回来</h2>
      <p class="login-subtitle">登录后管理你的博客内容与站点数据。</p>
    </div>

    <VForm
      id="kt_login_signin_form"
      class="form w-100"
      :validation-schema="login"
      :initial-values="{ username: '', password: '' }"
      @submit="onSubmitLogin"
    >
      <div class="fv-row mb-7">
        <label class="form-label fs-6 fw-semibold text-gray-800">用户名</label>
        <Field
          class="form-control form-control-lg form-control-solid login-input"
          type="text"
          name="username"
          autocomplete="username"
          placeholder="请输入用户名"
        />
        <div class="fv-plugins-message-container"><div class="fv-help-block"><ErrorMessage name="username" /></div></div>
      </div>

      <div class="fv-row mb-8">
        <label class="form-label fs-6 fw-semibold text-gray-800">密码</label>
        <Field
          class="form-control form-control-lg form-control-solid login-input"
          type="password"
          name="password"
          autocomplete="current-password"
          placeholder="请输入密码"
        />
        <div class="fv-plugins-message-container"><div class="fv-help-block"><ErrorMessage name="password" /></div></div>
      </div>

      <button ref="submitButton" id="kt_sign_in_submit" type="submit" class="btn btn-lg btn-primary w-100 login-submit">
        <span class="indicator-label">登录后台 <i class="bi bi-arrow-right ms-2"></i></span>
        <span class="indicator-progress">正在登录… <span class="spinner-border spinner-border-sm align-middle ms-2"></span></span>
      </button>

      <router-link :to="{ name: 'blog-home' }" class="back-home"><i class="bi bi-arrow-left"></i> 返回博客首页</router-link>
    </VForm>
  </div>
</template>

<script setup lang="ts">
import { ref } from "vue";
import { ErrorMessage, Field, Form as VForm } from "vee-validate";
import { useRouter } from "vue-router";
import Swal from "sweetalert2";
import * as Yup from "yup";
import { useAuthStore, type User } from "@/stores/auth";

const authStore = useAuthStore();
const router = useRouter();
const submitButton = ref<HTMLButtonElement | null>(null);

const login = Yup.object({
  username: Yup.string().required("请输入用户名"),
  password: Yup.string().min(4, "密码至少需要 4 位").required("请输入密码"),
});

const onSubmitLogin = async (values: Record<string, unknown>) => {
  authStore.logout();
  if (submitButton.value) {
    submitButton.value.disabled = true;
    submitButton.value.setAttribute("data-kt-indicator", "on");
  }

  try {
    await authStore.login(values as unknown as User);
    const error = Object.values(authStore.errors)[0];
    if (error) {
      await Swal.fire({
        text: String(error), icon: "error", buttonsStyling: false, confirmButtonText: "重试", heightAuto: false,
        customClass: { confirmButton: "btn fw-semibold btn-light-danger" },
      });
      authStore.errors = {};
      return;
    }

    await Swal.fire({
      text: "登录成功", icon: "success", buttonsStyling: false, confirmButtonText: "进入后台", heightAuto: false,
      customClass: { confirmButton: "btn fw-semibold btn-light-primary" },
    });
    const redirect = router.currentRoute.value.query.redirect;
    await router.push(typeof redirect === "string" && redirect.startsWith("/") ? redirect : { name: "dashboard" });
  } finally {
    submitButton.value?.removeAttribute("data-kt-indicator");
    if (submitButton.value) submitButton.value.disabled = false;
  }
};
</script>

<style scoped>
.login-panel { width: min(100%, 430px); padding: 2rem 0; }
.login-kicker { color: var(--kt-primary); font-size: .75rem; font-weight: 700; letter-spacing: .18em; }
.login-title { margin: 1rem 0 .65rem; color: #181c32; font-size: clamp(2rem, 4vw, 2.65rem); font-weight: 750; letter-spacing: -.04em; }
.login-subtitle { margin: 0; color: #7e8299; font-size: 1rem; }
.login-input { min-height: 54px; border: 1px solid transparent; border-radius: .75rem; transition: border-color .2s ease, box-shadow .2s ease, background-color .2s ease; }
.login-input:focus { border-color: rgba(0, 158, 247, .4); background: #fff; box-shadow: 0 0 0 .25rem rgba(0, 158, 247, .08); }
.login-submit { min-height: 52px; border-radius: .75rem; font-weight: 650; box-shadow: 0 .75rem 1.75rem rgba(0, 158, 247, .18); }
.back-home { display: flex; justify-content: center; align-items: center; gap: .5rem; margin-top: 1.5rem; color: #7e8299; font-size: .9rem; font-weight: 600; }
.back-home:hover { color: var(--kt-primary); }
</style>
