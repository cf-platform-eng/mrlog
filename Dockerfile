FROM harbor-repo.vmware.com/dockerhub-proxy-cache/library/ubuntu
LABEL maintainer="Pivotal Platform Engineering ISV-CI Team <cf-isv-dashboard@pivotal.io>"

COPY build/mrlog-linux /usr/local/bin/mrlog

ENTRYPOINT [ "mrlog" ]