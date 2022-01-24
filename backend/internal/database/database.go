package database

import (
	"context"
	"github.com/go-redis/redis/v8"
	"github.com/pkg/errors"
	"log"
)

type Database struct {
	address    string
	password   string
	connection *redis.Client
	ctx        context.Context
}

func NewDBConnection(address, password string) Database {
	return Database{
		address:  address,
		password: password,
		ctx:      context.Background(),
	}
}

func (d *Database) Connect() error {
	if d.connection != nil {
		log.Print("Connection to database is already initialized")
		return nil
	}

	d.connection = redis.NewClient(&redis.Options{
		Addr:     d.address,
		Password: d.password,
		DB:       0,
	})

	if err := d.connection.Ping(d.ctx).Err(); err != nil {
		return err
	}
	return nil
}

func (d Database) Set(key string, value string) error {
	if d.connection == nil {
		return errors.New("no connection to database")
	}
	_, err := d.connection.Set(d.ctx, key, value, 0).Result()
	return err
}

func (d Database) Get(key string) (interface{}, error) {
	if d.connection == nil {
		return nil, errors.New("no connection to database")
	}

	val, err := d.connection.Get(d.ctx, key).Result()

	switch {
	case err == redis.Nil:
		err = errors.New("key does not exist")
	}
	return val, err
}

func (d Database) Keys(pattern string) ([]string, error) {
	if d.connection == nil {
		return nil, errors.New("no connection to database")
	}

	keys, err := d.connection.Keys(d.ctx, pattern).Result()
	if err != nil {
		return nil, err
	}
	return keys, nil
}
