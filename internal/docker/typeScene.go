package docker

import (
	"sync"
)

// scene
//
// English: Determines the maximum number of stopped and paused containers per scene.
// The scene grants some control over the chaos in container testing, preventing all containers in a scene from pausing
// or stopping at the same time.
//
// Português: Determina os números máximos de containers parados e pausados por cena.
// A cena garante algum controle sobre o caos no teste dos container, impedindo que todos os containers de uma cena
// sejam pausando ou parados ao mesmo tempo.
type scene struct {
	StopedContainers                   int
	PausedContainers                   int
	MaxStopedContainers                int
	MaxPausedContainers                int
	MaxTotalPausedAndStoppedContainers int
}

// Theater
//
// English: Theater is the collection of scene
//
// Português: Teatro é a coleção de cenas
type Theater struct {
	s sync.Mutex
	m map[string]scene
}

// Init
//
// English: Initialization must always be the first function called.
//
// Português: A inicialização sempre deve ser a primeira função chamada.
func (e *Theater) Init() {
	e.s.Lock()
	defer e.s.Unlock()

	e.m = make(map[string]scene)
}

// ConfigScene
//
// English: Create and configure a new scene.
//
//	Input:
//	  sceneName: unique name of the scene
//	  maxStopedContainers: maximum number of stopped containers
//	  maxPausedContainers: maximum number of paused containers
//
// Português: Cria e configura uma cena nova.
//
//	Entrada:
//	  sceneName: nome único da cena
//	  maxStopedContainers: quantidade máxima de containers parados
//	  maxPausedContainers: quantidade máxima de containers pausados
func (e *Theater) ConfigScene(sceneName string, maxStopedContainers, maxPausedContainers, maxTotalPausedAndStoppedContainers int) {
	e.s.Lock()
	defer e.s.Unlock()

	e.m[sceneName] = scene{
		StopedContainers:                   0,
		PausedContainers:                   0,
		MaxStopedContainers:                maxStopedContainers,
		MaxPausedContainers:                maxPausedContainers,
		MaxTotalPausedAndStoppedContainers: maxTotalPausedAndStoppedContainers,
	}
}

// SetContainerUnPaused
//
// English: Decreases the paused containers counter
//
//	Input:
//	  sceneName: unique name of the scene
//
// Português: Decrementa o contador de containers pausados
//
//	Entrada:
//	  sceneName: nome único da cena
func (e *Theater) SetContainerUnPaused(sceneName string) {
	if sceneName == "" {
		return
	}

	e.s.Lock()
	defer e.s.Unlock()

	sc := e.m[sceneName]
	sc.PausedContainers = sc.PausedContainers - 1
	e.m[sceneName] = sc
}

// SetContainerPaused
//
// English: Increments the paused containers counter
//
//	Input:
//	  sceneName: unique name of the scene
//	Output:
//	  doNotPauseContainer: the maximum number of containers has been reached
//
// Português: Incrementa o contador de containers pausados
//
//	Entrada:
//	  sceneName: nome único da cena
//	Saída:
//	  doNotPauseContainer: a quantidade máxima de containers foi atingida
func (e *Theater) SetContainerPaused(sceneName string) (doNotPauseContainer bool) {
	if sceneName == "" {
		return
	}

	e.s.Lock()
	defer e.s.Unlock()

	sc := e.m[sceneName]

	if sc.MaxPausedContainers <= sc.PausedContainers ||
		sc.MaxTotalPausedAndStoppedContainers <= sc.PausedContainers+sc.StopedContainers {
		return true
	}

	sc.PausedContainers = sc.PausedContainers + 1
	e.m[sceneName] = sc

	return false
}

// SetContainerStopped
//
// English: Increments the stopped containers counter
//
//	Input:
//	  sceneName: unique name of the scene
//	Output:
//	  IsOnTheEdge: the maximum number of containers has been reached
//
// Português: Incrementa o contador de containers parados
//
//	Entrada:
//	  sceneName: nome único da cena
//	Saída:
//	  IsOnTheEdge: a quantidade máxima de containers foi atingida
func (e *Theater) SetContainerStopped(sceneName string) (IsOnTheEdge bool) {
	if sceneName == "" {
		return
	}

	e.s.Lock()
	defer e.s.Unlock()

	sc := e.m[sceneName]

	if sc.MaxStopedContainers <= sc.StopedContainers ||
		sc.MaxTotalPausedAndStoppedContainers <= sc.PausedContainers+sc.StopedContainers {
		return true
	}

	sc.StopedContainers = sc.StopedContainers + 1
	e.m[sceneName] = sc

	return false
}

// SetContainerUnStopped
//
// English: Decreases the stopped containers counter
//
//	Input:
//	  sceneName: unique name of the scene
//
// Português: Decrementa o contador de containers parados
//
//	Entrada:
//	  sceneName: nome único da cena
func (e *Theater) SetContainerUnStopped(sceneName string) {
	if sceneName == "" {
		return
	}

	e.s.Lock()
	defer e.s.Unlock()

	sc := e.m[sceneName]
	sc.StopedContainers = sc.StopedContainers - 1
	e.m[sceneName] = sc
}

var theater = Theater{}

// init
//
// English: Launch the test theater
//
// Português: Inicializa o teatro de teste
func init() {
	theater.Init()
}
