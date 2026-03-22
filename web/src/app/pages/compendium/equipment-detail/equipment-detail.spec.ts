import { render, screen, waitFor } from '@testing-library/angular';
import { provideRouter, ActivatedRoute, convertToParamMap } from '@angular/router';
import { provideLocationMocks } from '@angular/common/testing';
import { provideTanStackQuery, QueryClient } from '@tanstack/angular-query-experimental';
import { of } from 'rxjs';
import { EquipmentDetail } from './equipment-detail';
import { CompendiumApiClient } from '$core/api-clients/compendium-api-client';
import { UserApiClient } from '$core/api-clients/user-api-client';
import { provideTranslocoForTest } from '$core/testing/transloco-testing';
import { Equipment } from '$generated/models';
import { UserEquipment } from '$generated/user-models';

const EQUIPMENT: Equipment = {
  id: 1,
  createdAt: '',
  updatedAt: '',
  name: 'barbell',
  displayName: 'Barbell',
  description: 'A standard barbell',
  category: 'free_weights',
  templateId: 'tmpl-barbell',
  createdBy: 'seed',
  version: 1,
};

const USER_EQUIPMENT: UserEquipment = {
  id: 10,
  createdAt: '',
  updatedAt: '',
  owner: 'anon',
  compendiumEquipmentId: 'tmpl-barbell',
  compendiumVersion: 1,
};

function setup(userEquipment: UserEquipment[] = []) {
  const compendiumApi: Partial<CompendiumApiClient> = {
    fetchEquipmentItem: vi.fn().mockResolvedValue(EQUIPMENT),
    fetchEquipmentVersions: vi.fn().mockResolvedValue([EQUIPMENT]),
    deleteEquipment: vi.fn(),
  };
  const userApi: Partial<UserApiClient> = {
    fetchUserEquipment: vi.fn().mockResolvedValue(userEquipment),
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
    const { providers } = setup([USER_EQUIPMENT]);
    await render(EquipmentDetail, { providers });

    await waitFor(() => {
      expect(screen.getByText('compendium.equipment.alreadyAdded')).toBeTruthy();
    });
    expect(screen.queryByText('compendium.equipment.addToMine')).toBeNull();

    const link = screen.getByText('compendium.equipment.alreadyAdded');
    expect(link.getAttribute('href')).toBe('/user/equipment/10');
  });

  it('shows "compendium.equipment.addToMine" when user has other equipment but not this one', async () => {
    const otherEquipment: UserEquipment = {
      ...USER_EQUIPMENT,
      id: 99,
      compendiumEquipmentId: 'tmpl-other',
    };
    const { providers } = setup([otherEquipment]);
    await render(EquipmentDetail, { providers });

    await waitFor(() => {
      expect(screen.getByText('compendium.equipment.addToMine')).toBeTruthy();
    });
    expect(screen.queryByText('compendium.equipment.alreadyAdded')).toBeNull();
  });
});
