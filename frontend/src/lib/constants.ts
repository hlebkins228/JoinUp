import type { Category } from '../types';

// Hard-coded city list — the backend stores cities as plain strings and does
// not expose a directory endpoint, so we replicate the prototype's options.
export const CITIES = [
  'Москва',
  'Санкт-Петербург',
  'Казань',
  'Новосибирск',
  'Екатеринбург',
  'Нижний Новгород',
  'Краснодар',
  'Самара',
  'Ростов-на-Дону',
  'Уфа',
  'Воронеж',
  'Пермь',
  'Волгоград',
  'Калининград',
  'Сочи',
];

// Hard-coded categories. The numeric ids are arbitrary — when the backend
// gains a categories directory we can swap this for an API call without
// touching the rest of the UI.
export const CATEGORIES: Category[] = [
  { id: 1, slug: 'sport', name: 'спорт', hue: 25 },
  { id: 2, slug: 'music', name: 'музыка', hue: 290 },
  { id: 3, slug: 'games', name: 'игры', hue: 220 },
  { id: 4, slug: 'art', name: 'арт', hue: 340 },
  { id: 5, slug: 'walks', name: 'прогулки', hue: 140 },
  { id: 6, slug: 'study', name: 'учёба', hue: 200 },
  { id: 7, slug: 'food', name: 'еда', hue: 40 },
  { id: 8, slug: 'tech', name: 'tech', hue: 260 },
  { id: 9, slug: 'travel', name: 'путешествия', hue: 180 },
  { id: 10, slug: 'photo', name: 'фото', hue: 320 },
  { id: 11, slug: 'cinema', name: 'кино', hue: 0 },
];

export function findCategory(idOrSlug: number | string): Category | undefined {
  if (typeof idOrSlug === 'number') return CATEGORIES.find((c) => c.id === idOrSlug);
  return CATEGORIES.find((c) => c.slug === idOrSlug);
}

export function categoryColor(c: Category): { fg: string; bg: string; border: string } {
  return {
    fg: `oklch(0.30 0.10 ${c.hue})`,
    bg: `oklch(0.95 0.04 ${c.hue})`,
    border: `oklch(0.80 0.08 ${c.hue})`,
  };
}
