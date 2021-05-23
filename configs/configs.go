package configs

import "fmt"

type main struct {
	Host string
	Port int
}

type postgres struct {
	User     string
	Password string
	DBName   string
	Host     string
	Port     int
}

type config struct {
	Main     main
	Postgres postgres
}

var Configs config

func (c *config) GetMainHost() string {
	return c.Main.Host
}

func (c *config) GetMainPort() int {
	return c.Main.Port
}

func (c *config) GetPostgresConfig() string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		c.Postgres.Host, c.Postgres.Port, c.Postgres.User, c.Postgres.Password, c.Postgres.DBName)
}

func init() {
	Configs = config{
		Main: main{
			"localhost",
			5000,
		},
		Postgres: postgres{
			"docker",
			"docker",
			"docker",
			"localhost",
			5432,
		},
	}
}
