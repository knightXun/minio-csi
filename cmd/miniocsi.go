package main

import (
	"flag"
	"k8s.io/klog"
	"minio/minio-csi/pkg/minio"
)

const (
	minioDefaultName = "s3fs.csi.minio.com"
)


var (
	// common flags
	endpoint   = flag.String("endpoint", "unix://tmp/csi.sock", "CSI endpoint")
	nodeID     = flag.String("nodeid", "", "node id")
	instanceID = flag.String("instanceid", "", "Unique ID distinguishing this instance of Ceph CSI among other"+
		" instances, when sharing Ceph clusters across CSI instances for provisioning")

	// rbd related flags
	containerized = flag.Bool("containerized", true, "whether run as containerized")
)

func init() {
	klog.InitFlags(nil)
	if err := flag.Set("logtostderr", "true"); err != nil {
		klog.Exitf("failed to set logtostderr flag: %v", err)
	}
	flag.Parse()
}

func main() {
	minio.PluginFolder += "s3fs"
	driver := minio.NewDriver()
	driver.Run("s3fs.csi.minio.com", *nodeID, *endpoint, *instanceID, *containerized)
}

