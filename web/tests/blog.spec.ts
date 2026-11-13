import { test, expect, type Page } from "@playwright/test";

const avatar = "/test-author.svg";
const post = { ID: 123, fid: 1, title: "回归检查文章", content: "# Markdown 标题\n\n文章内容", contentFormat: "markdown", ctime: 1700000000, type: 1, pic: "", view: 1 };
const profile = { ID: 1, admin: 1, pic: avatar, truename: "回归作者", bio: "测试简介", email: "test@example.com" };

async function mockAPI(page: Page) {
  await page.route("**/test-author.svg", route => route.fulfill({ contentType: "image/svg+xml", body: '<svg xmlns="http://www.w3.org/2000/svg" width="40" height="40"><rect width="40" height="40" fill="green"/></svg>' }));
  await page.route("**/blog/**", route => {
    const url = new URL(route.request().url());
    let data: unknown = {};
    if (url.pathname.endsWith("/getContact")) data = { ...profile, trueName: profile.truename };
    if (url.pathname.endsWith("/profile")) data = profile;
    if (url.pathname.endsWith("/getSiteSetting") || url.pathname.endsWith("/getSetting")) data = { siteName: "测试博客", siteSlogan: "记录", description: "测试描述", author: "回归作者" };
    if (url.pathname.endsWith("/getCategoryList")) data = [{ ID: 1, name: "测试分类", fid: 0 }];
    if (url.pathname.endsWith("/getArticleList")) data = { list: [{ ...post, title: url.searchParams.get("page") === "2" ? "第二页文章" : post.title }], total: 20 };
    if (url.pathname.endsWith("/getArticle")) data = post;
    if (url.pathname.endsWith("/getArchive")) data = { list: [post] };
    if (url.pathname.endsWith("/getPostAside")) data = { related: [], hot: [] };
    if (url.pathname.endsWith("/getAllFriendlinks") || url.pathname.endsWith("/getSummary")) data = [];
    if (url.pathname.endsWith("/getCommentList")) data = { list: [], total: 0 };
    return route.fulfill({ contentType: "application/json", body: JSON.stringify({ code: 0, data, msg: "成功" }) });
  });
}

test.beforeEach(async ({ page }) => { await mockAPI(page); });

test("作者头像在首页、文章和关于页一致，Markdown 摘要无原始标记", async ({ page }) => {
  const errors: string[] = [];
  page.on("pageerror", error => errors.push(error.message));
  for (const path of ["/", "/content/1/123", "/about"]) {
    await page.goto(path, { waitUntil: "domcontentloaded" });
    const authorImages = page.locator(`img[src="${avatar}"]`);
    await expect(authorImages.first()).toBeVisible();
    await expect.poll(() => authorImages.count()).toBeGreaterThanOrEqual(3);
    await expect.poll(() => authorImages.first().evaluate((image: HTMLImageElement) => image.naturalWidth)).toBeGreaterThan(0);
    await expect(page.locator('img[src$="/media/avatars/300-1.jpg"]')).toHaveCount(0);
    if (path === "/") await expect(page.locator(".feed-content")).not.toContainText("#");
  }
  expect(errors).toEqual([]);
});

test("头像资源失效时使用统一兜底", async ({ page }) => {
  await page.route("**/test-author.svg", route => route.fulfill({ status: 404, body: "missing" }));
  await page.goto("/about", { waitUntil: "domcontentloaded" });
  await expect(page.locator(".about-avatar")).toHaveAttribute("src", /\/media\/avatars\/300-1.jpg$/);
  await expect.poll(() => page.locator(".about-avatar").evaluate((image: HTMLImageElement) => image.naturalWidth)).toBeGreaterThan(0);
});

test("分页保留 URL，刷新和返回后维持页码与 noindex", async ({ page }) => {
  await page.goto("/", { waitUntil: "domcontentloaded" });
  await page.locator(".pagination").getByText("2", { exact: true }).click();
  await expect(page).toHaveURL(/\?page=2$/);
  await expect(page.locator(".feed-title")).toContainText("第二页文章");
  await expect(page.locator('meta[name="robots"]')).toHaveAttribute("content", "noindex, follow");
  await page.reload({ waitUntil: "domcontentloaded" });
  await expect(page.locator(".feed-title")).toContainText("第二页文章");
  await page.goto("/?keyword=Go&page=2", { waitUntil: "domcontentloaded" });
  await expect(page.locator('link[rel="canonical"]')).toHaveAttribute("href", "http://127.0.0.1:4173/?keyword=Go&page=2");
});

test("后台按需加载 UI 组件，登录资料头像正确显示", async ({ page }) => {
  const errors: string[] = [];
  page.on("pageerror", error => errors.push(error.message));
  await page.addInitScript(() => {
    const payload = btoa(JSON.stringify({ exp: Math.floor(Date.now() / 1000) + 3600 }));
    localStorage.setItem("id_token", `test.${payload}.test`);
  });
  await page.goto("/admin/account/settings", { waitUntil: "domcontentloaded" });
  await expect(page.locator(".account-avatar-40")).toHaveAttribute("src", avatar);
  await expect(page.locator('input[name="avatar"]')).toBeAttached();
  // Other admin screens exercise template auto-registration after removing app.use(ElementPlus).
  for (const [path, selector] of [["/admin/category", ".category-page"], ["/admin/article/index", "table"], ["/admin/friendlink", ".el-table"], ["/admin/comment", "table"]]) {
    await page.goto(path, { waitUntil: "domcontentloaded" });
    if (path === "/admin/category") await expect(page.getByPlaceholder("搜索分类名称")).toBeVisible();
    else await expect(page.locator(selector).first()).toBeVisible();
  }
  expect(errors).toEqual([]);
});

test("保存资料后不刷新页面也能显示新头像", async ({ page }) => {
  await page.addInitScript(() => {
    const payload = btoa(JSON.stringify({ exp: Math.floor(Date.now() / 1000) + 3600 }));
    localStorage.setItem("id_token", `test.${payload}.test`);
  });
  await page.route("**/updated-author.svg", route => route.fulfill({ contentType: "image/svg+xml", body: '<svg xmlns="http://www.w3.org/2000/svg" width="40" height="40"><rect width="40" height="40" fill="blue"/></svg>' }));
  await page.route("**/blog/admin/user/profile", route => {
    const updated = route.request().method() === "PUT";
    return route.fulfill({ contentType: "application/json", body: JSON.stringify({ code: 0, data: { ...profile, pic: updated ? "/updated-author.svg" : avatar } }) });
  });
  await page.goto("/admin/account/settings", { waitUntil: "domcontentloaded" });
  await expect(page.locator(".account-avatar-40")).toHaveAttribute("src", avatar);
  await page.getByRole("button", { name: "保存更改", exact: true }).click();
  await expect(page.locator(".account-avatar-40")).toHaveAttribute("src", "/updated-author.svg");
  await page.getByRole("link", { name: "关于", exact: true }).click();
  await expect(page.locator(".about-avatar")).toHaveAttribute("src", "/updated-author.svg");
});

async function authenticate(page: Page) {
  await page.addInitScript(() => {
    const payload = btoa(JSON.stringify({ exp: Math.floor(Date.now() / 1000) + 3600 }));
    localStorage.setItem("id_token", `test.${payload}.test`);
  });
}

test("归档支持组合筛选和刷新保留条件", async ({ page }) => {
  await page.route("**/blog/base/getArchive", route => route.fulfill({ contentType: "application/json", body: JSON.stringify({ code: 0, data: { list: [
    { id: 1, fid: 1, title: "Go 入门", ctime: 1680000000 },
    { id: 2, fid: 1, title: "Vue 使用", ctime: 1710000000 },
    { id: 3, fid: 1, title: "Go 进阶", ctime: 1710000000 },
  ] } }) }));
  await page.goto("/archive?q=Go&year=2024", { waitUntil: "domcontentloaded" });
  await expect(page.getByRole("link", { name: "Go 进阶", exact: true })).toBeVisible();
  await expect(page.getByRole("link", { name: "Go 入门", exact: true })).toHaveCount(0);
  await expect(page.getByRole("link", { name: "Vue 使用", exact: true })).toHaveCount(0);
  await page.screenshot({ path: test.info().outputPath("archive.png"), fullPage: true });
  await page.getByLabel("搜索归档标题").fill("不存在");
  await expect(page).toHaveURL(/q=%E4%B8%8D%E5%AD%98%E5%9C%A8/);
  await page.reload({ waitUntil: "domcontentloaded" });
  await expect(page.getByLabel("搜索归档标题")).toHaveValue("不存在");
  await expect(page.getByText("没有符合条件的文章", { exact: false })).toBeVisible();
  await page.getByRole("button", { name: "清除筛选" }).click();
  await expect(page.getByRole("link", { name: "Vue 使用", exact: true })).toBeVisible();
});

test("文章显示阅读用时，剪贴板被拒绝时反馈错误", async ({ page }) => {
  await page.addInitScript(() => {
    Object.defineProperty(navigator, "clipboard", { value: { writeText: () => Promise.reject(new Error("permission denied")) }, configurable: true });
    document.execCommand = () => false;
  });
  const errors: string[] = [];
  page.on("pageerror", error => errors.push(error.message));
  await page.goto("/content/1/123", { waitUntil: "domcontentloaded" });
  await expect(page.getByTestId("reading-stats")).toContainText("分钟阅读");
  await page.screenshot({ path: test.info().outputPath("reading.png"), fullPage: true });
  await page.getByRole("button", { name: "复制文章链接" }).click();
  await expect(page.getByText("复制失败，请从地址栏复制链接")).toBeVisible();
  expect(errors).toEqual([]);
});

test("新建文章自动保存草稿，离开可取消，刷新后可恢复", async ({ page }) => {
  await authenticate(page);
  await page.goto("/admin/article/content", { waitUntil: "domcontentloaded" });
  await expect(page.getByRole("button", { name: "发布文章", exact: true })).toBeVisible();
  await page.getByPlaceholder("请输入文章标题").fill("未发布草稿");
  await page.locator('.w-e-text-container [contenteditable="true"]').fill("草稿正文不能丢失");
  await expect.poll(() => page.evaluate(() => JSON.parse(localStorage.getItem("feilong:article-draft:1:new") || "null")?.value.title)).toBe("未发布草稿");
  await page.getByRole("link", { name: "关于", exact: true }).click();
  await expect(page.getByText("离开编辑页面？", { exact: true })).toBeVisible();
  await page.getByRole("button", { name: "继续编辑", exact: true }).click();
  await expect(page).toHaveURL(/\/admin\/article\/content$/);
  page.on("dialog", dialog => dialog.accept());
  await page.reload({ waitUntil: "domcontentloaded" });
  await expect(page.getByRole("button", { name: "恢复草稿", exact: true })).toBeVisible();
  await page.getByRole("button", { name: "恢复草稿", exact: true }).click();
  await expect(page.getByPlaceholder("请输入文章标题")).toHaveValue("未发布草稿");
  await expect(page.locator('.w-e-text-container [contenteditable="true"]')).toContainText("草稿正文不能丢失");
  await page.screenshot({ path: test.info().outputPath("editor.png"), fullPage: true });
});

test("发布失败保留草稿，成功后清除且不重复创建", async ({ page }) => {
  await authenticate(page);
  let calls = 0;
  let saved: Record<string, unknown> = {};
  await page.route("**/blog/admin/article/createArticle", async route => {
    ++calls;
    if (calls === 1) return route.fulfill({ status: 500, contentType: "application/json", body: JSON.stringify({ msg: "模拟保存失败" }) });
    saved = route.request().postDataJSON();
    await new Promise(resolve => setTimeout(resolve, 250));
    return route.fulfill({ contentType: "application/json", body: JSON.stringify({ code: 0, data: { id: 456 } }) });
  });
  await page.route("**/blog/admin/article/getArticle", route => route.fulfill({ contentType: "application/json", body: JSON.stringify({ code: 0, data: { ...post, ...saved, ID: 456 } }) }));
  await page.goto("/admin/article/content", { waitUntil: "domcontentloaded" });
  await page.getByPlaceholder("请输入文章标题").fill("发布测试");
  await page.getByRole("button", { name: "发布文章", exact: true }).click();
  await expect(page.getByText("模拟保存失败", { exact: true })).toBeVisible();
  await expect.poll(() => page.evaluate(() => !!localStorage.getItem("feilong:article-draft:1:new"))).toBe(true);
  await page.getByRole("button", { name: "发布文章", exact: true }).evaluate((button: HTMLButtonElement) => { button.click(); button.click(); });
  await expect(page).toHaveURL(/\/admin\/article\/content\/456$/);
  await expect(page.getByRole("button", { name: "保存修改", exact: true })).toBeVisible();
  expect(calls).toBe(2);
  expect(saved.title).toBe("发布测试");
  await expect.poll(() => page.evaluate(() => localStorage.getItem("feilong:article-draft:1:new"))).toBeNull();
});

test("本机存储不足时明确提示，账号之间不恢复草稿", async ({ page }) => {
  await authenticate(page);
  await page.addInitScript(() => {
    localStorage.setItem("feilong:article-draft:2:new", JSON.stringify({ version: 1, savedAt: Date.now(), value: { title: "其他账号私有草稿", fid: 1, type: 1, contentFormat: "html", html: "", markdown: "", pic: "" } }));
    const original = Storage.prototype.setItem;
    Storage.prototype.setItem = function(key, value) {
      if (key.startsWith("feilong:article-draft:1:")) throw new DOMException("quota", "QuotaExceededError");
      return original.call(this, key, value);
    };
  });
  await page.goto("/admin/article/content", { waitUntil: "domcontentloaded" });
  await expect(page.getByRole("button", { name: "发布文章", exact: true })).toBeVisible();
  await expect(page.getByRole("button", { name: "恢复草稿", exact: true })).toHaveCount(0);
  await page.getByPlaceholder("请输入文章标题").fill("无法保存的草稿");
  await expect(page.getByText("本机草稿保存失败，请及时发布或备份内容。")).toBeVisible();
});

test("Markdown 编辑器使用同域资源初始化", async ({ page }) => {
  await authenticate(page);
  const external: string[] = [];
  page.on("request", request => { if (!request.url().startsWith("http://127.0.0.1:4173") && /^https?:/.test(request.url())) external.push(request.url()); });
  await page.goto("/admin/article/content", { waitUntil: "domcontentloaded" });
  await page.getByRole("button", { name: "Markdown", exact: true }).click();
  await expect(page.locator(".vditor-ir pre[contenteditable=true]")).toBeVisible({ timeout: 15000 });
  expect(external).toEqual([]);
});
