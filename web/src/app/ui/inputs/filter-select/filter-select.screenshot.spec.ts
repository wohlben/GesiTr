import { render } from '@testing-library/angular';
import { page } from 'vitest/browser';
import { FilterSelect } from './filter-select';

describe('FilterSelect screenshots', () => {
  afterEach(() => {
    document.documentElement.style.colorScheme = '';
  });

  it('default - light', async () => {
    const { fixture } = await render(FilterSelect, {
      inputs: { allLabel: 'All types', options: ['STRENGTH', 'CARDIO', 'STRETCHING'] },
    });
    const locator = page.elementLocator(fixture.nativeElement);
    await expect(locator).toMatchScreenshot('light');
  });

  it('with selection - light', async () => {
    const { fixture } = await render(FilterSelect, {
      inputs: { allLabel: 'All types', options: ['STRENGTH', 'CARDIO', 'STRETCHING'], value: 'STRENGTH' },
    });
    const locator = page.elementLocator(fixture.nativeElement);
    await expect(locator).toMatchScreenshot('selected-light');
  });

  it('default - dark', async () => {
    document.documentElement.style.colorScheme = 'dark';
    const { fixture } = await render(FilterSelect, {
      inputs: { allLabel: 'All types', options: ['STRENGTH', 'CARDIO', 'STRETCHING'] },
    });
    const locator = page.elementLocator(fixture.nativeElement);
    await expect(locator).toMatchScreenshot('dark');
    document.documentElement.style.colorScheme = '';
  });

  it('with selection - dark', async () => {
    document.documentElement.style.colorScheme = 'dark';
    const { fixture } = await render(FilterSelect, {
      inputs: { allLabel: 'All types', options: ['STRENGTH', 'CARDIO', 'STRETCHING'], value: 'STRENGTH' },
    });
    const locator = page.elementLocator(fixture.nativeElement);
    await expect(locator).toMatchScreenshot('selected-dark');
  });
});
