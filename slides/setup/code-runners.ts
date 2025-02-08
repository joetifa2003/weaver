import { defineCodeRunnersSetup } from "@slidev/types";

export default defineCodeRunnersSetup(() => {
  return {
    async weaver(code, ctx) {
      return {
        text: code,
      };
    },
  };
});
