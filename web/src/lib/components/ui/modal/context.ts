import { getContext, setContext } from "svelte";

const KEY = Symbol("modal");

export interface ModalContext {
  readonly open: boolean;
  setOpen: (value: boolean) => void;
}

export function setModalContext(ctx: ModalContext) {
  setContext(KEY, ctx);
}

export function getModalContext(): ModalContext {
  const ctx = getContext<ModalContext>(KEY);
  if (!ctx) throw new Error("Modal.Trigger/Content must be used inside Modal.Root");
  return ctx;
}
