import api, { apiCall } from './client';
import type {
	Student,
	CreateStudentRequest,
	UpdateStudentRequest,
	PaginatedData,
	StudentFilters,
	ImportResult
} from './types';

export async function getStudents(filters?: StudentFilters): Promise<PaginatedData<Student>> {
	const params = new URLSearchParams();
	if (filters?.status) params.set('status', filters.status);
	if (filters?.cohort) params.set('cohort', filters.cohort);
	if (filters?.search) params.set('search', filters.search);
	if (filters?.residence_country_id) params.set('residence_country_id', filters.residence_country_id);
	params.set('limit', String(filters?.limit ?? 20));
	params.set('offset', String(filters?.offset ?? 0));

	return apiCall<PaginatedData<Student>>(api.get(`/students?${params.toString()}`));
}

export async function getStudent(id: string): Promise<Student> {
	return apiCall<Student>(api.get(`/students/${id}`));
}

export async function createStudent(data: CreateStudentRequest): Promise<Student> {
	return apiCall<Student>(api.post('/students', data));
}

export async function updateStudent(id: string, data: UpdateStudentRequest): Promise<Student> {
	return apiCall<Student>(api.put(`/students/${id}`, data));
}

export async function deleteStudent(id: string): Promise<void> {
	await apiCall<void>(api.delete(`/students/${id}`));
}

export async function importStudents(file: File): Promise<ImportResult> {
	const formData = new FormData();
	formData.append('file', file);
	return apiCall<ImportResult>(
		api.post('/students/import', formData, {
			headers: { 'Content-Type': 'multipart/form-data' }
		})
	);
}
