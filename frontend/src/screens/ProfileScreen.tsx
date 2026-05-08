import { useEffect, useState, type FormEvent } from 'react';
import { Field } from '../components/Field';
import { NumberInput, Select, TextInput } from '../components/inputs';
import { UploadAvatar } from '../components/UploadAvatar';
import { CITIES } from '../lib/constants';
import { updateUser } from '../api/users';
import { ApiError } from '../api/client';
import { useAuth } from '../auth/AuthContext';

interface ProfileScreenProps {
  onBack(): void;
}

export function ProfileScreen({ onBack }: ProfileScreenProps) {
  const { user, refreshUser, signOut } = useAuth();
  const [name, setName] = useState('');
  const [age, setAge] = useState<number | ''>('');
  const [city, setCity] = useState('');
  const [tg, setTg] = useState('');
  const [avatarId, setAvatarId] = useState<number | null>(null);
  const [avatarPreview, setAvatarPreview] = useState<string | null>(null);
  const [busy, setBusy] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [success, setSuccess] = useState(false);

  useEffect(() => {
    if (!user) return;
    setName(user.name);
    setAge(user.age);
    setCity(user.city);
    setTg(user.telegram_login ?? '');
    setAvatarId(user.avatar_id ?? null);
    setAvatarPreview(null);
  }, [user]);

  if (!user) return null;

  async function onSubmit(e: FormEvent) {
    e.preventDefault();
    setError(null);
    setSuccess(false);
    if (!name.trim() || age === '' || !city) {
      setError('заполните обязательные поля');
      return;
    }
    setBusy(true);
    try {
      await updateUser({
        name: name.trim(),
        age: Number(age),
        city,
        telegram_login: tg.trim() || undefined,
        avatar_id: avatarId ?? undefined,
      });
      await refreshUser();
      setSuccess(true);
    } catch (err) {
      setError(err instanceof ApiError ? err.message : String(err));
    } finally {
      setBusy(false);
    }
  }

  const cityOptions = CITIES.map((c) => ({ value: c, label: c }));

  return (
    <div className="container">
      <div className="section-head">
        <h2>профиль</h2>
        <button type="button" className="btn btn-sm btn-ghost" onClick={onBack}>
          назад
        </button>
      </div>

      <div className="profile-grid">
        <div>
          <UploadAvatar
            src={avatarPreview}
            name={user.name}
            onUploaded={(id, preview) => {
              setAvatarId(id);
              setAvatarPreview(preview);
            }}
            onError={(message) => setError(message)}
          />
        </div>

        <form onSubmit={onSubmit} noValidate>
          {error && <div className="alert">{error}</div>}
          {success && (
            <div
              className="alert"
              style={{
                borderColor: 'var(--fg-30)',
                color: 'var(--fg)',
                background: 'transparent',
              }}
            >
              сохранено
            </div>
          )}

          <Field label="имя">
            <TextInput value={name} onChange={setName} />
          </Field>
          <Field label="логин" hint="логин нельзя поменять">
            <TextInput
              value={user.login}
              onChange={() => undefined}
              disabled
            />
          </Field>
          <Field label="возраст">
            <NumberInput value={age} onChange={setAge} min={1} max={120} />
          </Field>
          <Field label="город">
            <Select value={city} onChange={setCity} options={cityOptions} />
          </Field>
          <Field label="telegram" optional>
            <TextInput value={tg} onChange={setTg} placeholder="anna_k" />
          </Field>

          <div className="row-end" style={{ marginTop: 16 }}>
            <button
              type="button"
              className="btn btn-danger btn-sm"
              onClick={signOut}
            >
              выйти
            </button>
            <button type="submit" className="btn" disabled={busy}>
              {busy ? 'сохраняем…' : 'сохранить'}
            </button>
          </div>
        </form>
      </div>
    </div>
  );
}
