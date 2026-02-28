import React, { useState } from "react";
import { useAuth } from "../../context/AuthContext";
import { createEvaluationRequest } from "../../api/evaluator";
import Card from "../../components/Card";

const RequestEvaluation: React.FC = () => {
    const { user } = useAuth();
    const [formData, setFormData] = useState({
        course_id: "",
        semester: "",
        academic_year: "2025-2026",
        description: "",
    });
    const [loading, setLoading] = useState(false);
    const [success, setSuccess] = useState<string | null>(null);
    const [error, setError] = useState<string | null>(null);

    const handleChange = (e: React.ChangeEvent<HTMLInputElement | HTMLTextAreaElement | HTMLSelectElement>) => {
        setFormData({ ...formData, [e.target.name]: e.target.value });
    };

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault();
        if (!user?.user_id) {
            setError("User ID not found. Please relogin.");
            return;
        }

        setLoading(true);
        setError(null);
        setSuccess(null);

        try {
            await createEvaluationRequest({
                evaluator_id: user.user_id,
                ...formData,
            });
            setSuccess("Request submitted successfully!");
            setFormData({
                course_id: "",
                semester: "",
                academic_year: "2025-2026",
                description: "",
            });
        } catch (err: any) {
            console.error("Request failed", err);
            setError(err.response?.data || "Failed to submit request.");
        } finally {
            setLoading(false);
        }
    };

    return (
        <div className="max-w-2xl mx-auto space-y-6">
            <h1 className="text-2xl font-bold text-gray-800">Request Evaluation Access</h1>

            <Card>
                <form onSubmit={handleSubmit} className="space-y-4">
                    <div>
                        <label className="block text-sm font-medium text-gray-700">Evaluator ID</label>
                        <input
                            type="text"
                            value={user?.user_id || ""}
                            disabled
                            className="mt-1 block w-full px-3 py-2 bg-gray-100 border border-gray-300 rounded-md shadow-sm focus:outline-none sm:text-sm"
                        />
                    </div>

                    <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                        <div>
                            <label className="block text-sm font-medium text-gray-700">Course ID</label>
                            <input
                                type="text"
                                name="course_id"
                                value={formData.course_id}
                                onChange={handleChange}
                                required
                                placeholder="e.g. BCS501"
                                className="mt-1 block w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:ring-blue-500 focus:border-blue-500 sm:text-sm"
                            />
                        </div>

                        <div>
                            <label className="block text-sm font-medium text-gray-700">Semester</label>
                            <select
                                name="semester"
                                value={formData.semester}
                                onChange={handleChange}
                                required
                                className="mt-1 block w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:ring-blue-500 focus:border-blue-500 sm:text-sm"
                            >
                                <option value="">Select Semester</option>
                                {[1, 2, 3, 4, 5, 6, 7, 8].map((sem) => (
                                    <option key={sem} value={sem.toString()}>
                                        {sem}
                                    </option>
                                ))}
                            </select>
                        </div>
                    </div>

                    <div>
                        <label className="block text-sm font-medium text-gray-700">Academic Year</label>
                        <input
                            type="text"
                            name="academic_year"
                            value={formData.academic_year}
                            onChange={handleChange}
                            required
                            placeholder="e.g. 2025-2026"
                            className="mt-1 block w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:ring-blue-500 focus:border-blue-500 sm:text-sm"
                        />
                    </div>

                    <div>
                        <label className="block text-sm font-medium text-gray-700">Description / Reason</label>
                        <textarea
                            name="description"
                            value={formData.description}
                            onChange={handleChange}
                            rows={3}
                            placeholder="Briefly explain why you need access..."
                            className="mt-1 block w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:ring-blue-500 focus:border-blue-500 sm:text-sm"
                        />
                    </div>

                    {error && <div className="text-red-600 text-sm">{error}</div>}
                    {success && <div className="text-green-600 text-sm">{success}</div>}

                    <button
                        type="submit"
                        disabled={loading}
                        className={`w-full flex justify-center py-2 px-4 border border-transparent rounded-md shadow-sm text-sm font-medium text-white ${loading ? "bg-gray-400" : "bg-blue-600 hover:bg-blue-700"
                            } focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500`}
                    >
                        {loading ? "Submitting..." : "Submit Request"}
                    </button>
                </form>
            </Card>
        </div>
    );
};

export default RequestEvaluation;
