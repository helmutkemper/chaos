package docker

type VolumeMountType int

const (
	// KVolumeMountTypeBindString TypeBind is the type for mounting host dir (real folder inside computer where this code work)
	KVolumeMountTypeBindString = "bind"

	// KVolumeMountTypeVolumeString TypeVolume is the type for remote storage volumes
	KVolumeMountTypeVolumeString = "volume"

	// KVolumeMountTypeTmpfsString TypeTmpfs is the type for mounting tmpfs
	KVolumeMountTypeTmpfsString = "tmpfs"

	// KVolumeMountTypeNpipeString TypeNamedPipe TypeNamedPipe is the type for mounting Windows named pipes
	KVolumeMountTypeNpipeString = "npipe"
)

const (
	// KVolumeMountTypeBind TypeBind is the type for mounting host dir (real folder inside computer where this code work)
	KVolumeMountTypeBind VolumeMountType = iota

	// KVolumeMountTypeVolume TypeVolume is the type for remote storage volumes
	KVolumeMountTypeVolume

	// KVolumeMountTypeTmpfs TypeTmpfs is the type for mounting tmpfs
	KVolumeMountTypeTmpfs

	// KVolumeMountTypeNpipe TypeNamedPipe is the type for mounting Windows named pipes
	KVolumeMountTypeNpipe
)

func (el VolumeMountType) String() string {
	return volumeMountTypes[el]
}

var volumeMountTypes = [...]string{
	"bind",
	"volume",
	"tmpfs",
	"npipe",
}
