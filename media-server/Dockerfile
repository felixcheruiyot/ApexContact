# Build Nginx with RTMP module
FROM alpine:3.20 AS builder

RUN apk add --no-cache build-base pcre-dev openssl-dev zlib-dev wget

# Download Nginx and nginx-rtmp-module
ARG NGINX_VERSION=1.26.1
ARG RTMP_VERSION=1.2.2

WORKDIR /build
RUN wget -q https://nginx.org/download/nginx-${NGINX_VERSION}.tar.gz && \
    tar xf nginx-${NGINX_VERSION}.tar.gz && \
    wget -q https://github.com/arut/nginx-rtmp-module/archive/v${RTMP_VERSION}.tar.gz && \
    tar xf v${RTMP_VERSION}.tar.gz

RUN cd nginx-${NGINX_VERSION} && \
    ./configure \
      --prefix=/usr/local/nginx \
      --with-http_ssl_module \
      --with-http_v2_module \
      --add-module=../nginx-rtmp-module-${RTMP_VERSION} && \
    make -j$(nproc) && \
    make install

# Runtime
FROM alpine:3.20

RUN apk add --no-cache pcre openssl

COPY --from=builder /usr/local/nginx /usr/local/nginx
COPY nginx.conf /usr/local/nginx/conf/nginx.conf

RUN mkdir -p /tmp/hls

EXPOSE 1935 8888

CMD ["/usr/local/nginx/sbin/nginx", "-g", "daemon off;"]
