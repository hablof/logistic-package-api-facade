package messagewriter

import (
	"encoding/json"
	"fmt"

	"github.com/hablof/logistic-package-api-facade/internal/model"
	pb "github.com/hablof/logistic-package-api/pkg/kafka-proto"
	"github.com/rs/zerolog/log"
	"google.golang.org/protobuf/proto"
)

type Writer struct {
	channel <-chan []byte
}

func NewWriter(channel <-chan []byte) *Writer {
	return &Writer{
		channel: channel,
	}
}

func (w *Writer) Start() {
	go func() {
		for msg := range w.channel {
			scanUnit := pb.PackageEvent{}
			if err := proto.Unmarshal(msg, &scanUnit); err != nil {
				log.Err(err).Msg("unable to unmarshal protobuf message")
			}

			unit := model.PackageEvent{
				ID:        scanUnit.ID,
				PackageID: scanUnit.PackageID,
				Type:      0, // rewrite by switch statement below
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
			fmt.Println()
		}
	}()
}
