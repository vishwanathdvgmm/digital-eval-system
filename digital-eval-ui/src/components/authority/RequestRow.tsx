import React, { useState } from 'react';
import { RequestRow as RequestRowType } from '../../types/authority';
import Card from '../Card';

interface Props {
  request: RequestRowType;
  onApprove: (id: number, assignNum: number) => Promise<void>;
  onReject: (id: number) => Promise<void>;
}

const RequestRow: React.FC<Props> = ({ request, onApprove, onReject }) => {
  const [assignNum, setAssignNum] = useState(5);
  const [loading, setLoading] = useState(false);

  const handleApprove = async () => {
    setLoading(true);
    try {
      await onApprove(request.id, assignNum);
    } finally {
      setLoading(false);
    }
  };

  const handleReject = async () => {
    if (!confirm('Are you sure you want to reject this request?')) return;
    setLoading(true);
    try {
      await onReject(request.id);
    } finally {
      setLoading(false);
    }
  };

  return (
    <Card className="mb-4 border-l-4 border-l-indigo-500 transition-all hover:shadow-md">
      <div className="flex flex-col gap-4 md:flex-row md:items-center md:justify-between">
        <div className="flex-1 space-y-2">
          <div className="flex items-center gap-2">
            <span className="rounded-full bg-indigo-100 px-2 py-1 text-xs font-semibold text-indigo-700">
              ID: {request.id}
            </span>
            <span className="text-sm text-gray-500">
              {new Date(request.created_at).toLocaleString()}
            </span>
          </div>
          <h3 className="text-lg font-bold text-gray-800">
            {request.course_id} <span className="text-gray-400">|</span> {request.semester} - {request.academic_year}
          </h3>
          <p className="text-sm text-gray-600">
            <span className="font-medium text-gray-700">Evaluator:</span> {request.evaluator_id}
          </p>
          <p className="text-sm text-gray-600">
            <span className="font-medium text-gray-700">Description:</span> {request.description}
          </p>
        </div>

        <div className="flex flex-col gap-3 md:items-end">
          <div className="flex items-center gap-2">
            <label htmlFor={`assign-${request.id}`} className="text-sm font-medium text-gray-700">
              Scripts:
            </label>
            <input
              id={`assign-${request.id}`}
              type="number"
              min="1"
              value={assignNum}
              onChange={(e) => setAssignNum(parseInt(e.target.value) || 0)}
              className="w-20 rounded-md border border-gray-300 px-2 py-1 text-sm focus:border-indigo-500 focus:outline-none focus:ring-1 focus:ring-indigo-500"
            />
          </div>
          <div className="flex gap-2">
            <button
              onClick={handleReject}
              disabled={loading}
              className="rounded-md bg-red-50 px-4 py-2 text-sm font-medium text-red-600 hover:bg-red-100 disabled:opacity-50"
            >
              Reject
            </button>
            <button
              onClick={handleApprove}
              disabled={loading}
              className="rounded-md bg-indigo-600 px-4 py-2 text-sm font-medium text-white shadow-sm hover:bg-indigo-700 disabled:opacity-50"
            >
              {loading ? 'Processing...' : 'Approve'}
            </button>
          </div>
        </div>
      </div>
    </Card>
  );
};

export default RequestRow;
