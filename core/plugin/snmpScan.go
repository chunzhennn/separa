package plugin

import (
	"bytes"

	"separa/pkg"
	"separa/pkg/fingers"

	"github.com/chainreactors/parsers"
)

var snmpPublicData = parsers.UnHexlify("302902010104067075626c6963a01c02049acb0442020100020100300e300c06082b060102010101000500")

func snmpScan(result *pkg.Result) {
	result.Port = "161"
	conn, err := pkg.NewSocket("udp", result.GetTarget(), 2)
	if err != nil {
		result.Error = err.Error()
		return
	}
	data, err := conn.Request(snmpPublicData, 4096)
	if err != nil {
		result.Error = err.Error()
		return
	}
	if i := bytes.Index(data, []byte{0x0, 0x4}); i != -1 && len(data) > i+3 {
		result.Title = string(data[i+3:])
	}

	result.Open = true
	result.Protocol = "snmp"
	result.Status = "snmp"
	result.AddVuln(&parsers.Vuln{Name: "snmp_public_auth", Payload: map[string]interface{}{"auth": "public"}, SeverityLevel: fingers.INFO})
}
