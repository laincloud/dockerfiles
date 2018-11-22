FROM laincloud/gradle:4.7.0-jdk8

# Download and untar SDK
ENV ANDROID_SDK_URL https://dl.google.com/android/repository/sdk-tools-linux-3859397.zip
RUN wget $ANDROID_SDK_URL 
RUN unzip sdk-tools-linux*.zip -d /usr/local/android_sdk

ENV ANDROID_HOME /usr/local/android_sdk
ENV ANDROID_SDK /usr/local/android_sdk
ENV PATH ${ANDROID_HOME}/tools:$ANDROID_HOME/platform-tools:${ANDROID_HOME}/tools/bin:$PATH

# Install Android SDK components
RUN yes | sdkmanager "extras;android;m2repository" "extras;google;m2repository" "platforms;android-28" "build-tools;28.0.3" "platform-tools" | echo yes 

# Support Gradle
ENV TERM dumb
ENV JAVA_OPTS -Xms256m -Xmx512m