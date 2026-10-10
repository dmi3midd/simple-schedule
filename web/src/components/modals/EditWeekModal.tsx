import React, { useState, useEffect } from 'react';
import type { Week } from '../../types/schedule';
import { useScheduleStore } from '../../store/useScheduleStore';
import { X, Edit2, Trash2 } from 'lucide-react';

interface EditWeekModalProps {
  isOpen: boolean;
  onClose: () => void;
  week: Week | null;
}

export const EditWeekModal: React.FC<EditWeekModalProps> = ({
  isOpen,
  onClose,
  week,
}) => {
  const { updateWeek, deleteWeek } = useScheduleStore();
  const [title, setTitle] = useState('');
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [isDeleting, setIsDeleting] = useState(false);

  useEffect(() => {
    if (week) {
      setTitle(week.title);
    }
  }, [week]);

  if (!isOpen || !week) return null;

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!title.trim()) return;

    setIsSubmitting(true);
    const success = await updateWeek(week.id, title.trim());
    setIsSubmitting(false);

    if (success) {
      onClose();
    }
  };

  const handleDelete = async () => {
    if (!window.confirm(`Ви дійсно бажаєте видалити тиждень "${week.title}" та всі його заняття?`)) {
      return;
    }

    setIsDeleting(true);
    const success = await deleteWeek(week.id);
    setIsDeleting(false);

    if (success) {
      onClose();
    }
  };

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center p-4 bg-canvas-deep/80 backdrop-blur-sm animate-fade-in font-mono">
      <div className="bg-canvas border border-panel rounded-2xl w-full max-w-md p-6 shadow-2xl relative">
        <div className="flex items-center justify-between pb-4 mb-4 border-b border-panel">
          <div className="flex items-center gap-2.5">
            <div className="w-8 h-8 rounded-lg bg-panel border border-accent flex items-center justify-center text-sand">
              <Edit2 className="w-5 h-5" />
            </div>
            <div>
              <h3 className="text-base font-bold text-sand">&gt; Налаштування тижня</h3>
              <p className="text-xs text-sand-muted">// зміна назви або видалення</p>
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
              value={title}
              onChange={(e) => setTitle(e.target.value)}
              className="w-full px-3.5 py-2.5 bg-canvas-deep border border-panel rounded-xl text-sand placeholder-sand-dim focus:outline-none focus:border-sand text-sm transition-all"
            />
          </div>

          <div className="flex items-center justify-between pt-4 border-t border-panel">
            <button
              type="button"
              onClick={handleDelete}
              disabled={isDeleting}
              className="px-3.5 py-2 text-xs font-medium text-rose-300 hover:text-rose-100 bg-rose-950/40 hover:bg-rose-900/60 border border-rose-800/60 rounded-xl transition-all flex items-center gap-1.5"
            >
              <Trash2 className="w-3.5 h-3.5" />
              <span>{isDeleting ? 'Видалення...' : 'Видалити'}</span>
            </button>

            <div className="flex items-center gap-2">
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
                className="px-4 py-2 text-sm font-semibold text-canvas-deep bg-sand hover:bg-white disabled:opacity-50 rounded-xl shadow-md transition-all flex items-center gap-1.5"
              >
                {isSubmitting ? 'Збереження...' : 'Зберегти'}
              </button>
            </div>
          </div>
        </form>
      </div>
    </div>
  );
};
