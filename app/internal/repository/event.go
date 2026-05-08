package repository

import (
	"JoinUp/internal/exceptions"
	"JoinUp/internal/models"
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type EventRepo struct {
	Pool
}

func NewEventRepo(pool *pgxpool.Pool) EventRepo {
	return EventRepo{Pool: Pool{pool: pool}}
}

func (r *EventRepo) CreateEvent(ctx context.Context, event *models.Event) (int, error) {
	run, err := r.runFromCtx(ctx)
	if err != nil {
		return 0, err
	}

	const locationInsertQuery = `
		insert into location (name, longitude, latitude, address)
		values ($1, $2, $3, $4)
		returning id
	`

	const eventInsertQuery = `
		insert into event (name, description, created_at, updated_at, event_date, telegram_chat_url, city, creator_id, location_id, image_id)
		values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		returning id
	`

	var locationID int
	err = run.QueryRow(ctx, locationInsertQuery, event.Location.Name,
		event.Location.Longitude, event.Location.Latitude, event.Location.Address).Scan(&locationID)
	if err != nil {
		return 0, err
	}

	var eventID int
	err = run.QueryRow(ctx, eventInsertQuery, event.Name, event.Desc,
		event.CreatedAt, event.UpdatedAt, event.EventTime, event.TelegramChatURL,
		event.City, event.CreatorID, locationID, event.ImageID).Scan(&eventID)
	if err != nil {
		return 0, err
	}

	return eventID, nil
}

func (r *EventRepo) GetEvent(ctx context.Context, id int) (*models.Event, error) {
	run, err := r.runFromCtx(ctx)
	if err != nil {
		return nil, err
	}

	const query = `
	select e.id, e.name, e.description, e.created_at, e.updated_at, e.event_date,
	e.telegram_chat_url, e.city, e.creator_id, e.image_id, l.name, l.longitude, l.latitude, l.address, e.deleted
	from event e
	join
	location l
	on e.location_id = l.id
	where e.id = $1`

	event := models.Event{}

	err = run.QueryRow(ctx, query, id).Scan(
		&event.ID,
		&event.Name,
		&event.Desc,
		&event.CreatedAt,
		&event.UpdatedAt,
		&event.EventTime,
		&event.TelegramChatURL,
		&event.City,
		&event.CreatorID,
		&event.ImageID,
		&event.Location.Name,
		&event.Location.Longitude,
		&event.Location.Latitude,
		&event.Location.Address,
		&event.Deleted,
	)
	if err != nil {
		return nil, err
	}

	return &event, nil
}

func (r *EventRepo) UpdateEvent(ctx context.Context, event *models.Event) error {
	run, err := r.runFromCtx(ctx)
	if err != nil {
		return err
	}

	const locationIDQuery = `
		select location_id
		from event
		where id = $1 and deleted = false
	`

	const locationUpdateQuery = `
		update location
		set name = $1, longitude = $2, latitude = $3, address = $4
		where id = $5
	`

	const eventUpdateQuery = `
		update event
		set name = $1,
		    description = $2,
		    updated_at = $3,
		    event_date = $4,
		    telegram_chat_url = $5,
		    city = $6,
		    image_id = $7
		where id = $8 and deleted = false
	`

	var locationID int
	if err = run.QueryRow(ctx, locationIDQuery, event.ID).Scan(&locationID); err != nil {
		return err
	}

	if _, err = run.Exec(ctx, locationUpdateQuery,
		event.Location.Name, event.Location.Longitude, event.Location.Latitude, event.Location.Address, locationID,
	); err != nil {
		return err
	}

	tag, err := run.Exec(ctx, eventUpdateQuery,
		event.Name, event.Desc, event.UpdatedAt, event.EventTime, event.TelegramChatURL,
		event.City, event.ImageID, event.ID,
	)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func (r *EventRepo) DeleteEvent(ctx context.Context, id int) error {
	run, err := r.runFromCtx(ctx)
	if err != nil {
		return err
	}

	var deleted bool
	err = run.QueryRow(ctx, `select deleted from event where id = $1`, id).Scan(&deleted)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return exceptions.ErrNotFound
		}
		return err
	}
	if deleted {
		return exceptions.ErrAlreadyDeleted
	}

	now := time.Now().UTC()
	const query = `
		update event
		set deleted = true, updated_at = $2
		where id = $1 and deleted = false
	`

	_, err = run.Exec(ctx, query, id, now)
	if err != nil {
		return err
	}

	return nil
}

func (r *EventRepo) JoinEvent(ctx context.Context, eventID int, userID int) error {
	run, err := r.runFromCtx(ctx)
	if err != nil {
		return err
	}

	var eventExists bool
	err = run.QueryRow(ctx, `select exists(select 1 from event where id = $1 and deleted = false)`, eventID).Scan(&eventExists)
	if err != nil {
		return err
	}
	if !eventExists {
		return exceptions.ErrNotFound
	}

	var alreadyMember bool
	err = run.QueryRow(ctx, `select exists(select 1 from member where event_id = $1 and user_id = $2)`, eventID, userID).Scan(&alreadyMember)
	if err != nil {
		return err
	}
	if alreadyMember {
		return exceptions.ErrAlreadyExists
	}

	const query = `
		insert into member (user_id, event_id, role)
		values ($1, $2, 'member')
	`

	_, err = run.Exec(ctx, query, userID, eventID)
	return err
}

func (r *EventRepo) UpdateEventImage(ctx context.Context, eventID int, imageID int, updatedAt time.Time) error {
	run, err := r.runFromCtx(ctx)
	if err != nil {
		return err
	}

	const query = `
		update event
		set image_id = $1, updated_at = $2
		where id = $3 and deleted = false
	`

	tag, err := run.Exec(ctx, query, imageID, updatedAt, eventID)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return exceptions.ErrNotFound
	}

	return nil
}

func (r *EventRepo) AddEventCategory(ctx context.Context, eventID int, categoryID int) error {
	run, err := r.runFromCtx(ctx)
	if err != nil {
		return err
	}

	var eventExists bool
	err = run.QueryRow(ctx, `select exists(select 1 from event where id = $1 and deleted = false)`, eventID).Scan(&eventExists)
	if err != nil {
		return err
	}
	if !eventExists {
		return exceptions.ErrNotFound
	}

	var categoryName string
	err = run.QueryRow(ctx, `select name from category where id = $1`, categoryID).Scan(&categoryName)
	if errors.Is(err, sql.ErrNoRows) {
		return exceptions.ErrNotExists
	}
	if err != nil {
		return err
	}

	var alreadyAdded bool
	err = run.QueryRow(ctx, `
		select exists(
			select 1
			from category
			where event_id = $1
			  and (id = $2 or subcategory_id = $2)
		)
	`, eventID, categoryID).Scan(&alreadyAdded)
	if err != nil {
		return err
	}
	if alreadyAdded {
		return exceptions.ErrAlreadyExists
	}

	const query = `
		insert into category (event_id, subcategory_id, name)
		values ($1, $2, $3)
	`

	_, err = run.Exec(ctx, query, eventID, categoryID, categoryName)
	return err
}

func (r *EventRepo) SearchEvents(ctx context.Context, filter models.EventSearchFilter) ([]*models.Event, error) {
	run, err := r.runFromCtx(ctx)
	if err != nil {
		return nil, err
	}

	categoryIDs := make([]int32, 0, len(filter.CategoryIDs))
	for _, id := range filter.CategoryIDs {
		categoryIDs = append(categoryIDs, int32(id))
	}

	const query = `
		select e.id, e.name, e.description, e.created_at, e.updated_at, e.event_date,
		       e.telegram_chat_url, e.city, e.creator_id, e.image_id,
		       l.name, l.longitude, l.latitude, l.address, e.deleted
		from event e
		join location l on e.location_id = l.id
		where e.deleted = false
		  and ($1::text is null or e.name ilike '%' || $1 || '%')
		  and ($2::timestamp is null or e.event_date >= $2)
		  and ($3::timestamp is null or e.event_date <= $3)
		  and ($4::text is null or e.city = $4)
		  and (
			cardinality($5::int[]) = 0
			or e.id in (
				select c.event_id
				from category c
				where c.id = any($5::int[]) or c.subcategory_id = any($5::int[])
				group by c.event_id
				having count(distinct coalesce(c.subcategory_id, c.id)) = cardinality($5::int[])
			)
		  )
		order by e.event_date, e.id
	`

	rows, err := run.Query(ctx, query,
		filter.Name,
		filter.EventFrom,
		filter.EventTo,
		filter.City,
		categoryIDs,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	events := make([]*models.Event, 0)
	for rows.Next() {
		event := models.Event{}
		err = rows.Scan(
			&event.ID,
			&event.Name,
			&event.Desc,
			&event.CreatedAt,
			&event.UpdatedAt,
			&event.EventTime,
			&event.TelegramChatURL,
			&event.City,
			&event.CreatorID,
			&event.ImageID,
			&event.Location.Name,
			&event.Location.Longitude,
			&event.Location.Latitude,
			&event.Location.Address,
			&event.Deleted,
		)
		if err != nil {
			return nil, err
		}
		events = append(events, &event)
	}
	if rows.Err() != nil {
		return nil, rows.Err()
	}

	return events, nil
}
