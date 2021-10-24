package sms

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/xpartacvs/go-k1ng/core"
)

// A type to group several SMS channel
type Channel string

// A main struct of SMS client where you can configure and send SMS.
type Sms struct {
	core         core.Core
	module       core.Module
	channel      Channel
	destinations []string
	sid          string
	content      string
	template     string
}

const (
	// A dedicated channel for indonesia customer to send non OTP SMS
	ChannelRegular Channel = "NON-2FA"

	// A dedicated channel for indonesia customer to send OTP SMS
	ChannelOTP Channel = "2FA"

	// A dedicated channel for international customer to send SMS to indonesian number
	ChannelDefault Channel = "DEFAULT"

	// A dedicated channel to send SMS using longnumber as it's sender ID
	ChannelLongNumber Channel = "LONGNUMBER"

	endpointSend string = "api/v1/send"
)

func create(hosturl, apikey, password string, channelType Channel) (*Sms, error) {
	c, err := core.New(hosturl, apikey, password)
	if err != nil {
		return nil, err
	}

	return &Sms{
		core:    c,
		channel: channelType,
		module:  core.ModuleSms,
	}, nil
}

// Default func is a helper to create a SMS client with `ChannelDefault` as it's preferred channel.
// Please note that channel is mutable, meaning that you might change it in the future by `SetChannel` method.
func Default(hosturl, apikey, password string) (*Sms, error) {
	return create(hosturl, apikey, password, ChannelDefault)
}

// LongNumber func is a helper to create a SMS client with `ChannelLongNumber` as it's preferred channel.
// Please note that channel is mutable, meaning that you might change it in the future by `SetChannel` method.
func LongNumber(hosturl, apikey, password string) (*Sms, error) {
	return create(hosturl, apikey, password, ChannelLongNumber)
}

// OTP func is s helper to create a SMS client with `ChannelOTP` as it's preferred channel.
// Please note that channel is mutable, meaning that you might change it in the future by `SetChannel` method.
func OTP(hosturl, apikey, password string) (*Sms, error) {
	return create(hosturl, apikey, password, ChannelOTP)
}

// Regular func is helper to create a SMS client with `ChannelRegular` as it's preferred channel.
// Please note that channel is mutable, meaning that you might change it in the future by `SetChannel` method.
func Regular(hosturl, apikey, password string) (*Sms, error) {
	return create(hosturl, apikey, password, ChannelRegular)
}

// SetSenderId assign the given argument as sender ID of the SMS.
// Please note this method return itself meaning it can be used as chaining-function.
func (s *Sms) SetSenderId(senderId string) *Sms {
	s.sid = senderId
	return s
}

// SetChannel assign the given argument as the main SMS channel.
// Please note this method return itself meaning it can be used as chaining-function.
func (s *Sms) SetChannel(channelType Channel) *Sms {
	s.channel = channelType
	return s
}

// SetContent assign the given argument as SMS body message.
// A coma-separated string can be passed as it's argument in order to fill template placeholder
// Please note this method return itself meaning it can be used as chaining-function.
func (s *Sms) SetContent(strContent string) *Sms {
	s.content = strContent
	return s
}

// SetTemplate assign the given argument to be used as the main SMS template.
// Use `SetContent` method to fill it's placeholder
// Please note this method return itself meaning it can be used as chaining-function.
func (s *Sms) SetTemplate(templateName string) *Sms {
	s.template = templateName
	return s
}

func (s Sms) urlValues() (url.Values, error) {
	data := url.Values{}
	data.Add("module", string(s.module))
	data.Add("sub_module", string(s.channel))

	if len(s.sid) <= 0 {
		return data, errors.New("sender id must not empty")
	}
	data.Add("sid", s.sid)

	if len(s.sid) <= 0 {
		return data, errors.New("message content must not empty")
	}
	data.Add("content", s.content)

	if len(strings.TrimSpace(s.template)) > 0 {
		data.Add("template_name", s.template)
	}

	if len(s.destinations) <= 0 {
		return data, errors.New("destination content must not empty")
	}
	data.Add("destination", strings.Join(s.destinations, ","))

	return data, nil
}

// AddDestination fills the destiantion pool with the given arguments.
// Please note this method return itself meaning it can be used as chaining-function.
func (s *Sms) AddDestination(phoneNumbers ...string) *Sms {
	s.destinations = append(s.destinations, phoneNumbers...)
	return s
}

// EmptyDestination flush the destination pool.
// Please note this method return itself meaning it can be used as chaining-function.
func (s *Sms) EmptyDestination() *Sms {
	s.destinations = nil
	return s
}

// Reset will empty the sender id, body message, and destination pool.
// Please note this method return itself meaning it can be used as chaining-function.
func (s *Sms) Reset() *Sms {
	return s.SetSenderId("").SetContent("").SetTemplate("").EmptyDestination()
}

// Send the configured SMS to the target in destination pool immediately.
func (s Sms) Send() (*core.Response, error) {
	data, err := s.urlValues()
	if err != nil {
		return nil, err
	}

	resp, err := s.core.ConsumeAPI(core.MethodPost, endpointSend, data)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	bytesResp, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		errMsg := "http status code" + strconv.Itoa(resp.StatusCode)
		return nil, errors.New(errMsg)
	}

	ret := new(core.Response)
	if err = json.Unmarshal(bytesResp, ret); err != nil {
		return nil, err
	}

	return ret, nil
}

// Send the configured SMS to the target in destination pool at the given argument (time).
func (s Sms) SendAt(sendTime time.Time) (*core.Response, error) {
	data, err := s.urlValues()
	if err != nil {
		return nil, err
	}
	data.Add("schedule_time", sendTime.Format(core.TimeFormatMySQLDateTime))

	resp, err := s.core.ConsumeAPI(core.MethodPost, endpointSend, data)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	bytesResp, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		errMsg := "http status code" + strconv.Itoa(resp.StatusCode)
		return nil, errors.New(errMsg)
	}

	ret := new(core.Response)
	if err = json.Unmarshal(bytesResp, ret); err != nil {
		return nil, err
	}

	return ret, nil
}
