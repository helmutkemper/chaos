package docker

import (
	dockerfileGolang "github.com/helmutkemper/iotmaker.docker.builder.golang.dockerfile"
	"github.com/helmutkemper/util"
	"runtime"
	"time"
)

// Init
//
// English:
//
//	Initializes the object.
//
//	 Output:
//	   err: Standard error object
//
// Note:
//
//   - Should be called only after all settings have been configured
//
// Português:
//
//	Inicializa o objeto.
//
//	 Saída:
//	   err: Objeto de erro padrão
//
// Nota:
//
//   - Deve ser chamado apenas depois de toas as configurações serem definidas
func (e *ContainerBuilder) Init() (err error) {
	e.init = true

	var osName = runtime.GOOS
	if e.rowsToPrint == 0 && osName == "darwin" {
		e.rowsToPrint = KLogColumnMacOs
	} else if e.rowsToPrint == 0 && osName == "windows" {
		e.rowsToPrint = KLogColumnWindows
	} else if e.rowsToPrint == 0 {
		e.rowsToPrint = KLogColumnAll
	}

	if e.sshDefaultFileName == "" {
		e.sshDefaultFileName = "sshDefaultFileName"
	}

	e.chaos.event = make(chan Event, 1)

	if e.metadata == nil {
		e.metadata = make(map[string]interface{})
	}

	if e.csvValueSeparator == "" {
		e.csvValueSeparator = ","
	}

	if e.csvRowSeparator == "" {
		e.csvRowSeparator = "\n"
	}

	if e.imageCacheName == "" {
		e.imageCacheName = "cache:latest"
	}

	id, _ := e.ImageFindIdByName(e.imageCacheName)
	if id == "" {
		e.enableCache = false
	}

	e.restartPolicy = KRestartPolicyNo

	if e.autoDockerfile == nil {
		e.autoDockerfile = &dockerfileGolang.DockerfileGolang{}
	}

	for _, copyFile := range e.copyFile {
		e.autoDockerfile.AddCopyToFinalImage(copyFile.Src, copyFile.Dst)
	}

	if e.sshDefaultFileName != "" {
		e.autoDockerfile.SetDefaultSshFileName(e.sshDefaultFileName)
	}

	if e.environmentVar == nil {
		e.environmentVar = make([]string, 0)
	}

	onStart := make(chan bool, 1)
	e.onContainerReady = &onStart

	onInspect := make(chan bool, 1)
	e.onContainerInspect = &onInspect

	e.changePointer = *NewImagePullStatusChannel()

	e.dockerSys = DockerSystem{}
	err = e.dockerSys.Init()
	if err != nil {
		util.TraceToLog()
		return
	}

	if e.inspectInterval != 0 {
		e.ticker = time.NewTicker(e.inspectInterval)

		go func(e *ContainerBuilder) {
			var err error
			var logs []byte

			for {
				select {
				case <-e.ticker.C:

					if e.containerID == "" {
						e.containerID, err = e.dockerSys.ContainerFindIdByName(e.containerName)
						if err != nil {
							continue
						}
					}

					e.inspect, _ = e.dockerSys.ContainerInspectParsed(e.containerID)
					logs, _ = e.dockerSys.ContainerLogs(e.containerID)
					e.logs = string(logs)
					*e.onContainerInspect <- true
				}
			}

		}(e)
	}
	return
}
