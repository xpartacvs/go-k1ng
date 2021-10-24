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

type Channel string

type sms struct {
	core         core.Core
	module       core.Module
	channel      Channel
	destinations []string
	sid          string
	content      string
	template     string
}

type Sms interface {
	SetSenderId(senderId string) Sms
	SetChannel(channelType Channel) Sms
	SetContent(strContent string) Sms
	SetTemplate(templateName string) Sms
	AddDestinations(phoneNumbers ...string) Sms
	EmptyDestinations() Sms
	Reset() Sms
	Send() (*core.Response, error)
	SendAt(sendTime time.Time) (*core.Response, error)
}

const (
	ChannelRegular    Channel = "NON-2FA"
	ChannelOTP        Channel = "2FA"
	ChannelDefault    Channel = "DEFAULT"
	ChannelLongNumber Channel = "LONGNUMBER"

	endpointSend string = "api/v1/send"
)

func create(hosturl, apikey, password string, channelType Channel) (Sms, error) {
	c, err := core.New(hosturl, apikey, password)
	if err != nil {
		return nil, err
	}

	return &sms{
		core:    c,
		channel: channelType,
		module:  core.ModuleSms,
	}, nil
}

func Default(hosturl, apikey, password string) (Sms, error) {
	return create(hosturl, apikey, password, ChannelDefault)
}

func LongNumber(hosturl, apikey, password string) (Sms, error) {
	return create(hosturl, apikey, password, ChannelLongNumber)
}

func OTP(hosturl, apikey, password string) (Sms, error) {
	return create(hosturl, apikey, password, ChannelOTP)
}

func Regular(hosturl, apikey, password string) (Sms, error) {
	return create(hosturl, apikey, password, ChannelRegular)
}

func (s *sms) SetSenderId(senderId string) Sms {
	s.sid = senderId
	return s
}

func (s *sms) SetChannel(channelType Channel) Sms {
	s.channel = channelType
	return s
}

func (s *sms) SetContent(strContent string) Sms {
	s.content = strContent
	return s
}

func (s *sms) SetTemplate(templateName string) Sms {
	s.template = templateName
	return s
}

func (s *sms) urlValues() (url.Values, error) {
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

func (s *sms) AddDestinations(phoneNumbers ...string) Sms {
	s.destinations = append(s.destinations, phoneNumbers...)
	return s
}

func (s *sms) EmptyDestinations() Sms {
	s.destinations = nil
	return s
}

func (s *sms) Reset() Sms {
	return s.SetSenderId("").SetContent("").SetTemplate("").EmptyDestinations()
}

func (s *sms) Send() (*core.Response, error) {
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

func (s *sms) SendAt(sendTime time.Time) (*core.Response, error) {
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
