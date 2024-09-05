# 1. For build React app
FROM node:lts AS development

ARG VITE_GOOGLE_CLIENT_ID
ARG VITE_GOOGLE_REDIRECT_URI
ARG VITE_BACKEND_URL
ARG VITE_BACKEND_WS_URL

WORKDIR /app

COPY package.json /app/package.json
COPY package-lock.json /app/package-lock.json

RUN --mount=type=cache,target=/root/.npm npm ci

COPY . /app

ENV CI=true
ENV PORT=80

CMD [ "npm", "run", "dev"]

FROM development AS build

RUN npm run build

# 2. For Nginx setup
FROM nginx:alpine as production

# Copy config nginx
COPY --from=build /app/.nginx/nginx.conf /etc/nginx/conf.d/default.conf

WORKDIR /usr/share/nginx/html

# Remove default nginx static assets
RUN rm -rf ./*

# Copy static assets from builder stage
COPY --from=build /app/dist .

# Containers run nginx with global directives and daemon off
ENTRYPOINT ["nginx", "-g", "daemon off;"]
