FROM harbor-repo.vmware.com/dockerhub-proxy-cache/library/ubuntu
LABEL maintainer="Tanzu ISV Partner Engineering Team <tanzu-isv-engineering@groups.vmware.com>"

COPY build/mrlog-linux /usr/local/bin/mrlog

ENTRYPOINT [ "mrlog" ]