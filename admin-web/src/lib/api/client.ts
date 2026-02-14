import axios from 'axios';
import type { APIResponse } from './types';

const api = axios.create({
	baseURL: 'http://localhost:8080/api/v1',
	headers: {
		'Content-Type': 'application/json'
	}
});

export async function apiCall<T>(request: Promise<{ data: APIResponse<T> }>): Promise<T> {
	const response = await request;
	if (!response.data.success) {
		throw new Error(response.data.error || response.data.message);
	}
	return response.data.data as T;
}

export default api;
