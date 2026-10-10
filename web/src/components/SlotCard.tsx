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
  const tagColor = slot.tag?.hexColor || '#456882';

  return (
    <div
      className="group relative bg-panel/70 hover:bg-panel border border-accent/60 hover:border-sand/40 rounded-xl p-3 transition-all duration-200 shadow-sm flex flex-col justify-between overflow-hidden font-mono"
      style={{
        borderLeftWidth: '4px',
        borderLeftColor: tagColor,
      }}
    >
      <div>
        {/* Header: Time and actions */}
        <div className="flex items-center justify-between gap-1 mb-2">
          <div className="flex items-center gap-1.5 text-[11px] font-semibold text-sand bg-canvas-deep/80 border border-accent/40 px-2 py-0.5 rounded-md">
            <Clock className="w-3 h-3 text-sand-muted" />
            <span>{formatTimeRange(slot.startTime, slot.endTime)}</span>
            <span className="text-sand-dim font-normal">({duration}m)</span>
          </div>

          <div className="flex items-center gap-1 opacity-0 group-hover:opacity-100 transition-opacity">
            <button
              onClick={() => onEdit(slot)}
              title="Редагувати"
              className="p-1 text-sand-muted hover:text-sand hover:bg-accent/40 rounded-md transition-colors"
            >
              <Edit2 className="w-3.5 h-3.5" />
            </button>
            <button
              onClick={() => onDelete(slot.id)}
              title="Видалити"
              className="p-1 text-sand-muted hover:text-rose-300 hover:bg-rose-950/40 rounded-md transition-colors"
            >
              <Trash2 className="w-3.5 h-3.5" />
            </button>
          </div>
        </div>

        {/* Activity Title */}
        <h4 className="text-sm font-semibold text-sand tracking-tight line-clamp-2 mb-2">
          {slot.activity}
        </h4>
      </div>

      {/* Footer: Tag if present */}
      {slot.tag && (
        <div className="mt-1 pt-2 border-t border-accent/40 flex items-center gap-1.5">
          <span
            className="inline-flex items-center gap-1 px-2 py-0.5 rounded-md text-[10px] font-medium"
            style={{
              backgroundColor: `${slot.tag.hexColor}25`,
              color: slot.tag.hexColor,
              border: `1px solid ${slot.tag.hexColor}50`,
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

