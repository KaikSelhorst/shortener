import { redirect } from "@sveltejs/kit";
import { ACCESS_TOKEN_COOKIE } from "$lib/server/auth";
import type { LayoutServerLoad } from "./$types";

export const load: LayoutServerLoad = ({ cookies }) => {
  if (!cookies.get(ACCESS_TOKEN_COOKIE)) redirect(303, "/login");
};
