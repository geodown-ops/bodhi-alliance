# 建置官網或後台的靜態檔，用 Caddy 提供（SPA 路由一律回 index.html）。
# Railway：Dockerfile 路徑設 apps/Dockerfile，root 用 repo 根目錄，並設定變數 APP=web 或 APP=admin、
# VITE_API_BASE、VITE_GUIDE_BASE（Railway 會把服務變數當成同名的 build arg 傳進來）。
FROM node:22 AS build
ARG APP=web
ARG VITE_API_BASE=
ARG VITE_GUIDE_BASE=
WORKDIR /src
COPY package.json package-lock.json ./
COPY apps/web/package.json apps/web/
COPY apps/admin/package.json apps/admin/
RUN npm ci
COPY apps apps
# 官網首頁是 repo 根目錄原本整頁的招募頁，複製成 recruit.html，由 Caddyfile 在 / 提供
COPY index.html recruit.html
RUN npm run build -w apps/${APP} && cp -r apps/${APP}/dist /out \
 && if [ "$APP" = web ]; then cp recruit.html /out/recruit.html; fi

FROM caddy:2-alpine
COPY --from=build /out /srv
COPY apps/Caddyfile /etc/caddy/Caddyfile
