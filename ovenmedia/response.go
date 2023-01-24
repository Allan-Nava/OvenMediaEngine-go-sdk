package ovenmedia

// Response<VirtualHost>
type ResponseVirtualHost struct {
	VRHosts []VirtualHost
}

type VirtualHost struct {
	Message    string         `json:"message"`
	Response   VRHostResponse `json:"response"`
	StatusCode int            `json:"statusCode"`
}

type VRHostResponse struct {
	Name string `json:"name"`
}
