import { error } from "@sveltejs/kit";
import { apiFetch } from "$lib/server/api";
import type { PageServerLoad } from "./$types";

const PERIODS = ["7d", "30d", "90d"] as const;
type Period = (typeof PERIODS)[number];

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

interface LinkAnalytics {
  link_id: number;
  short_code: string;
  total_clicks: number;
  unique_clicks: number;
  over_time: { date: string; count: number }[];
  devices: Record<string, number>;
  referrers: Record<string, number>;
  browsers: Record<string, number>;
}

export const load: PageServerLoad = async (event) => {
  const requested = event.url.searchParams.get("period");
  const period: Period = (PERIODS as readonly string[]).includes(requested ?? "") ? (requested as Period) : "30d";

  const [linkRes, analyticsRes] = await Promise.all([
    apiFetch(event, `/projects/${event.params.slug}/links/${event.params.code}`),
    apiFetch(event, `/projects/${event.params.slug}/links/${event.params.code}/analytics?period=${period}`),
  ]);

  if (!linkRes.ok) error(linkRes.status, "Failed to load link");
  if (!analyticsRes.ok) error(analyticsRes.status, "Failed to load analytics");

  const link: Link = await linkRes.json();
  const analytics: LinkAnalytics = await analyticsRes.json();
  return { link, analytics, period, periods: PERIODS };
};
