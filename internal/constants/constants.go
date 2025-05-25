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
	AeroSpaceSocketClientVersion string = "0.15.2-Beta"
)
