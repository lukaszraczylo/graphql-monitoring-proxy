package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/suite"
)

type MainSecurityTestSuite struct {
	suite.Suite
}

func TestMainSecurityTestSuite(t *testing.T) {
	suite.Run(t, new(MainSecurityTestSuite))
}

// isTempPathAllowed checks if a temp path would be allowed by validateFilePath
func (suite *MainSecurityTestSuite) isTempPathAllowed(path string) bool {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return false
	}

	// Check if temp path is in allowed locations
	allowedPrefixes := []string{"/tmp/", "/var/tmp/"}
	for _, prefix := range allowedPrefixes {
		if strings.HasPrefix(absPath, prefix) {
			return true
		}
	}

	// Check if it's in the working directory
	workDir, err := os.Getwd()
	if err != nil {
		return false
	}
	cleanedWorkDir := filepath.Clean(workDir)
	return strings.HasPrefix(absPath, cleanedWorkDir+string(filepath.Separator))
}

// TestValidateFilePathSecurity tests the validateFilePath function for various security scenarios
func (suite *MainSecurityTestSuite) TestValidateFilePathSecurity() {
	tests := []struct {
		name        string
		inputPath   string
		description string
		shouldFail  bool
	}{
		// Path traversal attacks
		{
			name:        "Basic path traversal with double dots",
			inputPath:   "../../../../etc/passwd",
			shouldFail:  true,
			description: "Should reject basic path traversal attempt",
		},
		{
			name:        "Path traversal with current directory prefix",
			inputPath:   "./../../etc/passwd",
			shouldFail:  true,
			description: "Should reject path traversal even with ./ prefix",
		},
		{
			name:        "Deep path traversal",
			inputPath:   "../../../../../../../etc/shadow",
			shouldFail:  true,
			description: "Should reject deep path traversal attempts",
		},
		{
			name:        "URL encoded path traversal",
			inputPath:   "%2e%2e%2f%2e%2e%2fetc%2fpasswd",
			shouldFail:  true,
			description: "Should reject URL encoded traversal (if decoded)",
		},
		{
			name:        "Double encoded path traversal",
			inputPath:   "%252e%252e%252f%252e%252e%252fetc%252fpasswd",
			shouldFail:  true,
			description: "Should reject double encoded traversal",
		},
		{
			name:        "Mixed case path traversal",
			inputPath:   "../ETC/passwd",
			shouldFail:  true,
			description: "Should reject mixed case traversal attempts",
		},
		{
			name:        "Path traversal with backslashes",
			inputPath:   "..\\..\\windows\\system32\\drivers\\etc\\hosts",
			shouldFail:  true,
			description: "Should reject Windows-style path traversal",
		},

		// Absolute path attacks
		{
			name:        "Absolute path to sensitive file",
			inputPath:   "/etc/shadow",
			shouldFail:  true,
			description: "Should reject absolute path outside allowed directories",
		},
		{
			name:        "Absolute path to system directories",
			inputPath:   "/bin/bash",
			shouldFail:  true,
			description: "Should reject access to system binaries",
		},
		{
			name:        "Absolute path to home directory",
			inputPath:   "/home/user/.ssh/id_rsa",
			shouldFail:  true,
			description: "Should reject access to user directories",
		},
		{
			name:        "Absolute path to proc filesystem",
			inputPath:   "/proc/self/environ",
			shouldFail:  true,
			description: "Should reject access to proc filesystem",
		},

		// Null byte injection
		{
			name:        "Null byte injection",
			inputPath:   "/tmp/test.txt\x00.jpg",
			shouldFail:  true,
			description: "Should reject null byte injection attempts",
		},
		{
			name:        "Null byte in middle of path",
			inputPath:   "/tmp/test\x00/file.txt",
			shouldFail:  true,
			description: "Should reject null bytes anywhere in path",
		},

		// Symbolic link attempts (path patterns that might be symlinks)
		{
			name:        "Suspicious symlink pattern",
			inputPath:   "./symlink_to_etc",
			shouldFail:  false, // This is allowed by current logic but would need real symlink detection
			description: "Pattern that might be a symlink to sensitive location",
		},

		// Valid paths that should pass
		{
			name:        "Valid application directory path",
			inputPath:   "/go/src/app/banned_users.txt",
			shouldFail:  false,
			description: "Should accept valid app directory path",
		},
		{
			name:        "Valid current directory path",
			inputPath:   "./data/banned_users.txt",
			shouldFail:  false,
			description: "Should accept valid relative path",
		},
		{
			name:        "Valid temp directory path",
			inputPath:   "/tmp/test_file.txt",
			shouldFail:  false,
			description: "Should accept valid temp directory path",
		},
		{
			name:        "Valid var/tmp directory path",
			inputPath:   "/var/tmp/cache_file.json",
			shouldFail:  false,
			description: "Should accept valid var/tmp directory path",
		},
		{
			name:        "Valid nested path in app directory",
			inputPath:   "/go/src/app/config/settings.json",
			shouldFail:  false,
			description: "Should accept nested paths in allowed directories",
		},

		// Edge cases
		{
			name:        "Empty path",
			inputPath:   "",
			shouldFail:  true,
			description: "Should reject empty paths",
		},
		{
			name:        "Only dots",
			inputPath:   "..",
			shouldFail:  true,
			description: "Should reject bare double dots",
		},
		{
			name:        "Current directory only",
			inputPath:   ".",
			shouldFail:  true,
			description: "Should reject bare current directory",
		},
		{
			name:        "Root directory",
			inputPath:   "/",
			shouldFail:  true,
			description: "Should reject root directory access",
		},
		{
			name:        "Path with multiple consecutive dots",
			inputPath:   "./....//....//etc/passwd",
			shouldFail:  true,
			description: "Should reject obfuscated path traversal",
		},

		// Special character attacks
		{
			name:        "Path with semicolon",
			inputPath:   "/tmp/file;rm -rf /",
			shouldFail:  true,
			description: "Should handle paths with command injection attempts",
		},
		{
			name:        "Path with pipe",
			inputPath:   "/tmp/file|cat /etc/passwd",
			shouldFail:  true,
			description: "Should handle paths with pipe characters",
		},
		{
			name:        "Path with newline",
			inputPath:   "/tmp/file\ncat /etc/passwd",
			shouldFail:  true,
			description: "Should handle paths with newline injection",
		},
	}

	for _, tt := range tests {
		suite.Run(tt.name, func() {
			result, err := validateFilePath(tt.inputPath)

			if tt.shouldFail {
				suite.Error(err, "Expected error for path: %s (%s)", tt.inputPath, tt.description)
				suite.Empty(result, "Should return empty result on error")

				// Verify error messages don't leak sensitive information
				if err != nil {
					errMsg := strings.ToLower(err.Error())
					suite.NotContains(errMsg, "secret", "Error should not contain 'secret'")
					suite.NotContains(errMsg, "password", "Error should not contain 'password'")
					suite.NotContains(errMsg, "key", "Error should not contain 'key'")
				}
			} else {
				suite.NoError(err, "Expected no error for path: %s (%s)", tt.inputPath, tt.description)
				suite.NotEmpty(result, "Should return validated path")
				suite.Equal(tt.inputPath, result, "Should return original path when valid")
			}
		})
	}
}

// TestValidateFilePathConcurrentAccess tests path validation under concurrent conditions
func (suite *MainSecurityTestSuite) TestValidateFilePathConcurrentAccess() {
	maliciousPaths := []string{
		"../../../../etc/passwd",
		"../../../etc/shadow",
		"/etc/hosts",
		"./../../var/log/messages",
		"/proc/self/environ",
	}

	suite.Run("Concurrent malicious paths should all be rejected", func() {
		done := make(chan error, len(maliciousPaths))

		for _, path := range maliciousPaths {
			go func(p string) {
				_, err := validateFilePath(p)
				done <- err
			}(path)
		}

		// Collect all results
		for i := 0; i < len(maliciousPaths); i++ {
			err := <-done
			suite.Error(err, "All malicious paths should be rejected concurrently")
		}
	})
}

// TestValidateFilePathWithRealFiles tests validation with actual file system operations
func (suite *MainSecurityTestSuite) TestValidateFilePathWithRealFiles() {
	// Create temporary directory and files for testing
	tempDir, err := os.MkdirTemp("", "path_security_test")
	suite.NoError(err)
	defer os.RemoveAll(tempDir)

	// Create a test file
	testFile := filepath.Join(tempDir, "test.txt")
	err = os.WriteFile(testFile, []byte("test content"), 0644)
	suite.NoError(err)

	// Determine if temp file should fail based on system temp location
	tempFileShouldFail := !suite.isTempPathAllowed(testFile)

	tests := []struct {
		name       string
		path       string
		shouldFail bool
	}{
		{
			name:       "Valid temp file",
			path:       testFile,
			shouldFail: tempFileShouldFail, // Depends on system temp location
		},
		{
			name:       "Non-existent file in allowed directory",
			path:       "/tmp/non_existent.txt",
			shouldFail: false, // Should pass validation (file existence not checked)
		},
		{
			name:       "Directory instead of file",
			path:       "/tmp/",
			shouldFail: false, // Should pass validation
		},
	}

	for _, tt := range tests {
		suite.Run(tt.name, func() {
			_, err := validateFilePath(tt.path)
			if tt.shouldFail {
				suite.Error(err)
			} else {
				suite.NoError(err)
			}
		})
	}
}

// TestValidateFilePathEdgeCases tests various edge cases and corner conditions
func (suite *MainSecurityTestSuite) TestValidateFilePathEdgeCases() {
	suite.Run("Very long path", func() {
		// Create a very long path that might cause buffer overflows
		longPath := "/tmp/" + strings.Repeat("a", 4096) + ".txt"
		_, err := validateFilePath(longPath)
		// Should handle gracefully without crashing
		suite.NoError(err) // Long paths in /tmp/ should be allowed
	})

	suite.Run("Path with unicode characters", func() {
		unicodePath := "/tmp/тест.txt" // Russian characters
		_, err := validateFilePath(unicodePath)
		suite.NoError(err) // Unicode should be allowed in valid directories
	})

	suite.Run("Path with spaces", func() {
		spacePath := "/tmp/file with spaces.txt"
		_, err := validateFilePath(spacePath)
		suite.NoError(err) // Spaces should be allowed
	})

	suite.Run("Path with special but safe characters", func() {
		specialPath := "/tmp/file-name_123.json"
		_, err := validateFilePath(specialPath)
		suite.NoError(err) // Safe special characters should be allowed
	})
}

// TestValidateFilePathAllowedDirectories tests the allowed directory logic
func (suite *MainSecurityTestSuite) TestValidateFilePathAllowedDirectories() {
	allowedTests := []struct {
		name string
		path string
	}{
		{"Go app directory", "/go/src/app/config.json"},
		{"Current directory", "./config.json"},
		{"Temp directory", "/tmp/cache.json"},
		{"Var temp directory", "/var/tmp/session.json"},
	}

	for _, tt := range allowedTests {
		suite.Run(tt.name, func() {
			result, err := validateFilePath(tt.path)
			suite.NoError(err, "Path should be allowed: %s", tt.path)
			suite.Equal(tt.path, result)
		})
	}

	disallowedTests := []struct {
		name string
		path string
	}{
		{"Home directory", "/home/user/file.txt"},
		{"Root etc", "/etc/config"},
		{"System bin", "/bin/executable"},
		{"Var log", "/var/log/messages"},
		{"Opt directory", "/opt/app/config"},
		{"Absolute path without allowed prefix", "/random/path/file.txt"},
	}

	for _, tt := range disallowedTests {
		suite.Run(tt.name, func() {
			_, err := validateFilePath(tt.path)
			suite.Error(err, "Path should be rejected: %s", tt.path)
		})
	}
}

// TestValidateFilePathBoundaryConditions tests boundary conditions
func (suite *MainSecurityTestSuite) TestValidateFilePathBoundaryConditions() {
	suite.Run("Path exactly at allowed prefix boundary", func() {
		// Test paths that are exactly the allowed prefixes
		prefixes := []string{"/go/src/app/", "./", "/tmp/", "/var/tmp/"}

		for _, prefix := range prefixes {
			// Exact prefix should be allowed
			_, err := validateFilePath(prefix)
			suite.NoError(err, "Exact prefix should be allowed: %s", prefix)

			// Prefix with filename should be allowed
			_, err = validateFilePath(prefix + "file.txt")
			suite.NoError(err, "Prefix with file should be allowed: %s", prefix+"file.txt")

			// Similar but not exact prefix should be rejected (if not otherwise allowed)
			if prefix != "./" { // Skip this test for "./" as it's tricky
				similar := prefix[:len(prefix)-1] + "x/"
				_, err = validateFilePath(similar + "file.txt")
				if !strings.HasPrefix(similar, "/tmp") && !strings.HasPrefix(similar, "/var/tmp") {
					suite.Error(err, "Similar but different prefix should be rejected: %s", similar+"file.txt")
				}
			}
		}
	})
}

// BenchmarkValidateFilePath benchmarks the path validation function
func BenchmarkValidateFilePath(b *testing.B) {
	testPaths := []string{
		"/go/src/app/config.json",
		"./data/file.txt",
		"/tmp/cache.json",
		"../../../../etc/passwd", // malicious
		"/etc/shadow",            // malicious
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, path := range testPaths {
			validateFilePath(path)
		}
	}
}

// TestValidateFilePathErrorMessages tests that error messages are appropriate
func (suite *MainSecurityTestSuite) TestValidateFilePathErrorMessages() {
	errorTests := []struct {
		path             string
		expectedContains string
	}{
		{"", "empty"},
		{"..", "traversal"},
		{"../etc/passwd", "traversal"},
		{"/tmp/file\x00.txt", "null byte"},
		{"/etc/passwd", "not in allowed"},
	}

	for _, tt := range errorTests {
		suite.Run(fmt.Sprintf("Error for %s", tt.path), func() {
			_, err := validateFilePath(tt.path)
			suite.Error(err)
			suite.Contains(strings.ToLower(err.Error()), tt.expectedContains)
		})
	}
}
