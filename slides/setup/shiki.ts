/* ./setup/shiki.ts */
import { defineShikiSetup } from "@slidev/types";
import weaver from "./weaver.tmLanguage.json";

export default defineShikiSetup(() => {
  return {
    theme: "kanagawa-dragon",
    langs: ["rust", "js", "ts", weaver],
  };
});
