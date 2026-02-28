import React, { useState } from 'react';
import { releaseResults } from '../../api/authority';
import Card from '../../components/Card';

const ReleaseResults: React.FC = () => {
  const [semester, setSemester] = useState('');
  const [academicYear, setAcademicYear] = useState('');
  const [loading, setLoading] = useState(false);
  const [success, setSuccess] = useState<string | null>(null);
  const [error, setError] = useState<string | null>(null);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setLoading(true);
    setError(null);
    setSuccess(null);

    // TODO: Get actual user ID from auth context
    const releasedBy = 'authority_1'; 

    try {
      const res = await releaseResults({
        semester,
        academic_year: academicYear,
        released_by: releasedBy,
      });
      setSuccess(`Results released successfully! Block Hash: ${res.block_hash}`);
      setSemester('');
      setAcademicYear('');
    } catch (err: any) {
      setError(err.message || 'Failed to release results');
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="min-h-screen bg-gray-50 p-6">
      <div className="mx-auto max-w-2xl">
        <header className="mb-8 text-center">
          <h1 className="text-3xl font-bold text-gray-900">Release Results</h1>
          <p className="mt-2 text-gray-600">
            Publish evaluation results to the blockchain.
          </p>
        </header>

        <Card className="p-8 shadow-lg">
          <form onSubmit={handleSubmit} className="space-y-6">
            <div>
              <label htmlFor="semester" className="block text-sm font-medium text-gray-700">
                Semester
              </label>
              <input
                id="semester"
                type="text"
                required
                placeholder="e.g. 5"
                value={semester}
                onChange={(e) => setSemester(e.target.value)}
                className="mt-1 block w-full rounded-md border border-gray-300 px-3 py-2 shadow-sm focus:border-indigo-500 focus:outline-none focus:ring-1 focus:ring-indigo-500"
              />
            </div>

            <div>
              <label htmlFor="academicYear" className="block text-sm font-medium text-gray-700">
                Academic Year
              </label>
              <input
                id="academicYear"
                type="text"
                required
                placeholder="e.g. 2024-2025"
                value={academicYear}
                onChange={(e) => setAcademicYear(e.target.value)}
                className="mt-1 block w-full rounded-md border border-gray-300 px-3 py-2 shadow-sm focus:border-indigo-500 focus:outline-none focus:ring-1 focus:ring-indigo-500"
              />
            </div>

            {error && (
              <div className="rounded-md bg-red-50 p-4 text-sm text-red-700 border border-red-200">
                {error}
              </div>
            )}

            {success && (
              <div className="rounded-md bg-green-50 p-4 text-sm text-green-700 border border-green-200 break-all">
                {success}
              </div>
            )}

            <button
              type="submit"
              disabled={loading}
              className="w-full rounded-md bg-indigo-600 px-4 py-3 text-sm font-semibold text-white shadow-md hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:ring-offset-2 disabled:opacity-50"
            >
              {loading ? 'Releasing...' : 'Release Results'}
            </button>
          </form>
        </Card>
      </div>
    </div>
  );
};

export default ReleaseResults;
