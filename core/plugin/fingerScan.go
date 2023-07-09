package plugin

import (
	"separa/pkg"
	"separa/pkg/fingers"

	"github.com/chainreactors/parsers"
)

func fingerScan(result *pkg.Result) {
	//如果是http协议,则判断cms,如果是tcp则匹配规则库.暂时不考虑udp
	if result.IsHttp {
		sender := func(sendData []byte) ([]byte, bool) {
			conn := result.GetHttpConn(RunOpt.Delay)
			url := result.GetURL() + string(sendData)
			resp, err := conn.Get(url)
			if err == nil {
				return parsers.ReadRaw(resp), true
			} else {
				return nil, false
			}
		}

		getFramework(result, pkg.HttpFingers, sender)
	} else {
		if pkg.Proxy != nil {
			// 如果存在http代理，跳过tcp指纹识别
			return
		}
		sender := func(sendData []byte) ([]byte, bool) {
			conn, err := pkg.NewSocket("tcp", result.GetTarget(), RunOpt.Delay/2+1)
			if err != nil {
				return nil, false
			}
			defer conn.Close()
			data, err := conn.QuickRequest(sendData, 1024)
			if err != nil {
				return nil, false
			}

			return data, true
		}

		getFramework(result, pkg.TcpFingers, sender)
	}
	return
}

func getFramework(result *pkg.Result, fingermap fingers.FingerMapper, sender func(sendData []byte) ([]byte, bool)) {
	// 优先匹配默认端口,第一次循环只匹配默认端口
	var matcher func(result *pkg.Result, finger *fingers.Finger, sender func(sendData []byte) ([]byte, bool)) (*parsers.Framework, *parsers.Vuln)
	if result.IsHttp {
		matcher = httpFingerMatch
	} else {
		matcher = tcpFingerMatch
	}
	var alreadyFrameworks = make(map[string]bool)
	for _, finger := range fingermap[result.Port] {
		frame, vuln := matcher(result, finger, sender)
		alreadyFrameworks[finger.Name] = true
		if frame != nil {
			if vuln != nil {
				result.AddVuln(vuln)
			}
			result.AddFramework(frame)
			if result.Protocol == "tcp" {
				// 如果是tcp协议,并且识别到一个指纹,则退出.
				// 如果是http协议,可能存在多个指纹,则进行扫描
				return
			}
		}
	}

	for _, fs := range fingermap {
		for _, finger := range fs {
			if _, ok := alreadyFrameworks[finger.Name]; ok {
				continue
			} else {
				alreadyFrameworks[finger.Name] = true
			}

			frame, vuln := matcher(result, finger, sender)
			if frame != nil {
				result.AddFramework(frame)
				if vuln != nil {
					result.AddVuln(vuln)
				}
				if result.Protocol == "tcp" {
					return
				}
			}
		}
	}
	return
}

func httpFingerMatch(result *pkg.Result, finger *fingers.Finger, sender func(sendData []byte) ([]byte, bool)) (*parsers.Framework, *parsers.Vuln) {
	frame, vuln, ok := fingers.FingerMatcher(finger, result.ContentMap(), RunOpt.VersionLevel, sender)
	if ok {
		if len(frame.Data) != 0 {
			result.Title = parsers.MatchTitle(frame.Data)
		}
		return frame, vuln
	}
	return nil, nil
}

func tcpFingerMatch(result *pkg.Result, finger *fingers.Finger, sender func(sendData []byte) ([]byte, bool)) (*parsers.Framework, *parsers.Vuln) {
	frame, vuln, ok := fingers.FingerMatcher(finger, result.ContentMap(), RunOpt.VersionLevel, sender)
	if ok {
		return frame, vuln
	}
	return nil, nil
}
