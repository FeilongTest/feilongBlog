import type { Page } from "@playwright/test";
export const avatar = "/test-author.svg";
export const post = {
  status: 0,
  version: 1,
  tags: [{ id: 1, name: "Go", description: "Go 主题" }],
  tagNames: ["Go"],
  ID: 123,
  fid: 1,
  title: "回归检查文章",
  content: "# Markdown 标题\n\n文章内容",
  contentFormat: "markdown",
  ctime: 1700000000,
  type: 1,
  pic: "",
  view: 1,
};
export const profile = {
  ID: 1,
  admin: 1,
  pic: avatar,
  truename: "回归作者",
  bio: "测试简介",
  email: "test@example.com",
};

export async function mockAPI(page: Page) {
  await page.route("**/test-author.svg", (route) =>
    route.fulfill({
      contentType: "image/svg+xml",
      body: '<svg xmlns="http://www.w3.org/2000/svg" width="40" height="40"><rect width="40" height="40" fill="green"/></svg>',
    })
  );
  await page.route("**/blog/**", (route) => {
    const url = new URL(route.request().url());
    let data: unknown = {};
    if (url.pathname.endsWith("/getContact"))
      data = { ...profile, trueName: profile.truename };
    if (url.pathname.endsWith("/profile")) data = profile;
    if (
      url.pathname.endsWith("/getSiteSetting") ||
      url.pathname.endsWith("/getSetting")
    )
      data = {
        siteName: "测试博客",
        siteSlogan: "记录",
        description: "测试描述",
        author: "回归作者",
      };
    if (url.pathname.endsWith("/getCategoryList"))
      data = [{ ID: 1, name: "测试分类", fid: 0 }];
    if (url.pathname.endsWith("/getArticleList"))
      data = {
        list: [
          {
            ...post,
            title:
              url.searchParams.get("page") === "2" ? "第二页文章" : post.title,
          },
        ],
        total: 20,
      };
    if (url.pathname.endsWith("/getArticle")) data = post;
    if (url.pathname.endsWith("/getArchive"))
      data = { list: [{ ...post, id: post.ID }] };
    if (url.pathname.endsWith("/getTags") || url.pathname.endsWith("/tags"))
      data = [{ id: 1, name: "Go", description: "Go 主题", count: 1 }];
    if (url.pathname.endsWith("/getTag"))
      data = { id: 1, name: "Go", description: "Go 主题", count: 1 };
    if (url.pathname.endsWith("/getReadingArticles")) data = [post];
    if (url.pathname.endsWith("/getPostAside")) data = { related: [], hot: [] };
    if (
      url.pathname.endsWith("/getAllFriendlinks") ||
      url.pathname.endsWith("/getSummary")
    )
      data = [];
    if (url.pathname.endsWith("/getCommentList")) data = { list: [], total: 0 };
    return route.fulfill({
      contentType: "application/json",
      body: JSON.stringify({ code: 0, data, msg: "成功" }),
    });
  });
}
