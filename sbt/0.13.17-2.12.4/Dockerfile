# TAGS 0.13.17-2.12.4
FROM laincloud/scala:2.12.4

Run wget -O- "https://github.com/sbt/sbt/releases/download/v0.13.17/sbt-0.13.17.tgz" \
    |  tar xzf - -C /usr/local --strip-components=1 \
    && sbt exit

CMD ["sbt"]
