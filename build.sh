#!/bin/bash

set -e

docker login -u="$DOCKER_USER" -p="$DOCKER_PASSWORD"

docker build -t laincloud/centos centos/7

docker push laincloud/centos

docker tag laincloud/centos laincloud/centos:7

docker push laincloud/centos:7

docker tag laincloud/centos  laincloud/centos:7.4.1708

docker push laincloud/centos:7.4.1708

docker build -t laincloud/debian:stretch debian/stretch

docker push laincloud/debian:stretch

docker build -t laincloud/golang:1.8.3 golang/1.8.3

docker push laincloud/golang:1.8.3

docker build -t laincloud/golang:1.9.0 golang/1.9.0

docker push laincloud/golang:1.9.0

docker build -t laincloud/golang-node:1.9.0-8.4.0 golang-node/1.9.0-8.4.0

docker push laincloud/golang-node:1.9.0-8.4.0

docker build -t laincloud/java:8 java/8-jdk

docker push laincloud/java:8

docker tag laincloud/java:8 laincloud/java:8-jdk

docker push laincloud/java:8-jdk

docker tag laincloud/java:8 laincloud/java:8u152

docker push laincloud/java:8u152

docker tag laincloud/java:8 laincloud/java:8u152-jdk

docker push laincloud/java:8u152-jdk

docker build -t laincloud/jetty:9 jetty/9

docker push laincloud/jetty:9

docker tag laincloud/jetty:9 laincloud/jetty:9.4

docker push laincloud/jetty:9.4

docker tag laincloud/jetty:9 laincloud/jetty:9.4.7

docker push laincloud/jetty:9.4.7

docker build -t laincloud/maven:3 maven/3

docker push laincloud/maven:3

docker tag laincloud/maven:3 laincloud/maven:3.5

docker push laincloud/maven:3.5

docker tag laincloud/maven:3 laincloud/maven:3.5.0

docker push laincloud/maven:3.5.0

docker build -t laincloud/maven-kotlin:3.5.0-1.1.3-2 maven-kotlin/3.5.0-1.1.3-2

docker push laincloud/maven-kotlin:3.5.0-1.1.3-2

docker build -t laincloud/nginx:1.13.5 nginx/1.13.5

docker push laincloud/nginx:1.13.5

docker build -t laincloud/node:8.4.0 node/8.4.0

docker push laincloud/node:8.4.0

docker build -t laincloud/python:3.6 python/3.6

docker push laincloud/python:3.6

#docker build -t laincloud/tomcat tomcat/8
#
#docker push laincloud/tomcat
#
#docker tag laincloud/tomcat laincloud/tomcat:8
#
#docker push laincloud/tomcat:8
#
#docker tag laincloud/tomcat laincloud/tomcat:8.5
#
#docker push laincloud/tomcat:8.5
#
#docker tag laincloud/tomcat laincloud/tomcat:8.5.24
#
#docker push laincloud/tomcat:8.5.24
