package iotmakerdocker

// NewImagePullStatusChannel (English): Prepare a new channel for pull/build data
//
// NewImagePullStatusChannel (Português): Prepara um canal para os dados de pull/build
func NewImagePullStatusChannel() (
	chanPointer *chan ContainerPullStatusSendToChannel,
) {

	channel := make(chan ContainerPullStatusSendToChannel, 1)
	chanPointer = &channel
	return
}
