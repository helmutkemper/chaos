package docker

import (
	"github.com/docker/docker/api/types"
	"github.com/helmutkemper/util"
	"log"
	"os"
	"strconv"
)

func (e *ContainerBuilder) writeBlkioIoServiceTimeRecursive(file *os.File, stats *types.Stats) (tab bool, err error) {
	if e.rowsToPrint&KLogColumnBlkioIoServiceTimeRecursive == KLogColumnBlkioIoServiceTimeRecursive {
		length := len(stats.BlkioStats.IoServiceTimeRecursive)
		for i := 0; i != length; i += 1 {
			_, err = file.Write([]byte(strconv.FormatUint(stats.BlkioStats.IoServiceTimeRecursive[i].Major, 10)))
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

			_, err = file.Write([]byte(strconv.FormatUint(stats.BlkioStats.IoServiceTimeRecursive[i].Minor, 10)))
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

			_, err = file.Write([]byte(stats.BlkioStats.IoServiceTimeRecursive[i].Op))
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

			_, err = file.Write([]byte(strconv.FormatUint(stats.BlkioStats.IoServiceTimeRecursive[i].Value, 10)))
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
			tab = e.rowsToPrint&KBlkioIoServiceTimeRecursiveComa != 0
		}
	}

	return
}
