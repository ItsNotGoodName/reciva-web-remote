const { generateApi } = require("swagger-typescript-api");
const path = require("path");
const fs = require("fs");

/* NOTE: all fields are optional expect one of `output`, `url`, `spec` */
generateApi({
  name: "api.ts",
  // set to `false` to prevent the tool from writing to disk
  output: path.resolve(process.cwd(), "./src/"),
  input: path.resolve(process.cwd(), "../docs/swagger/swagger.json"),
})
  .then(({ files, configuration }) => {
    const envKey = "VITE_API_PATH";
    const envPath = path.resolve(process.cwd(), ".env");
    const lines = fs
      .readFileSync(envPath)
      .toString()
      .trim()
      .split("\n")
      .map((line) => {
        if (line.startsWith(envKey))
          return `${envKey}=${configuration.apiConfig.baseUrl}\n`;
        return line;
      });
    fs.writeFileSync(envPath, lines.join("\n"));
  })
  .catch((e) => console.error(e));
