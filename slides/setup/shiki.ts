/* ./setup/shiki.ts */
import { defineShikiSetup } from "@slidev/types";
import weaver from "./weaver.tmLanguage.json";

export default defineShikiSetup(() => {
  return {
    langs: ["rust", "js", "ts", weaver],
  };
});
