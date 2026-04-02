package adaptersinhttproutes

import (
	"net/http"

	adaptersinhttphandler "github.com/dhonanhibatullah/panzerbot/backend/internal/adapters/in/http/handler"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type Routes struct {
	router           *gin.RouterGroup
	peripheralHdl    *adaptersinhttphandler.Peripheral
	rtcSignallingHdl *adaptersinhttphandler.RTCSignalling
}

func New(
	router *gin.RouterGroup,
	peripheralHdl *adaptersinhttphandler.Peripheral,
	rtcSignallingHdl *adaptersinhttphandler.RTCSignalling,
) *Routes {
	return &Routes{
		router:           router,
		peripheralHdl:    peripheralHdl,
		rtcSignallingHdl: rtcSignallingHdl,
	}
}

func (r *Routes) Route() {
	docs := r.router.Group("/docs")

	docs.GET("/", func(c *gin.Context) { c.Redirect(http.StatusFound, "/docs/swagger/index.html") })
	docs.GET("/swagger", func(c *gin.Context) { c.Redirect(http.StatusFound, "/docs/swagger/index.html") })
	docs.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	v1 := r.router.Group("/v1")

	peripheral := v1.Group("/peripheral")
	peripheral.GET("/ws", r.peripheralHdl.Ws)
	peripheral.GET("/soundboard", r.peripheralHdl.GetSoundboardTrack)
	peripheral.POST("/soundboard/stop", r.peripheralHdl.PostSoundboardTrackStop)
	peripheral.POST("/soundboard/:track_idx", r.peripheralHdl.PostSoundboardTrack)

	rtc := v1.Group("/rtc")
	rtc.GET("/ws", r.rtcSignallingHdl.Ws)
}
