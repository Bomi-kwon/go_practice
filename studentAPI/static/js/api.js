// API 호출 관련 함수들을 분리
const API = {
    async getStudent(id) {
        const response = await fetch(`/student/${id}`);
        if (!response.ok) {
            throw new Error('학생 정보를 가져오는데 실패했습니다');
        }
        return response.json();
    },

    async deleteStudent(id) {
        const response = await fetch(`/student/${id}`, {
            method: 'DELETE'
        });
        if (!response.ok) {
            throw new Error('삭제에 실패했습니다');
        }
        return response;
    },

    async addStudent(student) {
        const response = await fetch('/student', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify(student)
        });
        if (!response.ok) {
            throw new Error('학생 추가에 실패했습니다');
        }
        return response;
    },

    async updateStudent(id, student) {
        const response = await fetch(`/student/${id}`, {
            method: 'PUT',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify(student)
        });
        if (!response.ok) {
            throw new Error('학생 정보 수정에 실패했습니다');
        }
        return response;
    }
};

export default API; 