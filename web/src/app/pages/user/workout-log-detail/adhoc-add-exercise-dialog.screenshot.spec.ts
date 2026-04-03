import { render } from '@testing-library/angular';
import { page } from 'vitest/browser';
import { provideTanStackQuery, QueryClient } from '@tanstack/angular-query-experimental';
import { UserApiClient } from '$core/api-clients/user-api-client';
import { CompendiumApiClient } from '$core/api-clients/compendium-api-client';
import { exerciseKeys, masteryKeys } from '$core/query-keys';
import { provideTranslocoForTest } from '$core/testing/transloco-testing';
import { AdhocAddExerciseDialog } from './adhoc-add-exercise-dialog';
import { Exercise } from '$generated/models';

describe('AdhocAddExerciseDialog screenshots', () => {
  afterEach(() => {
    document.documentElement.classList.remove('dark');
  });

  const mockUserExercises: Exercise[] = [
    {
      id: 1,
      owner: 'user',
      names: [{ id: 1, name: 'Bench Press' }],
      type: 'STRENGTH',
      force: [],
      primaryMuscles: [],
      secondaryMuscles: [],
      technicalDifficulty: 'intermediate',
      bodyWeightScaling: 0,
      suggestedMeasurementParadigms: [],
      description: '',
      instructions: [],
      images: [],
      public: false,
      version: 1,
      equipmentIds: [],
      createdAt: '2024-01-01T00:00:00Z',
      updatedAt: '2024-01-01T00:00:00Z',
    },
    {
      id: 2,
      owner: 'user',
      names: [{ id: 2, name: 'Bicep Curl' }],
      type: 'STRENGTH',
      force: [],
      primaryMuscles: [],
      secondaryMuscles: [],
      technicalDifficulty: 'beginner',
      bodyWeightScaling: 0,
      suggestedMeasurementParadigms: [],
      description: '',
      instructions: [],
      images: [],
      public: false,
      version: 1,
      equipmentIds: [],
      createdAt: '2024-01-01T00:00:00Z',
      updatedAt: '2024-01-01T00:00:00Z',
    },
  ];

  async function renderDialog() {
    const queryClient = new QueryClient({
      defaultOptions: { queries: { retry: false } },
    });

    queryClient.setQueryData(exerciseKeys.list({ limit: 1000 }), {
      items: mockUserExercises,
      total: mockUserExercises.length,
    });
    queryClient.setQueryData(masteryKeys.list(), []);

    const result = await render(AdhocAddExerciseDialog, {
      inputs: {
        open: true,
        sectionId: 1,
        logId: 1,
        exerciseCount: 0,
      },
      providers: [
        provideTranslocoForTest(),
        provideTanStackQuery(queryClient),
        {
          provide: UserApiClient,
          useValue: {
            fetchMasteryList: vi.fn().mockResolvedValue([]),
            createExerciseScheme: vi.fn(),
            createWorkoutLogExercise: vi.fn(),
          },
        },
        {
          provide: CompendiumApiClient,
          useValue: {
            fetchExercises: vi
              .fn()
              .mockResolvedValue({ items: mockUserExercises, total: mockUserExercises.length }),
          },
        },
      ],
    });

    // Wait for queries to settle
    await result.fixture.whenStable();

    // Set exercise selection to show the runner (Phase 2)
    const configComponent = result.fixture.componentInstance.exerciseConfig();
    configComponent.model.set({
      exerciseId: 1,
      measurementType: 'REP_BASED',
      sets: 3,
      reps: 10,
      weight: 60,
      restBetweenSets: 90,
      timePerRep: null,
      duration: null,
      distance: null,
      targetTime: null,
    });
    result.fixture.detectChanges();
    await result.fixture.whenStable();

    return result;
  }

  it('light', async () => {
    await renderDialog();
    const locator = page.elementLocator(document.body);
    await expect(locator).toMatchScreenshot('light');
  });

  it('dark', async () => {
    document.documentElement.classList.add('dark');
    await renderDialog();
    const locator = page.elementLocator(document.body);
    await expect(locator).toMatchScreenshot('dark');
  });
});
