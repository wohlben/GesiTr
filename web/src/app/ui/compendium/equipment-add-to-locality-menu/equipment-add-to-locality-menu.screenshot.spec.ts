import { render } from '@testing-library/angular';
import { page } from 'vitest/browser';
import { provideRouter } from '@angular/router';
import { provideTanStackQuery, QueryClient } from '@tanstack/angular-query-experimental';
import { provideTranslocoForTest } from '$core/testing/transloco-testing';
import { CompendiumApiClient } from '$core/api-clients/compendium-api-client';
import { localityKeys, localityAvailabilityKeys } from '$core/query-keys';
import { Locality, LocalityAvailability } from '$generated/models';
import { EquipmentAddToLocalityMenu } from './equipment-add-to-locality-menu';

const EQUIPMENT_ID = 42;

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

// Equipment is already added to "Home Gym"
const mockAvailabilities: LocalityAvailability[] = [
  {
    id: 10,
    localityId: 1,
    equipmentId: EQUIPMENT_ID,
    available: true,
    owner: 'user',
    createdAt: '2024-01-15T00:00:00Z',
    updatedAt: '2024-01-15T00:00:00Z',
  },
];

describe('EquipmentAddToLocalityMenu screenshots', () => {
  afterEach(() => {
    document.documentElement.classList.remove('dark');
  });

  async function renderMenu() {
    const queryClient = new QueryClient({
      defaultOptions: { queries: { retry: false } },
    });

    queryClient.setQueryData(localityKeys.list({ owner: 'me', limit: 100 }), {
      items: mockLocalities,
      total: mockLocalities.length,
    });

    queryClient.setQueryData(
      localityAvailabilityKeys.list({ equipmentId: EQUIPMENT_ID }),
      mockAvailabilities,
    );

    const result = await render(EquipmentAddToLocalityMenu, {
      inputs: { equipmentId: EQUIPMENT_ID },
      providers: [
        provideTranslocoForTest(),
        provideRouter([]),
        provideTanStackQuery(queryClient),
        {
          provide: CompendiumApiClient,
          useValue: {
            fetchLocalities: vi
              .fn()
              .mockResolvedValue({ items: mockLocalities, total: mockLocalities.length }),
            fetchLocalityAvailabilities: vi.fn().mockResolvedValue(mockAvailabilities),
            createLocalityAvailability: vi.fn(),
            deleteLocalityAvailability: vi.fn(),
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
    await renderMenu();
    const locator = page.elementLocator(document.body);
    await expect(locator).toMatchScreenshot('light');
  });

  it('dark', async () => {
    document.documentElement.classList.add('dark');
    await renderMenu();
    const locator = page.elementLocator(document.body);
    await expect(locator).toMatchScreenshot('dark');
  });
});
