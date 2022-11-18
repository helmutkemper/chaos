package iotmakerdocker

import (
	"bytes"
	"github.com/docker/docker/api/types"
	"io"
	"io/ioutil"
	"log"
	"strings"
	"sync"
	"time"
)

// ContainerLogsWaitText (English):
//
// ContainerLogsWaitText (Português):
func (el *DockerSystem) ContainerLogsWaitText(
	id string,
	text string,
	out io.Writer,
) (
	logContainer []byte,
	err error,
) {

	var wg sync.WaitGroup
	var reader io.ReadCloser
	var previousLog = make([]byte, 0)
	var cleanLog = make([]byte, 0)

	if out != nil {
		log.New(out, "", 0)
	}

	wg.Add(1)
	go func(el *DockerSystem, err *error, reader *io.ReadCloser, previousLog, cleanLog, logContainer *[]byte, text *string, id string) {
		defer wg.Done()

		for {
			*reader, *err = el.cli.ContainerLogs(el.ctx, id, types.ContainerLogsOptions{
				ShowStdout: true,
				ShowStderr: true,
				Timestamps: true,
				Follow:     false,
				Details:    false,
			})
			if *err != nil {
				return
			}

			*logContainer, *err = ioutil.ReadAll(*reader)
			if *err != nil {
				return
			}

			*cleanLog = bytes.Replace(*logContainer, *previousLog, []byte(""), -1)
			*previousLog = make([]byte, len(*logContainer))
			copy(*previousLog, *logContainer)

			//
			if out != nil && len(*cleanLog) != 0 {
				log.Printf("%s", *cleanLog)
			}

			if strings.Contains(string(*logContainer), *text) == true {
				return
			}

			time.Sleep(kWaitTextLoopSleep)
		}
	}(el, &err, &reader, &previousLog, &cleanLog, &logContainer, &text, id)
	wg.Wait()

	return
}
