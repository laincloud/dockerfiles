FROM buildpack-deps:sid

RUN sed -i 's/deb.debian.org/mirrors.ustc.edu.cn/g' /etc/apt/sources.list \
    && sed -i 's|security.debian.org|mirrors.ustc.edu.cn/debian-security|g' /etc/apt/sources.list \
    && apt-get update \
    && apt-get install -y --no-install-recommends \
                       vim \
                       net-tools \
                       tcpdump \
                       nmap \
                       lsof \
                       telnet \
                       iptraf \
                       traceroute \
                       unzip \
    && rm -rf /var/lib/apt/lists/*

ENV S6_OVERLAY_VERSION=v1.21.2.1 \
    S6_LOGGING_SCRIPT=n3\ s268435455 \
    S6_KEEP_ENV=1

RUN wget https://github.com/just-containers/s6-overlay/releases/download/${S6_OVERLAY_VERSION}/s6-overlay-amd64.tar.gz -O /tmp/s6-overlay-amd64.tar.gz \
    && tar xzf /tmp/s6-overlay-amd64.tar.gz -C / --exclude="./bin" \
    && tar xzf /tmp/s6-overlay-amd64.tar.gz -C /usr ./bin \
    && rm -f /tmp/s6-overlay-amd64.tar.gz

COPY rootfs /

ENTRYPOINT ["/app-init"]