package test

import (
	"fmt"
	"github.com/go-pansy/pansy"
	"testing"
)

func Test_Signatory(t *testing.T) {
	var (
		signer  = pansy.NewSignatory("abc")
		payload = make(map[string]any)
	)
	payload = map[string]any{
		"machine_no": "16327128",
		"password":   "332237",
		"timestamp":  "1722825479",
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
		"features": []map[string]any{
			{
				"name":  "feature1",
				"value": "enabled",
			},
			{
				"name":  "feature2",
				"value": "disabled",
			},
			{
				"name":  "feature3",
				"value": "enabled",
			},
		},
	}

	sign := signer.GenSignature(payload)

	data, err := signer.ToBase64String(payload)
	fmt.Println("To Base64 String:", data, err)

	isOK := signer.CheckSignature(payload, sign)
	fmt.Println(isOK)
}
