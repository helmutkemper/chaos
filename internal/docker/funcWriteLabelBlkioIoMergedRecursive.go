package docker

import (
	"github.com/docker/docker/api/types"
	"github.com/helmutkemper/util"
	"log"
	"os"
)

func (e *ContainerBuilder) writeLabelBlkioIoMergedRecursive(file *os.File, stats *types.Stats) (tab bool, err error) {
	if e.rowsToPrint&KLogColumnBlkioIoMergedRecursive == KLogColumnBlkioIoMergedRecursive {
		length := len(stats.BlkioStats.IoMergedRecursive)
		for i := 0; i != length; i += 1 {
			_, err = file.Write([]byte("BlkioStats stores All IO service stats for data read and write. Major. Io Merged Recursive."))
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

			_, err = file.Write([]byte("BlkioStats stores All IO service stats for data read and write. Minor. Io Merged Recursive."))
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

			_, err = file.Write([]byte("BlkioStats stores All IO service stats for data read and write. Op. Io Merged Recursive."))
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

			_, err = file.Write([]byte("BlkioStats stores All IO service stats for data read and write. Value. Io Merged Recursive."))
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
			tab = e.rowsToPrint&KBlkioIoMergedRecursiveComa != 0
		}
	}

	return
}
