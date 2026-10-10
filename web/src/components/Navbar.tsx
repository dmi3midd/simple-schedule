import React from 'react';
import { useScheduleStore } from '../store/useScheduleStore';
import {
  Plus,
  Tag as TagIcon,
  BookOpen,
  Edit3,
} from 'lucide-react';

interface NavbarProps {
  onOpenCreateSlot: () => void;
  onOpenCreateWeek: () => void;
  onOpenEditWeek: () => void;
  onOpenTagManager: () => void;
}

export const Navbar: React.FC<NavbarProps> = ({
  onOpenCreateSlot,
  onOpenCreateWeek,
  onOpenEditWeek,
  onOpenTagManager,
}) => {
  const { weeks, activeWeekId, selectWeek } = useScheduleStore();

  const activeWeek = weeks.find((w) => w.id === activeWeekId);

  return (
    <header className="sticky top-0 z-40 bg-canvas-deep/90 backdrop-blur-xl border-b border-panel">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 h-16 flex items-center justify-between gap-4 font-mono">
        {/* Brand / Logo */}
        <div className="flex items-center gap-3">
          {/* <div className="w-10 h-10 rounded-xl bg-panel border border-accent flex items-center justify-center text-sand shadow-md">
            <span className="text-sm font-bold text-sand">&gt;_</span>
          </div> */}
          <div>
            <div className="flex items-center gap-2">
              <span className="font-bold text-base text-sand tracking-tight">
                Simple Schedule
              </span>
            </div>
            <p className="text-xs text-sand-muted hidden sm:block">
              Розклад занять
            </p>
          </div>
        </div>

        {/* Center: Week Selector */}
        <div className="flex items-center gap-2">
          {weeks.length > 0 ? (
            <div className="flex items-center bg-panel border border-accent/80 rounded-xl p-1 shadow-inner">
              <span className="text-xs text-sand-dim pl-2 font-mono hidden sm:inline">Тиждень:</span>
              <select
                value={activeWeekId || ''}
                onChange={(e) => selectWeek(e.target.value)}
                className="bg-transparent text-sm font-medium text-sand px-2 py-1 pr-6 focus:outline-none cursor-pointer appearance-none"
                style={{
                  backgroundImage: `url("data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' fill='none' viewBox='0 0 24 24' stroke='%23D2C1B6'%3E%3Cpath stroke-linecap='round' stroke-linejoin='round' stroke-width='2' d='M19 9l-7 7-7-7'%3E%3C/path%3E%3C/svg%3E")`,
                  backgroundRepeat: 'no-repeat',
                  backgroundPosition: 'right 4px center',
                  backgroundSize: '16px',
                }}
              >
                {weeks.map((w) => (
                  <option key={w.id} value={w.id} className="bg-canvas-deep text-sand">
                    {w.title}
                  </option>
                ))}
              </select>

              <button
                onClick={onOpenEditWeek}
                title="Редагувати або видалити тиждень"
                className="p-1.5 text-sand-muted hover:text-sand rounded-lg hover:bg-accent/40 transition-colors ml-1"
              >
                <Edit3 className="w-3.5 h-3.5" />
              </button>
            </div>
          ) : (
            <div className="text-xs text-sand-muted">Немає тижнів</div>
          )}

          <button
            onClick={onOpenCreateWeek}
            title="Створити новий тиждень"
            className="p-2 text-sand hover:text-white bg-panel hover:bg-panel-light border border-accent/80 rounded-xl transition-all flex items-center gap-1.5 text-xs font-medium"
          >
            <Plus className="w-3.5 h-3.5 text-sand" />
            <span className="hidden md:inline">Тиждень</span>
          </button>
        </div>

        {/* Right side actions */}
        <div className="flex items-center gap-2.5">
          <button
            onClick={onOpenTagManager}
            className="px-3 py-2 text-xs font-medium text-sand hover:text-white bg-panel hover:bg-panel-light border border-accent/80 rounded-xl transition-all flex items-center gap-1.5"
          >
            <TagIcon className="w-3.5 h-3.5 text-sand" />
            <span className="hidden sm:inline">Теги</span>
          </button>

          <a
            href="/swagger/"
            target="_blank"
            rel="noreferrer"
            className="px-3 py-2 text-xs font-medium text-sand-muted hover:text-sand bg-panel/50 hover:bg-panel border border-accent/50 rounded-xl transition-all flex items-center gap-1.5"
            title="Swagger API Documentation"
          >
            <BookOpen className="w-3.5 h-3.5 text-sand-dim" />
            <span className="hidden lg:inline">swagger</span>
          </a>

          {activeWeek && (
            <button
              onClick={onOpenCreateSlot}
              className="px-3.5 py-2 text-xs sm:text-sm font-semibold text-canvas-deep bg-sand hover:bg-white rounded-xl shadow-md transition-all flex items-center gap-1.5 active:scale-95"
            >
              <Plus className="w-4 h-4 text-canvas-deep stroke-[3]" />
              <span>Слот</span>
            </button>
          )}
        </div>
      </div>
    </header>
  );
};
