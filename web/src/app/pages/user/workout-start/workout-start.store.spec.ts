import { TestBed } from '@angular/core/testing';
import { WorkoutStartStore, formatSchemeSummary } from './workout-start.store';
import { UserApiClient } from '$core/api-clients/user-api-client';
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

  beforeEach(() => {
    userApiMock = {
      fetchExerciseScheme: vi.fn(),
      fetchUserExercise: vi.fn(),
    };

    TestBed.configureTestingModule({
      providers: [WorkoutStartStore, { provide: UserApiClient, useValue: userApiMock }],
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
      exerciseId: 5,
      measurementType: 'REP_BASED',
      sets: 3,
      reps: 10,
      weight: 60,
    });
    userApiMock.fetchUserExercise.mockResolvedValue({
      id: 5,
      name: 'Bench Press',
      version: 1,
    });

    const sections = [
      {
        items: [{ type: 'exercise', exerciseSchemeId: 10 }],
      },
    ] as WorkoutSection[];

    await store.loadExerciseDisplay(sections);

    expect(store.isLoadingDisplay()).toBe(false);
    expect(store.exerciseDisplay()).toEqual({
      10: {
        name: 'Bench Press',
        summary: '3x10 @ 60kg',
        measurementType: 'REP_BASED',
        sets: [
          {
            setNumber: 1,
            targetReps: 10,
            targetWeight: 60,
            targetDuration: undefined,
            targetDistance: undefined,
            targetTime: undefined,
            restAfterSeconds: null,
          },
          {
            setNumber: 2,
            targetReps: 10,
            targetWeight: 60,
            targetDuration: undefined,
            targetDistance: undefined,
            targetTime: undefined,
            restAfterSeconds: null,
          },
          {
            setNumber: 3,
            targetReps: 10,
            targetWeight: 60,
            targetDuration: undefined,
            targetDistance: undefined,
            targetTime: undefined,
            restAfterSeconds: null,
          },
        ],
      },
    });
  });

  it('handles multiple exercises across sections', async () => {
    userApiMock.fetchExerciseScheme.mockImplementation((id: number) => {
      if (id === 10) {
        return Promise.resolve({
          id: 10,
          exerciseId: 5,
          measurementType: 'REP_BASED',
          sets: 5,
          reps: 5,
          weight: 100,
        });
      }
      return Promise.resolve({
        id: 20,
        exerciseId: 6,
        measurementType: 'TIME_BASED',
        duration: 30,
      });
    });

    userApiMock.fetchUserExercise.mockImplementation((id: number) => {
      if (id === 5) {
        return Promise.resolve({
          id: 5,
          name: 'Squat',
          version: 1,
        });
      }
      return Promise.resolve({
        id: 6,
        name: 'Plank',
        version: 2,
      });
    });

    const sections = [
      { items: [{ type: 'exercise', exerciseSchemeId: 10 }] },
      { items: [{ type: 'exercise', exerciseSchemeId: 20 }] },
    ] as WorkoutSection[];

    await store.loadExerciseDisplay(sections);

    expect(store.exerciseDisplay()).toEqual({
      10: {
        name: 'Squat',
        summary: '5x5 @ 100kg',
        measurementType: 'REP_BASED',
        sets: [
          {
            setNumber: 1,
            targetReps: 5,
            targetWeight: 100,
            targetDuration: undefined,
            targetDistance: undefined,
            targetTime: undefined,
            restAfterSeconds: null,
          },
          {
            setNumber: 2,
            targetReps: 5,
            targetWeight: 100,
            targetDuration: undefined,
            targetDistance: undefined,
            targetTime: undefined,
            restAfterSeconds: null,
          },
          {
            setNumber: 3,
            targetReps: 5,
            targetWeight: 100,
            targetDuration: undefined,
            targetDistance: undefined,
            targetTime: undefined,
            restAfterSeconds: null,
          },
          {
            setNumber: 4,
            targetReps: 5,
            targetWeight: 100,
            targetDuration: undefined,
            targetDistance: undefined,
            targetTime: undefined,
            restAfterSeconds: null,
          },
          {
            setNumber: 5,
            targetReps: 5,
            targetWeight: 100,
            targetDuration: undefined,
            targetDistance: undefined,
            targetTime: undefined,
            restAfterSeconds: null,
          },
        ],
      },
      20: {
        name: 'Plank',
        summary: '30s',
        measurementType: 'TIME_BASED',
        sets: [],
      },
    });
  });

  it('handles scheme fetch failure gracefully', async () => {
    userApiMock.fetchExerciseScheme.mockRejectedValue(new Error('not found'));

    const sections = [{ items: [{ type: 'exercise', exerciseSchemeId: 99 }] }] as WorkoutSection[];

    await store.loadExerciseDisplay(sections);

    expect(store.isLoadingDisplay()).toBe(false);
    expect(store.exerciseDisplay()).toEqual({});
  });

  it('handles exercise name fetch failure with fallback', async () => {
    userApiMock.fetchExerciseScheme.mockResolvedValue({
      id: 10,
      exerciseId: 5,
      measurementType: 'REP_BASED',
      reps: 8,
    });
    userApiMock.fetchUserExercise.mockRejectedValue(new Error('not found'));

    const sections = [{ items: [{ type: 'exercise', exerciseSchemeId: 10 }] }] as WorkoutSection[];

    await store.loadExerciseDisplay(sections);

    expect(store.exerciseDisplay()).toEqual({
      10: {
        name: 'Exercise #5',
        summary: '8',
        measurementType: 'REP_BASED',
        sets: [],
      },
    });
  });

  it('handles empty sections', async () => {
    await store.loadExerciseDisplay([]);

    expect(store.isLoadingDisplay()).toBe(false);
    expect(store.exerciseDisplay()).toEqual({});
  });

  it('deduplicates user exercise fetches', async () => {
    // Two schemes pointing to the same exerciseId
    userApiMock.fetchExerciseScheme.mockImplementation((id: number) =>
      Promise.resolve({
        id,
        exerciseId: 5,
        measurementType: 'REP_BASED',
        sets: id === 10 ? 3 : 4,
        reps: 10,
      }),
    );
    userApiMock.fetchUserExercise.mockResolvedValue({
      id: 5,
      name: 'Curl',
      version: 1,
    });

    const sections = [
      {
        items: [
          { type: 'exercise', exerciseSchemeId: 10 },
          { type: 'exercise', exerciseSchemeId: 11 },
        ],
      },
    ] as WorkoutSection[];

    await store.loadExerciseDisplay(sections);

    // Should only fetch user exercise once despite two schemes referencing it
    expect(userApiMock.fetchUserExercise).toHaveBeenCalledTimes(1);
    expect(userApiMock.fetchUserExercise).toHaveBeenCalledWith(5);

    expect(store.exerciseDisplay()[10].name).toBe('Curl');
    expect(store.exerciseDisplay()[11].name).toBe('Curl');
  });

  it('generates set previews with restBetweenSets', async () => {
    userApiMock.fetchExerciseScheme.mockResolvedValue({
      id: 10,
      exerciseId: 5,
      measurementType: 'REP_BASED',
      sets: 3,
      reps: 8,
      weight: 50,
      restBetweenSets: 90,
    });
    userApiMock.fetchUserExercise.mockResolvedValue({
      id: 5,
      name: 'Row',
      version: 1,
    });

    const sections = [{ items: [{ type: 'exercise', exerciseSchemeId: 10 }] }] as WorkoutSection[];

    await store.loadExerciseDisplay(sections);

    const info = store.exerciseDisplay()[10];
    expect(info.sets).toEqual([
      {
        setNumber: 1,
        targetReps: 8,
        targetWeight: 50,
        targetDuration: undefined,
        targetDistance: undefined,
        targetTime: undefined,
        restAfterSeconds: 90,
      },
      {
        setNumber: 2,
        targetReps: 8,
        targetWeight: 50,
        targetDuration: undefined,
        targetDistance: undefined,
        targetTime: undefined,
        restAfterSeconds: 90,
      },
      {
        setNumber: 3,
        targetReps: 8,
        targetWeight: 50,
        targetDuration: undefined,
        targetDistance: undefined,
        targetTime: undefined,
        restAfterSeconds: null,
      },
    ]);
  });
});
