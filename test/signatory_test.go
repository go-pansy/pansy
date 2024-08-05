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
		"a":         "123",
		"b":         "456",
		"c":         "789",
		"sign":      "0C59AD369ED0E8F5AD06610EBBE6F263",
		"timestamp": "1717899356",
		"goods": []map[string]any{{
			"x": "1",
			"y": "2",
			"z": "3",
		}},
	}

	body, err := signer.ToBase64String(payload)
	if err != nil {
		t.Log(err)
		return
	}

	fmt.Println(body)

	value, err := signer.DecryptBase64String(body)
	if err != nil {
		t.Log(err)
		return
	}

	isOK := signer.CheckSignature(value, payload["sign"].(string))

	t.Logf("解密的数据签名验证结果：%v", isOK)
}
