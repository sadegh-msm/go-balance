package balancers

import (
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

func TestLeastLoad_Balance(t *testing.T) {
	expect, err := Build(LeastLoadBalancer, []string{
		"127.0.0.1:8015",
		"127.0.0.1:8016",
		"127.0.0.1:8017",
		"127.0.0.1:8018"},
	)

	expect.Remove("127.0.0.1:8018")
	assert.Equal(t, err, nil)

	expect.Inc("127.0.0.1:8015")
	expect.Inc("127.0.0.1:8016")
	expect.Inc("127.0.0.1:8016")
	expect.Inc("127.0.0.1:8018")
	expect.Done("127.0.0.1:8018")
	expect.Done("127.0.0.1:8016")

	ll := NewLeastLoad([]string{
		"127.0.0.1:8016",
	})
	ll.Remove("127.0.0.1:8018")
	ll.Add("127.0.0.1:8015")
	ll.Add("127.0.0.1:8016")
	ll.Add("127.0.0.1:8017")
	ll.Inc("127.0.0.1:8015")
	ll.Inc("127.0.0.1:8016")
	ll.Inc("127.0.0.1:8016")
	ll.Done("127.0.0.1:8016")

	llHost, _ := ll.Balance("")
	expectHost, _ := expect.Balance("")
	assert.Equal(t, true, reflect.DeepEqual(llHost, expectHost))
}
