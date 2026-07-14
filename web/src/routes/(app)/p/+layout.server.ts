import { error, redirect } from "@sveltejs/kit";
import { env } from "$env/dynamic/private";
import { ACCESS_TOKEN_COOKIE } from "$lib/server/auth";
import type { Project } from "$lib/types";
import type { LayoutServerLoad } from "./$types";

export const load: LayoutServerLoad = async ({ cookies, fetch }) => {
  const res = await fetch(`${env.API_URL}/projects`, {
    headers: { authorization: `Bearer ${cookies.get(ACCESS_TOKEN_COOKIE)}` },
  });

  if (res.status === 401) redirect(303, "/login");
  if (!res.ok) error(res.status, "Failed to load projects");

  const projects: Project[] = await res.json();
  return { projects };
};
