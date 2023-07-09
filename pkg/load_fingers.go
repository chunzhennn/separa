package pkg

import (
	"separa/pkg/fingers"

	"github.com/chainreactors/parsers/iutils"
)

var (
	Md5Fingers     map[string]string
	Mmh3Fingers    map[string]string
	AllHttpFingers fingers.Fingers
	TcpFingers     fingers.FingerMapper
	HttpFingers    fingers.FingerMapper
)

// 加载指纹到全局变量
func LoadFinger(t string) fingers.Fingers {
	fs, err := fingers.LoadFingers(LoadConfig(t))
	if err != nil {
		iutils.Fatal(err.Error())
	}
	return fs
}

func LoadHashFinger(fs fingers.Fingers) (map[string]string, map[string]string) {
	md5hash := make(map[string]string)
	mmh3hash := make(map[string]string)
	for _, f := range fs {
		for _, rule := range f.Rules {
			if rule.Favicon != nil {
				for _, mmh3 := range rule.Favicon.Mmh3 {
					mmh3hash[mmh3] = f.Name
				}
				for _, md5 := range rule.Favicon.Md5 {
					md5hash[md5] = f.Name
				}
			}
		}
	}
	return mmh3hash, md5hash
}
