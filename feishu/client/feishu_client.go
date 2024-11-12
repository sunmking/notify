package client

import (
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/sunmking/notify/feishu/message"
	"io"
	"net/http"
	"strings"
	"time"
)

type Client struct {
	token              string
	keyWork            string // 自定义关键词
	timeOut            time.Duration
	insecureSkipVerify bool
}

func NewFeiShuClient(token string, keyWork string, timeOut time.Duration, insecureSkipVerify bool) *Client {
	if timeOut == 0 {
		timeOut = 5 * time.Second
	}
	return &Client{
		token:              token,
		keyWork:            keyWork,
		timeOut:            timeOut,
		insecureSkipVerify: insecureSkipVerify,
	}
}

const (
	Webhook = "https://open.feishu.cn/open-apis/bot/v2/hook/"
)

func (client *Client) Send(msg any) error {
	url := Webhook + client.token

	messageContent, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	payload := strings.NewReader(string(messageContent))
	request, err := http.NewRequest(http.MethodPost, url, payload)
	if err != nil {
		return err
	}

	httpClient := http.Client{
		Timeout: client.timeOut,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: client.insecureSkipVerify,
			},
		},
	}
	request.Header.Add("Content-Type", "application/json")
	response, err := httpClient.Do(request)
	if err != nil {
		return err
	}

	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		fmt.Println(err)
		return err
	}

	if response.StatusCode != http.StatusOK {
		errMessageResponse := &message.ErrMessageResponse{}
		err := json.Unmarshal(body, errMessageResponse)
		if err != nil {
			return err
		}
		return errors.New(errMessageResponse.Msg)
	}

	return nil
}

func (client *Client) SendTextMessage(text string) error {
	text = text + client.keyWork
	message := message.NewTextMessage(text)
	return client.Send(message)
}

func (client *Client) SendPostMessage(title string, content [][]message.PostMessageContentPostZhCnContent) error {
	title = title + client.keyWork
	message := message.NewPostMessage(title, content)
	return client.Send(message)
}

func (client *Client) SendImageMessage(imageKey string) error {
	message := message.NewImageMessage(imageKey)
	return client.Send(message)
}

func (client *Client) SendShareChatMessage(shareChatId string) error {
	message := message.NewShareChatMessage(shareChatId)
	return client.Send(message)
}

func (client *Client) SendInteractiveMessage() error {
	return nil
}
