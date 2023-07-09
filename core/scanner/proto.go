package scanner

import (
	"fmt"
	"net"
	"net/url"
	"separa/common/log"
	"separa/core/plugin"
	"separa/core/report"
	"separa/pkg"
	"strconv"
	"strings"

	"github.com/lcvvvv/appfinger"
)

type ProtoScanner struct {
	*client
}

type ipPort struct {
	addr net.IP
	port int
}

func AppFingerParse(result *pkg.Result) *appfinger.FingerPrint {
	var banner *appfinger.Banner
	var finger *appfinger.FingerPrint
	var err error
	URLRaw := fmt.Sprintf("%s://%s:%s", result.Protocol, result.Ip, result.Port)
	URL, _ := url.Parse(URLRaw)
	if result.Protocol == "http" || result.Protocol == "https" {
		banner, err = appfinger.GetBannerWithURL(URL, nil, result.HttpConn)
		if err != nil {
			log.Log.Println(err)
			return nil
		}
		finger = appfinger.Search(URL, banner)
	}
	return finger
}

func NewProtoScanner(config *Config, threads int) (ps *ProtoScanner) {
	ps = &ProtoScanner{
		client: newConfig(config, threads),
	}
	ps.pool.Function = func(in interface{}) {
		value := in.(ipPort)
		// 首先通过 Dispatch 进行指纹识别
		result := pkg.NewResult(value.addr.String(), strconv.Itoa(value.port))
		plugin.Dispatch(result)
		// 如果扫到东西了
		if result.Open {
			log.Log.Printf(result.FullOutput())
			port, _ := strconv.Atoi(result.Port)

			app := make(map[string]string, 0)
			// 遍历所有的指纹信息
			for _, v := range result.Frameworks {
				isHonypot := false
				// 遍历指纹信息的 tag
				for _, tag := range v.Tags {
					// 如果该指纹信息是蜜罐，那么就添加进 honypot
					if tag == "honeypot" {
						report.AppendHonypot(result.Ip, result.Port+"/"+v.Name)
						isHonypot = true
						continue
					}
					// 如果该指纹信息是设备，那么就更新设备信息
					if tag == "device" {
						report.UpdateDeviceinfo(result.Ip, v.Name)
					}
				}
				// 如果是蜜罐，那么就跳过后续加入到 serviceApp 中
				if isHonypot {
					continue
				}
				if v.Version != "" {
					app[v.Name] = v.Version
				} else {
					app[v.Name] = "N"
				}
			}
			if result.Midware != "" {
				name, version := report.AttachVersion(result.Midware)
				app[name] = version
			}
			if result.Os != "" {
				name, version := report.AttachVersion(result.Os)
				app[name] = version
			}
			// 对于 http, https 协议使用 appFinger 补充指纹信息
			finger := AppFingerParse(result)
			if finger != nil {
				for _, name := range finger.ProductName {
					// 去除finger.ProductName里的 '\t' 并小写
					name = strings.ToLower(strings.ReplaceAll(name, "\t", ""))
					// 可能有 version 信息
					index := strings.LastIndex(name, "/")
					// 如果 version 为空，则默认为 N
					version := "N"
					prod := name
					if index != -1 {
						prod = name[:index]
						version = name[index+1:]
					}

					// 如果 productMap 中已经存在 prod，则比较 version 是否为 N，为 N 且新的不为 N 则替换
					_, ok := app[prod]
					if ok && app[prod] != "N" {
						continue
					}
					app[prod] = version
				}
			}

			// 合并结果
			appVec := make([]string, 0)
			for k, v := range app {
				appVec = append(appVec, k+"/"+v)
			}
			service := report.NewServiceUnit(port, result.Protocol, appVec)
			report.AppendService(result.Ip, service)
		}
	}
	return
}

func (c *ProtoScanner) Push(ip net.IP, port int) {
	c.pool.Push(ipPort{ip, port})
}
