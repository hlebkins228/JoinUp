import { useEffect, useState } from 'react';
import { EventCard } from '../components/EventCard';
import { searchEvents } from '../api/events';
import type { EventResponse } from '../types';
import { ApiError } from '../api/client';
import { useAuth } from '../auth/AuthContext';

interface JoiningScreenProps {
  onOpenEvent(id: number): void;
}

export function JoiningScreen({ onOpenEvent }: JoiningScreenProps) {
  const { user } = useAuth();
  const [events, setEvents] = useState<EventResponse[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [warning, setWarning] = useState<string | null>(null);

  useEffect(() => {
    let cancelled = false;
    async function load() {
      if (!user) return;
      setLoading(true);
      setError(null);
      setWarning(null);
      try {
        const all = await searchEvents();
        if (cancelled) return;

        // The EventResponse `members` field is declared as []string in the
        // backend DTO but is not populated by the current repository queries
        // (no JOIN with the `member` table). Until that is fixed we can only
        // best-effort filter by login or numeric id present in the field.
        const haystackKeys: string[] = [user.login, String(user.id)];
        const joined = all.filter(
          (ev) =>
            !ev.deleted &&
            ev.members?.some((m) => haystackKeys.includes(m)),
        );
        setEvents(joined);
        if (joined.length === 0 && all.length > 0) {
          setWarning(
            'API не возвращает список участников события — фильтр «участвую» работать не будет, пока бэкенд не вернёт members.',
          );
        }
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
        <h2>я участвую</h2>
      </div>

      {error && <div className="alert">{error}</div>}
      {warning && <div className="alert" style={{ borderColor: 'var(--fg-30)', color: 'var(--fg-50)', background: 'transparent' }}>{warning}</div>}
      {loading ? (
        <div className="empty-state">загружаем…</div>
      ) : events.length === 0 ? (
        <div className="empty-state">вы ещё ни к чему не присоединились</div>
      ) : (
        events.map((ev) => <EventCard key={ev.id} event={ev} onOpen={onOpenEvent} />)
      )}
    </div>
  );
}
