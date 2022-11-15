package docker

import (
	"github.com/docker/docker/api/types"
	"log"
	"runtime"
	"sync"
)

// RemoveAllByNameContains remove trash after test.
// This function removes container, image and network by name, and unlinked volumes and
// imagens
func (el DockerSystem) RemoveAllByNameContains(name string) (err error) {
	var nameAndId []NameAndId
	var container types.ContainerJSON
	var wg sync.WaitGroup

	runtime.GOMAXPROCS(runtime.NumCPU() * 2)

	// quando tem algo em torno de 255 containers, este código falha, por isto, o laço
	for {
		nameAndId, err = el.ContainerFindIdByNameContains(name)
		if err != nil && err.Error() != "container name not found" {
			return err
		}

		if len(nameAndId) == 0 {
			break
		}

		for _, data := range nameAndId {
			wg.Add(1)

			go func(id string) {
				defer wg.Done()

				container, err = el.ContainerInspect(id)
				if err != nil {
					return
				}

				if container.State != nil && container.State.Running == true {
					err = el.ContainerStop(id)
					if err != nil {
						return
					}
				}
			}(data.ID)
		}
		wg.Wait()

		for _, data := range nameAndId {
			wg.Add(1)

			go func(data NameAndId) {
				defer wg.Done()

				container, err = el.ContainerInspect(data.ID)
				if err != nil {
					return
				}

				if container.State != nil && container.State.Running == true {
					err = el.ContainerStopAndRemove(data.ID, true, false, true)
					if err != nil {
						return
					}
				} else if container.State != nil && container.State.Running == false {
					err = el.ContainerRemove(data.ID, true, false, true)
					if err != nil {
						return
					}
				}

				log.Printf("remove: %v", data.Name)
			}(data)
		}
		wg.Wait()
	}

	nameAndId, err = el.ImageFindIdByNameContains(name)
	if err != nil && err.Error() != "image name not found" {
		return err
	}
	for _, data := range nameAndId {
		err = el.ImageRemove(data.ID, true, false)
		if err != nil {
			return
		}
	}

	nameAndId, err = el.NetworkFindIdByNameContains(name)
	if err != nil && err.Error() != "network name not found" {
		return err
	}
	for _, data := range nameAndId {
		err = el.NetworkRemove(data.ID)
		if err != nil {
			return
		}
	}

	err = el.VolumesUnreferencedRemove()
	if err != nil {
		return
	}

	err = el.ImageGarbageCollector()
	if err != nil {
		return
	}

	return
}
