package minio

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"minio/minio-csi/pkg/csi-common"

	"github.com/container-storage-interface/spec/lib/go/csi"
	"golang.org/x/net/context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"k8s.io/klog"
	"k8s.io/kubernetes/pkg/util/mount"
)

// NodeServer struct of ceph rbd driver with supported methods of CSI
// node server spec
type NodeServer struct {
	*csicommon.DefaultNodeServer
	mounter mount.Interface
}

// minioVolume represents a CSI volume
type minioVolume struct {
	minioBucket       string
	minioKey          string
	minioURL          string
}

// NodePublishVolume mounts the volume mounted to the device path to the target
// path
func (ns *NodeServer) NodePublishVolume(ctx context.Context, req *csi.NodePublishVolumeRequest) (*csi.NodePublishVolumeResponse, error) {

	targetPath := req.GetTargetPath()
	if targetPath == "" {
		return nil, status.Error(codes.InvalidArgument, "empty target path in request")
	}

	if req.GetVolumeCapability() == nil {
		return nil, status.Error(codes.InvalidArgument, "empty volume capability in request")
	}

	if req.GetVolumeId() == "" {
		return nil, status.Error(codes.InvalidArgument, "empty volume ID in request")
	}

	isBlock := req.GetVolumeCapability().GetBlock() != nil
	// Check if that target path exists properly
	notMnt, err := ns.createTargetPath(targetPath, isBlock)
	if err != nil {
		return nil, err
	}

	if !notMnt {
		return &csi.NodePublishVolumeResponse{}, fmt.Errorf("Already has a dir")
	}

	// Publish Path
	err = ns.mountVolume(req)
	if err != nil {
		return nil, err
	}

	return &csi.NodePublishVolumeResponse{}, nil
}

func genVolFromVolumeOptions(volOptions map[string]string) (*minioVolume, error) {
	var ok bool
	minioVol := &minioVolume{}

	minioVol.minioBucket, ok = volOptions["BucketName"]
	if !ok {
		return nil, fmt.Errorf("missing required parameter BucketName")
	}

	minioVol.minioKey, ok = volOptions["MinioKey"]
	if !ok {
		return nil, fmt.Errorf("missing required parameter MinioKey")
	}

	minioVol.minioURL, ok = volOptions["MinioURL"]
	if !ok {
		return nil, fmt.Errorf("missing required parameter MinioURL")
	}
	return minioVol, nil
}

func (ns *NodeServer) mountVolume(req *csi.NodePublishVolumeRequest) error {
	// Publish Path
	fsType := req.GetVolumeCapability().GetMount().GetFsType()
	readOnly := req.GetReadonly()
	attrib := req.GetVolumeContext()

	minioVol, err := genVolFromVolumeOptions(req.VolumeContext)

	if err != nil {
		return err
	}

	targetPath := req.GetTargetPath()

	klog.Infof("target %v\n BucketName %v\nfstype %v\nreadonly %v\nattributes %v\n MinioKey %v\n MinioURL%v\n",
		targetPath, minioVol.minioBucket, fsType, readOnly, attrib, minioVol.minioKey, minioVol.minioURL)

	tmpFile, err := ioutil.TempFile("/tmp", "minio")

	if err != nil {
		return fmt.Errorf("Create minio secret failed %v", err)
	}

	tmpFileName := tmpFile.Name()

	tmpFile.WriteString(minioVol.minioKey)

	tmpFile.Chmod(0600)

	tmpFile.Close()

	klog.Info("Exec s3fs : %v  %v", "/usr/bin/s3fs",  []string{"-o", "passwd_file="+tmpFileName,
		"-o", "url="+minioVol.minioURL,  "-o", "use_path_request_style", "-o", "bucket=" + minioVol.minioBucket ,
		targetPath, "-o", "curldbg", "-o", "no_check_certificate", "-o", "connect_timeout=5", "-o", "retries=1"} )


	_, err = execCommand("/usr/bin/s3fs", []string{"-o", "passwd_file="+tmpFileName,
		"-o", "url="+minioVol.minioURL,  "-o", "use_path_request_style", "-o", "bucket=" + minioVol.minioBucket ,
		targetPath, "-o", "curldbg", "-o", "no_check_certificate", "-o", "connect_timeout=5", "-o", "retries=1"})

	if err != nil {
		klog.Errorln("s3fs mount error: ", err.Error())
		return err
	}



	defer func() {
		_ = os.Remove(filepath.Join("/tmp", tmpFileName))
	}()
	return nil
}

func (ns *NodeServer) waitPathToShow(targetPath string) bool {
	time.Sleep(5)

	cmd := exec.Command("findmnt", "-n", "-o", "SOURCE", "--first-only", "--target", targetPath)
	out, err := cmd.CombinedOutput()
	if err != nil {
		klog.V(2).Infof("Failed findmnt command for path %s: %s %v", targetPath, out, err)
		return false
	}

	if string(out) == "s3fs" {
		return true
	}

	return false
}

func (ns *NodeServer) createTargetPath(targetPath string, isBlock bool) (bool, error) {
	// Check if that target path exists properly
	notMnt, err := mount.IsNotMountPoint(ns.mounter, targetPath)
	if err != nil {
		if os.IsNotExist(err) {
			if isBlock {
				// create an empty file
				// #nosec
				targetPathFile, e := os.OpenFile(targetPath, os.O_CREATE|os.O_RDWR, 0750)
				if e != nil {
					klog.V(4).Infof("Failed to create targetPath:%s with error: %v", targetPath, err)
					return notMnt, status.Error(codes.Internal, e.Error())
				}
				if err = targetPathFile.Close(); err != nil {
					klog.V(4).Infof("Failed to close targetPath:%s with error: %v", targetPath, err)
					return notMnt, status.Error(codes.Internal, err.Error())
				}
			} else {
				// Create a directory
				if err = os.MkdirAll(targetPath, 0750); err != nil {
					return notMnt, status.Error(codes.Internal, err.Error())
				}
			}
			notMnt = true
		} else {
			return false, status.Error(codes.Internal, err.Error())
		}
	}
	return notMnt, err

}

// NodeUnpublishVolume unmounts the volume from the target path
func (ns *NodeServer) NodeUnpublishVolume(ctx context.Context, req *csi.NodeUnpublishVolumeRequest) (*csi.NodeUnpublishVolumeResponse, error) {
	targetPath := req.GetTargetPath()

	if err := ns.unmount(targetPath); err != nil {
		return nil, err
	}

	return &csi.NodeUnpublishVolumeResponse{}, nil
}

func (ns *NodeServer) unmount(targetPath string) error {
	var err error

	err = ns.mounter.Unmount(targetPath)
	if err != nil {
		klog.V(3).Infof("failed to unmount targetPath: %s with error: %v", targetPath, err)
		return status.Error(codes.Internal, err.Error())
	}

	if err = os.RemoveAll(targetPath); err != nil {
		klog.V(3).Infof("failed to remove targetPath: %s with error: %v", targetPath, err)
	}
	return err
}
