import {
  test,
  expect,
  type APIRequestContext,
  type Page,
} from "@playwright/test";

let token = "";
let fid = 1;
const created: number[] = [];
const createdTags: number[] = [];
const stamp = () => `回归-${Date.now()}-${Math.floor(Math.random() * 1000)}`;
const privateHeaders = () => ({ "x-token": token });
async function create(
  request: APIRequestContext,
  overrides: Record<string, unknown> = {}
) {
  const title = stamp();
  const response = await request.post("/blog/admin/article/createArticle", {
    headers: privateHeaders(),
    data: {
      title,
      content: "<p>真实数据库文章正文</p>",
      contentFormat: "html",
      fid,
      type: 1,
      status: 2,
      tagNames: [],
      ...overrides,
    },
  });
  expect(response.status()).toBe(200);
  const payload = await response.json();
  expect(payload.code).toBe(0);
  created.push(payload.data.id);
  return payload.data as { id: number; version: number; status: number };
}
async function adminArticle(request: APIRequestContext, id: number) {
  const response = await request.post("/blog/admin/article/getArticle", {
    headers: privateHeaders(),
    data: { id },
  });
  const payload = await response.json();
  expect(payload.code).toBe(0);
  return payload.data;
}
async function signIn(page: Page) {
  await page.addInitScript(
    (value) => localStorage.setItem("id_token", value),
    token
  );
}

test.beforeAll(async ({ request, baseURL }) => {
  // A forwarded production site must never be used as this fixture accidentally.
  const root = await request.get("/");
  expect(await root.text()).toContain(`rel="canonical" href="${baseURL}/"`);
  const login = await request.post("/blog/base/login", {
    data: { username: "admin", password: "123456" },
  });
  const payload = await login.json();
  expect(payload.code).toBe(0);
  token = payload.data.token;
  const categories = await (
    await request.get("/blog/base/getCategoryList")
  ).json();
  fid = (
    categories.data.find((item: any) => item.fid !== 0) || categories.data[0]
  ).ID;
});
test.afterEach(async ({ request }) => {
  if (created.length)
    await request.delete("/blog/admin/article/delArticleByIds", {
      headers: privateHeaders(),
      data: { ids: created.splice(0) },
    });
  for (const id of createdTags.splice(0))
    await request.delete("/blog/admin/article/tags", {
      headers: privateHeaders(),
      data: { id },
    });
});

test("Nginx 直达页面、真实静态脚本、robots 与缺页状态正确", async ({
  request,
}) => {
  const response = await request.get("/");
  expect(response.status()).toBe(200);
  const html = await response.text();
  const script = html.match(/src="(\/assets\/[^\"]+\.js)"/);
  expect(script).toBeTruthy();
  const resource = await request.get(script![1]);
  expect(resource.status()).toBe(200);
  expect(resource.headers()["content-type"]).toMatch(/javascript/);
  for (const path of ["/about", "/archive", "/tags", "/reading-list"]) {
    expect((await request.get(path)).status()).toBe(200);
  }
  expect((await request.get("/missing-fixture-page")).status()).toBe(404);
  expect((await request.get("/tag/999999999")).status()).toBe(404);
  const robots = await (await request.get("/robots.txt")).text();
  expect(robots).toContain("User-agent: Baiduspider\nAllow: /");
  expect(robots).toContain("/sitemap.xml");
});

test("真实写作流程：服务器草稿到公开专题、SEO 与稍后阅读", async ({
  page,
  request,
}) => {
  await signIn(page);
  const title = stamp();
  const tagName = `专题${Date.now()}`;
  await page.goto("/admin/article/content", { waitUntil: "domcontentloaded" });
  await page.getByPlaceholder("请输入文章标题").fill(title);
  await page
    .locator('.w-e-text-container [contenteditable="true"]')
    .fill("真实写作正文，确认数据库与页面保持一致。");
  await page
    .getByLabel("列表摘要", { exact: true })
    .fill("作者自定义的文章摘要");
  await page.locator(".el-select__input").fill(tagName);
  await page.getByRole("option", { name: tagName, exact: true }).click();
  await page.getByPlaceholder("请输入文章标题").click();
  await page.getByText("搜索引擎展示", { exact: true }).click();
  await page.getByLabel("SEO 标题", { exact: true }).fill(`${title} 搜索标题`);
  await page.getByLabel("SEO 描述", { exact: true }).fill("独立 SEO 描述");
  await page.getByRole("button", { name: "保存草稿", exact: true }).click();
  await expect(page).toHaveURL(/\/admin\/article\/content\/\d+$/);
  const id = Number(new URL(page.url()).pathname.split("/").pop());
  created.push(id);
  const draft = await adminArticle(request, id);
  expect(draft.status).toBe(2);
  expect(draft.version).toBe(1);
  expect(draft.tagNames).toContain(tagName);
  for (const tag of draft.tags) createdTags.push(tag.id);
  expect(
    (await request.post("/blog/base/getArticle", { data: { id } })).status()
  ).toBe(404);
  expect((await request.get(`/content/${draft.fid}/${id}`)).status()).toBe(404);
  expect(await (await request.get("/sitemap.xml")).text()).not.toContain(
    `/content/${draft.fid}/${id}`
  );
  await page.getByLabel("发布状态").selectOption("0");
  await page.getByRole("button", { name: "发布文章", exact: true }).click();
  await expect(
    page.getByRole("button", { name: "保存修改", exact: true })
  ).toBeEnabled();
  const publicHTML = await request.get(`/content/${draft.fid}/${id}`);
  expect(publicHTML.status()).toBe(200);
  expect(await publicHTML.text()).toContain(`<title>${title} 搜索标题</title>`);
  await page.goto(`/content/${draft.fid}/${id}`, {
    waitUntil: "domcontentloaded",
  });
  await expect(
    page.getByRole("heading", { name: title, exact: true })
  ).toBeVisible();
  await expect(page.locator('meta[name="description"]')).toHaveAttribute(
    "content",
    "独立 SEO 描述"
  );
  await page.getByRole("button", { name: "稍后阅读", exact: true }).click();
  await page.getByRole("link", { name: "稍后阅读清单", exact: true }).click();
  await expect(
    page.getByRole("link", { name: title, exact: true })
  ).toBeVisible();
  const afterReading = await adminArticle(request, id);
  expect(afterReading.view).toBe(1);
  await page.goto(`/tag/${draft.tags[0].id}`, {
    waitUntil: "domcontentloaded",
  });
  await expect(
    page.getByRole("heading", { name: `# ${tagName}`, exact: true })
  ).toBeVisible();
  await expect(page.locator(".feed-title")).toContainText(title);
  const revisions = await (
    await request.get(`/blog/admin/article/revisions?id=${id}`, {
      headers: privateHeaders(),
    })
  ).json();
  expect(revisions.data.length).toBe(2);
  await page.screenshot({
    path: test.info().outputPath("topic-real.png"),
    fullPage: true,
  });
});

test("真实保存冲突保护本机内容，恢复后生成新历史版本", async ({
  page,
  request,
}) => {
  const record = await create(request);
  await signIn(page);
  await page.goto(`/admin/article/content/${record.id}`, {
    waitUntil: "domcontentloaded",
  });
  const original = await adminArticle(request, record.id);
  await expect(page.getByPlaceholder("请输入文章标题")).toHaveValue(
    original.title
  );
  const update = await request.put("/blog/admin/article/updateArticle", {
    headers: privateHeaders(),
    data: { ...original, title: "另一页面的更新" },
  });
  expect(update.status()).toBe(200);
  await page.getByPlaceholder("请输入文章标题").fill("当前页面保留的内容");
  await page.getByRole("button", { name: "保存草稿", exact: true }).click();
  await expect(
    page.getByRole("button", { name: "载入服务器最新版本", exact: true })
  ).toBeVisible();
  await page
    .getByRole("button", { name: "载入服务器最新版本", exact: true })
    .click();
  await expect(page.getByPlaceholder("请输入文章标题")).toHaveValue(
    "另一页面的更新"
  );
  await page.getByRole("button", { name: "恢复草稿", exact: true }).click();
  await page.getByRole("button", { name: "保存草稿", exact: true }).click();
  await expect(
    page.getByRole("button", { name: "保存草稿", exact: true })
  ).toBeEnabled();
  await expect
    .poll(async () => (await adminArticle(request, record.id)).title)
    .toBe("当前页面保留的内容");
  const final = await adminArticle(request, record.id);
  expect(final.version).toBe(3);
});

test("真实定时任务自动发布且刷新 RSS 和 Sitemap", async ({ request }) => {
  const when = Math.floor(Date.now() / 1000) + 2;
  const title = stamp();
  const record = await create(request, { title, status: 3, publishAt: when });
  expect((await request.get(`/content/${fid}/${record.id}`)).status()).toBe(
    404
  );
  expect(await (await request.get("/rss.xml")).text()).not.toContain(title);
  await expect
    .poll(
      async () =>
        (
          await request.post("/blog/base/getArticle", {
            data: { id: record.id },
          })
        ).status(),
      { timeout: 45000, intervals: [1000] }
    )
    .toBe(200);
  const published = await adminArticle(request, record.id);
  expect(published.status).toBe(0);
  expect(published.version).toBe(2);
  expect(published.ctime).toBe(when);
  expect(await (await request.get("/rss.xml")).text()).toContain(title);
  expect(await (await request.get("/sitemap.xml")).text()).toContain(
    `/content/${fid}/${record.id}`
  );
});

test("真实上传经 Nginx API 代理读取，同名文件不覆盖且类型正确", async ({
  request,
}) => {
  const image = Buffer.from(
    "iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAQAAAC1HAwCAAAAC0lEQVR42mP8/x8AAwMCAO+aDBkAAAAASUVORK5CYII=",
    "base64"
  );
  const urls: string[] = [];
  for (let i = 0; i < 2; i++) {
    const response = await request.post("/blog/admin/file/upload", {
      headers: privateHeaders(),
      multipart: {
        file: { name: "same-name.html", mimeType: "image/png", buffer: image },
      },
    });
    const payload = await response.json();
    expect(payload.code).toBe(0);
    urls.push(payload.data.url);
    const asset = await request.get(payload.data.url);
    expect(asset.status()).toBe(200);
    expect(asset.headers()["content-type"]).toBe("image/png");
    expect(await asset.body()).toEqual(image);
  }
  expect(urls[0]).not.toBe(urls[1]);
  expect(urls[0]).toMatch(/\.png$/);
  const legacy = urls[0].replace(
    "/blog/uploads/file/",
    "/public/uploads/file/"
  );
  expect((await request.get(legacy)).status()).toBe(200);
});

test("匿名访客不能读取版本、服务器草稿或写入管理接口", async ({ request }) => {
  for (const path of [
    "/blog/admin/article/revisions?id=1",
    "/blog/admin/article/revision?id=1&revisionId=1",
    "/blog/admin/article/tags",
  ]) {
    expect((await request.get(path)).status()).toBe(401);
  }
  expect(
    (
      await request.post("/blog/article/createArticle", {
        data: { title: "不得写入" },
      })
    ).status()
  ).toBe(404);
  const record = await create(request, { title: "不可公开的私有草稿" });
  expect(
    (
      await request.post("/blog/base/getArticle", { data: { id: record.id } })
    ).status()
  ).toBe(404);
  const reading = await (
    await request.post("/blog/base/getReadingArticles", {
      data: { ids: [record.id] },
    })
  ).json();
  expect(reading.data).toEqual([]);
});

test("手机编辑器先展示正文并固定提供保存入口", async ({ page }) => {
  await signIn(page);
  await page.setViewportSize({ width: 390, height: 844 });
  await page.goto("/admin/article/content", { waitUntil: "domcontentloaded" });
  await expect(page.getByPlaceholder("请输入文章标题")).toBeInViewport();
  await expect(page.locator(".mobile-publish-bar")).toBeInViewport();
  await expect(page.locator(".mobile-publish-bar")).toContainText(
    "尚未保存到服务器"
  );
  await page.screenshot({
    path: test.info().outputPath("editor-mobile-real.png"),
    fullPage: true,
  });
});

test("标签管理真实创建、重命名和移除，草稿标签不公开", async ({
  page,
  request,
}) => {
  await signIn(page);
  const originalName = `标签-${Date.now()}`;
  await page.goto("/admin/tags", { waitUntil: "domcontentloaded" });
  await page.getByRole("button", { name: "新建标签", exact: true }).click();
  await page.getByLabel("标签名称", { exact: true }).fill(originalName);
  await page.getByLabel("标签描述", { exact: true }).fill("后台管理的专题说明");
  await page.getByRole("button", { name: "保存标签", exact: true }).click();
  const originalRow = page.getByRole("row").filter({ hasText: originalName });
  await expect(originalRow).toContainText("后台管理的专题说明");
  const all = await (
    await request.get("/blog/admin/article/tags", { headers: privateHeaders() })
  ).json();
  const tag = all.data.find((item: any) => item.name === originalName);
  expect(tag.id).toBeGreaterThan(0);
  createdTags.push(tag.id);
  const record = await create(request, { tagNames: [originalName] });
  expect((await request.get(`/tag/${tag.id}`)).status()).toBe(404);
  await originalRow.getByRole("button", { name: "编辑", exact: true }).click();
  await page.getByLabel("标签名称", { exact: true }).fill(`${originalName}-新`);
  await page.getByRole("button", { name: "保存标签", exact: true }).click();
  const renamedRow = page
    .getByRole("row")
    .filter({ hasText: `${originalName}-新` });
  await expect(renamedRow).toBeVisible();
  expect((await adminArticle(request, record.id)).tagNames).toEqual([
    `${originalName}-新`,
  ]);
  await renamedRow.getByRole("button", { name: "移除", exact: true }).click();
  await page.getByRole("button", { name: "取消", exact: true }).click();
  await expect(renamedRow).toBeVisible();
  await renamedRow.getByRole("button", { name: "移除", exact: true }).click();
  const deletion = page.waitForResponse(
    (response) =>
      response.url().endsWith("/admin/article/tags") &&
      response.request().method() === "DELETE"
  );
  await page.getByRole("button", { name: "移除标签", exact: true }).click();
  expect((await deletion).status()).toBe(200);
  await expect(renamedRow).toHaveCount(0);
  const after = await adminArticle(request, record.id);
  expect(after.title).toBeTruthy();
  expect(after.tagNames).toEqual([]);
});
