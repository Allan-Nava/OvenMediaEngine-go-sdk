package ovenmedia

import (
	"log"

	"github.com/go-resty/resty/v2"
)

//type Builder OvenMedia
//
//
func BuildOven(url string, debug bool, header *HeaderConfigurator) (*OvenMedia, error) {
	ovenClient := &OvenMedia{
		Url:        url,
		restClient: resty.New(),
	}
	// You can override all below settings and options at request level if you want to
	//--------------------------------------------------------------------------------
	// Host URL for all request. So you can use relative URL in the request
	ovenClient.restClient.SetHostURL(url)
	// Headers for all request
	for h, v := range header.GetHeaders() {
		ovenClient.restClient.SetHeader(h, v)
	}
	//
	if debug {
		ovenClient.restClient.SetDebug(true)
		ovenClient.debug = true
		log.Println("Debug mode is enabled for the OvenMedia client ")
	}
	return ovenClient, nil
}

//

func (o *OvenMedia) debugPrint(data interface{}) {
	if o.debug {
		log.Println(data)
	}
}
