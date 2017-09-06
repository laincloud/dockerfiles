FROM laincloud/centos:7.3.1611

COPY rootfs /
RUN yum update \
    && yum install -y nginx-1.13.5 \
    && yum clean all

EXPOSE 80

STOPSIGNAL SIGTERM

CMD ["nginx", "-g", "daemon off;"]