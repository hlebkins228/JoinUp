// The backend stores images in the database but does not currently expose
// a `GET /image/:id` endpoint. We keep this helper as the single place that
// would change once such a route lands. For now any request to render an
// image by id falls back to `null`, which forces the UI to show initials /
// the empty-state placeholder.
export function imageUrl(_id: number | null | undefined): string | null {
  return null;
}
