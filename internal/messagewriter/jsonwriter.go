package messagewriter

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
)

type JsonWriter struct {
	channel <-chan []byte
}

func NewJsonWriter(channel <-chan []byte) *JsonWriter {
	return &JsonWriter{
		channel: channel,
	}
}

func (w *JsonWriter) Start(ctx context.Context) {
	go func() {
		for {
			select {
			case msg := <-w.channel:

				out := bytes.Buffer{}
				json.Indent(&out, msg, "", "    ")
				fmt.Println(out.String())
				// fmt.Println()

			case <-ctx.Done():
				return
			}
		}
	}()
}
