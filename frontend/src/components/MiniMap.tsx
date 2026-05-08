interface MiniMapProps {
  longitude: number;
  latitude: number;
}

export function MiniMap({ longitude, latitude }: MiniMapProps) {
  return (
    <div className="map" aria-hidden="true">
      <span className="pin" />
      <span className="map-coords">
        {latitude.toFixed(4)}, {longitude.toFixed(4)}
      </span>
    </div>
  );
}
