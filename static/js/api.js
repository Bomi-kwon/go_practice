const API_BASE_URL = 'http://localhost:8080/api/v1';

// API 호출 함수들
const api = {
    // 목록 조회
    async listResources() {
        const response = await fetch(`${API_BASE_URL}/resources`);
        if (!response.ok) throw new Error('리소스 목록 조회 실패');
        return response.json();
    },

    // 단일 리소스 조회
    async getResource(id) {
        const response = await fetch(`${API_BASE_URL}/resources/${id}`);
        if (!response.ok) throw new Error('리소스 조회 실패');
        return response.json();
    },

    // 리소스 생성
    async createResource(data) {
        const response = await fetch(`${API_BASE_URL}/resources`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify(data),
        });
        if (!response.ok) throw new Error('리소스 생성 실패');
        return response.json();
    },

    // 리소스 수정
    async updateResource(id, data) {
        const response = await fetch(`${API_BASE_URL}/resources/${id}`, {
            method: 'PUT',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify(data),
        });
        if (!response.ok) throw new Error('리소스 수정 실패');
        return response.json();
    },

    // 리소스 삭제
    async deleteResource(id) {
        const response = await fetch(`${API_BASE_URL}/resources/${id}`, {
            method: 'DELETE',
        });
        if (!response.ok) throw new Error('리소스 삭제 실패');
        return response.json();
    },
}; 