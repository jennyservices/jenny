package debug

import (
	"context"
	"log"

	"github.com/Typeform/jenny/debug/transport/v1"
	"sevki.org/lib/prettyprint"
)

func (ds *debugService) NewError(ctx context.Context, ID string, Packet v1.Packet) (Body *v1.ErrorResponse, err error) {
	log.Println(prettyprint.AsJSON(Packet))
	return &v1.ErrorResponse{Wrote: true}, nil
}
