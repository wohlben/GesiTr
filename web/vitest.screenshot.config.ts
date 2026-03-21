import { defineConfig } from 'vitest/config';

export default defineConfig({
  test: {
    browser: {
      expect: {
        toMatchScreenshot: {
          comparatorName: 'pixelmatch',
          comparatorOptions: {
            allowedMismatchedPixelRatio: 0.02,
          },
        },
      },
    },
  },
});
