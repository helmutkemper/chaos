package docker

import (
	"github.com/docker/docker/api/types"
	"github.com/helmutkemper/util"
	"io/fs"
	"log"
	"os"
)

// writeContainerLogToFile
//
// Português: Escreve um arquivo csv com dados capturados da saída padrão do container e dados estatísticos do container
func (e *ContainerBuilder) writeContainerLogToFile(path string, lineList [][]byte) (err error) {
	if path == "" {
		return
	}

	if lineList == nil {
		return
	}

	var makeLabel = false
	_, err = os.Stat(path)
	if err != nil {
		makeLabel = true
	}

	var file *os.File
	file, err = os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, fs.ModePerm)
	if err != nil {
		log.Printf("writeContainerLogToFile().error: %v", err.Error())
		util.TraceToLog()
		return
	}

	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Printf("writeContainerLogToFile().error: %v", err.Error())
			util.TraceToLog()
		}
	}(file)

	var stats = types.Stats{}
	stats, err = e.ContainerStatisticsOneShot()
	if err != nil {
		log.Printf("writeContainerLogToFile().error: %v", err.Error())
		util.TraceToLog()
		return
	}

	if makeLabel == true && e.csvConstHeader == true {
		err = e.writeContainerConstToFile(file, &stats)
		if err != nil {
			log.Printf("writeContainerLogToFile().error: %v", err.Error())
			util.TraceToLog()
			return
		}

		_, err = file.Write([]byte(e.csvRowSeparator))
		if err != nil {
			log.Printf("writeContainerLogToFile().error: %v", err.Error())
			util.TraceToLog()
			return
		}
	}

	if makeLabel == true {
		err = e.writeContainerLabelToFile(file, &stats)
		if err != nil {
			log.Printf("writeContainerLogToFile().error: %v", err.Error())
			util.TraceToLog()
			return
		}

		_, err = file.Write([]byte(e.csvRowSeparator))
		if err != nil {
			log.Printf("writeContainerLogToFile().error: %v", err.Error())
			util.TraceToLog()
			return
		}
	}

	err = e.writeContainerStatsToFile(file, &stats, &lineList)
	if err != nil {
		log.Printf("writeContainerLogToFile().error: %v", err.Error())
		util.TraceToLog()
		return
	}

	_, err = file.Write([]byte(e.csvRowSeparator))
	if err != nil {
		log.Printf("writeContainerLogToFile().error: %v", err.Error())
		util.TraceToLog()
		return
	}

	return
}
