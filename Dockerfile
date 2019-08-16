FROM ubuntu:xenial
LABEL maintainer="Pivotal Platform Engineering ISV-CI Team <cf-isv-dashboard@pivotal.io>"

COPY build/mrlog-linux /usr/local/bin/mrlog

ENTRYPOINT [ "mrlog" ]