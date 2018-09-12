#!/bin/sh

mkdir -p /var/log/journal

exec /usr/local/bin/fluentd $@
