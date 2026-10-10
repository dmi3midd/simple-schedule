export type DayOfWeek =
  | 'monday'
  | 'tuesday'
  | 'wednesday'
  | 'thursday'
  | 'friday'
  | 'saturday'
  | 'sunday';

export interface Tag {
  id: string;
  title: string;
  hexColor: string;
  createdAt: string;
  updatedAt: string;
}

export interface Slot {
  id: string;
  weekId: string;
  tagId?: string | null;
  dayOfWeek: DayOfWeek;
  activity: string;
  startTime: number;
  endTime: number;
  createdAt: string;
  updatedAt: string;
}

export interface SlotWithTag extends Slot {
  tag?: Tag | null;
}

export interface Week {
  id: string;
  title: string;
  createdAt: string;
  updatedAt: string;
}

export interface WeekSchedule {
  week: Week;
  slots: SlotWithTag[];
  days: Record<DayOfWeek, SlotWithTag[]>;
}

export interface CreateSlotPayload {
  weekId: string;
  tagId?: string | null;
  dayOfWeek: DayOfWeek;
  activity: string;
  startTime: number;
  endTime: number;
}

export interface UpdateSlotPayload {
  weekId: string;
  tagId?: string | null;
  dayOfWeek: DayOfWeek;
  activity: string;
  startTime: number;
  endTime: number;
}

export interface ApiUserError {
  code: number;
  message: string;
  timestamp: string;
}
