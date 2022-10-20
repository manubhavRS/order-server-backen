package redisProvider

import (
	"OrderServer/providers"
	"context"
	"github.com/gomodule/redigo/redis"
	"github.com/sirupsen/logrus"
)

type redisProvider struct {
	redisPub redis.Conn
	ctx      context.Context
	psc      redis.PubSubConn
}

func NewRedisProvider() providers.RedisProvider {
	RedisPub, err := redis.Dial("tcp", "localhost:6379")
	if err != nil {
		logrus.Info("Unable to connect to redis client")
	}
	RedisSub, err := redis.Dial("tcp", "localhost:6379")
	if err != nil {
		logrus.Info("Unable to connect to redis client")
	}
	ctx := context.Background()
	psc := redis.PubSubConn{Conn: RedisSub}
	err = psc.Subscribe("order")
	if err != nil {
		logrus.Info("Unable to subscribe to order channel")
	}
	err = psc.Subscribe("details")
	if err != nil {
		logrus.Info("Unable to subscribe to details channel")
	}
	return &redisProvider{
		redisPub: RedisPub,
		ctx:      ctx,
		psc:      psc,
	}
}
func (r redisProvider) Get() (string, string, error) {
	switch v := r.psc.Receive().(type) {
	case redis.Message:
		return v.Channel, string(v.Data), nil
	case error:
		return "", "", v
	}
	return "", "", nil
}

func (r redisProvider) Publish(key string, value interface{}) error {
	_, err := r.redisPub.Do("PUBLISH", key, value)
	if err != nil {
		return err
	}
	return nil
}
