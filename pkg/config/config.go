package config

type Config struct {
	Port       string
	Path       string
	DadJokeURL string
	ConnectionString string
}

func New() Config {
	return Config{
		Port:       "8080",
		Path:       "/idiogo/v1",
		DadJokeURL: "https://icanhazdadjoke.com/",
		ConnectionString: "sqlserver://localhost:1433?user id=sa&password=sOdifn3ijnvsd8!&database=dadjokes",
	}
}


