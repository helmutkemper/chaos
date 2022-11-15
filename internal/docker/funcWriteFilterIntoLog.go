package docker

import (
	"bytes"
	"github.com/helmutkemper/util"
	"log"
	"os"
	"regexp"
)

func (e *ContainerBuilder) writeFilterIntoLog(file *os.File, filter []LogFilter, lineList *[][]byte) (tab bool, err error) {
	var lineToFile = make([]byte, 0)
	var skipMatch = make([]bool, len(e.chaos.filterLog))

	for logLine := len(*lineList) - 1; logLine >= 0; logLine -= 1 {
		for filterLine := 0; filterLine != len(filter); filterLine += 1 {
			if skipMatch[filterLine] == true {
				continue
			}

			if filter[filterLine].Label == "" {
				continue
			}

			if bytes.Contains((*lineList)[logLine], []byte(filter[filterLine].Match)) == true {
				skipMatch[filterLine] = true

				var re *regexp.Regexp
				re, err = regexp.Compile(filter[filterLine].Filter)
				if err != nil {
					util.TraceToLog()
					log.Printf("regexp.Compile().error: %v", err)
					log.Printf("regexp.Compile().error.filter: %v", filter[filterLine].Filter)
					continue
				}

				var toFile []byte
				toFile = re.ReplaceAll((*lineList)[logLine], []byte("${valueToGet}"))

				if filter[filterLine].Search != "" {
					re, err = regexp.Compile(filter[filterLine].Search)
					if err != nil {
						util.TraceToLog()
						log.Printf("regexp.Compile().error: %v", err)
						log.Printf("regexp.Compile().error.filter: %v", filter[filterLine].Search)
						continue
					}

					toFile = re.ReplaceAll(toFile, []byte(filter[filterLine].Replace))
				}

				lineToFile = append(lineToFile, toFile...)
				lineToFile = append(lineToFile, []byte(e.csvValueSeparator)...)

			}
			tab = e.rowsToPrint&KLogColumnFilterLogComa != 0
		}
	}

	lineToFile = bytes.TrimSuffix(lineToFile, []byte(e.csvValueSeparator))
	_, err = file.Write(lineToFile)
	if err != nil {
		log.Printf("writeContainerLogToFile().error: %v", err.Error())
		util.TraceToLog()
		return
	}

	return
}
