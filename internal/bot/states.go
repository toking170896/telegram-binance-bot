package bot

import "strconv"

//states
const (
	AcceptTermsState = "acceptTerms"
	DenyTermsState = "denyTerms"
	AcceptFeesState = "acceptFees"
	DenyFeesState = "denyFees"
	InsertLicenseKeyState = "insertLicenseKey"
	ValidLicenseKey = "validLicenseKey"
	InvalidLicenseKey = "invalidLicenseKey"
	InsertBinanceKeysState = "insertBinanceKeys"
	ValidBinanceKeys = "validBinanceKeys"
	InvalidBinanceKeys = "invalidBinanceKeys"
	GeneratePaymentLinkState = "generatePaymentLink"
)

func (s *Svc) updateUserState(userID int64, state string) {
	s.States.Set(strconv.Itoa(int(userID)), state)
}

func (s *Svc) getUserState(userID int64) string {
	val, exists := s.States.Get(strconv.Itoa(int(userID)))
	if !exists {
		return ""
	}
	return val.(string)
}