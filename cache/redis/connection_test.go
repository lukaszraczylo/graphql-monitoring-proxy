package libpack_redis

import (
	"testing"
	"time"

	"github.com/gookit/goutil/envutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type RedisConfigSuite struct {
	suite.Suite
	redisConfig *RedisConfig
}

func (suite *RedisConfigSuite) SetupTest() {
	redis_server := envutil.Getenv("REDIS_SERVER", "localhost:6379")
	suite.redisConfig = NewClient(&RedisClientConfig{
		RedisServer:   redis_server,
		RedisPassword: "",
		RedisDB:       0,
	})
	suite.redisConfig.Delete("testkey")
}

func TestRedisConfigSuite(t *testing.T) {
	suite.Run(t, new(RedisConfigSuite))
}

func (suite *RedisConfigSuite) TestSet() {
	key := "testkey"
	value := []byte("testvalue")
	suite.redisConfig.Delete(key) // Ensure the key is deleted before the test

	// Test writing a new key-value pair
	suite.redisConfig.Set(key, value, 0)
	storedValue, found := suite.redisConfig.Get(key)
	assert.True(suite.T(), found)
	assert.Equal(suite.T(), value, storedValue)

	// Test overwriting an existing key-value pair
	newValue := []byte("newvalue")
	suite.redisConfig.Set(key, newValue, 0)
	storedValue, found = suite.redisConfig.Get(key)
	assert.True(suite.T(), found)
	assert.Equal(suite.T(), newValue, storedValue)

	suite.redisConfig.Delete(key) // Clean up after the test
}

func (suite *RedisConfigSuite) TestSetWithExpiry() {
	key := "testkey"
	value := []byte("testvalue")
	expiry := 1 * time.Second
	suite.redisConfig.Delete(key) // Ensure the key is deleted before the test

	// Test writing a new key-value pair
	suite.redisConfig.Set(key, value, expiry)
	storedValue, found := suite.redisConfig.Get(key)
	assert.True(suite.T(), found)
	assert.Equal(suite.T(), value, storedValue)

	// Test that key expires after the specified time
	time.Sleep(2 * time.Second)
	_, found = suite.redisConfig.Get(key)
	assert.False(suite.T(), found)

	suite.redisConfig.Delete(key) // Clean up after the test
}

func (suite *RedisConfigSuite) TestGet() {
	key := "testkey"
	value := []byte("testvalue")
	suite.redisConfig.Set(key, value, 0) // Set the key-value pair
	storedValue, found := suite.redisConfig.Get(key)
	assert.True(suite.T(), found)
	assert.Equal(suite.T(), value, storedValue)
}

func (suite *RedisConfigSuite) TestDeleteKey() {
	key := "testkey"
	value := []byte("testvalue")
	suite.redisConfig.Set(key, value, 0) // Set the key-value pair
	suite.redisConfig.Delete(key)
	_, found := suite.redisConfig.Get(key)
	assert.False(suite.T(), found)
}

func (suite *RedisConfigSuite) TestCheckIfKeyExists() {
	ttl := time.Duration(10) * time.Second
	key := "testkey"
	value := []byte("testvalue")
	suite.redisConfig.Set(key, value, ttl) // Set the key-value pair
	_, found := suite.redisConfig.Get(key)
	assert.True(suite.T(), found)

	suite.redisConfig.Delete(key)
	_, found = suite.redisConfig.Get(key)
	assert.False(suite.T(), found)
}

func (suite *RedisConfigSuite) TestGetKeys() {
	ttl := time.Duration(10) * time.Second
	suite.redisConfig.Set("testkey1", []byte("testvalue1"), ttl)
	suite.redisConfig.Set("testkey2", []byte("testvalue2"), ttl)
	suite.redisConfig.Set("otherkey", []byte("othervalue"), ttl)

	keys, _ := suite.redisConfig.client.Keys(suite.redisConfig.ctx, prependKeyName("testkey*")).Result()
	expectedKeys := []string{prependKeyName("testkey1"), prependKeyName("testkey2")}
	assert.ElementsMatch(suite.T(), expectedKeys, keys)

	suite.redisConfig.client.Del(suite.redisConfig.ctx, "testkey1", "testkey2", "otherkey")
}

func (suite *RedisConfigSuite) TestGetKeysCount() {
	ttl := time.Duration(10) * time.Second
	suite.redisConfig.Set("testkey1", []byte("testvalue1"), ttl)
	suite.redisConfig.Set("testkey2", []byte("testvalue2"), ttl)
	suite.redisConfig.Set("otherkey", []byte("othervalue"), ttl)

	assert.Equal(suite.T(), 2, suite.redisConfig.CountQueriesWithPattern("testkey*"))

	suite.redisConfig.client.Del(suite.redisConfig.ctx, "testkey1", "testkey2", "otherkey")
}
