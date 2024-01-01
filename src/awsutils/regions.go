package awsutils

var (
	supportedRegions = map[string]struct{}{
		"us-east-1": {},
		"us-east-2": {},
		// us-west-1 is intentionally omitted, see BUG-4233
		"us-west-2": {},
	}
)

// IsValidRegion returns if the AWS region is one that is supported by this application
func IsValidRegion(region string) bool {
	_, found := supportedRegions[region]
	return found
}
