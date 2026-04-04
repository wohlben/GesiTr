import { expect, test } from '../../base-test';
import {
  createExercise,
  deleteExercise,
  createWorkout,
  deleteWorkout,
  createExerciseScheme,
  deleteExerciseScheme,
  createWorkoutSection,
  deleteWorkoutSection,
  createWorkoutSectionItem,
  deleteWorkoutSectionItem,
  upsertSchemeSectionItem,
  fetchWorkoutLogs,
  deleteWorkoutLog,
} from '../../helpers';

interface TestExercise {
  name: string;
  scheme: { sets: number; reps: number; weight: number; restBetweenSets?: number };
}

interface Variant {
  workoutName: string;
  sectionLabel: string;
  exercises: TestExercise[];
}

async function createFixtures(
  request: Parameters<Parameters<typeof test>[2]>[0]['request'],
  v: Variant,
  sectionOverrides: Record<string, unknown> = {},
) {
  const cleanup: (() => Promise<void>)[] = [];

  const workout = await createWorkout(request, { name: v.workoutName });
  cleanup.push(() => deleteWorkout(request, workout.id));

  const section = await createWorkoutSection(request, {
    workoutId: workout.id,
    label: v.sectionLabel,
    restBetweenExercises: 90,
    ...sectionOverrides,
  });
  cleanup.push(() => deleteWorkoutSection(request, section.id));

  for (let i = 0; i < v.exercises.length; i++) {
    const ex = v.exercises[i];
    const exercise = await createExercise(request, { names: [ex.name] });
    cleanup.push(() => deleteExercise(request, exercise.id));
    const scheme = await createExerciseScheme(request, {
      exerciseId: exercise.id,
      ...ex.scheme,
    });
    cleanup.push(() => deleteExerciseScheme(request, scheme.id));
    const sectionExercise = await createWorkoutSectionItem(request, {
      workoutSectionId: section.id,
      exerciseId: exercise.id,
      position: i,
    });
    await upsertSchemeSectionItem(request, {
      exerciseSchemeId: scheme.id,
      workoutSectionItemId: sectionExercise.id,
    });
    cleanup.push(() => deleteWorkoutSectionItem(request, sectionExercise.id));
  }

  return { workout, cleanup };
}

async function cleanupPlanningLogs(
  request: Parameters<Parameters<typeof test>[2]>[0]['request'],
  workoutId: number,
) {
  const logs = await fetchWorkoutLogs(request, { workoutId, status: 'planning' });
  for (const log of logs) {
    await deleteWorkoutLog(request, log.id);
  }
}

test.describe('/compendium/workouts/[id]/start', () => {
  test('renders start page with workout exercises', async ({ request, page }) => {
    const v: Variant = {
      workoutName: 'Start Push Day',
      sectionLabel: 'Main Lifts',
      exercises: [
        {
          name: 'WS Bench Press',
          scheme: { sets: 5, reps: 5, weight: 100, restBetweenSets: 180 },
        },
        {
          name: 'WS Overhead Press',
          scheme: { sets: 4, reps: 8, weight: 50, restBetweenSets: 120 },
        },
      ],
    };
    const { workout, cleanup } = await createFixtures(request, v);

    await page.goto(`/compendium/workouts/${workout.id}/start`, { waitUntil: 'networkidle' });
    await expect(page.locator('h1')).toHaveText('Plan Workout');
    await expect(page.locator('input#name')).toHaveValue(v.workoutName);
    await expect(page.getByText('Section 1')).toBeVisible();
    for (const ex of v.exercises) {
      await expect(page.getByText(ex.name)).toBeVisible({ timeout: 10000 });
    }

    await cleanupPlanningLogs(request, workout.id);
    for (const fn of cleanup.reverse()) {
      await fn();
    }
  });
});
