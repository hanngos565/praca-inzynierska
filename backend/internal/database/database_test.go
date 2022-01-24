package database

import (
	"context"
	"errors"
	"github.com/go-redis/redismock/v8"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDatabase_Connect(t *testing.T) {
	t.Run("should return nil when there is already a connection to database", func(t *testing.T) {
		//given
		client, _ := redismock.NewClientMock()
		d := Database{address: "localhost:8081", password: "", connection: client, ctx: context.Background()}

		//when
		err := d.Connect()

		//then
		assert.NoError(t, err)
	})
	t.Run("should return error when couldn't ping database", func(t *testing.T) {

		//given
		d := Database{address: "localhost:8081", password: "", ctx: context.Background()}

		//when
		err := d.Connect()

		//then
		assert.Error(t, err)
		assert.NotEqual(t, nil, d.connection)
	})
}

func TestDatabase_Set(t *testing.T) {
	const key, value = "key", "value"
	t.Run("should return error when there is no connection to database", func(t *testing.T) {
		//given
		_, clientMock := redismock.NewClientMock()
		clientMock.ClearExpect()
		clientMock.MatchExpectationsInOrder(false)
		database := Database{connection: nil}

		//when
		err := database.Set(key, value)

		//then
		assert.Error(t, err, errors.New("no connection to database"))
		if err := clientMock.ExpectationsWereMet(); err != nil {
			t.Error(err)
		}
	})
	t.Run("should return error when couldn't set value in database", func(t *testing.T) {
		//given
		client, clientMock := redismock.NewClientMock()
		clientMock.ClearExpect()
		clientMock.MatchExpectationsInOrder(true)
		database := Database{connection: client}
		clientMock.ExpectSet(key, value, 0)

		//when
		err := database.Set(key, value)

		//then
		assert.Error(t, err)
		if err := clientMock.ExpectationsWereMet(); err != nil {
			t.Error(err)
		}
	})
	t.Run("should return no error when insert data successfully to database", func(t *testing.T) {

		//given
		client, clientMock := redismock.NewClientMock()
		clientMock.ClearExpect()
		clientMock.MatchExpectationsInOrder(true)
		database := Database{connection: client}
		clientMock.ExpectSet(key, value, 0).SetVal(value)

		//when
		err := database.Set(key, value)

		//then
		assert.NoError(t, err)
		if err := clientMock.ExpectationsWereMet(); err != nil {
			t.Error(err)
		}
	})
}

func TestDatabase_Get(t *testing.T) {
	const key = "key"
	t.Run("should return error when there is no connection to database", func(t *testing.T) {
		//given
		_, clientMock := redismock.NewClientMock()
		clientMock.ClearExpect()
		clientMock.MatchExpectationsInOrder(false)
		database := Database{connection: nil}

		//when
		_, err := database.Get(key)

		//then
		assert.Error(t, err, errors.New("no connection to database"))
		if err := clientMock.ExpectationsWereMet(); err != nil {
			t.Error(err)
		}
	})
	t.Run("should return error when key does not exist in database", func(t *testing.T) {
		//given
		client, clientMock := redismock.NewClientMock()
		clientMock.ClearExpect()
		clientMock.MatchExpectationsInOrder(true)
		database := Database{connection: client}
		clientMock.ExpectGet(key).RedisNil()

		//when
		_, err := database.Get(key)

		//then
		assert.Equal(t, "key does not exist", err.Error())
		if err := clientMock.ExpectationsWereMet(); err != nil {
			t.Error(err)
		}
	})
	t.Run("should return error when couldn't get data from database", func(t *testing.T) {

		//given
		client, clientMock := redismock.NewClientMock()
		clientMock.ClearExpect()
		clientMock.MatchExpectationsInOrder(true)
		database := Database{connection: client}
		clientMock.ExpectGet(key)

		//when
		_, err := database.Get(key)

		//then
		assert.Error(t, err)
		if err := clientMock.ExpectationsWereMet(); err != nil {
			t.Error(err)
		}
	})
	t.Run("should return no error when successfully got data from database", func(t *testing.T) {
		//given
		client, clientMock := redismock.NewClientMock()
		clientMock.ClearExpect()
		clientMock.MatchExpectationsInOrder(true)
		database := Database{connection: client}
		val := "val"
		clientMock.ExpectGet(key).SetVal(val)

		//when
		_, err := database.Get(key)

		//then
		assert.NoError(t, err)
		if err := clientMock.ExpectationsWereMet(); err != nil {
			t.Error(err)
		}
	})
}

func TestDatabase_Keys(t *testing.T) {
	t.Run("should return error when there is no connection to database", func(t *testing.T) {
		//given
		_, clientMock := redismock.NewClientMock()
		clientMock.ClearExpect()
		clientMock.MatchExpectationsInOrder(false)
		database := Database{connection: nil}

		//when
		_, err := database.Keys("*")

		//then
		assert.Error(t, err, errors.New("no connection to database"))
		if err := clientMock.ExpectationsWereMet(); err != nil {
			t.Error(err)
		}
	})
	t.Run("should return error when couldn't get keys from database", func(t *testing.T) {
		//given
		client, clientMock := redismock.NewClientMock()
		clientMock.ClearExpect()
		clientMock.MatchExpectationsInOrder(false)
		database := Database{connection: client}
		clientMock.ExpectKeys("*")

		//when
		val, err := database.Keys("*")

		//then
		assert.Error(t, err)
		var n []string = nil
		assert.Equal(t, n, val)
		if err := clientMock.ExpectationsWereMet(); err != nil {
			t.Error(err)
		}
	})
	t.Run("should return no error when successfully got keys from database", func(t *testing.T) {
		//given
		client, clientMock := redismock.NewClientMock()
		clientMock.ClearExpect()
		clientMock.MatchExpectationsInOrder(false)
		database := Database{connection: client}
		clientMock.ExpectKeys("*").SetVal([]string{"models", "images", "results"})

		//when
		_, err := database.Keys("*")

		//then
		assert.NoError(t, err)
		if err := clientMock.ExpectationsWereMet(); err != nil {
			t.Error(err)
		}
	})
}
