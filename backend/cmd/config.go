package main

type Config struct {
	LogLevel    int
	LogDir      string
	LogFileName string

	JwtSecret string

	Binding string

	MongoURL string
	RedisURL string

	MongoServer     string
	MongoDatabase   string
	MongoCollection string
	MongoSource     string
	MongoUser       string
	MongoPassword   string
}
