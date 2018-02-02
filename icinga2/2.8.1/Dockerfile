FROM laincloud/debian:stretch

ENV ICINGA2_VERSION 2.8.1-1.stretch

RUN wget -O - https://packages.icinga.com/icinga.key | apt-key add - && \
    echo 'deb https://packages.icinga.com/debian icinga-stretch main' >/etc/apt/sources.list.d/icinga.list && \
    apt-get update && \
    DEBIAN_FRONTEND=noninteractive apt-get install -qqy --no-install-recommends monitoring-plugins icinga2=$ICINGA2_VERSION icinga2-ido-mysql=$ICINGA2_VERSION && \
    icinga2 feature enable ido-mysql

COPY etc /

COPY start.sh /

CMD /start.sh