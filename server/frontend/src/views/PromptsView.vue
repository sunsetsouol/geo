<template>
  <div class="prompts-container">
    <div class="header">
      <h1>Prompt Management</h1>
      <button @click="showAddModal = true" class="btn-primary">Add Prompt</button>
    </div>

    <div v-if="loading" class="loading">Loading...</div>
    <div v-else-if="error" class="error">{{ error }}</div>
    <div v-else class="table-wrapper">
      <table class="prompts-table">
        <thead>
          <tr>
            <th>ID</th>
            <th>Category</th>
            <th>Content</th>
            <th>Actions</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="prompt in prompts" :key="prompt.id">
            <td>{{ prompt.id }}</td>
            <td><span class="category-tag">{{ prompt.category }}</span></td>
            <td class="content-cell">{{ prompt.content }}</td>
            <td class="actions">
              <button @click="editPrompt(prompt)" class="btn-edit">Edit</button>
              <button @click="deletePrompt(prompt.id)" class="btn-delete">Delete</button>
            </td>
          </tr>
        </tbody>
      </table>
    </div>

    <!-- Add/Edit Modal -->
    <div v-if="showAddModal || editingPrompt" class="modal-overlay" @click.self="closeModal">
      <div class="modal">
        <div class="modal-header">
          <h2>{{ editingPrompt ? 'Edit Prompt' : 'Add New Prompt' }}</h2>
          <button @click="closeModal" class="btn-close">&times;</button>
        </div>
        <form @submit.prevent="savePrompt">
          <div class="form-group">
            <label>Category</label>
            <input v-model="form.category" placeholder="e.g., default, brand, competitor" required />
          </div>
          <div class="form-group">
            <label>Content</label>
            <textarea v-model="form.content" rows="8" placeholder="Enter prompt content..." required></textarea>
          </div>
          <div class="modal-actions">
            <button type="button" @click="closeModal" class="btn-secondary">Cancel</button>
            <button type="submit" class="btn-primary" :disabled="saving">
              <span v-if="saving" class="spinner"></span>
              {{ saving ? 'Saving...' : 'Save' }}
            </button>
          </div>
        </form>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, reactive } from 'vue';

interface Prompt {
  id: number;
  content: string;
  category: string;
  created_at: string;
  updated_at: string;
}

const prompts = ref<Prompt[]>([]);
const loading = ref(true);
const error = ref('');
const saving = ref(false);
const showAddModal = ref(false);
const editingPrompt = ref<Prompt | null>(null);

const form = reactive({
  content: '',
  category: 'default'
});

const fetchPrompts = async () => {
  loading.value = true;
  try {
    const res = await fetch('/api/prompts');
    if (!res.ok) throw new Error('Failed to fetch prompts');
    prompts.value = await res.json();
  } catch (err: any) {
    error.value = err.message;
  } finally {
    loading.value = false;
  }
};

const savePrompt = async () => {
  saving.value = true;
  const url = editingPrompt.value ? `/api/prompts/${editingPrompt.value.id}` : '/api/prompts';
  const method = editingPrompt.value ? 'PUT' : 'POST';

  try {
    const res = await fetch(url, {
      method,
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(form)
    });
    if (!res.ok) throw new Error('Failed to save prompt');
    
    // 关闭弹窗
    closeModal();
    
    // 自动重新查询全部 prompt
    await fetchPrompts();
  } catch (err: any) {
    alert(err.message);
  } finally {
    saving.value = false;
  }
};

const editPrompt = (prompt: Prompt) => {
  editingPrompt.value = prompt;
  form.content = prompt.content;
  form.category = prompt.category;
};

const deletePrompt = async (id: number) => {
  if (!confirm('Are you sure you want to delete this prompt? This action cannot be undone.')) return;
  try {
    const res = await fetch(`/api/prompts/${id}`, { method: 'DELETE' });
    if (!res.ok) {
      const errorData = await res.json().catch(() => ({}));
      throw new Error(errorData.error || 'Failed to delete prompt');
    }
    // 成功提示（可选）
    console.log('Prompt deleted successfully');
    await fetchPrompts();
  } catch (err: any) {
    alert('Error: ' + err.message);
  }
};

const closeModal = () => {
  showAddModal.value = false;
  editingPrompt.value = null;
  form.content = '';
  form.category = 'default';
};

onMounted(fetchPrompts);
</script>

<style scoped>
.prompts-container {
  padding: 20px;
  max-width: 1200px;
  margin: 0 auto;
}

.header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}

.table-wrapper {
  background: white;
  border-radius: 8px;
  box-shadow: 0 2px 8px rgba(0,0,0,0.1);
  overflow: hidden;
}

.prompts-table {
  width: 100%;
  border-collapse: collapse;
  text-align: left;
}

.prompts-table th, .prompts-table td {
  padding: 12px 15px;
  border-bottom: 1px solid #eee;
}

.prompts-table th {
  background: #f8f9fa;
  font-weight: 600;
}

.content-cell {
  max-width: 400px;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.category-tag {
  background: #e9ecef;
  padding: 2px 8px;
  border-radius: 4px;
  font-size: 0.85em;
}

.actions {
  display: flex;
  gap: 8px;
}

.btn-primary {
  background: #42b883;
  color: white;
  border: none;
  padding: 8px 16px;
  border-radius: 4px;
  cursor: pointer;
}

.btn-secondary {
  background: #6c757d;
  color: white;
  border: none;
  padding: 8px 16px;
  border-radius: 4px;
  cursor: pointer;
}

.btn-edit {
  background: #3498db;
  color: white;
  border: none;
  padding: 4px 8px;
  border-radius: 4px;
  cursor: pointer;
}

.btn-delete {
  background: #e74c3c;
  color: white;
  border: none;
  padding: 4px 8px;
  border-radius: 4px;
  cursor: pointer;
}

.modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0,0,0,0.5);
  display: flex;
  justify-content: center;
  align-items: center;
}

.modal {
  background: white;
  padding: 24px;
  border-radius: 8px;
  width: 600px;
  max-width: 95%;
  box-shadow: 0 4px 20px rgba(0,0,0,0.2);
}

.modal-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}

.modal-header h2 {
  margin: 0;
}

.btn-close {
  background: none;
  border: none;
  font-size: 24px;
  cursor: pointer;
  color: #999;
}

.btn-close:hover {
  color: #333;
}

.spinner {
  display: inline-block;
  width: 12px;
  height: 12px;
  border: 2px solid rgba(255,255,255,0.3);
  border-radius: 50%;
  border-top-color: #fff;
  animation: spin 1s ease-in-out infinite;
  margin-right: 8px;
}

@keyframes spin {
  to { transform: rotate(360deg); }
}

.form-group {
  margin-bottom: 16px;
}

.form-group label {
  display: block;
  margin-bottom: 4px;
  font-weight: 600;
}

.form-group input, .form-group textarea {
  width: 100%;
  padding: 8px;
  border: 1px solid #ddd;
  border-radius: 4px;
  box-sizing: border-box;
}

.modal-actions {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
  margin-top: 20px;
}
</style>
