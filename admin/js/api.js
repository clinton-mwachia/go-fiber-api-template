const API_BASE = "http://localhost:8080/api";

// Generic helper for GET requests
async function apiGet(endpoint) {
  const res = await fetch(`${API_BASE}${endpoint}`);
  if (!res.ok) throw new Error("Failed to fetch " + endpoint);
  return await res.json();
}

// Generic helper for DELETE requests
async function apiDelete(endpoint) {
  const res = await fetch(`${API_BASE}${endpoint}`, { method: "DELETE" });
  if (!res.ok) throw new Error("Failed to delete " + endpoint);
  return await res.json();
}
