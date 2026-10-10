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
    <div className="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/70 backdrop-blur-sm animate-fade-in">
      <div className="bg-slate-900 border border-slate-800 rounded-2xl w-full max-w-md p-6 shadow-2xl relative">
        <div className="flex items-center justify-between pb-4 mb-4 border-b border-slate-800">
          <div className="flex items-center gap-2.5">
            <div className="w-8 h-8 rounded-lg bg-indigo-500/10 border border-indigo-500/30 flex items-center justify-center text-indigo-400">
              <Edit2 className="w-5 h-5" />
            </div>
            <div>
              <h3 className="text-lg font-bold text-white">Редагувати тиждень</h3>
              <p className="text-xs text-slate-400">Зміна назви або видалення тижня</p>
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
          <div>
            <label className="block text-xs font-semibold text-slate-300 uppercase tracking-wider mb-1.5">
              Назва тижня
            </label>
            <input
              type="text"
              required
              value={title}
              onChange={(e) => setTitle(e.target.value)}
              className="w-full px-3.5 py-2.5 bg-slate-950 border border-slate-700 rounded-xl text-white placeholder-slate-500 focus:outline-none focus:border-indigo-500 text-sm transition-all"
            />
          </div>

          <div className="flex items-center justify-between pt-4 border-t border-slate-800">
            <button
              type="button"
              onClick={handleDelete}
              disabled={isDeleting}
              className="px-3.5 py-2 text-sm font-medium text-rose-400 hover:text-rose-300 hover:bg-rose-950/40 rounded-xl transition-colors flex items-center gap-1.5"
            >
              <Trash2 className="w-4 h-4" />
              <span>{isDeleting ? 'Видалення...' : 'Видалити'}</span>
            </button>

            <div className="flex items-center gap-2">
              <button
                type="button"
                onClick={onClose}
                className="px-4 py-2 text-sm text-slate-400 hover:text-white rounded-xl hover:bg-slate-800 transition-colors"
              >
                Скасувати
              </button>
              <button
                type="submit"
                disabled={isSubmitting || !title.trim()}
                className="px-5 py-2 text-sm font-medium text-white bg-indigo-600 hover:bg-indigo-500 disabled:opacity-50 rounded-xl shadow-lg shadow-indigo-600/30 transition-all flex items-center gap-1.5"
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
