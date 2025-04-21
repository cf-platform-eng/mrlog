ARG ubuntu_image=tas-ecosystem-docker-virtual.usw1.packages.broadcom.com/ubuntu

FROM ${ubuntu_image}

LABEL maintainer="Tanzu ISV Partner Engineering Team <tanzu-isv-engineering@groups.vmware.com>"

COPY mrlog-build/mrlog-linux /usr/local/bin/mrlog

ENTRYPOINT [ "mrlog" ]