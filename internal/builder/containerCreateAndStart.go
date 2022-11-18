package builder

import (
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/go-connections/nat"
)

// ContainerCreateAndStart (English): Creates a container, automatically exposes the
// ports listed in the image and start then
//
//	imageName: image name for download and pull
//	containerName: unique container name
//	restartPolicy:
//	   KRestartPolicyNo - Do not automatically restart the container.
//	   KRestartPolicyOnFailure - Restart the container if it exits due to an error,
//	   which manifests as a non-zero exit code.
//	   KRestartPolicyAlways - Always restart the container if it stops. If it is
//	   manually stopped, it is restarted only when Docker daemon restarts or the
//	   container itself is manually restarted. (See the second bullet listed in restart
//	   policy details)
//	   KRestartPolicyUnlessStopped - Similar to always, except that when the container
//	   is stopped (manually or otherwise), it is not restarted even after Docker daemon
//	   restarts.
//	portExposedList: nat.PortMap exposed port list from container
//	   key: string (container port umber/protocol tcp|udp). Example: "3000/tcp"
//	   value: []nat.PortBinding
//	     key: one key per server host port
//	     value: struct{ HostPort: numeric string host port }. Example: {HostPort: "3000"}
//	mountVolumes: array of mount.Mount{}
//	   Type:
//	      KVolumeMountTypeBindString: Bind is the type for mounting host dir (real
//	      folder inside computer where this code work)
//	      KVolumeMountTypeVolumeString: Volume is the type for remote storage volumes
//	      KVolumeMountTypeTmpfsString: Tmpfs is the type for mounting tmpfs
//	      KVolumeMountTypeNpipeString: NPipe is the type for mounting Windows named pipes
//	   Source: path inside the host machine
//	   Target: path inside the container after container start
//	      Note: For a complete list of volumes exposed by image, use
//	      ImageListExposedVolumes(id) and ImageListExposedVolumesByName(name)
//	containerNetwork: container network configuration
//	  Note: please, use NetworkCreate() for correct configuration of network
//
// ContainerCreateAndStart (Português): Cria um container, automaticamente expões as
// portas listas na imagem e depois o inicia
//
//	imageName: nome da image para download e pull
//	containerName: nome do container (deve ser único)
//	restartPolicy:
//	   KRestartPolicyNo - Não reinicia o container automaticamente.
//	   KRestartPolicyOnFailure - Reinicia o container ser ele terminar por error,
//	   com o manifesto contendo um valore de erro diferente de zero.
//	   KRestartPolicyAlways - Sempre reinicia o container se ele parar. Caso o container
//	   seja parado, ele só reinicia se o Docker daemon reiniciar ou com reinício manual
//	   (Veja 'second bullet' nas políticas de reinício do container)
//	   KRestartPolicyUnlessStopped - Similar ao always, exceto quando o container é
//	   parado (manualmente ou não) e não é reiniciado, mesmo quando o Docker daemon
//	   reinicia
//	portExposedList: nat.PortMap lista de portas expostas do container
//	   key: string (número da porta do container/protocolo tcp|udp). Exemplo: "3000/tcp"
//	   value: []nat.PortBinding
//	     key: uma chave por porta no servidor hospedeiro
//	     value: struct{ HostPort: string numérica da porta no servidor hospedeiro }.
//	            Example: {HostPort: "3000"}
//	mountVolumes: array de mount.Mount{}
//	   Type:
//	      KVolumeMountTypeBindString: Vincula uma pasta do computador, host, com uma
//	      pasta do container
//	      KVolumeMountTypeVolumeString: Volume é usado para arquivos de armazenamento
//	      remoto
//	      KVolumeMountTypeTmpfsString: Tmpfs arquiva os dados em memória RAM de forma
//	      volátil
//	      KVolumeMountTypeNpipeString: NPipe é um tipo de volume do Windows chamado de
//	      pipes
//	   Source: caminho dentro do computador, host
//	   Target: caminho dentro do container quando ele inicia
//	      Nota: Para uma lista completa de volumes expostos na imagem, use
//	      ImageListExposedVolumes(id) e ImageListExposedVolumesByName(name)
//	containerNetwork: configuração de rede do container
//	  Nota: Por favor, use NetworkCreate() para a forma correta da configuração de rede
func (el *DockerSystem) ContainerCreateAndStart(
	imageName string,
	containerName string,
	restart RestartPolicy,
	portExposedList nat.PortMap,
	mountVolumes []mount.Mount,
	containerNetwork *network.NetworkingConfig,
) (
	containerID string,
	err error,
) {

	imageName = el.AdjustImageName(imageName)

	containerID, err = el.ContainerCreate(
		imageName,
		containerName,
		restart,
		portExposedList,
		mountVolumes,
		containerNetwork,
	)
	if err != nil {
		return
	}

	err = el.ContainerStart(containerID)
	return
}
