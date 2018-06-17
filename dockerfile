FROM ubuntu:latest

RUN apt-get update \
  && apt-get install -y python3-pip python3-dev libpng-dev libjpeg-dev golang-go \
  && cd /usr/local/bin \
  && ln -s /usr/bin/python3 python \
  && pip3 install --upgrade pillow python-slugify psutil pyqt5 raven KindleComicConverter
