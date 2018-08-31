#!/bin/bash

chown -R grafana:grafana /var/lib/grafana /var/log/grafana
chmod 777 /var/lib/grafana /var/log/grafana

exec gosu grafana /usr/sbin/grafana-server   \
  --homepath=/usr/share/grafana              \
  --config=/etc/grafana/grafana.ini          \
  cfg:default.paths.data=${GF_PATHS_DATA}    \
  cfg:default.paths.logs=${GF_PATHS_LOGS}    \
  cfg:default.paths.plugins=${GF_PLUGIN_DIR} \
  "$@"
