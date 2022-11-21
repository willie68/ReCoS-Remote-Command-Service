// +build !windows
package osdependent

import "wkla.no-ip.biz/remote-desk-service/internal/config"

func InitOSDependend(config config.Config) error {
	return nil
}

func DisposeOSDependend() {
}
