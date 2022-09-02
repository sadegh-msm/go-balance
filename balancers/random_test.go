package balancers

import (
	"github.com/stretchr/testify/assert"
	"math/rand"
	"testing"
	"time"
)

func TestRandom_Add(t *testing.T) {
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))

	cases := []struct {
		name   string
		lb     Balancer
		args   string
		expect Balancer
	}{
		{
			"test-1",
			&Random{
				hosts: []string{
					"http://127.0.0.1:8011",
					"http://127.0.0.1:8012",
					"http://127.0.0.1:8013",
				},
				rnd: rnd,
			},

			"http://127.0.0.1:8013",
			&Random{
				hosts: []string{
					"http://127.0.0.1:8011",
					"http://127.0.0.1:8012",
					"http://127.0.0.1:8013",
				},
				rnd: rnd,
			},
		},
		{
			"test-2",
			&Random{
				hosts: []string{
					"http://127.0.0.1:8011",
					"http://127.0.0.1:8012",
				},
				rnd: rnd,
			},
			"http://127.0.0.1:8012",
			&Random{
				hosts: []string{
					"http://127.0.0.1:8011",
					"http://127.0.0.1:8012",
				},
				rnd: rnd,
			},
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			c.lb.Add(c.args)
			assert.Equal(t, c.expect, c.lb)
		})
	}
}

func TestRandom_Remove(t *testing.T) {
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))

	cases := []struct {
		name   string
		lb     Balancer
		args   string
		expect Balancer
	}{
		{
			"test-1",
			&Random{
				hosts: []string{
					"http://127.0.0.1:8011",
					"http://127.0.0.1:8012",
					"http://127.0.0.1:8013",
				},
				rnd: rnd,
			},
			"http://127.0.0.1:8013",
			&Random{
				hosts: []string{
					"http://127.0.0.1:8011",
					"http://127.0.0.1:8012",
				},
				rnd: rnd,
			},
		},
		{
			"test-2",
			&Random{
				hosts: []string{
					"http://127.0.0.1:8011",
					"http://127.0.0.1:8012",
				},
				rnd: rnd,
			},
			"http://127.0.0.1:8013",
			&Random{
				hosts: []string{
					"http://127.0.0.1:8011",
					"http://127.0.0.1:8012",
				},
				rnd: rnd,
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

func TestRandom_Balance(t *testing.T) {
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
			NewRandom(
				[]string{
					"http://127.0.0.1:8011"},
			),
			"",
			expect{
				"http://127.0.0.1:8011",
				nil,
			},
		},
		{
			"test-2",
			NewRandom(
				[]string{}),
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
