import { getContext, setContext } from "svelte";

const KEY = Symbol("popover");

export interface PopoverContext {
  readonly open: boolean;
  setOpen: (value: boolean) => void;
}

export function setPopoverContext(ctx: PopoverContext) {
  setContext(KEY, ctx);
}

export function getPopoverContext(): PopoverContext {
  const ctx = getContext<PopoverContext>(KEY);
  if (!ctx) throw new Error("Popover.Trigger/Content must be used inside Popover.Root");
  return ctx;
}
