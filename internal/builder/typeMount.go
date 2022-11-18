package builder

// MountType:
//
//	KVolumeMountTypeBind - TypeBind is the type for mounting host dir
//	KVolumeMountTypeVolume - TypeVolume is the type for remote storage volumes
//	KVolumeMountTypeTmpfs - TypeTmpfs is the type for mounting tmpfs
//	KVolumeMountTypeNpipe - TypeNamedPipe is the type for mounting Windows named pipes
//
// Source: relative file/dir path in computer
// Destination: full path inside container
type Mount struct {
	MountType   VolumeMountType
	Source      string
	Destination string
}
