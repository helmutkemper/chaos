package docker

import (
	"github.com/docker/docker/api/types"
	"github.com/helmutkemper/util"
	"log"
	"os"
)

func (e *ContainerBuilder) writeConstBlkioIoServiceBytesRecursive(file *os.File, stats *types.Stats) (tab bool, err error) {
	if e.rowsToPrint&KLogColumnBlkioIoServiceBytesRecursive == KLogColumnBlkioIoServiceBytesRecursive {
		length := len(stats.BlkioStats.IoServiceBytesRecursive)
		for i := 0; i != length; i += 1 {
			_, err = file.Write([]byte("KBlkioIoServiceBytesRecursive"))
			if err != nil {
				log.Printf("writeContainerLogToFile().error: %v", err.Error())
				util.TraceToLog()
				return
			}

			_, err = file.Write([]byte(e.csvValueSeparator))
			if err != nil {
				log.Printf("writeContainerLogToFile().error: %v", err.Error())
				util.TraceToLog()
				return
			}

			_, err = file.Write([]byte("KBlkioIoServiceBytesRecursive"))
			if err != nil {
				log.Printf("writeContainerLogToFile().error: %v", err.Error())
				util.TraceToLog()
				return
			}

			_, err = file.Write([]byte(e.csvValueSeparator))
			if err != nil {
				log.Printf("writeContainerLogToFile().error: %v", err.Error())
				util.TraceToLog()
				return
			}

			_, err = file.Write([]byte("KBlkioIoServiceBytesRecursive"))
			if err != nil {
				log.Printf("writeContainerLogToFile().error: %v", err.Error())
				util.TraceToLog()
				return
			}

			_, err = file.Write([]byte(e.csvValueSeparator))
			if err != nil {
				log.Printf("writeContainerLogToFile().error: %v", err.Error())
				util.TraceToLog()
				return
			}

			_, err = file.Write([]byte("KBlkioIoServiceBytesRecursive"))
			if err != nil {
				log.Printf("writeContainerLogToFile().error: %v", err.Error())
				util.TraceToLog()
				return
			}

			if i != length-1 {
				_, err = file.Write([]byte(e.csvValueSeparator))
				if err != nil {
					log.Printf("writeContainerLogToFile().error: %v", err.Error())
					util.TraceToLog()
					return
				}
			}
		}

		if length > 0 {
			tab = e.rowsToPrint&KBlkioIoServiceBytesRecursiveComa != 0
		}
	}

	return
}
