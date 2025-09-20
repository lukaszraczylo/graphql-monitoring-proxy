package libpack_cache_redis

import (
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type RedisConfigSuite struct {
	suite.Suite
	redisConfig  *RedisConfig
	redis_server *miniredis.Miniredis
}

func (suite *RedisConfigSuite) SetupTest() {
	suite.redis_server, _ = miniredis.Run()
	var err error
	suite.redisConfig, err = New(&RedisClientConfig{
		RedisServer:   suite.redis_server.Addr(),
		RedisPassword: "",
		RedisDB:       0,
	})
	assert.NoError(suite.T(), err)
	suite.redisConfig.Delete("testkey")
}

func TestRedisConfigSuite(t *testing.T) {
	suite.Run(t, new(RedisConfigSuite))
}

func (suite *RedisConfigSuite) TestSet() {
	key := "testkeyset"
	value := []byte("testvalue")
	suite.redisConfig.Delete(key) // Ensure the key is deleted before the test

	// Test writing a new key-value pair
	err := suite.redisConfig.Set(key, value, 0)
	assert.NoError(suite.T(), err)
	storedValue, found, err := suite.redisConfig.Get(key)
	assert.NoError(suite.T(), err)
	assert.True(suite.T(), found)
	assert.Equal(suite.T(), value, storedValue)

	// Test overwriting an existing key-value pair
	newValue := []byte("newvalue")
	err = suite.redisConfig.Set(key, newValue, 0)
	assert.NoError(suite.T(), err)
	storedValue, found, err = suite.redisConfig.Get(key)
	assert.NoError(suite.T(), err)
	assert.True(suite.T(), found)
	assert.Equal(suite.T(), newValue, storedValue)

	suite.redisConfig.Delete(key) // Clean up after the test
}

func (suite *RedisConfigSuite) TestSetWithExpiry() {
	key := "testkey_with_expiry"
	value := []byte("testvaluewithexpiry")
	expiry := 2 * time.Second
	suite.redisConfig.Delete(key) // Ensure the key is deleted before the test

	// Test writing a new key-value pair
	err := suite.redisConfig.Set(key, value, expiry)
	assert.NoError(suite.T(), err)
	storedValue, found, err := suite.redisConfig.Get(key)
	assert.NoError(suite.T(), err)
	assert.True(suite.T(), found)
	assert.Equal(suite.T(), value, storedValue)
	_, found, err = suite.redisConfig.Get(key)
	assert.NoError(suite.T(), err)
	assert.True(suite.T(), found, "Key should exist")

	// Test that key expires after the specified time
	suite.redis_server.FastForward(3 * time.Second)
	_, found, err = suite.redisConfig.Get(key)
	assert.NoError(suite.T(), err)
	assert.False(suite.T(), found, "Key should have expired after 2 seconds")

	suite.redisConfig.Delete(key) // Clean up after the test
}

func (suite *RedisConfigSuite) TestGet() {
	key := "testkeyget"
	value := []byte("testvalue")
	err := suite.redisConfig.Set(key, value, 0) // Set the key-value pair
	assert.NoError(suite.T(), err)
	storedValue, found, err := suite.redisConfig.Get(key)
	assert.NoError(suite.T(), err)
	assert.True(suite.T(), found)
	assert.Equal(suite.T(), value, storedValue)
}

func (suite *RedisConfigSuite) TestDeleteKey() {
	key := "testkeydelete"
	value := []byte("testvalue")
	err := suite.redisConfig.Set(key, value, 0) // Set the key-value pair
	assert.NoError(suite.T(), err)
	err = suite.redisConfig.Delete(key)
	assert.NoError(suite.T(), err)
	_, found, err := suite.redisConfig.Get(key)
	assert.NoError(suite.T(), err)
	assert.False(suite.T(), found)
}

func (suite *RedisConfigSuite) TestCheckIfKeyExists() {
	ttl := time.Duration(10) * time.Second
	key := "testkeyifexists"
	value := []byte("testvalue")
	err := suite.redisConfig.Set(key, value, ttl) // Set the key-value pair
	assert.NoError(suite.T(), err)
	_, found, err := suite.redisConfig.Get(key)
	assert.NoError(suite.T(), err)
	assert.True(suite.T(), found)

	err = suite.redisConfig.Delete(key)
	assert.NoError(suite.T(), err)
	_, found, err = suite.redisConfig.Get(key)
	assert.NoError(suite.T(), err)
	assert.False(suite.T(), found)
}

func (suite *RedisConfigSuite) TestGetKeys() {
	ttl := time.Duration(10) * time.Second
	err := suite.redisConfig.Set("testkey1", []byte("testvalue1"), ttl)
	assert.NoError(suite.T(), err)
	err = suite.redisConfig.Set("testkey2", []byte("testvalue2"), ttl)
	assert.NoError(suite.T(), err)
	err = suite.redisConfig.Set("otherkey", []byte("othervalue"), ttl)
	assert.NoError(suite.T(), err)

	keys, _ := suite.redisConfig.client.Keys(suite.redisConfig.ctx, "testkey*").Result()
	expectedKeys := []string{"testkey1", "testkey2"}
	assert.ElementsMatch(suite.T(), expectedKeys, keys)

	suite.redisConfig.client.Del(suite.redisConfig.ctx, "testkey1", "testkey2", "otherkey")
}

func (suite *RedisConfigSuite) TestGetKeysCount() {
	ttl := time.Duration(10) * time.Second
	suite.redisConfig.Set("testkey1", []byte("testvalue1"), ttl)
	suite.redisConfig.Set("testkey2", []byte("testvalue2"), ttl)
	suite.redisConfig.Set("otherkey", []byte("othervalue"), ttl)

	count1, err := suite.redisConfig.CountQueriesWithPattern("testkey*")
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), 2, count1)
	count2, err := suite.redisConfig.CountQueriesWithPattern("otherkey*")
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), 1, count2)
	count3, err := suite.redisConfig.CountQueries()
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), int64(3), count3)

	suite.redisConfig.client.Del(suite.redisConfig.ctx, "testkey1", "testkey2", "otherkey")
}
