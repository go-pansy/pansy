package test

import (
	"fmt"
	"github.com/go-pansy/pansy"
	"testing"
)

func Test_Signatory(t *testing.T) {
	var (
		signer  = pansy.NewSignatory("ds069ed4223ac1660f")
		payload = make(map[string]any)
	)
	payload = map[string]any{
		"machine_no": "16327128", "password": "599240",
		"os": map[string]any{
			"version":   "12.0.0",
			"platform":  "linux",
			"family":    "unix",
			"locale":    "zh-CN",
			"hostname":  "debian",
			"os_type":   "linux",
			"serial_sn": "unknown",
			"model_no":  "unknown",
		},
		"timestamp": "1722846152",
	}

	data, err := signer.ToBase64String(payload)
	fmt.Println("To Base64 String:", data, err)

	isOK := signer.CheckSignature(payload, "AC7822FA3BA8DE95E73486CDC7B60CDF")
	fmt.Println(isOK)
}
