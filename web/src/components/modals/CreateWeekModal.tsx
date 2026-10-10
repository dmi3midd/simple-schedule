import React, { useState } from 'react';
import { useScheduleStore } from '../../store/useScheduleStore';
import { X, Calendar } from 'lucide-react';

interface CreateWeekModalProps {
  isOpen: boolean;
  onClose: () => void;
}

export const CreateWeekModal: React.FC<CreateWeekModalProps> = ({
  isOpen,
  onClose,
}) => {
  const { createWeek } = useScheduleStore();
  const [title, setTitle] = useState('');
  const [isSubmitting, setIsSubmitting] = useState(false);

  if (!isOpen) return null;

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!title.trim()) return;

    setIsSubmitting(true);
    const newId = await createWeek(title.trim());
    setIsSubmitting(false);

    if (newId) {
      setTitle('');
      onClose();
    }
  };

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center p-4 bg-canvas-deep/80 backdrop-blur-sm animate-fade-in font-mono">
      <div className="bg-canvas border border-panel rounded-2xl w-full max-w-md p-6 shadow-2xl relative">
        <div className="flex items-center justify-between pb-4 mb-4 border-b border-panel">
          <div className="flex items-center gap-2.5">
            <div className="w-8 h-8 rounded-lg bg-panel border border-accent flex items-center justify-center text-sand">
              <Calendar className="w-5 h-5" />
            </div>
            <div>
              <h3 className="text-base font-bold text-sand">&gt; Новий тиждень</h3>
              <p className="text-xs text-sand-muted">// створення навчального тижня</p>
            </div>
          </div>
          <button
            onClick={onClose}
            className="p-1.5 text-sand-muted hover:text-sand rounded-lg hover:bg-panel transition-colors"
          >
            <X className="w-5 h-5" />
          </button>
        </div>

        <form onSubmit={handleSubmit} className="space-y-4">
          <div>
            <label className="block text-xs font-semibold text-sand-muted uppercase tracking-wider mb-1.5">
              назва тижня
            </label>
            <input
              type="text"
              required
              autoFocus
              placeholder="напр., 1-й тиждень (Чисельник)"
              value={title}
              onChange={(e) => setTitle(e.target.value)}
              className="w-full px-3.5 py-2.5 bg-canvas-deep border border-panel rounded-xl text-sand placeholder-sand-dim focus:outline-none focus:border-sand text-sm transition-all"
            />
          </div>

          <div className="flex items-center justify-end gap-3 pt-4 border-t border-panel">
            <button
              type="button"
              onClick={onClose}
              className="px-4 py-2 text-sm text-sand-muted hover:text-sand rounded-xl hover:bg-panel transition-colors"
            >
              Скасувати
            </button>
            <button
              type="submit"
              disabled={isSubmitting || !title.trim()}
              className="px-5 py-2 text-sm font-semibold text-canvas-deep bg-sand hover:bg-white disabled:opacity-50 rounded-xl shadow-md transition-all flex items-center gap-1.5"
            >
              {isSubmitting ? 'Створення...' : '+ Створити'}
            </button>
          </div>
        </form>
      </div>
    </div>

  );
};
