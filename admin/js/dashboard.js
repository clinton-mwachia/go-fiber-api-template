async function loadDashboard() {
  const content = document.getElementById("content");
  content.innerHTML = "<h3>Dashboard</h3><p>Loading summary...</p>";

  try {
    const users = await apiGet("/users");
    const todos = await apiGet("/todos");

    content.innerHTML = `
      <h3>Dashboard Overview</h3>
      <div class="row">
        <div class="col-md-6">
          <div class="card text-bg-primary mb-3">
            <div class="card-body">
              <h5 class="card-title">Total Users</h5>
              <p class="card-text fs-3">${users.length}</p>
            </div>
          </div>
        </div>
        <div class="col-md-6">
          <div class="card text-bg-success mb-3">
            <div class="card-body">
              <h5 class="card-title">Total Todos</h5>
              <p class="card-text fs-3">${todos.length}</p>
            </div>
          </div>
        </div>
      </div>`;
  } catch (err) {
    content.innerHTML = `<p class="text-danger">Error loading dashboard: ${err.message}</p>`;
  }
}
