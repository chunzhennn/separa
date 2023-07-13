package plugin

import (
	"separa/pkg"

	"github.com/chainreactors/logs"
	"github.com/chainreactors/parsers/iutils"
)

type RunnerOpts struct {
	VersionLevel int
	Delay        int
	HttpsDelay   int
	SuffixStr    string
	Debug        bool
}

var RunOpt RunnerOpts

func Dispatch(result *pkg.Result) {
	defer func() {
		if err := recover(); err != nil {
			logs.Log.Errorf("scan %s unexcept error, %v", result.GetTarget(), err)
			panic(err)
		}
	}()

	// 初始化扫描
	initScan(result)

	// 如果没有开放便返回
	if !result.Open {
		return
	}

	// 指纹识别
	fingerScan(result)

	// 进行服务猜测
	if !result.IsHttp && result.NoFramework() {
		result.GuessFramework()
	}

	result.Title = iutils.AsciiEncode(result.Title)
}
