FROM node:24.1.0-alpine AS builder

ARG VITE_APP_API_BASE_URL # Declare build argument

ENV VITE_APP_API_BASE_URL=${VITE_APP_API_BASE_URL}

RUN echo "Build-time VITE_APP_API_BASE_URL is: $VITE_APP_API_BASE_URL"

WORKDIR /app


COPY package*.json ./

RUN npm ci


COPY . .

RUN npm run build


FROM nginx:stable-alpine3.21-slim

RUN rm -rf /usr/share/nginx/html/*

COPY --from=builder /app/dist /usr/share/nginx/html

COPY nginx.conf /etc/nginx/conf.d/default.conf

EXPOSE 80

CMD ["nginx", "-g", "daemon off;"]