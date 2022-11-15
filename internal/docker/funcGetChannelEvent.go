package docker

// GetChannelEvent (english):
//
// GetChannelEvent (português): Canal disparado durante o processo de image build ou container build e retorna
// informações como andamento do download da imagem, processo de extração da mesma entre outras informações
//
//	Waiting: Esperando o processo ser iniciado pelo docker
//	Downloading: Estado do download da imagem, caso a mesma não exista na máquina host
//	  Count: Quantidade de blocos a serem baixados
//	  Current: Total de bytes baixados até o momento
//	  Total: Total de bytes a serem baixados
//	  Percent: Percentual atual do processo com uma casa decimal de precisão
//	DownloadComplete: todo: fazer
//	Extracting: Estado da extração da imagem baixada
//	  Count: Quantidade de blocos a serem extraídos
//	  Current: Total de bytes extraídos até o momento
//	  Total: Total de bytes a serem extraídos
//	  Percent: Percentual atual do processo com uma casa decimal de precisão
//	PullComplete: todo: fazer
//	ImageName: nome da imagem baixada
//	ImageID: ID da imagem baixada. (Cuidado: este valor só é definido ao final do processo)
//	ContainerID: ID do container criado. (Cuidado: este valor só é definido ao final do processo)
//	Closed: todo: fazer
//	Stream: saída padrão do container durante o processo de build
//	SuccessfullyBuildContainer: sucesso ao fim do processo de build do container
//	SuccessfullyBuildImage: sucesso ao fim do processo de build da imagem
//	IdAuxiliaryImages: usado pelo coletor de lixo para apagar as imagens axiliares ao fim do processo de build
func (e *ContainerBuilder) GetChannelEvent() (channel *chan ContainerPullStatusSendToChannel) {
	return &e.changePointer
}
