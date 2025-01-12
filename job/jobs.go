package job

import (
	"github.com/robfig/cron/v3"
	"go.uber.org/zap"
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
	// 添加定时任务
	// AddCronJob(c, common.JobInferencePublicListCheckSpec, common.InferencePublicListCheckLockPrefix, common.JobResourceStatusCheckSecond, j.inferencePublicModelJob.Run)
	// 启动定时任务
	c.Start()
}

// AddCronFunc 不支持Delay
func AddCronFunc(c *cron.Cron, spec string, cmd func()) {
	_, err := c.AddFunc(spec, func() {
		if IsRunning {
			cmd()
			return
		} else {
			log.Log.Info("服务已经进行优雅关机, 新的定时任务不再执行")
		}
	})
	if err != nil {
		log.Log.Error("添加定时任务失败", zap.Error(err))
	}
}

// AddCronJob 支持Delay或Skip
func AddCronJob(c *cron.Cron, spec string, key string, timeout int, cmd func()) {
	_, err := c.AddJob(spec, cron.NewChain(cron.SkipIfStillRunning(cron.DefaultLogger)).Then(Wrapper{key: key, timeout: timeout, cmd: cmd}))
	if err != nil {
		log.Log.Error("添加定时任务失败", zap.Error(err))
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
		log.Log.Info("服务已经进行优雅关机, 新的定时任务不再执行")
	}
}
