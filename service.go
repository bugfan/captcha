package captcha

import (
	"github.com/bugfan/captcha/service"
)

func Start(addr string) error {
	return (service.NewServer(addr).Run())
}
