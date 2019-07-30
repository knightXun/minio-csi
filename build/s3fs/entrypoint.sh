#!/bin/sh

echo "Starting deploy s3fs plugin...."
s3fsVer="1.85-1"

# install S3FS
mkdir -p /host/etc/csi-oss-tool/
if [ ! `/nsenter --mount=/proc/1/ns/mnt which s3fs` ]; then
    echo "First install s3fs...."
    cp /root/s3fs-fuse-${s3fsVer}.el7.x86_64.rpm /host/etc/
    /nsenter --mount=/proc/1/ns/mnt yum localinstall -y /etc/s3fs-fuse-${s3fsVer}.el7.x86_64.rpm
    rm -rf /host/etc/s3fs-fuse-${s3fsVer}.el7.x86_64.rpm
fi


updateConnector="true"
if [ ! -f "/host/etc/csi-s3fs-tool/csiplugin-connector" ];then
    mkdir -p /host/etc/csi-s3fs-tool/
    echo "mkdir /etc/csi-s3fs-tool/ directory..."
else
    oldmd5=`md5sum /host/etc/csi-s3fs-tool/csiplugin-connector | awk '{print $1}'`
    newmd5=`md5sum /bin/csiplugin-connector | awk '{print $1}'`
    if [ "$oldmd5" = "$newmd5" ]; then
        updateConnector="false"
    else
        rm -rf /host/etc/csi-s3fs-tool/
        rm -rf /host/etc/csi-s3fs-tool/connector.sock
        mkdir -p /host/etc/csi-s3fs-tool/
    fi
fi
if [ "$updateConnector" = "true" ]; then
    echo "copy csiplugin-connector...."
    cp /bin/csiplugin-connector /host/etc/csi-s3fs-tool/csiplugin-connector
    chmod 755 /host/etc/csi-s3fs-tool/csiplugin-connector
fi


# install csiplugin connector service
updateConnectorService="true"
if [ -f "/host/usr/lib/systemd/system/csiplugin-connector.service" ];then
    echo "prepare install csiplugin-connector.service...."
    oldmd5=`md5sum /host/usr/lib/systemd/system/csiplugin-connector.service | awk '{print $1}'`
    newmd5=`md5sum /bin/csiplugin-connector.service | awk '{print $1}'`
    if [ "$oldmd5" = "$newmd5" ]; then
        updateConnectorService="false"
    else
        rm -rf /host/usr/lib/systemd/system/csiplugin-connector.service
    fi
fi

if [ "$updateConnectorService" = "true" ]; then
    echo "install csiplugin-connector...."
    cp /bin/csiplugin-connector.service /host/usr/lib/systemd/system/csiplugin-connector.service
fi

#/nsenter --mount=/proc/1/ns/mnt service csiplugin-connector-svc restart
rm -rf /var/log/alicloud/connector.pid
/nsenter --mount=/proc/1/ns/mnt systemctl enable csiplugin-connector.service
/nsenter --mount=/proc/1/ns/mnt systemctl restart csiplugin-connector.service

# start daemon
/bin/minio-csi $@
