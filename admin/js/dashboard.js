const API_BASE = "http://localhost:8080/api";

document.addEventListener("DOMContentLoaded", () => {
  loadStats();
  loadUsers();
  loadTodos();
});

async function loadStats() {
  const usersRes = await fetch(`${API_BASE}/users`);
  const todosRes = await fetch(`${API_BASE}/todos/count`);

  todos = await todosRes.json();

  document.getElementById("totalUsers").textContent = (
    await usersRes.json()
  ).length;
  document.getElementById("totalTodos").textContent = todos.count;
}

// USERS
async function loadUsers() {
  const res = await fetch(`${API_BASE}/users`);
  const users = await res.json();
  const tbody = document.querySelector("#usersTable tbody");
  tbody.innerHTML = "";
  users.forEach((u, i) => {
    tbody.innerHTML += `
      <tr>
        <td>${i + 1}</td>
        <td>${u.username}</td>
        <td>${u.email}</td>
        <td>${u.role}</td>
        <td>
          <button class="btn btn-sm btn-warning" onclick="editUser('${
            u.id
          }', '${u.username}', '${u.email}', '${u.role}')">Edit</button>
          <button class="btn btn-sm btn-danger" onclick="deleteUser('${
            u.id
          }')">Delete</button>
        </td>
      </tr>
    `;
  });
}

document.getElementById("userForm").addEventListener("submit", async (e) => {
  e.preventDefault();
  const id = document.getElementById("userId").value;
  const data = {
    username: document.getElementById("username").value,
    email: document.getElementById("email").value,
    role: document.getElementById("role").value,
  };

  const method = id ? "PUT" : "POST";
  const url = id ? `${API_BASE}/user/${id}` : `${API_BASE}/user/register`;

  await fetch(url, {
    method,
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(data),
  });
  bootstrap.Modal.getInstance(document.getElementById("userModal")).hide();
  loadUsers();
});

function editUser(id, username, email, role) {
  document.getElementById("userId").value = id;
  document.getElementById("username").value = username;
  document.getElementById("email").value = email;
  document.getElementById("role").value = role;
  new bootstrap.Modal(document.getElementById("userModal")).show();
}

async function deleteUser(id) {
  if (confirm("Delete this user?")) {
    await fetch(`${API_BASE}/user/${id}`, { method: "DELETE" });
    loadUsers();
  }
}

// TODOS
async function loadTodos() {
  const res = await fetch(`${API_BASE}/todos`);
  const todos = await res.json();
  const tbody = document.querySelector("#todosTable tbody");
  tbody.innerHTML = "";
  todos.forEach((t, i) => {
    tbody.innerHTML += `
      <tr>
        <td>${i + 1}</td>
        <td>${t.title}</td>
        <td>${t.userId}</td>
        <td>${t.completed ? "✅" : "❌"}</td>
        <td>
          <button class="btn btn-sm btn-warning" onclick="editTodo('${
            t.id
          }', '${t.title}', '${t.userId}', ${t.completed})">Edit</button>
          <button class="btn btn-sm btn-danger" onclick="deleteTodo('${
            t.id
          }')">Delete</button>
        </td>
      </tr>
    `;
  });
}

document.getElementById("todoForm").addEventListener("submit", async (e) => {
  e.preventDefault();
  const id = document.getElementById("todoId").value;
  const data = {
    title: document.getElementById("title").value,
    userId: document.getElementById("todoUserId").value,
    completed: document.getElementById("completed").checked,
  };

  const method = id ? "PUT" : "POST";
  const url = id ? `${API_BASE}/todo/${id}` : `${API_BASE}/todo`;

  await fetch(url, {
    method,
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(data),
  });
  bootstrap.Modal.getInstance(document.getElementById("todoModal")).hide();
  loadTodos();
});

function editTodo(id, title, userId, completed) {
  document.getElementById("todoId").value = id;
  document.getElementById("title").value = title;
  document.getElementById("todoUserId").value = userId;
  document.getElementById("completed").checked = completed;
  new bootstrap.Modal(document.getElementById("todoModal")).show();
}

async function deleteTodo(id) {
  if (confirm("Delete this todo?")) {
    await fetch(`${API_BASE}/todo/${id}`, { method: "DELETE" });
    loadTodos();
  }
}
