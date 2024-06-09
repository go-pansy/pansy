package pansy

import (
	"strconv"
	"strings"
	"sync"
)

var _ faker = (*Faker)(nil)

type faker interface {
	GenId(id uint, args ...bool) string
	RecoverId(id string) uint
}

type Faker struct {
	lock sync.Mutex
}

func (f *Faker) GenId(id uint, args ...bool) string {
	f.lock.Lock()
	defer f.lock.Unlock()

	withPrefix := true
	if len(args) > 0 {
		withPrefix = args[0]
	}

	s := strconv.FormatInt(int64(id), 31)
	reversed := reverse(s)
	if len(reversed) < 6 {
		reversed += strings.Repeat("0", 6-len(reversed))
	}

	if withPrefix {
		return "0x" + reversed
	}

	return reversed
}

func (f *Faker) RecoverId(id string) uint {
	f.lock.Lock()
	defer f.lock.Unlock()

	if id == "" {
		return 0
	}

	reversed := reverse(strings.TrimPrefix(id, "0x"))
	v, err := strconv.ParseInt(reversed, 31, 32)
	if err != nil {
		return 0
	}

	return uint(v)
}

func reverse(s string) string {
	last := len(s) - 1
	runes := []rune(s)
	for i := 0; i < len(s)/2; i++ {
		runes[i], runes[last-i] = runes[last-i], runes[i]
	}

	return string(runes)
}

func NewFaker() *Faker {
	return &Faker{}
}
