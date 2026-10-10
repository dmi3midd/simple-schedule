import type { SlotWithTag } from '../types/schedule';
import { formatTimeRange, getDurationMinutes } from '../utils/time';
import { Clock, Edit2, Trash2 } from 'lucide-react';

interface SlotCardProps {
  slot: SlotWithTag;
  onEdit: (slot: SlotWithTag) => void;
  onDelete: (id: string) => void;
}

export const SlotCard: React.FC<SlotCardProps> = ({ slot, onEdit, onDelete }) => {
  const duration = getDurationMinutes(slot.startTime, slot.endTime);
  const tagColor = slot.tag?.hexColor || '#6366f1';

  return (
    <div
      className="group relative bg-slate-900/70 hover:bg-slate-800/90 border border-slate-800/90 hover:border-slate-700/80 rounded-xl p-3.5 transition-all duration-200 shadow-sm hover:shadow-md flex flex-col justify-between overflow-hidden"
      style={{
        borderLeftWidth: '4px',
        borderLeftColor: tagColor,
      }}
    >
      <div>
        {/* Header: Time and actions */}
        <div className="flex items-center justify-between gap-1 mb-2">
          <div className="flex items-center gap-1.5 text-xs font-semibold text-slate-300 bg-slate-800/80 px-2 py-0.5 rounded-md">
            <Clock className="w-3 h-3 text-indigo-400" />
            <span>{formatTimeRange(slot.startTime, slot.endTime)}</span>
            <span className="text-slate-500 font-normal">({duration} хв)</span>
          </div>

          <div className="flex items-center gap-1 opacity-0 group-hover:opacity-100 transition-opacity">
            <button
              onClick={() => onEdit(slot)}
              title="Редагувати"
              className="p-1 text-slate-400 hover:text-indigo-300 hover:bg-indigo-950/60 rounded-md transition-colors"
            >
              <Edit2 className="w-3.5 h-3.5" />
            </button>
            <button
              onClick={() => onDelete(slot.id)}
              title="Видалити"
              className="p-1 text-slate-400 hover:text-rose-400 hover:bg-rose-950/60 rounded-md transition-colors"
            >
              <Trash2 className="w-3.5 h-3.5" />
            </button>
          </div>
        </div>

        {/* Activity Title */}
        <h4 className="text-sm font-semibold text-white tracking-tight line-clamp-2 mb-2">
          {slot.activity}
        </h4>
      </div>

      {/* Footer: Tag if present */}
      {slot.tag && (
        <div className="mt-1 pt-2 border-t border-slate-800/60 flex items-center gap-1.5">
          <span
            className="inline-flex items-center gap-1 px-2 py-0.5 rounded-full text-[11px] font-medium"
            style={{
              backgroundColor: `${slot.tag.hexColor}20`,
              color: slot.tag.hexColor,
              border: `1px solid ${slot.tag.hexColor}40`,
            }}
          >
            <span
              className="w-1.5 h-1.5 rounded-full"
              style={{ backgroundColor: slot.tag.hexColor }}
            />
            {slot.tag.title}
          </span>
        </div>
      )}
    </div>
  );
};
