package pkg

import (
	"wkla.no-ip.biz/remote-desk-service/pkg/audio"
	"wkla.no-ip.biz/remote-desk-service/pkg/hardware"
	"wkla.no-ip.biz/remote-desk-service/pkg/lighting"
	"wkla.no-ip.biz/remote-desk-service/pkg/models"
	"wkla.no-ip.biz/remote-desk-service/pkg/smarthome"
	"wkla.no-ip.biz/remote-desk-service/pkg/video"
)

var IntegInfos = []models.IntegInfo{
	audio.AudioPlayerIntegInfo,
	hardware.OpenHardwareMonitorIntegInfo,
	lighting.PhilipsHueIntegInfo,
	smarthome.HomematicIntegInfo,
	hardware.StreamdeckIntegInfo,
	video.OBSIntegInfo,
}
