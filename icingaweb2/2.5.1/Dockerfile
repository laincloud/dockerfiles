FROM laincloud/debian:stretch

ENV ICINGAWEB2_VERSION 2.5.1-1.stretch

RUN curl  http://packages.icinga.com/icinga.key | apt-key add - && \
    echo 'deb http://packages.icinga.com/debian icinga-stretch main' >/etc/apt/sources.list.d/icinga.list && \
    apt-get update && \
    apt-get install -qqy icingaweb2=$ICINGAWEB2_VERSION icingacli=$ICINGAWEB2_VERSION php7.0-gd && \
    icingacli setup config webserver apache && \
    echo 'date.timezone = Asia/Shanghai' > /etc/php/7.0/apache2/conf.d/timeszone.ini

COPY apache2-foreground /usr/local/bin/

EXPOSE 80

CMD ["apache2-foreground"]
