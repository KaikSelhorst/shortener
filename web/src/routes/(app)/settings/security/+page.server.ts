import { error, fail } from "@sveltejs/kit";
import { apiFetch } from "$lib/server/api";
import type { Actions, PageServerLoad } from "./$types";

interface Me {
  id: number;
  email: string;
  totp_enabled: boolean;
  created_at: string;
}

export const load: PageServerLoad = async (event) => {
  const res = await apiFetch(event, "/auth/me");

  if (!res.ok) error(res.status, "Failed to load account");

  const me: Me = await res.json();
  return { me };
};

export const actions: Actions = {
  setup: async (event) => {
    const res = await apiFetch(event, "/auth/totp/setup", {
      method: "POST",
      headers: { "content-type": "application/json" },
    });
    const body = await res.json().catch(() => ({}));

    if (!res.ok) {
      return fail(res.status, { error: body.error ?? "Failed to start 2FA setup." });
    }

    return { success: true, setup: true, uri: body.uri as string, secret: body.secret as string };
  },

  confirm: async (event) => {
    const data = await event.request.formData();
    const code = data.get("code");

    if (typeof code !== "string" || !code) {
      return fail(400, { error: "Code is required." });
    }

    const res = await apiFetch(event, "/auth/totp/confirm", {
      method: "POST",
      headers: { "content-type": "application/json" },
      body: JSON.stringify({ code }),
    });

    if (!res.ok) {
      const body = await res.json().catch(() => ({}));
      return fail(res.status, { error: body.error ?? "Invalid code." });
    }

    return { success: true, confirmed: true };
  },

  disable: async (event) => {
    const data = await event.request.formData();
    const code = data.get("code");

    if (typeof code !== "string" || !code) {
      return fail(400, { error: "Code is required." });
    }

    const res = await apiFetch(event, "/auth/totp", {
      method: "DELETE",
      headers: { "content-type": "application/json" },
      body: JSON.stringify({ code }),
    });

    if (!res.ok) {
      const body = await res.json().catch(() => ({}));
      return fail(res.status, { error: body.error ?? "Invalid code." });
    }

    return { success: true, disabled: true };
  },
};
