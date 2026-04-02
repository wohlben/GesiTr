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
import { Exercise, ExerciseRelationship } from '$generated/models';

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
  equipmentIds: [],
};

const FORKED_RELATIONSHIP: ExerciseRelationship = {
  id: 1,
  createdAt: '',
  updatedAt: '',
  fromExerciseId: 10,
  toExerciseId: 1,
  relationshipType: 'forked',
  strength: 0,
  owner: 'anon',
};

function setup(forkedRelationships: ExerciseRelationship[] = []) {
  const compendiumApi: Partial<CompendiumApiClient> = {
    fetchExercise: vi.fn().mockResolvedValue(EXERCISE),
    fetchExerciseVersions: vi.fn().mockResolvedValue([EXERCISE]),
    fetchExercisePermissions: vi
      .fn()
      .mockResolvedValue({ permissions: ['READ', 'MODIFY', 'DELETE'] }),
    fetchExerciseRelationships: vi.fn().mockResolvedValue(forkedRelationships),
    deleteExercise: vi.fn(),
  };
  const userApi: Partial<UserApiClient> = {
    createUserExercise: vi.fn().mockResolvedValue({ id: 10 }),
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

function setupWithPermissions(permissions: string[]) {
  const compendiumApi: Partial<CompendiumApiClient> = {
    fetchExercise: vi.fn().mockResolvedValue(EXERCISE),
    fetchExerciseVersions: vi.fn().mockResolvedValue([EXERCISE]),
    fetchExercisePermissions: vi.fn().mockResolvedValue({ permissions }),
    fetchExerciseRelationships: vi.fn().mockResolvedValue([]),
    deleteExercise: vi.fn(),
  };
  const userApi: Partial<UserApiClient> = {
    createUserExercise: vi.fn().mockResolvedValue({ id: 10 }),
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
  it('shows track button and fork button when exercise is not yet forked', async () => {
    const { providers } = setup([]);
    await render(ExerciseDetail, { providers });

    await waitFor(() => {
      expect(screen.getByText('compendium.exercises.track')).toBeTruthy();
      expect(screen.getByText('compendium.exercises.addToMine')).toBeTruthy();
    });
    expect(screen.queryByText('compendium.exercises.alreadyAdded')).toBeNull();
  });

  it('shows track button and "already forked" link when exercise is forked', async () => {
    const { providers } = setup([FORKED_RELATIONSHIP]);
    await render(ExerciseDetail, { providers });

    await waitFor(() => {
      expect(screen.getByText('compendium.exercises.track')).toBeTruthy();
      expect(screen.getByText('compendium.exercises.alreadyAdded')).toBeTruthy();
    });
    expect(screen.queryByText('compendium.exercises.addToMine')).toBeNull();

    const link = screen.getByText('compendium.exercises.alreadyAdded');
    expect(link.getAttribute('href')).toBe('/compendium/exercises/10');
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
