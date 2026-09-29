package test

import (
	"testing"

	"github.com/Allan-Nava/OvenMediaEngine-go-sdk/ovenmedia"
)

func Test_HeaderInit(t *testing.T) {
	h := ovenmedia.InitHeaderConfigurator()
	h.SetHeader("X-Test", "test")
	h.CreateBasicAuthHeader("admin", "secret")

	if got := h.GetHeader("Authorization"); got == nil || *got != "Basic YWRtaW46c2VjcmV0" {
		t.Errorf("Authorization = %v", got)
	}
	if !h.HasHeader("X-Test") || len(h.GetHeaderKeys()) != 2 {
		t.Errorf("headers = %v", h.GetHeaders())
	}
	h.DeleteHeaders()
	if h.HasHeaders() {
		t.Error("DeleteHeaders left headers behind")
	}
}
