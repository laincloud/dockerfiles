# TAGS 1.0.2-2.12.4
FROM laincloud/scala:2.12.4

Run wget -O- "https://github.com/sbt/sbt/releases/download/v1.0.2/sbt-1.0.2.tgz" \
    |  tar xzf - -C /usr/local --strip-components=1 \
    && sbt exit

CMD ["sbt"]
