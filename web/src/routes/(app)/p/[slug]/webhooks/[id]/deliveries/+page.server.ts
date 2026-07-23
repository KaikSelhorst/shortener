import { error } from "@sveltejs/kit";
import { apiFetch } from "$lib/server/api";
import type { PageServerLoad } from "./$types";

interface Webhook {
  id: string;
  project_id: number;
  url: string;
  events: string[];
  enabled: boolean;
  created_at: string;
}

interface Delivery {
  id: string;
  webhook_id: string;
  event: string;
  payload: unknown;
  status: string;
  attempts: number;
  response_status: number | null;
  next_retry_at: string;
  created_at: string;
}

interface DeliveriesResponse {
  data: Delivery[];
  has_more: boolean;
  page: number;
}

export const load: PageServerLoad = async (event) => {
  const page = Math.max(1, Number(event.url.searchParams.get("page")) || 1);

  const [webhooksRes, deliveriesRes] = await Promise.all([
    apiFetch(event, `/projects/${event.params.slug}/webhooks`),
    apiFetch(event, `/projects/${event.params.slug}/webhooks/${event.params.id}/deliveries?page=${page}`),
  ]);

  if (!webhooksRes.ok) error(webhooksRes.status, "Failed to load webhook");
  if (!deliveriesRes.ok) error(deliveriesRes.status, "Failed to load deliveries");

  const webhooks: Webhook[] = await webhooksRes.json();
  const webhook = webhooks.find((w) => w.id === event.params.id);
  if (!webhook) error(404, "Webhook not found");

  const deliveries: DeliveriesResponse = await deliveriesRes.json();
  return { webhook, deliveries };
};
