import axios from "axios";
import { getAccessToken } from "../utils/token";

const API_URL = "http://localhost:8443/api/v1/admin";

const getHeaders = () => {
  const token = getAccessToken();
  return {
    Authorization: `Bearer ${token}`,
    "Content-Type": "application/json",
  };
};

export const performServiceAction = async (serviceName: string, action: "start" | "stop" | "restart") => {
  const response = await axios.post(
    `${API_URL}/services/${serviceName}/${action}`,
    {},
    { headers: getHeaders() }
  );
  return response.data;
};

export const getServiceLogs = async (serviceName: string) => {
  const response = await axios.get(
    `${API_URL}/services/${serviceName}/logs`,
    { headers: getHeaders() }
  );
  return response.data.logs as string[];
};
