package okcoin

import (
	"net/url"
	"strings"
	"log"
	"errors"
	"github.com/zaviichen/goex/common"
)

func SendAuthenticatedHTTPRequest(api, secret, baseuri string, verbose bool, method string, v url.Values, result interface{}) (err error) {
	v.Set("api_key", api)
	hasher := common.GetMD5([]byte(v.Encode() + "&secret_key=" + secret))
	v.Set("sign", strings.ToUpper(common.HexEncodeToString(hasher)))

	encoded := v.Encode()
	path := baseuri + method

	if verbose {
		log.Printf("Sending POST request to %s with params %s\n", path, encoded)
	}

	headers := make(map[string]string)
	headers["Content-Type"] = "application/x-www-form-urlencoded"

	resp, err := common.SendHTTPRequest("POST", path, headers, strings.NewReader(encoded))

	if err != nil {
		return err
	}

	if verbose {
		log.Printf("Recieved raw: \n%s\n", resp)
	}

	err = common.JSONDecode([]byte(resp), &result)

	if err != nil {
		return errors.New("Unable to JSON Unmarshal response.")
	}

	return nil
}

func (o *OKSpot) PostRequest(method string, v url.Values, result interface{}) (err error) {
	return SendAuthenticatedHTTPRequest(o.APIKey, o.APISecret, o.APIUrl, o.Verbose, method, v, result)
}

func (o *OKEx) PostRequest(method string, v url.Values, result interface{}) (err error) {
	return SendAuthenticatedHTTPRequest(o.APIKey, o.APISecret, o.APIUrl, o.Verbose, method, v, result)
}

func RESTErrorMessage(code string) (string) {
	errMsgs := map[string]string{
		"10000": "Required field, can not be null",
		"10001": "Request frequency too high",
		"10002": "System error",
		"10003": "Not in reqest list, please try again later",
		"10004": "IP not allowed to access the resource",
		"10005": "'secretKey' does not exist",
		"10006": "'partner' does not exist",
		"10007": "Signature does not match",
		"10008": "Illegal parameter",
		"10009": "Order does not exist",
		"10010": "Insufficient funds",
		"10011": "Amount too low",
		"10012": "Only btc_usd/btc_cny ltc_usd,ltc_cny supported",
		"10013": "Only support https request",
		"10014": "Order price must be between 0 and 1,000,000",
		"10015": "Order price differs from current market price too much",
		"10016": "Insufficient coins balance",
		"10017": "API authorization error",
		"10018": "Borrow amount less than lower limit [usd/cny:100,btc:0.1,ltc:1]",
		"10019": "Loan agreement not checked",
		"10020": `Rate cannot exceed 1%`,
		"10021": `Rate cannot less than 0.01%`,
		"10023": "Fail to get latest ticker",
		"10024": "Balance not sufficient",
		"10025": "Quota is full, cannot borrow temporarily",
		"10026": "Loan (including reserved loan) and margin cannot be withdrawn",
		"10027": "Cannot withdraw within 24 hrs of authentication information modification",
		"10028": "Withdrawal amount exceeds daily limit",
		"10029": "Account has unpaid loan, please cancel/pay off the loan before withdraw",
		"10031": "Deposits can only be withdrawn after 6 confirmations",
		"10032": "Please enabled phone/google authenticator",
		"10033": "Fee higher than maximum network transaction fee",
		"10034": "Fee lower than minimum network transaction fee",
		"10035": "Insufficient BTC/LTC",
		"10036": "Withdrawal amount too low",
		"10037": "Trade password not set",
		"10040": "Withdrawal cancellation fails",
		"10041": "Withdrawal address not approved",
		"10042": "Admin password error",
		"10043": "Account equity error, withdrawal failure",
		"10044": "fail to cancel borrowing order",
		"10047": "This function is disabled for sub-account",
		"10100": "User account frozen",
		"10216": "Non-available API",
		"20001": "User does not exist",
		"20002": "Account frozen",
		"20003": "Account frozen due to liquidation",
		"20004": "Futures account frozen",
		"20005": "User futures account does not exist",
		"20006": "Required field missing",
		"20007": "Illegal parameter",
		"20008": "Futures account balance is too low",
		"20009": "Future contract status error",
		"20010": "Risk rate ratio does not exist",
		"20011": `Risk rate higher than 90% before opening position`,
		"20012": `Risk rate higher than 90% after opening position`,
		"20013": "Temporally no counter party price",
		"20014": "System error",
		"20015": "Order does not exist",
		"20016": "Close amount bigger than your open positions",
		"20017": "Not authorized/illegal operation",
		"20018": `Order price differ more than 5% from the price in the last minute`,
		"20019": "IP restricted from accessing the resource",
		"20020": "secretKey does not exist",
		"20021": "Index information does not exist",
		"20022": "Wrong API interface (Cross margin mode shall call cross margin API, fixed margin mode shall call fixed margin API)",
		"20023": "Account in fixed-margin mode",
		"20024": "Signature does not match",
		"20025": "Leverage rate error",
		"20026": "API Permission Error",
		"20027": "No transaction record",
		"20028": "No such contract",
	}
	return errMsgs[code]
}
