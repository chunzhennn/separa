package api

import (
	"separa/common"
	"separa/common/flag"
	"separa/core/plugin"
	"separa/core/report"
	"separa/core/run"

	"github.com/gin-gonic/gin"
)

// Node status
// - 0: 表示空闲
// - 1: 表示忙碌
// - 2: 表示完成
// 0 -> 1 -> 2 -> 0
var Status = 0

var FREE = 0
var BUSY = 1
var DONE = 2

func TaskRun(c *gin.Context) {
	var task Task
	if err := c.BindJSON(&task); err != nil {
		c.JSON(400, gin.H{
			"message": err.Error(),
		})
		return
	}

	switch Status {
	// 当前空闲，允许执行任务
	case FREE:
		{
			go start(task)
			Status = BUSY
			c.JSON(200, gin.H{
				"message": "ok",
			})
			return
		}
	// 当前忙碌，不允许执行任务
	case BUSY:
		{
			c.JSON(200, gin.H{
				"message": "busy",
			})
			return
		}
	// 当前任务已完成，需要等待收集
	case DONE:
		{
			c.JSON(200, gin.H{
				"message": "done",
			})
			return
		}
	}
}

func Collect(c *gin.Context) {
	if Status != DONE {
		c.JSON(200, gin.H{
			"message": "shoule be done but " + string(Status),
		})
		return
	}
	out, _ := c.GetQuery("out")
	err := report.Load(out)
	if err != nil {
		c.JSON(400, gin.H{
			"message": err.Error(),
		})
		return
	}
	Status = FREE
	c.JSON(200, report.ResultKV.KV)
}

func ConfigInit(task Task) {
	flag.Targets = task.Targets
	flag.Command.Scan.OutputFile = task.Output
	flag.Command.Scan.Port = task.Port
	flag.Command.Scan.Delay = task.Delay

	common.Setting = common.New()
	common.Setting.Target = flag.Targets
	common.Setting.Output = flag.Command.Scan.OutputFile
	common.Setting.LoadPort(flag.Command.Scan.Port)

	plugin.RunOpt.Delay = flag.Command.Scan.Delay
	plugin.RunOpt.HttpsDelay = flag.Command.Scan.Delay / 2
}

func start(task Task) {
	ConfigInit(task)
	run.Start()
	Status = DONE
}
