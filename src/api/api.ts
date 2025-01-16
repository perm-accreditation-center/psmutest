import axios from "axios";
import { Test, TestResult } from "../types";

const API_URL = "http://localhost:8080/api";

const api = axios.create({
  baseURL: "http://localhost:8080/api",
  headers: {
    "Content-Type": "application/json",
  },
});

export const getTests = () => api.get<{ tests: Test[] }>(`${API_URL}/tests`);
export const submitTestResult = (data: TestResult) =>
  api.post(`${API_URL}/results`, data);
export const generatePDF = (userId: string, testId: number) =>
  api.get(`${API_URL}/results/${userId}/${testId}/pdf`, {
    responseType: "blob",
  });

export const getTestByID = (testId: number) =>
  api.get<Test>(`${API_URL}/tests/${testId}`);

export const getTestResults = (userId: string, testId: number) =>
  api.get(
    `${API_URL}/results/${userId}/${testId}`
  )