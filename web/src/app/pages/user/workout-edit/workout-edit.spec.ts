import { render, screen, waitFor } from '@testing-library/angular';
import { provideRouter, ActivatedRoute } from '@angular/router';
import { provideLocationMocks } from '@angular/common/testing';
import { provideTanStackQuery, QueryClient } from '@tanstack/angular-query-experimental';
import { of } from 'rxjs';
import { convertToParamMap } from '@angular/router';
import { UserApiClient } from '$core/api-clients/user-api-client';
import { CompendiumApiClient } from '$core/api-clients/compendium-api-client';
import { provideTranslocoForTest } from '$core/testing/transloco-testing';
import { Workout } from '$generated/user-models';
import { WorkoutEdit } from './workout-edit';

function makeWorkout(overrides: Partial<Workout> = {}): Workout {
  return {
    id: 1,
    createdAt: '',
    updatedAt: '',
    owner: 'alice',
    name: 'Push Day',
    sections: [],
    ...overrides,
  };
}

function setup(workout: Workout) {
  const userApi: Partial<UserApiClient> = {
    fetchWorkout: vi.fn().mockResolvedValue(workout),
    fetchMasteryList: vi.fn().mockResolvedValue([]),
  };
  const compendiumApi: Partial<CompendiumApiClient> = {
    fetchExercises: vi.fn().mockResolvedValue({ items: [], total: 0 }),
    fetchExerciseGroups: vi.fn().mockResolvedValue({ items: [], total: 0 }),
  };

  return {
    providers: [
      provideRouter([]),
      provideLocationMocks(),
      provideTanStackQuery(new QueryClient({ defaultOptions: { queries: { retry: false } } })),
      {
        provide: ActivatedRoute,
        useValue: { paramMap: of(convertToParamMap({ id: '1' })) },
      },
      { provide: UserApiClient, useValue: userApi },
      { provide: CompendiumApiClient, useValue: compendiumApi },
      provideTranslocoForTest(),
    ],
  };
}

describe('WorkoutEdit', () => {
  it('shows workoutGroup banner for shared workouts', async () => {
    const { providers } = setup(
      makeWorkout({ workoutGroup: { groupName: 'Gym Buddies', membership: 'member' } }),
    );
    await render(WorkoutEdit, { providers });

    await waitFor(() => {
      expect(screen.getByText(/Gym Buddies/)).toBeTruthy();
      expect(screen.getByText(/enums\.workoutGroupRole\.member/)).toBeTruthy();
    });
  });

  it('does not show workoutGroup banner for owned workouts', async () => {
    const { providers } = setup(makeWorkout());
    await render(WorkoutEdit, { providers });

    await waitFor(() => {
      expect(screen.getByDisplayValue('Push Day')).toBeTruthy();
    });
    expect(screen.queryByText('Gym Buddies')).toBeNull();
  });
});
