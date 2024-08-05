package pansy

import (
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/bytedance/sonic"
	"sort"
	"strings"
	"time"
)

var _ signer = (*Signatory)(nil)

type signer interface {
	GenSignature(args map[string]any) string
	ToBase64String(args map[string]any) (string, error)
	DecryptBase64String(args string) (map[string]any, error)
	CheckSignature(args map[string]any, sign string) bool
}

type Signatory struct {
	appKey string
}

func NewSignatory(appKey string) *Signatory {
	return &Signatory{
		appKey: appKey,
	}
}

// GenSignature implements board.
func (s *Signatory) GenSignature(source map[string]any) string {
	// 直接删掉sign字段
	delete(source, "sign")

	if _, ok := source["timestamp"]; !ok {
		source["timestamp"] = fmt.Sprintf("%d", time.Now().Unix())
	}

	var (
		hasher = md5.New()
		code   = s.JoinMap(source)
	)
	code = fmt.Sprintf("%s&key=%v", code, s.appKey)

	hasher.Write([]byte(code))
	hb := hasher.Sum(nil)
	sign := hex.EncodeToString(hb)

	return strings.ToUpper(sign)
}

func (s *Signatory) JoinMap(source map[string]any) string {
	var (
		keys []string
		code string
	)
	for k, v := range source {
		switch v.(type) {
		case string, int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, float32, float64, bool:
			value := fmt.Sprintf("%v", v)
			if value == "" || value == "NULL" {
				continue
			}
			keys = append(keys, k)
		case []map[string]any, []any:
			// TODO: 递归
		case map[string]any, any:
			source[k] = s.JoinMap(v.(map[string]any))
		default:
			value := fmt.Sprintf("%v", v)
			if value == "" || value == "NULL" {
				continue
			}
			keys = append(keys, k)
		}
	}

	// 字典序 从小到大
	sort.Strings(keys)
	for _, v := range keys {
		code += fmt.Sprintf("%s=%v&", v, source[v])
	}

	return code[0 : len(code)-1]
}

// CheckSignature implements board.
func (s *Signatory) CheckSignature(source map[string]any, sign string) bool {
	return s.GenSignature(source) == sign
}

// DecryptBase64String implements board.
func (s *Signatory) DecryptBase64String(str string) (map[string]any, error) {
	var (
		payload map[string]any
	)

	if len(str) == 0 {
		return nil, errors.New("payload is empty")
	}

	body, err := base64.StdEncoding.DecodeString(str)
	if err != nil {
		return nil, err
	}

	if err = sonic.Unmarshal(body, &payload); err != nil {
		return nil, err
	}

	return payload, nil
}

// ToBase64String implements board.
func (s *Signatory) ToBase64String(source map[string]any) (string, error) {
	// 去掉存在空值的键
	for key, value := range source {
		if value == nil {
			delete(source, key)
		}

		var (
			body = fmt.Sprintf("%v", value)
		)

		// 过滤掉字符串 为 空 或者 为 NULL
		if body == "" || body == "NULL" {
			delete(source, key)
		}
	}

	if _, ok := source["sign"]; !ok {
		source["sign"] = s.GenSignature(source)
	}

	content, err := sonic.Marshal(source)
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(content), nil
}

func DecodeString[T any](args string, result *T) error {
	body, err := base64.StdEncoding.DecodeString(args)
	if err != nil {
		return err
	}

	if err = sonic.Unmarshal(body, result); err != nil {
		return err
	}

	return nil
}

func EncodeToString[T any]() {}
