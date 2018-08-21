FROM centos:7.3.1611

RUN yum update -y \
    && yum groupinstall -y "Development tools" \
    && yum install -y git \
                      iproute \
                      iptraf \
                      lsof \
                      net-tools \
                      nmap \
                      tcpdump \
                      telnet \
                      traceroute \
                      unzip \
                      vim \
                      wget \
                      xz

ENV S6_OVERLAY_VERSION=v1.21.2.1 \
    S6_LOGGING_SCRIPT=n3\ s268435455 \
    S6_KEEP_ENV=1

RUN wget https://github.com/just-containers/s6-overlay/releases/download/${S6_OVERLAY_VERSION}/s6-overlay-amd64.tar.gz -O /tmp/s6-overlay-amd64.tar.gz \
    && tar xzf /tmp/s6-overlay-amd64.tar.gz -C / --exclude="./bin" \
    && tar xzf /tmp/s6-overlay-amd64.tar.gz -C /usr ./bin \
    && rm -f /tmp/s6-overlay-amd64.tar.gz

COPY rootfs /

ENTRYPOINT ["/app-init"]
