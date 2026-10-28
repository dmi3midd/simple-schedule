import { useState } from 'react';
import type { DayOfWeek, SlotWithTag } from '../types/schedule';
import { DAYS_OF_WEEK } from '../utils/time';
import { SlotCard } from './SlotCard';
import { useScheduleStore } from '../store/useScheduleStore';
import { Plus, Calendar } from 'lucide-react';

interface ScheduleGridProps {
  onEditSlot: (slot: SlotWithTag) => void;
  onDeleteSlot: (id: string) => void;
  onOpenCreateWeek: () => void;
}

export const ScheduleGrid: React.FC<ScheduleGridProps> = ({
  onEditSlot,
  onDeleteSlot,
  onOpenCreateWeek,
}) => {
  const { weeks, activeSchedule, isScheduleLoading, isLoading } = useScheduleStore();
  const [dayFilter, setDayFilter] = useState<'first_half' | 'second_half'>('first_half');

  if (isLoading) {
    return (
      <div className="flex flex-col items-center justify-center min-h-[60vh] gap-3 font-mono">
        <div className="w-10 h-10 border-4 border-panel border-t-sand rounded-full animate-spin" />
        <p className="text-sm text-sand-muted">Завантажуємо розклад...</p>
      </div>
    );
  }

  if (weeks.length === 0) {
    return (
      <div className="flex flex-col items-center justify-center min-h-[60vh] text-center px-4 font-mono">
        <div className="w-16 h-16 rounded-2xl bg-panel border border-accent flex items-center justify-center text-sand mb-4 shadow-xl">
          <Calendar className="w-8 h-8" />
        </div>
        <h2 className="text-xl font-bold text-sand mb-2">Жодного тижня не знайдено</h2>
        <p className="text-sm text-sand-muted max-w-md mb-6 leading-relaxed">
          Розклад порожній. Створіть перший навчальний тиждень.
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

  // Filter days based on tab: Mon-Wed or Thu-Sun
  const displayedDays = DAYS_OF_WEEK.filter((day) => {
    if (dayFilter === 'first_half') {
      return ['monday', 'tuesday', 'wednesday'].includes(day.id);
    }
    return ['thursday', 'friday', 'saturday', 'sunday'].includes(day.id);
  });

  return (
    <div className="py-6 px-4 sm:px-6 lg:px-8 max-w-7xl mx-auto font-mono">
      {/* Title banner */}
      <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-3 mb-6">
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
            Всього слотів: {activeSchedule?.slots?.length || 0}
          </p>
        </div>

        {/* Day Range Segmented Control (like in screenshot) */}
        <div className="inline-flex p-1 bg-panel/80 border border-accent/40 rounded-xl">
          <button
            onClick={() => setDayFilter('first_half')}
            className={`px-3.5 py-1.5 text-xs font-semibold rounded-lg transition-all ${dayFilter === 'first_half'
                ? 'bg-sand text-canvas-deep shadow-sm'
                : 'text-sand-muted hover:text-sand'
              }`}
          >
            ПН–СР
          </button>
          <button
            onClick={() => setDayFilter('second_half')}
            className={`px-3.5 py-1.5 text-xs font-semibold rounded-lg transition-all ${dayFilter === 'second_half'
                ? 'bg-sand text-canvas-deep shadow-sm'
                : 'text-sand-muted hover:text-sand'
              }`}
          >
            ЧТ–НД
          </button>
        </div>
      </div>

      {/* Grid columns: wider columns */}
      <div className={`grid gap-4 ${dayFilter === 'first_half'
          ? 'grid-cols-1 md:grid-cols-2 lg:grid-cols-3'
          : 'grid-cols-1 md:grid-cols-2 lg:grid-cols-4'
        }`}>
        {displayedDays.map((day) => {
          const slotsForDay = (daysMap[day.id] || []).slice().sort((a, b) => a.startTime - b.startTime);
          const hasSlots = slotsForDay.length > 0;

          return (
            <div
              key={day.id}
              className="flex flex-col bg-canvas-deep/70 border border-panel rounded-2xl p-4 backdrop-blur-sm min-h-[500px]"
            >
              {/* Day Header */}
              <div className="flex items-center justify-between pb-3 mb-3 border-b border-panel">
                <div className="flex items-center gap-2">
                  <span className="text-sm font-bold text-sand tracking-wide">
                    {day.label}
                  </span>
                </div>
                <span className="text-xs font-mono px-2 py-0.5 rounded-md bg-panel text-sand-muted border border-accent/40">
                  {slotsForDay.length}
                </span>
              </div>

              {/* Slots List */}
              <div className="flex-1 flex flex-col gap-3 overflow-y-auto pr-0.5">
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
                  <div className="flex-1 flex flex-col items-center justify-center py-12 border border-dashed border-panel/80 rounded-xl text-center px-2">
                    <p className="text-xs text-sand-dim">Немає занять</p>
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
