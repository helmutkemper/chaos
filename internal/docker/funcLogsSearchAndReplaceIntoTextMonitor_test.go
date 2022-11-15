package docker

import (
	"bytes"
	"io/fs"
	"io/ioutil"
	"os"
	"testing"
)

func TestLogsSearchAndReplaceIntoTextMonitor(t *testing.T) {

	t.Cleanup(func() {
		var list []fs.FileInfo
		list, _ = ioutil.ReadDir("./test_log")
		for _, file := range list {
			_ = os.Remove("./test_log/" + file.Name())
		}
		_ = os.Remove("./test_log")
	})

	var err error
	var db = &ContainerBuilder{}
	var data []byte

	err = db.Init()
	if err != nil {
		t.FailNow()
	}

	err = db.AddMonitorMatchFlagToFileLog("texto monitorado", "./test_log")
	if err != nil {
		t.FailNow()
	}

	logs := make([]byte, 0)
	logs = append(logs, []byte("Linha 1\n")...)
	logs = append(logs, []byte("Linha 2\n")...)
	logs = append(logs, []byte("Linha 3\n")...)
	logs = append(logs, []byte("texto monitorado\n")...)
	logs = append(logs, []byte("Linha 4\n")...)
	logs = append(logs, []byte("Linha 5\n")...)
	db.logsSearchAndReplaceIntoTextMonitor(&logs, db.chaos.filterMonitor)
	logs = append(logs, []byte("Linha 7\n")...)
	logs = append(logs, []byte("Linha 8\n")...)
	logs = append(logs, []byte("Linha 9\n")...)
	logs = append(logs, []byte("Linha 10\n")...)
	db.logsSearchAndReplaceIntoTextMonitor(&logs, db.chaos.filterMonitor)
	logs = append(logs, []byte("Linha 11\n")...)
	logs = append(logs, []byte("Linha 12\n")...)
	logs = append(logs, []byte("texto monitorado\n")...)
	db.logsSearchAndReplaceIntoTextMonitor(&logs, db.chaos.filterMonitor)

	var list []fs.FileInfo
	list, err = ioutil.ReadDir("./test_log")
	if err != nil {
		t.FailNow()
	}

	if len(list) != 2 {
		t.FailNow()
	}

	data, err = ioutil.ReadFile("./test_log/log.0.log")
	equal := bytes.Compare(data, []byte("Linha 1\nLinha 2\nLinha 3\ntexto monitorado\nLinha 4\nLinha 5\n"))
	if equal != 0 {
		t.FailNow()
	}

	data, err = ioutil.ReadFile("./test_log/log.1.log")
	equal = bytes.Compare(data, []byte("Linha 7\nLinha 8\nLinha 9\nLinha 10\nLinha 11\nLinha 12\ntexto monitorado\n"))
	if equal != 0 {
		t.FailNow()
	}
}
