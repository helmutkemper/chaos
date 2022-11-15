package docker

import (
	"bytes"
	"github.com/helmutkemper/util"
	"io/fs"
	"io/ioutil"
	"log"
	"regexp"
	"strconv"
)

func (e *ContainerBuilder) logsSearchAndReplaceIntoText(logs *[]byte, lineList [][]byte, configuration []LogFilter) (line []byte, found bool) {
	var err error
	var dirList []fs.FileInfo

	if configuration == nil {
		return
	}

	for logLine := len(lineList) - 1; logLine >= 0; logLine -= 1 {

		for filterLine := 0; filterLine != len(configuration); filterLine += 1 {
			line = lineList[logLine]
			if bytes.Contains(line, []byte(configuration[filterLine].Match)) == true {

				if configuration[filterLine].LogPath != "" {
					dirList, err = ioutil.ReadDir(configuration[filterLine].LogPath)
					if err != nil {
						log.Printf("ioutil.ReadDir().error: %v", err.Error())
						util.TraceToLog()
						return
					}
					var totalOfFiles = strconv.Itoa(len(dirList))
					err = ioutil.WriteFile(configuration[filterLine].LogPath+"log."+totalOfFiles+".log", *logs, fs.ModePerm)
					if err != nil {
						log.Printf("ioutil.WriteFile().error: %v", err.Error())
						util.TraceToLog()
						return
					}
				}

				if configuration[filterLine].Filter != "" {

					var re *regexp.Regexp
					re, err = regexp.Compile(configuration[filterLine].Filter)
					if err != nil {
						util.TraceToLog()
						log.Printf("regexp.Compile().error: %v", err)
						log.Printf("regexp.Compile().error.filter: %v", configuration[filterLine].Filter)
						continue
					}

					line = re.ReplaceAll(lineList[logLine], []byte("${valueToGet}"))

					if configuration[filterLine].Search != "" {
						re, err = regexp.Compile(configuration[filterLine].Search)
						if err != nil {
							util.TraceToLog()
							log.Printf("regexp.Compile().error: %v", err)
							log.Printf("regexp.Compile().error.filter: %v", configuration[filterLine].Search)
							continue
						}

						line = re.ReplaceAll(line, []byte(configuration[filterLine].Replace))
					}
				}

				found = true
				return
			}
		}
	}

	return
}
