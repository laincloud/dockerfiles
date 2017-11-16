FROM laincloud/centos:7.4.1708

RUN yum update -y \
    && yum install -y python36u \
                      python36u-devel \
                      python36u-pip \
    && yum clean all \
    && ln -sf /usr/bin/python3.6 /usr/bin/python3 \
    && ln -sf /usr/bin/pip3.6 /usr/bin/pip

CMD ["python3"]