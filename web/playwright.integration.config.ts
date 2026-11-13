import { defineConfig } from "@playwright/test";
const baseURL = process.env.BLOG_INTEGRATION_URL || "http://127.0.0.1:4176";
const target = new URL(baseURL);
if (
  target.protocol !== "http:" ||
  !["127.0.0.1", "localhost"].includes(target.hostname)
)
  throw new Error(
    "Integration tests only run against a local isolated fixture"
  );
export default defineConfig({
  testDir: "./integration",
  workers: 1,
  timeout: 70000,
  expect: { timeout: 10000 },
  outputDir: "test-results/integration",
  use: {
    baseURL,
    viewport: { width: 1440, height: 1000 },
    launchOptions: {
      executablePath: process.env.PLAYWRIGHT_CHROMIUM_EXECUTABLE_PATH,
    },
  },
});
