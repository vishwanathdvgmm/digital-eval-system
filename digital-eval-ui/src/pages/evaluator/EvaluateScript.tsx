import React, { useEffect, useState, useRef } from "react";
import { useParams, useNavigate } from "react-router-dom";
import { useAuth } from "../../context/AuthContext";
import { getScriptMetadata, submitEvaluation } from "../../api/evaluator";
import { ScriptMetadata } from "../../types/evaluator";
import Card from "../../components/Card";

const EvaluateScript: React.FC = () => {
  const { scriptId } = useParams<{ scriptId: string }>();
  const { user } = useAuth();
  const navigate = useNavigate();

  const [metadata, setMetadata] = useState<ScriptMetadata | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [submitting, setSubmitting] = useState(false);

  // Evaluation Form State
  const [totalQuestions, setTotalQuestions] = useState(10); // Fixed to 10 for 5 modules x 2
  const [marksPerQuestion, setMarksPerQuestion] = useState(20); // Default 20 (20x5=100)
  const [marksScored, setMarksScored] = useState<number[]>(Array(10).fill(0));
  const [courseCredits, setCourseCredits] = useState(4);
  const [showScratchpad, setShowScratchpad] = useState(false);

  // Canvas Ref
  const canvasRef = useRef<HTMLCanvasElement>(null);
  const [isDrawing, setIsDrawing] = useState(false);

  useEffect(() => {
    const fetchMetadata = async () => {
      if (!scriptId) return;
      try {
        const data = await getScriptMetadata(scriptId);
        setMetadata(data);
      } catch (err) {
        console.error("Failed to fetch script metadata", err);
        setError("Failed to load script details.");
      } finally {
        setLoading(false);
      }
    };

    fetchMetadata();
  }, [scriptId]);

  // Handle Marks Change
  const handleMarkChange = (index: number, value: string) => {
    const newMarks = [...marksScored];
    newMarks[index] = parseInt(value) || 0;
    setMarksScored(newMarks);
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

  // Handle Submit
  const handleSubmit = async () => {
    if (!scriptId || !user?.user_id || !metadata) return;

    // Validation: Min 5 questions attempted
    const attemptedCount = marksScored.filter(m => m > 0).length;
    if (attemptedCount < 5) {
      alert(`Minimum 5 questions must be attempted. You have attempted ${attemptedCount}.`);
      return;
    }

    const totalScore = calculateTotalScore();
    if (totalScore > 100) {
      alert(`Total score cannot exceed 100. Current: ${totalScore}`);
      return;
    }

    setSubmitting(true);
    try {
      const marksAllotted = Array(totalQuestions).fill(marksPerQuestion);
      
      await submitEvaluation({
        script_id: scriptId,
        evaluator_id: user.user_id,
        total_questions: totalQuestions,
        marks_per_question: marksPerQuestion,
        total_marks: 100, // Fixed max marks
        questions_answered: attemptedCount,
        marks_allotted_per_question: marksAllotted,
        marks_scored: marksScored,
        course_id: metadata.CourseID || "",
        semester: metadata.Semester || "",
        academic_year: metadata.AcademicYear || "2025-2026",
        course_credits: courseCredits,
        additional_metadata: {},
      });

      alert("Evaluation submitted successfully!");
      navigate("/dashboard/evaluator/assigned");
    } catch (err: any) {
      console.error("Submission failed", err);
      alert("Submission failed: " + (err.response?.data?.error || err.message));
    } finally {
      setSubmitting(false);
    }
  };

  // Canvas Drawing Logic
  const startDrawing = (e: React.MouseEvent) => {
    const canvas = canvasRef.current;
    if (!canvas) return;
    const ctx = canvas.getContext("2d");
    if (!ctx) return;

    setIsDrawing(true);
    ctx.beginPath();
    ctx.moveTo(e.nativeEvent.offsetX, e.nativeEvent.offsetY);
  };

  const draw = (e: React.MouseEvent) => {
    if (!isDrawing) return;
    const canvas = canvasRef.current;
    if (!canvas) return;
    const ctx = canvas.getContext("2d");
    if (!ctx) return;

    ctx.lineTo(e.nativeEvent.offsetX, e.nativeEvent.offsetY);
    ctx.stroke();
  };

  const stopDrawing = () => {
    setIsDrawing(false);
  };

  const clearCanvas = () => {
    const canvas = canvasRef.current;
    if (!canvas) return;
    const ctx = canvas.getContext("2d");
    if (!ctx) return;
    ctx.clearRect(0, 0, canvas.width, canvas.height);
  };

  if (loading) return <div className="p-6 text-center">Loading script...</div>;
  if (error) return <div className="p-6 text-center text-red-500">{error}</div>;
  if (!metadata) return <div className="p-6 text-center">No metadata found.</div>;

  const pdfUrl = metadata.pdf_cid ? `http://127.0.0.1:8080/ipfs/${metadata.pdf_cid}` : "";

  return (
    <div className="h-[calc(100vh-100px)] flex flex-col">
      <div className="flex justify-between items-center mb-4 px-4">
        <h1 className="text-xl font-bold text-gray-800">
          Evaluating: {metadata.CourseID} ({metadata.USN})
        </h1>
        <div className="space-x-2">
          <button
            onClick={() => setShowScratchpad(true)}
            className="bg-purple-600 hover:bg-purple-700 text-white px-4 py-2 rounded"
          >
            Open Scratchpad
          </button>
          <button
            onClick={handleSubmit}
            disabled={submitting}
            className={`bg-green-600 hover:bg-green-700 text-white px-4 py-2 rounded ${
              submitting ? "opacity-50 cursor-not-allowed" : ""
            }`}
          >
            {submitting ? "Submitting..." : "Submit Evaluation"}
          </button>
        </div>
      </div>

      <div className="flex-1 flex overflow-hidden border-t border-gray-200">
        {/* Left: PDF Viewer */}
        <div className="w-2/3 bg-gray-100 border-r border-gray-200 relative">
          {pdfUrl ? (
            <iframe
              src={pdfUrl}
              className="w-full h-full"
              title="Script PDF"
            />
          ) : (
            <div className="flex items-center justify-center h-full text-gray-500">
              PDF not available
            </div>
          )}
        </div>

        {/* Right: Evaluation Form */}
        <div className="w-1/3 bg-white p-4 overflow-y-auto">
          <h2 className="text-lg font-semibold mb-4">Marks Entry</h2>
          
          <div className="space-y-4">
            <div className="grid grid-cols-2 gap-4">
              <div>
                <label className="block text-xs font-medium text-gray-500">Course Credits</label>
                <input
                  type="number"
                  value={courseCredits}
                  onChange={(e) => setCourseCredits(parseInt(e.target.value) || 0)}
                  className="mt-1 block w-full border border-gray-300 rounded px-2 py-1"
                />
              </div>
              <div>
                <label className="block text-xs font-medium text-gray-500">Max Marks / Q</label>
                <input
                  type="number"
                  value={marksPerQuestion}
                  onChange={(e) => setMarksPerQuestion(parseInt(e.target.value) || 0)}
                  className="mt-1 block w-full border border-gray-300 rounded px-2 py-1"
                />
              </div>
            </div>

            <div className="border-t border-gray-200 pt-4">
              <h3 className="text-sm font-medium mb-2">Module-wise Questions</h3>
              <div className="space-y-4">
                {[0, 1, 2, 3, 4].map((moduleIndex) => (
                  <div key={moduleIndex} className="bg-gray-50 p-3 rounded border border-gray-200">
                    <h4 className="text-sm font-bold text-gray-700 mb-2">Module {String(moduleIndex + 1).padStart(2, '0')}</h4>
                    <div className="grid grid-cols-2 gap-2">
                      {[0, 1].map((qOffset) => {
                        const qIndex = moduleIndex * 2 + qOffset;
                        return (
                          <div key={qIndex} className="flex flex-col">
                            <label className="text-xs text-gray-500 mb-1">Q.NO {qIndex + 1}</label>
                            <input
                              type="number"
                              min="0"
                              max={marksPerQuestion}
                              value={marksScored[qIndex]}
                              onChange={(e) => handleMarkChange(qIndex, e.target.value)}
                              className={`w-full border rounded px-2 py-1 ${
                                marksScored[qIndex] > marksPerQuestion ? "border-red-500 bg-red-50" : "border-gray-300"
                              }`}
                            />
                          </div>
                        );
                      })}
                    </div>
                    <div className="text-right mt-1 text-xs text-gray-500">
                      Module Score: <span className="font-semibold text-blue-600">
                        {Math.max(marksScored[moduleIndex * 2], marksScored[moduleIndex * 2 + 1])}
                      </span>
                    </div>
                  </div>
                ))}
              </div>
            </div>

            <div className="border-t border-gray-200 pt-4">
              <div className="flex justify-between font-bold text-gray-800 text-lg">
                <span>Total Score:</span>
                <span>{calculateTotalScore()} / 100</span>
              </div>
              <p className="text-xs text-gray-500 mt-1">
                * Score calculated as Best of 2 per Module.
              </p>
            </div>
          </div>
        </div>
      </div>

      {/* Scratchpad Modal */}
      {showScratchpad && (
        <div className="fixed inset-0 bg-black bg-opacity-50 z-50 flex items-center justify-center">
          <div className="bg-white rounded-lg shadow-xl w-[800px] h-[600px] flex flex-col">
            <div className="flex justify-between items-center p-4 border-b">
              <h3 className="font-bold text-lg">Scratchpad</h3>
              <div className="space-x-2">
                <button onClick={clearCanvas} className="text-sm text-red-600 hover:text-red-800">Clear</button>
                <button onClick={() => setShowScratchpad(false)} className="text-gray-500 hover:text-gray-700">Close</button>
              </div>
            </div>
            <div className="flex-1 bg-white relative cursor-crosshair overflow-hidden">
              <canvas
                ref={canvasRef}
                width={800}
                height={540}
                onMouseDown={startDrawing}
                onMouseMove={draw}
                onMouseUp={stopDrawing}
                onMouseLeave={stopDrawing}
                className="w-full h-full"
              />
            </div>
            <div className="p-2 bg-gray-50 text-xs text-gray-500 text-center">
              Use this area for rough calculations or notes. Notes are not saved.
            </div>
          </div>
        </div>
      )}
    </div>
  );
};

export default EvaluateScript;
