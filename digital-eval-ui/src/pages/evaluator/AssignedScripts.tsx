import React, { useEffect, useState } from "react";
import { useNavigate } from "react-router-dom";
import { useAuth } from "../../hooks/useAuth";
import { getAssignedScripts } from "../../api/evaluator";
import { AssignedScript } from "../../types/evaluator";
import Button from "../../components/Button";
import Card from "../../components/Card";

const AssignedScripts: React.FC = () => {
  const { user } = useAuth();
  const navigate = useNavigate();
  const [scripts, setScripts] = useState<AssignedScript[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");

  useEffect(() => {
    if (user?.user_id) {
      loadScripts(user.user_id);
    }
  }, [user]);

  const loadScripts = async (evaluatorId: string) => {
    try {
      const data = await getAssignedScripts(evaluatorId);
      setScripts(data);
    } catch (err) {
      setError("Failed to load assigned scripts");
    } finally {
      setLoading(false);
    }
  };

  if (loading) return (
    <div className="flex items-center justify-center h-64">
      <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-indigo-600"></div>
    </div>
  );

  if (error) return (
    <div className="p-4 rounded-lg bg-red-50 border border-red-200 text-red-700 text-center">
      {error}
    </div>
  );

  return (
    <div className="max-w-6xl mx-auto space-y-8">
      <div className="flex justify-between items-center">
        <div>
          <h1 className="text-3xl font-bold text-slate-800">Assigned Scripts</h1>
          <p className="text-slate-500 mt-1">Manage and evaluate your assigned student scripts</p>
        </div>
        <div className="bg-indigo-50 text-indigo-700 px-4 py-2 rounded-lg text-sm font-medium">
          Total Assigned: {scripts.length}
        </div>
      </div>

      <Card className="overflow-hidden border-0 shadow-xl ring-1 ring-slate-900/5" noPadding>
        <div className="overflow-x-auto">
          <table className="w-full text-left border-collapse">
            <thead>
              <tr className="bg-slate-50 border-b border-slate-200 text-xs uppercase text-slate-500 font-semibold tracking-wider">
                <th className="px-6 py-4">Script ID</th>
                <th className="px-6 py-4">Course ID</th>
                <th className="px-6 py-4 text-center">Semester</th>
                <th className="px-6 py-4 text-center">Status</th>
                <th className="px-6 py-4 text-right">Action</th>
              </tr>
            </thead>
            <tbody className="divide-y divide-slate-100">
              {scripts.map((script) => (
                <tr key={script.ScriptID} className="hover:bg-slate-50/80 transition-colors">
                  <td className="px-6 py-4 font-medium text-slate-700 font-mono text-xs">{script.ScriptID}</td>
                  <td className="px-6 py-4 text-slate-600">{script.CourseID}</td>
                  <td className="px-6 py-4 text-center text-slate-600">{script.Semester}</td>
                  <td className="px-6 py-4 text-center">
                    <span className={`inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium ${
                      script.Status === "evaluated" 
                        ? "bg-emerald-100 text-emerald-800" 
                        : "bg-amber-100 text-amber-800"
                    }`}>
                      {script.Status}
                    </span>
                  </td>
                  <td className="px-6 py-4 text-right">
                    <Button
                      size="sm"
                      onClick={() => navigate(`/dashboard/evaluator/evaluate/${script.ScriptID}`)}
                      disabled={script.Status === "evaluated"}
                      variant={script.Status === "evaluated" ? "secondary" : "primary"}
                    >
                      {script.Status === "evaluated" ? "Completed" : "Evaluate"}
                    </Button>
                  </td>
                </tr>
              ))}
              {scripts.length === 0 && (
                <tr>
                  <td colSpan={5} className="px-6 py-12 text-center text-slate-400">
                    <p className="text-lg font-medium">No scripts assigned yet</p>
                    <p className="text-sm mt-1">Check back later or request new scripts</p>
                  </td>
                </tr>
              )}
            </tbody>
          </table>
        </div>
      </Card>
    </div>
  );
};

export default AssignedScripts;
