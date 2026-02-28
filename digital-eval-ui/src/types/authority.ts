export interface RequestRow {
	id: number;
	evaluator_id: string;
	course_id: string;
	semester: string;
	academic_year: string;
	description: string;
	status: string;
	created_at: string;
}

export interface ApprovePayload {
	assign_num: number;
}

export interface ReleasePayload {
	semester: string;
	academic_year: string;
	released_by: string;
}

export interface ReleaseResponse {
	block_hash: string;
}
