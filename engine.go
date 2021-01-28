package fixlib

import (
	"time"

	//"github.com/beevik/guid"
	"github.com/quickfixgo/enum"
	"github.com/quickfixgo/field"
	fix44slr "github.com/quickfixgo/fix44/SecurityListRequest"
	fix44mdr "github.com/quickfixgo/fix44/marketdatarequest"
	fix44nos "github.com/quickfixgo/fix44/newordersingle"
	fix44ocrr "github.com/quickfixgo/fix44/ordercancelreplacerequest"
	fix44cxl "github.com/quickfixgo/fix44/ordercancelrequest"
	fix44oms "github.com/quickfixgo/fix44/ordermassstatusrequest"
	fix44osr "github.com/quickfixgo/fix44/orderstatusrequest"
	fix44trq "github.com/quickfixgo/fix44/testrequest"
	fix44tmr "github.com/quickfixgo/fix44/tradecapturereport"
	"github.com/quickfixgo/quickfix"
	"github.com/shopspring/decimal"
)

const (
	scale int32 = 2

	accountSummaryReqType string = "UASQ"
)

func queryOrigClOrdID() field.OrigClOrdIDField {
	return field.NewOrigClOrdID(("OrigClOrdID"))
}

func queryDecimal(s string) decimal.Decimal {
	if v, err := decimal.NewFromString(s); err != nil {
		panic(err)
	} else {
		return v
	}
}

func querySenderCompID(val string) field.SenderCompIDField {
	return field.NewSenderCompID(val)
}

func queryTargetCompID(val string) field.TargetCompIDField {
	return field.NewTargetCompID(val)
}

func queryTradeDate(val string) field.TradeDateField {
	return field.NewTradeDate(val)
}

type header interface {
	Set(f quickfix.FieldWriter) *quickfix.FieldMap
}

func queryAccount(val string) field.AccountField {
	return field.NewAccount(val)
}

func queryCurrency(val string) field.CurrencyField {
	return field.NewCurrency(val)
}

func queryQuantity(val decimal.Decimal) field.QuantityField {
	return field.NewQuantity(val, scale)
}

func queryPrice(val decimal.Decimal) field.PriceField {
	return field.NewPrice(val, scale)
}


func queryHeader(h header, targetCompID, senderCompID string) {
	h.Set(querySenderCompID(senderCompID))
	h.Set(queryTargetCompID(targetCompID))
	h.Set(queryBeginString())
	//h.Set(queryTargetSubID())
}

func queryMsgType(val enum.MsgType) field.MsgTypeField {
	return field.NewMsgType(val)
}

func queryBeginString() field.BeginStringField {
	return field.NewBeginString("FIX.4.4")
}

func queryTargetSubID() field.TargetSubIDField {
	return field.NewTargetSubID("5")
}

func queryRequestID(ID string) field.TestReqIDField {
	return field.NewTestReqID(ID)
}

func queryTransactTime() field.TransactTimeField {
	return field.NewTransactTime(time.Now())
}

func querySide(val enum.Side) field.SideField {
	return field.NewSide(val)
}

func queryNewCloID(val string) field.ClOrdIDField {
	return field.NewClOrdID(val)
}

func queryCFI(val string) field.CFICodeField {
	return field.NewCFICode(val)
}

func querySymbol(val string) field.SymbolField {
	return field.NewSymbol(val)
}

func queryOrderQty(val string) field.OrderQtyField {
	return field.NewOrderQty(queryDecimal(val), scale)
}

func queryOrdType(val enum.OrdType) field.OrdTypeField {
	return field.NewOrdType(val)
}

func queryTradeRequestType() field.TradeRequestTypeField {
	return field.NewTradeRequestType(enum.TradeRequestType_ALL_TRADES)
}

func queryNoDates(val bool) field.NoDatesField {
	if val {
		return field.NewNoDates(2)
	}
	return field.NewNoDates(1)
}

func generateReqID() string { return "guid.New().String()" }

func setReqID(msg *quickfix.Message, ID string) {
	msg.Body.SetField(20020, quickfix.FIXString(ID))
}

func queryTestRequest44(requestID string) (msg *quickfix.Message) {
	testRequest := fix44trq.New(field.NewTestReqID(requestID))
	t := field.NewTargetCompID("EXANTE_TRADE_UAT")
	testRequest.Header.Set(t)
	msg = testRequest.ToMessage()
	return
}

func queryNewOrderSingle44(symbol, orderQty, account, price, stopPx, sendercompID, targetCompID string,
		ordType enum.OrdType) (msg *quickfix.Message) {
	query := fix44nos.New(
		queryNewCloID(generateReqID()),
		querySide(enum.Side_BUY),
		queryTransactTime(),
		queryOrdType(ordType),
	)
	query.SetHandlInst("1")
	query.Set(field.NewAccount(account))
	query.Set(querySymbol(symbol))
	query.Set(queryOrderQty(orderQty))
	switch ordType {
	case enum.OrdType_LIMIT, enum.OrdType_STOP_LIMIT:
		query.Set(field.NewPrice(
			queryDecimal(price), scale),
		)
	}
	switch ordType {
	case enum.OrdType_STOP, enum.OrdType_STOP_LIMIT:
		query.Set(field.NewStopPx(queryDecimal(stopPx), scale))
	}
	msg = query.ToMessage()
	queryHeader(&msg.Header, targetCompID, sendercompID)
	return
}

func queryOrderCancelRequest44(cloID, symbol, orderQty string,
	side enum.Side) (msg *quickfix.Message) {
	cancel := fix44cxl.New(queryOrigClOrdID(),
		queryNewCloID(cloID),
		querySide(side),
		queryTransactTime(),
	)
	cancel.Set(querySymbol(symbol))
	cancel.Set(queryOrderQty(orderQty))

	msg = cancel.ToMessage()
	return quickfix.NewMessage()
}

func queryTradesCuptureRequest(tradeReportID, tradeData string, prevReportID bool,
	lastQty, lastPX decimal.Decimal) (msg *quickfix.Message) {
	query := fix44tmr.New(field.NewTradeReportID(tradeReportID),
		field.NewPreviouslyReported(prevReportID),
		field.NewLastQty(lastQty, scale),
		field.NewLastPx(lastPX, scale),
		field.NewTradeDate(tradeData),
		queryTransactTime(),
	)
	msg = query.ToMessage()
	return
}

func querySecurityListSymbolRequest44(symbol, sendercompID, targetCompID string) (msg *quickfix.Message) {
	reqID := field.NewSecurityReqID(generateReqID())
	rType := field.NewSecurityListRequestType(
		enum.SecurityListRequestType_SYMBOL)
	query := fix44slr.New(reqID, rType)
	query.Body.Set(querySymbol(symbol))
	msg = query.ToMessage()
	queryHeader(&msg.Header, targetCompID, sendercompID)
	return
}

func querySecurityListCFICODERequest44(cfi, sendercompID, targetCompID string) (msg *quickfix.Message) {
	reqID := field.NewSecurityReqID(generateReqID())
	rType := field.NewSecurityListRequestType(
		enum.SecurityListRequestType_SECURITYTYPE_AND_OR_CFICODE)
	query := fix44slr.New(reqID, rType)
	query.Body.Set(queryCFI(cfi))
	msg = query.ToMessage()
	queryHeader(&msg.Header, targetCompID, sendercompID)
	return
}

func querySecurityListRequest44(sendercompID, targetCompID string) (msg *quickfix.Message) {
	reqID := field.NewSecurityReqID(generateReqID())
	rType := field.NewSecurityListRequestType(
		enum.SecurityListRequestType_ALL_SECURITIES)
	query := fix44slr.New(reqID, rType)
	msg = query.ToMessage()
	queryHeader(&msg.Header, targetCompID, sendercompID)
	return
}

func queryTradesCapture44(sendercompID, targetCompID string,
		startTime, stopTime time.Time) (msg *quickfix.Message) {
	query := quickfix.NewMessage()
	setReqID(query, generateReqID())
	query.Body.Set(queryTradeRequestType())
	query.Body.Set(queryNoDates(true))
	query.Header.Set(field.NewTradeRequestID(generateReqID()))
	query.Body.Set(field.NewTradSesStartTime(startTime))
	query.Body.Set(field.NewTradSesCloseTime(stopTime))
	query.Header.SetField(35, quickfix.FIXString("AD"))
	query.Body.Set(queryAccount("test"))
	query.Body.SetField(58, quickfix.FIXString("AD"))
	queryHeader(&query.Header, targetCompID, sendercompID)
	msg = query.ToMessage()
	return
}

func queryTradesCaptureStopDate44(
	targetCompID, sendercompID, stopDate string) (msg *quickfix.Message) {
	query := quickfix.NewMessage()
	setReqID(query, generateReqID())
	query.Body.Set(queryTradeRequestType())
	query.Body.Set(queryNoDates(true))
	query.Body.Set(queryTradeDate(stopDate))
	query.Header.SetField(35, quickfix.FIXString("UTMQ"))
	queryHeader(&query.Header, targetCompID, sendercompID)
	msg = query.ToMessage()
	return
}


func queryTradeMarginRequest44(symbol, account, currency, sendercompID, 
	targetCompID  string,
	quantity, price decimal.Decimal) (msg *quickfix.Message){
	query := quickfix.NewMessage()
	setReqID(query, generateReqID())
	//TODO: dublicate
	query.Header.SetField(35, quickfix.FIXString("UTMQ"))
	query.Body.Set(querySymbol(symbol))
	query.Body.Set(queryAccount(account))
	query.Body.Set(queryCurrency(currency))
	query.Body.Set(queryQuantity(quantity))
	if price.BigInt().Int64() >= 0 {
		query.Body.Set(queryPrice(price))
	}
	// request ID
	query.Body.SetField(20050, quickfix.FIXString(generateReqID()))
	queryHeader(&query.Header, targetCompID, sendercompID)
	msg = query.ToMessage()
	return
}

func queryAccountSummaryRequest44(sendercompID, targetCompID string) (msg *quickfix.Message) {
	query := quickfix.NewMessage()
	//TODO: dublicate
	query.Header.SetField(35, quickfix.FIXString("UASQ"))

	// request ID
	query.Body.SetField(20020, quickfix.FIXString("325"))
	queryHeader(&query.Header, targetCompID, sendercompID)
	msg = query.ToMessage()
	return
}

func queryOrderMassStatusRequest44(sendercompID, targetCompID string,
	r enum.MassStatusReqType) (msg *quickfix.Message) {
	query := fix44oms.New(field.NewMassStatusReqID(generateReqID()),
		field.NewMassStatusReqType(r),
	)
	queryHeader(&query.Header, targetCompID, sendercompID)
	msg = query.ToMessage()
	return
}

func queryMarketDataRequest44(marketDepth int, symbol, 
	sendercompID, targetCompID string) (msg *quickfix.Message) {
	query := fix44mdr.New(field.NewMDReqID(generateReqID()),
		field.NewSubscriptionRequestType(enum.SubscriptionRequestType_SNAPSHOT_PLUS_UPDATES),
		field.NewMarketDepth(marketDepth),
	)

	entryTypes := fix44mdr.NewNoMDEntryTypesRepeatingGroup()
	entryTypes.Add().SetMDEntryType(enum.MDEntryType_BID)
	query.SetNoMDEntryTypes(entryTypes, )

	relatedSymGroup := fix44mdr.NewNoRelatedSymRepeatingGroup()
	relatedSymGroup.Add().SetSymbol(symbol)
	query.SetNoRelatedSym(relatedSymGroup)

	queryHeader(query.Header, targetCompID, sendercompID)
	msg = query.ToMessage()
	return msg
}

func queryOrderStatus44(senderCompID, targetCompID, symbol string) (msg *quickfix.Message) {
	query := fix44osr.New(queryNewCloID(generateReqID()), querySide(enum.Side_BUY))
	query.Body.Set(querySymbol(symbol))
	queryHeader(query.Header, targetCompID, senderCompID)
	msg = query.ToMessage()
	return msg
}

func queryReplaceOrder44(origClo, cloID string,
	side enum.Side, ordType enum.OrdType) (msg *quickfix.Message) {
	replace := fix44ocrr.New(
		field.NewOrigClOrdID(origClo),
		queryNewCloID(cloID),
		querySide(side),
		queryTransactTime(),
		queryOrdType(ordType),
	)
	msg = replace.ToMessage()
	return quickfix.NewMessage()
}

// ModifyMsg ...
func ModifyMsg(msg *quickfix.Message) {
	var msgType quickfix.FIXString
	if err := msg.Header.GetField(35, &msgType); err != nil {
	}
	if msgType == "A" {
		msg.Header.SetField(554, quickfix.FIXString("Bs3CHeEMSk"))
		msg.Body.SetField(141, quickfix.FIXBoolean(true))
	}
}
