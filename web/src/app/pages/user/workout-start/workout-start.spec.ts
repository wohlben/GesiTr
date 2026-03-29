import { render, screen, waitFor } from '@testing-library/angular';
import { provideRouter, ActivatedRoute } from '@angular/router';
import { provideLocationMocks } from '@angular/common/testing';
import { provideTanStackQuery, QueryClient } from '@tanstack/angular-query-experimental';
import { of } from 'rxjs';
import { convertToParamMap } from '@angular/router';
import { UserApiClient } from '$core/api-clients/user-api-client';
import { CompendiumApiClient } from '$core/api-clients/compendium-api-client';
import { provideTranslocoForTest } from '$core/testing/transloco-testing';
import { Workout, WorkoutLog } from '$generated/user-models';
import { WorkoutStart } from './workout-start';

// brn-select uses ResizeObserver
beforeAll(() => {
  globalThis.ResizeObserver ??= class {
    // eslint-disable-next-line @typescript-eslint/no-empty-function
    observe() {}
    // eslint-disable-next-line @typescript-eslint/no-empty-function
    unobserve() {}
    // eslint-disable-next-line @typescript-eslint/no-empty-function
    disconnect() {}
  } as unknown as typeof ResizeObserver;
});

const WORKOUT: Workout = {
  id: 1,
  createdAt: '',
  updatedAt: '',
  owner: 'alice',
  name: 'Push Day',
  sections: [],
};

const PLANNING_LOG: WorkoutLog = {
  id: 10,
  createdAt: '',
  updatedAt: '',
  owner: 'alice',
  name: 'Push Day',
  status: 'planning',
  sections: [
    {
      id: 1,
      createdAt: '',
      updatedAt: '',
      workoutLogId: 10,
      type: 'main',
      position: 0,
      status: 'planning',
      exercises: [
        {
          id: 1,
          createdAt: '',
          updatedAt: '',
          workoutLogSectionId: 1,
          workoutLogId: 10,
          sourceExerciseSchemeId: 10,
          position: 0,
          status: 'planning',
          targetMeasurementType: 'REP_BASED',
          sets: [
            {
              id: 1,
              createdAt: '',
              updatedAt: '',
              workoutLogExerciseId: 1,
              workoutLogId: 10,
              setNumber: 1,
              status: 'planning',
              targetReps: 10,
              targetWeight: 60,
            },
          ],
        },
      ],
    },
  ],
};

function setup(permissions: string[]) {
  const userApi: Partial<UserApiClient> = {
    fetchWorkout: vi.fn().mockResolvedValue(WORKOUT),
    fetchWorkoutPermissions: vi.fn().mockResolvedValue({ permissions }),
    fetchWorkoutLogs: vi.fn().mockResolvedValue([PLANNING_LOG]),
    fetchExerciseScheme: vi.fn().mockResolvedValue({
      id: 10,
      exerciseId: 5,
      measurementType: 'REP_BASED',
      sets: 3,
      reps: 10,
      weight: 60,
    }),
    fetchUserExercise: vi.fn().mockResolvedValue({ id: 5, name: 'Bench Press', version: 1 }),
  };

  const compendiumApi: Partial<CompendiumApiClient> = {
    fetchExercises: vi.fn().mockResolvedValue({ items: [], total: 0 }),
  };

  return {
    userApi,
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

describe('WorkoutStart', () => {
  describe('with MODIFY permission (owner)', () => {
    it('shows add section button', async () => {
      const { providers } = setup(['READ', 'MODIFY', 'DELETE']);
      await render(WorkoutStart, { providers });

      await waitFor(() => {
        expect(screen.getByText('user.workouts.addSection')).toBeTruthy();
      });
    });

    it('shows start workout button', async () => {
      const { providers } = setup(['READ', 'MODIFY', 'DELETE']);
      await render(WorkoutStart, { providers });

      await waitFor(() => {
        expect(screen.getByText('user.workoutStart.startWorkout')).toBeTruthy();
      });
    });

    it('renders section components', async () => {
      const { providers } = setup(['READ', 'MODIFY', 'DELETE']);
      const { fixture } = await render(WorkoutStart, { providers });

      await waitFor(() => {
        const sections = fixture.nativeElement.querySelectorAll('app-workout-start-section');
        expect(sections.length).toBe(1);
      });
    });
  });

  describe('with READ-only permission (group member)', () => {
    it('hides add section button', async () => {
      const { providers } = setup(['READ']);
      await render(WorkoutStart, { providers });

      await waitFor(() => {
        // Wait for the form to render (sections visible)
        expect(screen.getByText('user.workoutStart.startWorkout')).toBeTruthy();
      });
      expect(screen.queryByText('user.workouts.addSection')).toBeNull();
    });

    it('still shows start workout button', async () => {
      const { providers } = setup(['READ']);
      await render(WorkoutStart, { providers });

      await waitFor(() => {
        expect(screen.getByText('user.workoutStart.startWorkout')).toBeTruthy();
      });
    });

    it('still renders section components', async () => {
      const { providers } = setup(['READ']);
      const { fixture } = await render(WorkoutStart, { providers });

      await waitFor(() => {
        const sections = fixture.nativeElement.querySelectorAll('app-workout-start-section');
        expect(sections.length).toBe(1);
      });
    });

    it('passes readonly to section components', async () => {
      const { providers } = setup(['READ']);
      const { fixture } = await render(WorkoutStart, { providers });

      await waitFor(() => {
        const section = fixture.nativeElement.querySelector('app-workout-start-section');
        expect(section).toBeTruthy();
        // In readonly mode, the section should not have drag handles
        expect(section.querySelector('[cdkDragHandle]')).toBeNull();
      });
    });

    it('keeps set inputs editable', async () => {
      const { providers } = setup(['READ']);
      const { fixture } = await render(WorkoutStart, { providers });

      await waitFor(() => {
        const repsInputs = fixture.nativeElement.querySelectorAll('input[data-field="targetReps"]');
        expect(repsInputs.length).toBeGreaterThan(0);
      });
    });

    it('keeps name and notes editable', async () => {
      const { providers } = setup(['READ']);
      const { fixture } = await render(WorkoutStart, { providers });

      await waitFor(() => {
        const nameInput = fixture.nativeElement.querySelector('input#name') as HTMLInputElement;
        expect(nameInput).toBeTruthy();
        expect(nameInput.disabled).toBe(false);
      });
    });
  });

  describe('permissions edge cases', () => {
    it('treats READ+MODIFY as editable', async () => {
      const { providers } = setup(['READ', 'MODIFY']);
      await render(WorkoutStart, { providers });

      await waitFor(() => {
        expect(screen.getByText('user.workouts.addSection')).toBeTruthy();
      });
    });
  });
});
