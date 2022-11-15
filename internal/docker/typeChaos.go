package docker

import "time"

// chaos
//
// English:
//
//	Object chaos manager
//
// Português:
//
//	Objeto gerenciador de caos
type chaos struct {
	foundSuccess               bool
	foundFail                  bool
	filterToStart              []LogFilter
	filterRestart              []LogFilter
	filterSuccess              []LogFilter
	filterFail                 []LogFilter
	filterMonitor              []LogFilter
	filterLog                  []LogFilter
	sceneName                  string
	logPath                    string
	serviceStartedAt           time.Time
	minimumTimeBeforeRestart   time.Duration
	maximumTimeBeforeRestart   time.Duration
	minimumTimeToStartChaos    time.Duration
	maximumTimeToStartChaos    time.Duration
	minimumTimeToPause         time.Duration
	maximumTimeToPause         time.Duration
	minimumTimeToUnpause       time.Duration
	maximumTimeToUnpause       time.Duration
	minimumTimeToRestart       time.Duration
	maximumTimeToRestart       time.Duration
	restartProbability         float64
	restartChangeIpProbability float64
	restartLimit               int
	enableChaos                bool
	event                      chan Event
	monitorStop                chan struct{}
	monitorRunning             bool
	//containerStarted         bool
	containerPaused          bool
	containerStopped         bool
	linear                   bool
	chaosStarted             bool
	chaosCanRestartContainer bool
	//chaosCanRestartEnd       bool
	eventNext time.Time

	disableStopContainer  bool
	disablePauseContainer bool
}
