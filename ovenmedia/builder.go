package ovenmedia

import (
	"log"

	"github.com/go-resty/resty/v2"
)

//type Builder OvenMedia
//
//
func BuildOven(url string, debug bool, header *HeaderConfigurator) (IOvenMediaClient, error) {
	ovenClient := &ovenMedia{
		Url:        url,
		restClient: resty.New(),
	}
	// You can override all below settings and options at request level if you want to
	//--------------------------------------------------------------------------------
	// Host URL for all request. So you can use relative URL in the request
	ovenClient.restClient.SetBaseURL(url)
	if header != nil {
		// Headers for all request
		for h, v := range header.GetHeaders() {
			ovenClient.restClient.SetHeader(h, v)
		}
	}
	//
	if debug {
		ovenClient.restClient.SetDebug(true)
		ovenClient.debug = true
		log.Println("Debug mode is enabled for the IOvenMediaClient client ")
	}
	return ovenClient, nil
}

//

func (o *ovenMedia) debugPrint(data interface{}) {
	if o.debug {
		log.Println(data)
	}
}
