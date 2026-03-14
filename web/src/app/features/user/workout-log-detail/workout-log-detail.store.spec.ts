import { TestBed } from '@angular/core/testing';
import { WorkoutLogDetailStore } from './workout-log-detail.store';
import { UserApiClient } from '$core/api-clients/user-api-client';
import { CompendiumApiClient } from '$core/api-clients/compendium-api-client';
import { WorkoutLogSection } from '$generated/user-models';

describe('WorkoutLogDetailStore', () => {
  let store: InstanceType<typeof WorkoutLogDetailStore>;
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
        WorkoutLogDetailStore,
        { provide: UserApiClient, useValue: userApiMock },
        { provide: CompendiumApiClient, useValue: compendiumApiMock },
      ],
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
      userExerciseId: 5,
      measurementType: 'REP_BASED',
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
        return Promise.resolve({ id: 10, userExerciseId: 5, measurementType: 'REP_BASED' });
      }
      return Promise.resolve({ id: 20, userExerciseId: 6, measurementType: 'TIME_BASED' });
    });

    userApiMock.fetchUserExercise.mockImplementation((id: number) => {
      if (id === 5) {
        return Promise.resolve({ id: 5, compendiumExerciseId: 'tmpl-a', compendiumVersion: 1 });
      }
      return Promise.resolve({ id: 6, compendiumExerciseId: 'tmpl-b', compendiumVersion: 2 });
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
      userExerciseId: 5,
      measurementType: 'REP_BASED',
    });
    userApiMock.fetchUserExercise.mockResolvedValue({
      id: 5,
      compendiumExerciseId: 'tmpl-gone',
      compendiumVersion: 1,
    });
    compendiumApiMock.fetchExerciseVersion.mockRejectedValue(new Error('not found'));

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
      userExerciseId: 5,
      measurementType: 'REP_BASED',
    });
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
      Promise.resolve({ id, userExerciseId: 5, measurementType: 'REP_BASED' }),
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
      { exercises: [{ sourceExerciseSchemeId: 10 }, { sourceExerciseSchemeId: 11 }] },
    ] as WorkoutLogSection[];

    await store.loadExerciseNames(sections);

    // Two schemes point to same user exercise — should fetch only once
    expect(userApiMock.fetchUserExercise).toHaveBeenCalledTimes(1);
    expect(store.exerciseNames()[10]).toBe('Curl');
    expect(store.exerciseNames()[11]).toBe('Curl');
  });
});
