package main

import (
	"io"
	"os"
	"path"

	database "ctigroupjsc.com/phuocnn/gps-management/database/mongo"
	dbRedis "ctigroupjsc.com/phuocnn/gps-management/database/redis"
	"ctigroupjsc.com/phuocnn/gps-management/database/repository"
	"ctigroupjsc.com/phuocnn/gps-management/uitilities/providers/jwt"
	"ctigroupjsc.com/phuocnn/gps-management/uitilities/providers/logger"
	"ctigroupjsc.com/phuocnn/gps-management/uitilities/providers/mongo"
	"ctigroupjsc.com/phuocnn/gps-management/uitilities/providers/redis"

	"gopkg.in/natefinch/lumberjack.v2"
)

const DEFAULT_LOG_FILE_NAME = "general.log"

type Provider struct {
	*logger.LoggerProvider
	*mongo.MongoProvider
	*redis.RedisProvider
	*jwt.JWTService
}

type Container struct {
	*Provider

	Config Config

	UserRepository         repository.UserRepository
	UserAccessIDRepository repository.UserAccessIDRepository
	AccountRepository      repository.AccountRepository
	DeviceRepository       repository.DeviceRepository
	ActivityLogRepository  repository.ActivityLogRepository
	ReportRepository       repository.ReportRepository
	ContactRepository      repository.ContactRepository
}

func NewContainer(config Config) (*Container, error) {
	container := new(Container)
	err := container.InitContainer(config)
	if err != nil {
		return nil, err
	}

	container.Config = config

	return container, nil
}

func (container *Container) InitContainer(config Config) error {
	// Load providers into container
	err := container.LoadProviders(config)
	if err != nil {
		return err
	}

	// Load repositories
	container.LoadRepositoryImplementations(config)

	return nil
}

func (container *Container) LoadProviders(config Config) error {
	loggerProvider := new(logger.LoggerProvider)
	logFileName := DEFAULT_LOG_FILE_NAME
	if config.LogFileName != "" {
		logFileName = config.LogFileName
	}

	loggerProvider = logger.NewLoggerProvider(config.LogLevel)
	if config.LogDir != "" {
		loggerProvider = logger.NewLoggerProvider(config.LogLevel, io.MultiWriter(
			os.Stdout,
			&lumberjack.Logger{
				Filename:   path.Join(config.LogDir, logFileName),
				MaxSize:    100, // megabytes
				MaxBackups: 3,
				MaxAge:     30,   // days
				Compress:   true, // disabled by default
			},
		))
	}

	redisProvider := redis.NewRedisProviderFromURL(config.RedisURL)

	container.Provider = &Provider{
		LoggerProvider: loggerProvider,
		MongoProvider:  mongo.NewMongoProviderFromURL(config.MongoURL),
		RedisProvider:  redisProvider,
		JWTService:     jwt.NewJwtService(config.JwtSecret),
	}
	return nil
}

func (container *Container) LoadRepositoryImplementations(config Config) {
	container.UserRepository = database.NewUserMongoRepository(container.MongoProvider)
	container.UserAccessIDRepository = dbRedis.NewUserAccessIDRedisRepository(container.RedisProvider)
	container.AccountRepository = database.NewAccountMongoRepository(container.MongoProvider)
	container.DeviceRepository = database.NewDeviceMongoRepository(container.MongoProvider)
	container.ActivityLogRepository = database.NewActivityLogMongoRepository(container.MongoProvider)
	container.ReportRepository = database.NewReportMongoRepository(container.MongoProvider)
	container.ContactRepository = database.NewContactMongoRepository(container.MongoProvider)
}
