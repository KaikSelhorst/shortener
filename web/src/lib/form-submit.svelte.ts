import type { SubmitFunction } from "@sveltejs/kit";

export function createFormSubmit(onSuccess?: () => void) {
  let submitting = $state(false);

  const submit: SubmitFunction = () => {
    submitting = true;
    return async ({ update, result }) => {
      await update();
      submitting = false;
      if (result.type === "success") onSuccess?.();
    };
  };

  return {
    get submitting() {
      return submitting;
    },
    submit,
  };
}
