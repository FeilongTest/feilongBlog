<template>
  <div class="seo-setting-page">
    <div class="card mb-5 mb-xl-10">
      <div class="card-header border-0">
        <div class="card-title m-0">
          <h3 class="fw-bold m-0">站点信息</h3>
        </div>
      </div>
      <div class="card-body border-top p-9">
        <div v-if="loading" class="placeholder-glow">
          <span v-for="item in 4" :key="item" class="placeholder col-12 d-block mb-6 py-4"></span>
        </div>
        <form v-else class="form" novalidate @submit.prevent="save">
          <div class="row mb-6">
            <label class="col-lg-4 col-form-label fw-semibold fs-6">站点名称</label>
            <div class="col-lg-8 fv-row">
              <input v-model.trim="form.siteName" type="text" maxlength="60" class="form-control form-control-lg form-control-solid" placeholder="例如：飞龙博客" />
              <div class="form-text">出现在浏览器标题与搜索结果的第一行。</div>
            </div>
          </div>

          <div class="row mb-6">
            <label class="col-lg-4 col-form-label fw-semibold fs-6">站点副标题</label>
            <div class="col-lg-8 fv-row">
              <input v-model.trim="form.siteSlogan" type="text" maxlength="120" class="form-control form-control-lg form-control-solid" placeholder="例如：记录技术与生活" />
            </div>
          </div>

          <div class="row mb-6">
            <label class="col-lg-4 col-form-label fw-semibold fs-6">站点描述</label>
            <div class="col-lg-8 fv-row">
              <textarea v-model.trim="form.description" rows="3" maxlength="200" class="form-control form-control-lg form-control-solid" placeholder="一句话介绍这个博客的内容方向"></textarea>
              <div class="form-text">对应搜索引擎结果里的摘要，建议 60 到 120 个字符。</div>
            </div>
          </div>

          <div class="row mb-6">
            <label class="col-lg-4 col-form-label fw-semibold fs-6">站点关键词</label>
            <div class="col-lg-8 fv-row">
              <input v-model.trim="form.keywords" type="text" maxlength="200" class="form-control form-control-lg form-control-solid" placeholder="用英文逗号分隔，例如：个人博客,Go,Vue" />
            </div>
          </div>

          <div class="row mb-6">
            <label class="col-lg-4 col-form-label fw-semibold fs-6">作者</label>
            <div class="col-lg-8 fv-row">
              <input v-model.trim="form.author" type="text" maxlength="60" class="form-control form-control-lg form-control-solid" placeholder="文章结构化数据中的作者名称" />
            </div>
          </div>

          <div class="row mb-9">
            <label class="col-lg-4 col-form-label fw-semibold fs-6">备案信息</label>
            <div class="col-lg-8 fv-row">
              <input v-model.trim="form.icp" type="text" maxlength="80" class="form-control form-control-lg form-control-solid" placeholder="选填" />
            </div>
          </div>

          <div class="separator separator-dashed my-8"></div>

          <h4 class="fw-bold mb-2">搜索引擎验证</h4>
          <p class="text-muted fs-7 mb-6">填写站长平台给出的验证码后，前台页面会自动输出对应的验证 meta 标签。</p>

          <div class="row mb-6">
            <label class="col-lg-4 col-form-label fw-semibold fs-6">百度验证码</label>
            <div class="col-lg-8 fv-row">
              <input v-model.trim="form.baiduVerify" type="text" maxlength="120" class="form-control form-control-lg form-control-solid" placeholder="codeva-xxxxxxxx" />
            </div>
          </div>

          <div class="row mb-6">
            <label class="col-lg-4 col-form-label fw-semibold fs-6">Google 验证码</label>
            <div class="col-lg-8 fv-row">
              <input v-model.trim="form.googleVerify" type="text" maxlength="120" class="form-control form-control-lg form-control-solid" placeholder="来自 Search Console 的 HTML 标记" />
            </div>
          </div>

          <div class="row mb-9">
            <label class="col-lg-4 col-form-label fw-semibold fs-6">Bing 验证码</label>
            <div class="col-lg-8 fv-row">
              <input v-model.trim="form.bingVerify" type="text" maxlength="120" class="form-control form-control-lg form-control-solid" placeholder="msvalidate.01 的值" />
            </div>
          </div>

          <button type="submit" class="btn btn-primary px-6" :disabled="saving">
            <span v-if="saving" class="spinner-border spinner-border-sm me-2"></span>
            {{ saving ? "正在保存..." : "保存设置" }}
          </button>
        </form>
      </div>
    </div>

    <div class="card">
      <div class="card-header border-0">
        <div class="card-title m-0">
          <h3 class="fw-bold m-0">收录入口</h3>
        </div>
      </div>
      <div class="card-body border-top p-9">
        <p class="text-muted fs-7">把下面的地址提交到各搜索引擎的站长平台，可以更快发现新文章。文章发布或更新后，系统也会自动推送。</p>
        <div class="d-flex flex-wrap gap-4">
          <a class="btn btn-sm btn-light-primary" :href="sitemapUrl" target="_blank" rel="noopener">{{ sitemapUrl }}</a>
          <a class="btn btn-sm btn-light" :href="robotsUrl" target="_blank" rel="noopener">{{ robotsUrl }}</a>
          <a class="btn btn-sm btn-light" :href="rssUrl" target="_blank" rel="noopener">{{ rssUrl }}</a>
        </div>
        <div class="notice d-flex bg-light-info rounded border-info border border-dashed p-6 mt-8">
          <i class="bi bi-info-circle-fill text-info fs-2x me-5"></i>
          <div class="text-gray-700 fs-7">
            搜索收录需要时间，通常几天到数周。保存设置后请刷新前台页面查看效果，并在站长平台重新抓取首页。
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { onMounted, reactive, ref } from "vue";
import { ElNotification } from "element-plus";
import service from "@/utils/request";

const loading = ref(true);
const saving = ref(false);
const form = reactive({
  siteName: "",
  siteSlogan: "",
  description: "",
  keywords: "",
  author: "",
  icp: "",
  baiduVerify: "",
  googleVerify: "",
  bingVerify: "",
});

const origin = window.location.origin;
const sitemapUrl = `${origin}/sitemap.xml`;
const robotsUrl = `${origin}/robots.txt`;
const rssUrl = `${origin}/rss.xml`;

const loadSetting = async () => {
  loading.value = true;
  try {
    const res: any = await service.get("/admin/setting/getSetting");
    const data = res.data || {};
    form.siteName = data.siteName || "";
    form.siteSlogan = data.siteSlogan || "";
    form.description = data.description || "";
    form.keywords = data.keywords || "";
    form.author = data.author || "";
    form.icp = data.icp || "";
    form.baiduVerify = data.baiduVerify || "";
    form.googleVerify = data.googleVerify || "";
    form.bingVerify = data.bingVerify || "";
  } finally {
    loading.value = false;
  }
};

const save = async () => {
  saving.value = true;
  try {
    await service.put("/admin/setting/updateSetting", { ...form });
    ElNotification({ title: "保存成功", message: "站点信息已更新，前台刷新后即可生效", type: "success" });
  } finally {
    saving.value = false;
  }
};

onMounted(loadSetting);
</script>
