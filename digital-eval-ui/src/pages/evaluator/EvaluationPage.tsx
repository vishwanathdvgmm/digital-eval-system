import React, { useEffect, useState } from "react";
import { useParams, useNavigate } from "react-router-dom";
import { useAuth } from "../../hooks/useAuth";
import { getScriptMetadata, submitEvaluation, uploadEvaluatedScript } from "../../api/evaluator";
import { ScriptMetadata } from "../../types/evaluator";
import Button from "../../components/Button";
import Card from "../../components/Card";

const EvaluationPage: React.FC = () => {
  const { scriptId } = useParams<{ scriptId: string }>();
  const { user } = useAuth();
  const navigate = useNavigate();

  const [metadata, setMetadata] = useState<ScriptMetadata | null>(null);
  const [loading, setLoading] = useState(true);
  const [submitting, setSubmitting] = useState(false);
  const [error, setError] = useState("");

  // Form State
  const [totalQuestions, setTotalQuestions] = useState(10);
  const [marksPerQuestion, setMarksPerQuestion] = useState(20); // Fixed to 20
  const [totalMarks, setTotalMarks] = useState(100);
  const [marksScored, setMarksScored] = useState<number[]>(Array(10).fill(0));
  const [courseCredits, setCourseCredits] = useState(4);
  const [annotatedFile, setAnnotatedFile] = useState<File | null>(null);

  useEffect(() => {
    if (scriptId) {
      loadMetadata(scriptId);
    }
  }, [scriptId]);

  const loadMetadata = async (id: string) => {
    try {
      const data = await getScriptMetadata(id);
      setMetadata(data);
    } catch (err) {
      setError("Failed to load script metadata");
    } finally {
      setLoading(false);
    }
  };

  const handleMarkChange = (index: number, value: string) => {
    const newMarks = [...marksScored];
    newMarks[index] = parseInt(value) || 0;
    setMarksScored(newMarks);
  };

  const handleFileChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    if (e.target.files && e.target.files[0]) {
      setAnnotatedFile(e.target.files[0]);
    }
  };

  // Calculate Total Score (Best of 2 per Module)
  const calculateTotalScore = () => {
    let sum = 0;
    for (let i = 0; i < 10; i += 2) {
      const m1 = marksScored[i] || 0;
      const m2 = marksScored[i + 1] || 0;
      sum += Math.max(m1, m2);
    }
    return sum;
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!scriptId || !user || !metadata) return;

    // Validation: Min 5 questions attempted
    const attemptedCount = marksScored.filter(m => m > 0).length;
    if (attemptedCount < 5) {
      setError(`Minimum 5 questions must be attempted. You have attempted ${attemptedCount}.`);
      return;
    }

    const calculatedScore = calculateTotalScore();
    if (calculatedScore > 100) {
      setError(`Total score cannot exceed 100. Current: ${calculatedScore}`);
      return;
    }

    setSubmitting(true);
    setError("");

    try {
      // 1. Upload Annotated Script (if provided)
      if (annotatedFile) {
        await uploadEvaluatedScript(annotatedFile);
      }

      // 2. Submit Marks
      const marksAllotted = Array(totalQuestions).fill(marksPerQuestion);

      await submitEvaluation({
        script_id: scriptId,
        evaluator_id: user.user_id,
        total_questions: totalQuestions,
        marks_per_question: marksPerQuestion,
        total_marks: 100,
        questions_answered: attemptedCount,
        marks_allotted_per_question: marksAllotted,
        marks_scored: marksScored,
        course_id: metadata.course_id || metadata.CourseID || "UNKNOWN",
        semester: metadata.semester || metadata.Semester || "1",
        academic_year: "2025-2026",
        course_credits: courseCredits,
        additional_metadata: {},
      });

      alert("Evaluation Submitted Successfully!");
      navigate("/dashboard/evaluator/assigned");
    } catch (err: any) {
      console.error(err);
      setError(err.response?.data || "Submission failed");
    } finally {
      setSubmitting(false);
    }
  };

  if (loading) return (
    <div className="flex items-center justify-center h-screen bg-slate-50">
      <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-indigo-600"></div>
    </div>
  );

  if (!metadata) return (
    <div className="flex items-center justify-center h-screen bg-slate-50 text-slate-500">
      Script not found
    </div>
  );

  const pdfCid = metadata.pdf_cid || metadata.PDFCid || "";
  const pdfUrl = pdfCid ? `http://127.0.0.1:8080/ipfs/${pdfCid}` : "";

  return (
    <div className="flex h-[calc(100vh-64px)] overflow-hidden bg-slate-100">
      {/* Left: PDF Viewer */}
      <div className="w-1/2 h-full flex flex-col border-r border-slate-200 bg-slate-200/50">
        <div className="bg-white px-6 py-3 border-b border-slate-200 flex justify-between items-center shadow-sm z-10">
          <h2 className="font-semibold text-slate-700 flex items-center gap-2">
            <span className="bg-slate-100 px-2 py-1 rounded text-xs font-mono text-slate-500">SCRIPT</span>
            {scriptId}
          </h2>
          <span className="text-xs text-slate-400">PDF Viewer</span>
        </div>
        <div className="flex-1 relative">
          {pdfUrl ? (
            <iframe
              src={pdfUrl}
              className="absolute inset-0 w-full h-full"
              title="Script PDF"
            />
          ) : (
            <div className="flex items-center justify-center h-full text-slate-400">
              PDF not available
            </div>
          )}
        </div>
      </div>

      {/* Right: Evaluation Form */}
      <div className="w-1/2 h-full overflow-y-auto bg-white">
        <div className="max-w-2xl mx-auto p-8">
          <div className="mb-8">
            <h2 className="text-2xl font-bold text-slate-800">Evaluation Form</h2>
            <p className="text-slate-500 mt-1">Enter marks and upload annotated script</p>
          </div>

          {error && (
            <div className="mb-6 p-4 rounded-lg bg-red-50 border border-red-200 text-red-700 text-sm">
              {error}
            </div>
          )}

          <form onSubmit={handleSubmit} className="space-y-8">
            <Card className="bg-slate-50/50">
              <h3 className="text-sm font-semibold text-slate-900 uppercase tracking-wider mb-4">Configuration</h3>
              <div className="grid grid-cols-2 gap-6">
                <div className="space-y-1">
                  <label className="block text-xs font-medium text-slate-500 uppercase">Course Credits</label>
                  <input
                    type="number"
                    value={courseCredits}
                    onChange={(e) => setCourseCredits(parseInt(e.target.value) || 0)}
                    className="w-full px-3 py-2 bg-white border border-slate-300 rounded-lg text-sm focus:ring-2 focus:ring-indigo-500 outline-none"
                  />
                </div>
                <div className="space-y-1">
                  <label className="block text-xs font-medium text-slate-500 uppercase">Marks/Q (Fixed)</label>
                  <input
                    type="number"
                    value={marksPerQuestion}
                    disabled
                    className="w-full px-3 py-2 bg-slate-100 border border-slate-300 rounded-lg text-sm text-slate-500 cursor-not-allowed"
                  />
                </div>
              </div>
            </Card>

            <div>
              <h3 className="text-sm font-semibold text-slate-900 uppercase tracking-wider mb-4">Module-wise Questions</h3>
              <div className="space-y-4">
                {[0, 1, 2, 3, 4].map((moduleIndex) => (
                  <div key={moduleIndex} className="bg-slate-50 p-4 rounded-lg border border-slate-200">
                    <h4 className="text-sm font-bold text-slate-700 mb-3">Module {String(moduleIndex + 1).padStart(2, '0')}</h4>
                    <div className="grid grid-cols-2 gap-4">
                      {[0, 1].map((qOffset) => {
                        const qIndex = moduleIndex * 2 + qOffset;
                        return (
                          <div key={qIndex} className="space-y-1">
                            <label className="block text-xs font-medium text-slate-500">Q.NO {qIndex + 1}</label>
                            <input
                              type="number"
                              min="0"
                              max={marksPerQuestion}
                              value={marksScored[qIndex]}
                              onChange={(e) => handleMarkChange(qIndex, e.target.value)}
                              className={`w-full px-3 py-2 border rounded-lg text-center font-mono text-sm outline-none focus:ring-2 focus:ring-indigo-500 ${marksScored[qIndex] > marksPerQuestion ? "border-red-300 bg-red-50" : "border-slate-300"
                                }`}
                            />
                          </div>
                        );
                      })}
                    </div>
                    <div className="text-right mt-2 text-xs text-slate-500">
                      Module Score (Best of 2): <span className="font-semibold text-indigo-600">
                        {Math.max(marksScored[moduleIndex * 2], marksScored[moduleIndex * 2 + 1])}
                      </span>
                    </div>
                  </div>
                ))}
              </div>
            </div>

            <div className="flex justify-between items-center p-4 bg-indigo-50 rounded-lg border border-indigo-100">
              <span className="text-sm font-medium text-indigo-900">Total Score (Max 100)</span>
              <span className="text-2xl font-bold text-indigo-700">{calculateTotalScore()}</span>
            </div>

            <div className="pt-4 border-t border-slate-100">
              <h3 className="text-sm font-semibold text-slate-900 uppercase tracking-wider mb-4">Annotation</h3>
              <div className="flex items-center justify-center w-full">
                <label className="flex flex-col items-center justify-center w-full h-32 border-2 border-slate-300 border-dashed rounded-xl cursor-pointer bg-slate-50 hover:bg-slate-100 transition-colors">
                  <div className="flex flex-col items-center justify-center pt-5 pb-6">
                    <svg className="w-8 h-8 mb-3 text-slate-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M7 16a4 4 0 01-.88-7.903A5 5 0 1115.9 6L16 6a5 5 0 011 9.9M15 13l-3-3m0 0l-3 3m3-3v12"></path>
                    </svg>
                    <p className="text-sm text-slate-500">
                      <span className="font-semibold">Click to upload</span> annotated PDF
                    </p>
                    <p className="text-xs text-slate-400 mt-1">
                      {annotatedFile ? annotatedFile.name : "PDF only (MAX. 10MB)"}
                    </p>
                  </div>
                  <input type="file" className="hidden" accept="application/pdf" onChange={handleFileChange} />
                </label>
              </div>
            </div>

            <div className="pt-4">
              <Button
                type="submit"
                isLoading={submitting}
                className="w-full py-4 text-lg shadow-xl shadow-indigo-500/20"
              >
                Submit Evaluation
              </Button>
            </div>
          </form>
        </div>
      </div>
    </div>
  );
};

export default EvaluationPage;
