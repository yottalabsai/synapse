package job

import (
	"github.com/robfig/cron/v3"
	"go.uber.org/zap"
	"synapse/common"
	"synapse/log"
)

type SynapseJobManager struct {
	inferencePublicModelJob *InferencePublicModelJob
}

func NewSynapseJobManager(inferencePublicModelJob *InferencePublicModelJob) *SynapseJobManager {
	return &SynapseJobManager{inferencePublicModelJob: inferencePublicModelJob}
}

var IsRunning = true

func (j *SynapseJobManager) StartJobs() {
	log.Log.Info("Start jobs......")
	var c = cron.New()
	// add cron job
	AddCronJob(c, common.JobInferencePublicListCheckSpec, common.InferencePublicListCheckLockPrefix, common.JobResourceStatusCheckSecond, j.inferencePublicModelJob.Run)
	// start cron
	c.Start()
}

// AddCronFunc not support Delay
func AddCronFunc(c *cron.Cron, spec string, cmd func()) {
	_, err := c.AddFunc(spec, func() {
		if IsRunning {
			cmd()
			return
		} else {
			log.Log.Info("The service has gracefully shut down, new scheduled tasks will no longer be executed.")
		}
	})
	if err != nil {
		log.Log.Errorw("Failed to add scheduled tas", zap.Error(err))
	}
}

// AddCronJob support Delay or Skip
func AddCronJob(c *cron.Cron, spec string, key string, timeout int, cmd func()) {
	_, err := c.AddJob(spec, cron.NewChain(cron.SkipIfStillRunning(cron.DefaultLogger)).Then(Wrapper{key: key, timeout: timeout, cmd: cmd}))
	if err != nil {
		log.Log.Errorw("Failed to add scheduled tas", zap.Error(err))
	}
}

type Wrapper struct {
	key     string
	timeout int
	cmd     func()
}

func (t Wrapper) Run() {
	if IsRunning {
		t.cmd()
		return
	} else {
		log.Log.Info("The service has gracefully shut down, new scheduled tasks will no longer be executed.")
	}
}
