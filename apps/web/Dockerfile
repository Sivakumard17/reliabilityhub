FROM node:20-alpine AS builder

WORKDIR /app

COPY apps/web/package*.json ./
RUN npm ci

COPY apps/web/ .

ENV NEXT_TELEMETRY_DISABLED=1

ARG NEXT_PUBLIC_API_URL=http://localhost:9090
ENV NEXT_PUBLIC_API_URL=$NEXT_PUBLIC_API_URL

RUN npm run build

FROM node:20-alpine AS runner

WORKDIR /app
ENV NODE_ENV=production
ENV NEXT_TELEMETRY_DISABLED=1

RUN addgroup -g 1001 nodejs && \
    adduser -u 1001 -G nodejs -s /bin/sh -D nextjs

COPY --from=builder --chown=nextjs:nodejs /app/.next/standalone ./
COPY --from=builder --chown=nextjs:nodejs /app/.next/static ./.next/static
RUN mkdir -p /app/public

USER nextjs
EXPOSE 3000
ENV PORT=3000
ENV HOSTNAME="0.0.0.0"

CMD ["node", "server.js"]
