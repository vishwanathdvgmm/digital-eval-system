import api from "./http";
import {
  EvaluationRequestCreate,
  EvaluationRequest,
  AssignedScript,
  ScriptMetadata,
  EvaluationSubmit,
  SubmitResponse,
} from "../types/evaluator";

export async function createEvaluationRequest(data: EvaluationRequestCreate): Promise<{ request_id: number }> {
  const response = await api.post<{ request_id: number }>("/evaluator/requests", data);
  return response.data;
}

export async function getRequestHistory(evaluatorId: string): Promise<EvaluationRequest[]> {
  const response = await api.get<EvaluationRequest[]>(`/evaluator/requests/history?evaluator_id=${evaluatorId}`);
  return response.data;
}

export async function getAssignedScripts(evaluatorId: string): Promise<AssignedScript[]> {
  const response = await api.get<AssignedScript[]>(`/evaluator/assigned?evaluator_id=${evaluatorId}`);
  return response.data;
}

export async function getScriptMetadata(scriptId: string): Promise<ScriptMetadata> {
  const response = await api.get<ScriptMetadata>(`/evaluator/script/${scriptId}`);
  return response.data;
}

export async function submitEvaluation(data: EvaluationSubmit): Promise<SubmitResponse> {
  const response = await api.post<SubmitResponse>("/evaluator/submit", data);
  return response.data;
}

export async function uploadEvaluatedScript(file: File): Promise<{ cid: string; pdf_path: string; status: string }> {
  const formData = new FormData();
  formData.append("file", file);
  const response = await api.post<{ cid: string; pdf_path: string; status: string }>("/evaluator/upload", formData, {
    headers: {
      "Content-Type": "multipart/form-data",
    },
  });
  return response.data;
}