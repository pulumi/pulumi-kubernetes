package await

// canonicalizeDeploymentAPIVersion unifies the various pre-release apiVerion values for a
// Deployment into "apps/v1".
func canonicalizeDeploymentAPIVersion(ver string) string {
	switch ver {
	case "extensions/v1beta1", "apps/v1beta1", "apps/v1beta2", "apps/v1":
		// Canonicalize all of these to "apps/v1".
		return "apps/v1"
	default:
		// If the input version was not a version we understand, just return it as-is.
		return ver
	}
}

// canonicalizeStatefulSetAPIVersion unifies the various pre-release apiVerion values for a
// StatefulSet into "apps/v1".
func canonicalizeStatefulSetAPIVersion(ver string) string {
	switch ver {
	case "apps/v1beta1", "apps/v1beta2", "apps/v1":
		// Canonicalize all of these to "apps/v1".
		return "apps/v1"
	default:
		// If the input version was not a version we understand, just return it as-is.
		return ver
	}
}

