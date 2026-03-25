import { render } from '@testing-library/angular';
import { page } from 'vitest/browser';
import { provideTanStackQuery, QueryClient } from '@tanstack/angular-query-experimental';
import { UserApiClient } from '$core/api-clients/user-api-client';
import { userExerciseKeys } from '$core/query-keys';
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
      name: 'Bench Press',
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
      alternativeNames: [],
      public: false,
      version: 1,
      equipmentIds: [],
      createdAt: '2024-01-01T00:00:00Z',
      updatedAt: '2024-01-01T00:00:00Z',
    },
    {
      id: 2,
      owner: 'user',
      name: 'Bicep Curl',
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
      alternativeNames: [],
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

    queryClient.setQueryData(userExerciseKeys.list(), mockUserExercises);

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
            fetchUserExercises: vi.fn().mockResolvedValue(mockUserExercises),
            createExerciseScheme: vi.fn(),
            createWorkoutLogExercise: vi.fn(),
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
