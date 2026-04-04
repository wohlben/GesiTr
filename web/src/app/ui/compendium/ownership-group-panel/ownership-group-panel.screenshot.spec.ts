import { render } from '@testing-library/angular';
import { page } from 'vitest/browser';
import { provideTanStackQuery, QueryClient } from '@tanstack/angular-query-experimental';
import { CompendiumApiClient } from '$core/api-clients/compendium-api-client';
import { ownershipGroupKeys } from '$core/query-keys';
import { OwnershipGroupMembership } from '$generated/ownershipgroup';
import { OwnershipGroupPanel } from './ownership-group-panel';

describe('OwnershipGroupPanel screenshots', () => {
  afterEach(() => {
    document.documentElement.classList.remove('dark');
  });

  const GROUP_ID = 42;

  const ownerOnly: OwnershipGroupMembership[] = [
    {
      id: 1,
      groupId: GROUP_ID,
      userId: 'alice',
      role: 'owner',
      createdAt: '2024-01-01T00:00:00Z',
      updatedAt: '2024-01-01T00:00:00Z',
    },
  ];

  const withMembers: OwnershipGroupMembership[] = [
    ...ownerOnly,
    {
      id: 2,
      groupId: GROUP_ID,
      userId: 'bob',
      role: 'member',
      createdAt: '2024-01-01T00:00:00Z',
      updatedAt: '2024-01-01T00:00:00Z',
    },
    {
      id: 3,
      groupId: GROUP_ID,
      userId: 'carol',
      role: 'admin',
      createdAt: '2024-01-01T00:00:00Z',
      updatedAt: '2024-01-01T00:00:00Z',
    },
  ];

  function renderPanel(memberships: OwnershipGroupMembership[]) {
    const queryClient = new QueryClient({
      defaultOptions: { queries: { retry: false } },
    });
    queryClient.setQueryData(ownershipGroupKeys.memberships(GROUP_ID), memberships);

    const api: Partial<CompendiumApiClient> = {
      fetchOwnershipGroupMemberships: vi.fn().mockResolvedValue(memberships),
      createOwnershipGroupMembership: vi.fn(),
      deleteOwnershipGroupMembership: vi.fn(),
    };

    return render(OwnershipGroupPanel, {
      inputs: { ownershipGroupId: GROUP_ID, open: true },
      providers: [
        provideTanStackQuery(queryClient),
        { provide: CompendiumApiClient, useValue: api },
      ],
    });
  }

  it('owner only - light', async () => {
    await renderPanel(ownerOnly);
    const locator = page.elementLocator(document.body);
    await expect(locator).toMatchScreenshot('owner-only-light');
  });

  it('owner only - dark', async () => {
    document.documentElement.classList.add('dark');
    await renderPanel(ownerOnly);
    const locator = page.elementLocator(document.body);
    await expect(locator).toMatchScreenshot('owner-only-dark');
  });

  it('with members - light', async () => {
    await renderPanel(withMembers);
    const locator = page.elementLocator(document.body);
    await expect(locator).toMatchScreenshot('with-members-light');
  });

  it('with members - dark', async () => {
    document.documentElement.classList.add('dark');
    await renderPanel(withMembers);
    const locator = page.elementLocator(document.body);
    await expect(locator).toMatchScreenshot('with-members-dark');
  });
});
