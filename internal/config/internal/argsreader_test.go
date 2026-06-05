package configreader

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReader(t *testing.T) {
	t.Run("args_parse_test", func(t *testing.T) {
		tests := []struct {
			l string
			p string
			c string
			d string
			a []string
		}{
			{
				l: "172.0.0.1",
				p: "127.0.0.1:1080",
				a: []string{"-l", "172.0.0.1", "-p", "127.0.0.1:1080"},
			},
			{
				c: "/home/user/conf.json",
				a: []string{"-c", "/home/user/conf.json"},
			},
			{
				l: "172.0.0.1",
				p: "127.0.0.1:1080",
				c: "/home/user/conf.json",
				a: []string{"-l", "172.0.0.1", "-p", "127.0.0.1:1080", "-c", "/home/user/conf.json"},
			},
			{
				l: "172.0.0.1",
				p: "127.0.0.1:1080",
				d: "doh.example.com",
				a: []string{"-l", "172.0.0.1", "-p", "127.0.0.1:1080", "-doh", "doh.example.com"},
			},
		}

		for n, test := range tests {
			t.Run(fmt.Sprintf("args_parse_test_%d", n), func(t *testing.T) {
				t.Parallel()

				cmds := cmdArgsDefaults
				cmdArgs := &cmds

				cmdArgs.parse(test.a)

				assert.Equal(t, test.p, cmdArgs.proxyAddr.val)
				assert.Equal(t, test.l, cmdArgs.localAddr.val)
				assert.Equal(t, test.c, cmdArgs.jsonConfig.val)
			})
		}
	})
}
