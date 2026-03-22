import { render } from '@testing-library/angular';
import { page } from 'vitest/browser';
import { provideRouter } from '@angular/router';
import { provideTranslocoForTest } from '$core/testing/transloco-testing';
import { Pagination } from './pagination';

describe('Pagination screenshots', () => {
  afterEach(() => {
    document.documentElement.classList.remove('dark');
  });

  const providers = [provideTranslocoForTest(), provideRouter([])];

  const multiPage = { items: new Array(50), total: 874, limit: 50, offset: 0 };
  const midPage = { items: new Array(50), total: 874, limit: 50, offset: 100 };
  const singlePage = { items: new Array(12), total: 12, limit: 50, offset: 0 };
  const empty = { items: [], total: 0, limit: 50, offset: 0 };

  it('first page – light', async () => {
    const { fixture } = await render(Pagination, {
      inputs: { page: multiPage },
      providers,
    });
    const locator = page.elementLocator(fixture.nativeElement);
    await expect(locator).toMatchScreenshot('first-page-light');
  });

  it('first page – dark', async () => {
    document.documentElement.classList.add('dark');
    const { fixture } = await render(Pagination, {
      inputs: { page: multiPage },
      providers,
    });
    const locator = page.elementLocator(fixture.nativeElement);
    await expect(locator).toMatchScreenshot('first-page-dark');
  });

  it('mid page – light', async () => {
    const { fixture } = await render(Pagination, {
      inputs: { page: midPage },
      providers,
    });
    const locator = page.elementLocator(fixture.nativeElement);
    await expect(locator).toMatchScreenshot('mid-page-light');
  });

  it('single page – light', async () => {
    const { fixture } = await render(Pagination, {
      inputs: { page: singlePage },
      providers,
    });
    const locator = page.elementLocator(fixture.nativeElement);
    await expect(locator).toMatchScreenshot('single-page-light');
  });

  it('empty – light', async () => {
    const { fixture } = await render(Pagination, {
      inputs: { page: empty, emptyLabel: 'No exercises found' },
      providers,
    });
    const locator = page.elementLocator(fixture.nativeElement);
    await expect(locator).toMatchScreenshot('empty-light');
  });
});
