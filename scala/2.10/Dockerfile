# TAGS 2.10.7 2.10
FROM laincloud/openjdk:8

RUN wget -O- "http://downloads.lightbend.com/scala/2.10.7/scala-2.10.7.tgz" \
    | tar xzf - -C /usr/local --strip-components=1

CMD ["scala"]
