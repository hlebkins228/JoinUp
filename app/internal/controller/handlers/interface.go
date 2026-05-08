package handlers

import (
	"JoinUp/internal/controller/dto"
	"context"
)

type EventService interface {
	CreateEvent(ctx context.Context, req *dto.CreateEventRequest, userID int) (int, error)
	GetEvent(ctx context.Context, id int) (*dto.EventResponse, error)
	UpdateEvent(ctx context.Context, req *dto.UpdateEventRequest) error
	DeleteEvent(ctx context.Context, id int) error
	JoinEvent(ctx context.Context, eventID int, userID int) error
	UploadEventImage(ctx context.Context, eventID int, name string, data []byte) (int, error)
	AddEventCategory(ctx context.Context, req *dto.AddEventCategoryRequest) error
}

type UserService interface {
	CreateUser(ctx context.Context, req *dto.CreateUserRequest) (int, error)
	GetUser(ctx context.Context, id int) (*dto.UserResponse, error)
	UpdateUser(ctx context.Context, id int, req *dto.UpdateUserRequest) error
	UploadImage(ctx context.Context, name string, data []byte) (int, error)
	Auth(ctx context.Context, login, password string) (string, error)
	CheckPerms(ctx context.Context, eventID int, userID int, role string) (bool, error)
}

type SearchEngine interface {
	SearchEvents(ctx context.Context, req *dto.EventSearchRequest) (*dto.EventsResponse, error)
}
