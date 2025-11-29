package constants

const (
	// Environment variables names constants

	// EnvAeroSpaceSock is the environment variable for the AeroSpace socket path
	//  Default: `/tmp/bobko.aerospace-$USER.sock`
	EnvAeroSpaceSock string = "AEROSPACESOCK"

	// Other constants

	// AerspaceSocketClientVersion is the minimum version of the AeroSpace socket client
	//
	// Minimum version of the AeroSpace socket client required for compatibility
	// AeroSpace 0.15.0 till 0.19.x use <v0.2.1
	// AeroSpace 0.20.0 onwards use >v0.3.0
	AeroSpaceSocketClientMajor int = 0
	AeroSpaceSocketClientMinor int = 20
)
