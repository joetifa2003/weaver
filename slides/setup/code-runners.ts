import { defineCodeRunnersSetup } from "@slidev/types";

export default defineCodeRunnersSetup(() => {
  return {
    async weaver(code, ctx) {
      const res = await fetch("http://localhost:8080", {
        method: "POST",
        body: code,
        cache: "no-store"
      });

      return {
        text: await res.text(),
      };
    },
  };
});
