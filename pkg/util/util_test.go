package util

import (
	"testing"

	"github.com/matryer/is"
)

func TestLookup(t *testing.T) {
	is := is.New(t)
	var domains = []string{`cisco.com`, `ibm.com`, `microsoft.com`}

	for _, domain := range domains {
		addresses, err := GetDomainAddresses(domain)
		is.NoErr(err)
		for _, addr := range addresses {
			t.Logf("Domain %s", domain)
			if addr.Is4() {
				t.Log("ip4: ", addr.String())
				continue
			}
			t.Log("ip6: ", addr.String())
		}
	}
}
