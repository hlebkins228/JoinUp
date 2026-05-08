import type { EventResponse } from '../types';
import { Avatar } from './Avatar';
import { fmtDateTime } from '../lib/format';
import { imageUrl } from '../lib/imageUrl';

interface EventCardProps {
  event: EventResponse;
  onOpen(id: number): void;
}

export function EventCard({ event, onOpen }: EventCardProps) {
  const cover = imageUrl(event.image_id);
  return (
    <button
      type="button"
      className="event-card"
      onClick={() => onOpen(event.id)}
      aria-label={event.name}
    >
      <Avatar src={cover} name={event.name} size="event" />
      <div className="event-card-meta">
        <h3 className="event-card-title">{event.name}</h3>
        <div className="event-card-line">
          <span>{fmtDateTime(event.event_time)}</span>
          <span>·</span>
          <span>{event.city}</span>
          {event.location?.name ? (
            <>
              <span>·</span>
              <span>{event.location.name}</span>
            </>
          ) : null}
        </div>
      </div>
      <div className="event-card-side">
        {event.members?.length ? (
          <span>
            {event.members.length} участн{event.members.length === 1 ? 'ик' : 'иков'}
          </span>
        ) : (
          <span className="muted">создатель #{event.creator_id}</span>
        )}
      </div>
    </button>
  );
}
