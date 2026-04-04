package domainappcore

import (
	"context"
	"math"

	adaptersoutlogginggeneric "github.com/dhonanhibatullah/panzerbot/backend/internal/adapters/out/logging/generic"
	adaptersoutperipheralmotor "github.com/dhonanhibatullah/panzerbot/backend/internal/adapters/out/peripheral/motor"
	adaptersoutperipheralrtcmedia "github.com/dhonanhibatullah/panzerbot/backend/internal/adapters/out/peripheral/rtc_media"
	adaptersoutperipheralservo "github.com/dhonanhibatullah/panzerbot/backend/internal/adapters/out/peripheral/servo"
	adaptersoutperipheralsoundboard "github.com/dhonanhibatullah/panzerbot/backend/internal/adapters/out/peripheral/soundboard"
	"github.com/dhonanhibatullah/panzerbot/backend/internal/config"
	domainmodel "github.com/dhonanhibatullah/panzerbot/backend/internal/domain/model"
	domainservicehttpperipheral "github.com/dhonanhibatullah/panzerbot/backend/internal/domain/service/http/peripheral"
	domainservicehttprtcsignalling "github.com/dhonanhibatullah/panzerbot/backend/internal/domain/service/http/rtc_signalling"
	portsinhttp "github.com/dhonanhibatullah/panzerbot/backend/internal/ports/in/http"
	portsoutlogging "github.com/dhonanhibatullah/panzerbot/backend/internal/ports/out/logging"
	portsoutperipheral "github.com/dhonanhibatullah/panzerbot/backend/internal/ports/out/peripheral"
)

type Wiring struct {
	logPort              portsoutlogging.Log
	motorRightPort       portsoutperipheral.Motor
	motorLeftPort        portsoutperipheral.Motor
	servoPanPort         portsoutperipheral.Servo
	servoTiltPort        portsoutperipheral.Servo
	soundboardPort       portsoutperipheral.Soundboard
	rtcMediaPort         portsoutperipheral.RTCMedia
	peripheralSvcPort    portsinhttp.Peripheral
	rtcSignallingSvcPort portsinhttp.RTCSignalling
}

func (c *Core) NewWiring(ctx context.Context) (err error) {
	const tag = path + "/NewWiring"

	logPort := adaptersoutlogginggeneric.New(
		c.infrastructure.logger,
	)

	motorRightPort := adaptersoutperipheralmotor.NewNoPwm(
		logPort,
		&config.MotorRightAPin,
		&config.MotorRightBPin,
	)
	motorLeftPort := adaptersoutperipheralmotor.NewNoPwm(
		logPort,
		&config.MotorLeftAPin,
		&config.MotorLeftBPin,
	)
	servoPanPort := adaptersoutperipheralservo.NewPwm(
		logPort,
		&config.ServoPanPin,
	)
	servoTiltPort := adaptersoutperipheralservo.NewPwm(
		logPort,
		&config.ServoTiltPin,
	)
	servoPanPort.SetAngle(ctx, math.Pi/2.0)
	servoTiltPort.SetAngle(ctx, math.Pi/2.0)

	speakerSampleRate := c.infrastructure.soundTracks[0].Format.SampleRate

	soundboardPort := adaptersoutperipheralsoundboard.New(
		logPort,
		c.infrastructure.soundTracks,
		speakerSampleRate,
	)
	rtcMediaPort, err := adaptersoutperipheralrtcmedia.New(
		logPort,
		speakerSampleRate,
	)
	if err != nil {
		logPort.Error(
			ctx, tag,
			"Failed to initiate RTC Media Port",
			domainmodel.LogMeta{
				"error": err.Error(),
			},
		)
		return err
	}

	peripheralSvcPort := domainservicehttpperipheral.New(
		logPort,
		motorRightPort,
		motorLeftPort,
		servoPanPort,
		servoTiltPort,
		soundboardPort,
	)
	rtcSignallingSvcPort := domainservicehttprtcsignalling.New(
		logPort,
		rtcMediaPort,
	)

	c.wiring = &Wiring{
		logPort:              logPort,
		motorRightPort:       motorRightPort,
		motorLeftPort:        motorLeftPort,
		servoPanPort:         servoPanPort,
		servoTiltPort:        servoTiltPort,
		soundboardPort:       soundboardPort,
		rtcMediaPort:         rtcMediaPort,
		peripheralSvcPort:    peripheralSvcPort,
		rtcSignallingSvcPort: rtcSignallingSvcPort,
	}
	return nil
}
