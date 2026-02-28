import React, { useEffect, useState } from "react";
import { getUploadHistory } from "../../api/examiner";
import { ScriptRecord } from "../../types/examiner";
import ScriptUploadItem from "../../components/ScriptUploadItem";

const UploadHistory: React.FC = () => {
    const [history, setHistory] = useState<ScriptRecord[]>([]);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState<string | null>(null);

    useEffect(() => {
        const fetchHistory = async () => {
            try {
                const data = await getUploadHistory();
                setHistory(data);
            } catch (err) {
                console.error("Failed to fetch history", err);
                setError("Failed to load upload history.");
            } finally {
                setLoading(false);
            }
        };

        fetchHistory();
    }, []);

    if (loading) {
        return <div className="p-6 text-center text-gray-500">Loading history...</div>;
    }

    if (error) {
        return <div className="p-6 text-center text-red-500">{error}</div>;
    }

    return (
        <div className="space-y-6">
            <div className="flex justify-between items-center">
                <div className="flex items-center gap-4">
                    <h1 className="text-2xl font-bold text-gray-800">Upload History</h1>
                    <button 
                        onClick={() => window.location.reload()} 
                        className="p-2 text-slate-400 hover:text-indigo-600 transition-colors"
                        title="Refresh"
                    >
                        <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15"></path>
                        </svg>
                    </button>
                </div>
                <span className="text-sm text-gray-500">
                    Total Uploads: {history.length}
                </span>
            </div>

            {history.length === 0 ? (
                <div className="text-center py-10 bg-gray-50 rounded-lg border border-gray-200 text-gray-500">
                    No uploads found.
                </div>
            ) : (
                <div className="space-y-4">
                    {history.map((script) => (
                        <ScriptUploadItem key={script.script_id} script={script} />
                    ))}
                </div>
            )}
        </div>
    );
};

export default UploadHistory;
