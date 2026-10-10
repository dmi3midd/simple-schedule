import { request } from './client';
import type {
  Week,
  WeekSchedule,
  SlotWithTag,
  Tag,
  CreateSlotPayload,
  UpdateSlotPayload,
} from '../types/schedule';

export const api = {
  // Weeks
  async getWeeks(): Promise<Week[]> {
    const res = await request<{ weeks: Week[] }>('/api/weeks');
    return res.weeks || [];
  },

  async getWeek(id: string): Promise<Week> {
    const res = await request<{ week: Week }>(`/api/weeks/${id}`);
    return res.week;
  },

  async getWeekSchedule(id: string): Promise<WeekSchedule> {
    const res = await request<{ schedule: WeekSchedule }>(`/api/weeks/${id}/schedule`);
    return res.schedule;
  },

  async createWeek(title: string): Promise<string> {
    const res = await request<{ id: string }>('/api/weeks', {
      method: 'POST',
      body: JSON.stringify({ title }),
    });
    return res.id;
  },

  async updateWeek(id: string, title: string): Promise<string> {
    const res = await request<{ id: string }>(`/api/weeks/${id}`, {
      method: 'PUT',
      body: JSON.stringify({ title }),
    });
    return res.id;
  },

  async deleteWeek(id: string): Promise<void> {
    await request<void>(`/api/weeks/${id}`, {
      method: 'DELETE',
    });
  },

  // Slots
  async getSlotsByWeek(weekId: string): Promise<SlotWithTag[]> {
    const res = await request<{ slots: SlotWithTag[] }>(`/api/weeks/${weekId}/slots`);
    return res.slots || [];
  },

  async getSlot(id: string): Promise<SlotWithTag> {
    const res = await request<{ slot: SlotWithTag }>(`/api/slots/${id}`);
    return res.slot;
  },

  async createSlot(payload: CreateSlotPayload): Promise<string> {
    const res = await request<{ id: string }>('/api/slots', {
      method: 'POST',
      body: JSON.stringify(payload),
    });
    return res.id;
  },

  async updateSlot(id: string, payload: UpdateSlotPayload): Promise<void> {
    await request<void>(`/api/slots/${id}`, {
      method: 'PUT',
      body: JSON.stringify(payload),
    });
  },

  async deleteSlot(id: string): Promise<void> {
    await request<void>(`/api/slots/${id}`, {
      method: 'DELETE',
    });
  },

  // Tags
  async getTags(): Promise<Tag[]> {
    const res = await request<{ tags: Tag[] }>('/api/tags');
    return res.tags || [];
  },

  async createTag(title: string, hexColor: string): Promise<string> {
    const res = await request<{ id: string }>('/api/tags', {
      method: 'POST',
      body: JSON.stringify({ title, hexColor }),
    });
    return res.id;
  },

  async updateTag(id: string, title: string, hexColor: string): Promise<string> {
    const res = await request<{ id: string }>(`/api/tags/${id}`, {
      method: 'PUT',
      body: JSON.stringify({ title, hexColor }),
    });
    return res.id;
  },

  async deleteTag(id: string): Promise<void> {
    await request<void>(`/api/tags/${id}`, {
      method: 'DELETE',
    });
  },
};
