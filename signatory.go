package pansy

import (
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/bytedance/sonic"
	"net/url"
	"reflect"
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

	values := url.Values{}

	var addValues func(key string, value any)
	addValues = func(key string, value any) {
		v := reflect.ValueOf(value)
		switch v.Kind() {
		case reflect.String, reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
			reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
			reflect.Float32, reflect.Float64, reflect.Bool:

			strValue := fmt.Sprintf("%v", value)
			if strValue == "" || strValue == "NULL" {
				return
			}
			values.Add(key, strValue)

		case reflect.Map:
			for _, k := range v.MapKeys() {
				addValues(fmt.Sprintf("%s[%s]", key, k.String()), v.MapIndex(k).Interface())
			}

		case reflect.Slice, reflect.Array:
			for i := 0; i < v.Len(); i++ {
				addValues(fmt.Sprintf("%s[%d]", key, i), v.Index(i).Interface())
			}

		default:
			strValue := fmt.Sprintf("%v", value)
			if strValue == "" || strValue == "NULL" {
				return
			}
			values.Add(key, strValue)
		}
	}

	for k, v := range source {
		addValues(k, v)
	}

	// 对键进行排序
	sortedKeys := make([]string, 0, len(values))
	for key := range values {
		sortedKeys = append(sortedKeys, key)
	}
	sort.Strings(sortedKeys)

	// 拼接排序后的键值对
	result := url.Values{}
	for _, key := range sortedKeys {
		for _, value := range values[key] {
			result.Add(key, value)
		}
	}

	// appKey
	var (
		payload = result.Encode()
	)
	payload = fmt.Sprintf("%s&key=%s", payload, s.appKey)

	hash := md5.Sum([]byte(payload))
	return strings.ToUpper(hex.EncodeToString(hash[:]))
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

func Decode[T any](args string) (T, error) {
	var result T
	return result, DecodeString(args, &result)
}
