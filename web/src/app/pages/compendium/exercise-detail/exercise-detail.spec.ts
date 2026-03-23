import { render, screen, waitFor } from '@testing-library/angular';
import { provideRouter, ActivatedRoute } from '@angular/router';
import { provideLocationMocks } from '@angular/common/testing';
import { provideTanStackQuery, QueryClient } from '@tanstack/angular-query-experimental';
import { of } from 'rxjs';
import { convertToParamMap } from '@angular/router';
import { ExerciseDetail } from './exercise-detail';
import { CompendiumApiClient } from '$core/api-clients/compendium-api-client';
import { UserApiClient } from '$core/api-clients/user-api-client';
import { provideTranslocoForTest } from '$core/testing/transloco-testing';
import { Exercise } from '$generated/models';

const EXERCISE: Exercise = {
  id: 1,
  createdAt: '',
  updatedAt: '',
  name: 'Bench Press',
  type: 'STRENGTH',
  force: [],
  primaryMuscles: ['CHEST'],
  secondaryMuscles: ['TRICEPS'],
  technicalDifficulty: 'intermediate',
  bodyWeightScaling: 0,
  suggestedMeasurementParadigms: [],
  description: 'A compound exercise',
  instructions: [],
  images: [],
  alternativeNames: [],
  owner: 'seed',
  public: true,
  version: 1,
  templateId: 'tmpl-bench',
  equipmentIds: [],
};

const USER_EXERCISE: Exercise = {
  id: 10,
  createdAt: '',
  updatedAt: '',
  name: 'Bench Press',
  type: 'STRENGTH',
  force: [],
  primaryMuscles: ['CHEST'],
  secondaryMuscles: ['TRICEPS'],
  technicalDifficulty: 'intermediate',
  bodyWeightScaling: 0,
  suggestedMeasurementParadigms: [],
  description: 'A compound exercise',
  instructions: [],
  images: [],
  alternativeNames: [],
  owner: 'anon',
  public: false,
  version: 1,
  templateId: 'tmpl-bench',
  equipmentIds: [],
};

function setup(userExercises: Exercise[] = []) {
  const compendiumApi: Partial<CompendiumApiClient> = {
    fetchExercise: vi.fn().mockResolvedValue(EXERCISE),
    fetchExerciseVersions: vi.fn().mockResolvedValue([EXERCISE]),
    fetchExercisePermissions: vi
      .fn()
      .mockResolvedValue({ permissions: ['READ', 'MODIFY', 'DELETE'] }),
    deleteExercise: vi.fn(),
  };
  const userApi: Partial<UserApiClient> = {
    fetchUserExercises: vi.fn().mockResolvedValue(userExercises),
    createUserExercise: vi.fn().mockResolvedValue(USER_EXERCISE),
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

function setupWithPermissions(permissions: string[], userExercises: Exercise[] = []) {
  const compendiumApi: Partial<CompendiumApiClient> = {
    fetchExercise: vi.fn().mockResolvedValue(EXERCISE),
    fetchExerciseVersions: vi.fn().mockResolvedValue([EXERCISE]),
    fetchExercisePermissions: vi.fn().mockResolvedValue({ permissions }),
    deleteExercise: vi.fn(),
  };
  const userApi: Partial<UserApiClient> = {
    fetchUserExercises: vi.fn().mockResolvedValue(userExercises),
    createUserExercise: vi.fn().mockResolvedValue(USER_EXERCISE),
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

describe('ExerciseDetail', () => {
  it('shows "compendium.exercises.addToMine" when exercise is not yet added', async () => {
    const { providers } = setup([]);
    await render(ExerciseDetail, { providers });

    await waitFor(() => {
      expect(screen.getByText('compendium.exercises.addToMine')).toBeTruthy();
    });
    expect(screen.queryByText('compendium.exercises.alreadyAdded')).toBeNull();
  });

  it('shows "compendium.exercises.alreadyAdded" link when exercise is already imported', async () => {
    const { providers } = setup([USER_EXERCISE]);
    await render(ExerciseDetail, { providers });

    await waitFor(() => {
      expect(screen.getByText('compendium.exercises.alreadyAdded')).toBeTruthy();
    });
    expect(screen.queryByText('compendium.exercises.addToMine')).toBeNull();

    const link = screen.getByText('compendium.exercises.alreadyAdded');
    expect(link.getAttribute('href')).toBe('/user/exercises/10');
  });

  it('shows "compendium.exercises.addToMine" when user has other exercises but not this one', async () => {
    const otherExercise: Exercise = {
      ...USER_EXERCISE,
      id: 99,
      templateId: 'tmpl-other',
    };
    const { providers } = setup([otherExercise]);
    await render(ExerciseDetail, { providers });

    await waitFor(() => {
      expect(screen.getByText('compendium.exercises.addToMine')).toBeTruthy();
    });
    expect(screen.queryByText('compendium.exercises.alreadyAdded')).toBeNull();
  });

  it('shows edit button when user has MODIFY permission', async () => {
    const { providers } = setupWithPermissions(['READ', 'MODIFY', 'DELETE']);
    await render(ExerciseDetail, { providers });

    await waitFor(() => {
      expect(screen.getByText('common.edit')).toBeTruthy();
    });
  });

  it('hides edit button when user lacks MODIFY permission', async () => {
    const { providers } = setupWithPermissions(['READ']);
    await render(ExerciseDetail, { providers });

    await waitFor(() => {
      expect(screen.getByText('Bench Press')).toBeTruthy();
    });
    expect(screen.queryByText('common.edit')).toBeNull();
  });

  it('shows delete button when user has DELETE permission', async () => {
    const { providers } = setupWithPermissions(['READ', 'MODIFY', 'DELETE']);
    await render(ExerciseDetail, { providers });

    await waitFor(() => {
      expect(screen.getByText('common.delete')).toBeTruthy();
    });
  });

  it('hides delete button when user lacks DELETE permission', async () => {
    const { providers } = setupWithPermissions(['READ']);
    await render(ExerciseDetail, { providers });

    await waitFor(() => {
      expect(screen.getByText('Bench Press')).toBeTruthy();
    });
    expect(screen.queryByText('common.delete')).toBeNull();
  });
});
