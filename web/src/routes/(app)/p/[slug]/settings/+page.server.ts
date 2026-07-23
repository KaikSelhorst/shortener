import { fail, redirect } from "@sveltejs/kit";
import { apiFetch } from "$lib/server/api";
import type { Actions } from "./$types";

export const actions: Actions = {
  rename: async (event) => {
    const data = await event.request.formData();
    const name = data.get("name");

    if (typeof name !== "string" || !name) {
      return fail(400, { error: "Name is required." });
    }

    const res = await apiFetch(event, `/projects/${event.params.slug}`, {
      method: "PUT",
      headers: { "content-type": "application/json" },
      body: JSON.stringify({ name }),
    });
    const body = await res.json();

    if (!res.ok) {
      return fail(res.status, { error: body.error ?? "Something went wrong." });
    }

    redirect(303, `/p/${body.slug}/settings`);
  },

  delete: async (event) => {
    const res = await apiFetch(event, `/projects/${event.params.slug}`, { method: "DELETE" });

    if (!res.ok) {
      const body = await res.json().catch(() => ({}));
      return fail(res.status, { error: body.error ?? "Failed to delete project." });
    }

    redirect(303, "/p");
  },
};
