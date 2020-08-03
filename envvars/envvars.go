package envvars

type EnvVars struct {
	Port               string
	Path               string
	DadJokeURL         string
	DBConnectionString string
}

func New() EnvVars {
	return EnvVars{
		Port:               "8080",
		Path:               "/idiogo/v1",
		DadJokeURL:         "https://icanhazdadjoke.com/",
		DBConnectionString: "sqlserver://localhost:1433?user id=sa&password=sOdifn3ijnvsd8!&database=dadjokes",
	}
}
