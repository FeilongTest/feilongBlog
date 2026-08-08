# Feilong Blog

一个前后端分离的个人博客系统，提供文章信息流、富文本内容、评论互动和后台内容管理。前端使用 Vue 3 与 Vite，后端使用 Gin 与 MySQL，并支持本地文件存储和 Cloudflare R2 对象存储。

## 在线预览

- 博客地址：[https://teshh.com](https://teshh.com)
- 后台入口：[https://teshh.com/#/sign-in](https://teshh.com/#/sign-in)

> 后台需要管理员账号，项目不提供公开测试账号。

## 功能特性

### 博客前台

- 文章与纯文本动态信息流
- 文章缩略图和富文本图片展示
- 一级、二级分类导航
- 标题搜索与分页
- 文章详情和评论
- 基于访客标识的点赞功能
- 响应式布局与深色模式

### 管理后台

- JWT 登录鉴权
- 文章新增、编辑、隐藏、置顶和删除
- 按标题、分类和状态筛选文章
- WangEditor 富文本编辑器
- 本地或 Cloudflare R2 图片上传
- 树形分类管理
- 评论与友情链接管理
- PV、UV、文章、评论和点赞统计

## 技术栈

| 模块 | 技术 |
| --- | --- |
| 前端 | Vue 3、TypeScript、Vite、Pinia、Vue Router |
| UI | Bootstrap 5、Element Plus、Bootstrap Icons |
| 编辑器 | WangEditor |
| 后端 | Go、Gin、GORM、Viper、Zap |
| 数据库 | MySQL |
| 文件存储 | 本地存储、Cloudflare R2 |

## 项目结构

```text
blog/
├─ server/                 # Gin 后端
│  ├─ api/                 # API 控制器
│  ├─ config/              # 配置结构
│  ├─ initialize/          # 路由、数据库初始化
│  ├─ middleware/          # JWT、CORS 等中间件
│  ├─ model/               # 数据模型
│  ├─ router/              # 路由注册
│  ├─ service/             # 业务逻辑
│  ├─ utils/               # 通用工具与存储实现
│  ├─ config.yaml          # 本地开发配置
│  └─ config.production.yaml
├─ web/                    # Vue 前端
├─ blog.sql                # 初始数据库
└─ README.md
```

## 环境要求

- Node.js 22+
- Go 1.19+
- MySQL 5.7+ 或兼容版本（生产环境推荐 MySQL 8.0+）

Linux、macOS 和 Windows 均可用于本地开发。生产环境推荐使用常见的 Linux 发行版，x86_64 与 ARM64 均可运行。

## 快速开始

### 1. 获取代码

```bash
git clone <repository-url> feilong-blog
cd feilong-blog
```

### 2. 初始化数据库

```sql
CREATE DATABASE blog CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
```

导入项目根目录的 `blog.sql`：

```bash
mysql -u root -p blog < blog.sql
```

### 3. 配置并启动后端

复制配置模板，再修改数据库、JWT 和存储配置。`config.yaml` 包含敏感信息，已被 Git 忽略。

```bash
cd server
cp config.example.yaml config.yaml
```

使用本地文件存储时：

```yaml
system:
  env: public
  addr: 8889
  db-type: mysql
  oss-type: local

local:
  path: uploads/file
  store-path: uploads/file
```

使用 Cloudflare R2 时：

```yaml
system:
  env: public
  addr: 8889
  db-type: mysql
  oss-type: r2

r2:
  account-id: your-account-id
  access-key-id: your-access-key-id
  secret-access-key: your-secret-access-key
  bucket: your-bucket-name
  public-url: https://img.example.com
  max-size-mb: 10
```

MySQL 和 JWT 配置示例：

```yaml
mysql:
  path: 127.0.0.1
  port: "3306"
  config: charset=utf8mb4&parseTime=True&loc=Local
  db-name: blog
  username: root
  password: your-password

jwt:
  signing-key: replace-with-a-random-string-at-least-32-characters
  expires-time: 24h
  buffer-time: 2h
  issuer: blog
```

启动后端：

```bash
cd server
go mod download
go run . -c config.yaml
```

检查服务：

```bash
curl http://127.0.0.1:8889/health
```

### 4. 启动前端

```bash
cd web
npm install
npm run dev
```

访问 [http://localhost:5173](http://localhost:5173)。开发服务器会把 `/blog` 请求代理到 `http://127.0.0.1:8889`。

## 配置说明

后端通过 YAML 文件读取配置，可使用以下方式指定配置文件：

```bash
./blog-server -c /path/to/config.yaml
```

也可以设置 `BLOG_CONFIG`：

```bash
export BLOG_CONFIG=/path/to/config.yaml
./blog-server
```

配置文件优先级为：命令行 `-c`、`BLOG_CONFIG`、默认配置文件。生产环境建议始终显式指定 `config.production.yaml`。

前端 API 地址通过 `VITE_BLOG_API_URL` 设置：

```env
VITE_BLOG_API_URL=/blog
```

推荐在生产环境保持 `/blog`，再由 Nginx 或其他反向代理转发至 Gin。这样前端和 API 使用同一域名，无需额外处理跨域。

## 文件存储

### 本地存储

设置 `system.oss-type: local`。上传文件会保存到 `local.store-path`，并通过 `local.path` 对外提供访问。

本地存储适合开发、小型站点或已经具备共享磁盘的环境。部署时应备份上传目录，并确保运行后端的用户具有写入权限。

### Cloudflare R2

设置 `system.oss-type: r2`，并填写 R2 S3 凭据、桶名和公开访问地址。

建议：

- 为博客单独创建只拥有目标桶对象读写权限的凭据。
- 使用自定义域名提供公开图片，不在生产环境使用 `r2.dev` 地址。
- R2 密钥只存放在后端配置中，不能写入 `VITE_*` 前端变量。
- 根据实际前端域名配置桶的 CORS 规则。

参考文档：[R2 公共桶](https://developers.cloudflare.com/r2/buckets/public-buckets/)、[R2 CORS](https://developers.cloudflare.com/r2/buckets/cors/)。

## 构建

### 前端

```bash
cd web
npm ci
npm run build
```

构建结果位于 `web/dist`。

### 后端

```bash
cd server
go build -trimpath -ldflags="-s -w" -o blog-server .
```

交叉编译示例（后端不依赖 CGO）：

```bash
# Linux x86_64
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -trimpath -ldflags="-s -w" -o blog-server .

# Linux ARM64
CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -trimpath -ldflags="-s -w" -o blog-server .
```

在 PowerShell 中交叉编译时使用 `$env:CGO_ENABLED`、`$env:GOOS` 和 `$env:GOARCH` 设置环境变量。如需进一步缩小文件，可在构建后执行 `upx --best blog-server`，并使用 `upx -t blog-server` 检查压缩产物。

## 推荐部署方式

适合大多数用户的部署组合：

- 一台 Linux VPS 或云服务器
- Nginx 提供前端静态文件和 HTTPS
- systemd 管理 Gin 服务
- MySQL 使用本机实例或托管数据库
- 图片使用本地存储或兼容 S3 的对象存储

```text
Browser
   │
   ▼
Nginx :80 / :443
   ├─ /           → /var/www/feilong-blog/web
   └─ /blog/*     → Gin 127.0.0.1:8889
                          │
                          ├─ MySQL
                          └─ Local / R2 storage
```

### 1. 准备生产配置

在本地复制并修改生产配置，确保以下内容不再使用示例值：

```bash
cp server/config.example.yaml server/config.production.yaml
```

- MySQL 地址、账号和密码
- JWT signing key
- CORS 域名
- R2 凭据和公开地址（使用 R2 时）

建议限制配置文件权限：

```bash
chmod 600 server/config.production.yaml
```

生产配置不要提交到 Git。推荐在本地完成构建，然后仅上传后端二进制、生产配置和前端 `dist` 成品。服务器目录可整理为：

```text
/var/www/feilong-blog/
├─ blog-server
├─ config.production.yaml
├─ web/
├─ log/
└─ backups/
```

### 2. 使用 systemd 运行后端

创建 `/etc/systemd/system/feilong-blog.service`：

```ini
[Unit]
Description=Feilong Blog API
After=network-online.target
Wants=network-online.target

[Service]
Type=simple
User=www-data
Group=www-data
WorkingDirectory=/var/www/feilong-blog
ExecStart=/var/www/feilong-blog/blog-server -c /var/www/feilong-blog/config.production.yaml
Restart=on-failure
RestartSec=5
NoNewPrivileges=true
PrivateTmp=true

[Install]
WantedBy=multi-user.target
```

启动服务：

```bash
sudo systemctl daemon-reload
sudo systemctl enable --now feilong-blog
sudo systemctl status feilong-blog
```

确保 `www-data` 对工作目录、日志目录以及本地上传目录拥有所需权限，并将生产配置权限设为 `600`。

### 3. 配置 Nginx

```nginx
server {
    listen 80;
    server_name blog.example.com;

    root /var/www/feilong-blog/web;
    index index.html;
    client_max_body_size 10m;

    location / {
        try_files $uri $uri/ /index.html;
    }

    location /blog/ {
        proxy_pass http://127.0.0.1:8889/;
        proxy_http_version 1.1;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
}
```

复制前端产物并重载 Nginx：

```bash
sudo mkdir -p /var/www/feilong-blog/web
sudo cp -a web/dist/. /var/www/feilong-blog/web/
sudo nginx -t
sudo systemctl reload nginx
```

使用 Certbot、Caddy 自动 HTTPS 或云平台证书为站点启用 TLS。若域名接入 Cloudflare，建议源站具备有效证书并使用 **Full (strict)** 模式。

### 其他部署选择

- 前端可以单独发布到 Cloudflare Pages、Vercel、Netlify 等静态托管平台。
- 后端可以部署到任何能够长期运行 Go 程序并访问 MySQL 的平台。
- 前后端使用不同域名时，需要同时配置 `VITE_BLOG_API_URL` 和后端 CORS 白名单。
- 使用托管数据库时，应开启 TLS、限制来源地址并使用专用数据库账号。

## 更新部署

在本地重新构建并完成测试，只把成品上传到服务器。替换后端二进制前先停止服务，避免覆盖正在执行的文件：

```bash
# 在服务器执行，文件名可按实际上传路径调整
sudo systemctl stop feilong-blog
sudo install -m 0755 /tmp/blog-server /var/www/feilong-blog/blog-server
sudo systemctl start feilong-blog

sudo cp -a /tmp/web-dist/. /var/www/feilong-blog/web/
sudo nginx -t
sudo systemctl reload nginx
```

更新后检查：

```bash
curl http://127.0.0.1:8889/health
sudo systemctl status feilong-blog
journalctl -u feilong-blog --since "10 minutes ago"
```

## 测试

```bash
# 后端
cd server
go test ./...

# 前端类型检查和生产构建
cd ../web
npm run type-check
npm run build
```

## 安全建议

- 不要提交数据库密码、JWT 密钥或对象存储密钥。
- 为数据库和对象存储创建最小权限账号。
- 不要将 MySQL 和 Gin 的内部端口直接暴露到公网。
- 定期备份 MySQL 和本地上传文件。
- 限制上传大小和允许的文件类型。
- 首次部署后立即修改默认管理员密码。
- 已经出现在提交记录、日志或聊天中的凭据应立即轮换。

## License

项目暂未声明开源许可证。使用、修改或分发前请先获得作者许可。
