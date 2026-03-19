# Stage 1: Build + test Angular
FROM node:22-slim AS web-builder
RUN npm install -g npm@11 && \
    apt-get update && apt-get install -y --no-install-recommends \
    libnss3 libatk1.0-0 libatk-bridge2.0-0 libcups2 libdrm2 libxkbcommon0 \
    libxcomposite1 libxdamage1 libxfixes3 libxrandr2 libgbm1 libpango-1.0-0 \
    libcairo2 libasound2 libxshmfence1 && rm -rf /var/lib/apt/lists/*
WORKDIR /app/web
COPY web/package*.json ./
RUN --mount=type=cache,target=/root/.npm npm ci
RUN npx playwright install chromium
COPY web/ ./
RUN npm run lint && npm run format:check
RUN npm test
RUN npx ng build --configuration=production

# Stage 2: Build + test Go
FROM golang:1.25 AS go-builder
RUN apt-get update && apt-get install -y gcc libc6-dev && rm -rf /var/lib/apt/lists/*
RUN useradd -m tester
WORKDIR /app
COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod go mod download
COPY . .
COPY --from=web-builder /app/web/dist ./web/dist
RUN test -z "$(gofmt -l .)" || (echo "Go files not formatted:" && gofmt -l . && exit 1)
RUN chown -R tester:tester /app
USER tester
RUN --mount=type=cache,target=/home/tester/.cache/go-build,uid=1000 \
    --mount=type=cache,target=/go/pkg/mod \
    go test ./...
RUN --mount=type=cache,target=/home/tester/.cache/go-build,uid=1000 \
    --mount=type=cache,target=/go/pkg/mod \
    CGO_ENABLED=1 go build -o gesitr .
RUN --mount=type=cache,target=/home/tester/.cache/go-build,uid=1000 \
    --mount=type=cache,target=/go/pkg/mod \
    CGO_ENABLED=1 go build -o seed ./cmd/seed

# Stage 3: E2E tests — Playwright against the production binary
FROM node:22 AS e2e-tester
RUN npm install -g npm@11
RUN apt-get update && apt-get install -y --no-install-recommends \
    libnss3 libatk1.0-0 libatk-bridge2.0-0 libcups2 libdrm2 libxkbcommon0 \
    libxcomposite1 libxdamage1 libxfixes3 libxrandr2 libgbm1 libpango-1.0-0 \
    libcairo2 libasound2 libxshmfence1 && rm -rf /var/lib/apt/lists/*
WORKDIR /app/web
RUN chown node:node /app /app/web
COPY --chown=node:node web/package*.json ./
USER node
RUN --mount=type=cache,target=/home/node/.npm,uid=1000 npm ci
RUN npx playwright install chromium
COPY --chown=node:node web/playwright.config.ts ./
COPY --chown=node:node web/e2e/ ./e2e/
COPY --from=go-builder --chown=node:node /app/gesitr /app/gesitr
ENV PLAYWRIGHT_TEST_BASE_URL=http://localhost:8080
ENV AUTH_FALLBACK_USER=e2e-tester
WORKDIR /app
RUN ./gesitr & SERVER_PID=$! && \
    sleep 2 && \
    cd web && npx playwright test ; \
    TEST_EXIT=$? ; kill $SERVER_PID 2>/dev/null ; exit $TEST_EXIT
RUN date -u '+%Y-%m-%dT%H:%M:%SZ' > /tmp/.e2e-passed

# Stage 4: Runtime
# COPY --from=e2e-tester forces BuildKit to run the e2e stage before finalizing
FROM debian:bookworm-slim
RUN useradd -m app
RUN apt-get update && apt-get install -y ca-certificates && rm -rf /var/lib/apt/lists/*
WORKDIR /app
COPY --from=go-builder --chown=app:app /app/gesitr .
COPY --from=go-builder --chown=app:app /app/seed .
COPY --from=go-builder --chown=app:app /app/data ./data/
COPY --chown=app:app entrypoint.sh .
COPY --from=e2e-tester /tmp/.e2e-passed /tmp/.e2e-passed
RUN chmod +x entrypoint.sh && mkdir -p /app/db && chown app:app /app/db
ENV GIN_MODE=release
ENV DATABASE_PATH=/app/db/gesitr.db
VOLUME /app/db
USER app
EXPOSE 8080
ENTRYPOINT ["./entrypoint.sh"]
