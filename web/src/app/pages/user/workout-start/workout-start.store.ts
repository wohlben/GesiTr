import { inject } from '@angular/core';
import { signalStore, withState, withMethods, patchState } from '@ngrx/signals';
import { UserApiClient } from '$core/api-clients/user-api-client';
import { WorkoutSection, WorkoutLogSection } from '$generated/user-models';
import { ExerciseScheme } from '$generated/user-exercisescheme';
import { formatSchemeSummary } from '$core/scheme-utils';
export { formatSchemeSummary };

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

export const WorkoutStartStore = signalStore(
  withState(initialState),
  withMethods((store) => {
    const userApi = inject(UserApiClient);

    return {
      async loadExerciseDisplay(sections: WorkoutSection[]) {
        patchState(store, { isLoadingDisplay: true });

        const display: Record<number, ExerciseDisplayInfo> = {};

        // 1. Collect exercise items and fetch scheme assignments via join table
        const exerciseItems = sections.flatMap((s) =>
          (s.items ?? []).filter((item) => item.type === 'exercise' && item.exerciseId != null),
        );
        const itemIds = exerciseItems.map((i) => i.id);
        const assignments =
          itemIds.length > 0 ? await userApi.fetchSchemeSectionItems(itemIds) : [];
        const assignmentByItemId = new Map(assignments.map((a) => [a.workoutSectionItemId, a]));

        // 2. Fetch schemes for assigned items
        const schemeResults = await Promise.all(
          assignments.map((a) =>
            userApi
              .fetchExerciseScheme(a.exerciseSchemeId)
              .then((scheme) => ({ itemId: a.workoutSectionItemId, scheme }))
              .catch(() => null),
          ),
        );
        const schemeByItemId = new Map(
          schemeResults
            .filter((r): r is { itemId: number; scheme: ExerciseScheme } => r !== null)
            .map((r) => [r.itemId, r.scheme]),
        );

        // 3. Fetch unique exercises to get names
        const uniqueExerciseIds = [...new Set(exerciseItems.map((i) => i.exerciseId!))];
        const exerciseResults = await Promise.all(
          uniqueExerciseIds.map((id) => userApi.fetchUserExercise(id).catch(() => null)),
        );
        const exerciseNames: Record<number, string> = {};
        for (const exercise of exerciseResults) {
          if (exercise) {
            exerciseNames[exercise.id] = exercise.names?.[0]?.name ?? `Exercise #${exercise.id}`;
          }
        }

        // 4. Build display map — keyed by scheme ID (for log exercise references)
        for (const item of exerciseItems) {
          const assignment = assignmentByItemId.get(item.id);
          const scheme = schemeByItemId.get(item.id);
          if (!assignment || !scheme) continue;

          const numSets = scheme.sets ?? 0;
          const sets: SetPreview[] = [];
          for (let i = 1; i <= numSets; i++) {
            sets.push({
              setNumber: i,
              targetReps: scheme.reps,
              targetWeight: scheme.weight,
              targetDuration: scheme.duration,
              targetDistance: scheme.distance,
              targetTime: scheme.targetTime,
              restAfterSeconds: i < numSets ? (scheme.restBetweenSets ?? 0) : null,
            });
          }
          display[assignment.exerciseSchemeId] = {
            name: exerciseNames[item.exerciseId!] ?? `Exercise #${item.exerciseId}`,
            summary: formatSchemeSummary(scheme),
            measurementType: scheme.measurementType,
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
            exerciseNames[exercise.id] = exercise.names?.[0]?.name ?? `Exercise #${exercise.id}`;
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
