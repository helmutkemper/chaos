package docker

import (
	"fmt"
)

// NewPullStatusMonitor (English): Small example how to use channel do monitoring
// image/container build
//
// NewPullStatusMonitor (Português): Pequeno exemplo de como usar o canal para ver
// imagem/container sendo criados
func NewPullStatusMonitor() (pullStatusChannel *chan ContainerPullStatusSendToChannel) {
	pullStatusChannel = NewImagePullStatusChannel()

	go func(c chan ContainerPullStatusSendToChannel) {

		for {
			select {
			case status := <-c:
				// English: remove this comment to see all build status
				// Português: remova este comentário para vê _todo o status da criação da imagem
				//fmt.Printf("image pull status: %+v\n", status)

				if status.Closed == true {
					fmt.Println("image pull complete!")

					// English: Eliminate this goroutine after process end
					// Português: Elimina a goroutine após o fim do processo
					return
				}
			}
		}

	}(*pullStatusChannel)

	return
}
