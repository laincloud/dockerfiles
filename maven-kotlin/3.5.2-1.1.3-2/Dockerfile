FROM laincloud/maven:3.5.2

ENV KOTLIN_VERSION=1.1.3-2 \
    KOTLIN_HOME=/opt/kotlin

# Install kotlin
RUN cd /tmp \
    && wget -q -k "https://github.com/JetBrains/kotlin/releases/download/v${KOTLIN_VERSION}/kotlin-compiler-${KOTLIN_VERSION}.zip" \
    && unzip "kotlin-compiler-${KOTLIN_VERSION}.zip" \
    && mv "/tmp/kotlinc" "${KOTLIN_HOME}" \
    && rm "${KOTLIN_HOME}"/bin/*.bat \
    && chmod +x ${KOTLIN_HOME}/bin/* \
    && ln -s "${KOTLIN_HOME}/bin/"* "/usr/bin/" \
    && rm -rf /tmp/*