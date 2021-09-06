package daaz

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"milliard-easy/daaz_simbank/context"
	"milliard-easy/daaz_simbank/log"
	"net/http"
	"net/url"
	"strings"

	"github.com/sirupsen/logrus"
)

// API for daazweb
type API struct {
	cfg *context.Config
}

// NewAPI instance
func NewAPI(cfg *context.Config) *API {
	return &API{
		cfg: cfg,
	}
}

// DeleteAllWallets clear all
func (a *API) DeleteAllWallets() error {
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/deleteAllAutonumbers/", a.cfg.Daazweb.BaseURL), nil)
	if err != nil {
		return err
	}
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", a.cfg.Daazweb.Token))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	respBody := struct {
		Status        string `json:"status"`
		WalletsStatus bool   `json:"wallets_status"`
	}{}

	response, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	logrus.Debugf(log.DebugColor, fmt.Sprintf("Delete all wallets response: %s", string(response)))

	err = json.Unmarshal(response, &respBody)
	if err != nil {
		return err
	}

	if respBody.Status != "success" {
		return fmt.Errorf("cannot delete all wallets")
	}

	return nil
}

// AddWallet create new wallet
func (a *API) AddWallet(wallet string, port string) error {
	if !strings.HasPrefix(wallet, "+") {
		wallet = fmt.Sprintf("+%s", wallet)
	}

	formData := url.Values{
		"wallet":         {wallet},
		"action":         {"add"},
		"limit":          {a.cfg.Daazweb.Limit},
		"autonumber_off": {"true"},
		"notify_off":     {"true"},
		"description":    {fmt.Sprintf("SIM (порт: %s)", port)},
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/editAutonumber/", a.cfg.Daazweb.BaseURL), strings.NewReader(formData.Encode()))
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", a.cfg.Daazweb.Token))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	respBody := struct {
		Status string `json:"status"`
	}{}

	response, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	logrus.WithFields(logrus.Fields{
		"wallet": wallet,
		"port":   port,
	}).Debugf(log.DebugColor, fmt.Sprintf("Add wallet response: %s", string(response)))

	err = json.Unmarshal(response, &respBody)
	if err != nil {
		return err
	}
	// TODO: check status

	return nil
}

// AddPayment when new payment fetched
func (a *API) AddPayment(lastFour, text string) error {
	params := url.Values{}
	params.Add("bank", a.cfg.Daazweb.Operator)
	params.Add("call_type", "sms")
	params.Add("text", text)
	params.Add("lastfour", lastFour)
	params.Add("token", a.cfg.Daazweb.Token)

	url := fmt.Sprintf("%s?%s", a.cfg.Daazweb.PaymentURL, params.Encode())
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	respBody := struct {
		Status      string `json:"status"`
		Description string `json:"description"`
	}{}

	response, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	logrus.WithField("lastFour", lastFour).Debugf(log.DebugColor, fmt.Sprintf("Add payment response: %s", string(response)))

	err = json.Unmarshal(response, &respBody)
	if err != nil {
		return err
	}

	if respBody.Status != "success" {
		return fmt.Errorf("cannot add payment %s", respBody.Description)
	}

	return nil
}

// ToggleWalletState enable or disable wallet
func (a *API) ToggleWalletState(wallet string, state bool) error {
	if !strings.HasPrefix(wallet, "+") {
		wallet = fmt.Sprintf("+%s", wallet)
	}

	enable := "off"
	if state {
		enable = "on"
	}

	req, err := http.NewRequest("GET", fmt.Sprintf("%s/toggleEnableWallet/?wallet=%s&enable=%s", a.cfg.Daazweb.BaseURL, wallet, enable), nil)
	if err != nil {
		return err
	}
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", a.cfg.Daazweb.Token))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	respBody := struct {
		Status        string `json:"status"`
		WalletsStatus bool   `json:"wallets_status"`
	}{}

	response, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	logrus.WithFields(logrus.Fields{
		"wallet": wallet,
		"state":  state,
	}).Debugf(log.DebugColor, fmt.Sprintf("Toggle wallet state response: %s", string(response)))

	err = json.Unmarshal(response, &respBody)
	if err != nil {
		return err
	}

	if respBody.Status != "success" {
		return fmt.Errorf("cannot toggle wallet state")
	}

	return nil
}
