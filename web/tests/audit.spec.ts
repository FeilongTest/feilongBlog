import { test, expect, type Page } from "@playwright/test";
import { mockAPI, post } from "./fixtures";
async function writer(page: Page) {
  await page.addInitScript(() =>
    localStorage.setItem(
      "id_token",
      `test.${btoa(
        JSON.stringify({ exp: Math.floor(Date.now() / 1000) + 3600 })
      )}.test`
    )
  );
  await mockAPI(page);
}
test("Markdown 快速输入后立即保存，首个请求必须含最新正文", async ({
  page,
}) => {
  await writer(page);
  let saved: any;
  await page.route("**/blog/admin/article/createArticle", async (route) => {
    saved = route.request().postDataJSON();
    await route.fulfill({
      json: { code: 0, data: { id: 456, version: 1, status: 2 } },
    });
  });
  await page.route("**/blog/admin/article/getArticle", (route) =>
    route.fulfill({
      json: {
        code: 0,
        data: { ...post, ...saved, ID: 456, version: 1, status: 2 },
      },
    })
  );
  await page.goto("/admin/article/content");
  await page.getByPlaceholder("请输入文章标题").fill("快速 Markdown 草稿");
  await page.getByRole("button", { name: "Markdown", exact: true }).click();
  const editor = page.locator(".vditor-ir pre[contenteditable=true]");
  await expect(editor).toBeVisible({ timeout: 15000 });
  await editor.fill("快速输入后不等待防抖就保存");
  await page.getByRole("button", { name: "保存草稿", exact: true }).click();
  await expect
    .poll(() => saved?.content)
    .toContain("快速输入后不等待防抖就保存");
  await expect(page).toHaveURL(/content\/456$/);
  await expect(page.getByText("离开编辑页面？", { exact: true })).toHaveCount(
    0
  );
});

async function markdownDraft(page: Page) {
  await writer(page);
  await page.goto("/admin/article/content");
  await page.getByLabel("文章标题", { exact: true }).fill("Markdown 即时操作");
  await page.getByRole("button", { name: "Markdown", exact: true }).click();
  const editor = page.locator(".vditor-ir pre[contenteditable=true]");
  await expect(editor).toBeVisible({ timeout: 15000 });
  return editor;
}

test("Markdown 输入后立即预览，不依赖防抖或保存", async ({ page }) => {
  const editor = await markdownDraft(page);
  await editor.fill("尚未保存的即时预览正文");
  await page.getByRole("button", { name: "预览文章", exact: true }).click();
  await expect(page.locator(".el-dialog .preview-content")).toContainText(
    "尚未保存的即时预览正文"
  );
});

test("Markdown 首次保存失败，本机草稿仍保留最新正文", async ({ page }) => {
  const editor = await markdownDraft(page);
  await page.route("**/blog/admin/article/createArticle", (route) =>
    route.fulfill({
      status: 503,
      json: { code: 7, msg: "隔离测试：保存不可用" },
    })
  );
  await editor.fill("失败时也不能丢失的正文");
  await page.getByRole("button", { name: "保存草稿", exact: true }).click();
  await expect(
    page.getByText("隔离测试：保存不可用", { exact: true })
  ).toBeVisible();
  await expect
    .poll(() =>
      page.evaluate(
        () =>
          JSON.parse(
            localStorage.getItem("feilong:article-draft:1:new") || "null"
          )?.value.markdown
      )
    )
    .toContain("失败时也不能丢失的正文");
  await expect(editor).toContainText("失败时也不能丢失的正文");
});

test("Markdown 快速输入后离开，提示并保留正文；取消继续编辑", async ({
  page,
}) => {
  const editor = await markdownDraft(page);
  await editor.fill("准备离开时的最后一行");
  await page.getByRole("link", { name: "查看文章", exact: true }).click();
  await expect(page.getByText("离开编辑页面？", { exact: true })).toBeVisible();
  await expect
    .poll(() =>
      page.evaluate(
        () =>
          JSON.parse(
            localStorage.getItem("feilong:article-draft:1:new") || "null"
          )?.value.markdown
      )
    )
    .toContain("准备离开时的最后一行");
  await page.getByRole("button", { name: "继续编辑", exact: true }).click();
  await expect(editor).toContainText("准备离开时的最后一行");
});

test("标签 API 网络失败可重试，不能误报不存在或空结果", async ({ page }) => {
  await mockAPI(page);
  let failed = true;
  await page.route("**/blog/base/getTag?*", (route) =>
    failed
      ? route.fulfill({
          status: 503,
          json: { code: 7, msg: "隔离测试：服务暂不可用" },
        })
      : route.fulfill({
          json: {
            code: 0,
            data: { id: 1, name: "Go", description: "Go 主题", count: 1 },
          },
        })
  );
  await page.goto("/tag/1");
  await expect(
    page.getByRole("button", { name: "重新加载", exact: true })
  ).toBeVisible();
  await expect(
    page.getByText("标签不存在，或暂时没有公开文章。", { exact: false })
  ).toHaveCount(0);
  await expect(page.getByText("没有找到相关文章", { exact: true })).toHaveCount(
    0
  );
  failed = false;
  await page.getByRole("button", { name: "重新加载", exact: true }).click();
  await expect(page.locator(".feed-title")).toContainText(post.title);
});

test("专题有效但文章 API 返回 404，不能说标签不存在", async ({ page }) => {
  await mockAPI(page);
  await page.route("**/blog/base/getArticleList?*", (route) =>
    route.fulfill({
      status: 404,
      json: { code: 7, msg: "隔离测试：列表端点不可用" },
    })
  );
  await page.goto("/tag/1");
  await expect(
    page.getByRole("button", { name: "重新加载", exact: true })
  ).toBeVisible();
  await expect(
    page.getByText("标签不存在，或暂时没有公开文章。", { exact: false })
  ).toHaveCount(0);
});

test("分页支持键盘，提供页码和边界状态并保留搜索条件", async ({ page }) => {
  await mockAPI(page);
  await page.goto("/?keyword=Go&scope=all");
  const pagination = page.getByRole("navigation", {
    name: "文章分页",
    exact: true,
  });
  await expect(
    pagination.getByRole("button", { name: "上一页", exact: true })
  ).toBeDisabled();
  const next = pagination.getByRole("button", { name: "下一页", exact: true });
  await next.focus();
  await page.keyboard.press("Enter");
  await expect(page).toHaveURL(/keyword=Go&scope=all&page=2/);
  await expect(
    pagination.getByRole("button", { name: "第 2 页", exact: true })
  ).toHaveAttribute("aria-current", "page");
  await expect(
    pagination.getByRole("button", { name: "下一页", exact: true })
  ).toBeDisabled();
});

test("手机导航可用键盘开关并同步展开状态，表单标签关联控件", async ({
  page,
}) => {
  await page.setViewportSize({ width: 390, height: 844 });
  await writer(page);
  await page.goto("/admin/article/content");
  await expect(page.getByLabel("文章标题", { exact: true })).toBeVisible();
  await expect(page.getByLabel("文章分类", { exact: true })).toBeVisible();
  await expect(page.getByLabel("文章样式", { exact: true })).toBeVisible();
  const toggle = page.getByRole("button", { name: "导航菜单", exact: true });
  await toggle.focus();
  await page.keyboard.press("Space");
  await expect(toggle).toHaveAttribute("aria-expanded", "true");
  await expect(page.locator("#kt_aside")).toHaveClass(/drawer-on/);
  await page.locator(".drawer-overlay").click({ position: { x: 380, y: 80 } });
  await expect(toggle).toHaveAttribute("aria-expanded", "false");
});

test("手机文章封面保留比例限高，浮条不遮挡，页脚仍可打开友链", async ({
  page,
}) => {
  await page.setViewportSize({ width: 390, height: 844 });
  await mockAPI(page);
  await page.route("**/cover-audit.svg", (route) =>
    route.fulfill({
      contentType: "image/svg+xml",
      body: '<svg xmlns="http://www.w3.org/2000/svg" width="1200" height="600"><rect width="1200" height="600" fill="green"/></svg>',
    })
  );
  await page.route("**/blog/base/getArticle", (route) =>
    route.fulfill({
      json: { code: 0, data: { ...post, pic: "/cover-audit.svg" } },
    })
  );
  await page.goto("/content/1/123");
  const image = page.locator(".article-cover img");
  await expect(image).toBeVisible();
  const dimensions = await image.evaluate((el) => ({
    height: el.getBoundingClientRect().height,
    width: el.getBoundingClientRect().width,
    fit: getComputedStyle(el).objectFit,
  }));
  expect(dimensions.height).toBeLessThanOrEqual(280);
  expect(dimensions.width / dimensions.height).toBeCloseTo(2, 1);
  expect(dimensions.fit).toBe("contain");
  await expect(page.locator(".engage-toolbar")).not.toBeVisible();
  await page.locator(".article-cover").click();
  await expect(page.locator(".image-lightbox img")).toBeVisible();
  await page.keyboard.press("Escape");
  await page.getByRole("button", { name: "友情链接", exact: true }).click();
  await expect(page.locator("#kt_engage_demos_label")).toHaveClass(/drawer-on/);
  expect(
    await page.evaluate(
      () => document.documentElement.scrollWidth <= innerWidth
    )
  ).toBe(true);
});

test("文章侧栏与顶部搜索统一正文范围", async ({ page }) => {
  await mockAPI(page);
  await page.goto("/content/1/123");
  await expect(
    page.getByRole("heading", { name: post.title, exact: true })
  ).toBeVisible();
  const input = page.locator(".card").getByPlaceholder("搜索标题、摘要或正文…");
  await input.fill("副作用");
  await input.press("Enter");
  await expect(page).toHaveURL(/keyword=%E5%89%AF%E4%BD%9C%E7%94%A8&scope=all/);
});

test("有文章的分类删除被拒绝时保留分类，没有未处理异常", async ({ page }) => {
  await writer(page);
  const errors: string[] = [];
  page.on("pageerror", (error) => errors.push(error.message));
  await page.route("**/blog/admin/category/delCategory", (route) =>
    route.fulfill({
      status: 400,
      json: { code: 7, msg: "分类仍有文章，请先迁移文章" },
    })
  );
  await page.goto("/admin/category");
  const row = page.locator(".root-row");
  await expect(row).toContainText("测试分类");
  await row
    .getByRole("button", { name: "删除分类 测试分类", exact: true })
    .click();
  await page.getByRole("button", { name: "确认删除", exact: true }).click();
  await expect(
    page.getByText("分类仍有文章，请先迁移文章", { exact: true })
  ).toBeVisible();
  await expect(row).toContainText("测试分类");
  expect(errors).toEqual([]);
});

test("分类加载失败不冒充没有分类，可重试恢复", async ({ page }) => {
  await writer(page);
  let fail = true;
  await page.route("**/blog/admin/category/getCategoryList", (route) =>
    fail
      ? route.fulfill({ status: 503, json: { code: 7, msg: "分类暂时不可用" } })
      : route.fulfill({
          json: { code: 0, data: [{ ID: 1, fid: 0, name: "恢复后的分类" }] },
        })
  );
  await page.goto("/admin/category");
  await expect(
    page.getByRole("button", { name: "重新加载分类", exact: true })
  ).toBeVisible();
  await expect(page.getByText("没有找到分类", { exact: true })).toHaveCount(0);
  fail = false;
  await page.getByRole("button", { name: "重新加载分类", exact: true }).click();
  await expect(page.locator(".category-name")).toContainText("恢复后的分类");
});

test("Markdown 图片上传未完成时禁止保存，完成后保存包含图片引用", async ({
  page,
}) => {
  const editor = await markdownDraft(page);
  let release!: () => void;
  const gate = new Promise<void>((resolve) => {
    release = resolve;
  });
  await page.route("**/blog/admin/file/upload", async (route) => {
    await gate;
    await route.fulfill({
      json: { code: 0, data: { url: "/audit-upload.png" } },
    });
  });
  let payload: any;
  await page.route("**/blog/admin/article/createArticle", (route) => {
    payload = route.request().postDataJSON();
    return route.fulfill({
      json: { code: 0, data: { id: 456, version: 1, status: 2 } },
    });
  });
  await page.route("**/blog/admin/article/getArticle", (route) =>
    route.fulfill({
      json: {
        code: 0,
        data: { ...post, ...payload, ID: 456, status: 2, version: 1 },
      },
    })
  );
  const uploading = page.waitForRequest((request) =>
    request.url().endsWith("/admin/file/upload")
  );
  await page.locator(".markdown-editor input[type=file]").setInputFiles({
    name: "audit.png",
    mimeType: "image/png",
    buffer: Buffer.from(
      "iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAQAAAC1HAwCAAAAC0lEQVR42mP8/x8AAwMCAO+aWQAAAABJRU5ErkJggg==",
      "base64"
    ),
  });
  await uploading;
  try {
    await expect(
      page.getByRole("button", { name: "图片上传中…", exact: true })
    ).toBeDisabled();
  } finally {
    release();
  }
  await expect(
    page.getByRole("button", { name: "保存草稿", exact: true })
  ).toBeEnabled();
  await page.getByRole("button", { name: "保存草稿", exact: true }).click();
  await expect.poll(() => payload?.content).toContain("/audit-upload.png");
});

test("富文本选择图片上传时禁止保存，完成后正文包含图片", async ({ page }) => {
  await writer(page);
  await page.goto("/admin/article/content");
  const editor = page.locator('.w-e-text-container [contenteditable="true"]');
  await expect(editor).toBeVisible();
  let release!: () => void;
  const gate = new Promise<void>((resolve) => {
    release = resolve;
  });
  await page.route("**/blog/admin/file/upload", async (route) => {
    await gate;
    await route.fulfill({
      json: { code: 0, data: { url: "/audit-upload.png" } },
    });
  });
  const uploading = page.waitForRequest((request) =>
    request.url().endsWith("/admin/file/upload")
  );
  await editor.locator("p").first().click();
  await page.keyboard.type("Rich text with image");
  await page.locator('button[data-menu-key="group-image"]').click();
  await expect(
    page.locator('button[data-menu-key="uploadImage"]')
  ).not.toHaveClass(/disabled/);
  const [chooser] = await Promise.all([
    page.waitForEvent("filechooser"),
    page.locator('button[data-menu-key="uploadImage"]').click(),
  ]);
  await chooser.setFiles({
    name: "audit.png",
    mimeType: "image/png",
    buffer: Buffer.from(
      "iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAQAAAC1HAwCAAAAC0lEQVR42mP8/x8AAwMCAO+aWQAAAABJRU5ErkJggg==",
      "base64"
    ),
  });
  await uploading;
  try {
    await expect(
      page.getByRole("button", { name: "图片上传中…", exact: true })
    ).toBeDisabled();
  } finally {
    release();
  }
  await expect(
    page.getByRole("button", { name: "保存草稿", exact: true })
  ).toBeEnabled();
  await expect(editor.locator("img")).toHaveAttribute(
    "src",
    /audit-upload.png/
  );
});

test("Markdown 图片上传失败释放保存按钮并保留正文", async ({ page }) => {
  const editor = await markdownDraft(page);
  await editor.fill("失败后仍保留的正文");
  const errors: string[] = [];
  page.on("pageerror", (error) => errors.push(error.message));
  await page.route("**/blog/admin/file/upload", (route) =>
    route.fulfill({ status: 503, json: { code: 7, msg: "隔离测试：上传失败" } })
  );
  await page.locator(".markdown-editor input[type=file]").setInputFiles({
    name: "audit.png",
    mimeType: "image/png",
    buffer: Buffer.from(
      "iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAQAAAC1HAwCAAAAC0lEQVR42mP8/x8AAwMCAO+aWQAAAABJRU5ErkJggg==",
      "base64"
    ),
  });
  await expect(
    page.getByText("隔离测试：上传失败", { exact: true })
  ).toBeVisible();
  await expect(
    page.getByRole("button", { name: "保存草稿", exact: true })
  ).toBeEnabled();
  await expect(editor).toContainText("失败后仍保留的正文");
  expect(errors).toEqual([]);
});

test("Markdown 超大图片在发请求前提示，保存不被卡住", async ({ page }) => {
  await markdownDraft(page);
  let uploads = 0;
  await page.route("**/blog/admin/file/upload", (route) => {
    uploads++;
    return route.fulfill({
      json: { code: 0, data: { url: "/unexpected.png" } },
    });
  });
  await page
    .locator(".markdown-editor input[type=file]")
    .setInputFiles({
      name: "too-large.png",
      mimeType: "image/png",
      buffer: Buffer.alloc(10 * 1024 * 1024 + 1),
    });
  await expect(
    page.getByText("单张图片不能超过 10 MB", { exact: true })
  ).toBeVisible();
  await expect(
    page.getByRole("button", { name: "保存草稿", exact: true })
  ).toBeEnabled();
  expect(uploads).toBe(0);
});
