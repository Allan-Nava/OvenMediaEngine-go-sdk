package test

import (
	"testing"

	"github.com/Allan-Nava/OvenMediaEngine-go-sdk/ovenmedia"
)

func TestMain(m *testing.M) {
	m.Run()
}

//
func GetOvenMedia() ovenmedia.IOvenMediaClient {
	//
	o, err := ovenmedia.BuildOven("https://airensoft.gitbook.io/ovenmediaengine/rest-api/v1/virtualhost/application/stream", true, nil)
	if err != nil {
		panic(err)
	}
	return o
}

func Test_InitOvenMedia(t *testing.T) {
	o := GetOvenMedia()
	if o == nil {
		panic("Init Oven Media problem")
	}
}

func Test_HealthCheck(t *testing.T) {
	o := GetOvenMedia()
	if err := o.HealthCheck(); err != nil {
		panic(err)
	}
}
