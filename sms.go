package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

//SMSClient for https://rapidapi.com/JayemithLLC/api/quick-easy-sms
type SMSClient struct {
	//APIKey API Key can be obtained from https://rapidapi.com/JayemithLLC/api/quick-easy-sms
	//Choose a plan based on your need.
	APIKey string
	Client *http.Client
}

//SendMessage send message. Returns messageId and error message
func (c *SMSClient) SendMessage(phoneNumber, message string) (uint, error) {
	err := c.verifyCreds()
	if err != nil {
		return 0, err
	}

	form := url.Values{}
	form.Add("message", message)
	form.Add("toNumber", phoneNumber)
	req, err := http.NewRequest("POST", "https://quick-easy-sms.p.rapidapi.com/send", strings.NewReader(form.Encode()))
	if err != nil {
		return 0, errors.New("Failed to create new request to /send")
	}
	req.Header.Set("X-RapidAPI-Key", c.APIKey)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rawResp, err := c.Client.Do(req)
	if err != nil {
		return 0, fmt.Errorf("Error encountered while sending request to /send. Error: " + err.Error())
	}
	defer rawResp.Body.Close()
	bytesContent, err := ioutil.ReadAll(rawResp.Body)
	if err != nil {
		return 0, errors.New("Failed to read response content. Error: " + err.Error())
	}
	if rawResp.StatusCode != 200 {
		return 0, fmt.Errorf("Response status code %d received for /send. Message: %s", rawResp.StatusCode, string(bytesContent))
	}
	resp := &sendResponse{}

	err = json.Unmarshal(bytesContent, &resp)
	if err != nil {
		return 0, errors.New("Failed to unmarshal response. " + err.Error())
	}
	return resp.ID, nil
}

//Status check message is sent or not
func (c *SMSClient) Status(messageID uint) (bool, error) {
	err := c.verifyCreds()
	if err != nil {
		return false, err
	}
	req, err := http.NewRequest("GET", fmt.Sprintf("https://quick-easy-sms.p.rapidapi.com/status/%d", messageID), nil)
	if err != nil {
		return false, errors.New("Failed to create new request to /status")
	}
	req.Header.Set("X-RapidAPI-Key", c.APIKey)
	rawResp, err := c.Client.Do(req)
	if err != nil {
		return false, fmt.Errorf("Error encountered while sending request to /status. Error: " + err.Error())
	}
	defer rawResp.Body.Close()
	bytesContent, err := ioutil.ReadAll(rawResp.Body)
	if err != nil {
		return false, errors.New("Failed to read response content. Error: " + err.Error())
	}
	if rawResp.StatusCode != 200 {
		return false, fmt.Errorf("Response status code %d received for /status. Message: %s", rawResp.StatusCode, string(bytesContent))
	}
	resp := &statusResponse{}

	err = json.Unmarshal(bytesContent, &resp)
	if err != nil {
		return false, errors.New("Failed to unmarshal response. " + err.Error())
	}
	return resp.Message == "Sent", nil
}

func (c *SMSClient) verifyCreds() error {
	if strings.Trim(c.APIKey, " ") == "" {
		return errors.New("Invalid RapidAPI API key")
	} else if c.Client == nil {
		c.Client = &http.Client{}
	}
	return nil
}
