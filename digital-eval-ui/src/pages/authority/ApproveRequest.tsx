import React, { useEffect, useState } from "react";
import { fetchRequestHistory } from "../../api/authority";
import { RequestRow } from "../../types/authority";
import Card from "../../components/Card";

const ApproveRequest: React.FC = () => {
    const [history, setHistory] = useState<RequestRow[]>([]);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState<string | null>(null);

    useEffect(() => {
        loadHistory();
    }, []);

    const loadHistory = async () => {
        try {
            const data = await fetchRequestHistory();
            setHistory(data || []);
        } catch (err: any) {
            setError(err.message || "Failed to load history");
        } finally {
            setLoading(false);
        }
    };

    return (
        <div className="p-6 max-w-6xl mx-auto">
            <div className="mb-8">
                <h1 className="text-3xl font-bold text-slate-800">Request History</h1>
                <p className="text-slate-500 mt-1">View past evaluator requests and their status.</p>
            </div>

            {error && (
                <div className="mb-6 p-4 rounded-lg bg-red-50 border border-red-200 text-red-700">
                    {error}
                </div>
            )}

            {loading ? (
                <div className="flex justify-center py-12">
                    <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-indigo-600"></div>
                </div>
            ) : (
                <Card className="overflow-hidden" noPadding>
                    <table className="w-full text-left border-collapse">
                        <thead>
                            <tr className="bg-slate-50 border-b border-slate-200 text-xs uppercase text-slate-500 font-semibold">
                                <th className="px-6 py-4">ID</th>
                                <th className="px-6 py-4">Evaluator</th>
                                <th className="px-6 py-4">Course</th>
                                <th className="px-6 py-4">Semester</th>
                                <th className="px-6 py-4">Status</th>
                                <th className="px-6 py-4">Date</th>
                            </tr>
                        </thead>
                        <tbody className="divide-y divide-slate-100">
                            {history.map((req) => (
                                <tr key={req.id} className="hover:bg-slate-50/50">
                                    <td className="px-6 py-4 font-mono text-xs text-slate-500">#{req.id}</td>
                                    <td className="px-6 py-4 font-medium text-slate-700">{req.evaluator_id}</td>
                                    <td className="px-6 py-4 text-slate-600">{req.course_id}</td>
                                    <td className="px-6 py-4 text-slate-600">{req.semester}</td>
                                    <td className="px-6 py-4">
                                        <span className={`inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium ${
                                            req.status === 'approved' 
                                                ? 'bg-emerald-100 text-emerald-800' 
                                                : 'bg-red-100 text-red-800'
                                        }`}>
                                            {req.status}
                                        </span>
                                    </td>
                                    <td className="px-6 py-4 text-slate-500 text-sm">
                                        {new Date(req.created_at).toLocaleDateString()}
                                    </td>
                                </tr>
                            ))}
                            {history.length === 0 && (
                                <tr>
                                    <td colSpan={6} className="px-6 py-12 text-center text-slate-400">
                                        No history found
                                    </td>
                                </tr>
                            )}
                        </tbody>
                    </table>
                </Card>
            )}
        </div>
    );
};

export default ApproveRequest;
