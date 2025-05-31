package model

type SystemInfo struct {
	Hostname    string `json:"hostname"`
	Version     string `json:"version"`
	ServiceName string `json:"serviceName"`
	Build       string `json:"build"`
	Env         string `json:"env"`
	BuildDate   string `json:"buildDate"`
}
