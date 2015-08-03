package main

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/garyburd/redigo/redis"
)

type Alerter struct {
	pool *redis.Pool
}

func NewAlerter(server, password string) (*Alerter, error) {
	return &Alerter{
		pool: &redis.Pool{
			MaxIdle:     3,
			IdleTimeout: 240 * time.Second,
			Dial: func() (redis.Conn, error) {
				c, err := redis.Dial("tcp", server)
				if err != nil {
					return nil, err
				}
				if password != "" {
					if _, err := c.Do("AUTH", password); err != nil {
						c.Close()
						return nil, err
					}
				}
				return c, err
			},
			TestOnBorrow: func(c redis.Conn, t time.Time) error {
				_, err := c.Do("PING")
				return err
			},
		},
	}, nil
}

func (a *Alerter) Set(key string, value []byte, expire int) {
	conn := a.pool.Get()
	defer conn.Close()
	if _, err := conn.Do("SET", key, value); err != nil {
		fmt.Printf("ERROR: unable to set redis key %s\n", err.Error())
	}
}

func (a *Alerter) SetWithExpires(key string, value []byte, expire int) {
	conn := a.pool.Get()
	defer conn.Close()
	if _, err := conn.Do("SET", key, value, "EX", expire); err != nil {
		fmt.Printf("ERROR: unable to set redis key %s\n", err.Error())
	}
}

func (a *Alerter) SetHash(key, field string, value []byte) {
	conn := a.pool.Get()
	defer conn.Close()
	// expire the key after 5 minutes
	conn.Send("MULTI")
	conn.Send("HSET", key, field, value)
	conn.Send("EXPIRE", key, 60*30)
	_, err := redis.Ints(conn.Do("EXEC"))
	if err != nil {
		fmt.Printf("ERROR: unable to set redis key %s\n", err.Error())
	}
}

func (a *Alerter) Resolve(key, field string) {
	conn := a.pool.Get()
	defer conn.Close()
	redis.Int(conn.Do("HDEL", key, field))
}

func (a *Alerter) Get(key string) (interface{}, error) {
	conn := a.pool.Get()
	defer conn.Close()
	return conn.Do("GET", key)
}

func (a *Alerter) GetAll(key_search string) (map[string][]map[string]interface{}, error) {
	var alerts = make(map[string][]map[string]interface{})
	conn := a.pool.Get()
	defer conn.Close()
	keys, _ := redis.Strings(conn.Do("KEYS", key_search))
	for _, v := range keys {
		capsule_id := strings.Split(v, ":")[1]
		current, _ := redis.StringMap(conn.Do("HGETALL", v))
		for _, alert := range current {
			var dat map[string]interface{}
			if err := json.Unmarshal([]byte(alert), &dat); err != nil {
				return alerts, err
			}
			alerts[capsule_id] = append(alerts[capsule_id], dat)
		}
	}
	return alerts, nil
}
