import { fail, redirect } from "@sveltejs/kit";
import { env } from "$env/dynamic/private";
import { ACCESS_TOKEN_COOKIE } from "$lib/server/auth";
import type { Actions, PageServerLoad } from "./$types";

export const load: PageServerLoad = async ({ parent }) => {
  const { projects } = await parent();
  if (projects.length > 0) redirect(307, `/p/${projects[0].slug}/links`);
};

export const actions: Actions = {
  default: async ({ request, cookies, fetch }) => {
    const data = await request.formData();
    const name = data.get("name");

    if (typeof name !== "string" || !name) {
      return fail(400, { error: "Project name is required." });
    }

    const res = await fetch(`${env.API_URL}/projects`, {
      method: "POST",
      headers: {
        "content-type": "application/json",
        authorization: `Bearer ${cookies.get(ACCESS_TOKEN_COOKIE)}`,
      },
      body: JSON.stringify({ name }),
    });
    const body = await res.json();

    if (!res.ok) {
      return fail(res.status, { error: body.error ?? "Something went wrong." });
    }

    redirect(303, `/p/${body.slug}/links`);
  },
};
