package minio

import (
	"k8s.io/utils/keymutex"

	"github.com/container-storage-interface/spec/lib/go/csi"
	"github.com/kubernetes-csi/csi-lib-utils/protosanitizer"
	"golang.org/x/net/context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"k8s.io/klog"
	"minio/minio-csi/pkg/csi-common"
)

const (
	oneGB = 1073741824
)

var (
	volumeNameMutex = keymutex.NewHashed(0)
	targetPathMutex = keymutex.NewHashed(0)
)

type ControllerServer struct {
	*csicommon.DefaultControllerServer
}

func (cs *ControllerServer) validateVolumeReq(req *csi.CreateVolumeRequest) error {
	if err := cs.Driver.ValidateControllerServiceRequest(csi.ControllerServiceCapability_RPC_CREATE_DELETE_VOLUME); err != nil {
		klog.V(3).Infof("invalid create volume req: %v", protosanitizer.StripSecrets(req))
		return err
	}
	// Check sanity of request Name, Volume Capabilities
	if req.Name == "" {
		return status.Error(codes.InvalidArgument, "volume Name cannot be empty")
	}
	if req.VolumeCapabilities == nil {
		return status.Error(codes.InvalidArgument, "volume Capabilities cannot be empty")
	}
	options := req.GetParameters()
	if value, ok := options["BucketName"]; !ok || value == "" {
		return status.Error(codes.InvalidArgument, "missing or empty BucketName to provision volume from")
	}
	if value, ok := options["MinioKey"]; !ok || value == "" {
		return status.Error(codes.InvalidArgument, "missing or empty MinioKey to provision volume from")
	}

	if value, ok := options["MinioURL"]; !ok || value == "" {
		return status.Error(codes.InvalidArgument, "missing or empty MinioURL to provision volume from")
	}

	return nil
}

func (cs *ControllerServer) parseVolCreateRequest(req *csi.CreateVolumeRequest) (*minioVolume, error) {
	// TODO (sbezverk) Last check for not exceeding total storage capacity
	return &minioVolume{}, nil
}

func (cs *ControllerServer) CreateVolume(ctx context.Context, req *csi.CreateVolumeRequest) (*csi.CreateVolumeResponse, error) {
	return &csi.CreateVolumeResponse{}, nil
}

func (cs *ControllerServer) DeleteVolume(ctx context.Context, req *csi.DeleteVolumeRequest) (*csi.DeleteVolumeResponse, error) {
	return &csi.DeleteVolumeResponse{}, nil
}

func (cs *ControllerServer) ValidateVolumeCapabilities(ctx context.Context, req *csi.ValidateVolumeCapabilitiesRequest) (*csi.ValidateVolumeCapabilitiesResponse, error) {
	return &csi.ValidateVolumeCapabilitiesResponse{}, nil
}

func (cs *ControllerServer) CreateSnapshot(ctx context.Context, req *csi.CreateSnapshotRequest) (*csi.CreateSnapshotResponse, error) {
	return &csi.CreateSnapshotResponse{}, nil
}

func (cs *ControllerServer) DeleteSnapshot(ctx context.Context, req *csi.DeleteSnapshotRequest) (*csi.DeleteSnapshotResponse, error) {
	return &csi.DeleteSnapshotResponse{}, nil
}
