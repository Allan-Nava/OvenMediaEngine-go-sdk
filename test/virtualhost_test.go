import (
	"fmt"
	"testing"
)

func Test_GetAllVirtualHost(t *testing.T) {
	o := GetOvenMedia()
	vr, err := o.GetAllVirtualHosts()
	if err != nil {
		panic(err)
	}
	fmt.Println(vr)
}

