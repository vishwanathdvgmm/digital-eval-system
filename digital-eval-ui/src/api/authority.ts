import { RequestRow, ApprovePayload, ReleasePayload, ReleaseResponse } from '../types/authority';

const API_BASE_URL = 'http://127.0.0.1:8443/api/v1';

const getAuthHeaders = () => {
  const token = localStorage.getItem('token');
  return {
    'Content-Type': 'application/json',
    'Authorization': `Bearer ${token}`,
  };
};

export const fetchPendingRequests = async (): Promise<RequestRow[]> => {
  const response = await fetch(`${API_BASE_URL}/authority/requests/pending`, {
    headers: getAuthHeaders(),
  });
  if (!response.ok) {
    throw new Error('Failed to fetch pending requests');
  }
  return response.json();
};

export const fetchRequestHistory = async (): Promise<RequestRow[]> => {
  const response = await fetch(`${API_BASE_URL}/authority/requests/history`, {
    headers: getAuthHeaders(),
  });
  if (!response.ok) {
    throw new Error('Failed to fetch request history');
  }
  return response.json();
};

export const approveRequest = async (id: number, payload: ApprovePayload): Promise<{ assigned: number }> => {
  const response = await fetch(`${API_BASE_URL}/authority/requests/${id}/approve`, {
    method: 'POST',
    headers: getAuthHeaders(),
    body: JSON.stringify(payload),
  });
  if (!response.ok) {
    throw new Error('Failed to approve request');
  }
  return response.json();
};

export const rejectRequest = async (id: number): Promise<{ status: string }> => {
  const response = await fetch(`${API_BASE_URL}/authority/requests/${id}/reject`, {
    method: 'POST',
    headers: getAuthHeaders(),
  });
  if (!response.ok) {
    throw new Error('Failed to reject request');
  }
  return response.json();
};

export const releaseResults = async (payload: ReleasePayload): Promise<ReleaseResponse> => {
  const response = await fetch(`${API_BASE_URL}/authority/results/release`, {
    method: 'POST',
    headers: getAuthHeaders(),
    body: JSON.stringify(payload),
  });
  if (!response.ok) {
    const errorText = await response.text();
    throw new Error(`Failed to release results: ${errorText}`);
  }
  return response.json();
};