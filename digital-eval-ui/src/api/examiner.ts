import api from "./http";
import { ScriptUploadResponse, ScriptRecord } from "../types/examiner";

const HISTORY_KEY = "examiner_upload_history";

export async function uploadScript(file: File): Promise<ScriptUploadResponse> {
    const formData = new FormData();
    formData.append("file", file);

    const response = await api.post<ScriptUploadResponse>("/examiner/upload", formData, {
        headers: {
            "Content-Type": "multipart/form-data",
        },
    });

    // Save to local history
    const newRecord: ScriptRecord = {
        script_id: response.data.script_id,
        usn: response.data.metadata.USN || "Unknown",
        course_id: response.data.metadata.CourseID || "Unknown",
        semester: response.data.metadata.Semester || "Unknown",
        course_name: response.data.metadata.CourseName,
        pdf_cid: response.data.pdf_cid,
        pdf_path: response.data.pdf_path,
        metadata: response.data.metadata,
        status: "validated", // Assuming validated if upload succeeds
        created_at: new Date().toISOString(),
    };

    const currentHistory = await getUploadHistory();
    const updatedHistory = [newRecord, ...currentHistory];
    localStorage.setItem(HISTORY_KEY, JSON.stringify(updatedHistory));

    return response.data;
}

export async function getUploadHistory(): Promise<ScriptRecord[]> {
    // Fetch from localStorage to simulate "live" data persistence
    const stored = localStorage.getItem(HISTORY_KEY);
    if (stored) {
        try {
            return JSON.parse(stored);
        } catch (e) {
            console.error("Failed to parse upload history", e);
            return [];
        }
    }
    return [];
}
