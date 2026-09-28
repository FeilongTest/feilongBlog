import {
  createRouter,
  createWebHistory,
  type RouteRecordRaw,
} from "vue-router";
import { useAuthStore } from "@/stores/auth";
import { useConfigStore } from "@/stores/config";
import { useSiteStore } from "@/stores/site";
import JwtService from "@/core/services/JwtService";

const routes: Array<RouteRecordRaw> = [
  {
    // 首页统一使用根路径，旧地址永久跳转
    path: "/index",
    redirect: "/",
  },
  {
    path: "/admin",
    redirect: "/dashboard",
    component: () => import("@/layouts/main-layout/MainLayout.vue"),
    meta: {
      middleware: "auth",
    },
    children: [
      {
        path: "/dashboard",
        name: "dashboard",
        component: () => import("@/views/Dashboard.vue"),
        meta: {
          pageTitle: "Dashboard",
          breadcrumbs: ["Dashboards"],
        },
      },
      {
        path: "/admin/category",
        name: "AdminCategory",
        component: () => import("@/views/admin/category/index.vue"),
        meta: {
          pageTitle: "分类管理",
          breadcrumbs: ["category"],
        },
      },
      {
        path: "/admin/article/",
        name: "AdminArticle",
        meta: {
          pageTitle: "文章管理",
          breadcrumbs: ["article"],
        },
        children: [
          {
            path: "index",
            name: "admin-article-index",
            component: () => import("@/views/admin/article/index.vue"),
          },
          {
            path: "content/:id?",
            name: "admin-article-content",
            component: () => import("@/views/admin/article/content.vue"),
          },
        ],
      },
      {
        path: "/admin/comment",
        name: "AdminComment",
        component: () => import("@/views/admin/comment/index.vue"),
        meta: {
          pageTitle: "评论管理",
          breadcrumbs: ["comment"],
        },
      },
      {
        path: "/admin/friendlink",
        name: "AdminFriendlink",
        component: () => import("@/views/admin/friendlink/index.vue"),
        meta: {
          pageTitle: "友链管理",
          breadcrumbs: ["friendlink"],
        },
      },
      {
        path: "/admin/account/settings",
        name: "admin-account-settings",
        component: () => import("@/views/admin/account/Settings.vue"),
        meta: {
          pageTitle: "账户设置",
          breadcrumbs: ["account", "settings"],
        },
      },
      {
        path: "/admin/setting/seo",
        name: "admin-setting-seo",
        component: () => import("@/views/admin/setting/SeoSetting.vue"),
        meta: {
          pageTitle: "站点设置",
          breadcrumbs: ["setting", "seo"],
        },
      },

    ],
  },
  {
    path: "/",
    component: () => import("@/layouts/main-layout/MainLayout.vue"),
    children: [
      {
        path: "",
        name: "blog-home",
        component: () =>
          import("@/views/blog/Feeds.vue"),
        meta: {
          pageTitle: "首页",
          breadcrumbs: ["Pages", "Index"],
        },
      },
      {
        path: "content/:fid/:id",
        name: "blog-content",
        component: () =>
          import("@/views/blog/Post.vue"),
        meta: {
          pageTitle: "文章详情",
          breadcrumbs: ["Pages", "Article"],
        },
      },
      {
        path: "category/:id",
        name: "category",
        component: () => import("@/views/blog/Feeds.vue"),
        meta: {
          pageTitle: "分类",
          breadcrumbs: ["Pages", "Article"],
        },
      },
    ],
  },
  {
    path: "/",
    component: () => import("@/layouts/AuthLayout.vue"),
    children: [
      {
        path: "/sign-in",
        name: "sign-in",
        component: () =>
          import("@/views/crafted/authentication/basic-flow/SignIn.vue"),
        meta: {
          pageTitle: "Sign In",
        },
      },
    ],
  },
  {
    path: "/",
    component: () => import("@/layouts/SystemLayout.vue"),
    children: [
      {
        // the 404 route, when none of the above matches
        path: "/404",
        name: "404",
        component: () => import("@/views/crafted/authentication/Error404.vue"),
        meta: {
          pageTitle: "Error 404",
        },
      },
      {
        path: "/500",
        name: "500",
        component: () => import("@/views/crafted/authentication/Error500.vue"),
        meta: {
          pageTitle: "Error 500",
        },
      },
    ],
  },
  {
    path: "/:pathMatch(.*)*",
    redirect: "/404",
  },
];

// 兼容旧版 hash 路由的分享地址，例如 /#/content/1/2 会跳转到 /content/1/2
// 使用 history 模式后地址对搜索引擎才是有意义的独立页面
if (window.location.hash.startsWith("#/")) {
  window.history.replaceState(null, "", window.location.hash.slice(1) || "/");
}

const router = createRouter({
  history: createWebHistory(),
  routes,
});

router.beforeEach((to) => {
  const authStore = useAuthStore();
  const configStore = useConfigStore();

  // 更新认证状态（基于token是否存在）
	const authenticated = JwtService.isTokenValid();
	authStore.isAuthenticated = authenticated;
	if (!authenticated) JwtService.destroyToken();

  // 根路径"/"就是博客首页，管理员通过菜单进入后台

  // current page view title
  const siteStore = useSiteStore();
  const appName = siteStore.setting.siteName || import.meta.env.VITE_APP_NAME || "飞龙博客";
	document.title = to.meta.pageTitle ? `${to.meta.pageTitle} - ${appName}` : appName;

  // reset config to initial state
  configStore.resetLayoutConfig();

  // before page access check if page requires authentication
  if (to.meta.middleware == "auth") {
    // verify auth token before each page change
	if (!authStore.isAuthenticated) {
	  return { name: "sign-in", query: { redirect: to.fullPath } };
	}
  }

  // Scroll page to top on every route change
  window.scrollTo({
    top: 0,
    left: 0,
    behavior: "smooth",
  });
	return true;
});

export default router;
