package fixlib

import (
	"os"
	"path"
	"time"

	"github.com/quickfixgo/enum"
	"github.com/quickfixgo/quickfix"
	"github.com/shopspring/decimal"
)

// FixAPI struct provide communication with xnt services
type FixAPI struct {
	initiator *quickfix.Initiator
	RespChan  chan string
}

// Run FixAPI initiator
func (f FixAPI) Run() (err error) { err = f.initiator.Start(); return }

// Stop FixAPI initiator
func (f FixAPI) Stop() { f.initiator.Stop() }

// Send message to acceptor
func (f FixAPI) send(msg *quickfix.Message) (err error) {
	err = quickfix.Send(msg)
	return
}

// NewOrder request FIX44
func (f FixAPI) NewOrder(symbol, orderQty, account, price, stopPx, 
	sendercompID, targetCompID string,ordType enum.OrdType) error {
	return f.send(queryNewOrderSingle44(symbol, orderQty, account, price, 
		stopPx, sendercompID, targetCompID, ordType))
}

// CancelOrder request FIX44
func (f FixAPI) CancelOrder(cloID, symbol, orderQty string, side enum.Side) error {
	return f.send(queryOrderCancelRequest44(
		cloID, symbol, orderQty, side))
}

// OrderStatus request FIX44
func (f FixAPI) OrderStatus(senderCompID, targetCompID, symbol string) error {
	return f.send(queryOrderStatus44(senderCompID, targetCompID, symbol))
}

// ReplaceOrder request FIX44
func (f FixAPI) ReplaceOrder(origClo, cloID string,
	side enum.Side, ordType enum.OrdType) error {
	return f.send(queryReplaceOrder44(origClo, cloID, side, ordType))
}

// AccountSummary request FIX44
func (f FixAPI) AccountSummary(senderCompID, targetCompID string) error {
	return f.send(queryAccountSummaryRequest44(senderCompID, targetCompID))
}

// OrderMassStatus request FIX44
func (f FixAPI) OrderMassStatus(sendercompID, targetCompID string,
	reqType enum.MassStatusReqType) error {
	return f.send(queryOrderMassStatusRequest44(sendercompID, targetCompID,
		enum.MassStatusReqType_STATUS_FOR_ALL_ORDERS))
}

// SecurityList request FIX44
func (f FixAPI) SecurityList(sendercompID, targetCompID string) error {
	return f.send(querySecurityListRequest44(sendercompID, targetCompID))
}

// SecurityListSymbol request FIX44
func (f FixAPI) SecurityListSymbol(sym, sendercompID, targetCompID string) error {
	return f.send(querySecurityListSymbolRequest44(
		sym, sendercompID, targetCompID))
}

// SecurityListCFI request FIX44
func (f FixAPI) SecurityListCFI(cfi, sendercompID, targetCompID string,
	secList enum.SecurityListRequestType) error {
	return f.send(querySecurityListCFICODERequest44(
		cfi, sendercompID, targetCompID))
}

// TradesCapture request 44
func (f FixAPI) TradesCapture(sendercompID, targetCompID string,
		startTime, stopTime time.Time) error {
	return f.send(queryTradesCapture44(
		sendercompID, targetCompID, startTime, stopTime))
}

// TradesCaptureStopDate request 44
func (f FixAPI) TradesCaptureStopDate(
	targetCompID, sendercompID, stopDate string) error {
	return f.send(queryTradesCaptureStopDate44(
		targetCompID, sendercompID, stopDate))
}

// TradesMargin ...
func (f FixAPI) TradesMargin(symbol, account, currency, sendercompID,
	targetCompID string,
	quantity, price decimal.Decimal) error {
	return f.send(queryTradeMarginRequest44(symbol, account, currency,
		sendercompID, targetCompID, quantity, price))
}

// TestRequest ...
func (f FixAPI) TestRequest(requestID string) error {
	msg := queryTestRequest44(requestID)
	err := f.send(msg)
	return err
}

// MarketData ...
func (f FixAPI) MarketData(marketDepth int, symbol, sendercompID,targetCompID string) error {
	return f.send(queryMarketDataRequest44(marketDepth, sendercompID, symbol,
		targetCompID ))
}

// NewFixAPI constructor
func NewFixAPI(sessionCfgPath string, respChanMaxSize int) (FixAPI, error) {
	cfgFileName := path.Join(sessionCfgPath)
	var cfg *os.File
	var err error
	var fileLogFactory quickfix.LogFactory
	var appSettings *quickfix.Settings
	var initiator *quickfix.Initiator
	if cfg, err = os.Open(cfgFileName); err != nil {
		return FixAPI{}, err
	}
	defer cfg.Close()
	if appSettings, err = quickfix.ParseSettings(cfg); err != nil {
		return FixAPI{}, err
	}
	if fileLogFactory, err = quickfix.NewFileLogFactory(appSettings); err != nil {
		return FixAPI{}, err
	}
	fileStoreFactory := quickfix.NewFileStoreFactory(appSettings)
	respChan := make(chan string, respChanMaxSize)
	tradeClient := TradeClient{respChan}
	if initiator, err = quickfix.NewInitiator(tradeClient, fileStoreFactory,
		appSettings, fileLogFactory); err != nil {
		return FixAPI{}, err
	}
	return FixAPI{
		initiator: initiator,
		RespChan:  respChan,
	}, nil
}
