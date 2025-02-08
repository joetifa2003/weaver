import { defineMonacoSetup } from "@slidev/types";

export default defineMonacoSetup(async (monaco) => {
  monaco.languages.register({ id: "weaver" });
  return {
    editorOptions: {
      lineNumbers: "on",
      renderWhitespace: "none",
      experimentalWhitespaceRendering: "off",
      fontLigatures: true,
      automaticLayout: true,
    },
  };
});
