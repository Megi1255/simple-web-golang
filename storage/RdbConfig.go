package storage

const (
	DefaultRdbHost     = "127.0.0.1"
	DefaultRdbPort     = 3306
	DefaultRdbDatabase = "gin"
	DefaultRdbUser     = "user"
	DefaultRdbPasswd   = "password"

	TbNameUser    = "USER"
	TbNameGiftBox = "GIFTBOX"
)

type Config struct {
	Host   string
	Port   int
	DbName string
	User   string
	Passwd string
}

func DefaultConfig() *Config {
	return &Config{
		Host:   DefaultRdbHost,
		Port:   DefaultRdbPort,
		DbName: DefaultRdbDatabase,
		User:   DefaultRdbUser,
		Passwd: DefaultRdbPasswd,
	}
}
