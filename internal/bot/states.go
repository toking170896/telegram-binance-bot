package bot

//states
const (
	AcceptTermsState         = "acceptTerms"
	DenyTermsState           = "denyTerms"
	AcceptFeesState          = "acceptFees"
	DenyFeesState            = "denyFees"
	InsertLicenseKeyState    = "insertLicenseKey"
	ValidLicenseKey          = "validLicenseKey"
	InvalidLicenseKey        = "invalidLicenseKey"
	InsertBinanceKeysState   = "insertBinanceKeys"
	ValidBinanceKeys         = "validBinanceKeys"
	InvalidBinanceKeys       = "invalidBinanceKeys"
	GeneratePaymentLinkState = "generatePaymentLink"
)

func (s *Svc) updateUserState(userID string, state string) {
	s.States.Set(userID, state)
}

func (s *Svc) getUserState(userID string) string {
	val, exists := s.States.Get(userID)
	if !exists {
		return ""
	}
	return val.(string)
}
