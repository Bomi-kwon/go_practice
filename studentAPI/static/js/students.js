import API from './api.js';
import UI from './ui.js';

let currentStudent = null;

async function showStudentDetail(id, event) {
    if (event.target.classList.contains('delete-btn') || event.target.classList.contains('edit-btn')) {
        return;
    }
    
    try {
        const student = await API.getStudent(id);
        currentStudent = student;
        UI.showStudentDetail(student);
    } catch (error) {
        alert(error.message);
    }
}

async function deleteStudent(id, name) {
    if (!confirm(`${name} 학생을 삭제하시겠습니까?`)) {
        return;
    }
    
    try {
        await API.deleteStudent(id);
        alert('삭제되었습니다');
        window.location.reload();
    } catch (error) {
        alert(error.message);
    }
}

async function openEditModal(id) {
    try {
        const student = await API.getStudent(id);
        currentStudent = student;
        UI.setFormData('edit', student);
        document.getElementById('editStudentId').value = id;
        UI.openModal('editStudentModal');
    } catch (error) {
        alert(error.message);
    }
}

document.addEventListener('DOMContentLoaded', function() {
    // 이벤트 리스너 설정
    setupEventListeners();
});

function setupEventListeners() {
    document.getElementById('addStudentForm').addEventListener('submit', handleAddStudent);
    document.getElementById('editStudentForm').addEventListener('submit', handleEditStudent);
}

async function handleAddStudent(e) {
    e.preventDefault();
    const student = UI.getFormData('');
    UI.closeModal('addStudentModal');

    try {
        await API.addStudent(student);
        alert('학생이 추가되었습니다.');
        window.location.reload();
    } catch (error) {
        console.error('Error:', error);
        alert('오류가 발생했습니다.');
    }
}

async function handleEditStudent(e) {
    e.preventDefault();
    const id = document.getElementById('editStudentId').value;
    const student = UI.getFormData('edit');
    UI.closeModal('editStudentModal');

    try {
        await API.updateStudent(id, student);
        alert('학생 정보가 수정되었습니다.');
        window.location.reload();
    } catch (error) {
        console.error('Error:', error);
        alert('오류가 발생했습니다.');
    }
}

// 전역 함수로 노출
window.showStudentDetail = showStudentDetail;
window.deleteStudent = deleteStudent;
window.openEditModal = openEditModal;
window.openModal = () => UI.openModal('addStudentModal');
window.closeModal = () => UI.closeModal('addStudentModal');
window.closeEditModal = () => UI.closeModal('editStudentModal'); 