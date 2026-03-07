import { render } from '@testing-library/angular';
import { page } from 'vitest/browser';
import { LoadingSpinner } from './loading-spinner';

describe('LoadingSpinner screenshots', () => {
  afterEach(() => {
    document.documentElement.style.colorScheme = '';
  });

  it('light', async () => {
    const { fixture } = await render(LoadingSpinner);
    const locator = page.elementLocator(fixture.nativeElement);
    await expect(locator).toMatchScreenshot('light');
  });

  it('dark', async () => {
    document.documentElement.style.colorScheme = 'dark';
    const { fixture } = await render(LoadingSpinner);
    const locator = page.elementLocator(fixture.nativeElement);
    await expect(locator).toMatchScreenshot('dark');
  });
});
