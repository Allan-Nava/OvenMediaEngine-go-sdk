package test
//
import (
	"fmt"
	"testing"
)

func Test_GetAllVirtualHost(t *testing.T) {
	o := GetOvenMedia()
	vr, _ := o.GetAllVirtualHosts()
	/*if err != nil {
		panic(err)
	}
	fmt.Println(vr)*/
	fmt.Println(vr)
}

