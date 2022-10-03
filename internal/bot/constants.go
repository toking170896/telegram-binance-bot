package bot

//messages
const (
	EmptyUsernameErr      = "Your username is not set up in Telegram, in order to continue please setup a username after that type /start"
	AcceptTermsMsg        = "In order to activate you agree to our Terms of Service and Disclaimer which can be found here:\n\nTOS: www.domain.com/tos\nDisclaimer: www.domain.com/disclaimer\n"
	DeniedTermsErr        = "You have to accept our tos and disclaimer to continue."
	InsertLicenseKeyMsg   = "Please insert a license key:"
	ValidLicenseKeyMsg    = "Perfect! Your license key is valid!\n\n***Payment Fees and how its working***\nWe connect our fee-bot to your binance account with a read only function to fetch all trades from our network which are exectued on your binance account.\n\n***Important***\nThere is a difference between the Futures Market and the Spot Market.\n\n***Future Market***\nFor the future market we are connecting our fee-bot to your future read only binance api access. If you make profit we will automatically calculate the profit you have made and based on that the fee will be calculated.\n\n***Spot Market***\nFor the spot market the profit calculation will be a little bit different. We connect your account to our trading network using your binance read only api access.The calculation for the spot market will be: (SellPrice-EntryPrice)/EntryPrice . We also add 10% dust trades in your behaviour.This means the profit e.g from cornix does not match the profit inside your report, actually the profit in our report will be less than inside cornix and binance.\n\n***Payments and Report***\nWe have a fixed trading flat fee of 15%, this means we calculate the fees for each trade which is closed only with profit on your account. A report will be send out weekly where you can get a summary of all the fees and a payment link to binance pay.The link will be available for one hour, if you miss the payment window you simple can re-create a new one. If we have not recived the payment within 72 hours your account will be permantnly locked from our network.\n\n***Help***\nIf you have questions regarding your payment please send us a message to: billing@cashways.ai including your telegram username and license key."
	AcceptFeesRetryMsg    = "Please *Accept/Deny* our rules.\n\n***Payment Fees and how its working***\nWe connect our fee-bot to your binance account with a read only function to fetch all trades from our network which are exectued on your binance account.\n\n***Important***\nThere is a difference between the Futures Market and the Spot Market.\n\n***Future Market***\nFor the future market we are connecting our fee-bot to your future read only binance api access. If you make profit we will automatically calculate the profit you have made and based on that the fee will be calculated.\n\n***Spot Market***\nFor the spot market the profit calculation will be a little bit different. We connect your account to our trading network using your binance read only api access.The calculation for the spot market will be: (SellPrice-EntryPrice)/EntryPrice . We also add 10% dust trades in your behaviour.This means the profit e.g from cornix does not match the profit inside your report, actually the profit in our report will be less than inside cornix and binance.\n\n***Payments and Report***\nWe have a fixed trading flat fee of 15%, this means we calculate the fees for each trade which is closed only with profit on your account. A report will be send out weekly where you can get a summary of all the fees and a payment link to binance pay.The link will be available for one hour, if you miss the payment window you simple can re-create a new one. If we have not recived the payment within 72 hours your account will be permantnly locked from our network.\n\n***Help***\nIf you have questions regarding your payment please send us a message to: billing@cashways.ai including your telegram username and license key."
	DeniedFeesMsg         = "It’s required to get started to accept out TOS and Disclaimer which can be found here:\n\nTOS: www.domain.com/tos\nDisclaimer: www.domain.com/disclaimer\n"
	InvalidLicenseKeyMsg  = "Unfortunately your license key is invalid, please check it and try once again."
	InsertBinanceKeyMsg   = "Create a new binance api key with READ ONLY and insert your Binance api key and api secret in such a format: *apikey_apisecret*"
	InvalidBinanceKeysMsg = "Unfortunately your Binance api keys are invalid, please check if the format is correct and try once again."
	FeeLineMsgStructure   = "TRADE: %s\n\nCLOSED DATE: %s\n\nPROFIT: %.4f\n\nFEE: %.4f\n-----------------------------------\n\n"
	ReportStartMsg        = "FEE OVERVIEW - %s - %s\n\n"
	ReportEndMsg          = "\nSummary:\n\nBased on %d amount of trades your profit: $%.2f, open fees: $%.2f\n\nPlease transfer the open fee within the next 72 hours." +
		" Otherwise your account will be banned permanently.\nThe payment link will be valid for 1 hour. In case you missed the timeframe you can create a new payment link."
	NewlyGeneratedPaymentLinkMsg = "New payment link was generated. Please transfer the open fee. Otherwise your account will be banned."
	PaymentReminderMsg           = "We have not received your payment. If we won't receive your payment within the next 12 hours your account will be banned. \n\n"
)

var (
	ValidBinanceKeysMsg = "You Binance keys are valid!\n\n You can now join our telegram channel:\n %s"
)
