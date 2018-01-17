FROM laincloud/php:7.2-apache

ENV LAST_RELEASE_URL http://dl.cnezsoft.com/zentao/9.7/ZenTaoPMS.9.7.stable.zip
ENV LAST_RELEASE_FILENAME ZenTaoPMS.9.7.stable

# configure timezone
RUN echo "Asia/Shanghai" > /etc/timezone;dpkg-reconfigure -f noninteractive tzdata

RUN curl -s -fSL $LAST_RELEASE_URL -o /tmp/$LAST_RELEASE_FILENAME && \
    cd /tmp && unzip -q $LAST_RELEASE_FILENAME && \
    mv zentaopms /var/www/html && \
    mkdir -p /var/www/html/zentaopms/tmp/php && \
    chmod o=rwx -R /var/www/html/zentaopms/tmp/php

RUN docker-php-ext-install pdo pdo_mysql

WORKDIR /var/www/html

VOLUME /data

COPY docker-entrypoint.sh /

COPY php.ini /usr/local/etc/php/

EXPOSE 80

CMD ["/docker-entrypoint.sh"]
