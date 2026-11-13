import { test, expect, type Page } from "@playwright/test";
import { mockAPI, post } from "./fixtures";

async function authenticate(page: Page) {
  await page.addInitScript(() => {
    const payload = btoa(
      JSON.stringify({ exp: Math.floor(Date.now() / 1000) + 3600 })
    );
    localStorage.setItem("id_token", `test.${payload}.test`);
  });
}
test.beforeEach(async ({ page }) => {
  await mockAPI(page);
});

test("新文章默认保存服务器草稿，标签摘要与 SEO 一并提交", async ({ page }) => {
  await authenticate(page);
  let saved: any;
  let creates = 0;
  await page.route("**/blog/admin/article/createArticle", (route) => {
    creates++;
    saved = route.request().postDataJSON();
    return route.fulfill({
      contentType: "application/json",
      body: JSON.stringify({
        code: 0,
        data: { id: 456, version: 1, status: 2 },
      }),
    });
  });
  await page.route("**/blog/admin/article/getArticle", (route) =>
    route.fulfill({
      contentType: "application/json",
      body: JSON.stringify({
        code: 0,
        data: { ...post, ...saved, ID: 456, version: 1, status: 2 },
      }),
    })
  );
  await page.goto("/admin/article/content", { waitUntil: "domcontentloaded" });
  await expect(page.getByLabel("发布状态")).toHaveValue("2");
  await page.getByPlaceholder("请输入文章标题").fill("保存在服务器的草稿");
  await page.getByLabel("列表摘要", { exact: true }).fill("手写摘要");
  await page.getByText("搜索引擎展示", { exact: true }).click();
  await page.getByLabel("SEO 标题", { exact: true }).fill("搜索展示标题");
  await page.getByLabel("SEO 描述", { exact: true }).fill("独立描述");
  await page.getByRole("button", { name: "保存草稿", exact: true }).click();
  await expect(page).toHaveURL(/content\/456$/);
  expect(creates).toBe(1);
  expect(saved.status).toBe(2);
  expect(saved.excerpt).toBe("手写摘要");
  expect(saved.seoTitle).toBe("搜索展示标题");
  expect(saved.seoDescription).toBe("独立描述");
});

test("定时发布校验未来时间并提交 Unix 时间戳", async ({ page }) => {
  await authenticate(page);
  let payload: any;
  await page.route("**/blog/admin/article/createArticle", (route) => {
    payload = route.request().postDataJSON();
    return route.fulfill({
      contentType: "application/json",
      body: JSON.stringify({ code: 0, data: { id: 456, version: 1 } }),
    });
  });
  await page.goto("/admin/article/content", { waitUntil: "domcontentloaded" });
  await page.getByPlaceholder("请输入文章标题").fill("定时文章");
  await page.getByLabel("发布状态").selectOption("3");
  await page.getByLabel("发布时间", { exact: true }).fill("2020-01-01T10:00");
  await page.getByRole("button", { name: "保存定时发布", exact: true }).click();
  await expect(
    page.getByText("请选择未来的发布时间", { exact: true })
  ).toBeVisible();
  expect(payload).toBeUndefined();
  const target = new Date(Date.now() + 86400000);
  const local = new Date(target.getTime() - target.getTimezoneOffset() * 60000)
    .toISOString()
    .slice(0, 16);
  await page.getByLabel("发布时间", { exact: true }).fill(local);
  await page.getByRole("button", { name: "保存定时发布", exact: true }).click();
  await expect.poll(() => payload?.status).toBe(3);
  expect(payload.publishAt).toBe(Math.floor(new Date(local).getTime() / 1000));
});

test("保存冲突保留草稿，重新载入后明确选择恢复", async ({ page }) => {
  await authenticate(page);
  let version = 2;
  await page.route("**/blog/admin/article/getArticle", (route) =>
    route.fulfill({
      contentType: "application/json",
      body: JSON.stringify({
        code: 0,
        data: {
          ...post,
          status: 2,
          version,
          contentFormat: "html",
          content: "<p>服务器内容</p>",
          title: version === 2 ? "原始草稿" : "服务器最新草稿",
        },
      }),
    })
  );
  await page.route("**/blog/admin/article/updateArticle", (route) =>
    route.fulfill({
      status: 409,
      contentType: "application/json",
      body: JSON.stringify({ code: 7, msg: "文章已更新，本机草稿仍然保留" }),
    })
  );
  await page.goto("/admin/article/content/123", {
    waitUntil: "domcontentloaded",
  });
  await page.getByPlaceholder("请输入文章标题").fill("本机未覆盖的修改");
  await page.getByRole("button", { name: "保存草稿", exact: true }).click();
  await expect(
    page.getByRole("button", { name: "载入服务器最新版本", exact: true })
  ).toBeVisible();
  await expect
    .poll(() =>
      page.evaluate(
        () =>
          JSON.parse(
            localStorage.getItem("feilong:article-draft:1:123") || "null"
          )?.value.title
      )
    )
    .toBe("本机未覆盖的修改");
  version = 3;
  await page
    .getByRole("button", { name: "载入服务器最新版本", exact: true })
    .click();
  await expect(page.getByPlaceholder("请输入文章标题")).toHaveValue(
    "服务器最新草稿"
  );
  await page.getByRole("button", { name: "恢复草稿", exact: true }).click();
  await expect(page.getByPlaceholder("请输入文章标题")).toHaveValue(
    "本机未覆盖的修改"
  );
});

test("历史版本只载入编辑器，不自动修改线上文章，预览过滤脚本", async ({
  page,
}) => {
  await authenticate(page);
  let writes = 0;
  let payload: any;
  const errors: string[] = [];
  page.on("pageerror", (error) => errors.push(error.message));
  await page.route("**/blog/admin/article/getArticle", (route) =>
    route.fulfill({
      contentType: "application/json",
      body: JSON.stringify({
        code: 0,
        data: {
          ...post,
          status: 0,
          version: 2,
          contentFormat: "html",
          content: "<p>当前内容</p>",
        },
      }),
    })
  );
  await page.route("**/blog/admin/article/revisions?*", (route) =>
    route.fulfill({
      contentType: "application/json",
      body: JSON.stringify({
        code: 0,
        data: [
          {
            id: 9,
            articleId: 123,
            version: 1,
            title: "旧版本标题",
            savedAt: 1700000000,
            reason: "create",
          },
        ],
      }),
    })
  );
  await page.route("**/blog/admin/article/revision?*", (route) =>
    route.fulfill({
      contentType: "application/json",
      body: JSON.stringify({
        code: 0,
        data: {
          title: "旧版本标题",
          content:
            '<p>旧版本正文</p><script>window.__historyXss=1</script><img src="x" onerror="window.__historyXss=1">',
          contentFormat: "html",
          fid: 1,
          type: 1,
          pic: "",
          excerpt: "旧摘要",
          seoTitle: "",
          seoDescription: "",
          tagNames: ["Go"],
        },
      }),
    })
  );
  await page.route("**/blog/admin/article/updateArticle", (route) => {
    writes++;
    payload = route.request().postDataJSON();
    return route.fulfill({
      contentType: "application/json",
      body: JSON.stringify({
        code: 0,
        data: { id: 123, version: 3, status: 0 },
      }),
    });
  });
  await page.goto("/admin/article/content/123", {
    waitUntil: "domcontentloaded",
  });
  await page.getByRole("button", { name: "历史版本", exact: true }).click();
  await page.getByRole("button", { name: /旧版本标题/ }).click();
  await expect(page.locator(".preview-content")).toContainText("旧版本正文");
  expect(
    await page.evaluate(() => (window as any).__historyXss)
  ).toBeUndefined();
  await page
    .getByRole("button", { name: "载入此版本到编辑器", exact: true })
    .click();
  await expect(page.getByPlaceholder("请输入文章标题")).toHaveValue(
    "旧版本标题"
  );
  expect(writes).toBe(0);
  expect(
    await page.evaluate(() => (window as any).__historyXss)
  ).toBeUndefined();
  await page.getByRole("button", { name: "保存修改", exact: true }).click();
  await expect.poll(() => writes).toBe(1);
  expect(payload.version).toBe(2);
  expect(payload.status).toBe(0);
  expect(payload.title).toBe("旧版本标题");
  expect(errors).toEqual([]);
});

test("撤下已发布文章需要确认，取消后不会写入", async ({ page }) => {
  await authenticate(page);
  let writes = 0;
  await page.route("**/blog/admin/article/updateArticle", (route) => {
    writes++;
    return route.fulfill({
      contentType: "application/json",
      body: JSON.stringify({ code: 0, data: { id: 123, version: 2 } }),
    });
  });
  await page.goto("/admin/article/content/123", {
    waitUntil: "domcontentloaded",
  });
  await page.getByLabel("发布状态").selectOption("2");
  await page.getByRole("button", { name: "保存草稿", exact: true }).click();
  await expect(
    page.getByText("暂时撤下已发布文章？", { exact: true })
  ).toBeVisible();
  await page.getByRole("button", { name: "继续编辑", exact: true }).click();
  expect(writes).toBe(0);
});

test("预览手机尺寸不发布文章", async ({ page }) => {
  await authenticate(page);
  let writes = 0;
  await page.route("**/blog/admin/article/createArticle", (route) => {
    writes++;
    return route.abort();
  });
  await page.goto("/admin/article/content", { waitUntil: "domcontentloaded" });
  await page.getByPlaceholder("请输入文章标题").fill("只预览不发布");
  await page.getByRole("button", { name: "预览文章", exact: true }).click();
  await expect(page.locator(".preview-body")).toContainText("只预览不发布");
  await page.getByRole("button", { name: "手机", exact: true }).click();
  await expect(page.locator(".preview-body")).toHaveClass(/preview-mobile/);
  expect(writes).toBe(0);
});

test("标签专题可搜索并进入对应文章列表", async ({ page }) => {
  await page.goto("/tags", { waitUntil: "domcontentloaded" });
  await expect(page.getByRole("link", { name: /# Go/ })).toBeVisible();
  await page.getByLabel("搜索标签", { exact: true }).fill("Rust");
  await expect(
    page.getByText("没有找到这个标签，试试其他关键词。")
  ).toBeVisible();
  await page.getByLabel("搜索标签", { exact: true }).fill("Go");
  await page.getByRole("link", { name: /# Go/ }).click();
  await expect(page).toHaveURL(/\/tag\/1$/);
  await expect(
    page.getByRole("heading", { name: "# Go", exact: true })
  ).toBeVisible();
  await expect(page.locator(".feed-title")).toContainText(post.title);
});

test("标签不存在时显示可恢复提示，不展示旧标签内容", async ({ page }) => {
  await page.route("**/blog/base/getTag?*", (route) =>
    route.fulfill({
      status: 404,
      contentType: "application/json",
      body: JSON.stringify({ code: 7, msg: "标签不存在" }),
    })
  );
  await page.goto("/tag/999", { waitUntil: "domcontentloaded" });
  await expect(
    page.getByText("标签不存在，或暂时没有公开文章。", { exact: false })
  ).toBeVisible();
  await expect(page.locator(".feed-title")).toHaveCount(0);
});

test("公开搜索范围可切换到正文并保持 URL", async ({ page }) => {
  const scopes: string[] = [];
  await page.route("**/blog/base/getArticleList?*", (route) => {
    scopes.push(
      new URL(route.request().url()).searchParams.get("scope") || "title"
    );
    return route.fulfill({
      contentType: "application/json",
      body: JSON.stringify({ code: 0, data: { list: [post], total: 1 } }),
    });
  });
  await page.goto("/?keyword=正文", { waitUntil: "domcontentloaded" });
  await page.getByLabel("搜索范围").selectOption("all");
  await expect(page).toHaveURL(/scope=all/);
  await expect.poll(() => scopes.includes("all")).toBe(true);
  await expect(page.locator('meta[name="robots"]')).toHaveAttribute(
    "content",
    "noindex, follow"
  );
});

test("稍后阅读可保存、刷新、移除，不上传个人标题和时间", async ({ page }) => {
  const bodies: any[] = [];
  await page.route("**/blog/base/getReadingArticles", (route) => {
    bodies.push(route.request().postDataJSON());
    return route.fulfill({
      contentType: "application/json",
      body: JSON.stringify({ code: 0, data: [post] }),
    });
  });
  await page.goto("/content/1/123", { waitUntil: "domcontentloaded" });
  await page.getByRole("button", { name: "稍后阅读", exact: true }).click();
  await expect(
    page.getByRole("button", { name: "已收藏", exact: true })
  ).toHaveAttribute("aria-pressed", "true");
  await page.getByRole("link", { name: "稍后阅读清单", exact: true }).click();
  await expect(
    page.getByRole("link", { name: post.title, exact: true })
  ).toBeVisible();
  await page.reload({ waitUntil: "domcontentloaded" });
  await expect(
    page.getByRole("link", { name: post.title, exact: true })
  ).toBeVisible();
  await expect.poll(() => bodies).toEqual([{ ids: [123] }, { ids: [123] }]);
  await page
    .getByRole("button", { name: `移除 ${post.title}`, exact: true })
    .click();
  await expect(
    page.getByText("给下一次阅读留个位置", { exact: true })
  ).toBeVisible();
});

test("下线的收藏不再提供公开链接", async ({ page }) => {
  await page.addInitScript(() =>
    localStorage.setItem(
      "feilong:reading-list:v1",
      JSON.stringify([
        { id: 123, fid: 1, title: "已下线文章", savedAt: Date.now() },
      ])
    )
  );
  await page.route("**/blog/base/getReadingArticles", (route) =>
    route.fulfill({
      contentType: "application/json",
      body: JSON.stringify({ code: 0, data: [] }),
    })
  );
  await page.goto("/reading-list", { waitUntil: "domcontentloaded" });
  await expect(page.getByText("文章已下线", { exact: true })).toBeVisible();
  await expect(page.locator('a[href="/content/1/123"]')).toHaveCount(0);
});

test("阅读清单存储失败不会假装收藏成功", async ({ page }) => {
  await page.addInitScript(() => {
    const original = Storage.prototype.setItem;
    Storage.prototype.setItem = function (key, value) {
      if (key === "feilong:reading-list:v1")
        throw new DOMException("quota", "QuotaExceededError");
      return original.call(this, key, value);
    };
  });
  await page.goto("/content/1/123", { waitUntil: "domcontentloaded" });
  await page.getByRole("button", { name: "稍后阅读", exact: true }).click();
  await expect(
    page.getByText("保存失败，请检查浏览器存储空间。", { exact: true })
  ).toBeVisible();
  await expect(
    page.getByRole("button", { name: "稍后阅读", exact: true })
  ).toHaveAttribute("aria-pressed", "false");
});

test("后台状态筛选写入 URL，置顶仅提交字段和版本", async ({ page }) => {
  await authenticate(page);
  const params: string[] = [];
  let patch: any;
  let featured = false;
  await page.route("**/blog/admin/article/getArticleList?*", (route) => {
    params.push(
      new URL(route.request().url()).searchParams.get("state") || "all"
    );
    return route.fulfill({
      contentType: "application/json",
      body: JSON.stringify({
        code: 0,
        data: { list: [{ ...post, isTop: featured ? 1 : 0 }], total: 1 },
      }),
    });
  });
  await page.route("**/blog/admin/article/state", (route) => {
    patch = route.request().postDataJSON();
    featured = true;
    return route.fulfill({
      contentType: "application/json",
      body: JSON.stringify({
        code: 0,
        data: { ...post, version: 2, isTop: 1 },
      }),
    });
  });
  await page.goto("/admin/article/index", { waitUntil: "domcontentloaded" });
  await page.getByRole("button", { name: "草稿", exact: true }).click();
  await expect(page).toHaveURL(/state=draft/);
  await expect.poll(() => params.includes("draft")).toBe(true);
  await page
    .getByRole("button", { name: `置顶 ${post.title}`, exact: true })
    .click();
  await expect.poll(() => patch?.isTop).toBe(1);
  expect(Object.keys(patch).sort()).toEqual(["id", "isTop", "version"]);
  await expect(
    page.getByRole("button", { name: `取消置顶 ${post.title}`, exact: true })
  ).toBeEnabled();
});

test("无效文章地址展示可恢复的缺页状态，并禁止索引", async ({ page }) => {
  await page.goto("/content/1/not-a-number", { waitUntil: "domcontentloaded" });
  await expect(page.getByText("文章暂不可访问", { exact: true })).toBeVisible();
  await expect(page.locator('meta[name="robots"]')).toHaveAttribute(
    "content",
    /noindex/
  );
  await expect(
    page.getByRole("link", { name: "返回首页", exact: true })
  ).toBeVisible();
});

test("快捷搜索支持单个中文字符并进入正文搜索", async ({ page }) => {
  await page.goto("/", { waitUntil: "domcontentloaded" });
  await expect(page.locator(".feed-title").first()).toBeVisible();
  await page.keyboard.press("Control+k");
  const input = page.getByPlaceholder("搜索标题、摘要或正文…");
  await expect(input).toBeFocused();
  await input.fill("云");
  await input.press("Enter");
  await expect(page).toHaveURL(/keyword=%E4%BA%91.*scope=all/);
  await expect(input).not.toBeVisible();
});

test("阅读清单同步其他标签页的增删", async ({ page, context }) => {
  await page.goto("/reading-list", { waitUntil: "domcontentloaded" });
  const other = await context.newPage();
  await mockAPI(other);
  await other.goto("/", { waitUntil: "domcontentloaded" });
  await other.evaluate(() =>
    localStorage.setItem(
      "feilong:reading-list:v1",
      JSON.stringify([
        { id: 123, fid: 1, title: "跨标签页收藏", savedAt: Date.now() },
      ])
    )
  );
  await expect(
    page.getByRole("link", { name: post.title, exact: true })
  ).toBeVisible();
  await other.evaluate(() =>
    localStorage.removeItem("feilong:reading-list:v1")
  );
  await expect(
    page.getByText("给下一次阅读留个位置", { exact: true })
  ).toBeVisible();
});
