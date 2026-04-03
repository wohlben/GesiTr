import { render } from '@testing-library/angular';
import { page } from 'vitest/browser';
import { provideRouter } from '@angular/router';
import { provideTanStackQuery, QueryClient } from '@tanstack/angular-query-experimental';
import { provideTranslocoForTest } from '$core/testing/transloco-testing';
import { CompendiumApiClient } from '$core/api-clients/compendium-api-client';
import { localityKeys } from '$core/query-keys';
import { Locality } from '$generated/models';
import { EquipmentAddToLocalityMenu } from './equipment-add-to-locality-menu';

const mockLocalities: Locality[] = [
  {
    id: 1,
    name: 'Home Gym',
    owner: 'user',
    public: false,
    createdAt: '2024-01-01T00:00:00Z',
    updatedAt: '2024-01-01T00:00:00Z',
  },
  {
    id: 2,
    name: 'Commercial Gym',
    owner: 'user',
    public: false,
    createdAt: '2024-02-01T00:00:00Z',
    updatedAt: '2024-02-01T00:00:00Z',
  },
];

describe('EquipmentAddToLocalityMenu screenshots', () => {
  afterEach(() => {
    document.documentElement.classList.remove('dark');
  });

  async function renderMenu(localities: Locality[]) {
    const queryClient = new QueryClient({
      defaultOptions: { queries: { retry: false } },
    });

    queryClient.setQueryData(localityKeys.list({ owner: 'me', limit: 100 }), {
      items: localities,
      total: localities.length,
    });

    const result = await render(EquipmentAddToLocalityMenu, {
      inputs: { equipmentId: 42 },
      providers: [
        provideTranslocoForTest(),
        provideRouter([]),
        provideTanStackQuery(queryClient),
        {
          provide: CompendiumApiClient,
          useValue: {
            fetchLocalities: vi
              .fn()
              .mockResolvedValue({ items: localities, total: localities.length }),
            createLocalityAvailability: vi.fn(),
          },
        },
      ],
    });

    await result.fixture.whenStable();

    // Open the popover
    result.fixture.componentInstance.open.set(true);
    result.fixture.detectChanges();
    await result.fixture.whenStable();

    return result;
  }

  it('light', async () => {
    await renderMenu(mockLocalities);
    const locator = page.elementLocator(document.body);
    await expect(locator).toMatchScreenshot('light');
  });

  it('dark', async () => {
    document.documentElement.classList.add('dark');
    await renderMenu(mockLocalities);
    const locator = page.elementLocator(document.body);
    await expect(locator).toMatchScreenshot('dark');
  });
});
