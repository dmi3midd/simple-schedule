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
      <div className="flex flex-col items-center justify-center min-h-[60vh] gap-3 font-mono">
        <div className="w-10 h-10 border-4 border-panel border-t-sand rounded-full animate-spin" />
        <p className="text-sm text-sand-muted">// loading schedule...</p>
      </div>
    );
  }

  if (weeks.length === 0) {
    return (
      <div className="flex flex-col items-center justify-center min-h-[60vh] text-center px-4 font-mono">
        <div className="w-16 h-16 rounded-2xl bg-panel border border-accent flex items-center justify-center text-sand mb-4 shadow-xl">
          <Calendar className="w-8 h-8" />
        </div>
        <h2 className="text-xl font-bold text-sand mb-2">&gt; [!] Жодного тижня не знайдено</h2>
        <p className="text-sm text-sand-muted max-w-md mb-6 leading-relaxed">
          Розклад порожній. Створіть перший навчальний тиждень (наприклад: "Чисельник", "Знаменник" або "Семестр 1").
        </p>
        <button
          onClick={onOpenCreateWeek}
          className="px-5 py-2.5 text-sm font-semibold text-canvas-deep bg-sand hover:bg-white rounded-xl shadow-md transition-all flex items-center gap-2"
        >
          <Plus className="w-4 h-4 stroke-[3]" />
          <span>Створити тиждень</span>
        </button>
      </div>
    );
  }

  const daysMap = activeSchedule?.days || {} as Record<DayOfWeek, SlotWithTag[]>;

  return (
    <div className="py-6 px-4 sm:px-6 lg:px-8 max-w-7xl mx-auto font-mono">
      {/* Title banner */}
      <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-2 mb-6">
        <div>
          <div className="flex items-center gap-2">
            <span className="text-sand-dim font-bold text-lg">#</span>
            <h1 className="text-xl sm:text-2xl font-bold text-sand tracking-tight">
              {activeSchedule?.week.title || 'Розклад'}
            </h1>
            {isScheduleLoading && (
              <span className="w-4 h-4 border-2 border-panel border-t-sand rounded-full animate-spin" />
            )}
          </div>
          <p className="text-xs text-sand-muted mt-0.5">
            // total slots: {activeSchedule?.slots?.length || 0}
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
              className="flex flex-col bg-canvas-deep/70 border border-panel rounded-2xl p-3 backdrop-blur-sm min-h-[480px]"
            >
              {/* Day Header */}
              <div className="flex items-center justify-between pb-3 mb-3 border-b border-panel">
                <div className="flex items-center gap-2">
                  <span className="text-xs font-bold text-sand tracking-wide uppercase">
                    {day.label}
                  </span>
                  <span className="text-[10px] font-mono px-1.5 py-0.2 rounded bg-panel text-sand-muted border border-accent/40">
                    {slotsForDay.length}
                  </span>
                </div>

                <button
                  onClick={() => onAddSlotForDay(day.id)}
                  title={`Додати слот на ${day.label}`}
                  className="p-1 text-sand-muted hover:text-sand hover:bg-panel rounded-lg transition-colors"
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
                  <div className="flex-1 flex flex-col items-center justify-center py-10 border border-dashed border-panel/80 rounded-xl text-center px-2">
                    <p className="text-xs text-sand-dim mb-2">// empty</p>
                    <button
                      onClick={() => onAddSlotForDay(day.id)}
                      className="text-xs text-sand hover:underline font-medium flex items-center gap-1"
                    >
                      <Plus className="w-3.5 h-3.5" />
                      <span>додати</span>
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
