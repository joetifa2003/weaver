import { type Config } from "prettier";

const config: Config = {
  overrides: [
    {
      files: ["slides.md", "pages/**/*.md"],
      options: {
        parser: "slidev",
        plugins: ["prettier-plugin-slidev"],
      },
    },
  ],
};

export default config;
