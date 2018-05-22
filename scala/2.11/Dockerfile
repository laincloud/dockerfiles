# TAGS 2.11.12 2.11
FROM laincloud/openjdk:8

RUN wget -O- "http://downloads.lightbend.com/scala/2.11.12/scala-2.11.12.tgz" \
    | tar xzf - -C /usr/local --strip-components=1

CMD ["scala"]
