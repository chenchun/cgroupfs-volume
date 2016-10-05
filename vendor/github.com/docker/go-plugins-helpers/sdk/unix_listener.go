// +build linux freebsd

package sdk

import (
	"net"
	"fmt"
)

func newUnixListener(pluginName string, group string) (net.Listener, string, error) {
	return nil, "", fmt.Errorf("not implemented")
}
