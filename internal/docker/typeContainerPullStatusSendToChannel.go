package docker

import (
	"github.com/helmutkemper/iotmaker.docker/util"
	"math"
)

type ContainerPullStatusSendToChannel struct {
	Waiting                    int
	Downloading                ContainerPullStatusSendToChannelCount
	VerifyingChecksum          int
	DownloadComplete           int
	Extracting                 ContainerPullStatusSendToChannelCount
	PullComplete               int
	ImageName                  string
	ImageID                    string
	ContainerID                string
	Closed                     bool
	Stream                     string
	SuccessfullyBuildContainer bool
	SuccessfullyBuildImage     bool
	IdAuxiliaryImages          []string
}

func (el *ContainerPullStatusSendToChannel) calcPercentage() {
	var percent = float64(el.Downloading.Current) / float64(el.Downloading.Total) * 100.0
	if math.IsNaN(percent) == true {
		percent = 0.0
	}
	percent = util.Round(percent, 0.5, 2.0)
	el.Downloading.Percent = percent

	percent = float64(el.Extracting.Current) / float64(el.Extracting.Total) * 100.0
	if math.IsNaN(percent) == true {
		percent = 0.0
	}
	percent = util.Round(percent, 0.5, 2.0)
	el.Extracting.Percent = percent
}

func (el *ContainerPullStatusSendToChannel) SetAuxiliaryImageList(list []string) {
	el.IdAuxiliaryImages = list
}
