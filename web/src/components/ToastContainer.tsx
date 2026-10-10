import React from 'react';
import { useScheduleStore } from '../store/useScheduleStore';
import { AlertCircle, CheckCircle, Info, X } from 'lucide-react';

export const ToastContainer: React.FC = () => {
  const { toasts, removeToast } = useScheduleStore();

  if (toasts.length === 0) return null;

  return (
    <div className="fixed bottom-5 right-5 z-50 flex flex-col gap-2 max-w-md w-full px-4 sm:px-0 pointer-events-none">
      {toasts.map((toast) => {
        const isError = toast.type === 'error';
        const isSuccess = toast.type === 'success';

        return (
          <div
            key={toast.id}
            className={`pointer-events-auto flex items-start gap-3 p-4 rounded-xl shadow-xl backdrop-blur-md border transition-all duration-300 animate-slide-up ${
              isError
                ? 'bg-rose-950/80 border-rose-800/80 text-rose-200'
                : isSuccess
                ? 'bg-emerald-950/80 border-emerald-800/80 text-emerald-200'
                : 'bg-slate-900/90 border-slate-700/80 text-slate-200'
            }`}
          >
            <div className="shrink-0 mt-0.5">
              {isError && <AlertCircle className="w-5 h-5 text-rose-400" />}
              {isSuccess && <CheckCircle className="w-5 h-5 text-emerald-400" />}
              {!isError && !isSuccess && <Info className="w-5 h-5 text-sky-400" />}
            </div>
            <div className="flex-1 text-sm font-medium leading-relaxed">
              {toast.message}
            </div>
            <button
              onClick={() => removeToast(toast.id)}
              className="shrink-0 text-slate-400 hover:text-white transition-colors p-0.5 rounded-lg hover:bg-white/10"
            >
              <X className="w-4 h-4" />
            </button>
          </div>
        );
      })}
    </div>
  );
};
