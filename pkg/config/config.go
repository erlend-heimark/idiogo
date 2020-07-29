package config

type Config struct {
	Port       string
	Path       string
	DadJokeURL string
}

func New() Config {
	return Config{
		Port:       "8080",
		Path:       "/idiogo/v1",
		DadJokeURL: "https://icanhazdadjoke.com/",
	}
}


