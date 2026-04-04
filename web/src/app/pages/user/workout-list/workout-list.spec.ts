import { render, screen, waitFor } from '@testing-library/angular';
import { provideRouter } from '@angular/router';
import { provideLocationMocks } from '@angular/common/testing';
import { provideTanStackQuery, QueryClient } from '@tanstack/angular-query-experimental';
import { UserApiClient } from '$core/api-clients/user-api-client';
import { provideTranslocoForTest } from '$core/testing/transloco-testing';
import { Workout } from '$generated/user-models';
import { WorkoutList } from './workout-list';

function makeWorkout(overrides: Partial<Workout> = {}): Workout {
  return {
    id: 1,
    createdAt: '',
    updatedAt: '',
    ownershipGroupId: 0,
    name: 'Push Day',
    public: false,
    version: 0,
    sections: [],
    ...overrides,
  };
}

function setup(workouts: Workout[]) {
  const userApi: Partial<UserApiClient> = {
    fetchWorkouts: vi.fn().mockResolvedValue(workouts),
    startAdhocWorkoutLog: vi.fn(),
  };

  return {
    providers: [
      provideRouter([]),
      provideLocationMocks(),
      provideTanStackQuery(new QueryClient({ defaultOptions: { queries: { retry: false } } })),
      { provide: UserApiClient, useValue: userApi },
      provideTranslocoForTest(),
    ],
  };
}

describe('WorkoutList', () => {
  it('shows membership badge for shared workouts', async () => {
    const { providers } = setup([
      makeWorkout({ workoutGroup: { groupName: 'Gym Buddies', membership: 'member' } }),
    ]);
    await render(WorkoutList, { providers });

    await waitFor(() => {
      expect(screen.getByText('enums.workoutGroupRole.member')).toBeTruthy();
    });
  });

  it('does not show membership badge for owned workouts', async () => {
    const { providers } = setup([makeWorkout()]);
    await render(WorkoutList, { providers });

    await waitFor(() => {
      expect(screen.getByText('Push Day')).toBeTruthy();
    });
    expect(screen.queryByText('enums.workoutGroupRole.member')).toBeNull();
    expect(screen.queryByText('enums.workoutGroupRole.admin')).toBeNull();
    expect(screen.queryByText('enums.workoutGroupRole.invited')).toBeNull();
  });

  it('shows correct badge for each role type', async () => {
    const { providers } = setup([
      makeWorkout({
        id: 1,
        name: 'Workout A',
        workoutGroup: { groupName: 'Group A', membership: 'invited' },
      }),
      makeWorkout({
        id: 2,
        name: 'Workout B',
        workoutGroup: { groupName: 'Group B', membership: 'admin' },
      }),
      makeWorkout({ id: 3, name: 'Workout C' }),
    ]);
    await render(WorkoutList, { providers });

    await waitFor(() => {
      expect(screen.getByText('enums.workoutGroupRole.invited')).toBeTruthy();
      expect(screen.getByText('enums.workoutGroupRole.admin')).toBeTruthy();
    });
  });
});
