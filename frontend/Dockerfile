FROM node:22-alpine AS builder

WORKDIR /app

COPY package.json package-lock.json ./
RUN npm ci

COPY . .

# Build the Next.js application
RUN npm run build

FROM node:22-alpine

WORKDIR /app

# Install only production dependencies
COPY package.json package-lock.json ./
RUN npm ci --only=production

# Copy built application from builder
COPY --from=builder /app/.next ./.next
COPY --from=builder /app/public ./public

EXPOSE 3000

CMD ["npm", "start"]
