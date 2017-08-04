FROM centos:7.3.1611

RUN yum update -y \
    && yum install -y telnet traceroute net-tools tcpdump iptraf nmap \
                      lsof unzip xz vim git wget \
    && yum clean all

ENV S6_OVERLAY_VERSION=v1.19.1.1 \
    S6_LOGGING_SCRIPT=n3\ s268435455

RUN wget https://keybase.io/justcontainers/key.asc --no-check-certificate -O /tmp/s6-overlay-key.asc \
    && wget https://github.com/just-containers/s6-overlay/releases/download/${S6_OVERLAY_VERSION}/s6-overlay-amd64.tar.gz --no-check-certificate -O /tmp/s6-overlay-amd64.tar.gz \
    && wget https://github.com/just-containers/s6-overlay/releases/download/${S6_OVERLAY_VERSION}/s6-overlay-amd64.tar.gz.sig --no-check-certificate -O /tmp/s6-overlay-amd64.tar.gz.sig \
    && gpg --import /tmp/s6-overlay-key.asc \
    && gpg --verify /tmp/s6-overlay-amd64.tar.gz.sig /tmp/s6-overlay-amd64.tar.gz \
    && tar xzf /tmp/s6-overlay-amd64.tar.gz -C / --exclude="./bin" \
    && tar xzf /tmp/s6-overlay-amd64.tar.gz -C /usr ./bin \
    && rm -f /tmp/s6-overlay-key.asc \
    && rm -f /tmp/s6-overlay-amd64.tar.gz \
    && rm -f /tmp/s6-overlay-amd64.tar.gz.sig

COPY rootfs /

ENTRYPOINT ["/app-init"]