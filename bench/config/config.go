package config

import (
	"io"

	"github.com/recruit-tech/RISUCON2022Summer/bench/report"
)

type Config struct {
	PublicWriter  io.Writer
	PrivateWriter io.Writer
	TargetUrl     string
	SnapshotDir   string
	Debug         bool
	Webhook       report.WebhookData
}
