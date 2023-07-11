package plugin

import (
	"bytes"

	"github.com/M09ic/go-ntlmssp"

	"separa/pkg"

	"github.com/chainreactors/parsers/iutils"
)

var data = pkg.Decode("YmXgZhZgYGCoYNBgYGZgYNghsAPEZWAEY0aGBSAGAwPDAQjlBiJYYju6XsucFJz/goNBW8AjgYmBgYGLCaLAL8THNzg4AKyfvYljEQMaYGPcKMvAwMAPAAAA//8=")

func wmiScan(result *pkg.Result) {
	result.Port = "135"
	target := result.GetTarget()
	conn, err := pkg.NewSocket("tcp", target, RunOpt.Delay)
	if err != nil {
		return
	}
	defer conn.Close()

	result.Open = true
	ret, err := conn.Request(data, 4096)
	if err != nil {
		return
	}

	if bytes.Index(ret, []byte("NTLMSSP")) != -1 {
		result.Protocol = "wmi"
		result.Status = "WMI"
		tinfo := iutils.ToStringMap(ntlmssp.NTLMInfo(ret))
		result.AddNTLMInfo(tinfo, "wmi")
	}
}
