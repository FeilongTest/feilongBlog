#!/usr/bin/env bash
# Build a reviewable Linux release; never reads/copies the deployed config or uploads.
set -euo pipefail
root=$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)
arch=${1:-arm64}
case "$arch" in arm64|amd64) ;; *) echo 'Usage: scripts/release.sh [arm64|amd64] [output-directory]' >&2; exit 2 ;; esac
output=${2:-"$root/../feilongBlog-releases"}
go_bin=${GO_BIN:-go}
"$go_bin" version
mkdir -p "$output"
output=$(cd "$output" && pwd)
revision=$(git -C "$root" rev-parse --short=12 HEAD)
if [[ -n $(git -C "$root" status --porcelain) ]]; then revision="$revision-dirty"; fi
name="feilongBlog-linux-$arch-$revision"
if [[ -e "$output/$name.tar.gz" ]]; then echo "Release already exists: $output/$name.tar.gz" >&2; exit 1; fi
stage=$(mktemp -d)
trap 'rm -rf "$stage"' EXIT
mkdir -p "$stage/$name/server" "$stage/$name/web"
(cd "$root/web" && npm ci && npm run build)
(cd "$root/server" && CGO_ENABLED=0 GOOS=linux GOARCH="$arch" "$go_bin" build -trimpath -ldflags='-s -w' -o "$stage/$name/server/blog-server" .)
cp -R "$root/web/dist/." "$stage/$name/web/"
cp "$root/server/config.example.yaml" "$stage/$name/server/config.example.yaml"
# Example expects web/dist in source checkout; release has a flat web directory.
python3 - "$stage/$name/server/config.example.yaml" <<'PY'
import sys
p = sys.argv[1]
s = open(p).read().replace('web-root: ../web/dist', 'web-root: ../web')
open(p, 'w').write(s)
PY
cp "$root/web/nginx.conf.example" "$stage/$name/nginx.conf.example"
printf 'revision=%s\nplatform=linux/%s\nbuilt_at=%s\n' "$revision" "$arch" "$(date -u +%FT%TZ)" > "$stage/$name/BUILD.txt"
cat > "$stage/$name/UPGRADE.txt" <<'TEXT'
Back up the existing database, private config, uploads and old executable first.
Stop the old backend; replace the executable and public web assets together.
Keep the production config and uploads. Do not import blog.sql into an existing database.
Check seo.web-root and Nginx root against the actual installed web directory.
Start the backend with an explicit -c /absolute/path/to/config.production.yaml.
Startup adds publishing columns and three tables, and widens legacy article columns.
Check /blog/health, /tags, an article URL and /sitemap.xml through Nginx.
This package contains example configs only; it does not contain credentials or data.
TEXT
find "$stage/$name" -type d -exec chmod 755 {} +
find "$stage/$name" -type f -exec chmod 644 {} +
chmod 755 "$stage/$name/server/blog-server"
(cd "$stage" && tar -czf "$output/$name.tar.gz" "$name")
(cd "$output" && sha256sum "$name.tar.gz" > "$name.tar.gz.sha256")
printf 'Release: %s/%s.tar.gz\n' "$output" "$name"
