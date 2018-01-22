FROM buildpack-deps:stretch

ENV LANG C.UTF-8

RUN apt-get -qq -y update \
    && apt-get -qq -y --no-install-recommends install \
                       dnsutils \
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

ENV S6_OVERLAY_VERSION=v1.21.2.2 \
    S6_LOGGING_SCRIPT=n3\ s268435455 \
    S6_KEEP_ENV=1

RUN wget -q https://github.com/just-containers/s6-overlay/releases/download/${S6_OVERLAY_VERSION}/s6-overlay-amd64.tar.gz -O /tmp/s6-overlay-amd64.tar.gz \
    && tar xfz /tmp/s6-overlay-amd64.tar.gz -C / \
    && rm -f /tmp/s6-overlay-amd64.tar.gz

COPY rootfs /

ENTRYPOINT ["/app-init"]