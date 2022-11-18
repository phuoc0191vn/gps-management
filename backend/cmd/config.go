package main

type Config struct {
	LogLevel    int
	LogDir      string
	LogFileName string

	JwtSecret string

	Binding string

	MongoURL string
	RedisURL string
}
