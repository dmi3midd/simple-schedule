import React, { useState } from 'react';
import type { Tag } from '../../types/schedule';
import { useScheduleStore } from '../../store/useScheduleStore';
import { X, Tag as TagIcon, Plus, Trash2, Edit2, Check } from 'lucide-react';

interface TagManagerModalProps {
  isOpen: boolean;
  onClose: () => void;
}

const PRESET_COLORS = [
  '#6366f1', // Indigo
  '#3b82f6', // Blue
  '#06b6d4', // Cyan
  '#10b981', // Emerald
  '#84cc16', // Lime
  '#eab308', // Yellow
  '#f97316', // Orange
  '#ef4444', // Red
  '#ec4899', // Pink
  '#a855f7', // Purple
  '#8b5cf6', // Violet
  '#64748b', // Slate
];

export const TagManagerModal: React.FC<TagManagerModalProps> = ({
  isOpen,
  onClose,
}) => {
  const { tags, createTag, updateTag, deleteTag } = useScheduleStore();

  const [title, setTitle] = useState('');
  const [hexColor, setHexColor] = useState(PRESET_COLORS[0]);
  const [editingTag, setEditingTag] = useState<Tag | null>(null);
  const [isSubmitting, setIsSubmitting] = useState(false);

  if (!isOpen) return null;

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!title.trim()) return;

    setIsSubmitting(true);
    let success = false;
    if (editingTag) {
      success = await updateTag(editingTag.id, title.trim(), hexColor);
      if (success) setEditingTag(null);
    } else {
      success = await createTag(title.trim(), hexColor);
    }
    setIsSubmitting(false);

    if (success) {
      setTitle('');
    }
  };

  const startEdit = (tag: Tag) => {
    setEditingTag(tag);
    setTitle(tag.title);
    setHexColor(tag.hexColor);
  };

  const cancelEdit = () => {
    setEditingTag(null);
    setTitle('');
    setHexColor(PRESET_COLORS[0]);
  };

  const handleDelete = async (id: string, tagTitle: string) => {
    if (window.confirm(`Видалити тег "${tagTitle}"?`)) {
      await deleteTag(id);
      if (editingTag?.id === id) {
        cancelEdit();
      }
    }
  };

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/70 backdrop-blur-sm animate-fade-in">
      <div className="bg-slate-900 border border-slate-800 rounded-2xl w-full max-w-lg p-6 shadow-2xl relative max-h-[90vh] flex flex-col">
        {/* Header */}
        <div className="flex items-center justify-between pb-4 mb-4 border-b border-slate-800 shrink-0">
          <div className="flex items-center gap-2.5">
            <div className="w-8 h-8 rounded-lg bg-indigo-500/10 border border-indigo-500/30 flex items-center justify-center text-indigo-400">
              <TagIcon className="w-5 h-5" />
            </div>
            <div>
              <h3 className="text-lg font-bold text-white">Керування тегами</h3>
              <p className="text-xs text-slate-400">Категорії занять з колірним маркуванням</p>
            </div>
          </div>
          <button
            onClick={onClose}
            className="p-1.5 text-slate-400 hover:text-white rounded-lg hover:bg-slate-800 transition-colors"
          >
            <X className="w-5 h-5" />
          </button>
        </div>

        {/* Form: Create or Edit */}
        <form onSubmit={handleSubmit} className="mb-5 p-4 bg-slate-950/70 border border-slate-800 rounded-xl space-y-3 shrink-0">
          <div className="flex items-center justify-between text-xs font-semibold text-slate-300 uppercase tracking-wider">
            <span>{editingTag ? 'Редагувати тег' : 'Створити новий тег'}</span>
            {editingTag && (
              <button
                type="button"
                onClick={cancelEdit}
                className="text-indigo-400 hover:underline capitalize"
              >
                Скасувати режим редагування
              </button>
            )}
          </div>

          <div className="flex gap-2">
            <input
              type="text"
              required
              placeholder="Назва тегу (напр., Лекція, Лабораторна)"
              value={title}
              onChange={(e) => setTitle(e.target.value)}
              className="flex-1 px-3 py-2 bg-slate-900 border border-slate-700 rounded-xl text-white placeholder-slate-500 text-sm focus:outline-none focus:border-indigo-500"
            />
            <button
              type="submit"
              disabled={isSubmitting || !title.trim()}
              className="px-4 py-2 text-sm font-medium text-white bg-indigo-600 hover:bg-indigo-500 disabled:opacity-50 rounded-xl transition-all shrink-0 flex items-center gap-1.5"
            >
              {editingTag ? <Check className="w-4 h-4" /> : <Plus className="w-4 h-4" />}
              <span>{editingTag ? 'Оновити' : 'Додати'}</span>
            </button>
          </div>

          {/* Color Picker presets */}
          <div>
            <div className="text-[11px] text-slate-400 mb-1.5 flex items-center justify-between">
              <span>Оберіть колір:</span>
              <span className="font-mono text-slate-500 uppercase">{hexColor}</span>
            </div>
            <div className="flex flex-wrap items-center gap-2">
              {PRESET_COLORS.map((color) => (
                <button
                  type="button"
                  key={color}
                  onClick={() => setHexColor(color)}
                  className={`w-6 h-6 rounded-full border-2 transition-transform ${
                    hexColor.toLowerCase() === color.toLowerCase()
                      ? 'scale-125 border-white shadow-md'
                      : 'border-transparent hover:scale-110'
                  }`}
                  style={{ backgroundColor: color }}
                />
              ))}
              <input
                type="color"
                value={hexColor}
                onChange={(e) => setHexColor(e.target.value)}
                className="w-7 h-7 rounded-md cursor-pointer bg-transparent border-0"
                title="Власний колір"
              />
            </div>
          </div>
        </form>

        {/* Tags List */}
        <div className="flex-1 overflow-y-auto pr-1 space-y-2">
          <div className="text-xs font-semibold text-slate-400 uppercase tracking-wider mb-2">
            Наявні теги ({tags.length})
          </div>

          {tags.length === 0 ? (
            <div className="text-center py-8 text-sm text-slate-500">
              Тегів поки немає. Створіть перший тег вище!
            </div>
          ) : (
            tags.map((tag) => (
              <div
                key={tag.id}
                className="flex items-center justify-between p-3 rounded-xl bg-slate-950/40 border border-slate-800/80 hover:border-slate-700/80 transition-colors"
              >
                <div className="flex items-center gap-2.5">
                  <span
                    className="w-3.5 h-3.5 rounded-full shrink-0"
                    style={{ backgroundColor: tag.hexColor }}
                  />
                  <span className="text-sm font-medium text-slate-200">
                    {tag.title}
                  </span>
                  <span className="text-xs font-mono text-slate-500">
                    {tag.hexColor}
                  </span>
                </div>

                <div className="flex items-center gap-1">
                  <button
                    onClick={() => startEdit(tag)}
                    title="Редагувати"
                    className="p-1.5 text-slate-400 hover:text-indigo-400 hover:bg-slate-800 rounded-lg transition-colors"
                  >
                    <Edit2 className="w-3.5 h-3.5" />
                  </button>
                  <button
                    onClick={() => handleDelete(tag.id, tag.title)}
                    title="Видалити"
                    className="p-1.5 text-slate-400 hover:text-rose-400 hover:bg-slate-800 rounded-lg transition-colors"
                  >
                    <Trash2 className="w-3.5 h-3.5" />
                  </button>
                </div>
              </div>
            ))
          )}
        </div>

        {/* Footer */}
        <div className="pt-4 mt-4 border-t border-slate-800 flex justify-end shrink-0">
          <button
            type="button"
            onClick={onClose}
            className="px-4 py-2 text-sm text-slate-300 hover:text-white rounded-xl hover:bg-slate-800 transition-colors"
          >
            Закрити
          </button>
        </div>
      </div>
    </div>
  );
};
