# 本次修复与 Oracle ARM 服务器部署

本次修复基于 main 的 814a325，在云端完成代码修改与测试，尚未发布到线上。

本轮新增草稿、归档筛选、阅读工具与信息流优化，见 [第二轮迭代说明](iteration-2026-11.md)。

## 已完成

- 删除公开文章管理路由，后台接口统一验证 JWT 和管理员身份；未登录返回 401，非管理员返回 403。
- 公开文章详情、评论和点赞只允许访问已发布文章；后台编辑仍能读取隐藏文章。
- 服务端正文采用 bluemonday HTML 解析允许列表，过滤实体编码脚本链接、事件属性、SVG 等危险内容；Markdown 和 RSS 共用清洗逻辑，不再在 HTML 中间截断正文。
- 统一首页、文章、关于页、侧栏和账户菜单的作者头像，图片失效自动兜底；登录资料、公开资料和保存后的缓存统一处理。
- 新本地上传返回 /blog/... 根路径，旧本地图片路径兼容 API 代理；正文中的旧相对图片路径也会修正。保留原数据库内容，不进行批量替换。
- Nginx 示例使用 `location ^~ /blog/`，防止扩展名规则截走 API 图片；兼容历史 /public/uploads/file/ 地址。
- 数据库故障返回服务错误，真正缺页才返回 404；Nginx 服务故障回退保留 503，并附 noindex、Retry-After。
- 分页同步 URL；刷新和分享保留页码；服务端搜索和分页内容、canonical、robots 与前端保持一致。
- Markdown 摘要去除原始标记；Element Plus 改为按需加载；升级 npm 锁定依赖。本项目统一使用 npm/package-lock.json，移除重复 yarn.lock。
- robots 的具名搜索爬虫和通配组均排除后台/API/搜索/分页地址，同时允许公开内容并声明 Sitemap。

## 验证

后端执行 `go test ./...` 和 `go vet ./...`。SQL mock 回归覆盖公开管理路径移除、登录/非管理员拦截、管理员可达、隐藏文章隔离、脚本清洗绕过、数据库故障与缺页区分、分页和搜索直出。

前端执行 `npm ci`、`npm run build`、`npm audit`。浏览器回归使用模拟 API，覆盖公开页头像一致性、失败兜底、Markdown 摘要、分页刷新、后台组件和保存头像后的缓存同步。

```bash
cd web
npm ci
npx playwright install chromium
npm run build
npm run test:e2e
```

已有系统 Chromium 时可设置 `PLAYWRIGHT_CHROMIUM_EXECUTABLE_PATH=/usr/bin/chromium`。

Linux ARM64 程序已交叉编译并检查为静态链接 aarch64 ELF；云端机器不是 ARM，未直接运行该 ARM 程序。真实 MySQL/R2 上传集成和 Nginx 配置检测须在服务器验证。

## 部署包内容

- `blog-server-linux-arm64`：Oracle ARM（aarch64）后端程序。
- `web/`：生产前端静态文件，部署到现有 Nginx root 和 `seo.web-root` 所指目录。
- `nginx.conf.example`：配置参考；合并相应 location 到现有站点配置，保留证书、域名和现有代理参数。
- `source.patch`：相对于 814a325 的完整源码变更，包含新测试文件。
- `source.tar.gz`：源码快照，不包含私有配置、密钥、数据库业务数据或 node_modules。
- `SHA256SUMS`：包内产物校验值。

## 上传与替换

本机 `.ssh/config` 留在你的设备上。将下例 `YOUR_SSH_ALIAS` 换成其中的服务器别名：

```bash
scp feilongBlog-arm64-release.tar.gz YOUR_SSH_ALIAS:~/
ssh YOUR_SSH_ALIAS
uname -m
```

架构应为 `aarch64`。解压到独立目录并校验：

```bash
mkdir -p ~/feilongBlog-release
cd ~/feilongBlog-release
tar -xzf ~/feilongBlog-arm64-release.tar.gz
sha256sum -c SHA256SUMS
```

通过宝塔或现有运维方式记录后端程序路径、工作目录、启动命令、服务名称、Nginx root、`seo.web-root` 和真实本地上传目录。备份现有程序、静态目录、Nginx 配置和数据库。本次没有新增表/字段迁移；程序启动仍会执行项目已有的 AutoMigrate。

1. 前端：将包内 `web/` 的内容同步到原静态目录；保留旧 `/assets/` 文件一段时间，避免仍打开旧页面的访客在懒加载时请求已删除文件。HTML 与静态资源按完整构建成套发布。
2. 后端：停止现有进程后替换原程序，保留原工作目录、启动参数和配置文件。若原命令使用 `-c config.production.yaml`，继续保留。不要用示例配置覆盖真实配置。确认文件可执行，并用原宝塔进程管理或 systemd 服务启动。
3. Nginx：合并示例中的 `/blog/`、旧图片兼容和服务故障状态调整。先执行 `nginx -t`；通过后再用现有方式重载。在宝塔使用其配置检测和重载功能也可以。
4. 检查 `seo.web-root` 与 Nginx root 指向同一套新构建文件，确认页面响应的 JS/CSS 在当前 `/assets/` 中存在。
5. 若 Cloudflare 对 HTML 有强制缓存，清除旧 HTML/robots 缓存，再验证；保留哈希资源缓存。

本地图片即使修正了路由，也要求原文件仍存在于 `local.store-path`，且运行后端的用户可读取。

## Cloudflare 与百度：必须单独处理

仅部署代码无法覆盖 Cloudflare 对根路径 robots 的托管处理。当前线上根路径禁止 Baiduspider，而 `/blog/robots.txt` 返回的是应用允许抓取的版本。

1. 在 `teshh.com` 的 Cloudflare 控制台查找托管 robots.txt 设置，关闭覆盖行为或调整策略，使根路径最终返回应用规则。
2. 同时检查 Worker、转换规则、缓存和源站是否还有其他 robots 覆盖。最终响应不得出现 `User-agent: Baiduspider` 下的 `Disallow: /`，且应包含 Sitemap。
3. 在安全事件中按失败请求的 `cf-ray`/时间查询拦截记录。此次 Python 默认客户端返回 `403`、正文 `error code: 1010`。核实 Browser Integrity Check、WAF/Bot 等实际命中原因，仅调整需要的策略；不要只按可伪造的 User-Agent 全站放行。
4. 在百度搜索资源平台完成站点验证、抓取诊断和 Sitemap/文章 URL 提交。解除阻碍不等于立即收录。

```bash
curl -fsS https://teshh.com/robots.txt
curl -fsS https://teshh.com/blog/robots.txt
curl -fsS https://teshh.com/sitemap.xml
```

## 上线验收

```bash
# 公开管理地址应为 404；未登录后台管理应为 401
curl -sS -o /dev/null -w '%{http_code}\n' 'https://teshh.com/blog/article/getArticleList?page=1&pageSize=1'
curl -sS -o /dev/null -w '%{http_code}\n' 'https://teshh.com/blog/admin/article/getArticleList?page=1&pageSize=1'

# 公共页面应正常返回，真正不存在的地址应返回 404
curl -sS -o /dev/null -w '%{http_code}\n' https://teshh.com/about
curl -sS -o /dev/null -w '%{http_code}\n' https://teshh.com/archive
curl -sS -o /dev/null -w '%{http_code}\n' https://teshh.com/definitely-missing-page
```

浏览器确认首页、文章、关于、后台头像一致；测试一张新上传头像、保存后站内跳转、分页刷新、文章直接打开和旧 hash 链接。自动化访问优先等待 DOM 和关键内容，非关键第三方图片失败应单独记录，避免全部依赖 `networkidle`。

## 仍需原始文件或账号处理的项目

- 旧图片 `/public/uploads/file/2023-05-31/871dc746bad328d95d3966e14a0ba62a.png`：原路径和 API 路径此次线上检查均为 404。需要到服务器或备份中查找原图，恢复到本地上传目录，或迁移至 R2 后更新该记录。代码无法生成丢失的原文件。
- 第三方图片 `https://pic.smwl.top/temp/1808/bfc54b60f2182508.png` 曾返回 503。需取得原图后转存可控存储。
- 百度解除抓取限制与自动化客户端的 Cloudflare 拦截须在控制台完成。云端没有该账户管理连接。

## 回滚

若验收失败，停止新后端，恢复旧程序、原静态目录和原 Nginx 配置，再按原启动命令启动。检查配置语法、健康接口和站点可访问性。保留新旧上传文件与原配置，不覆盖上传目录或回滚数据库业务记录。

注意旧版本包含未鉴权文章管理漏洞；如暂时回滚，应先在 Nginx/Cloudflare 阻断 `/blog/article/` 和根路径 `/article/` 等公开管理入口，保留 `/blog/base/` 正常阅读接口，避免恢复暴露。
