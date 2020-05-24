package cmd

import (
	"fmt"
	"github.com/linuxsuren/net-proxy/nexus"
	"github.com/spf13/cobra"
	"net/http"
	"net/url"
	"os"
)

// NewRootCmd returns the root command
func NewRootCmd() (cmd *cobra.Command) {
	opt := &RootOptions{}

	cmd = &cobra.Command{
		Use:   "net-proxy",
		Short: "Network proxy",
		RunE:  opt.RunE,
	}

	flags := cmd.Flags()
	flags.StringVarP(&opt.Bind, "bind", "", "127.0.0.1", "The address which server will bind")
	flags.IntVarP(&opt.BindPort, "bind-port", "", 8089, "The port which will bind")
	flags.IntVarP(&opt.ProxyRedirectCode, "proxy-redirect-code", "", http.StatusMovedPermanently,
		"The status code when redirect the HTTP request")
	flags.StringVarP(&opt.NexusRawProxyPrefix, "nexus-raw-proxy-prefix", "", "raw-proxy",
		"The prefix of nexus raw proxy repository")
	flags.IntVarP(&opt.ProxyBuffer, "proxy-buffer", "", 2048*1024,
		"The buffer when transfer data")
	flags.StringVarP(&opt.Remote, "remote", "", "127.0.0.1:8081", "The remote proxy server")
	flags.StringVarP(&opt.NexusServer, "nexus-server", "", "http://localhost:8081", "The Nexus server address")
	flags.StringVarP(&opt.Username, "username", "", "admin", "Username of server(Nexus)")
	flags.StringVarP(&opt.Password, "password", "", "admin", "Password of server(Nexus)")

	cmd.AddCommand(NewVersionCmd())
	return
}

func (o *RootOptions) findRawProxy(remote *url.URL) (proxy string, ok bool) {
	client := nexus.NexusClient{
		URL:      o.NexusServer,
		Username: o.Username,
		Password: o.Password,
	}

	if repo, err := client.ListRepositories(); err == nil {
		fmt.Println(repo)
		for _, re := range repo {
			if re.Name == o.getProxyID(remote) {
				proxy = fmt.Sprintf("/repository/%s", o.getProxyID(remote))
				ok = true
				break
			}
		}
	} else {
		fmt.Printf("%#v\n", err)
	}
	return
}

func (o *RootOptions) addRawProxy(remote *url.URL) (proxy string) {
	proxy = o.getProxyID(remote)

	client := nexus.NexusClient{
		URL:      o.NexusServer,
		Username: o.Username,
		Password: o.Password,
	}
	rawProxy := nexus.RawProxy{
		Name:   proxy,
		Online: true,
		Format: "format",
		Recipe: "raw-proxy",
		Storage: nexus.Storage{
			BlobStoreName:               "default",
			StrictContentTypeValidation: false,
		},
		NegativeCache: nexus.NegativeCache{
			Enabled:    false,
			TimeToLive: 0,
		},
		Proxy: nexus.Proxy{
			RemoteURL: fmt.Sprintf("https://%s", remote.Host),
		},
		HTTPClient: nexus.HTTPClient{
			Blocked:   false,
			AutoBlock: false,
			Connection: nexus.HTTPClientConnection{
				Timeout: 60,
			},
			Authentication: nexus.HTTPClientAuthentication{
				Type: "none",
			},
		},
	}
	if err := client.CreateProxy(rawProxy); err != nil {
		fmt.Println("add raw-proxy repository error", err)
	}
	return
}

func (o *RootOptions) getProxyID(remote *url.URL) string {
	return fmt.Sprintf("%s-%s", o.NexusRawProxyPrefix, remote.Host)
}

// RunE is the main entry of root command
func (o *RootOptions) RunE(cmd *cobra.Command, args []string) (err error) {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("server", r.URL.Host, o.Remote)
		if r.URL.Host == o.Remote {
			client := &http.Client{}
			req := &http.Request{
				Method: r.Method,
				URL:    r.URL,
			}
			fmt.Println("start to connect with proxy server", r.URL)
			if rs, err := client.Do(req); err == nil {
				fmt.Println("get response", rs)
				for key, val := range rs.Header {
					for _, item := range val {
						fmt.Println(key, item)
						w.Header().Add(key, item)
					}
				}

				buf := make([]byte, o.ProxyBuffer)
				for count, _ := rs.Body.Read(buf); count > 0; count, err = rs.Body.Read(buf) {
					w.Write(buf[:count])
				}
				w.WriteHeader(rs.StatusCode)
			} else {
				cmd.Println(err)
			}
		} else {
			var rawProxy string
			var ok bool
			if rawProxy, ok = o.findRawProxy(r.URL); !ok {
				cmd.Println("cannot found", rawProxy)
				rawProxy = o.addRawProxy(r.URL)
			}

			proxy := url.URL{
				Scheme: "http",
				Host:   o.Remote,
				Path:   fmt.Sprintf("%s%s", rawProxy, r.URL.Path),
			}
			w.Header().Set("Location", proxy.String())
			w.WriteHeader(o.ProxyRedirectCode)
		}
	})
	err = http.ListenAndServe(fmt.Sprintf("%s:%d", o.Bind, o.BindPort), nil)
	return
}

// Execute will execute the command
func Execute() {
	if err := NewRootCmd().Execute(); err != nil {
		os.Exit(1)
	}
}
