package upnp

import (
	"context"
	"log"

	"github.com/ItsNotGoodName/go-upnpsub"
)

type BackgroundControlPoint struct {
	controlPoint upnpsub.ControlPoint
}

func NewBackgroundControlPoint(controlPoint upnpsub.ControlPoint) *BackgroundControlPoint {
	return &BackgroundControlPoint{
		controlPoint: controlPoint,
	}
}

func (bcp *BackgroundControlPoint) Background(ctx context.Context, doneC chan<- struct{}) {
	go func() {
		if err := upnpsub.ListenAndServe("", bcp.controlPoint); err != nil {
			log.Fatalln("Failed to start control point:", err)
		}
	}()
	doneC <- struct{}{}
}
