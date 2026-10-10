import type { DayOfWeek } from '../types/schedule';

export interface DayInfo {
  id: DayOfWeek;
  label: string;
  short: string;
}

export const DAYS_OF_WEEK: DayInfo[] = [
  { id: 'monday', label: 'Понеділок', short: 'Пн' },
  { id: 'tuesday', label: 'Вівторок', short: 'Вт' },
  { id: 'wednesday', label: 'Середа', short: 'Ср' },
  { id: 'thursday', label: 'Четвер', short: 'Чт' },
  { id: 'friday', label: "П'ятниця", short: 'Пт' },
  { id: 'saturday', label: 'Субота', short: 'Сб' },
  { id: 'sunday', label: 'Неділя', short: 'Нд' },
];

export function secondsToTimeString(totalSeconds: number): string {
  const hours = Math.floor(totalSeconds / 3600);
  const minutes = Math.floor((totalSeconds % 3600) / 60);
  const pad = (n: number) => n.toString().padStart(2, '0');
  return `${pad(hours)}:${pad(minutes)}`;
}

export function timeStringToSeconds(timeStr: string): number {
  const [hoursStr, minutesStr] = timeStr.split(':');
  const hours = parseInt(hoursStr || '0', 10);
  const minutes = parseInt(minutesStr || '0', 10);
  return hours * 3600 + minutes * 60;
}

export function formatTimeRange(startSeconds: number, endSeconds: number): string {
  return `${secondsToTimeString(startSeconds)} – ${secondsToTimeString(endSeconds)}`;
}

export function getDurationMinutes(startSeconds: number, endSeconds: number): number {
  return Math.round((endSeconds - startSeconds) / 60);
}
