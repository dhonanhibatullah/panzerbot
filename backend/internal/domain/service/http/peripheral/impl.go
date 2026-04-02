package domainservicehttpperipheral

import (
	"context"

	domainmodel "github.com/dhonanhibatullah/panzerbot/backend/internal/domain/model"
	portsinhttp "github.com/dhonanhibatullah/panzerbot/backend/internal/ports/in/http"
	portsoutlogging "github.com/dhonanhibatullah/panzerbot/backend/internal/ports/out/logging"
	portsoutperipheral "github.com/dhonanhibatullah/panzerbot/backend/internal/ports/out/peripheral"
)

const path = "domain/service/http/peripheral"

type peripheral struct {
	log        portsoutlogging.Log
	motorRight portsoutperipheral.Motor
	motorLeft  portsoutperipheral.Motor
	servoPan   portsoutperipheral.Servo
	servoTilt  portsoutperipheral.Servo
	soundboard portsoutperipheral.Soundboard
}

func New(
	log portsoutlogging.Log,
	motorRight portsoutperipheral.Motor,
	motorLeft portsoutperipheral.Motor,
	servoPan portsoutperipheral.Servo,
	servoTilt portsoutperipheral.Servo,
	soundboard portsoutperipheral.Soundboard,
) portsinhttp.Peripheral {
	return &peripheral{
		log:        log,
		motorRight: motorRight,
		motorLeft:  motorLeft,
		servoPan:   servoPan,
		servoTilt:  servoTilt,
		soundboard: soundboard,
	}
}

func (p *peripheral) SetMotorSpeedScale(ctx context.Context, right float64, left float64) (err error) {
	const tag = path + "/SetMotorSpeedScale"

	err = p.motorRight.SetSpeedScale(ctx, right)
	if err != nil {
		p.log.Error(
			ctx, tag,
			"Failed to set right motor speed",
			domainmodel.LogMeta{
				"error": err.Error(),
			},
		)
		return
	}

	err = p.motorLeft.SetSpeedScale(ctx, left)
	if err != nil {
		p.log.Error(
			ctx, tag,
			"Failed to set left motor speed",
			domainmodel.LogMeta{
				"error": err.Error(),
			},
		)
		return
	}

	p.log.Info(
		ctx, tag,
		"Motor speed scale set successfully",
		domainmodel.LogMeta{
			"right": right,
			"left":  left,
		},
	)

	return
}

func (p *peripheral) SetServoAngle(ctx context.Context, pan float64, tilt float64) (err error) {
	const tag = path + "/SetServoAngle"

	err = p.servoPan.SetAngle(ctx, pan)
	if err != nil {
		p.log.Error(
			ctx, tag,
			"Failed to set pan servo angle",
			domainmodel.LogMeta{
				"error": err.Error(),
			},
		)
		return
	}

	err = p.servoTilt.SetAngle(ctx, tilt)
	if err != nil {
		p.log.Error(
			ctx, tag,
			"Failed to set tilt servo angle",
			domainmodel.LogMeta{
				"error": err.Error(),
			},
		)
		return
	}

	p.log.Info(
		ctx, tag,
		"Servo angle set successfully",
		domainmodel.LogMeta{
			"pan":  pan,
			"tilt": tilt,
		},
	)

	return
}

func (p *peripheral) PlaySoundboardTrack(ctx context.Context, idx int) (err error) {
	const tag = path + "/PlayTrack"

	err = p.soundboard.PlayTrack(ctx, idx)
	if err != nil {
		p.log.Error(
			ctx, tag,
			"Failed to play track",
			domainmodel.LogMeta{
				"error": err.Error(),
			},
		)
		return
	}

	p.log.Info(
		ctx, tag,
		"Track played successfully",
		domainmodel.LogMeta{
			"index": idx,
		},
	)

	return
}

func (p *peripheral) StopSoundboardTracks(ctx context.Context) (err error) {
	const tag = path + "/StopTracks"

	err = p.soundboard.StopTracks(ctx)
	if err != nil {
		p.log.Error(
			ctx, tag,
			"Failed to stop tracks",
			domainmodel.LogMeta{
				"error": err.Error(),
			},
		)
		return
	}

	p.log.Info(
		ctx, tag,
		"Track stopped successfully",
		nil,
	)

	return
}

func (p *peripheral) GetSoundboardTrackNames(ctx context.Context) (trackNames []string, err error) {
	const tag = path + "/GetTrackNames"

	trackNames, err = p.soundboard.GetTrackNames(ctx)
	if err != nil {
		p.log.Error(
			ctx, tag,
			"Failed to get track names",
			domainmodel.LogMeta{
				"error": err.Error(),
			},
		)
	}

	p.log.Info(
		ctx, tag,
		"Track names received",
		domainmodel.LogMeta{
			"count":       len(trackNames),
			"track_names": trackNames,
		},
	)

	return
}

func (p *peripheral) GetSoundboardTrackFileNames(ctx context.Context) (trackFileNames []string, err error) {
	const tag = path + "/GetTrackFileNames"

	trackFileNames, err = p.soundboard.GetTrackFileNames(ctx)
	if err != nil {
		p.log.Error(
			ctx, tag,
			"Failed to get track file names",
			domainmodel.LogMeta{
				"error": err.Error(),
			},
		)
	}

	p.log.Info(
		ctx, tag,
		"Track names received",
		domainmodel.LogMeta{
			"count":       len(trackFileNames),
			"track_names": trackFileNames,
		},
	)

	return
}
