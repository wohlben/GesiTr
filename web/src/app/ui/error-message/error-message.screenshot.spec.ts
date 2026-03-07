import { render } from '@testing-library/angular';
import { page } from 'vitest/browser';
import { ErrorMessage } from './error-message';

describe('ErrorMessage screenshots', () => {
  afterEach(() => {
    document.documentElement.classList.remove('dark');
  });

  it('light', async () => {
    const { fixture } = await render(ErrorMessage, {
      inputs: { message: 'Failed to load exercises. Please try again.' },
    });
    const locator = page.elementLocator(fixture.nativeElement);
    await expect(locator).toMatchScreenshot('light');
  });

  it('dark', async () => {
    document.documentElement.classList.add('dark');
    const { fixture } = await render(ErrorMessage, {
      inputs: { message: 'Failed to load exercises. Please try again.' },
    });
    const locator = page.elementLocator(fixture.nativeElement);
    await expect(locator).toMatchScreenshot('dark');
  });
});
