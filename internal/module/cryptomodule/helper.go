package cryptomodule

func boolToSeverity(condition bool, weight int, riskScore *int, maxRiskScore *int) uint32 {
	*maxRiskScore += 5
	if condition {
		*riskScore += weight
		return uint32(weight) // severity is directly based on the risk weight
	}
	return 0 // no risk if not true
}

func boolToReason(condition bool, reason string) string {
	if condition {
		return reason
	}
	return ""
}
