package service

import (
	"JoinUp/internal/controller/dto"
	"JoinUp/internal/exceptions"
	"JoinUp/internal/models"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"
)

type EventService struct {
	eventRepo EventRepo
	userRepo  UserRepo
	imageRepo ImageRepo
	uow       UOW
}

func NewEventService(uow UOW, eventRepo EventRepo, userRepo UserRepo, imageRepo ImageRepo) EventService {
	return EventService{uow: uow, eventRepo: eventRepo, userRepo: userRepo, imageRepo: imageRepo}
}

func (s *EventService) CreateEvent(ctx context.Context, req *dto.CreateEventRequest, userID int) (int, error) {
	event := req.ToModel()
	event.CreatedAt = time.Now().UTC()
	event.UpdatedAt = event.CreatedAt
	event.Deleted = false
	event.CreatorID = userID

	ctx, err := s.uow.BeginTx(ctx)
	if err != nil {
		return 0, err
	}
	defer s.uow.Rollback(ctx)

	ok, err := s.userRepo.CheckExists(ctx, event.CreatorID)
	if err != nil {
		return 0, err
	}
	if !ok {
		return 0, fmt.Errorf("%w: creator with such id not found, creator_id=%d", exceptions.ErrNotExists, event.CreatorID)
	}

	if event.ImageID != nil {
		ok, err := s.imageRepo.CheckExists(ctx, *event.ImageID)
		if err != nil {
			return 0, err
		}
		if !ok {
			return 0, fmt.Errorf("%w: image with such id not found, image_id=%d", exceptions.ErrNotExists, *event.ImageID)
		}
	}

	id, err := s.eventRepo.CreateEvent(ctx, &event)
	if err != nil {
		return 0, err
	}

	return id, s.uow.Commit(ctx)
}

func (s *EventService) GetEvent(ctx context.Context, id int) (*dto.EventResponse, error) {
	event, err := s.eventRepo.GetEvent(ctx, id)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("%w: no event with id = %d", exceptions.ErrNotFound, id)
	} else if err != nil {
		return nil, err
	}
	return modelToEventResponseDTO(event), nil
}

func (s *EventService) UpdateEvent(ctx context.Context, req *dto.UpdateEventRequest) error {
	event := req.ToModel()

	if event.ImageID != nil {
		ok, err := s.imageRepo.CheckExists(ctx, *event.ImageID)
		if err != nil {
			return err
		}
		if !ok {
			return fmt.Errorf("%w: image with such id not found, image_id=%d", exceptions.ErrNotExists, *event.ImageID)
		}
	}

	event.UpdatedAt = time.Now().UTC()

	if err := s.eventRepo.UpdateEvent(ctx, &event); errors.Is(err, sql.ErrNoRows) {
		return fmt.Errorf("%w: no event with id = %d", exceptions.ErrNotFound, event.ID)
	} else if err != nil {
		return err
	}

	return nil
}

func (s *EventService) DeleteEvent(ctx context.Context, id int) error {
	err := s.eventRepo.DeleteEvent(ctx, id)
	if errors.Is(err, exceptions.ErrNotFound) {
		return fmt.Errorf("%w: no event with id = %d", exceptions.ErrNotFound, id)
	} else if errors.Is(err, exceptions.ErrAlreadyDeleted) {
		return fmt.Errorf("%w: event with id = %d is already deleted", exceptions.ErrAlreadyDeleted, id)
	}
	return err
}

func (s *EventService) JoinEvent(ctx context.Context, eventID int, userID int) error {
	err := s.eventRepo.JoinEvent(ctx, eventID, userID)
	if errors.Is(err, exceptions.ErrNotFound) {
		return fmt.Errorf("%w: no event with id = %d", exceptions.ErrNotFound, eventID)
	}
	if errors.Is(err, exceptions.ErrAlreadyExists) {
		return fmt.Errorf("%w: user id = %d already participates in event id = %d", exceptions.ErrAlreadyExists, userID, eventID)
	}
	return err
}

func (s *EventService) UploadEventImage(ctx context.Context, eventID int, name string, data []byte) (int, error) {
	ctx, err := s.uow.BeginTx(ctx)
	if err != nil {
		return 0, err
	}
	defer s.uow.Rollback(ctx)

	imageID, err := s.imageRepo.AddImage(ctx, &models.Image{Name: name, Data: data})
	if err != nil {
		return 0, err
	}

	err = s.eventRepo.UpdateEventImage(ctx, eventID, imageID, time.Now().UTC())
	if errors.Is(err, exceptions.ErrNotFound) {
		return 0, fmt.Errorf("%w: no event with id = %d", exceptions.ErrNotFound, eventID)
	}
	if err != nil {
		return 0, err
	}

	return imageID, s.uow.Commit(ctx)
}

func (s *EventService) AddEventCategory(ctx context.Context, req *dto.AddEventCategoryRequest) error {
	err := s.eventRepo.AddEventCategory(ctx, req.EventID, req.CategoryID)
	if errors.Is(err, exceptions.ErrNotFound) {
		return fmt.Errorf("%w: no event with id = %d", exceptions.ErrNotFound, req.EventID)
	}
	if errors.Is(err, exceptions.ErrNotExists) {
		return fmt.Errorf("%w: category with id = %d not found", exceptions.ErrNotExists, req.CategoryID)
	}
	if errors.Is(err, exceptions.ErrAlreadyExists) {
		return fmt.Errorf("%w: category id = %d already added to event id = %d", exceptions.ErrAlreadyExists, req.CategoryID, req.EventID)
	}
	return err
}

func modelToEventResponseDTO(event *models.Event) *dto.EventResponse {
	return &dto.EventResponse{
		ID:              event.ID,
		CreatorID:       event.CreatorID,
		Name:            event.Name,
		Desc:            event.Desc,
		CreatedAt:       event.CreatedAt,
		UpdatedAt:       event.UpdatedAt,
		EventTime:       event.EventTime,
		TelegramChatURL: event.TelegramChatURL,
		City:            event.City,
		Members:         event.Members,
		Location:        modelLocationToDTO(event.Location),
		ImageID:         event.ImageID,
		Deleted:         event.Deleted,
	}
}

func modelLocationToDTO(loc models.Location) *dto.EventLocation {
	return &dto.EventLocation{
		Name:      loc.Name,
		Longitude: &loc.Longitude,
		Latitude:  &loc.Latitude,
		Address:   loc.Address,
	}
}
