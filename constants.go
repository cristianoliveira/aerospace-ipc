package aerospace

const (
  // Environment variables names constants

  // EnvAeroSpaceSock is the environment variable for the AeroSpace socket path
  //  Default: `/tmp/bobko.aerospace-$USER.sock`
  EnvAeroSpaceSock string = "AEROSPACESOCK"

  // EnvWarnVersionMismatch is the environment variable for warning about
  // socket version mismatch
  //  Default: "true" (enabled by default, set to "false" to disable warnings about version mismatch
  EnvAeroSpaceVersion string = "AEROSPACE_WARN_VERSION_MISMATCH" 

  // Other constants

  // AerspaceSocketClientVersion is the minimum version of the AeroSpace socket client
  //
  // Minimum version of the AeroSpace socket client required for compatibility
  AeroSpaceSocketClientVersion string = "0.15.2-Beta" 
)
