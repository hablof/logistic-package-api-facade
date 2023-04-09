package messagewriter

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hablof/logistic-package-api-facade/internal/model"
	pb "github.com/hablof/logistic-package-api/pkg/kafka-proto"
	"github.com/rs/zerolog/log"
	"google.golang.org/protobuf/proto"
)

type PbWriter struct {
	channel <-chan []byte
}

func NewPbWriter(channel <-chan []byte) *PbWriter {
	return &PbWriter{
		channel: channel,
	}
}

func (w *PbWriter) Start(ctx context.Context) {
	go func() {
		for {
			select {
			case msg := <-w.channel:
				handleMsg(msg)

			case <-ctx.Done():
				return
			}
		}
	}()
}

func handleMsg(msg []byte) {
	scanUnit := pb.PackageEvent{}
	if err := proto.Unmarshal(msg, &scanUnit); err != nil {
		log.Err(err).Msg("unable to unmarshal protobuf message")
	}

	unit := model.PackageEvent{
		ID:        scanUnit.ID,
		PackageID: scanUnit.PackageID,
		Type:      0,
		Created:   scanUnit.Created.AsTime(),
		Payload:   scanUnit.Payload,
	}

	switch scanUnit.Type {
	case pb.EventType_Created:
		unit.Type = model.Created

	case pb.EventType_Updated:
		unit.Type = model.Updated

	case pb.EventType_Removed:
		unit.Type = model.Removed
	}

	b, err := json.MarshalIndent(unit, "", "    ")
	if err != nil {
		log.Err(err).Msg("unable to marshal json")
	}
	fmt.Println(string(b))
	// fmt.Println()
}
