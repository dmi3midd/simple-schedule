import React, { useState } from 'react';
import type { DayOfWeek } from '../../types/schedule';
import { DAYS_OF_WEEK, timeStringToSeconds } from '../../utils/time';
import { useScheduleStore } from '../../store/useScheduleStore';
import { X, Plus, Sparkles } from 'lucide-react';

interface CreateSlotModalProps {
  isOpen: boolean;
  onClose: () => void;
  defaultDay?: DayOfWeek;
}

const PRESET_TIMES = [
  { label: '1 пара (08:30 – 10:05)', start: '08:30', end: '10:05' },
  { label: '2 пара (10:20 – 11:55)', start: '10:20', end: '11:55' },
  { label: '3 пара (12:10 – 13:45)', start: '12:10', end: '13:45' },
  { label: '4 пара (14:15 – 15:50)', start: '14:15', end: '15:50' },
  { label: '5 пара (16:05 – 17:40)', start: '16:05', end: '17:40' },
];

export const CreateSlotModal: React.FC<CreateSlotModalProps> = ({
  isOpen,
  onClose,
  defaultDay = 'monday',
}) => {
  const { activeWeekId, tags, createSlot } = useScheduleStore();

  const [dayOfWeek, setDayOfWeek] = useState<DayOfWeek>(defaultDay);
  const [activity, setActivity] = useState('');
  const [startTime, setStartTime] = useState('08:30');
  const [endTime, setEndTime] = useState('10:05');
  const [selectedTagId, setSelectedTagId] = useState<string>('');
  const [isSubmitting, setIsSubmitting] = useState(false);

  React.useEffect(() => {
    if (isOpen) {
      setDayOfWeek(defaultDay);
    }
  }, [isOpen, defaultDay]);

  if (!isOpen) return null;

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!activeWeekId || !activity.trim()) return;

    const startSeconds = timeStringToSeconds(startTime);
    const endSeconds = timeStringToSeconds(endTime);

    if (startSeconds >= endSeconds) {
      alert('Час початку має бути раніше за час завершення');
      return;
    }

    setIsSubmitting(true);
    const success = await createSlot({
      weekId: activeWeekId,
      dayOfWeek,
      activity: activity.trim(),
      startTime: startSeconds,
      endTime: endSeconds,
      tagId: selectedTagId ? selectedTagId : null,
    });
    setIsSubmitting(false);

    if (success) {
      setActivity('');
      setSelectedTagId('');
      onClose();
    }
  };

  const applyPreset = (start: string, end: string) => {
    setStartTime(start);
    setEndTime(end);
  };

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/70 backdrop-blur-sm animate-fade-in">
      <div className="bg-slate-900 border border-slate-800 rounded-2xl w-full max-w-lg p-6 shadow-2xl relative">
        <div className="flex items-center justify-between pb-4 mb-4 border-b border-slate-800">
          <div className="flex items-center gap-2.5">
            <div className="w-8 h-8 rounded-lg bg-indigo-500/10 border border-indigo-500/30 flex items-center justify-center text-indigo-400">
              <Plus className="w-5 h-5" />
            </div>
            <div>
              <h3 className="text-lg font-bold text-white">Додати заняття / слот</h3>
              <p className="text-xs text-slate-400">Створення нового слоту в розкладі</p>
            </div>
          </div>
          <button
            onClick={onClose}
            className="p-1.5 text-slate-400 hover:text-white rounded-lg hover:bg-slate-800 transition-colors"
          >
            <X className="w-5 h-5" />
          </button>
        </div>

        <form onSubmit={handleSubmit} className="space-y-4">
          {/* Day of Week */}
          <div>
            <label className="block text-xs font-semibold text-slate-300 uppercase tracking-wider mb-2">
              День тижня
            </label>
            <div className="grid grid-cols-4 sm:grid-cols-7 gap-1.5">
              {DAYS_OF_WEEK.map((d) => (
                <button
                  type="button"
                  key={d.id}
                  onClick={() => setDayOfWeek(d.id)}
                  className={`py-2 px-1 text-center text-xs font-medium rounded-lg border transition-all ${
                    dayOfWeek === d.id
                      ? 'bg-indigo-600 text-white border-indigo-500 shadow-sm'
                      : 'bg-slate-800/60 text-slate-400 border-slate-700/60 hover:text-white hover:bg-slate-800'
                  }`}
                >
                  {d.short}
                </button>
              ))}
            </div>
          </div>

          {/* Activity */}
          <div>
            <label className="block text-xs font-semibold text-slate-300 uppercase tracking-wider mb-1.5">
              Назва дисципліни / активності
            </label>
            <input
              type="text"
              required
              placeholder="напр., Вища математика, Лекція"
              value={activity}
              onChange={(e) => setActivity(e.target.value)}
              className="w-full px-3.5 py-2.5 bg-slate-950 border border-slate-700 rounded-xl text-white placeholder-slate-500 focus:outline-none focus:border-indigo-500 focus:ring-1 focus:ring-indigo-500 text-sm transition-all"
            />
          </div>

          {/* Presets */}
          <div>
            <div className="flex items-center gap-1.5 text-xs font-semibold text-slate-400 uppercase tracking-wider mb-1.5">
              <Sparkles className="w-3.5 h-3.5 text-amber-400" />
              <span>Швидкий вибір пари</span>
            </div>
            <div className="flex flex-wrap gap-1.5">
              {PRESET_TIMES.map((preset, idx) => (
                <button
                  type="button"
                  key={idx}
                  onClick={() => applyPreset(preset.start, preset.end)}
                  className={`text-[11px] px-2.5 py-1 rounded-md border transition-all ${
                    startTime === preset.start && endTime === preset.end
                      ? 'bg-indigo-950 text-indigo-300 border-indigo-700'
                      : 'bg-slate-800/40 text-slate-400 border-slate-800 hover:text-slate-200 hover:bg-slate-800'
                  }`}
                >
                  {preset.label}
                </button>
              ))}
            </div>
          </div>

          {/* Times */}
          <div className="grid grid-cols-2 gap-3">
            <div>
              <label className="block text-xs font-semibold text-slate-300 uppercase tracking-wider mb-1.5">
                Час початку
              </label>
              <div className="relative">
                <input
                  type="time"
                  required
                  value={startTime}
                  onChange={(e) => setStartTime(e.target.value)}
                  className="w-full px-3 py-2 bg-slate-950 border border-slate-700 rounded-xl text-white text-sm focus:outline-none focus:border-indigo-500"
                />
              </div>
            </div>
            <div>
              <label className="block text-xs font-semibold text-slate-300 uppercase tracking-wider mb-1.5">
                Час завершення
              </label>
              <div className="relative">
                <input
                  type="time"
                  required
                  value={endTime}
                  onChange={(e) => setEndTime(e.target.value)}
                  className="w-full px-3 py-2 bg-slate-950 border border-slate-700 rounded-xl text-white text-sm focus:outline-none focus:border-indigo-500"
                />
              </div>
            </div>
          </div>

          {/* Tag */}
          <div>
            <label className="block text-xs font-semibold text-slate-300 uppercase tracking-wider mb-1.5">
              Тег (категорія)
            </label>
            <div className="flex flex-wrap gap-2">
              <button
                type="button"
                onClick={() => setSelectedTagId('')}
                className={`text-xs px-3 py-1.5 rounded-lg border transition-all ${
                  selectedTagId === ''
                    ? 'bg-slate-700 text-white border-slate-500'
                    : 'bg-slate-950 text-slate-400 border-slate-800 hover:bg-slate-800'
                }`}
              >
                Без тегу
              </button>
              {tags.map((tag) => {
                const isSelected = selectedTagId === tag.id;
                return (
                  <button
                    type="button"
                    key={tag.id}
                    onClick={() => setSelectedTagId(tag.id)}
                    className="text-xs px-3 py-1.5 rounded-lg border transition-all flex items-center gap-1.5"
                    style={{
                      backgroundColor: isSelected ? `${tag.hexColor}33` : 'rgba(2, 6, 23, 0.7)',
                      borderColor: isSelected ? tag.hexColor : 'rgba(51, 65, 85, 0.6)',
                      color: isSelected ? tag.hexColor : '#cbd5e1',
                    }}
                  >
                    <span
                      className="w-2 h-2 rounded-full"
                      style={{ backgroundColor: tag.hexColor }}
                    />
                    {tag.title}
                  </button>
                );
              })}
            </div>
          </div>

          {/* Buttons */}
          <div className="flex items-center justify-end gap-3 pt-4 border-t border-slate-800 mt-6">
            <button
              type="button"
              onClick={onClose}
              className="px-4 py-2 text-sm text-slate-400 hover:text-white rounded-xl hover:bg-slate-800 transition-colors"
            >
              Скасувати
            </button>
            <button
              type="submit"
              disabled={isSubmitting}
              className="px-5 py-2 text-sm font-medium text-white bg-indigo-600 hover:bg-indigo-500 disabled:opacity-50 rounded-xl shadow-lg shadow-indigo-600/30 transition-all flex items-center gap-1.5"
            >
              {isSubmitting ? 'Збереження...' : 'Додати слот'}
            </button>
          </div>
        </form>
      </div>
    </div>
  );
};
