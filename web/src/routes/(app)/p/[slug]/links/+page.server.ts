import { error, fail } from "@sveltejs/kit";
import { apiFetch } from "$lib/server/api";
import type { Actions, PageServerLoad } from "./$types";

interface Link {
  id: number;
  project_id: number;
  short_code: string;
  original_url: string;
  title: string | null;
  description: string | null;
  og_image: string | null;
  expires_at: string | null;
  max_clicks: number | null;
  created_at: string;
  short_url: string;
  total_clicks: number;
}

interface LinksPage {
  data: Link[];
  next_cursor: string | null;
  prev_cursor: string | null;
  limit: number;
}

export const load: PageServerLoad = async (event) => {
  const res = await apiFetch(event, `/projects/${event.params.slug}/links`);

  if (!res.ok) error(res.status, "Failed to load links");

  const links: LinksPage = await res.json();
  return { links };
};

export const actions: Actions = {
  default: async (event) => {
    const data = await event.request.formData();
    const url = data.get("url");
    const title = data.get("title");
    const customCode = data.get("custom_code");

    if (typeof url !== "string" || !url) {
      return fail(400, { error: "URL is required." });
    }

    const body: Record<string, unknown> = { url };
    if (typeof title === "string" && title) body.title = title;
    if (typeof customCode === "string" && customCode) body.custom_code = customCode;

    const res = await apiFetch(event, `/projects/${event.params.slug}/links`, {
      method: "POST",
      headers: { "content-type": "application/json" },
      body: JSON.stringify(body),
    });
    const responseBody = await res.json();

    if (!res.ok) {
      return fail(res.status, { error: responseBody.error ?? "Something went wrong." });
    }

    return { success: true };
  },
};
