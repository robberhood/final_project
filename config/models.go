package config

type Config struct {
	DBPath string `yaml:"dbpath"`
	Port   string `yaml:"port"`
}

type Task struct {
	ID      string `json:"id"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

type TasksResp struct {
	Tasks []*Task `json:"tasks"`
}
