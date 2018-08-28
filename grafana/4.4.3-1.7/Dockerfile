# TAGS 4.4-1.7 4.4.3-1.7
FROM laincloud/debian:jessie

ARG GRAFANA_VERSION=4.4.3
ARG GRAFANA_DEB_URL=https://s3-us-west-2.amazonaws.com/grafana-releases/release/grafana_${GRAFANA_VERSION}_amd64.deb
ARG GOSU_VERSION=1.7
ARG GOSU_BIN_URL=https://github.com/tianon/gosu/releases/download/${GOSU_VERSION}/gosu-amd64

ENV \
  GRAFANA_VERSION=${GRAFANA_VERSION} \
  GRAFANA_DEB_URL=${GRAFANA_DEB_URL} \
  GOSU_VERSION=${GOSU_VERSION} \
  GOSU_BIN_URL=${GOSU_BIN_URL} \
  GF_PLUGIN_DIR=/grafana-plugins \
  GF_PATHS_LOGS=/var/log/grafana \
  GF_PATHS_DATA=/var/lib/grafana

COPY ./run.sh /run.sh

RUN \
  apt-get update && \
  apt-get -y --force-yes --no-install-recommends install libfontconfig curl ca-certificates git jq && \
  curl -L ${GRAFANA_DEB_URL} > /tmp/grafana.deb && \
  dpkg -i /tmp/grafana.deb && \
  rm -f /tmp/grafana.deb && \
  curl -L ${GOSU_BIN_URL} > /usr/sbin/gosu && \
  chmod +x /usr/sbin/gosu && \
  chmod +x /run.sh && \
  apt-get remove -y --force-yes curl git jq && \
  apt-get autoremove -y --force-yes && \
  apt-get clean && \
  rm -rf /var/lib/apt/lists/*

VOLUME ["/var/lib/grafana", "/var/log/grafana", "/etc/grafana"]

EXPOSE 3000

ENTRYPOINT ["/run.sh"]
