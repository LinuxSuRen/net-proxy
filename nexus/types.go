package nexus

type NexusClient struct {
	URL      string
	Username string
	Password string
}

// Repository represents a repository from a Nexus server
type Repository struct {
	Name   string
	Format string
	Type   string
	URL    string
	Online bool
}

type RawProxy struct {
	Name          string        `json:"name"`
	Format        string        `json:"format"`
	Online        bool          `json:"online"`
	Storage       Storage       `json:"storage"`
	Cleanup       Cleanup       `json:"cleanup"`
	Proxy         Proxy         `json:"proxy"`
	NegativeCache NegativeCache `json:"negativeCache"`
	HTTPClient    HTTPClient    `json:"httpclient"`
	RoutingRule   string        `json:"routingRule"`
	Recipe        string        `json:"recipe"`
}

type Storage struct {
	BlobStoreName               string `json:"blobStoreName"`
	StrictContentTypeValidation bool   `json:"strictContentTypeValidation"`
}

type Cleanup struct {
	PolicyName []string `json:"policyName"`
}

type Proxy struct {
	RemoteURL      string `json:"remoteUrl"`
	ContentMaxAge  int    `json:"contentMaxAge"`
	MetadataMaxAge int    `json:"metadataMaxAge"`
}

type NegativeCache struct {
	Enabled    bool `json:"enabled"`
	TimeToLive int  `json:"timeToLive"`
}

type HTTPClient struct {
	Blocked        bool                     `json:"blocked"`
	AutoBlock      bool                     `json:"autoBlock"`
	Connection     HTTPClientConnection     `json:"connection"`
	Authentication HTTPClientAuthentication `json:"authentication"`
}

type HTTPClientConnection struct {
	Retries                 int    `json:"retries"`
	UserAgentSuffix         string `json:"userAgentSuffix"`
	Timeout                 int    `json:"timeout"`
	EnableCircularRedirects bool   `json:"enableCircularRedirects"`
	EnableCookies           bool   `json:"enableCookies"`
}

type HTTPClientAuthentication struct {
	Type       string `json:"type"`
	Username   string `json:"username"`
	NtlmHost   string `json:"ntlmHost"`
	NtlmDomain string `json:"ntlmDomain"`
}

type UIPayload struct {
	Action string   `json:"action"`
	Data   []UIData `json:"data"`
	Method string   `json:"method"`
	TID    int      `json:"tid"`
	Type   string   `json:"type"`
}

type UIData struct {
	RawProxy   `json:",inline"`
	Attributes UIDataAttributes `json:"attributes"`
}

type UIDataAttributes struct {
	Cleanup       Cleanup       `json:"cleanup"`
	HTTPClient    HTTPClient    `json:"httpcclient"`
	NegativeCache NegativeCache `json:"negativeCache"`
	Proxy         Proxy         `json:"proxy"`
	Storage       Storage       `json:"storage"`
}
