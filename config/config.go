package config

import (
	"flag"
	"fmt"
	"net/url"
)

type Config struct {
	InitMode          *bool
	VolumeDirs        volumeDirsFlag
	VolumeDirsArchive volumeDirsArchiveFlag
	DirForUnarchive   *string
	Webhook           Webhook
}

type Webhook struct {
	Urls       urlsFlag
	Method     *string
	StatusCode *int
	Retries    *int
}

func New() (*Config, error) {
	cfg := &Config{}
	cfg.InitMode = flag.Bool("init-mode", false, "Init mode for unarchive files. Works only if volume-dir-archive exist. Default - false")
	cfg.DirForUnarchive = flag.String("dir-for-unarchive", "/tmp/unatchive", "Directory where the archives will be unpacked")
	cfg.Webhook.Method = flag.String("webhook-method", "POST", "the HTTP method url to use to send the webhook")
	cfg.Webhook.StatusCode = flag.Int("webhook-status-code", 200, "the HTTP status code indicating successful triggering of reload")
	cfg.Webhook.Retries = flag.Int("webhook-retries", 1, "the amount of times to retry the webhook reload request")

	flag.Var(&cfg.VolumeDirs, "volume-dir", "the config map volume directory to watch for updates; may be used multiple times")
	flag.Var(&cfg.VolumeDirsArchive, "volume-dir-archive", "the config map volume directory to watch for updates and unarchiving; may be used multiple times")
	flag.Var(&cfg.Webhook.Urls, "webhook-url", "the url to send a request to when the specified config map volume directory has been updated")
	flag.Parse()

	return cfg, nil
}

type volumeDirsFlag []string
type volumeDirsArchiveFlag []string
type urlsFlag []*url.URL

func (v *volumeDirsFlag) Set(value string) error {
	*v = append(*v, value)
	return nil
}
func (v *volumeDirsFlag) String() string {
	return fmt.Sprint(*v)
}

func (v *volumeDirsArchiveFlag) Set(value string) error {
	*v = append(*v, value)
	return nil
}
func (v *volumeDirsArchiveFlag) String() string {
	return fmt.Sprint(*v)
}

func (v *urlsFlag) Set(value string) error {
	u, err := url.Parse(value)
	if err != nil {
		return fmt.Errorf("invalid URL: %v", err)
	}
	*v = append(*v, u)
	return nil
}

func (v *urlsFlag) String() string {
	return fmt.Sprint(*v)
}
