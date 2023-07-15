package api

type Task struct {
	Targets []string `json:"targets"`
	Port    string   `json:"port"`
	Delay   int      `json:"delay"`
	Output  string   `json:"output"`
}
