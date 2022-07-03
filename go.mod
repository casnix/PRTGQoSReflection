module github.com/casnix/PRTGQoSReflection

go 1.18

require (
	github.com/BurntSushi/toml v1.1.0 // indirect
	github.com/TwiN/go-color v1.1.0 // indirect
	github.com/jessevdk/go-flags v1.5.0 // indirect
	golang.org/x/sys v0.0.0-20210320140829-1e4c9ba3b0c4 // indirect
)

// My own packaged
require github.com/casnix/PRTGQoSReflection/buildinfo v0.0.0

// Local packages
replace github.com/casnix/PRTGQoSReflection/buildinfo => ./dynamicgeneration/buildinfo
