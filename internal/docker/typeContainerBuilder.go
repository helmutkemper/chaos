package docker

import (
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	dockerfileGolang "github.com/helmutkemper/iotmaker.docker.builder.golang.dockerfile"
	isolatedNetwork "github.com/helmutkemper/iotmaker.docker.builder.network.interface"
	"time"
)

// ContainerBuilder
//
// English:
//
//	Docker manager
//
// Português:
//
//	Gerenciador de containers e imagens docker
type ContainerBuilder struct {
	buildFolder bool
	buildServer bool

	metadata                map[string]interface{}
	problem                 string
	csvValueSeparator       string
	csvRowSeparator         string
	csvConstHeader          bool
	logCpus                 int
	rowsToPrint             int64
	chaos                   chaos
	enableCache             bool
	network                 isolatedNetwork.ContainerBuilderNetworkInterface
	dockerSys               DockerSystem
	changePointer           chan ContainerPullStatusSendToChannel
	onContainerReady        *chan bool
	onContainerInspect      *chan bool
	imageInspected          bool
	imageInstallExtras      bool
	imageCacheName          string
	imageName               string
	imageID                 string
	imageRepoTags           []string
	imageRepoDigests        []string
	imageParent             string
	imageComment            string
	imageCreated            time.Time
	imageContainer          string
	imageAuthor             string
	imageArchitecture       string
	imageVariant            string
	imageOs                 string
	imageOsVersion          string
	imageSize               int64
	imageVirtualSize        int64
	containerName           string
	buildPath               string
	environmentVar          []string
	changePorts             []dockerfileGolang.ChangePort
	openPorts               []string
	exposePortsOnDockerfile []string
	openAllPorts            bool
	waitString              string
	waitStringTimeout       time.Duration
	containerID             string
	ticker                  *time.Ticker
	inspect                 ContainerInspect
	logs                    string
	logsLastSize            int
	inspectInterval         time.Duration
	gitData                 gitData
	volumes                 []mount.Mount
	IPV4Address             string
	autoDockerfile          DockerfileAuto
	containerConfig         container.Config
	restartPolicy           RestartPolicy

	replaceDockerfile          string
	addFileToServerBeforeBuild []CopyFile

	makeDefaultDockerfile bool
	printBuildOutput      bool
	init                  bool
	startedAfterBuild     bool

	sshDefaultFileName string

	contentIdEcdsaFile             string
	contentIdEcdsaFileWithScape    string
	contentIdRsaFile               string
	contentIdRsaFileWithScape      string
	contentKnownHostsFile          string
	contentKnownHostsFileWithScape string
	contentGitConfigFile           string
	contentGitConfigFileWithScape  string

	gitPathPrivateRepository string

	buildOptions        types.ImageBuildOptions
	imageExpirationTime time.Duration

	copyFile []CopyFile
}
