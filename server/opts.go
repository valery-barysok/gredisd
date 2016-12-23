package server

// Options struct is configuration settings of Server
type Options struct {
	Host string `json:"addr"`
	Port int    `json:"port"`
}
