package bot

//messages
const (
	EmptyUsernameErr = "Your username is not set up in Telegram, in order to continue please setup a username after that type /start or click [hyperlink-start]"
	AcceptTermsMsg = "In order to activate you agree to our Terms of Service and Disclaimer which can be found here:\n\nTOS: www.domain.com/tos\nDisclaimer: www.domain.com/disclaimer\n"
	DeniedTermsErr = "You have to accept our tos and disclaimer to continue."
	InsertLicenseKeyMsg = "Please insert a license key:"
	ValidLicenseKeyMsg = "Perfect! Your license key is valid!\n\n Fees structure test message"
	AcceptFeesRetryMsg = "Please *Accept/Deny* our rules.\n\n Fees structure test message"
	DeniedFeesMsg = "It’s required to get started to accept out TOS and Disclaimer which can be found here:\n\nTOS: www.domain.com/tos\nDisclaimer: www.domain.com/disclaimer\n"
	InvalidLicenseKeyMsg = "Unfortunately your license key is invalid, please check it and try once again."
	InsertBinanceKeyMsg = "Please insert your Binance api key and api secret in such a format: *key_secret*"
	InvalidBinanceKeysMsg = "Unfortunately your Binance api keys are invalid, please check if the format is correct and try once again."
	FeeLineMsgStructure = "TRADE: %s\n\nCLOSED DATE: %s\n\nPROFIT: %.4f\n\nFEE: %.4f\n-----------------------------------\n\n"
	ReportStartMsg = "FEE OVERVIEW - %s - %s\n\n"
	ReportEndMsg = "\nSummary:\n\nBased on %d amount of trades your open fees: $%.2f\n\nPlease transfer the open fee through this link: %s within the next 45 hours." +
		" Otherwise your account will be banned."
	NewlyGeneratedPaymentLinkMsg = "It's your new payment link: %s"
	PaymentReminderMsg = "We have not received your payment. If we won't receive your payment within the next 12 hours your account will be banned. It's your new payment link: %s"
)

var (
	ValidBinanceKeysMsg = "You Binance keys are valid!\n\n You can now join our telegram channel:\n %s"
)