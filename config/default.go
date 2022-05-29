package config

//go:generate stringer -type CipherSuit -linecomment
//CipherSuit cipher suit
type CipherSuit int

//cipher suit enum
//nolint
const (
	//TLS_AES_128_GCM_SHA256
	TLS_AES_128_GCM_SHA256 CipherSuit = iota
	CipherSuit_MAX
)

var parseCipherSuit = make(map[string]CipherSuit, int(CipherSuit_MAX))

func init() {
	for i := 0; i < int(CipherSuit_MAX); i++ {
		c := CipherSuit(i)
		parseCipherSuit[c.String()] = c
	}
}

var defaultConfig = &Struct{
	TLS: TLSConfigStruct{
		CLientAuth:            true,
		CipherSuit:            []string{"TLS_AES_128_GCM_SHA256"},
		ClientCertCommandName: "*.",
		DisableTLSLow:         true,
	},
	IP: IPConfigStruct{
		Out: "0.0.0.0",
		In:  "127.0.0.1",
	},
	PortMap: []PortMapItem{
		{FromPort: 8080, ToPort: 8081},
	},
	Alerts: []AlertItem{},
}
