import { error } from "@sveltejs/kit";
import { setLastProjectCookie } from "$lib/server/projects";
import type { LayoutServerLoad } from "./$types";

export const load: LayoutServerLoad = async ({ params, parent, cookies }) => {
  const { projects } = await parent();
  const project = projects.find((p) => p.slug === params.slug);
  if (!project) error(404, "Project not found");

  setLastProjectCookie(cookies, project.slug);
  return { project };
};
