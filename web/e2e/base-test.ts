import { test as base, request as playwrightRequest } from '@playwright/test';

/**
 * Extended test fixture that resets the database before each test.
 * Import { test, expect } from this file instead of '@playwright/test'.
 */
export const test = base.extend({
  page: async ({ page }, use) => {
    const baseURL = process.env['PLAYWRIGHT_TEST_BASE_URL'] ?? 'http://localhost:4200';
    const ctx = await playwrightRequest.newContext({
      baseURL,
      extraHTTPHeaders: { 'X-User-Id': 'devuser' },
    });
    await ctx.post('/api/ci/reset-db');
    await ctx.post('/api/profile');
    await ctx.dispose();
    await use(page);
  },
});

export { expect } from '@playwright/test';
