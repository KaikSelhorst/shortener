// A fixed locale (not the runtime default) so the format doesn't depend on
// whether this renders on the server or in the visitor's browser, or on
// either one's own locale — every visitor sees the same unambiguous date.
export function formatDate(value: string | null | undefined): string {
  if (!value) return "—";
  return new Date(value).toLocaleDateString("en-GB", { day: "2-digit", month: "short", year: "numeric" });
}
