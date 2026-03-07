import { render } from '@testing-library/angular';
import { page } from 'vitest/browser';
import { SearchInput } from './search-input';

describe('SearchInput screenshots', () => {
  afterEach(() => {
    document.documentElement.style.colorScheme = '';
  });

  it('default - light', async () => {
    const { fixture } = await render(SearchInput, {
      inputs: { placeholder: 'Search exercises...' },
    });
    const locator = page.elementLocator(fixture.nativeElement);
    await expect(locator).toMatchScreenshot('light');
  });

  it('with value - light', async () => {
    const { fixture } = await render(SearchInput, {
      inputs: { placeholder: 'Search exercises...', value: 'bench press' },
    });
    const locator = page.elementLocator(fixture.nativeElement);
    await expect(locator).toMatchScreenshot('filled-light');
  });

  it('default - dark', async () => {
    document.documentElement.style.colorScheme = 'dark';
    const { fixture } = await render(SearchInput, {
      inputs: { placeholder: 'Search exercises...' },
    });
    const locator = page.elementLocator(fixture.nativeElement);
    await expect(locator).toMatchScreenshot('dark');
    document.documentElement.style.colorScheme = '';
  });

  it('with value - dark', async () => {
    document.documentElement.style.colorScheme = 'dark';
    const { fixture } = await render(SearchInput, {
      inputs: { placeholder: 'Search exercises...', value: 'bench press' },
    });
    const locator = page.elementLocator(fixture.nativeElement);
    await expect(locator).toMatchScreenshot('filled-dark');
  });
});
