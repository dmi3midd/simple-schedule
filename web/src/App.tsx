import React, { useEffect, useState } from 'react';
import { useScheduleStore } from './store/useScheduleStore';
import type { DayOfWeek, SlotWithTag } from './types/schedule';
import { Navbar } from './components/Navbar';
import { ScheduleGrid } from './components/ScheduleGrid';
import { CreateSlotModal } from './components/modals/CreateSlotModal';
import { EditSlotModal } from './components/modals/EditSlotModal';
import { CreateWeekModal } from './components/modals/CreateWeekModal';
import { EditWeekModal } from './components/modals/EditWeekModal';
import { TagManagerModal } from './components/modals/TagManagerModal';
import { ToastContainer } from './components/ToastContainer';

export const App: React.FC = () => {
  const { fetchWeeks, fetchTags, deleteSlot, weeks, activeWeekId } = useScheduleStore();

  const [createSlotOpen, setCreateSlotOpen] = useState(false);
  const [slotDay, setSlotDay] = useState<DayOfWeek>('monday');
  const [editingSlot, setEditingSlot] = useState<SlotWithTag | null>(null);

  const [createWeekOpen, setCreateWeekOpen] = useState(false);
  const [editWeekOpen, setEditWeekOpen] = useState(false);
  const [tagManagerOpen, setTagManagerOpen] = useState(false);

  useEffect(() => {
    fetchWeeks();
    fetchTags();
  }, [fetchWeeks, fetchTags]);

  const activeWeek = weeks.find((w) => w.id === activeWeekId) || null;

  const handleAddSlotForDay = (day: DayOfWeek) => {
    setSlotDay(day);
    setCreateSlotOpen(true);
  };

  const handleEditSlot = (slot: SlotWithTag) => {
    setEditingSlot(slot);
  };

  const handleDeleteSlot = async (id: string) => {
    if (window.confirm('Видалити цей слот?')) {
      await deleteSlot(id);
    }
  };

  return (
    <div className="min-h-screen bg-slate-950 text-slate-100 flex flex-col font-sans selection:bg-indigo-500 selection:text-white">
      {/* Navigation Bar */}
      <Navbar
        onOpenCreateSlot={() => {
          setSlotDay('monday');
          setCreateSlotOpen(true);
        }}
        onOpenCreateWeek={() => setCreateWeekOpen(true)}
        onOpenEditWeek={() => setEditWeekOpen(true)}
        onOpenTagManager={() => setTagManagerOpen(true)}
      />

      {/* Main Content Area */}
      <main className="flex-1">
        <ScheduleGrid
          onAddSlotForDay={handleAddSlotForDay}
          onEditSlot={handleEditSlot}
          onDeleteSlot={handleDeleteSlot}
          onOpenCreateWeek={() => setCreateWeekOpen(true)}
        />
      </main>

      {/* Footer */}
      <footer className="border-t border-slate-900 py-6 text-center text-xs text-slate-500">
        <div className="max-w-7xl mx-auto px-4 flex flex-col sm:flex-row items-center justify-between gap-2">
          <span>SimpleSchedule &copy; 2026. Усі права захищено.</span>
          <div className="flex items-center gap-4 text-slate-400">
            <a href="/swagger/" target="_blank" rel="noreferrer" className="hover:text-indigo-400 transition-colors">
              API Docs (Swagger)
            </a>
          </div>
        </div>
      </footer>

      {/* Modals */}
      <CreateSlotModal
        isOpen={createSlotOpen}
        onClose={() => setCreateSlotOpen(false)}
        defaultDay={slotDay}
      />

      <EditSlotModal
        isOpen={editingSlot !== null}
        onClose={() => setEditingSlot(null)}
        slot={editingSlot}
      />

      <CreateWeekModal
        isOpen={createWeekOpen}
        onClose={() => setCreateWeekOpen(false)}
      />

      <EditWeekModal
        isOpen={editWeekOpen}
        onClose={() => setEditWeekOpen(false)}
        week={activeWeek}
      />

      <TagManagerModal
        isOpen={tagManagerOpen}
        onClose={() => setTagManagerOpen(false)}
      />

      {/* Toast notifications */}
      <ToastContainer />
    </div>
  );
};

export default App;
