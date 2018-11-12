# TAGS 5.2 5.2.3
FROM laincloud/debian:stretch

ARG GRAFANA_DEB_URL="https://s3-us-west-2.amazonaws.com/grafana-releases/release/grafana_5.2.3_amd64.deb"
ARG GOSU_BIN_URL="https://github.com/tianon/gosu/releases/download/1.10/gosu-amd64"

ENV PATH=/usr/share/grafana/bin:/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin \
    GF_PATHS_DATA="/var/lib/grafana" \
    GF_PATHS_LOGS="/var/log/grafana" \
    GF_PATHS_PLUGINS="/var/lib/grafana/plugins" \
    GF_PATHS_PROVISIONING="/etc/grafana/provisioning"

COPY ./run.sh /run.sh

RUN apt-get update && apt-get install -qq -y libfontconfig curl ca-certificates && \
    mkdir -p "$GF_PATHS_PROVISIONING/datasources" \
             "$GF_PATHS_PROVISIONING/dashboards" \
             "$GF_PATHS_PLUGINS" \
             "$GF_PATHS_LOGS" \
             "$GF_PATHS_DATA" && \
    curl -L "$GRAFANA_DEB_URL" > /tmp/grafana.deb && \
    dpkg -i /tmp/grafana.deb && \
    rm -f /tmp/grafana.deb && \
    curl -L "$GOSU_BIN_URL" > /usr/sbin/gosu && \
    chmod +x /usr/sbin/gosu && \
    apt-get remove -y --force-yes curl && \
    apt-get autoremove -y --force-yes && \
    apt-get clean && \
    rm -rf /var/lib/apt/lists/* && \
    chmod +x /run.sh

VOLUME ["/var/lib/grafana", "/var/log/grafana", "/etc/grafana"]

EXPOSE 3000

CMD [ "/run.sh" ]
