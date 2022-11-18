package builder

func (el *DockerSystem) imagePullWriteChannel(
	progressChannel *chan ContainerPullStatusSendToChannel,
	data ContainerPullStatusSendToChannel,
) {

	if *progressChannel == nil {
		return
	}

	l := len(*progressChannel)
	if l != 0 {
		return
	}

	*progressChannel <- data
}
