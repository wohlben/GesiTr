import { render, screen, waitFor } from '@testing-library/angular';
import { provideRouter, ActivatedRoute } from '@angular/router';
import { provideLocationMocks } from '@angular/common/testing';
import { provideTanStackQuery, QueryClient } from '@tanstack/angular-query-experimental';
import { of } from 'rxjs';
import { convertToParamMap } from '@angular/router';
import { ExerciseDetail } from './exercise-detail';
import { CompendiumApiClient } from '$core/api-clients/compendium-api-client';
import { UserApiClient } from '$core/api-clients/user-api-client';
import { Exercise } from '$generated/models';
import { UserExercise } from '$generated/user-models';

const EXERCISE: Exercise = {
  id: 1,
  createdAt: '',
  updatedAt: '',
  name: 'Bench Press',
  slug: 'bench-press',
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
  createdBy: 'seed',
  version: 1,
  templateId: 'tmpl-bench',
  equipmentIds: [],
};

const USER_EXERCISE: UserExercise = {
  id: 10,
  createdAt: '',
  updatedAt: '',
  owner: 'anon',
  compendiumExerciseId: 'tmpl-bench',
  compendiumVersion: 1,
};

function setup(userExercises: UserExercise[] = []) {
  const compendiumApi: Partial<CompendiumApiClient> = {
    fetchExercise: vi.fn().mockResolvedValue(EXERCISE),
    fetchExerciseVersions: vi.fn().mockResolvedValue([EXERCISE]),
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
    ],
  };
}

describe('ExerciseDetail', () => {
  it('shows "Add to My Exercises" when exercise is not yet added', async () => {
    const { providers } = setup([]);
    await render(ExerciseDetail, { providers });

    await waitFor(() => {
      expect(screen.getByText('Add to My Exercises')).toBeTruthy();
    });
    expect(screen.queryByText('Already Added')).toBeNull();
  });

  it('shows "Already Added" link when exercise is already imported', async () => {
    const { providers } = setup([USER_EXERCISE]);
    await render(ExerciseDetail, { providers });

    await waitFor(() => {
      expect(screen.getByText('Already Added')).toBeTruthy();
    });
    expect(screen.queryByText('Add to My Exercises')).toBeNull();

    const link = screen.getByText('Already Added');
    expect(link.getAttribute('href')).toBe('/user/exercises/10');
  });

  it('shows "Add to My Exercises" when user has other exercises but not this one', async () => {
    const otherExercise: UserExercise = {
      ...USER_EXERCISE,
      id: 99,
      compendiumExerciseId: 'tmpl-other',
    };
    const { providers } = setup([otherExercise]);
    await render(ExerciseDetail, { providers });

    await waitFor(() => {
      expect(screen.getByText('Add to My Exercises')).toBeTruthy();
    });
    expect(screen.queryByText('Already Added')).toBeNull();
  });
});
