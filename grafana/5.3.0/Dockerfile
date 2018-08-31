# TAGS 5.3 5.3.0
FROM laincloud/debian:stretch

ARG GRAFANA_URL="https://s3-us-west-2.amazonaws.com/grafana-releases/master/grafana-latest.linux-x64.tar.gz"
ARG GF_UID="472"
ARG GF_GID="472"
ARG GOSU_BIN_URL=https://github.com/tianon/gosu/releases/download/1.10/gosu-amd64

ENV PATH=/usr/share/grafana/bin:/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin \
    GF_PATHS_HOME="/usr/share/grafana" \
    GF_PATHS_DATA="/var/lib/grafana" \
    GF_PATHS_LOGS="/var/log/grafana" \
    GF_PATHS_PLUGINS="/var/lib/grafana/plugins" \
    GF_PATHS_PROVISIONING="/etc/grafana/provisioning"

COPY ./run.sh /run.sh

RUN apt-get update && apt-get install -qq -y tar libfontconfig curl ca-certificates && \
    mkdir -p "$GF_PATHS_HOME" \
             "$GF_PATHS_PROVISIONING/datasources" \
             "$GF_PATHS_PROVISIONING/dashboards" \
             "$GF_PATHS_PLUGINS" \
             "$GF_PATHS_LOGS" \
             "$GF_PATHS_DATA" && \
    curl -L "$GRAFANA_URL" | tar zxvf - --strip-components=1 -C "$GF_PATHS_HOME" && \
    cp "$GF_PATHS_HOME/conf/sample.ini" /etc/grafana/grafana.ini && \
    cp "$GF_PATHS_HOME/conf/ldap.toml" /etc/grafana/ldap.toml && \
    curl -L "$GOSU_BIN_URL" > /usr/sbin/gosu && \
    chmod +x /usr/sbin/gosu && \
    chmod +x /run.sh && \
    apt-get autoremove -y && \
    rm -rf /var/lib/apt/lists/* && \
    groupadd -r -g $GF_GID grafana && \
    useradd -r -u $GF_UID -g grafana grafana

VOLUME ["/usr/share/grafana", "/var/lib/grafana", "/var/log/grafana", "/etc/grafana"]

EXPOSE 3000

ENTRYPOINT [ "/run.sh" ]
