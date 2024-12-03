package job

import (
	"github.com/robfig/cron/v3"
	"go.uber.org/zap"
	"synapse/log"
)

var IsRunning = true

func StartJobs() {
	var c = cron.New()
	// 添加定时任务
	// AddCronJob(c, common.JobResourceStatusCheckSpec, common.ResourceStatusCheckLockPrefix, common.JobResourceStatusCheckSecond, ResourceStatusCheckJob)
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
	//if IsRunning {
	//	lock := utils.Lock{Key: t.key, Conn: config.GetRDB(), Timeout: t.timeout}
	//	defer func(Conn redis.Conn) {
	//		err := Conn.Close()
	//		if err != nil {
	//			config.Logger.Error("lock/unlock: close redis connection error:", zap.Error(err))
	//		}
	//	}(lock.Conn)
	//	ok, err := lock.TryLock()
	//	if err != nil {
	//		config.Logger.Error("try lock error:", zap.String("key", t.key), zap.Int("timeout", t.timeout), zap.Error(err))
	//		return
	//	}
	//	if !ok {
	//		config.Logger.Warn("get lock failed:", zap.String("key", t.key), zap.Int("timeout", t.timeout))
	//		return
	//	}
	//	t.cmd()
	//	err = lock.Unlock()
	//	if err != nil {
	//		config.Logger.Error("unlock error:", zap.String("key", t.key), zap.Int("timeout", t.timeout), zap.Error(err))
	//	}
	//	return
	//} else {
	//	config.Logger.Info("服务已经进行优雅关机, 新的定时任务不再执行")
	//}
}
