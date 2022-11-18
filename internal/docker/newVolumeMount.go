package iotmakerdocker

import (
	"errors"
	"github.com/docker/docker/api/types/mount"
	"github.com/helmutkemper/iotmaker.docker/util"
)

func NewVolumeMount(
	list []Mount,
) (
	mountVolumesList []mount.Mount,
	err error,
) {

	mountVolumesList = make([]mount.Mount, 0)

	var found bool
	var fileAbsolutePath string

	for _, v := range list {
		found = util.VerifyFileExists(v.Source)
		if found == false {
			err = errors.New(v.Source + ": source file not found")
			return
		}

		err, fileAbsolutePath = util.FileGetAbsolutePath(v.Source)
		if err != nil {
			return
		}

		mountVolumesList = append(
			mountVolumesList,
			mount.Mount{
				Type:   mount.Type(v.MountType.String()),
				Source: fileAbsolutePath,
				Target: v.Destination,
			},
		)
	}

	return
}
