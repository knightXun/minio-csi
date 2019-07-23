package minio

import (
	"context"

	csicommon "minio/minio-csi/pkg/csi-common"

	"github.com/container-storage-interface/spec/lib/go/csi"
)

// IdentityServer struct of rbd CSI driver with supported methods of CSI
// identity server spec.
type IdentityServer struct {
	*csicommon.DefaultIdentityServer
}

// GetPluginCapabilities returns available capabilities of the rbd driver
func (is *IdentityServer) GetPluginCapabilities(ctx context.Context, req *csi.GetPluginCapabilitiesRequest) (*csi.GetPluginCapabilitiesResponse, error) {
	return &csi.GetPluginCapabilitiesResponse{
		Capabilities: []*csi.PluginCapability{
			{
				Type: &csi.PluginCapability_Service_{
					Service: &csi.PluginCapability_Service{
						Type: csi.PluginCapability_Service_CONTROLLER_SERVICE,
					},
				},
			},
		},
	}, nil
}
