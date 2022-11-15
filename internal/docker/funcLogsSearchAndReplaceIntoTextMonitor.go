package docker

import (
	"bytes"
	"github.com/helmutkemper/util"
	"io/fs"
	"io/ioutil"
	"log"
	"strconv"
)

func (e *ContainerBuilder) logsSearchAndReplaceIntoTextMonitor(logs *[]byte, configuration []LogFilter) {
	var err error
	var dirList []fs.FileInfo
	var line []byte

	if configuration == nil {
		return
	}

	for filterLine := 0; filterLine != len(configuration); filterLine += 1 {
		size := len(*logs)
		changeSize := false
		// faz o log só lê a parte mais recente do mesmo
		logFiltered := (*logs)[configuration[filterLine].size:]

		logFiltered = bytes.ReplaceAll(logFiltered, []byte("\r"), []byte(""))
		lineList := bytes.Split(logFiltered, []byte("\n"))

		for logLine := len(lineList) - 1; logLine >= 0; logLine -= 1 {
			line = lineList[logLine]
			if configuration[filterLine].LogPath != "" && bytes.Contains(line, []byte(configuration[filterLine].Match)) == true {
				changeSize = true
				dirList, err = ioutil.ReadDir(configuration[filterLine].LogPath)
				if err != nil {
					log.Printf("ioutil.ReadDir().error: %v", err.Error())
					util.TraceToLog()
					return
				}
				var totalOfFiles = strconv.Itoa(len(dirList))
				err = ioutil.WriteFile(configuration[filterLine].LogPath+"log."+totalOfFiles+".log", (*logs)[configuration[filterLine].size:], fs.ModePerm)
				if err != nil {
					log.Printf("ioutil.WriteFile().error: %v", err.Error())
					util.TraceToLog()
					return
				}
				break
			}
		}

		if changeSize == true {
			configuration[filterLine].size = size
		}
	}

	return
}
