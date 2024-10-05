package GGCaptcha

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type SMSProvider string

const (
	ALI     SMSProvider = "Ali"
	TENCENT SMSProvider = "Tencent"
	BAIDU   SMSProvider = "Baidu"
	HUAXIN  SMSProvider = "HuaXin"
	GETUI   SMSProvider = "GeTui"
)

// SMSClient 短信客户端结构体
type SMSClient struct {
	Provider  SMSProvider
	APIKey    string
	APISecret string
	SMSUrl    string
}

// NewSMSClient 初始化 SMS 客户端
func NewSMSClient(provider SMSProvider, apiKey, apiSecret, smsUrl string) *SMSClient {
	return &SMSClient{
		Provider:  provider,
		APIKey:    apiKey,
		APISecret: apiSecret,
		SMSUrl:    smsUrl,
	}
}

// SendSMS 发送短信验证码
func (client *SMSClient) SendSMS(to, code string) (string, error) {
	var req = &http.Request{}
	switch client.Provider {
	case ALI:
		req = client.sendAliSMS(to, code)
	case TENCENT:
		req = client.sendTencentSMS(to, code)
	case BAIDU:
		req = client.sendBaiduSMS(to, code)
	case HUAXIN:
		req = client.sendHuaXinSMS(to, code)
	case GETUI:
		req = client.sendGeTuiSMS(to, code)
	default:
		return "", fmt.Errorf("unsupported SMS provider: %s", client.Provider)
	}

	// Check the error after calling doRequest
	if err := client.doRequest(req); err != nil {
		return "", err
	}
	return "短信验证码发送成功", nil
}

// 阿里云发送短信
func (client *SMSClient) sendAliSMS(to, code string) *http.Request {
	// 构造请求体
	body := map[string]interface{}{
		"PhoneNumbers":  to,
		"TemplateCode":  "SMS_12345678", // 模板代码
		"TemplateParam": fmt.Sprintf(`{"code":"%s"}`, code),
	}
	jsonData, _ := json.Marshal(body)

	req, _ := http.NewRequest("POST", client.SMSUrl, bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", client.APIKey)

	//client.doRequest(req)
	return req
}

// 腾讯云发送短信
func (client *SMSClient) sendTencentSMS(to, code string) *http.Request {
	body := map[string]interface{}{
		"PhoneNumberSet":   []string{to},
		"TemplateID":       "123456", // 模板 ID
		"TemplateParamSet": []string{code},
	}
	jsonData, _ := json.Marshal(body)

	req, _ := http.NewRequest("POST", client.SMSUrl, bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", client.APIKey)

	//client.doRequest(req)
	return req
}

// 百度云发送短信
func (client *SMSClient) sendBaiduSMS(to, code string) *http.Request {
	body := map[string]interface{}{
		"mobile": to,
		"code":   code,
	}
	jsonData, _ := json.Marshal(body)

	req, _ := http.NewRequest("POST", client.SMSUrl, bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("API-Key", client.APIKey)

	//client.doRequest(req)
	return req
}

// 华信发送短信
func (client *SMSClient) sendHuaXinSMS(to, code string) *http.Request {
	body := map[string]interface{}{
		"mobile":  to,
		"content": fmt.Sprintf("您的验证码是：%s", code),
	}
	jsonData, _ := json.Marshal(body)

	req, _ := http.NewRequest("POST", client.SMSUrl, bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", client.APIKey)

	//client.doRequest(req)
	return req
}

// 个推发送短信
func (client *SMSClient) sendGeTuiSMS(to, code string) *http.Request {
	body := map[string]interface{}{
		"recipient": to,
		"message":   fmt.Sprintf("您的验证码是：%s", code),
	}
	jsonData, _ := json.Marshal(body)

	req, _ := http.NewRequest("POST", client.SMSUrl, bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("API-Key", client.APIKey)

	//client.doRequest(req)
	return req
}

// 执行 HTTP 请求
func (client *SMSClient) doRequest(req *http.Request) error {
	clients := &http.Client{}
	resp, err := clients.Do(req)
	if err != nil {
		return fmt.Errorf("error sending SMS: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to send SMS, status: %s", resp.Status)
	}

	fmt.Printf("Response Status: %s\n", resp.Status)
	return nil
}
