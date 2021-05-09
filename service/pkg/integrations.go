package pkg

import (
	"wkla.no-ip.biz/remote-desk-service/pkg/hardware"
	"wkla.no-ip.biz/remote-desk-service/pkg/lighting"
	"wkla.no-ip.biz/remote-desk-service/pkg/models"
	"wkla.no-ip.biz/remote-desk-service/pkg/smarthome"
)

var IntegInfos = []models.IntegInfo{
	hardware.OpenHardwareMonitorIntegInfo,
	lighting.PhilipsHueIntegInfo,
	smarthome.HomematicIntegInfo,
}
