import { render } from '@testing-library/angular';
import { page } from 'vitest/browser';
import { provideTanStackQuery, QueryClient } from '@tanstack/angular-query-experimental';
import { UserApiClient } from '$core/api-clients/user-api-client';
import { CompendiumApiClient } from '$core/api-clients/compendium-api-client';
import { userExerciseKeys, exerciseKeys } from '$core/query-keys';
import { AdhocAddExerciseDialog } from './adhoc-add-exercise-dialog';
import { UserExercise } from '$generated/user-models';
import { Exercise } from '$generated/models';

describe('AdhocAddExerciseDialog screenshots', () => {
  afterEach(() => {
    document.documentElement.classList.remove('dark');
  });

  const mockUserExercises: UserExercise[] = [
    {
      id: 1,
      owner: 'user',
      compendiumExerciseId: 'bench-press',
      compendiumVersion: 1,
      createdAt: '2024-01-01T00:00:00Z',
      updatedAt: '2024-01-01T00:00:00Z',
    },
    {
      id: 2,
      owner: 'user',
      compendiumExerciseId: 'bicep-curl',
      compendiumVersion: 1,
      createdAt: '2024-01-01T00:00:00Z',
      updatedAt: '2024-01-01T00:00:00Z',
    },
  ];

  const mockExerciseVersions: Record<string, { version: number; snapshot: Partial<Exercise> }> = {
    'bench-press': { version: 1, snapshot: { name: 'Bench Press' } },
    'bicep-curl': { version: 1, snapshot: { name: 'Bicep Curl' } },
  };

  async function renderDialog() {
    const queryClient = new QueryClient({
      defaultOptions: { queries: { retry: false } },
    });

    queryClient.setQueryData(userExerciseKeys.list(), mockUserExercises);
    for (const ue of mockUserExercises) {
      queryClient.setQueryData(
        exerciseKeys.version(ue.compendiumExerciseId, ue.compendiumVersion),
        mockExerciseVersions[ue.compendiumExerciseId],
      );
    }

    const result = await render(AdhocAddExerciseDialog, {
      inputs: {
        open: true,
        sectionId: 1,
        logId: 1,
        exerciseCount: 0,
      },
      providers: [
        provideTanStackQuery(queryClient),
        {
          provide: UserApiClient,
          useValue: {
            fetchUserExercises: vi.fn().mockResolvedValue(mockUserExercises),
            createExerciseScheme: vi.fn(),
            createWorkoutLogExercise: vi.fn(),
          },
        },
        {
          provide: CompendiumApiClient,
          useValue: {
            fetchExerciseVersion: vi
              .fn()
              .mockImplementation((templateId: string) =>
                Promise.resolve(mockExerciseVersions[templateId]),
              ),
          },
        },
      ],
    });

    // Wait for queries to settle
    await result.fixture.whenStable();

    // Set exercise selection to show the runner (Phase 2)
    const configComponent = result.fixture.componentInstance.exerciseConfig;
    configComponent.userExerciseId.set(1);
    configComponent.sets.set(3);
    configComponent.reps.set(10);
    configComponent.weight.set(60);
    configComponent.restBetweenSets.set(90);
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
