import { render } from '@testing-library/angular';
import { page } from 'vitest/browser';
import { provideTanStackQuery, QueryClient } from '@tanstack/angular-query-experimental';
import { provideTranslocoForTest } from '$core/testing/transloco-testing';
import { CompendiumApiClient } from '$core/api-clients/compendium-api-client';
import { localityKeys } from '$core/query-keys';
import { Locality } from '$generated/models';
import { LocalityToggleGroup } from './locality-toggle-group';

const mockLocalities: Locality[] = [
  {
    id: 1,
    name: 'Home Gym',
    ownershipGroupId: 0,
    public: false,
    createdAt: '2024-01-01T00:00:00Z',
    updatedAt: '2024-01-01T00:00:00Z',
  },
  {
    id: 2,
    name: 'Commercial Gym',
    ownershipGroupId: 0,
    public: false,
    createdAt: '2024-02-01T00:00:00Z',
    updatedAt: '2024-02-01T00:00:00Z',
  },
];

describe('LocalityToggleGroup screenshots', () => {
  afterEach(() => {
    document.documentElement.classList.remove('dark');
  });

  async function renderGroup(selectedId: number | null = null) {
    const queryClient = new QueryClient({
      defaultOptions: { queries: { retry: false } },
    });

    queryClient.setQueryData(localityKeys.list({ owner: 'me', limit: 100 }), {
      items: mockLocalities,
      total: mockLocalities.length,
    });

    const result = await render(LocalityToggleGroup, {
      providers: [
        provideTranslocoForTest(),
        provideTanStackQuery(queryClient),
        {
          provide: CompendiumApiClient,
          useValue: {
            fetchLocalities: vi
              .fn()
              .mockResolvedValue({ items: mockLocalities, total: mockLocalities.length }),
          },
        },
      ],
    });

    await result.fixture.whenStable();

    if (selectedId !== null) {
      result.fixture.componentInstance.select(selectedId);
      result.fixture.detectChanges();
      await result.fixture.whenStable();
    }

    return result;
  }

  it('all selected - light', async () => {
    const { fixture } = await renderGroup();
    const locator = page.elementLocator(fixture.nativeElement);
    await expect(locator).toMatchScreenshot('all-selected-light');
  });

  it('all selected - dark', async () => {
    document.documentElement.classList.add('dark');
    const { fixture } = await renderGroup();
    const locator = page.elementLocator(fixture.nativeElement);
    await expect(locator).toMatchScreenshot('all-selected-dark');
  });

  it('locality selected - light', async () => {
    const { fixture } = await renderGroup(1);
    const locator = page.elementLocator(fixture.nativeElement);
    await expect(locator).toMatchScreenshot('locality-selected-light');
  });

  it('locality selected - dark', async () => {
    document.documentElement.classList.add('dark');
    const { fixture } = await renderGroup(1);
    const locator = page.elementLocator(fixture.nativeElement);
    await expect(locator).toMatchScreenshot('locality-selected-dark');
  });
});
