package sms_provider

import (
	// "encoding/json"
	// "errors"
	"fmt"
	// "net/http"
	// "net/url"
	// "strings"

	"github.com/supabase/gotrue/internal/conf"
	// "github.com/supabase/gotrue/internal/utilities"
	// "golang.org/x/exp/utf8string"

	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	terrors "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	sms "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/sms/v20210111"
)

const (
	defaultVonageApiBase = "https://rest.nexmo.com"
)

type VonageProvider struct {
	Config  *conf.VonageProviderConfiguration
	APIPath string
}

type VonageResponseMessage struct {
	Status    string `json:"status"`
	ErrorText string `json:"error-text"`
}

type VonageResponse struct {
	Messages []VonageResponseMessage `json:"messages"`
}

// Creates a SmsProvider with the Vonage Config
func NewVonageProvider(config conf.VonageProviderConfiguration) (SmsProvider, error) {
	if err := config.Validate(); err != nil {
		return nil, err
	}

	apiPath := defaultVonageApiBase + "/sms/json"
	return &VonageProvider{
		Config:  &config,
		APIPath: apiPath,
	}, nil
}

func (t *VonageProvider) SendMessage(phone string, message string, channel string) error {
	switch channel {
	case SMSProvider:
		return t.SendSms(phone, message)
	default:
		return fmt.Errorf("channel type %q is not supported for Vonage", channel)
	}
}

// Send an SMS containing the OTP with Vonage's API
func (t *VonageProvider) SendSms(phone string, message string) error {
	credential := common.NewCredential(
		t.Config.ApiKey,
		t.Config.ApiSecret,
	)
	cpf := profile.NewClientProfile()
	cpf.HttpProfile.Endpoint = "sms.tencentcloudapi.com"
	client, _ := sms.NewClient(credential, "ap-beijing", cpf)

	request := sms.NewSendSmsRequest()

	request.PhoneNumberSet = common.StringPtrs([]string{"+86" + phone})
	request.SmsSdkAppId = common.StringPtr("1400817746")
	request.SignName = common.StringPtr("寄云科技")
	request.TemplateId = common.StringPtr("977982")
	request.TemplateParamSet = common.StringPtrs([]string{message, "5"})

	_, err := client.SendSms(request)
	if _, ok := err.(*terrors.TencentCloudSDKError); ok {
		return fmt.Errorf("An API error has returned: %s", err)
	}

	return nil
}
