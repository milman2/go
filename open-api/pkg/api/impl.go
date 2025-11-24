package api

import (
	"context"

	"github.com/oklog/ulid/v2"
)

type Server struct{}

func NewServer() Server {
	return Server{}
}

func (Server) GetDevice(ctx context.Context, request GetDeviceRequestObject) (GetDeviceResponseObject, error) {
	d := Device{
		Id:   ulid.Make(),
		Name: "Test Device",
	}
	return GetDevice200JSONResponse(d), nil
}
