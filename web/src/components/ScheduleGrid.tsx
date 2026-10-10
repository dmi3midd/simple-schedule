import type { DayOfWeek, SlotWithTag } from '../types/schedule';
import { DAYS_OF_WEEK } from '../utils/time';
import { SlotCard } from './SlotCard';
import { useScheduleStore } from '../store/useScheduleStore';
import { Plus, Calendar } from 'lucide-react';

interface ScheduleGridProps {
  onAddSlotForDay: (day: DayOfWeek) => void;
  onEditSlot: (slot: SlotWithTag) => void;
  onDeleteSlot: (id: string) => void;
  onOpenCreateWeek: () => void;
}

export const ScheduleGrid: React.FC<ScheduleGridProps> = ({
  onAddSlotForDay,
  onEditSlot,
  onDeleteSlot,
  onOpenCreateWeek,
}) => {
  const { weeks, activeSchedule, isScheduleLoading, isLoading } = useScheduleStore();

  if (isLoading) {
    return (
      <div className="flex flex-col items-center justify-center min-h-[60vh] gap-3">
        <div className="w-10 h-10 border-4 border-indigo-500/20 border-t-indigo-500 rounded-full animate-spin" />
        <p className="text-sm text-slate-400">Завантаження розкладу...</p>
      </div>
    );
  }

  if (weeks.length === 0) {
    return (
      <div className="flex flex-col items-center justify-center min-h-[60vh] text-center px-4">
        <div className="w-16 h-16 rounded-3xl bg-indigo-500/10 border border-indigo-500/30 flex items-center justify-center text-indigo-400 mb-4 shadow-xl">
          <Calendar className="w-8 h-8" />
        </div>
        <h2 className="text-2xl font-bold text-white mb-2">Немає жодного тижня</h2>
        <p className="text-sm text-slate-400 max-w-md mb-6 leading-relaxed">
          Розклад поки що порожній. Створіть перший навчальний тиждень (наприклад, "1-й семестр" або "Чисельник"), щоб додати заняття.
        </p>
        <button
          onClick={onOpenCreateWeek}
          className="px-5 py-2.5 text-sm font-semibold text-white bg-indigo-600 hover:bg-indigo-500 rounded-xl shadow-lg shadow-indigo-600/30 transition-all flex items-center gap-2"
        >
          <Plus className="w-4 h-4" />
          <span>Створити перший тиждень</span>
        </button>
      </div>
    );
  }

  const daysMap = activeSchedule?.days || {} as Record<DayOfWeek, SlotWithTag[]>;

  return (
    <div className="py-6 px-4 sm:px-6 lg:px-8 max-w-7xl mx-auto">
      {/* Title banner */}
      <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-2 mb-6">
        <div>
          <div className="flex items-center gap-2">
            <h1 className="text-xl sm:text-2xl font-bold text-white tracking-tight">
              {activeSchedule?.week.title || 'Розклад'}
            </h1>
            {isScheduleLoading && (
              <span className="w-4 h-4 border-2 border-indigo-400/20 border-t-indigo-400 rounded-full animate-spin" />
            )}
          </div>
          <p className="text-xs text-slate-400 mt-0.5">
            Загалом занять на тиждень: {activeSchedule?.slots?.length || 0}
          </p>
        </div>
      </div>

      {/* Grid columns for 7 days */}
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 xl:grid-cols-7 gap-4">
        {DAYS_OF_WEEK.map((day) => {
          const slotsForDay = (daysMap[day.id] || []).slice().sort((a, b) => a.startTime - b.startTime);
          const hasSlots = slotsForDay.length > 0;

          return (
            <div
              key={day.id}
              className="flex flex-col bg-slate-900/40 border border-slate-800/80 rounded-2xl p-3 backdrop-blur-sm min-h-[480px]"
            >
              {/* Day Header */}
              <div className="flex items-center justify-between pb-3 mb-3 border-b border-slate-800/70">
                <div className="flex items-center gap-2">
                  <span className="text-sm font-bold text-slate-200">
                    {day.label}
                  </span>
                  <span className="text-[11px] font-mono px-1.5 py-0.5 rounded-md bg-slate-800 text-slate-400">
                    {slotsForDay.length}
                  </span>
                </div>

                <button
                  onClick={() => onAddSlotForDay(day.id)}
                  title={`Додати слот на ${day.label}`}
                  className="p-1 text-slate-400 hover:text-indigo-400 hover:bg-slate-800/80 rounded-lg transition-colors"
                >
                  <Plus className="w-4 h-4" />
                </button>
              </div>

              {/* Slots List */}
              <div className="flex-1 flex flex-col gap-2.5 overflow-y-auto pr-0.5">
                {hasSlots ? (
                  slotsForDay.map((slot) => (
                    <SlotCard
                      key={slot.id}
                      slot={slot}
                      onEdit={onEditSlot}
                      onDelete={onDeleteSlot}
                    />
                  ))
                ) : (
                  <div className="flex-1 flex flex-col items-center justify-center py-10 border border-dashed border-slate-800/60 rounded-xl text-center px-2">
                    <p className="text-xs text-slate-500 mb-2">Немає занять</p>
                    <button
                      onClick={() => onAddSlotForDay(day.id)}
                      className="text-xs text-indigo-400 hover:text-indigo-300 font-medium hover:underline flex items-center gap-1"
                    >
                      <Plus className="w-3.5 h-3.5" />
                      <span>Додати</span>
                    </button>
                  </div>
                )}
              </div>
            </div>
          );
        })}
      </div>
    </div>
  );
};
