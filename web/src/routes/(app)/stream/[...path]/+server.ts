import { error } from "@sveltejs/kit";
import { apiFetch } from "$lib/server/api";
import type { RequestHandler } from "./$types";

// EventSource can't send an Authorization header, so every SSE stream in the
// app connects here instead: this attaches the access token and forwards to
// the matching API path. The API itself still enforces per-endpoint
// authorization, so proxying any path through is safe.
export const GET: RequestHandler = async (event) => {
  const res = await apiFetch(event, `/${event.params.path}`, { signal: event.request.signal });

  if (!res.ok || !res.body) error(res.status, "Failed to open stream");

  return new Response(res.body, {
    headers: {
      "content-type": "text/event-stream",
      "cache-control": "no-cache",
      connection: "keep-alive",
    },
  });
};
