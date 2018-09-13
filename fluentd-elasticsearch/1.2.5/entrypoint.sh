#!/bin/sh


set -o errexit

mkdir -p /var/log/journal

if [ -z ${FLUENT_ELASTICSEARCH_USER} ] ; then
   sed -i  '/FLUENT_ELASTICSEARCH_USER/d' /etc/fluent/fluent.conf
fi

if [ -z ${FLUENT_ELASTICSEARCH_PASSWORD} ] ; then
   sed -i  '/FLUENT_ELASTICSEARCH_PASSWORD/d' /etc/fluent/fluent.conf
fi

exec /usr/local/bin/fluentd $@
