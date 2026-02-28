export interface ScriptUploadResponse {
    script_id: string;
    block_hash: string;
    pdf_cid: string;
    pdf_path: string;
    metadata: Record<string, string>;
}

export interface ScriptRecord {
    script_id: string;
    usn: string;
    course_id: string;
    semester: string;
    course_name?: string;
    pdf_cid: string;
    pdf_path: string;
    metadata: Record<string, string>;
    status: 'uploaded' | 'validated' | 'errored';
    created_at: string;
}
