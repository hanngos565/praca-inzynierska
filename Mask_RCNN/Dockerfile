#MASK-RCNN
FROM ubuntu:18.04 as builder
MAINTAINER Hanna Gościniak <hanngos10@gmail.com>

ARG DEBIAN_FRONTEND=noninteractive

RUN apt-get update && \
    apt-get install -y build-essential && \
    apt-get install --no-install-recommends -y \
    python3-dev python3.5 python3-pip

RUN pip3 install --upgrade pip setuptools

RUN apt-get install -y python3-opencv

COPY requirements.txt .
RUN pip3 --no-cache-dir install -r requirements.txt && \
    rm requirements.txt

RUN apt-get install -y vim

RUN mkdir app
COPY /app /app
WORKDIR /app
ENTRYPOINT python3 app.py