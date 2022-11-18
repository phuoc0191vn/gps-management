package redis

import (
	"fmt"
	"log"
	"net"
	"net/url"
	"strconv"
	"strings"
	"time"

	"ctigroupjsc.com/phuocnn/gps-management/uitilities/conversion"

	"github.com/pkg/errors"
	"gopkg.in/redis.v5"
)

type PubSub *redis.PubSub

type RedisProvider struct {
	url      string
	server   string
	password string
	db       int
	redis    *RedisClient
}

type DatabaseExecutionError struct {
	Message string
}

func (e DatabaseExecutionError) Error() string {
	return e.Message
}

func NewRedisProviderFromURL(u string) *RedisProvider {
	client := newRedisClientFromURL(u)
	if client == nil {
		log.Fatalln("Redis server connected unsuccessfully")
	}
	return &RedisProvider{
		url:   u,
		redis: client,
	}
}

func NewRedisProvider(server string, password string, db int) *RedisProvider {
	client := newRedisClient(server, password, db)
	if client == nil {
		log.Fatalln("Redis server connected unsuccessfully")
	}
	return &RedisProvider{
		server:   server,
		password: password,
		db:       db,
		redis:    client,
	}
}

func (provider *RedisProvider) RedisClient() *RedisClient {
	return provider.redis
}

func (provider *RedisProvider) NewRedisClient() *RedisClient {
	if provider.url != "" {
		return newRedisClientFromURL(provider.url)
	}
	return newRedisClient(provider.server, provider.password, provider.db)
}

func (provider *RedisProvider) NewError(e error) error {
	if e == nil {
		return nil
	}
	return DatabaseExecutionError{Message: fmt.Sprintf("Redis execution error: %s", e.Error())}
}

type RedisClient struct {
	client *redis.Client
}

func newRedisClient(server string, password string, db int) *RedisClient {
	result := new(RedisClient)
	result.client = redis.NewClient(&redis.Options{
		Addr:         server,
		Password:     password,
		DB:           db,
		ReadTimeout:  5 * time.Minute,
		WriteTimeout: 5 * time.Minute,
	})

	return result
}

func newRedisClientFromURL(u string) *RedisClient {
	tempU, err := url.Parse(u)
	if err != nil {
		log.Fatalln("Redis server connected unsuccessfully")
	}
	result := new(RedisClient)

	switch tempU.Scheme {
	case "redis-sentinel":
		options, err := parseURLSentinel(u)
		if err != nil {
			log.Fatalln("Redis server connected unsuccessfully")
		}

		result.client = redis.NewFailoverClient(options)

	default:
		options, err := redis.ParseURL(u)
		if err != nil {
			log.Fatalln("Redis server connected unsuccessfully")
		}

		result.client = redis.NewClient(options)
	}

	return result
}

func parseURLSentinel(redisURL string) (*redis.FailoverOptions, error) {
	result := new(redis.FailoverOptions)

	o := &redis.Options{Network: "tcp"}
	u, err := url.Parse(redisURL)
	if err != nil {
		return nil, err
	}

	if u.Scheme != "redis-sentinel" {
		return nil, errors.New("invalid redis sentinel URL scheme: " + u.Scheme)
	}

	if u.User == nil {
		return nil, errors.New("invalid redis sentinel URL User: " + u.Scheme)
	}

	if masterName := u.User.Username(); masterName == "" {
		return nil, errors.New("invalid redis sentinel Master Name: " + u.Scheme)
	}

	result.MasterName = u.User.Username()

	if p, ok := u.User.Password(); ok {
		result.Password = p
	}

	if len(u.Query()) > 0 {
		return nil, errors.New("no options supported")
	}

	sentinelHostPorts := strings.Split(u.Host, ",")

	var sentinelAddresses []string

	for _, v := range sentinelHostPorts {
		h, p, err := net.SplitHostPort(v)
		if err != nil {
			return nil, errors.New("error sentinel address")
		}

		sentinelAddresses = append(sentinelAddresses, net.JoinHostPort(h, p))
	}
	result.SentinelAddrs = sentinelAddresses

	f := strings.FieldsFunc(u.Path, func(r rune) bool {
		return r == '/'
	})
	switch len(f) {
	case 0:
		o.DB = 0
	case 1:
		if o.DB, err = strconv.Atoi(f[0]); err != nil {
			return nil, fmt.Errorf("invalid redis database number: %q", f[0])
		}
	default:
		return nil, errors.New("invalid redis URL path: " + u.Path)
	}

	result.ReadTimeout = 5 * time.Minute
	result.WriteTimeout = 5 * time.Minute

	return result, nil
}

// Set key for blacklist/whitelist
func (r *RedisClient) HSet(key, field, value string) error {
	return r.client.HSet(key, field, value).Err()
}

func (r *RedisClient) HDel(key string, fields ...string) error {
	return r.client.HDel(key, fields...).Err()
}

func (r *RedisClient) HGetAll(key string) (map[string]string, error) {
	cmd := r.client.HGetAll(key)
	if err := cmd.Err(); err != nil {
		return map[string]string{}, err
	}

	return cmd.Val(), nil
}

func (r *RedisClient) HGet(key, field string) (string, error) {
	cmd := r.client.HGet(key, field)
	if err := cmd.Err(); err != nil {
		return "", err
	}

	return cmd.Val(), nil
}

func (r *RedisClient) HExists(key, field string) (bool, error) {
	cmd := r.client.HExists(key, field)
	if err := cmd.Err(); err != nil {
		return false, err
	}
	return cmd.Val(), nil
}

func (r *RedisClient) SMembers(key string) ([]string, error) {
	cmd := r.client.SMembers(key)
	if err := cmd.Err(); err != nil {
		return []string{}, err
	}

	return cmd.Val(), nil
}

func (r *RedisClient) SAdd(key string, values ...interface{}) error {
	return r.client.SAdd(key, values...).Err()
}

func (r *RedisClient) SRem(key string, values ...interface{}) error {
	return r.client.SRem(key, values...).Err()
}

func (r *RedisClient) SIsMember(key string, value interface{}) (bool, error) {
	cmd := r.client.SIsMember(key, value)
	if err := cmd.Err(); err != nil {
		return false, err
	}

	return cmd.Val(), nil
}

func (r *RedisClient) Del(key ...string) error {
	return r.client.Del(key...).Err()
}

func (r *RedisClient) HMGet(key string, fields ...string) ([]interface{}, error) {
	cmd := r.client.HMGet(key, fields...)

	if err := cmd.Err(); err != nil {
		return nil, err
	}

	return cmd.Val(), nil
}

func (r *RedisClient) HMSet(key string, fields map[string]string) error {
	return r.client.HMSet(key, fields).Err()
}

func (r *RedisClient) RPush(key string, values ...interface{}) error {
	return r.client.RPush(key, values...).Err()
}

func (r *RedisClient) BLPop(timeout time.Duration, keys ...string) ([]string, error) {
	return r.client.BLPop(timeout, keys...).Result()
}

func (r *RedisClient) Set(key string, val interface{}, expire time.Duration) error {
	return r.client.Set(key, val, expire).Err()
}

func (r *RedisClient) Get(key string) (string, error) {
	return r.client.Get(key).Result()
}

func (r *RedisClient) Expire(key string, expire time.Duration) error {
	return r.client.Expire(key, expire).Err()
}

func (r *RedisClient) Publish(channel string, command interface{}) error {
	return r.client.Publish(channel, conversion.ToJson(command)).Err()
}

func (r *RedisClient) Subscribe(channels ...string) (*redis.PubSub, error) {
	return r.client.Subscribe(channels...)
}

func (r *RedisClient) Keys(pattern string) []string {
	return r.client.Keys(pattern).Val()
}

func (r *RedisClient) LRange(key string, start, stop int64) ([]string, error) {
	return r.client.LRange(key, start, stop).Result()
}

func (r *RedisClient) LTrim(key string, start, stop int64) error {
	return r.client.LTrim(key, start, stop).Err()
}

func (r *RedisClient) Ping() error {
	return r.client.Ping().Err()
}
