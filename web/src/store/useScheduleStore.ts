import { create } from 'zustand';
import { api } from '../api/endpoints';
import type {
  Week,
  WeekSchedule,
  Tag,
  CreateSlotPayload,
  UpdateSlotPayload,
} from '../types/schedule';
import { ApiError } from '../api/client';

export interface ToastMessage {
  id: string;
  message: string;
  type: 'success' | 'error' | 'info';
}

interface ScheduleState {
  weeks: Week[];
  activeWeekId: string | null;
  activeSchedule: WeekSchedule | null;
  tags: Tag[];
  isLoading: boolean;
  isScheduleLoading: boolean;
  toasts: ToastMessage[];

  // Actions
  fetchWeeks: () => Promise<void>;
  selectWeek: (id: string) => Promise<void>;
  refreshSchedule: () => Promise<void>;
  fetchTags: () => Promise<void>;

  createWeek: (title: string) => Promise<string | null>;
  updateWeek: (id: string, title: string) => Promise<boolean>;
  deleteWeek: (id: string) => Promise<boolean>;

  createSlot: (payload: CreateSlotPayload) => Promise<boolean>;
  updateSlot: (id: string, payload: UpdateSlotPayload) => Promise<boolean>;
  deleteSlot: (id: string) => Promise<boolean>;

  createTag: (title: string, hexColor: string) => Promise<boolean>;
  updateTag: (id: string, title: string, hexColor: string) => Promise<boolean>;
  deleteTag: (id: string) => Promise<boolean>;

  addToast: (message: string, type?: 'success' | 'error' | 'info') => void;
  removeToast: (id: string) => void;
}

export const useScheduleStore = create<ScheduleState>((set, get) => ({
  weeks: [],
  activeWeekId: null,
  activeSchedule: null,
  tags: [],
  isLoading: false,
  isScheduleLoading: false,
  toasts: [],

  addToast: (message: string, type = 'info') => {
    const id = Math.random().toString(36).substring(2, 9);
    set((state) => ({
      toasts: [...state.toasts, { id, message, type }],
    }));
    setTimeout(() => {
      get().removeToast(id);
    }, 4500);
  },

  removeToast: (id: string) => {
    set((state) => ({
      toasts: state.toasts.filter((t) => t.id !== id),
    }));
  },

  fetchWeeks: async () => {
    set({ isLoading: true });
    try {
      const weeks = await api.getWeeks();
      set({ weeks });

      const currentActive = get().activeWeekId;
      if (weeks.length > 0) {
        if (!currentActive || !weeks.some((w) => w.id === currentActive)) {
          await get().selectWeek(weeks[0].id);
        } else {
          await get().refreshSchedule();
        }
      } else {
        set({ activeWeekId: null, activeSchedule: null });
      }
    } catch (err) {
      const message = err instanceof Error ? err.message : 'Не вдалося завантажити тижні';
      get().addToast(message, 'error');
    } finally {
      set({ isLoading: false });
    }
  },

  selectWeek: async (id: string) => {
    set({ activeWeekId: id, isScheduleLoading: true });
    try {
      const schedule = await api.getWeekSchedule(id);
      set({ activeSchedule: schedule });
    } catch (err) {
      const message = err instanceof Error ? err.message : 'Не вдалося завантажити розклад тижня';
      get().addToast(message, 'error');
      set({ activeSchedule: null });
    } finally {
      set({ isScheduleLoading: false });
    }
  },

  refreshSchedule: async () => {
    const activeId = get().activeWeekId;
    if (!activeId) return;
    try {
      const schedule = await api.getWeekSchedule(activeId);
      set({ activeSchedule: schedule });
    } catch (err) {
      const message = err instanceof Error ? err.message : 'Не вдалося оновити розклад';
      get().addToast(message, 'error');
    }
  },

  fetchTags: async () => {
    try {
      const tags = await api.getTags();
      set({ tags });
    } catch (err) {
      const message = err instanceof Error ? err.message : 'Не вдалося завантажити теги';
      get().addToast(message, 'error');
    }
  },

  createWeek: async (title: string) => {
    try {
      const newId = await api.createWeek(title);
      get().addToast('Тиждень успішно створено', 'success');
      await get().fetchWeeks();
      if (newId) {
        await get().selectWeek(newId);
      }
      return newId;
    } catch (err) {
      const message = err instanceof ApiError ? err.message : 'Помилка при створенні тижня';
      get().addToast(message, 'error');
      return null;
    }
  },

  updateWeek: async (id: string, title: string) => {
    try {
      await api.updateWeek(id, title);
      get().addToast('Назву тижня оновлено', 'success');
      await get().fetchWeeks();
      return true;
    } catch (err) {
      const message = err instanceof ApiError ? err.message : 'Помилка при оновленні тижня';
      get().addToast(message, 'error');
      return false;
    }
  },

  deleteWeek: async (id: string) => {
    try {
      await api.deleteWeek(id);
      get().addToast('Тиждень видалено', 'success');
      await get().fetchWeeks();
      return true;
    } catch (err) {
      const message = err instanceof ApiError ? err.message : 'Помилка при видаленні тижня';
      get().addToast(message, 'error');
      return false;
    }
  },

  createSlot: async (payload: CreateSlotPayload) => {
    try {
      await api.createSlot(payload);
      get().addToast('Слот успішно додано', 'success');
      await get().refreshSchedule();
      return true;
    } catch (err) {
      const message = err instanceof ApiError ? err.message : 'Помилка при створенні слоту';
      get().addToast(message, 'error');
      return false;
    }
  },

  updateSlot: async (id: string, payload: UpdateSlotPayload) => {
    try {
      await api.updateSlot(id, payload);
      get().addToast('Слот оновлено', 'success');
      await get().refreshSchedule();
      return true;
    } catch (err) {
      const message = err instanceof ApiError ? err.message : 'Помилка при оновленні слоту';
      get().addToast(message, 'error');
      return false;
    }
  },

  deleteSlot: async (id: string) => {
    try {
      await api.deleteSlot(id);
      get().addToast('Слот видалено', 'success');
      await get().refreshSchedule();
      return true;
    } catch (err) {
      const message = err instanceof ApiError ? err.message : 'Помилка при видаленні слоту';
      get().addToast(message, 'error');
      return false;
    }
  },

  createTag: async (title: string, hexColor: string) => {
    try {
      await api.createTag(title, hexColor);
      get().addToast('Тег створено', 'success');
      await get().fetchTags();
      return true;
    } catch (err) {
      const message = err instanceof ApiError ? err.message : 'Помилка при створенні тегу';
      get().addToast(message, 'error');
      return false;
    }
  },

  updateTag: async (id: string, title: string, hexColor: string) => {
    try {
      await api.updateTag(id, title, hexColor);
      get().addToast('Тег оновлено', 'success');
      await get().fetchTags();
      await get().refreshSchedule();
      return true;
    } catch (err) {
      const message = err instanceof ApiError ? err.message : 'Помилка при оновленні тегу';
      get().addToast(message, 'error');
      return false;
    }
  },

  deleteTag: async (id: string) => {
    try {
      await api.deleteTag(id);
      get().addToast('Тег видалено', 'success');
      await get().fetchTags();
      await get().refreshSchedule();
      return true;
    } catch (err) {
      const message = err instanceof ApiError ? err.message : 'Помилка при видаленні тегу';
      get().addToast(message, 'error');
      return false;
    }
  },
}));
