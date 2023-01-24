package ovenmedia

import (
	"log"

	"github.com/go-resty/resty/v2"
)

//type Builder OvenMedia
//
//
func BuildOven(url string, debug bool) (*OvenMedia, error) {
	ovenClient := &OvenMedia{
		Url:        url,
		RestClient: resty.New(),
	}
	// You can override all below settings and options at request level if you want to
	//--------------------------------------------------------------------------------
	// Host URL for all request. So you can use relative URL in the request
	ovenClient.RestClient.SetHostURL(url)
	// Headers for all request
	ovenClient.RestClient.SetHeader("ome-access-token", "Basic ")
	//
	if debug {
		ovenClient.RestClient.SetDebug(true)
		ovenClient.Debug = true
		log.Println("Debug mode is enabled for the Haproxy client ")
	}
	return ovenClient, nil
}

//

func (o *OvenMedia) DebugPrint(data interface{}) {
	if o.Debug {
		log.Println(data)
	}
}
