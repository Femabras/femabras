FROM node:22-alpine
WORKDIR /app

# Install pnpm for enterprise speed
RUN corepack enable && corepack prepare pnpm@latest --activate

COPY package.json pnpm-lock.yaml* ./
RUN pnpm install --frozen-lockfile

COPY . .

# We build in production mode to trigger the strict locale checks
RUN pnpm run build

EXPOSE 3000
CMD ["pnpm", "start"]