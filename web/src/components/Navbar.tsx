import React from 'react';
import { useScheduleStore } from '../store/useScheduleStore';
import {
  Calendar,
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
    <header className="sticky top-0 z-40 bg-slate-950/80 backdrop-blur-xl border-b border-slate-800/80">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 h-16 flex items-center justify-between gap-4">
        {/* Brand / Logo */}
        <div className="flex items-center gap-3">
          <div className="w-10 h-10 rounded-xl bg-gradient-to-br from-indigo-500 to-purple-600 flex items-center justify-center text-white shadow-lg shadow-indigo-500/25">
            <Calendar className="w-5 h-5" />
          </div>
          <div>
            <div className="flex items-center gap-2">
              <span className="font-bold text-base text-white tracking-tight">
                SimpleSchedule
              </span>
              <span className="text-[10px] font-semibold tracking-wider uppercase px-2 py-0.5 rounded-full bg-indigo-500/10 text-indigo-400 border border-indigo-500/20">
                v1.0
              </span>
            </div>
            <p className="text-xs text-slate-400 hidden sm:block">
              Розклад занять та дисциплін
            </p>
          </div>
        </div>

        {/* Center: Week Selector */}
        <div className="flex items-center gap-2">
          {weeks.length > 0 ? (
            <div className="flex items-center bg-slate-900 border border-slate-800 rounded-xl p-1 shadow-inner">
              <select
                value={activeWeekId || ''}
                onChange={(e) => selectWeek(e.target.value)}
                className="bg-transparent text-sm font-medium text-white px-2 py-1 pr-6 focus:outline-none cursor-pointer appearance-none"
                style={{
                  backgroundImage: `url("data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' fill='none' viewBox='0 0 24 24' stroke='%2394a3b8'%3E%3Cpath stroke-linecap='round' stroke-linejoin='round' stroke-width='2' d='M19 9l-7 7-7-7'%3E%3C/path%3E%3C/svg%3E")`,
                  backgroundRepeat: 'no-repeat',
                  backgroundPosition: 'right 4px center',
                  backgroundSize: '16px',
                }}
              >
                {weeks.map((w) => (
                  <option key={w.id} value={w.id} className="bg-slate-900 text-white">
                    {w.title}
                  </option>
                ))}
              </select>

              <button
                onClick={onOpenEditWeek}
                title="Редагувати або видалити тиждень"
                className="p-1.5 text-slate-400 hover:text-white rounded-lg hover:bg-slate-800 transition-colors ml-1"
              >
                <Edit3 className="w-3.5 h-3.5" />
              </button>
            </div>
          ) : (
            <div className="text-xs text-slate-400">Тижнів ще немає</div>
          )}

          <button
            onClick={onOpenCreateWeek}
            title="Створити новий тиждень"
            className="p-2 text-slate-300 hover:text-white bg-slate-900 hover:bg-slate-800 border border-slate-800 rounded-xl transition-all flex items-center gap-1.5 text-xs font-medium"
          >
            <Plus className="w-3.5 h-3.5 text-indigo-400" />
            <span className="hidden md:inline">Новий тиждень</span>
          </button>
        </div>

        {/* Right side actions */}
        <div className="flex items-center gap-2.5">
          <button
            onClick={onOpenTagManager}
            className="px-3 py-2 text-xs font-medium text-slate-300 hover:text-white bg-slate-900 hover:bg-slate-800 border border-slate-800 rounded-xl transition-all flex items-center gap-1.5"
          >
            <TagIcon className="w-3.5 h-3.5 text-purple-400" />
            <span className="hidden sm:inline">Теги</span>
          </button>

          <a
            href="/swagger/"
            target="_blank"
            rel="noreferrer"
            className="px-3 py-2 text-xs font-medium text-slate-400 hover:text-white bg-slate-900/50 hover:bg-slate-800 border border-slate-800/80 rounded-xl transition-all flex items-center gap-1.5"
            title="Swagger API Documentation"
          >
            <BookOpen className="w-3.5 h-3.5 text-slate-400" />
            <span className="hidden lg:inline">Swagger</span>
          </a>

          {activeWeek && (
            <button
              onClick={onOpenCreateSlot}
              className="px-3.5 py-2 text-xs sm:text-sm font-semibold text-white bg-indigo-600 hover:bg-indigo-500 rounded-xl shadow-lg shadow-indigo-600/25 transition-all flex items-center gap-1.5 active:scale-95"
            >
              <Plus className="w-4 h-4" />
              <span>Додати слот</span>
            </button>
          )}
        </div>
      </div>
    </header>
  );
};
