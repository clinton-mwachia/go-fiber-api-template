async function loadUsers() {
  const content = document.getElementById("content");
  content.innerHTML = "<h3>Users</h3><p>Loading...</p>";

  try {
    const users = await apiGet("/users");
    let html = `
      <h3>Users</h3>
      <table class="table table-striped">
        <thead>
          <tr><th>ID</th><th>Name</th><th>Email</th><th>Actions</th></tr>
        </thead><tbody>`;

    users.forEach((u) => {
      html += `
        <tr>
          <td>${u._id || "N/A"}</td>
          <td>${u.name}</td>
          <td>${u.email}</td>
          <td>
            <button class="btn btn-danger btn-sm" onclick="deleteUser('${
              u._id
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

async function deleteUser(id) {
  if (!confirm("Delete this user?")) return;
  await apiDelete(`/users/${id}`);
  loadUsers();
}
