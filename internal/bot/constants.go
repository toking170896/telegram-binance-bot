package bot

//messages
const (
	EmptyUsernameErr      = "Your username is not set up in Telegram, in order to continue please setup a username after that type /start"
	AcceptTermsMsg        = "In order to activate you agree to our Terms of Service and Disclaimer which can be found here:\n\nTOS: www.domain.com/tos\nDisclaimer: www.domain.com/disclaimer\n"
	DeniedTermsErr        = "You have to accept our tos and disclaimer to continue."
	InsertLicenseKeyMsg   = "Please insert a license key:"
	ValidLicenseKeyMsg    = "Perfect! Your license key is valid!\n\n Fees structure test message"
	AcceptFeesRetryMsg    = "Please *Accept/Deny* our rules.\n\n We have a fixed fee of 20% for all users. You will recive a message once a week with the open fee amount which can be paid with binance pay."
	DeniedFeesMsg         = "It’s required to get started to accept out TOS and Disclaimer which can be found here:\n\nTOS: www.domain.com/tos\nDisclaimer: www.domain.com/disclaimer\n"
	InvalidLicenseKeyMsg  = "Unfortunately your license key is invalid, please check it and try once again."
	InsertBinanceKeyMsg   = "Please insert your Binance api key and api secret in such a format: *key_secret*"
	InvalidBinanceKeysMsg = "Unfortunately your Binance api keys are invalid, please check if the format is correct and try once again."
	FeeLineMsgStructure   = "TRADE: %s\n\nCLOSED DATE: %s\n\nPROFIT: %.4f\n\nFEE: %.4f\n-----------------------------------\n\n"
	ReportStartMsg        = "FEE OVERVIEW - %s - %s\n\n"
	ReportEndMsg          = "\nSummary:\n\nBased on %d amount of trades your profit: $%.2f, open fees: $%.2f\n\nPlease transfer the open fee within the next 72 hours." +
		" Otherwise your account will be banned.\nThe payment link will be valid for 1 hour. In case you missed the timeframe you can create a new payment link."
	NewlyGeneratedPaymentLinkMsg = "New payment link was generated. Please transfer the open fee. Otherwise your account will be banned."
	PaymentReminderMsg           = "We have not received your payment. If we won't receive your payment within the next 12 hours your account will be banned. \n\n"
)

var (
	ValidBinanceKeysMsg = "You Binance keys are valid!\n\n You can now join our telegram channel:\n %s"
)
