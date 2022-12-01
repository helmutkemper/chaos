package builder

import (
	"encoding/json"
	"fmt"
	"io"
	"regexp"
	"strings"
)

type auxId struct {
	ID string `json:"ID"`
}

type aux struct {
	Aux auxId `json:"aux"`
}

func (el *DockerSystem) processBuildAndPullReaders(
	reader *io.Reader,
	channel chan ContainerPullStatusSendToChannel,
) (
	successfully bool,
	err error,
) {

	var imageName string
	var imageId string
	var bufferReader = make([]byte, 1)
	var bufferDataInput = make([]byte, 0)
	var channelOut ContainerPullProgress
	var toChannel ContainerPullStatusSendToChannel
	var toProcess = make(map[string]ContainerPullProgress)
	var auxIdList = make([]string, 0)

	if reader == nil {
		el.imagePullWriteChannel(channel, ContainerPullStatusSendToChannel{Closed: true})
		return
	}

	if *reader == nil {
		el.imagePullWriteChannel(channel, ContainerPullStatusSendToChannel{Closed: true})
		return
	}

	for {
		_, err = (*reader).Read(bufferReader)
		if err != nil {
			if err.Error() == "EOF" {
				err = nil

				//>>>>> send to channel
				toChannel.calcPercentage()

				if imageName != "" {
					toChannel.ImageName = imageName
				}

				if imageId != "" {
					toChannel.ImageID = imageId
				} else if imageName != "" {
					toChannel.ImageID, err = el.ImageFindIdByName(imageName)
					if err != nil {
						el.imagePullWriteChannel(channel, ContainerPullStatusSendToChannel{Closed: true})
						return
					}

					channelOut.SuccessfullyBuildImage = true
					successfully = true
				}

				if len(auxIdList) != 0 {
					channelOut.SuccessfullyBuildImage = true
					successfully = true
					imageId, auxIdList = auxIdList[len(auxIdList)-1], auxIdList[:len(auxIdList)-1]
					toChannel.SetAuxiliaryImageList(auxIdList)
				}

				toChannel.Closed = true
				toChannel.SuccessfullyBuildImage = channelOut.SuccessfullyBuildImage
				toChannel.SuccessfullyBuildContainer = channelOut.SuccessfullyBuildContainer
				el.imagePullWriteChannel(channel, toChannel)

				return
			}

			el.imagePullWriteChannel(channel, ContainerPullStatusSendToChannel{Closed: true})
			return
		}

		bufferDataInput = append(bufferDataInput, bufferReader[0])

		if bufferReader[0] == byte(0x0A) {
			err = json.Unmarshal(bufferDataInput, &channelOut)
			r := regexp.MustCompile("successful")
			if r.Match(bufferDataInput) == true {
				fmt.Printf("-- %v --", bufferDataInput)
			}

			bufferDataInput = make([]byte, 0)

			if strings.Contains(channelOut.Stream, kContainerBuildImageStatusSuccessContainer) {
				channelOut.SysStatus = KContainerPullStatusComplete
				channelOut.SuccessfullyBuildContainer = true
				successfully = true

			} else if channelOut.Stream != "" {
				channelOut.SysStatus = KContainerPullStatusBuilding
			}

			if strings.Contains(channelOut.Stream, kContainerBuildImageStatusSuccessImage) {
				channelOut.SysStatus = KContainerPullStatusComplete
				channelOut.SuccessfullyBuildImage = true
				successfully = true

				imageName = strings.Replace(channelOut.Stream, kContainerBuildImageStatusSuccessImage, "", 1)
				imageName = strings.Replace(imageName, "\n", "", -1)
				imageName = strings.Replace(imageName, "\r", "", -1)
				imageName = strings.TrimSpace(imageName)

			} else if channelOut.Stream != "" {
				channelOut.SysStatus = KContainerPullStatusBuilding
			}

			//Successfully tagged delete_remote_server:latest
			if strings.Contains(channelOut.Status, kContainerPullStatusDownloadedNewerImageText) {
				imageName = strings.Replace(channelOut.Status, kContainerPullStatusDownloadedNewerImageText, "", 1)
			}

			if strings.Contains(channelOut.Status, kContainerPullStatusAuxId) { // {"aux":{"ID":"sha256:262c77a02e05b41efe2097f62f0e687b323d140ca72948a85bd5a4d7dc50e483"}} {"aux":{"ID":"sha256:bc032e1e78666df7d8d084c18e35752ac821942f5c8ddfe0a790afec33d13eb2"}}
				var aux aux
				err = json.Unmarshal([]byte(channelOut.Status), &aux)
				if err != nil {
					el.imagePullWriteChannel(channel, ContainerPullStatusSendToChannel{Closed: true})
					return
				}

				auxIdList = append(auxIdList, aux.Aux.ID)
			}

			if strings.Contains(channelOut.Status, kContainerPullStatusDigestText) {
				imageId = strings.Replace(channelOut.Status, kContainerPullStatusDigestText, "", 1)
			}

			if strings.Contains(channelOut.Status, kContainerPullStatusPullCompleteText) {
				channelOut.SysStatus = KContainerPullStatusPullComplete
			}

			if strings.Contains(channelOut.Status, kContainerPullStatusExtractingText) {
				channelOut.SysStatus = KContainerPullStatusExtracting
			}

			if strings.Contains(channelOut.Status, kContainerPullStatusWaitingText) {
				channelOut.SysStatus = KContainerPullStatusWaiting
			}

			if strings.Contains(channelOut.Status, kContainerPullStatusDownloadingText) {
				channelOut.SysStatus = KContainerPullStatusDownloading
			}

			if strings.Contains(channelOut.Status, kContainerPullStatusVerifyingChecksumText) {
				channelOut.SysStatus = KContainerPullStatusVerifyingChecksum
			}

			if strings.Contains(channelOut.Status, kContainerPullStatusDownloadCompleteText) {
				channelOut.SysStatus = KContainerPullStatusDownloadComplete
			}

			if strings.Contains(channelOut.Status, kContainerPullStatusImageIsUpToDate) {
				imageName = strings.Replace(channelOut.Status, kContainerPullStatusImageIsUpToDate, "", 1)
				channelOut.SuccessfullyBuildImage = true
				successfully = true
			}

			toProcess[channelOut.ID] = channelOut

			for _, v := range toProcess {
				if v.SysStatus == KContainerPullStatusPullComplete {
					toChannel.PullComplete += 1
				} else if v.SysStatus == KContainerPullStatusExtracting {
					toChannel.Extracting.Count += 1
					toChannel.Extracting.Total += v.ProgressDetail.Total
					toChannel.Extracting.Current += v.ProgressDetail.Current
				} else if v.SysStatus == KContainerPullStatusWaiting {
					toChannel.Waiting += 1
				} else if v.SysStatus == KContainerPullStatusDownloading {
					toChannel.Downloading.Count += 1
					toChannel.Downloading.Total += v.ProgressDetail.Total
					toChannel.Downloading.Current += v.ProgressDetail.Current
				} else if v.SysStatus == KContainerPullStatusVerifyingChecksum {
					toChannel.VerifyingChecksum += 1
				} else if v.SysStatus == KContainerPullStatusDownloadComplete {
					toChannel.DownloadComplete += 1
				}
			}

			toChannel.calcPercentage()

			if imageName != "" {
				toChannel.ImageName = imageName
			}

			if imageId != "" {
				toChannel.ImageID = imageId
			}

			toChannel.Stream = TerminalToHtml(channelOut.Stream)
			toChannel.SuccessfullyBuildImage = channelOut.SuccessfullyBuildImage
			toChannel.SuccessfullyBuildContainer = channelOut.SuccessfullyBuildContainer

			//>>>>> send to channel
			el.imagePullWriteChannel(channel, toChannel)

			toChannel = ContainerPullStatusSendToChannel{}
		}
	}
}
