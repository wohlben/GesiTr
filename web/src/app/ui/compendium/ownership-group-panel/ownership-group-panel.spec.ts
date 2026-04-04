import { render, screen } from '@testing-library/angular';
import userEvent from '@testing-library/user-event';
import { provideTanStackQuery, QueryClient } from '@tanstack/angular-query-experimental';
import { CompendiumApiClient } from '$core/api-clients/compendium-api-client';
import { ownershipGroupKeys } from '$core/query-keys';
import { OwnershipGroupMembership } from '$generated/ownershipgroup';
import { OwnershipGroupPanel } from './ownership-group-panel';

describe('OwnershipGroupPanel', () => {
  const GROUP_ID = 42;

  const mockMemberships: OwnershipGroupMembership[] = [
    {
      id: 1,
      groupId: GROUP_ID,
      userId: 'alice',
      role: 'owner',
      createdAt: '2024-01-01T00:00:00Z',
      updatedAt: '2024-01-01T00:00:00Z',
    },
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

  function setup(memberships: OwnershipGroupMembership[] = []) {
    const queryClient = new QueryClient({
      defaultOptions: { queries: { retry: false } },
    });
    queryClient.setQueryData(ownershipGroupKeys.memberships(GROUP_ID), memberships);

    const api: Partial<CompendiumApiClient> = {
      fetchOwnershipGroupMemberships: vi.fn().mockResolvedValue(memberships),
      createOwnershipGroupMembership: vi.fn().mockResolvedValue({
        id: 99,
        groupId: GROUP_ID,
        userId: 'dave',
        role: 'member',
        createdAt: '2024-01-01T00:00:00Z',
        updatedAt: '2024-01-01T00:00:00Z',
      }),
      deleteOwnershipGroupMembership: vi.fn().mockResolvedValue(undefined),
    };

    return {
      api,
      render: () =>
        render(OwnershipGroupPanel, {
          inputs: { ownershipGroupId: GROUP_ID, open: true },
          providers: [
            provideTanStackQuery(queryClient),
            { provide: CompendiumApiClient, useValue: api },
          ],
        }),
    };
  }

  it('shows dialog title', async () => {
    const s = setup(mockMemberships);
    await s.render();
    expect(screen.getByText('Shared with')).toBeTruthy();
  });

  it('lists all members including owner', async () => {
    const s = setup(mockMemberships);
    await s.render();
    expect(screen.getByText('alice')).toBeTruthy();
    expect(screen.getByText('bob')).toBeTruthy();
    expect(screen.getByText('carol')).toBeTruthy();
  });

  it('shows role badges for all roles', async () => {
    const s = setup(mockMemberships);
    await s.render();
    expect(screen.getByText('owner')).toBeTruthy();
    expect(screen.getByText('member')).toBeTruthy();
    expect(screen.getByText('admin')).toBeTruthy();
  });

  it('does not show Remove button for owner', async () => {
    const s = setup(mockMemberships);
    await s.render();
    const removeButtons = screen.getAllByText('Remove');
    // only bob and carol get Remove buttons, not alice (owner)
    expect(removeButtons.length).toBe(2);
  });

  it('calls add member mutation on form submit', async () => {
    const user = userEvent.setup();
    const s = setup(mockMemberships);
    await s.render();

    const input = screen.getByPlaceholderText('User ID');
    await user.type(input, 'dave');
    await user.click(screen.getByText('Add'));

    expect(s.api.createOwnershipGroupMembership).toHaveBeenCalledWith(GROUP_ID, {
      userId: 'dave',
    });
  });

  it('calls remove mutation when clicking Remove', async () => {
    const user = userEvent.setup();
    const s = setup(mockMemberships);
    await s.render();

    const removeButtons = screen.getAllByText('Remove');
    await user.click(removeButtons[0]);

    expect(s.api.deleteOwnershipGroupMembership).toHaveBeenCalledWith(2);
  });

  it('disables Add button when input is empty', async () => {
    const s = setup(mockMemberships);
    await s.render();
    const addBtn = screen.getByText('Add');
    expect(addBtn.hasAttribute('disabled') || (addBtn as HTMLButtonElement).disabled).toBe(true);
  });
});
