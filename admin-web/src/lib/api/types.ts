export interface APIResponse<T> {
	success: boolean;
	message: string;
	data?: T;
	error?: string;
}

export interface PaginatedData<T> {
	items: T[];
	total: number;
	limit: number;
	offset: number;
}

export interface Student {
	id: string;
	full_name: string;
	document_id?: string;
	birth_date: string;
	profile_photo_url?: string;
	nationality_country_id: string;
	residence_country_id: string;
	residence_city_id?: string;
	emails: string[];
	phones?: string[];
	company_id?: string;
	student_code?: string;
	status: 'active' | 'graduated' | 'withdrawn' | 'suspended';
	cohort: string;
	enrollment_date: string;
	graduation_date?: string;
	created_at: string;
	created_by?: string;
	updated_at: string;
	updated_by?: string;
}

export interface CreateStudentRequest {
	full_name: string;
	document_id?: string;
	birth_date: string;
	profile_photo_url?: string;
	nationality_country_id: string;
	residence_country_id: string;
	residence_city_id?: string;
	emails: string[];
	phones?: string[];
	company_id?: string;
	student_code?: string;
	status: string;
	cohort: string;
	enrollment_date: string;
}

export interface UpdateStudentRequest {
	full_name?: string;
	document_id?: string;
	profile_photo_url?: string;
	emails?: string[];
	phones?: string[];
	company_id?: string;
	student_code?: string;
	status?: string;
}

export interface StudentFilters {
	status?: string;
	cohort?: string;
	search?: string;
	residence_country_id?: string;
	limit?: number;
	offset?: number;
}
