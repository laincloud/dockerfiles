# TAGS 3.7
FROM alpine:3.7

RUN apk add wget --update-cache --repository http://mirrors.ustc.edu.cn/alpine/v3.7/main/ --allow-untrusted

ENV S6_OVERLAY_VERSION=v1.21.2.2 \
    S6_LOGGING_SCRIPT=n3\ s268435455 \
    S6_KEEP_ENV=1

RUN wget -q https://github.com/just-containers/s6-overlay/releases/download/${S6_OVERLAY_VERSION}/s6-overlay-amd64.tar.gz -O /tmp/s6-overlay-amd64.tar.gz \
    && tar xfz /tmp/s6-overlay-amd64.tar.gz -C / \
    && rm -f /tmp/s6-overlay-amd64.tar.gz

COPY rootfs /

ENTRYPOINT ["/app-init"]
