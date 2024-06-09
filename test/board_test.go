package test

import (
	"github.com/go-pansy/pansy"
	"testing"
)

func Test_Board(t *testing.T) {
	var (
		board   = pansy.NewBoard("abc")
		payload = make(map[string]any)
	)
	payload = map[string]any{
		"a":         "123",
		"b":         "456",
		"c":         "789",
		"sign":      "0C59AD369ED0E8F5AD06610EBBE6F263",
		"timestamp": "1717899356",
	}

	body, err := board.ToBase64String(payload)
	if err != nil {
		t.Log(err)
		return
	}

	value, err := board.DecryptBase64String(body)
	if err != nil {
		t.Log(err)
		return
	}

	isOK := board.CheckSignature(value, payload["sign"].(string))

	t.Logf("解密的数据签名验证结果：%v", isOK)
}
