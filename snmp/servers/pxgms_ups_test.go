package servers

import (
	"fmt"
	"testing"
)

// TestPxgmsUps is the first PxgmsUps test.
func TestPxgmsUps(t *testing.T) {
	fmt.Printf("TestPxgmUps start\n")
	fmt.Printf("t: %+v\n", t)

	data := make(map[string]interface{})
	data["contextName"] = "public"
	data["endpoint"] = "127.0.0.1"
	data["userName"] = "simulator"
	data["privacyProtocol"] = "AES"
	data["privacyPassphrase"] = "privatus"
	data["port"] = "1024"
	data["authenticationProtocol"] = "SHA"
	data["authenticationPassphrase"] = "auctoritas"
	data["model"] = "PXGMS UPS + EATON 93PM"
	data["version"] = "v3"

	pxgmsUps, err := NewPxgmsUps(data)
	if err != nil {
		t.Fatal(err)
	}
	// TODO: Need to do more with this, but at least excersizes the code for now.
	fmt.Printf("pxgmsUps: %+v\n", pxgmsUps)
}
