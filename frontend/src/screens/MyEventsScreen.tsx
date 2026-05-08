import { useEffect, useState } from 'react';
import { EventCard } from '../components/EventCard';
import { searchEvents } from '../api/events';
import type { EventResponse } from '../types';
import { ApiError } from '../api/client';
import { useAuth } from '../auth/AuthContext';

interface MyEventsScreenProps {
  onOpenEvent(id: number): void;
  onCreateEvent(): void;
}

export function MyEventsScreen({ onOpenEvent, onCreateEvent }: MyEventsScreenProps) {
  const { user } = useAuth();
  const [events, setEvents] = useState<EventResponse[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  // The backend search does not accept a `creator_id` filter, so we pull every
  // event and filter on the client. This is a known limitation noted in the
  // README — see "Missing endpoints".
  useEffect(() => {
    let cancelled = false;
    async function load() {
      if (!user) return;
      setLoading(true);
      setError(null);
      try {
        const all = await searchEvents();
        if (cancelled) return;
        setEvents(all.filter((ev) => ev.creator_id === user.id && !ev.deleted));
      } catch (err) {
        if (cancelled) return;
        setError(err instanceof ApiError ? err.message : String(err));
      } finally {
        if (!cancelled) setLoading(false);
      }
    }
    load();
    return () => {
      cancelled = true;
    };
  }, [user]);

  return (
    <div className="container">
      <div className="section-head">
        <h2>мои события</h2>
        <button type="button" className="btn btn-sm btn-accent" onClick={onCreateEvent}>
          + создать
        </button>
      </div>

      {error && <div className="alert">{error}</div>}
      {loading ? (
        <div className="empty-state">загружаем…</div>
      ) : events.length === 0 ? (
        <div className="empty-state">вы ещё не создавали событий</div>
      ) : (
        events.map((ev) => <EventCard key={ev.id} event={ev} onOpen={onOpenEvent} />)
      )}
    </div>
  );
}
