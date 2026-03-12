import { TestBed } from '@angular/core/testing';
import { WorkoutStartStore, formatSchemeSummary } from './workout-start.store';
import { UserApiClient } from '$core/api-clients/user-api-client';
import { CompendiumApiClient } from '$core/api-clients/compendium-api-client';
import { WorkoutSection } from '$generated/user-models';

describe('formatSchemeSummary', () => {
  it('formats rep-based with sets, reps, and weight', () => {
    expect(
      formatSchemeSummary({ measurementType: 'REP_BASED', sets: 3, reps: 10, weight: 50 }),
    ).toBe('3x10 @ 50kg');
  });

  it('formats rep-based with sets and reps only', () => {
    expect(formatSchemeSummary({ measurementType: 'REP_BASED', sets: 4, reps: 8 })).toBe('4x8');
  });

  it('formats rep-based with reps only', () => {
    expect(formatSchemeSummary({ measurementType: 'REP_BASED', reps: 12 })).toBe('12');
  });

  it('formats rep-based with no details', () => {
    expect(formatSchemeSummary({ measurementType: 'REP_BASED' })).toBe('Rep based');
  });

  it('formats time-based with duration', () => {
    expect(formatSchemeSummary({ measurementType: 'TIME_BASED', duration: 60 })).toBe('60s');
  });

  it('formats time-based with no duration', () => {
    expect(formatSchemeSummary({ measurementType: 'TIME_BASED' })).toBe('Time based');
  });

  it('formats distance-based with distance', () => {
    expect(formatSchemeSummary({ measurementType: 'DISTANCE_BASED', distance: 5000 })).toBe(
      '5000m',
    );
  });

  it('formats distance-based with no distance', () => {
    expect(formatSchemeSummary({ measurementType: 'DISTANCE_BASED' })).toBe('Distance based');
  });

  it('falls back to measurement type string', () => {
    expect(formatSchemeSummary({ measurementType: 'CUSTOM' })).toBe('CUSTOM');
  });
});

describe('WorkoutStartStore', () => {
  let store: InstanceType<typeof WorkoutStartStore>;
  let userApiMock: {
    fetchExerciseScheme: ReturnType<typeof vi.fn>;
    fetchUserExercise: ReturnType<typeof vi.fn>;
  };
  let compendiumApiMock: {
    fetchExerciseVersion: ReturnType<typeof vi.fn>;
  };

  beforeEach(() => {
    userApiMock = {
      fetchExerciseScheme: vi.fn(),
      fetchUserExercise: vi.fn(),
    };
    compendiumApiMock = {
      fetchExerciseVersion: vi.fn(),
    };

    TestBed.configureTestingModule({
      providers: [
        WorkoutStartStore,
        { provide: UserApiClient, useValue: userApiMock },
        { provide: CompendiumApiClient, useValue: compendiumApiMock },
      ],
    });

    store = TestBed.inject(WorkoutStartStore);
  });

  it('has correct initial state', () => {
    expect(store.exerciseDisplay()).toEqual({});
    expect(store.isLoadingDisplay()).toBe(false);
  });

  it('loads exercise display data from sections', async () => {
    userApiMock.fetchExerciseScheme.mockResolvedValue({
      id: 10,
      userExerciseId: 5,
      measurementType: 'REP_BASED',
      sets: 3,
      reps: 10,
      weight: 60,
    });
    userApiMock.fetchUserExercise.mockResolvedValue({
      id: 5,
      compendiumExerciseId: 'tmpl-abc',
      compendiumVersion: 1,
    });
    compendiumApiMock.fetchExerciseVersion.mockResolvedValue({
      version: 1,
      snapshot: { name: 'Bench Press' },
    });

    const sections = [
      {
        exercises: [{ userExerciseSchemeId: 10 }],
      },
    ] as WorkoutSection[];

    await store.loadExerciseDisplay(sections);

    expect(store.isLoadingDisplay()).toBe(false);
    expect(store.exerciseDisplay()).toEqual({
      10: { name: 'Bench Press', summary: '3x10 @ 60kg' },
    });
  });

  it('handles multiple exercises across sections', async () => {
    userApiMock.fetchExerciseScheme.mockImplementation((id: number) => {
      if (id === 10) {
        return Promise.resolve({
          id: 10,
          userExerciseId: 5,
          measurementType: 'REP_BASED',
          sets: 5,
          reps: 5,
          weight: 100,
        });
      }
      return Promise.resolve({
        id: 20,
        userExerciseId: 6,
        measurementType: 'TIME_BASED',
        duration: 30,
      });
    });

    userApiMock.fetchUserExercise.mockImplementation((id: number) => {
      if (id === 5) {
        return Promise.resolve({
          id: 5,
          compendiumExerciseId: 'tmpl-a',
          compendiumVersion: 1,
        });
      }
      return Promise.resolve({
        id: 6,
        compendiumExerciseId: 'tmpl-b',
        compendiumVersion: 2,
      });
    });

    compendiumApiMock.fetchExerciseVersion.mockImplementation(
      (templateId: string, version: number) => {
        if (templateId === 'tmpl-a' && version === 1) {
          return Promise.resolve({ version: 1, snapshot: { name: 'Squat' } });
        }
        return Promise.resolve({ version: 2, snapshot: { name: 'Plank' } });
      },
    );

    const sections = [
      { exercises: [{ userExerciseSchemeId: 10 }] },
      { exercises: [{ userExerciseSchemeId: 20 }] },
    ] as WorkoutSection[];

    await store.loadExerciseDisplay(sections);

    expect(store.exerciseDisplay()).toEqual({
      10: { name: 'Squat', summary: '5x5 @ 100kg' },
      20: { name: 'Plank', summary: '30s' },
    });
  });

  it('handles scheme fetch failure gracefully', async () => {
    userApiMock.fetchExerciseScheme.mockRejectedValue(new Error('not found'));

    const sections = [{ exercises: [{ userExerciseSchemeId: 99 }] }] as WorkoutSection[];

    await store.loadExerciseDisplay(sections);

    expect(store.isLoadingDisplay()).toBe(false);
    expect(store.exerciseDisplay()).toEqual({});
  });

  it('handles exercise name fetch failure with fallback', async () => {
    userApiMock.fetchExerciseScheme.mockResolvedValue({
      id: 10,
      userExerciseId: 5,
      measurementType: 'REP_BASED',
      reps: 8,
    });
    userApiMock.fetchUserExercise.mockResolvedValue({
      id: 5,
      compendiumExerciseId: 'tmpl-gone',
      compendiumVersion: 1,
    });
    compendiumApiMock.fetchExerciseVersion.mockRejectedValue(new Error('not found'));

    const sections = [{ exercises: [{ userExerciseSchemeId: 10 }] }] as WorkoutSection[];

    await store.loadExerciseDisplay(sections);

    expect(store.exerciseDisplay()).toEqual({
      10: { name: 'Exercise #5', summary: '8' },
    });
  });

  it('handles empty sections', async () => {
    await store.loadExerciseDisplay([]);

    expect(store.isLoadingDisplay()).toBe(false);
    expect(store.exerciseDisplay()).toEqual({});
  });

  it('deduplicates user exercise fetches', async () => {
    // Two schemes pointing to the same userExerciseId
    userApiMock.fetchExerciseScheme.mockImplementation((id: number) =>
      Promise.resolve({
        id,
        userExerciseId: 5,
        measurementType: 'REP_BASED',
        sets: id === 10 ? 3 : 4,
        reps: 10,
      }),
    );
    userApiMock.fetchUserExercise.mockResolvedValue({
      id: 5,
      compendiumExerciseId: 'tmpl-x',
      compendiumVersion: 1,
    });
    compendiumApiMock.fetchExerciseVersion.mockResolvedValue({
      version: 1,
      snapshot: { name: 'Curl' },
    });

    const sections = [
      { exercises: [{ userExerciseSchemeId: 10 }, { userExerciseSchemeId: 11 }] },
    ] as WorkoutSection[];

    await store.loadExerciseDisplay(sections);

    // Should only fetch user exercise once despite two schemes referencing it
    expect(userApiMock.fetchUserExercise).toHaveBeenCalledTimes(1);
    expect(userApiMock.fetchUserExercise).toHaveBeenCalledWith(5);

    expect(store.exerciseDisplay()[10].name).toBe('Curl');
    expect(store.exerciseDisplay()[11].name).toBe('Curl');
  });
});
