FROM laincloud/golang:1.9.4

ENV NPM_CONFIG_LOGLEVEL info
ENV NODE_VERSION 8.4.0

RUN ARCH='x64' \
  && curl -SLO "https://nodejs.org/dist/v$NODE_VERSION/node-v$NODE_VERSION-linux-$ARCH.tar.xz" \
  && tar -xJf "node-v$NODE_VERSION-linux-$ARCH.tar.xz" -C /usr/local --strip-components=1 \
  && ln -s /usr/local/bin/node /usr/local/bin/nodejs

RUN go get github.com/tools/godep