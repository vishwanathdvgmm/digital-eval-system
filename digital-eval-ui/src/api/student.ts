import api from "./http";

export interface EvaluationRow {
  ID: number;
  ScriptID: string;
  Evaluator: string;
  CourseID: string;
  Semester: string;
  AcademicYear: string;
  TotalMarks: number;
  Result: string;
  CreatedAt: string;
  Marks: any; // The backend returns the marks object directly
  CourseCredits?: { Int32: number; Valid: boolean };
}

export interface StudentResultResponse {
  usn: string;
  semester: string;
  academic_year: string;
  rows: EvaluationRow[];
  sgpa: number;
}

export async function fetchStudentResults(
  usn: string,
  semester: string,
  academicYear: string
): Promise<StudentResultResponse> {
  const response = await api.get<StudentResultResponse>("/student/results", {
    params: { usn, semester, academic_year: academicYear },
  });
  return response.data;
}

export async function downloadResultPDF(
  usn: string,
  semester: string,
  academicYear: string
): Promise<void> {
  try {
    console.log("Requesting PDF download for:", { usn, semester, academicYear });

    const response = await api.get("/student/download", {
      params: { usn, semester, academic_year: academicYear },
      responseType: "blob", // Important for file download
    });

    console.log("PDF download response:", response);

    // Check if the response is actually a PDF
    const contentType = response.headers['content-type'];
    console.log("Content-Type:", contentType);

    // If we got an error response as JSON instead of PDF
    if (contentType?.includes('application/json')) {
      const text = await response.data.text();
      const error = JSON.parse(text);
      throw new Error(error.message || "Server returned an error instead of PDF");
    }

    // Create a blob link to trigger download
    const url = window.URL.createObjectURL(new Blob([response.data], { type: 'application/pdf' }));
    const link = document.createElement("a");
    link.href = url;
    link.setAttribute("download", `result_${usn}_${semester}.pdf`);
    document.body.appendChild(link);
    link.click();
    link.remove();

    // Clean up the URL object
    window.URL.revokeObjectURL(url);
  } catch (error: any) {
    console.error("Download PDF error details:", error);

    // If the error response is a blob, try to read it as text
    if (error.response?.data instanceof Blob) {
      const text = await error.response.data.text();
      console.error("Error response body:", text);
      throw new Error(`Download failed: ${text}`);
    }

    throw error;
  }
}
