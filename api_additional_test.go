package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	libpack_logger "github.com/lukaszraczylo/graphql-monitoring-proxy/logging"
	"github.com/stretchr/testify/assert"
)

func (suite *Tests) Test_PeriodicallyReloadBannedUsers() {
	// Setup
	cfg = &config{}
	parseConfig()
	cfg.Logger = libpack_logger.New()
	cfg.Api.BannedUsersFile = filepath.Join(os.TempDir(), "banned_users_reload_test.json")

	// Initial empty banned users
	replaceBannedUsers(map[string]string{})

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
		_ = os.Remove(cfg.Api.BannedUsersFile)
		_ = os.Remove(fmt.Sprintf("%s.lock", cfg.Api.BannedUsersFile))

		// Ensure banned users map is empty
		replaceBannedUsers(map[string]string{})

		// Execute reloader once
		go testPeriodicallyReloadBannedUsers()
		<-done

		// Verify file was created
		_, err := os.Stat(cfg.Api.BannedUsersFile)
		assert.NoError(suite.T(), err)

		// Safely check the map
		mapSize := len(snapshotBannedUsers())

		// Verify map is still empty
		assert.Equal(suite.T(), 0, mapSize)
	})

	// Run the test with a populated banned users file
	suite.Run("reload with populated file", func() {
		// Create file with test data
		testData := map[string]string{
			"test-user-reload-1": "reason reload 1",
			"test-user-reload-2": "reason reload 2",
		}
		data, _ := json.Marshal(testData)
		err := os.WriteFile(cfg.Api.BannedUsersFile, data, 0o644)
		assert.NoError(suite.T(), err)

		// Clear the banned users map
		replaceBannedUsers(map[string]string{})

		// Execute reloader once
		go testPeriodicallyReloadBannedUsers()
		<-done

		// Safely check the map
		snap := snapshotBannedUsers()
		mapSize := len(snap)
		value1 := snap["test-user-reload-1"]
		value2 := snap["test-user-reload-2"]

		// Verify banned users map was loaded
		assert.Equal(suite.T(), 2, mapSize)
		assert.Equal(suite.T(), "reason reload 1", value1)
		assert.Equal(suite.T(), "reason reload 2", value2)
	})

	// Test updating banned users file while reloader is running
	suite.Run("reload with updated file", func() {
		// Start with initial data
		initialData := map[string]string{
			"test-user-initial": "initial reason",
		}
		data, _ := json.Marshal(initialData)
		err := os.WriteFile(cfg.Api.BannedUsersFile, data, 0o644)
		assert.NoError(suite.T(), err)

		// Clear the banned users map
		replaceBannedUsers(map[string]string{})

		// Execute reloader once to load initial data
		go testPeriodicallyReloadBannedUsers()
		<-done

		// Safely check the map
		snap := snapshotBannedUsers()
		mapSize := len(snap)
		initialValue := snap["test-user-initial"]

		// Verify initial data was loaded
		assert.Equal(suite.T(), 1, mapSize)
		assert.Equal(suite.T(), "initial reason", initialValue)

		// Update the file with new data
		updatedData := map[string]string{
			"test-user-updated-1": "updated reason 1",
			"test-user-updated-2": "updated reason 2",
		}
		data, _ = json.Marshal(updatedData)
		err = os.WriteFile(cfg.Api.BannedUsersFile, data, 0o644)
		assert.NoError(suite.T(), err)

		// Execute reloader again to load updated data
		go testPeriodicallyReloadBannedUsers()
		<-done

		// Safely check the map
		snap = snapshotBannedUsers()
		mapSize = len(snap)
		value1 := snap["test-user-updated-1"]
		value2 := snap["test-user-updated-2"]
		_, exists := snap["test-user-initial"]

		// Verify updated data was loaded
		assert.Equal(suite.T(), 2, mapSize)
		assert.Equal(suite.T(), "updated reason 1", value1)
		assert.Equal(suite.T(), "updated reason 2", value2)
		assert.False(suite.T(), exists)
	})

	// Cleanup
	_ = os.Remove(cfg.Api.BannedUsersFile)
	_ = os.Remove(fmt.Sprintf("%s.lock", cfg.Api.BannedUsersFile))
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
	err := os.WriteFile(cfg.Api.BannedUsersFile, data, 0o644)
	assert.NoError(suite.T(), err)
	defer func() { _ = os.Remove(cfg.Api.BannedUsersFile) }()
	defer func() { _ = os.Remove(fmt.Sprintf("%s.lock", cfg.Api.BannedUsersFile)) }()

	// Test loading banned users
	suite.Run("load banned users", func() {
		// Clear the banned users map
		replaceBannedUsers(map[string]string{})

		// Load banned users
		loadBannedUsers()

		// Check the banned users map
		snap := snapshotBannedUsers()
		count := len(snap)
		reason1 := snap["user1"]
		reason2 := snap["user2"]

		assert.Equal(suite.T(), 2, count)
		assert.Equal(suite.T(), "reason1", reason1)
		assert.Equal(suite.T(), "reason2", reason2)
	})

	// Test updating banned users
	suite.Run("update banned users", func() {
		// Update the banned users map
		replaceBannedUsers(map[string]string{
			"user3": "reason3",
			"user4": "reason4",
		})

		// Store the updated banned users
		err := storeBannedUsers()
		assert.NoError(suite.T(), err)

		// Clear the banned users map
		replaceBannedUsers(map[string]string{})

		// Load banned users again
		loadBannedUsers()

		// Check the banned users map
		snap := snapshotBannedUsers()
		count := len(snap)
		reason3 := snap["user3"]
		reason4 := snap["user4"]
		_, user1Exists := snap["user1"]

		assert.Equal(suite.T(), 2, count)
		assert.Equal(suite.T(), "reason3", reason3)
		assert.Equal(suite.T(), "reason4", reason4)
		assert.False(suite.T(), user1Exists)
	})
}
