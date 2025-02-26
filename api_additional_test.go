package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	libpack_logger "github.com/lukaszraczylo/graphql-monitoring-proxy/logging"
)

func (suite *Tests) Test_PeriodicallyReloadBannedUsers() {
	// Setup
	cfg = &config{}
	parseConfig()
	cfg.Logger = libpack_logger.New()
	cfg.Api.BannedUsersFile = filepath.Join(os.TempDir(), "banned_users_reload_test.json")

	// Initial empty banned users
	bannedUsersIDsMutex.Lock()
	bannedUsersIDs = make(map[string]string)
	bannedUsersIDsMutex.Unlock()

	// Create a test version of periodicallyReloadBannedUsers that executes once and signals completion
	done := make(chan bool)
	testPeriodicallyReloadBannedUsers := func() {
		// Just call loadBannedUsers once
		loadBannedUsers()
		done <- true
	}

	// Run the test with initial empty banned users file
	suite.Run("reload with empty file", func() {
		// Clear existing file if any
		os.Remove(cfg.Api.BannedUsersFile)
		os.Remove(fmt.Sprintf("%s.lock", cfg.Api.BannedUsersFile))

		// Ensure banned users map is empty
		bannedUsersIDsMutex.Lock()
		bannedUsersIDs = make(map[string]string)
		bannedUsersIDsMutex.Unlock()

		// Execute reloader once
		go testPeriodicallyReloadBannedUsers()
		<-done

		// Verify file was created
		_, err := os.Stat(cfg.Api.BannedUsersFile)
		assert.NoError(err)

		// Safely check the map
		bannedUsersIDsMutex.RLock()
		mapSize := len(bannedUsersIDs)
		bannedUsersIDsMutex.RUnlock()

		// Verify map is still empty
		assert.Equal(0, mapSize)
	})

	// Run the test with a populated banned users file
	suite.Run("reload with populated file", func() {
		// Create file with test data
		testData := map[string]string{
			"test-user-reload-1": "reason reload 1",
			"test-user-reload-2": "reason reload 2",
		}
		data, _ := json.Marshal(testData)
		err := os.WriteFile(cfg.Api.BannedUsersFile, data, 0644)
		assert.NoError(err)

		// Clear the banned users map
		bannedUsersIDsMutex.Lock()
		bannedUsersIDs = make(map[string]string)
		bannedUsersIDsMutex.Unlock()

		// Execute reloader once
		go testPeriodicallyReloadBannedUsers()
		<-done

		// Safely check the map
		bannedUsersIDsMutex.RLock()
		mapSize := len(bannedUsersIDs)
		value1 := bannedUsersIDs["test-user-reload-1"]
		value2 := bannedUsersIDs["test-user-reload-2"]
		bannedUsersIDsMutex.RUnlock()

		// Verify banned users map was loaded
		assert.Equal(2, mapSize)
		assert.Equal("reason reload 1", value1)
		assert.Equal("reason reload 2", value2)
	})

	// Test updating banned users file while reloader is running
	suite.Run("reload with updated file", func() {
		// Start with initial data
		initialData := map[string]string{
			"test-user-initial": "initial reason",
		}
		data, _ := json.Marshal(initialData)
		err := os.WriteFile(cfg.Api.BannedUsersFile, data, 0644)
		assert.NoError(err)

		// Clear the banned users map
		bannedUsersIDsMutex.Lock()
		bannedUsersIDs = make(map[string]string)
		bannedUsersIDsMutex.Unlock()

		// Execute reloader once to load initial data
		go testPeriodicallyReloadBannedUsers()
		<-done

		// Safely check the map
		bannedUsersIDsMutex.RLock()
		mapSize := len(bannedUsersIDs)
		initialValue := bannedUsersIDs["test-user-initial"]
		bannedUsersIDsMutex.RUnlock()

		// Verify initial data was loaded
		assert.Equal(1, mapSize)
		assert.Equal("initial reason", initialValue)

		// Update the file with new data
		updatedData := map[string]string{
			"test-user-updated-1": "updated reason 1",
			"test-user-updated-2": "updated reason 2",
		}
		data, _ = json.Marshal(updatedData)
		err = os.WriteFile(cfg.Api.BannedUsersFile, data, 0644)
		assert.NoError(err)

		// Execute reloader again to load updated data
		go testPeriodicallyReloadBannedUsers()
		<-done

		// Safely check the map
		bannedUsersIDsMutex.RLock()
		mapSize = len(bannedUsersIDs)
		value1 := bannedUsersIDs["test-user-updated-1"]
		value2 := bannedUsersIDs["test-user-updated-2"]
		_, exists := bannedUsersIDs["test-user-initial"]
		bannedUsersIDsMutex.RUnlock()

		// Verify updated data was loaded
		assert.Equal(2, mapSize)
		assert.Equal("updated reason 1", value1)
		assert.Equal("updated reason 2", value2)
		assert.False(exists)
	})

	// Cleanup
	os.Remove(cfg.Api.BannedUsersFile)
	os.Remove(fmt.Sprintf("%s.lock", cfg.Api.BannedUsersFile))
}

// This is a better approach instead of the ticker-based test
func (suite *Tests) Test_LoadUnloadBannedUsers() {
	// Setup
	cfg = &config{}
	parseConfig()
	cfg.Logger = libpack_logger.New()
	cfg.Api.BannedUsersFile = filepath.Join(os.TempDir(), "banned_users_update_test.json")

	// Create a test banned users file with initial content
	initialData := map[string]string{
		"user1": "reason1",
		"user2": "reason2",
	}
	data, _ := json.Marshal(initialData)
	err := os.WriteFile(cfg.Api.BannedUsersFile, data, 0644)
	assert.NoError(err)
	defer os.Remove(cfg.Api.BannedUsersFile)
	defer os.Remove(fmt.Sprintf("%s.lock", cfg.Api.BannedUsersFile))

	// Test loading banned users
	suite.Run("load banned users", func() {
		// Clear the banned users map
		bannedUsersIDsMutex.Lock()
		bannedUsersIDs = make(map[string]string)
		bannedUsersIDsMutex.Unlock()

		// Load banned users
		loadBannedUsers()

		// Check the banned users map
		bannedUsersIDsMutex.RLock()
		count := len(bannedUsersIDs)
		reason1 := bannedUsersIDs["user1"]
		reason2 := bannedUsersIDs["user2"]
		bannedUsersIDsMutex.RUnlock()

		assert.Equal(2, count)
		assert.Equal("reason1", reason1)
		assert.Equal("reason2", reason2)
	})

	// Test updating banned users
	suite.Run("update banned users", func() {
		// Update the banned users map
		bannedUsersIDsMutex.Lock()
		bannedUsersIDs = map[string]string{
			"user3": "reason3",
			"user4": "reason4",
		}
		bannedUsersIDsMutex.Unlock()

		// Store the updated banned users
		err := storeBannedUsers()
		assert.NoError(err)

		// Clear the banned users map
		bannedUsersIDsMutex.Lock()
		bannedUsersIDs = make(map[string]string)
		bannedUsersIDsMutex.Unlock()

		// Load banned users again
		loadBannedUsers()

		// Check the banned users map
		bannedUsersIDsMutex.RLock()
		count := len(bannedUsersIDs)
		reason3 := bannedUsersIDs["user3"]
		reason4 := bannedUsersIDs["user4"]
		_, user1Exists := bannedUsersIDs["user1"]
		bannedUsersIDsMutex.RUnlock()

		assert.Equal(2, count)
		assert.Equal("reason3", reason3)
		assert.Equal("reason4", reason4)
		assert.False(user1Exists)
	})
}
