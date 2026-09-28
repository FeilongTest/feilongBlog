import { ref } from "vue";
import { defineStore } from "pinia";
import service from "@/utils/request";

export interface SiteSetting {
  siteName: string;
  siteSlogan: string;
  description: string;
  keywords: string;
  author: string;
  icp: string;
}

export interface SiteCategory {
  ID: number;
  name: string;
  fid: number;
}

export const defaultSiteSetting: SiteSetting = {
  siteName: "飞龙博客",
  siteSlogan: "记录技术与生活",
  description: "一个分享编程技术、折腾经验与日常记录的个人博客。",
  keywords: "个人博客,技术博客,编程,Go,Vue",
  author: "飞龙",
  icp: "",
};

let siteSettingRequest: Promise<SiteSetting> | null = null;
let categoryRequest: Promise<SiteCategory[]> | null = null;

export const useSiteStore = defineStore("site", () => {
  const setting = ref<SiteSetting>({ ...defaultSiteSetting });
  const categories = ref<SiteCategory[]>([]);
  const loaded = ref(false);

  // 站点信息同一会话内只请求一次
  async function loadSiteSetting() {
    if (loaded.value) return setting.value;
    if (!siteSettingRequest) {
      siteSettingRequest = service
        .get("/base/getSiteSetting")
        .then(({ data }) => {
          setting.value = { ...defaultSiteSetting, ...(data || {}) };
          loaded.value = true;
          return setting.value;
        })
        .catch(() => setting.value)
        .finally(() => {
          siteSettingRequest = null;
        });
    }
    return siteSettingRequest;
  }

  // 分类列表用于生成分类页标题与列表
  async function loadCategories() {
    if (categories.value.length) return categories.value;
    if (!categoryRequest) {
      categoryRequest = service
        .get("/base/getCategoryList")
        .then(({ data }) => {
          categories.value = data || [];
          return categories.value;
        })
        .catch(() => categories.value)
        .finally(() => {
          categoryRequest = null;
        });
    }
    return categoryRequest;
  }

  function categoryName(id: number) {
    return categories.value.find((item) => Number(item.ID) === Number(id))?.name || "";
  }

  return { setting, categories, loaded, loadSiteSetting, loadCategories, categoryName };
});
