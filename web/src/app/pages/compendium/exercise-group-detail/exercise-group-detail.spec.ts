import { render, screen, waitFor } from '@testing-library/angular';
import { provideRouter, ActivatedRoute, convertToParamMap } from '@angular/router';
import { provideLocationMocks } from '@angular/common/testing';
import { provideTanStackQuery, QueryClient } from '@tanstack/angular-query-experimental';
import { of } from 'rxjs';
import { ExerciseGroupDetail } from './exercise-group-detail';
import { CompendiumApiClient } from '$core/api-clients/compendium-api-client';
import { provideTranslocoForTest } from '$core/testing/transloco-testing';
import { ExerciseGroup } from '$generated/user-models';

const GROUP: ExerciseGroup = {
  id: 1,
  createdAt: '',
  updatedAt: '',
  name: 'Push Exercises',
  owner: 'seed',
};

function setup(permissions: string[]) {
  const compendiumApi: Partial<CompendiumApiClient> = {
    fetchExerciseGroup: vi.fn().mockResolvedValue(GROUP),
    fetchExerciseGroupPermissions: vi.fn().mockResolvedValue({ permissions }),
    deleteExerciseGroup: vi.fn(),
  };

  return {
    compendiumApi,
    providers: [
      provideRouter([]),
      provideLocationMocks(),
      provideTanStackQuery(new QueryClient({ defaultOptions: { queries: { retry: false } } })),
      {
        provide: ActivatedRoute,
        useValue: { paramMap: of(convertToParamMap({ id: '1' })) },
      },
      { provide: CompendiumApiClient, useValue: compendiumApi },
      provideTranslocoForTest(),
    ],
  };
}

describe('ExerciseGroupDetail', () => {
  it('shows edit button when user has MODIFY permission', async () => {
    const { providers } = setup(['READ', 'MODIFY', 'DELETE']);
    await render(ExerciseGroupDetail, { providers });

    await waitFor(() => {
      expect(screen.getByText('common.edit')).toBeTruthy();
    });
  });

  it('hides edit button when user lacks MODIFY permission', async () => {
    const { providers } = setup(['READ']);
    await render(ExerciseGroupDetail, { providers });

    await waitFor(() => {
      expect(screen.getByText('Push Exercises')).toBeTruthy();
    });
    expect(screen.queryByText('common.edit')).toBeNull();
  });

  it('shows delete button when user has DELETE permission', async () => {
    const { providers } = setup(['READ', 'MODIFY', 'DELETE']);
    await render(ExerciseGroupDetail, { providers });

    await waitFor(() => {
      expect(screen.getByText('common.delete')).toBeTruthy();
    });
  });

  it('hides delete button when user lacks DELETE permission', async () => {
    const { providers } = setup(['READ']);
    await render(ExerciseGroupDetail, { providers });

    await waitFor(() => {
      expect(screen.getByText('Push Exercises')).toBeTruthy();
    });
    expect(screen.queryByText('common.delete')).toBeNull();
  });
});
