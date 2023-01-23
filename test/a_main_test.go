package test

import (
	"testing"

	"github.com/Allan-Nava/OvenMediaEngine-go-sdk/ovenmedia"
)

func TestMain(m *testing.M) {
	m.Run()
}

//
func GetOvenMedia() *ovenmedia.OvenMedia {
	//
	o, err := ovenmedia.BuildOven("https://airensoft.gitbook.io/ovenmediaengine/rest-api/v1/virtualhost/application/stream", true)
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
	return
}
