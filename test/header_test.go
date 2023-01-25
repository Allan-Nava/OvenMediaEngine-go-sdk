package test

import (
	"fmt"
	"testing"

	"github.com/Allan-Nava/OvenMediaEngine-go-sdk/ovenmedia"
)

func Test_HeaderInit(t *testing.T) {
	fmt.Println("header init()")
	//
	headerconfig := ovenmedia.InitHeaderConfigurator()
	headerconfig.SetHeader("Content-Type", "application/json")
	headerconfig.SetHeaders(
		map[string]string{
			"Accept": "application/json",
			"X-Test": "test",
		},
	)
	fmt.Println(headerconfig.GetHeaders())
}
