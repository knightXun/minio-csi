package csicommon

import (
	"fmt"

	"github.com/container-storage-interface/spec/lib/go/csi"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"k8s.io/klog"
)

// CSIDriver stores driver information
type CSIDriver struct {
	name    string
	nodeID  string
	version string
	cap     []*csi.ControllerServiceCapability
	vc      []*csi.VolumeCapability_AccessMode
}

// NewCSIDriver Creates a NewCSIDriver object. Assumes vendor
// version is equal to driver version &  does not support optional
// driver plugin info manifest field. Refer to CSI spec for more details.
func NewCSIDriver(name, v, nodeID string) *CSIDriver {
	if name == "" {
		klog.Errorf("Driver name missing")
		return nil
	}

	if nodeID == "" {
		klog.Errorf("NodeID missing")
		return nil
	}
	// TODO version format and validation
	if v == "" {
		klog.Errorf("Version argument missing")
		return nil
	}

	driver := CSIDriver{
		name:    name,
		version: v,
		nodeID:  nodeID,
	}

	return &driver
}

// ValidateControllerServiceRequest validates the controller
// plugin capabilities
func (d *CSIDriver) ValidateControllerServiceRequest(c csi.ControllerServiceCapability_RPC_Type) error {
	if c == csi.ControllerServiceCapability_RPC_UNKNOWN {
		return nil
	}

	for _, cap := range d.cap {
		if c == cap.GetRpc().GetType() {
			return nil
		}
	}
	return status.Error(codes.InvalidArgument, fmt.Sprintf("%s", c)) //nolint
}

// AddControllerServiceCapabilities stores the controller capabilities
// in driver object
func (d *CSIDriver) AddControllerServiceCapabilities(cl []csi.ControllerServiceCapability_RPC_Type) {
	var csc []*csi.ControllerServiceCapability

	for _, c := range cl {
		klog.Infof("Enabling controller service capability: %v", c.String())
		csc = append(csc, NewControllerServiceCapability(c))
	}

	d.cap = csc

}

// AddVolumeCapabilityAccessModes stores volume access modes
func (d *CSIDriver) AddVolumeCapabilityAccessModes(vc []csi.VolumeCapability_AccessMode_Mode) []*csi.VolumeCapability_AccessMode {
	var vca []*csi.VolumeCapability_AccessMode
	for _, c := range vc {
		klog.Infof("Enabling volume access mode: %v", c.String())
		vca = append(vca, NewVolumeCapabilityAccessMode(c))
	}
	d.vc = vca
	return vca
}

// GetVolumeCapabilityAccessModes returns access modes
func (d *CSIDriver) GetVolumeCapabilityAccessModes() []*csi.VolumeCapability_AccessMode {
	return d.vc
}
