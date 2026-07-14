import { fail, redirect } from "@sveltejs/kit";
import { apiFetch } from "$lib/server/api";
import { LAST_PROJECT_COOKIE } from "$lib/server/projects";
import type { Actions, PageServerLoad } from "./$types";

export const load: PageServerLoad = async ({ parent, cookies }) => {
  const { projects } = await parent();
  if (projects.length === 0) return;

  const lastSlug = cookies.get(LAST_PROJECT_COOKIE);
  const target = projects.find((p) => p.slug === lastSlug) ?? projects[0];
  redirect(307, `/p/${target.slug}/links`);
};

export const actions: Actions = {
  default: async (event) => {
    const data = await event.request.formData();
    const name = data.get("name");

    if (typeof name !== "string" || !name) {
      return fail(400, { error: "Project name is required." });
    }

    const res = await apiFetch(event, "/projects", {
      method: "POST",
      headers: { "content-type": "application/json" },
      body: JSON.stringify({ name }),
    });
    const body = await res.json();

    if (!res.ok) {
      return fail(res.status, { error: body.error ?? "Something went wrong." });
    }

    redirect(303, `/p/${body.slug}/links`);
  },
};
