import { TestBed } from '@angular/core/testing';
import { WorkoutLogDetailStore } from './workout-log-detail.store';
import { UserApiClient } from '$core/api-clients/user-api-client';
import { WorkoutLogSection } from '$generated/user-models';

describe('WorkoutLogDetailStore', () => {
  let store: InstanceType<typeof WorkoutLogDetailStore>;
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
      providers: [WorkoutLogDetailStore, { provide: UserApiClient, useValue: userApiMock }],
    });

    store = TestBed.inject(WorkoutLogDetailStore);
  });

  it('has correct initial state', () => {
    expect(store.exerciseNames()).toEqual({});
    expect(store.isLoading()).toBe(false);
  });

  it('loads exercise names from log sections', async () => {
    userApiMock.fetchExerciseScheme.mockResolvedValue({
      id: 10,
      exerciseId: 5,
      measurementType: 'REP_BASED',
    });
    userApiMock.fetchUserExercise.mockResolvedValue({
      id: 5,
      name: 'Bench Press',
      version: 1,
    });

    const sections = [
      {
        exercises: [{ sourceExerciseSchemeId: 10 }],
      },
    ] as WorkoutLogSection[];

    await store.loadExerciseNames(sections);

    expect(store.isLoading()).toBe(false);
    expect(store.exerciseNames()).toEqual({
      10: 'Bench Press',
    });
  });

  it('handles multiple exercises across sections', async () => {
    userApiMock.fetchExerciseScheme.mockImplementation((id: number) => {
      if (id === 10) {
        return Promise.resolve({ id: 10, exerciseId: 5, measurementType: 'REP_BASED' });
      }
      return Promise.resolve({ id: 20, exerciseId: 6, measurementType: 'TIME_BASED' });
    });

    userApiMock.fetchUserExercise.mockImplementation((id: number) => {
      if (id === 5) {
        return Promise.resolve({ id: 5, name: 'Squat', version: 1 });
      }
      return Promise.resolve({ id: 6, name: 'Plank', version: 2 });
    });

    const sections = [
      { exercises: [{ sourceExerciseSchemeId: 10 }] },
      { exercises: [{ sourceExerciseSchemeId: 20 }] },
    ] as WorkoutLogSection[];

    await store.loadExerciseNames(sections);

    expect(store.exerciseNames()).toEqual({
      10: 'Squat',
      20: 'Plank',
    });
  });

  it('handles scheme fetch failure gracefully', async () => {
    userApiMock.fetchExerciseScheme.mockRejectedValue(new Error('not found'));

    const sections = [{ exercises: [{ sourceExerciseSchemeId: 99 }] }] as WorkoutLogSection[];

    await store.loadExerciseNames(sections);

    expect(store.isLoading()).toBe(false);
    expect(store.exerciseNames()).toEqual({});
  });

  it('handles exercise name fetch failure with fallback', async () => {
    userApiMock.fetchExerciseScheme.mockResolvedValue({
      id: 10,
      exerciseId: 5,
      measurementType: 'REP_BASED',
    });
    userApiMock.fetchUserExercise.mockRejectedValue(new Error('not found'));

    const sections = [{ exercises: [{ sourceExerciseSchemeId: 10 }] }] as WorkoutLogSection[];

    await store.loadExerciseNames(sections);

    expect(store.exerciseNames()).toEqual({
      10: 'Exercise #5',
    });
  });

  it('handles empty sections', async () => {
    await store.loadExerciseNames([]);

    expect(store.isLoading()).toBe(false);
    expect(store.exerciseNames()).toEqual({});
  });

  it('deduplicates scheme fetches for same scheme ID across sections', async () => {
    userApiMock.fetchExerciseScheme.mockResolvedValue({
      id: 10,
      exerciseId: 5,
      measurementType: 'REP_BASED',
    });
    userApiMock.fetchUserExercise.mockResolvedValue({
      id: 5,
      name: 'Curl',
      version: 1,
    });

    const sections = [
      { exercises: [{ sourceExerciseSchemeId: 10 }] },
      { exercises: [{ sourceExerciseSchemeId: 10 }] },
    ] as WorkoutLogSection[];

    await store.loadExerciseNames(sections);

    // Should only fetch scheme once despite appearing in both sections
    expect(userApiMock.fetchExerciseScheme).toHaveBeenCalledTimes(1);
    expect(store.exerciseNames()).toEqual({ 10: 'Curl' });
  });

  it('deduplicates user exercise fetches', async () => {
    userApiMock.fetchExerciseScheme.mockImplementation((id: number) =>
      Promise.resolve({ id, exerciseId: 5, measurementType: 'REP_BASED' }),
    );
    userApiMock.fetchUserExercise.mockResolvedValue({
      id: 5,
      name: 'Curl',
      version: 1,
    });

    const sections = [
      { exercises: [{ sourceExerciseSchemeId: 10 }, { sourceExerciseSchemeId: 11 }] },
    ] as WorkoutLogSection[];

    await store.loadExerciseNames(sections);

    // Two schemes point to same user exercise — should fetch only once
    expect(userApiMock.fetchUserExercise).toHaveBeenCalledTimes(1);
    expect(store.exerciseNames()[10]).toBe('Curl');
    expect(store.exerciseNames()[11]).toBe('Curl');
  });
});
