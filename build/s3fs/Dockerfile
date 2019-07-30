FROM centos
LABEL maintainers="Shulan tech"
LABEL description="Shulan tech minio csi"

COPY s3fs-fuse-1.85-1.el7.x86_64.rpm /root/s3fs-fuse-1.85-1.el7.x86_64.rpm
COPY nsenter /nsenter
COPY entrypoint.sh /entrypoint.sh
COPY csiplugin-connector /bin/csiplugin-connector
COPY csiplugin-connector.service /bin/csiplugin-connector.service
COPY minio-csi /bin/minio-csi
RUN chmod +x /bin/minio-csi && chmod +x /entrypoint.sh && chmod +x /bin/csiplugin-connector && chmod +x /bin/csiplugin-connector.service

ENTRYPOINT ["/entrypoint.sh"]