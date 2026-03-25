import { render, screen, waitFor } from '@testing-library/angular';
import { provideRouter, ActivatedRoute, convertToParamMap } from '@angular/router';
import { provideLocationMocks } from '@angular/common/testing';
import { provideTanStackQuery, QueryClient } from '@tanstack/angular-query-experimental';
import { of } from 'rxjs';
import { EquipmentDetail } from './equipment-detail';
import { CompendiumApiClient } from '$core/api-clients/compendium-api-client';
import { UserApiClient } from '$core/api-clients/user-api-client';
import { provideTranslocoForTest } from '$core/testing/transloco-testing';
import { Equipment, EquipmentRelationship } from '$generated/models';

const EQUIPMENT: Equipment = {
  id: 1,
  createdAt: '',
  updatedAt: '',
  name: 'barbell',
  displayName: 'Barbell',
  description: 'A standard barbell',
  category: 'free_weights',
  owner: 'seed',
  public: true,
  version: 1,
};

const USER_EQUIPMENT: Equipment = {
  id: 10,
  createdAt: '',
  updatedAt: '',
  name: 'barbell',
  displayName: 'Barbell',
  description: 'A standard barbell',
  category: 'free_weights',
  owner: 'anon',
  public: false,
  version: 1,
};

function setup(forkedRelationships: EquipmentRelationship[] = []) {
  const compendiumApi: Partial<CompendiumApiClient> = {
    fetchEquipmentItem: vi.fn().mockResolvedValue(EQUIPMENT),
    fetchEquipmentVersions: vi.fn().mockResolvedValue([EQUIPMENT]),
    fetchEquipmentPermissions: vi
      .fn()
      .mockResolvedValue({ permissions: ['READ', 'MODIFY', 'DELETE'] }),
    fetchEquipmentRelationships: vi.fn().mockResolvedValue(forkedRelationships),
    deleteEquipment: vi.fn(),
  };
  const userApi: Partial<UserApiClient> = {
    createUserEquipment: vi.fn().mockResolvedValue(USER_EQUIPMENT),
  };

  return {
    compendiumApi,
    userApi,
    providers: [
      provideRouter([]),
      provideLocationMocks(),
      provideTanStackQuery(new QueryClient({ defaultOptions: { queries: { retry: false } } })),
      {
        provide: ActivatedRoute,
        useValue: { paramMap: of(convertToParamMap({ id: '1' })) },
      },
      { provide: CompendiumApiClient, useValue: compendiumApi },
      { provide: UserApiClient, useValue: userApi },
      provideTranslocoForTest(),
    ],
  };
}

function setupWithPermissions(
  permissions: string[],
  forkedRelationships: EquipmentRelationship[] = [],
) {
  const compendiumApi: Partial<CompendiumApiClient> = {
    fetchEquipmentItem: vi.fn().mockResolvedValue(EQUIPMENT),
    fetchEquipmentVersions: vi.fn().mockResolvedValue([EQUIPMENT]),
    fetchEquipmentPermissions: vi.fn().mockResolvedValue({ permissions }),
    fetchEquipmentRelationships: vi.fn().mockResolvedValue(forkedRelationships),
    deleteEquipment: vi.fn(),
  };
  const userApi: Partial<UserApiClient> = {
    createUserEquipment: vi.fn().mockResolvedValue(USER_EQUIPMENT),
  };

  return {
    compendiumApi,
    userApi,
    providers: [
      provideRouter([]),
      provideLocationMocks(),
      provideTanStackQuery(new QueryClient({ defaultOptions: { queries: { retry: false } } })),
      {
        provide: ActivatedRoute,
        useValue: { paramMap: of(convertToParamMap({ id: '1' })) },
      },
      { provide: CompendiumApiClient, useValue: compendiumApi },
      { provide: UserApiClient, useValue: userApi },
      provideTranslocoForTest(),
    ],
  };
}

describe('EquipmentDetail', () => {
  it('shows "compendium.equipment.addToMine" when equipment is not yet added', async () => {
    const { providers } = setup([]);
    await render(EquipmentDetail, { providers });

    await waitFor(() => {
      expect(screen.getByText('compendium.equipment.addToMine')).toBeTruthy();
    });
    expect(screen.queryByText('compendium.equipment.alreadyAdded')).toBeNull();
  });

  it('shows "compendium.equipment.alreadyAdded" link when equipment is already imported', async () => {
    const forkedRel: EquipmentRelationship = {
      id: 1,
      createdAt: '',
      updatedAt: '',
      relationshipType: 'forked',
      strength: 1,
      owner: 'anon',
      fromEquipmentId: 10,
      toEquipmentId: 1,
    };
    const { providers } = setup([forkedRel]);
    await render(EquipmentDetail, { providers });

    await waitFor(() => {
      expect(screen.getByText('compendium.equipment.alreadyAdded')).toBeTruthy();
    });
    expect(screen.queryByText('compendium.equipment.addToMine')).toBeNull();

    const link = screen.getByText('compendium.equipment.alreadyAdded');
    expect(link.getAttribute('href')).toBe('/user/equipment/10');
  });

  it('shows edit button when user has MODIFY permission', async () => {
    const { providers } = setupWithPermissions(['READ', 'MODIFY', 'DELETE']);
    await render(EquipmentDetail, { providers });

    await waitFor(() => {
      expect(screen.getByText('common.edit')).toBeTruthy();
    });
  });

  it('hides edit button when user lacks MODIFY permission', async () => {
    const { providers } = setupWithPermissions(['READ']);
    await render(EquipmentDetail, { providers });

    await waitFor(() => {
      expect(screen.getByText('Barbell')).toBeTruthy();
    });
    expect(screen.queryByText('common.edit')).toBeNull();
  });

  it('shows delete button when user has DELETE permission', async () => {
    const { providers } = setupWithPermissions(['READ', 'MODIFY', 'DELETE']);
    await render(EquipmentDetail, { providers });

    await waitFor(() => {
      expect(screen.getByText('common.delete')).toBeTruthy();
    });
  });

  it('hides delete button when user lacks DELETE permission', async () => {
    const { providers } = setupWithPermissions(['READ']);
    await render(EquipmentDetail, { providers });

    await waitFor(() => {
      expect(screen.getByText('Barbell')).toBeTruthy();
    });
    expect(screen.queryByText('common.delete')).toBeNull();
  });
});
