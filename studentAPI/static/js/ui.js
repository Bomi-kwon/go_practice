// UI 조작 관련 함수들을 분리
const UI = {
    showStudentDetail(student) {
        const detailDiv = document.getElementById('studentDetail');
        detailDiv.innerHTML = `
            <h3>${student.Name} 학생의 상세 정보</h3>
            <p>나이: ${student.Age}세</p>
            <p>점수: ${student.Score}점</p>
        `;
        detailDiv.style.display = 'block';
    },

    openModal(modalId) {
        document.getElementById(modalId).style.display = 'block';
    },

    closeModal(modalId) {
        document.getElementById(modalId).style.display = 'none';
    },

    getFormData(formId) {
        return {
            name: document.getElementById(`${formId}Name`).value,
            age: parseInt(document.getElementById(`${formId}Age`).value),
            score: parseInt(document.getElementById(`${formId}Score`).value)
        };
    },

    setFormData(prefix, student) {
        document.getElementById(`${prefix}Name`).value = student.Name;
        document.getElementById(`${prefix}Age`).value = student.Age;
        document.getElementById(`${prefix}Score`).value = student.Score;
    }
};

export default UI; 