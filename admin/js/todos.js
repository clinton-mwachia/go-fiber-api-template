async function loadTodos() {
  const content = document.getElementById("content");
  content.innerHTML = "<h3>Todos</h3><p>Loading...</p>";

  try {
    const todos = await apiGet("/todos");
    let html = `
      <h3>Todos</h3>
      <table class="table table-striped">
        <thead>
          <tr><th>ID</th><th>Title</th><th>Completed</th><th>User ID</th><th>Actions</th></tr>
        </thead><tbody>`;

    todos.forEach((t) => {
      html += `
        <tr>
          <td>${t._id || "N/A"}</td>
          <td>${t.title}</td>
          <td>${t.completed ? "✅" : "❌"}</td>
          <td>${t.user_id || "N/A"}</td>
          <td>
            <button class="btn btn-danger btn-sm" onclick="deleteTodo('${
              t._id
            }')">Delete</button>
          </td>
        </tr>`;
    });

    html += "</tbody></table>";
    content.innerHTML = html;
  } catch (err) {
    content.innerHTML = `<p class="text-danger">Error: ${err.message}</p>`;
  }
}

async function deleteTodo(id) {
  if (!confirm("Delete this todo?")) return;
  await apiDelete(`/todos/${id}`);
  loadTodos();
}
