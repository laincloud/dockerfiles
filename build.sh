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

docker build -t laincloud/debian:jessie debian/jessie

docker push laincloud/debian:jessie

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

docker build -t laincloud/jira:6 jira/6

docker push laincloud/jira:6

docker tag laincloud/jira:6 laincloud/jira:6.4.14

docker push laincloud/jira:6.4.14

docker build -t laincloud/jira:7 jira/7

docker push laincloud/jira:7

docker tag laincloud/jira:7 laincloud/jira:7.6.1

docker push laincloud/jira:7.6.1

docker build -t laincloud/openjdk:7-jre openjdk/7-jre

docker push laincloud/openjdk:7-jre

docker tag laincloud/openjdk:7-jre laincloud/openjdk:7u151-jre

docker push laincloud/openjdk:7u151-jre

docker build -t laincloud/openjdk:7-jdk openjdk/7-jdk

docker push laincloud/openjdk:7-jdk

docker tag laincloud/openjdk:7-jdk laincloud/7u151-jdk

docker push laincloud/7u151-jdk

docker build -t laincloud/jetty:9.2-jre7 jetty/9.2-jre7

docker push laincloud/jetty:9.2-jre7
