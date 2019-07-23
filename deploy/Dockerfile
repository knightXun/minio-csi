FROM centos

ENV CSIBIN=/usr/local/bin/csi

COPY csi $CSIBIN

RUN chmod +x $CSIBIN

RUN yum install -y wget curl net-tools fuse fuse-libs

ADD s3fs-fuse-1.85-1.el7.x86_64.rpm /root/

RUN rpm -ivh /root/s3fs-fuse-1.85-1.el7.x86_64.rpm

ENTRYPOINT ["/usr/local/bin/csi"]
