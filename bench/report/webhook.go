package report

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"strings"

	"github.com/recruit-tech/RISUCON2022Summer/bench/logger"
)

type WebhookData struct {
	Channels string
	Token    string
}

func (r Report) SendToWebhook(reportID string, webhook WebhookData) error {
	logs, err := logger.GetPrivateLogs()
	if err != nil {
		return err
	}

	text := strings.Trim(logs, "\n")

	var body bytes.Buffer
	writer := multipart.NewWriter(&body)
	part, err := writer.CreateFormFile("file", fmt.Sprintf("Report ID: %s", reportID))
	if err != nil {
		return err
	}
	_, err = part.Write([]byte(text))
	if err != nil {
		return err
	}

	w, err := writer.CreateFormField("channels")
	if err != nil {
		return err
	}

	_, err = w.Write([]byte(webhook.Channels))
	if err != nil {
		return err
	}

	contentType := writer.FormDataContentType()

	err = writer.Close()
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPut, "https://slack.com/api/files.upload", &body)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", contentType)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", webhook.Token))

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	defer res.Body.Close()
	defer io.Copy(ioutil.Discard, res.Body)

	rb, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}
	logger.Private.Printf("Webhook response: %q\n", string(rb))

	return nil
}
