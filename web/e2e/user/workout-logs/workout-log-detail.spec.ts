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

async function createFixturesAndStartLog(
  request: Parameters<Parameters<typeof test>[2]>[0]['request'],
  page: Parameters<Parameters<typeof test>[2]>[0]['page'],
  v: Variant,
) {
  const cleanup: (() => Promise<void>)[] = [];

  const workout = await createWorkout(request, { name: v.workoutName });
  cleanup.push(() => deleteWorkout(request, workout.id));

  const section = await createWorkoutSection(request, {
    workoutId: workout.id,
    label: v.sectionLabel,
    restBetweenExercises: 90,
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

  // Navigate to start page, wait for planning log creation, click Start
  await page.goto(`/compendium/workouts/${workout.id}/start`, { waitUntil: 'networkidle' });
  await expect(page.locator('h1')).toHaveText('Plan Workout');
  for (const ex of v.exercises) {
    await expect(page.getByText(ex.name)).toBeVisible({ timeout: 10000 });
  }
  await page.getByRole('button', { name: 'Start Workout' }).click();
  await page.waitForURL(/\/user\/workout-logs\/\d+/, { timeout: 10000 });

  const logId = Number(page.url().match(/\/workout-logs\/(\d+)/)![1]);
  cleanup.push(() => deleteWorkoutLog(request, logId));

  return { logId, cleanup };
}

test.describe('/user/workout-logs/[id]', () => {
  test('renders workout log detail with exercises', async ({ request, page }) => {
    const v: Variant = {
      workoutName: 'Log Push Day',
      sectionLabel: 'Main Lifts',
      exercises: [
        {
          name: 'WLD Bench Press',
          scheme: { sets: 4, reps: 6, weight: 90, restBetweenSets: 150 },
        },
        {
          name: 'WLD Incline DB Press',
          scheme: { sets: 3, reps: 10, weight: 30, restBetweenSets: 90 },
        },
      ],
    };
    const { cleanup } = await createFixturesAndStartLog(request, page, v);

    await expect(page.locator('h1')).toHaveText(v.workoutName, { timeout: 10000 });
    for (const ex of v.exercises) {
      await expect(page.getByText(ex.name).first()).toBeVisible({ timeout: 10000 });
    }

    for (const fn of cleanup.reverse()) {
      await fn();
    }
  });
});
