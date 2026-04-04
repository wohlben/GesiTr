import { request } from '@playwright/test';

export default async function globalSetup() {
  const baseURL = process.env['PLAYWRIGHT_TEST_BASE_URL'] ?? 'http://localhost:4200';
  const ctx = await request.newContext({
    baseURL,
    extraHTTPHeaders: { 'X-User-Id': 'devuser' },
  });
  try {
    await ctx.post('/api/ci/reset-db');
    await ctx.post('/api/profile');
  } catch {
    // Endpoint may not exist in production mode — DB is fresh in Docker anyway
  }
  await ctx.dispose();
}
