package minio

// ErrImageNotFound is returned when image name is not found in the cluster on the given pool
type ErrImageNotFound struct {
	imageName string
	err       error
}

func (e ErrImageNotFound) Error() string {
	return e.err.Error()
}

// ErrSnapNotFound is returned when snap name passed is not found in the list of snapshots for the
// given image
type ErrSnapNotFound struct {
	snapName string
	err      error
}

func (e ErrSnapNotFound) Error() string {
	return e.err.Error()
}

// ErrVolNameConflict is generated when a requested CSI volume name already exists on RBD but with
// different properties, and hence is in conflict with the passed in CSI volume name
type ErrVolNameConflict struct {
	requestName string
	err         error
}

func (e ErrVolNameConflict) Error() string {
	return e.err.Error()
}
