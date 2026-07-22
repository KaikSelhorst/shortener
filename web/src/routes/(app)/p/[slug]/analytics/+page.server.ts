import { error } from "@sveltejs/kit";
import { apiFetch } from "$lib/server/api";
import type { PageServerLoad } from "./$types";

const PERIODS = ["7d", "30d", "90d"] as const;
type Period = (typeof PERIODS)[number];

interface ProjectAnalytics {
  total_clicks: number;
  unique_clicks: number;
  over_time: { date: string; count: number }[];
  devices: Record<string, number>;
  referrers: Record<string, number>;
  browsers: Record<string, number>;
  top_links: { short_code: string; original_url: string; title: string | null; total_clicks: number }[];
}

export const load: PageServerLoad = async (event) => {
  const requested = event.url.searchParams.get("period");
  const period: Period = (PERIODS as readonly string[]).includes(requested ?? "") ? (requested as Period) : "30d";

  const res = await apiFetch(event, `/projects/${event.params.slug}/analytics?period=${period}`);
  if (!res.ok) error(res.status, "Failed to load analytics");

  const analytics: ProjectAnalytics = await res.json();
  return { analytics, period, periods: PERIODS };
};
