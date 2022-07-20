package report

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/recruit-tech/RISUCON2022Summer/bench/fails"
	"github.com/recruit-tech/RISUCON2022Summer/bench/logger"
	"github.com/recruit-tech/RISUCON2022Summer/bench/random"
)

type Report struct {
	Pass       bool
	Reason     string
	ScoreSum   int64
	Level      int64
	FailReport fails.FailReport
	Language   string
	Webhook    WebhookData
}

type PublicReport struct {
	Pass     bool     `json:"pass"`
	Score    int      `json:"score"`
	Messages []string `json:"messages"`
}

func (r Report) PublicReport() (string, PublicReport) {
	id := random.ID()
	messages := make([]string, 0, 1+1+1+len(r.FailReport.Critical)+len(r.FailReport.Application))
	messages = append(messages, fmt.Sprintf("Report ID: %s", id))
	messages = append(messages, r.Reason)
	messages = append(messages, fmt.Sprintf("Load Level: %d", r.Level))
	for _, err := range r.FailReport.Critical {
		messages = append(messages, fmt.Sprintf("Critical Error: %v", err))
	}
	for _, err := range r.FailReport.Application {
		messages = append(messages, fmt.Sprintf("Application Error: %v", err))
	}
	return id, PublicReport{
		Pass:     r.Pass,
		Score:    int(r.ScoreSum),
		Messages: messages,
	}
}

func (r Report) Send() error {
	logger.Private.Printf("Language: %s\n", r.Language)
	logger.Private.Printf("Reason: %s\n", r.Reason)
	logger.Private.Printf("Score: %d\n", r.ScoreSum)
	logger.Private.Printf("Load Level: %d\n", r.Level)

	logger.Private.Println()
	logger.Private.Printf("Critical Error count: %d\n", len(r.FailReport.Critical))
	for _, err := range r.FailReport.Critical {
		logger.Private.Printf("\t%+v\n", err)
	}

	logger.Private.Println()
	logger.Private.Printf("Application Error count: %d\n", len(r.FailReport.Application))
	for _, err := range r.FailReport.Application {
		logger.Private.Printf("\t%+v\n", err)
	}

	logger.Private.Println()
	logger.Private.Printf("Trivial Error count: %d\n", len(r.FailReport.Trivial))
	for _, err := range r.FailReport.Trivial {
		logger.Private.Printf("\t%+v\n", err)
	}

	reportID, publicReport := r.PublicReport()
	b, err := json.Marshal(publicReport)
	if err != nil {
		return err
	}
	var buf bytes.Buffer
	err = json.Indent(&buf, []byte(b), "", "  ")
	if err != nil {
		return err
	}
	logger.Public.Println(buf.String())

	if r.Webhook.Token != "" {
		err = r.SendToWebhook(reportID, r.Webhook)
		if err != nil {
			return err
		}
	}

	return nil
}
