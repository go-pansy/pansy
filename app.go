package pansy

import (
	"fmt"
	"github.com/thanhpk/randstr"
)

var _ app = (*App)(nil)

type App struct {
}

func (a *App) GenAppId(size int, prefix ...string) string {
	return fmt.Sprintf("%s%v", prefix[0], randstr.Hex(size))
}

func (a *App) GenAppSecret(size int) string {
	return randstr.Hex(size)
}

func (a *App) GenAppKey(size int) string {
	return randstr.Hex(size)
}

func (a *App) GenUsername(prefix ...string) string {
	return fmt.Sprintf("%v%v", prefix[0], randstr.Hex(6))
}

func (a *App) GenCode(size int) string {
	return randstr.Dec(size)
}

type app interface {
	GenAppId(size int, prefix ...string) string
	GenAppSecret(size int) string
	GenAppKey(size int) string
	GenUsername(prefix ...string) string
	GenCode(size int) string
}

func NewApp() *App {
	return &App{}
}
