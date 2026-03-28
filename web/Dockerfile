FROM oven/bun:1 AS builder

WORKDIR /src

COPY package.json bun.lock ./
RUN bun install --frozen-lockfile

COPY . .
RUN bun run build

FROM alpine:3.21

COPY --from=builder /src/dist /dist

CMD ["cp", "-r", "/dist/.", "/srv/"]
