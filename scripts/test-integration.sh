#!/usr/bin/env bash
# Disposable MariaDB + Go + Nginx fixture. Never connects to the deployed site.
set -euo pipefail
root=$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)
go_bin=${GO_BIN:-go}
for command in docker python3 node npm; do command -v "$command" >/dev/null; done
"$go_bin" version
fixture=$(mktemp -d)
chmod 755 "$fixture"
name="feilong-integration-$$-$(date +%s)"
cleanup() {
  if [[ ${1:-0} != 0 ]]; then docker logs "$name-web" 2>/dev/null || true; docker exec "$name-web" tail -60 /tmp/backend.log 2>/dev/null || true; fi
  docker rm -fv "$name-web" "$name-db" >/dev/null 2>&1 || true
  docker network rm "$name" >/dev/null 2>&1 || true
  rm -rf "$fixture"
}
trap 'cleanup $?' EXIT
# Bind ephemeral ports to loopback only, with a network owned by this invocation.
docker network create "$name" >/dev/null
docker run -d --name "$name-db" --network "$name" --network-alias db \
  -p 127.0.0.1::3306 -e MARIADB_ALLOW_EMPTY_ROOT_PASSWORD=1 -e MARIADB_DATABASE=blog \
  mariadb:11.4 >/dev/null
for attempt in {1..60}; do
  if docker exec "$name-db" mariadb-admin -h127.0.0.1 ping --silent >/dev/null 2>&1; then break; fi
  if [[ $attempt == 60 ]]; then echo 'MariaDB failed to start' >&2; exit 1; fi
  sleep 1
done
docker exec -i "$name-db" mariadb -uroot blog < "$root/blog.sql"
docker exec "$name-db" mariadb -uroot -e 'CREATE DATABASE feilong_integration_product CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci'
db_port=$(docker port "$name-db" 3306/tcp | awk -F: '{print $NF}')
(cd "$root/server" && BLOG_TEST_MYSQL_DSN="root@tcp(127.0.0.1:$db_port)/feilong_integration_product?charset=utf8mb4&parseTime=true&loc=UTC" GIN_MODE=release "$go_bin" test ./... -count=1)
(cd "$root/web" && npm ci && npm run build)
case $(docker info --format '{{.Architecture}}') in aarch64|arm64) arch=arm64 ;; x86_64|amd64) arch=amd64 ;; *) echo 'Unsupported Docker architecture' >&2; exit 1 ;; esac
(cd "$root/server" && CGO_ENABLED=0 GOOS=linux GOARCH="$arch" "$go_bin" build -trimpath -o "$fixture/blog-server" .)
chmod 755 "$fixture/blog-server"
# Reserve a container port before generating canonical URLs.
docker create --name "$name-web" --init --network "$name" -p 127.0.0.1::8080 \
  -e GIN_MODE=release -w /tmp -v "$fixture:/app/runtime:ro" -v "$root/web/dist:/app/web:ro" \
  -v "$fixture/nginx.conf:/etc/nginx/conf.d/default.conf:ro" nginx:1.28-alpine \
  sh -c 'while [ ! -f /app/runtime/ready ]; do sleep 1; done; /app/runtime/blog-server -c /app/runtime/config.yaml >/tmp/backend.log 2>&1 & exec nginx -g "daemon off;"' >/dev/null
# Generate placeholder configs before mounting files; start the processes after the port is known.
python3 - "$root" "$fixture" <<'PY'
import sys
from pathlib import Path
root, target = map(Path, sys.argv[1:])
s = (root/'server/config.example.yaml').read_text().replace('path: 127.0.0.1', 'path: db').replace('username: blog', 'username: root').replace('password: replace-with-your-database-password', 'password: ""').replace('web-root: ../web/dist', 'web-root: /app/web').replace('store-path: uploads/file', 'store-path: /tmp/uploads/file')
(target/'config.yaml').write_text(s)
n = (root/'web/nginx.conf.example').read_text().replace('listen 80;', 'listen 8080;').replace('root /var/www/feilong-blog/web;', 'root /app/web;')
(target/'nginx.conf').write_text(n)
for p in (target/'config.yaml', target/'nginx.conf'): p.chmod(0o644)
PY
docker start "$name-web" >/dev/null
web_port=$(docker port "$name-web" 8080/tcp | awk -F: '{print $NF}')
base_url="http://127.0.0.1:$web_port"
python3 - "$fixture" "$base_url" <<'PY'
import sys
from pathlib import Path
target, base = Path(sys.argv[1]), sys.argv[2]
for name in ('config.yaml', 'nginx.conf'):
    p = target/name
    p.write_text(p.read_text().replace('https://blog.example.com', base))
PY
touch "$fixture/ready"
python3 - "$base_url" <<'PY'
import sys, time, urllib.request
for attempt in range(60):
    try:
        with urllib.request.urlopen(sys.argv[1]+'/blog/health', timeout=2) as response:
            if response.status == 200: break
    except Exception: pass
    time.sleep(1)
else: raise SystemExit('Backend failed to start')
PY
(cd "$root/web" && BLOG_INTEGRATION_URL="$base_url" npm run test:integration)
docker exec "$name-web" nginx -t
