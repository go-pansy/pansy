package test

import (
	"github.com/go-pansy/pansy"
	"testing"
)

func Test_Faker(t *testing.T) {
	var (
		faker      = pansy.NewFaker()
		id    uint = 1
	)
	fakerId := faker.GenId(id, true)

	t.Log(id)

	if id == faker.RecoverId(fakerId) {
		t.Log("OK")
		return
	}

	t.Log("Fail")
}
