// UI 관련 함수들
async function loadResources() {
    try {
        const response = await api.listResources();
        const tableBody = document.getElementById('resourceTableBody');
        tableBody.innerHTML = '';

        if (response.data && response.data.length > 0) {
            response.data.forEach(resource => {
                const row = createResourceRow(resource);
                tableBody.appendChild(row);
            });
        } else {
            tableBody.innerHTML = '<tr><td colspan="5" style="text-align: center;">리소스가 없습니다.</td></tr>';
        }
    } catch (error) {
        showError(error.message);
    }
}

function formatDate(dateString) {
    const date = new Date(dateString);
    return new Intl.DateTimeFormat('ko-KR', {
        year: 'numeric',
        month: '2-digit',
        day: '2-digit',
        hour: '2-digit',
        minute: '2-digit',
        second: '2-digit'
    }).format(date);
}

function createResourceRow(resource) {
    const tr = document.createElement('tr');
    tr.innerHTML = `
        <td>${resource.id}</td>
        <td>${resource.name || 'Unnamed Resource'}</td>
        <td>${formatDate(resource.created_at)}</td>
        <td>${formatDate(resource.updated_at)}</td>
        <td class="actions">
            <button onclick="editResource(${resource.id})" class="btn btn-warning">수정</button>
            <button onclick="deleteResource(${resource.id})" class="btn btn-danger">삭제</button>
        </td>
    `;
    return tr;
}

async function createResource() {
    const nameInput = document.getElementById('resourceName');
    const name = nameInput.value.trim();
    
    if (!name) {
        showError('리소스 이름을 입력해주세요.');
        return;
    }

    try {
        await api.createResource({ name });
        nameInput.value = '';
        loadResources();
    } catch (error) {
        showError(error.message);
    }
}

async function editResource(id) {
    const newName = prompt('새로운 이름을 입력하세요:');
    if (newName === null) return;

    try {
        await api.updateResource(id, { name: newName });
        loadResources();
    } catch (error) {
        showError(error.message);
    }
}

async function deleteResource(id) {
    if (!confirm('정말 삭제하시겠습니까?')) return;

    try {
        await api.deleteResource(id);
        loadResources();
    } catch (error) {
        showError(error.message);
    }
}

function showError(message) {
    const errorDiv = document.createElement('div');
    errorDiv.className = 'error';
    errorDiv.textContent = message;
    document.body.appendChild(errorDiv);
    setTimeout(() => errorDiv.remove(), 3000);
}

// 페이지 로드 시 리소스 목록 불러오기
document.addEventListener('DOMContentLoaded', loadResources); 