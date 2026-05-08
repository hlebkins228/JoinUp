import { useState } from 'react';
import { AuthProvider, useAuth } from './auth/AuthContext';
import { Header } from './components/Header';
import { LoginScreen } from './screens/LoginScreen';
import { RegisterScreen } from './screens/RegisterScreen';
import { SearchScreen } from './screens/SearchScreen';
import { MyEventsScreen } from './screens/MyEventsScreen';
import { JoiningScreen } from './screens/JoiningScreen';
import { ProfileScreen } from './screens/ProfileScreen';
import { EventDetailScreen } from './screens/EventDetailScreen';
import type { Route, RouteState } from './routes';

function Shell() {
  const { token, loading, user } = useAuth();
  const [authRoute, setAuthRoute] = useState<'login' | 'register'>('login');
  const [state, setState] = useState<RouteState>({ route: 'search' });

  if (loading) {
    return <div className="empty-state">загружаем…</div>;
  }

  if (!token || !user) {
    if (authRoute === 'register') {
      return <RegisterScreen onGoLogin={() => setAuthRoute('login')} />;
    }
    return <LoginScreen onGoRegister={() => setAuthRoute('register')} />;
  }

  function navigate(route: Route, eventId: number | null = null) {
    setState({ route, eventId });
  }

  let body: JSX.Element;
  switch (state.route) {
    case 'search':
      body = (
        <SearchScreen
          onOpenEvent={(id) => navigate('event-detail', id)}
          onCreateEvent={() => navigate('event-new', null)}
        />
      );
      break;
    case 'mine':
      body = (
        <MyEventsScreen
          onOpenEvent={(id) => navigate('event-detail', id)}
          onCreateEvent={() => navigate('event-new', null)}
        />
      );
      break;
    case 'joining':
      body = (
        <JoiningScreen
          onOpenEvent={(id) => navigate('event-detail', id)}
        />
      );
      break;
    case 'profile':
      body = <ProfileScreen onBack={() => navigate('search')} />;
      break;
    case 'event-detail':
      body = (
        <EventDetailScreen
          eventId={state.eventId ?? null}
          onBack={() => navigate('search')}
          onAfterDelete={() => navigate('mine')}
        />
      );
      break;
    case 'event-new':
      body = (
        <EventDetailScreen
          eventId={null}
          onBack={() => navigate('search')}
          onAfterDelete={() => navigate('mine')}
        />
      );
      break;
    default:
      body = (
        <SearchScreen
          onOpenEvent={(id) => navigate('event-detail', id)}
          onCreateEvent={() => navigate('event-new', null)}
        />
      );
  }

  return (
    <div className="app-shell">
      <Header
        route={state.route}
        user={user}
        onNavigate={(r) => navigate(r)}
      />
      <main style={{ flex: 1, paddingBottom: 64 }}>{body}</main>
    </div>
  );
}

export default function App() {
  return (
    <AuthProvider>
      <Shell />
    </AuthProvider>
  );
}
