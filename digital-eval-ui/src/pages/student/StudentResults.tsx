import React, { useState } from "react";
import { fetchStudentResults, downloadResultPDF, StudentResultResponse } from "../../api/student";
import Button from "../../components/Button";
import Card from "../../components/Card";

const StudentResults: React.FC = () => {
  const [usn, setUsn] = useState("");
  const [semester, setSemester] = useState("5");
  const [academicYear, setAcademicYear] = useState("2025-2026");
  const [result, setResult] = useState<StudentResultResponse | null>(null);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState("");

  const handleSearch = async (e: React.FormEvent) => {
    e.preventDefault();
    setLoading(true);
    setError("");
    setResult(null);

    try {
      const data = await fetchStudentResults(usn, semester, academicYear);
      setResult(data);
    } catch (err) {
      setError("No results found or error fetching results.");
    } finally {
      setLoading(false);
    }
  };

  const handleDownload = async () => {
    try {
      await downloadResultPDF(usn, semester, academicYear);
    } catch (err: any) {
      console.error("PDF download error:", err);
      const errorMessage = err.response?.data?.message || err.message || "Failed to download PDF";
      alert(errorMessage);
    }
  };

  return (
    <div className="max-w-5xl mx-auto space-y-8">
      <div className="text-center">
        <h1 className="text-3xl font-bold text-slate-800">Student Results</h1>
        <p className="text-slate-500 mt-2">View and download your academic performance reports</p>
      </div>

      <Card className="border-t-4 border-t-indigo-500">
        <form onSubmit={handleSearch} className="grid grid-cols-1 md:grid-cols-4 gap-6 items-end">
          <div className="space-y-1">
            <label className="block text-sm font-medium text-slate-700">USN</label>
            <input
              type="text"
              value={usn}
              onChange={(e) => setUsn(e.target.value.toUpperCase())}
              placeholder="e.g. 4BD23AI104"
              className="w-full px-4 py-2 border border-slate-300 rounded-lg focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 outline-none transition-all"
              required
            />
          </div>
          <div className="space-y-1">
            <label className="block text-sm font-medium text-slate-700">Semester</label>
            <select
              value={semester}
              onChange={(e) => setSemester(e.target.value)}
              className="w-full px-4 py-2 border border-slate-300 rounded-lg focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 outline-none transition-all bg-white"
            >
              {[1, 2, 3, 4, 5, 6, 7, 8].map((s) => (
                <option key={s} value={s}>
                  {s}
                </option>
              ))}
            </select>
          </div>
          <div className="space-y-1">
            <label className="block text-sm font-medium text-slate-700">Academic Year</label>
            <input
              type="text"
              value={academicYear}
              onChange={(e) => setAcademicYear(e.target.value)}
              className="w-full px-4 py-2 border border-slate-300 rounded-lg focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 outline-none transition-all"
            />
          </div>
          <Button
            type="submit"
            isLoading={loading}
            className="w-full"
          >
            View Results
          </Button>
        </form>
      </Card>

      {error && (
        <div className="p-4 rounded-lg bg-red-50 border border-red-200 text-red-700 text-center">
          {error}
        </div>
      )}

      {result && (
        <Card className="overflow-hidden border-0 shadow-xl ring-1 ring-slate-900/5" noPadding>
          <div className="p-6 bg-slate-50 border-b border-slate-200 flex justify-between items-center">
            <div>
              <h2 className="text-2xl font-bold text-slate-800">{result.usn}</h2>
              <p className="text-slate-500 text-sm mt-1">
                Semester {result.semester} • {result.academic_year}
              </p>
            </div>
            <div className="text-right">
              <p className="text-xs font-semibold text-slate-400 uppercase tracking-wider">SGPA</p>
              <p className="text-4xl font-bold text-emerald-600">{result.sgpa}</p>
            </div>
          </div>

          <div className="overflow-x-auto">
            <table className="w-full text-left border-collapse">
              <thead>
                <tr className="bg-slate-100/50 border-b border-slate-200 text-xs uppercase text-slate-500 font-semibold tracking-wider">
                  <th className="px-6 py-4">Course ID</th>
                  <th className="px-6 py-4 text-center">Credits</th>
                  <th className="px-6 py-4 text-center">Marks Scored</th>
                  <th className="px-6 py-4 text-center">Total Marks</th>
                  <th className="px-6 py-4 text-center">Result</th>
                </tr>
              </thead>
              <tbody className="divide-y divide-slate-100">
                {result.rows.map((row) => {
                  // Calculate marks scored from the Marks JSON object
                  let marksScored = 0;
                  try {
                    // Handle if Marks is a string (JSON string) or object
                    const marksObj = typeof row.Marks === 'string' ? JSON.parse(row.Marks) : row.Marks;
                    if (marksObj && Array.isArray(marksObj.marks_scored)) {
                      // Calculate score using Best of 2 per Module logic
                      const scores = marksObj.marks_scored;
                      for (let i = 0; i < scores.length; i += 2) {
                        const m1 = scores[i] || 0;
                        const m2 = scores[i + 1] || 0;
                        marksScored += Math.max(m1, m2);
                      }
                    }
                  } catch (e) {
                    console.error("Failed to parse marks", e);
                  }

                  return (
                    <tr key={row.ID} className="hover:bg-slate-50/80 transition-colors">
                      <td className="px-6 py-4 font-medium text-slate-700">{row.CourseID}</td>
                      <td className="px-6 py-4 text-center text-slate-600">
                        {row.CourseCredits?.Valid ? row.CourseCredits.Int32 : "-"}
                      </td>
                      <td className="px-6 py-4 text-center font-semibold text-indigo-600">
                        {marksScored}
                      </td>
                      <td className="px-6 py-4 text-center text-slate-600">{row.TotalMarks}</td>
                      <td className="px-6 py-4 text-center">
                        <span
                          className={`inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium ${row.Result === "PASS"
                            ? "bg-emerald-100 text-emerald-800"
                            : "bg-red-100 text-red-800"
                            }`}
                        >
                          {row.Result}
                        </span>
                      </td>
                    </tr>
                  );
                })}
              </tbody>
            </table>
          </div>

          <div className="p-6 bg-slate-50 border-t border-slate-200 text-center">
            <Button
              onClick={handleDownload}
              variant="secondary"
              className="!rounded-full px-8"
            >
              Download PDF Report
            </Button>
          </div>
        </Card>
      )}
    </div>
  );
};

export default StudentResults;
