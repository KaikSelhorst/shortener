import { error } from "@sveltejs/kit";
import { apiFetch } from "$lib/server/api";
import type { Project } from "$lib/types";
import type { LayoutServerLoad } from "./$types";

export const load: LayoutServerLoad = async (event) => {
  const res = await apiFetch(event, "/projects");

  if (!res.ok) error(res.status, "Failed to load projects");

  const projects: Project[] = await res.json();
  return { projects };
};
