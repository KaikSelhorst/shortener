import { error, fail } from "@sveltejs/kit";
import { apiFetch } from "$lib/server/api";
import type { Project } from "$lib/types";
import type { Actions, PageServerLoad } from "./$types";

interface Credential {
  id: number;
  user_id: number;
  project_id: number | null;
  name: string;
  key_prefix: string;
  scopes: string[];
  last_used_at: string | null;
  created_at: string;
}

export const load: PageServerLoad = async (event) => {
  const [keysRes, projectsRes] = await Promise.all([apiFetch(event, "/api-keys"), apiFetch(event, "/projects")]);

  if (!keysRes.ok) error(keysRes.status, "Failed to load credentials");
  if (!projectsRes.ok) error(projectsRes.status, "Failed to load projects");

  const credentials: Credential[] = await keysRes.json();
  const projects: Project[] = await projectsRes.json();
  return { credentials, projects };
};

export const actions: Actions = {
  create: async (event) => {
    const data = await event.request.formData();
    const name = data.get("name");
    const scopes = data.getAll("scopes").filter((value): value is string => typeof value === "string");
    const projectIdRaw = data.get("project_id");

    if (typeof name !== "string" || !name) {
      return fail(400, { error: "Name is required." });
    }
    if (scopes.length === 0) {
      return fail(400, { error: "Select at least one scope." });
    }

    const projectId = typeof projectIdRaw === "string" && projectIdRaw ? Number(projectIdRaw) : undefined;

    const res = await apiFetch(event, "/api-keys", {
      method: "POST",
      headers: { "content-type": "application/json" },
      body: JSON.stringify({ name, scopes, project_id: projectId }),
    });
    const responseBody = await res.json();

    if (!res.ok) {
      return fail(res.status, { error: responseBody.error ?? "Something went wrong." });
    }

    return { success: true, token: responseBody.token as string };
  },

  delete: async (event) => {
    const data = await event.request.formData();
    const id = data.get("id");

    if (typeof id !== "string" || !id) {
      return fail(400, { error: "Missing credential id." });
    }

    const res = await apiFetch(event, `/api-keys/${id}`, { method: "DELETE" });

    if (!res.ok) {
      const body = await res.json().catch(() => ({}));
      return fail(res.status, { error: body.error ?? "Failed to delete credential." });
    }

    return { success: true };
  },
};
