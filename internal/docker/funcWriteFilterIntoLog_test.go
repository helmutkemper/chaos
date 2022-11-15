package docker

import (
	"bytes"
	"github.com/helmutkemper/util"
	"io/ioutil"
	"log"
	"os"
	"testing"
)

func TestContainerBuilder_writeFilterIntoLog(t *testing.T) {
	var err error
	var data []byte
	var lineList [][]byte
	var file *os.File

	file, err = os.OpenFile("./test.file.must.be.deleted.txt", os.O_CREATE|os.O_RDWR|os.O_APPEND, 0666)
	if err != nil {
		log.Printf("os.OpenFile().error: %v", err)
		util.TraceToLog()
		t.FailNow()
	}

	data, err = ioutil.ReadFile("./container.output.test.txt")
	if err != nil {
		log.Printf("ioutil.ReadFile().error: %v", err)
		util.TraceToLog()
		t.FailNow()
	}

	e := &ContainerBuilder{}
	lineList = e.logsCleaner(data)

	_, err = e.writeFilterIntoLog(
		file,
		[]LogFilter{
			{
				Label:   "Contador",
				Match:   "counter",
				Filter:  "^.*counter: (?P<valueToGet>[0-9]+)",
				Search:  "\\.",
				Replace: ",",
				LogPath: "",
			},
		},
		&lineList,
	)
	if err != nil {
		log.Printf("writeContainerLogToFile().error: %v", err.Error())
		util.TraceToLog()
		t.FailNow()
	}

	_ = file.Close()

	data, err = ioutil.ReadFile("./test.file.must.be.deleted.txt")
	if err != nil {
		log.Printf("os.OpenFile().error: %v", err)
		util.TraceToLog()
		t.FailNow()
	}

	if bytes.Compare(data, []byte{54, 52, 44, 48}) != 0 {
		log.Printf("result must be a 64.0 float")
		log.Printf("%v", data)
		t.FailNow()
	}

	_ = os.Remove("./test.file.must.be.deleted.txt")
}
