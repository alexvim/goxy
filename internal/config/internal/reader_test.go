package configreader

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReader(t *testing.T) {
	t.Run("read_test", func(t *testing.T) {
		tests := []struct {
			l  string
			p  string
			j  string
			ok bool
		}{
			{
				l:  "172.0.0.1",
				p:  "127.0.0.1:1080",
				j:  `{"proxy_address": "127.0.0.1:1080", "local_address": "172.0.0.1"}`,
				ok: true,
			},
			{
				l:  "172.0.0.4",
				p:  "127.2.0.1:1080",
				j:  `{"proxy_address": "127.2.0.1:1080", "local_address": "172.0.0.4"}`,
				ok: true,
			},
			{
				l:  "172.0.0.1",
				p:  "",
				j:  `{"proxy_addrss": "127.0.0.1:1080", "local_address": "172.0.0.1"}`,
				ok: true,
			},
			{
				l:  "",
				p:  "127.0.0.1:1080",
				j:  `{"proxy_address": "127.0.0.1:1080", "loca_address": "172.0.0.1"}`,
				ok: true,
			},
			{
				l:  "",
				p:  "127.0.0.1:1080",
				j:  `{"proxy_address": "127.0.0.1:1080", "loca_address": "172.0.0.1"`,
				ok: false,
			},
		}

		for n, test := range tests {
			t.Run(fmt.Sprintf("read_test_%d", n), func(t *testing.T) {
				t.Parallel()

				cfg, err := readFromJson(strings.NewReader(test.j))

				if !test.ok {
					assert.Error(t, err)
					return
				}

				assert.NoError(t, err)

				assert.Equal(t, test.p, cfg.ProxyAddr)
				assert.Equal(t, test.l, cfg.LocalAddr)
			})
		}
	})

	t.Run("args_parse_test", func(t *testing.T) {
		tests := []struct {
			l string
			p string
			c string
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
