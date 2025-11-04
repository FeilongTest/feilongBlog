import { getAssetPath } from "@/core/helpers/assets";

export interface MenuItem {
  heading?: string;
  sectionTitle?: string;
  route?: string;
  pages?: Array<MenuItem>;
  svgIcon?: string;
  fontIcon?: string;
  sub?: Array<MenuItem>;
}

const MainMenuConfig: Array<MenuItem> = [
  {
    heading: "仪表盘",
    route: "/dashboard",
    svgIcon: getAssetPath("/media/icons/duotune/general/gen025.svg"),
    fontIcon: "bi-app-indicator",
  },
  {
    heading: "分类管理",
    route: "/admin/category",
    svgIcon: getAssetPath("/media/icons/duotune/general/gen022.svg"),
    fontIcon: "bi-folder",
  },
  {
    sectionTitle: "文章管理",
    route: "/admin/article",
    svgIcon: getAssetPath("/media/icons/duotune/files/fil003.svg"),
    fontIcon: "bi-file-text",
    pages: [
      {
        heading: "查看文章",
        route: "/admin/article/index",
      },
      {
        heading: "添加文章",
        route: "/admin/article/content",
      }
    ]
  },
  {
    heading: "评论管理",
    route: "/admin/comment",
    svgIcon: getAssetPath("/media/icons/duotune/communication/com012.svg"),
    fontIcon: "bi-chat-dots",
  },
  {
    heading: "友链管理",
    route: "/admin/friendlink",
    svgIcon: getAssetPath("/media/icons/duotune/general/gen016.svg"),
    fontIcon: "bi-link-45deg",
  },
];

export default MainMenuConfig;
