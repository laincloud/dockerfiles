# TAGS 2.12.4 2.12 latest
FROM laincloud/openjdk:8

RUN wget -O- "http://downloads.lightbend.com/scala/2.12.4/scala-2.12.4.tgz" \
    | tar xzf - -C /usr/local --strip-components=1

CMD ["scala"]
