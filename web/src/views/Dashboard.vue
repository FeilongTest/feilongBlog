<!-- eslint-disable vue/multi-word-component-names -->
<template>
  <div class="row g-5 g-xl-8 mb-8">
    <div v-for="(item, index) in statCards" :key="item.label" class="col-sm-6 col-xl-4 stat-card-wrap" :style="{ '--stat-index': index }">
      <router-link :to="item.to" class="card h-100 card-hover border-0">
        <div class="card-body p-6">
          <div class="d-flex align-items-center justify-content-between mb-6">
            <span class="symbol symbol-45px"><span class="symbol-label" :class="`bg-light-${item.color}`"><i :class="`${item.icon} text-${item.color} fs-2`"></i></span></span>
            <span class="badge badge-light">查看详情</span>
          </div>
          <div class="fs-2qx fw-bolder text-gray-900 mb-1"><AnimatedNumber :value="item.value" /></div>
          <div class="text-muted fw-semibold">{{ item.label }}</div>
        </div>
      </router-link>
    </div>
  </div>

  <div class="card mb-8">
    <div class="card-header border-0 pt-3"><div class="card-title flex-column align-items-start"><h3 class="fw-bold mb-1">最近 7 天访问趋势</h3><span class="text-muted fs-7">页面浏览量（PV）与独立访客（UV）</span></div></div>
    <div class="card-body pt-2">
      <div class="visit-chart d-flex align-items-end gap-3 gap-md-6">
        <div v-for="item in stats.trend" :key="item.date" class="visit-column flex-grow-1 text-center">
          <div class="d-flex align-items-end justify-content-center gap-1 visit-bars">
            <div class="visit-bar bg-primary" :style="{ height: `${Math.max(4, item.pv / maxPv * 150)}px` }" :title="`PV ${item.pv}`"></div>
            <div class="visit-bar bg-success" :style="{ height: `${Math.max(4, item.uv / maxPv * 150)}px` }" :title="`UV ${item.uv}`"></div>
          </div>
          <div class="text-gray-600 fs-8 mt-3">{{ item.date.slice(5) }}</div><div class="text-muted fs-9">{{ item.pv }}/{{ item.uv }}</div>
        </div>
      </div>
      <div class="d-flex gap-5 mt-5 fs-7 text-muted"><span><i class="bullet bullet-dot bg-primary me-2"></i>PV</span><span><i class="bullet bullet-dot bg-success me-2"></i>UV</span></div>
    </div>
  </div>

  <div class="row g-5 g-xl-8 mb-8">
    <div class="col-xl-8">
      <div class="card h-100">
        <div class="card-header border-0 pt-3"><div class="card-title flex-column align-items-start"><h3 class="fw-bold mb-1">最近文章</h3><span class="text-muted fs-7">快速检查内容状态与阅读表现</span></div><div class="card-toolbar"><router-link to="/admin/article/content" class="btn btn-sm btn-primary"><i class="bi bi-plus-lg me-2"></i>新建文章</router-link></div></div>
        <div class="card-body pt-2">
          <div class="table-responsive">
            <table class="table table-row-dashed align-middle gy-5">
              <thead><tr class="text-muted fw-bold fs-7 text-uppercase"><th>文章</th><th>状态</th><th class="text-end">阅读</th><th class="text-end">操作</th></tr></thead>
              <tbody><tr v-for="article in recentArticles" :key="article.ID"><td><div class="d-flex align-items-center"><span class="symbol symbol-50px me-4"><span class="symbol-label bg-light-primary"><i class="bi bi-file-earmark-text text-primary fs-2"></i></span></span><div><router-link :to="`/admin/article/content/${article.ID}`" class="text-gray-900 text-hover-primary fw-bold d-block mw-300px text-truncate">{{ article.title }}</router-link><span class="text-muted fs-7">{{ dayjs.unix(article.ctime).format('YYYY-MM-DD HH:mm') }}</span></div></div></td><td><span :class="article.status===0?'badge badge-light-success':'badge badge-light-secondary'">{{ article.status===0?'已发布':'已隐藏' }}</span></td><td class="text-end fw-bold">{{ article.view }}</td><td class="text-end"><router-link :to="`/admin/article/content/${article.ID}`" class="btn btn-sm btn-icon btn-light-primary"><i class="bi bi-pencil"></i></router-link></td></tr></tbody>
            </table>
          </div>
        </div>
      </div>
    </div>

    <div class="col-xl-4">
      <div class="card h-100 admin-welcome overflow-hidden">
        <div class="card-body p-8 position-relative">
          <div class="welcome-circle"></div><span class="badge badge-light-primary mb-5">BLOG CONSOLE</span><h2 class="fw-bolder text-gray-900 mb-3">欢迎回来</h2><p class="text-gray-600 lh-lg mb-8">内容、评论和对象存储都可以从这里统一管理。</p>
          <div class="d-grid gap-3"><router-link to="/admin/article/content" class="btn btn-primary"><i class="bi bi-pencil-square me-2"></i>撰写新文章</router-link><router-link to="/index" class="btn btn-light-primary"><i class="bi bi-box-arrow-up-right me-2"></i>查看博客首页</router-link></div>
        </div>
      </div>
    </div>
  </div>

  <div class="card">
    <div class="card-header border-0"><div class="card-title flex-column align-items-start"><h3 class="fw-bold mb-1">快捷管理</h3><span class="text-muted fs-7">常用后台功能</span></div></div>
    <div class="card-body pt-2"><div class="row g-4"><div v-for="action in quickActions" :key="action.label" class="col-6 col-md-4 col-xl-2"><router-link :to="action.to" class="quick-action d-flex flex-column align-items-center justify-content-center rounded p-5 h-100"><span class="symbol symbol-50px mb-3"><span class="symbol-label" :class="`bg-light-${action.color}`"><i :class="`${action.icon} text-${action.color} fs-2`"></i></span></span><span class="text-gray-800 fw-bold text-center">{{ action.label }}</span></router-link></div></div></div>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, ref } from "vue";
import dayjs from "dayjs";
import AnimatedNumber from "@/components/admin/AnimatedNumber.vue";
import type { ArticleList } from "@/core/blog/ArticleTypes";
import service from "@/utils/request";

const stats=ref({articleCount:0,categoryCount:0,commentCount:0,friendlinkCount:0,likeCount:0,totalPv:0,totalUv:0,trend:[] as Array<{date:string;pv:number;uv:number}>});
const recentArticles=ref<ArticleList[]>([]);
const statCards=computed(()=>[
  {label:'文章总数',value:stats.value.articleCount,icon:'bi bi-file-earmark-text',color:'primary',to:'/admin/article/index'},
  {label:'分类数量',value:stats.value.categoryCount,icon:'bi bi-folder2-open',color:'success',to:'/admin/category'},
  {label:'评论数量',value:stats.value.commentCount,icon:'bi bi-chat-square-dots',color:'info',to:'/admin/comment'},
  {label:'友情链接',value:stats.value.friendlinkCount,icon:'bi bi-link-45deg',color:'warning',to:'/admin/friendlink'},
  {label:'累计浏览',value:stats.value.totalPv,icon:'bi bi-eye',color:'danger',to:'/dashboard'},
  {label:'独立访客',value:stats.value.totalUv,icon:'bi bi-people',color:'dark',to:'/dashboard'},
]);
const maxPv=computed(()=>Math.max(1,...stats.value.trend.map(item=>item.pv)));
const quickActions=[
  {label:'文章管理',icon:'bi bi-file-text',color:'primary',to:'/admin/article/index'}, {label:'新建文章',icon:'bi bi-plus-circle',color:'success',to:'/admin/article/content'}, {label:'分类管理',icon:'bi bi-folder',color:'info',to:'/admin/category'}, {label:'评论审核',icon:'bi bi-chat-dots',color:'warning',to:'/admin/comment'}, {label:'友情链接',icon:'bi bi-link',color:'danger',to:'/admin/friendlink'}, {label:'博客首页',icon:'bi bi-house',color:'dark',to:'/index'},
];

onMounted(async()=>{
  const [statsRes,articleRes]=await Promise.all([
    service.get('/admin/statistics/dashboard'), service.get('/base/getArticleList',{params:{page:1,pageSize:6}}),
  ]);
  stats.value={...stats.value,...statsRes.data}; recentArticles.value=articleRes.data?.list||[];
});
</script>

<style scoped>
.card-hover,.quick-action { transition: transform .2s ease,box-shadow .2s ease; }
.stat-card-wrap { animation: stat-card-enter .45s ease both; animation-delay: calc(var(--stat-index) * 70ms); }
.card-hover:hover,.quick-action:hover { transform: translateY(-3px); box-shadow:0 .75rem 2rem rgba(31,38,135,.08); }
.quick-action { background:var(--bs-gray-100); }
.admin-welcome { background:linear-gradient(145deg,#f1f7ff,#fbf9ff); }
.welcome-circle { position:absolute;width:190px;height:190px;border-radius:50%;background:rgba(62,151,255,.12);right:-55px;top:-65px; }
.fs-2qx{font-size:2.25rem!important}
.visit-chart{min-height:190px;border-bottom:1px dashed var(--bs-gray-300)}
.visit-bars{height:160px}.visit-bar{width:10px;border-radius:6px 6px 2px 2px;transition:height .45s ease}
@keyframes stat-card-enter { from { opacity:0; transform:translateY(10px); } to { opacity:1; transform:translateY(0); } }
@media(prefers-reduced-motion:reduce){.card-hover,.quick-action{transition:none}.stat-card-wrap{animation:none}}
</style>
