import React, { useEffect, useState } from "react";
import { useAuth } from "../../context/AuthContext";
import { getRequestHistory } from "../../api/evaluator";
import { EvaluationRequest } from "../../types/evaluator";
import Card from "../../components/Card";

const RequestHistory: React.FC = () => {
    const { user } = useAuth();
    const [requests, setRequests] = useState<EvaluationRequest[]>([]);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState<string | null>(null);

    useEffect(() => {
        const fetchRequests = async () => {
            if (!user?.user_id) return;
            try {
                const data = await getRequestHistory(user.user_id);
                setRequests(data);
            } catch (err) {
                console.error("Failed to fetch history", err);
                setError("Failed to load request history.");
            } finally {
                setLoading(false);
            }
        };

        fetchRequests();
    }, [user?.user_id]);

    if (loading) return <div className="p-6 text-center">Loading history...</div>;
    if (error) return <div className="p-6 text-center text-red-500">{error}</div>;

    return (
        <div className="space-y-6">
            <h1 className="text-2xl font-bold text-gray-800">Request History</h1>

            {requests.length === 0 ? (
                <div className="text-center py-10 bg-gray-50 rounded-lg border border-gray-200 text-gray-500">
                    No requests found.
                </div>
            ) : (
                <div className="space-y-4">
                    {requests.map((req) => (
                        <Card key={req.id} className={`border-l-4 ${req.status === 'approved' ? 'border-l-green-500' :
                                req.status === 'rejected' ? 'border-l-red-500' :
                                    'border-l-yellow-500'
                            }`}>
                            <div className="flex justify-between items-start">
                                <div>
                                    <h3 className="text-lg font-semibold text-gray-800">
                                        {req.course_id} - Semester {req.semester}
                                    </h3>
                                    <p className="text-sm text-gray-600">Year: {req.academic_year}</p>
                                    <p className="text-sm text-gray-500 mt-1">{req.description}</p>
                                </div>
                                <div className="text-right">
                                    <span className={`px-2 py-1 rounded text-xs font-semibold ${req.status === 'approved' ? 'bg-green-100 text-green-800' :
                                            req.status === 'rejected' ? 'bg-red-100 text-red-800' :
                                                'bg-yellow-100 text-yellow-800'
                                        }`}>
                                        {req.status.toUpperCase()}
                                    </span>
                                    <p className="text-xs text-gray-400 mt-2">
                                        {new Date(req.created_at).toLocaleDateString()}
                                    </p>
                                </div>
                            </div>
                        </Card>
                    ))}
                </div>
            )}
        </div>
    );
};

export default RequestHistory;
