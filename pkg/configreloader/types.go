package configreloader

import (
	"fmt"
	"net/url"
)

type ConfigReloader struct {
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
