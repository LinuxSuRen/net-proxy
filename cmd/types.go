package cmd

import "github.com/linuxsuren/net-proxy/cmd/common"

// RootOptions is the option of root command
type RootOptions struct {
	common.Options

	Bind              string
	BindPort          int
	Remote            string
	ProxyBuffer       int
	ProxyRedirectCode int

	NexusRawProxyPrefix string

	NexusServer string
	Username    string
	Password    string
}
