package fixlib

import (
	"fmt"

	"github.com/quickfixgo/quickfix"
)

// TradeClient implements the quickfix.Application interface
type TradeClient struct{ respChan chan string }

// handleResp provide interface to response for main app
func (t TradeClient) handleResp(val string) { t.respChan <- val }

//OnCreate implemented as part of Application interface
func (t TradeClient) OnCreate(sessionID quickfix.SessionID) {
	t.handleResp(fmt.Sprintf("OnCreate. SessionID: %s\n", sessionID))
	return
}

//OnLogon implemented ƒas part of Application interface
func (t TradeClient) OnLogon(sessionID quickfix.SessionID) {
	t.handleResp(fmt.Sprintf("OnLogon. SessionID: %s\n", sessionID))
	return
}

//OnLogout implemented as part of Application interface
func (t TradeClient) OnLogout(sessionID quickfix.SessionID) {
	t.handleResp(fmt.Sprintf("OnLogout. SessionID: %s\n", sessionID))
	return
}

//FromAdmin implemented as part of Application interface
func (t TradeClient) FromAdmin(
	msg *quickfix.Message,
	sessionID quickfix.SessionID) (reject quickfix.MessageRejectError) {

	t.handleResp(fmt.Sprintf(
		"FromAdmin. SessionID:%s  Msg:%s \n", sessionID, msg.String()))
	return
}

//ToAdmin implemented as part of Application interface
func (t TradeClient) ToAdmin(msg *quickfix.Message, sessionID quickfix.SessionID) {
	t.respChan <- fmt.Sprintf("ToAdmin. SessionID: %s\n", sessionID)
	ModifyMsg(msg)
	return
}

//ToApp implemented as part of Application interface
func (t TradeClient) ToApp(
	msg *quickfix.Message, sessionID quickfix.SessionID) (err error) {
	t.respChan <- fmt.Sprintf(
		"ToApp. SessionID: %s Sending: %s\n", sessionID, msg.String())
	return
}

//FromApp implemented as part of Application interface. This is the callback for all Application level messages from the counter party.
func (t TradeClient) FromApp(msg *quickfix.Message,
	sessionID quickfix.SessionID) (reject quickfix.MessageRejectError) {
	t.respChan <- msg.String()
	return
}
