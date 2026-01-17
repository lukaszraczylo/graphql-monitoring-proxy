// Package libpack_config provides build-time configuration variables
// for package name and version, which are set during the build process
// using ldflags.
package libpack_config

var (
	PKG_NAME    string = "not-specified"
	PKG_VERSION string = "0.0.0-dev"
)
