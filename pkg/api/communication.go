package api

import (
	"context"

	"github.com/greenboxal/aip/pkg/collective/comms"
	"github.com/greenboxal/aip/pkg/daemon"
)

type CommunicationAPI struct {
	daemon  *daemon.Daemon
	manager *comms.Manager
}

type CreateRoomRequest struct {
	RoomName string `json:"room_name"`
}

type CreateRoomResponse struct {
	Room *Room `json:"room"`
}

type JoinRoomRequest struct {
	RoomName string `json:"room_name"`
	PortName string `json:"port_name"`
}

type JoinRoomResponse struct {
	Room *Room `json:"room"`
}

type LeaveRoomRequest struct {
	RoomName string `json:"room_name"`
	PortName string `json:"port_name"`
}

type LeaveRoomResponse struct {
	Room *Room `json:"room"`
}

type GetRoomRequest struct {
	RoomName string `json:"room_name"`
}

type GetRoomResponse struct {
	Room *Room `json:"room"`
}

type Room struct {
	Name    string   `json:"name"`
	Members []string `json:"members"`
}

func NewCommunicationApi(d *daemon.Daemon, m *comms.Manager) *CommunicationAPI {
	return &CommunicationAPI{daemon: d, manager: m}
}

func buildRoom(room *comms.Room) *Room {
	return &Room{
		Name:    room.Name(),
		Members: room.Members(),
	}
}

func (api *CommunicationAPI) CreateRoom(ctx context.Context, req *CreateRoomRequest) (*CreateRoomResponse, error) {
	room := api.manager.CreateRoom(req.RoomName)

	return &CreateRoomResponse{Room: buildRoom(room)}, nil
}

func (api *CommunicationAPI) GetRoom(ctx context.Context, req *GetRoomRequest) (*GetRoomResponse, error) {
	room := api.manager.GetRoom(req.RoomName)

	return &GetRoomResponse{Room: buildRoom(room)}, nil
}

func (api *CommunicationAPI) JoinRoom(ctx context.Context, req *JoinRoomRequest) (*JoinRoomResponse, error) {
	room := api.manager.JoinRoom(req.RoomName, req.PortName)

	return &JoinRoomResponse{Room: buildRoom(room)}, nil
}

func (api *CommunicationAPI) LeaveRoom(ctx context.Context, req *LeaveRoomRequest) (*LeaveRoomResponse, error) {
	room := api.manager.LeaveRoom(req.RoomName, req.PortName)

	return &LeaveRoomResponse{Room: buildRoom(room)}, nil
}
