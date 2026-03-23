import { inject } from '@angular/core';
import { signalStore, withState, withMethods, patchState } from '@ngrx/signals';
import { UserApiClient } from '$core/api-clients/user-api-client';
import { WorkoutSection, WorkoutLogSection } from '$generated/user-models';
import { ExerciseScheme } from '$generated/models';

export interface SetPreview {
  setNumber: number;
  targetReps?: number | null;
  targetWeight?: number | null;
  targetDuration?: number | null;
  targetDistance?: number | null;
  targetTime?: number | null;
  restAfterSeconds?: number | null;
}

export interface ExerciseDisplayInfo {
  name: string;
  summary: string;
  measurementType: string;
  sets: SetPreview[];
}

interface WorkoutStartState {
  exerciseDisplay: Record<number, ExerciseDisplayInfo>;
  isLoadingDisplay: boolean;
}

const initialState: WorkoutStartState = {
  exerciseDisplay: {},
  isLoadingDisplay: false,
};

export function formatSchemeSummary(scheme: {
  measurementType: string;
  sets?: number | null;
  reps?: number | null;
  weight?: number | null;
  duration?: number | null;
  distance?: number | null;
  targetTime?: number | null;
}): string {
  if (scheme.measurementType === 'REP_BASED') {
    const parts: string[] = [];
    if (scheme.sets) parts.push(`${scheme.sets}x`);
    if (scheme.reps) parts.push(`${scheme.reps}`);
    const setsReps = parts.join('');
    if (scheme.weight) return `${setsReps} @ ${scheme.weight}kg`;
    return setsReps || 'Rep based';
  }
  if (scheme.measurementType === 'TIME_BASED') {
    if (scheme.duration) return `${scheme.duration}s`;
    return 'Time based';
  }
  if (scheme.measurementType === 'DISTANCE_BASED') {
    if (scheme.distance) return `${scheme.distance}m`;
    return 'Distance based';
  }
  return scheme.measurementType;
}

export const WorkoutStartStore = signalStore(
  withState(initialState),
  withMethods((store) => {
    const userApi = inject(UserApiClient);

    return {
      async loadExerciseDisplay(sections: WorkoutSection[]) {
        patchState(store, { isLoadingDisplay: true });

        const display: Record<number, ExerciseDisplayInfo> = {};

        // 1. Fetch all schemes in parallel
        const schemeResults = await Promise.all(
          sections.flatMap((s) =>
            (s.exercises ?? []).map((ex) =>
              userApi
                .fetchExerciseScheme(ex.exerciseSchemeId)
                .then((scheme) => ({ schemeId: ex.exerciseSchemeId, scheme }))
                .catch(() => null),
            ),
          ),
        );
        const schemes = schemeResults.filter(
          (r): r is { schemeId: number; scheme: ExerciseScheme } => r !== null,
        );

        // 2. Fetch unique exercises to get names directly
        const uniqueExerciseIds = [...new Set(schemes.map((s) => s.scheme.exerciseId))];
        const exerciseResults = await Promise.all(
          uniqueExerciseIds.map((id) => userApi.fetchUserExercise(id).catch(() => null)),
        );

        // 3. Build exercise name map
        const exerciseNames: Record<number, string> = {};
        for (const exercise of exerciseResults) {
          if (exercise) {
            exerciseNames[exercise.id] = exercise.name;
          }
        }

        // 4. Build display map with set previews
        for (const item of schemes) {
          const numSets = item.scheme.sets ?? 0;
          const sets: SetPreview[] = [];
          for (let i = 1; i <= numSets; i++) {
            sets.push({
              setNumber: i,
              targetReps: item.scheme.reps,
              targetWeight: item.scheme.weight,
              targetDuration: item.scheme.duration,
              targetDistance: item.scheme.distance,
              targetTime: item.scheme.targetTime,
              restAfterSeconds: i < numSets ? (item.scheme.restBetweenSets ?? null) : null,
            });
          }
          display[item.schemeId] = {
            name: exerciseNames[item.scheme.exerciseId] ?? `Exercise #${item.scheme.exerciseId}`,
            summary: formatSchemeSummary(item.scheme),
            measurementType: item.scheme.measurementType,
            sets,
          };
        }

        patchState(store, { exerciseDisplay: display, isLoadingDisplay: false });
      },

      async loadExerciseDisplayFromLog(sections: WorkoutLogSection[]) {
        patchState(store, { isLoadingDisplay: true });

        const display: Record<number, ExerciseDisplayInfo> = {};

        // Collect unique scheme IDs
        const schemeIds = [
          ...new Set(
            sections.flatMap((s) => (s.exercises ?? []).map((ex) => ex.sourceExerciseSchemeId)),
          ),
        ];

        // Fetch schemes to resolve names via exercise
        const schemeResults = await Promise.all(
          schemeIds.map((id) =>
            userApi
              .fetchExerciseScheme(id)
              .then((scheme) => ({ schemeId: id, scheme }))
              .catch(() => null),
          ),
        );
        const schemes = schemeResults.filter(
          (r): r is { schemeId: number; scheme: ExerciseScheme } => r !== null,
        );

        // Fetch unique exercises to get names directly
        const uniqueExerciseIds = [...new Set(schemes.map((s) => s.scheme.exerciseId))];
        const exerciseResults = await Promise.all(
          uniqueExerciseIds.map((id) => userApi.fetchUserExercise(id).catch(() => null)),
        );

        // Build exercise name map
        const exerciseNames: Record<number, string> = {};
        for (const exercise of exerciseResults) {
          if (exercise) {
            exerciseNames[exercise.id] = exercise.name;
          }
        }

        // Build display map using log exercise data (already snapshotted)
        for (const section of sections) {
          for (const ex of section.exercises ?? []) {
            const schemeItem = schemes.find((s) => s.schemeId === ex.sourceExerciseSchemeId);
            const sets: SetPreview[] = (ex.sets ?? []).map((s) => ({
              setNumber: s.setNumber,
              targetReps: s.targetReps,
              targetWeight: s.targetWeight,
              targetDuration: s.targetDuration,
              targetDistance: s.targetDistance,
              targetTime: s.targetTime,
              restAfterSeconds: s.breakAfterSeconds ?? null,
            }));
            display[ex.id] = {
              name: schemeItem
                ? (exerciseNames[schemeItem.scheme.exerciseId] ??
                  `Exercise #${schemeItem.scheme.exerciseId}`)
                : `Exercise #${ex.sourceExerciseSchemeId}`,
              summary: schemeItem
                ? formatSchemeSummary(schemeItem.scheme)
                : ex.targetMeasurementType,
              measurementType: ex.targetMeasurementType,
              sets,
            };
          }
        }

        patchState(store, { exerciseDisplay: display, isLoadingDisplay: false });
      },

      addExerciseDisplay(
        exerciseLogId: number,
        exerciseName: string,
        scheme: ExerciseScheme,
        sets: SetPreview[],
      ) {
        const current = store.exerciseDisplay();
        patchState(store, {
          exerciseDisplay: {
            ...current,
            [exerciseLogId]: {
              name: exerciseName,
              summary: formatSchemeSummary(scheme),
              measurementType: scheme.measurementType,
              sets,
            },
          },
        });
      },
    };
  }),
);
