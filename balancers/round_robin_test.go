package balancers

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

// TestRoundRobin_Add
func TestRoundRobin_Add(t *testing.T) {
	cases := []struct {
		name   string
		lb     Balancer
		args   string
		expect Balancer
	}{
		{
			"test-1",
			NewRoundRobin([]string{
				"http://127.0.0.1:1011",
				"http://127.0.0.1:1012",
				"http://127.0.0.1:1013"},
			),
			"http://127.0.0.1:1013",
			&RoundRobin{
				hosts: []string{
					"http://127.0.0.1:1011",
					"http://127.0.0.1:1012",
					"http://127.0.0.1:1013"},
				num: 0},
		},
		{
			"test-2",
			NewRoundRobin([]string{
				"http://127.0.0.1:1011",
				"http://127.0.0.1:1012"},
			),
			"http://127.0.0.1:1012",
			&RoundRobin{
				hosts: []string{
					"http://127.0.0.1:1011",
					"http://127.0.0.1:1012"},
				num: 0},
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			c.lb.Add(c.args)
			assert.Equal(t, c.expect, c.lb)
		})
	}
}

// TestRoundRobin_Remove
func TestRoundRobin_Remove(t *testing.T) {
	cases := []struct {
		name   string
		lb     Balancer
		args   string
		expect Balancer
	}{
		{
			"test-1",
			NewRoundRobin([]string{
				"http://127.0.0.1:1011",
				"http://127.0.0.1:1012"},
			),
			"http://127.0.0.1:1013",
			&RoundRobin{hosts: []string{
				"http://127.0.0.1:1011",
				"http://127.0.0.1:1012"},
			},
		},
		{
			"test-2",
			NewRoundRobin([]string{
				"http://127.0.0.1:1011",
				"http://127.0.0.1:1012"},
			),
			"http://127.0.0.1:1012",
			&RoundRobin{hosts: []string{
				"http://127.0.0.1:1011"},
			},
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			c.lb.Remove(c.args)
			assert.Equal(t, c.expect, c.lb)
		})
	}
}

// TestRoundRobin_Balance .
func TestRoundRobin_Balance(t *testing.T) {
	type expect struct {
		reply string
		err   error
	}
	cases := []struct {
		name   string
		lb     Balancer
		args   string
		expect expect
	}{
		{
			"test-1",
			NewRoundRobin([]string{
				"http://127.0.0.1:1011"},
			),
			"",
			expect{
				"http://127.0.0.1:1011",
				nil,
			},
		},
		{
			"test-2",
			NewRoundRobin([]string{}),
			"",
			expect{
				"",
				NoHostError,
			},
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			reply, err := c.lb.Balance(c.args)
			assert.Equal(t, c.expect.reply, reply)
			assert.Equal(t, c.expect.err, err)
		})
	}
}
