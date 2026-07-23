import { error, fail } from "@sveltejs/kit";
import { apiFetch } from "$lib/server/api";
import type { Actions, PageServerLoad } from "./$types";

interface Webhook {
  id: string;
  project_id: number;
  url: string;
  events: string[];
  enabled: boolean;
  created_at: string;
}

export const load: PageServerLoad = async (event) => {
  const res = await apiFetch(event, `/projects/${event.params.slug}/webhooks`);

  if (!res.ok) error(res.status, "Failed to load webhooks");

  const webhooks: Webhook[] = await res.json();
  return { webhooks };
};

export const actions: Actions = {
  create: async (event) => {
    const data = await event.request.formData();
    const url = data.get("url");
    const events = data.getAll("events").filter((value): value is string => typeof value === "string");

    if (typeof url !== "string" || !url) {
      return fail(400, { error: "URL is required." });
    }
    if (events.length === 0) {
      return fail(400, { error: "Select at least one event." });
    }

    const res = await apiFetch(event, `/projects/${event.params.slug}/webhooks`, {
      method: "POST",
      headers: { "content-type": "application/json" },
      body: JSON.stringify({ url, events }),
    });
    const responseBody = await res.json();

    if (!res.ok) {
      return fail(res.status, { error: responseBody.error ?? "Something went wrong." });
    }

    return { success: true, secret: responseBody.secret as string };
  },

  delete: async (event) => {
    const data = await event.request.formData();
    const id = data.get("id");

    if (typeof id !== "string" || !id) {
      return fail(400, { error: "Missing webhook id." });
    }

    const res = await apiFetch(event, `/projects/${event.params.slug}/webhooks/${id}`, {
      method: "DELETE",
    });

    if (!res.ok) {
      const body = await res.json().catch(() => ({}));
      return fail(res.status, { error: body.error ?? "Failed to delete webhook." });
    }

    return { success: true };
  },
};
