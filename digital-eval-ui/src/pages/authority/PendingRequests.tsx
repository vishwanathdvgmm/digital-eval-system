import React, { useEffect, useState } from 'react';
import { fetchPendingRequests, approveRequest, rejectRequest } from '../../api/authority';
import { RequestRow as RequestRowType } from '../../types/authority';
import RequestRow from '../../components/authority/RequestRow';

const PendingRequests: React.FC = () => {
  const [requests, setRequests] = useState<RequestRowType[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  const loadRequests = async () => {
    setLoading(true);
    setError(null);
    try {
      const data = await fetchPendingRequests();
      setRequests(data || []);
    } catch (err: any) {
      setError(err.message || 'Failed to load requests');
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    loadRequests();
  }, []);

  const handleApprove = async (id: number, assignNum: number) => {
    try {
      await approveRequest(id, { assign_num: assignNum });
      // Remove from list
      setRequests((prev) => prev.filter((r) => r.id !== id));
    } catch (err: any) {
      alert(`Error approving request: ${err.message}`);
    }
  };

  const handleReject = async (id: number) => {
    try {
      await rejectRequest(id);
      // Remove from list
      setRequests((prev) => prev.filter((r) => r.id !== id));
    } catch (err: any) {
      alert(`Error rejecting request: ${err.message}`);
    }
  };

  return (
    <div className="min-h-screen bg-gray-50 p-6">
      <div className="mx-auto max-w-5xl">
        <header className="mb-8">
          <h1 className="text-3xl font-bold text-gray-900">Pending Requests</h1>
          <p className="mt-2 text-gray-600">
            Manage evaluator requests for script assignments.
          </p>
        </header>

        {error && (
          <div className="mb-6 rounded-md bg-red-50 p-4 text-red-700 border border-red-200">
            {error}
          </div>
        )}

        {loading ? (
          <div className="flex items-center justify-center py-12">
            <div className="h-8 w-8 animate-spin rounded-full border-4 border-indigo-500 border-t-transparent"></div>
          </div>
        ) : requests.length === 0 ? (
          <div className="rounded-lg border border-dashed border-gray-300 bg-white p-12 text-center">
            <p className="text-lg text-gray-500">No pending requests found.</p>
          </div>
        ) : (
          <div className="space-y-4">
            {requests.map((req) => (
              <RequestRow
                key={req.id}
                request={req}
                onApprove={handleApprove}
                onReject={handleReject}
              />
            ))}
          </div>
        )}
      </div>
    </div>
  );
};

export default PendingRequests;
