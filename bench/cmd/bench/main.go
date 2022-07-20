package main

import (
	"flag"
	"io/ioutil"
	"os"

	"github.com/recruit-tech-tech/RISUCON2022Summer/bench/config"
	"github.com/recruit-tech-tech/RISUCON2022Summer/bench/report"
	"github.com/recruit-tech/RISUCON2022Summer/bench"
)

func main() {
	var (
		targetUrl       string
		snapshotsDir    string
		webhookChannels string
		webhookToken    string
		debug           bool
	)

	flag.StringVar(&targetUrl, "target", "http://localhost:3000", "target url")
	flag.StringVar(&snapshotsDir, "snapshots-dir", "../snapshots", "path for snapshots directory")
	flag.StringVar(&webhookChannels, "webhook-channels", "", "Target channels for sending private report over webhook")
	flag.StringVar(&webhookToken, "webhook-token", "", "Webhook token for sending private report")
	flag.BoolVar(&debug, "debug", false, "debug mode")
	flag.Parse()

	config := config.Config{
		PublicWriter:  os.Stdout,
		PrivateWriter: os.Stderr,
		TargetUrl:     targetUrl,
		SnapshotDir:   snapshotsDir,
		Debug:         debug,
		Webhook: report.WebhookData{
			Channels: webhookChannels,
			Token:    webhookToken,
		},
	}

	if debug {
		config.PublicWriter = ioutil.Discard
		config.PrivateWriter = os.Stdout
	}

	passed := bench.Run(config)
	if !passed {
		os.Exit(1)
	}
	os.Exit(0)
}
