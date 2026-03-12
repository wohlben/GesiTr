import { inject } from '@angular/core';
import { signalStore, withState, withMethods, patchState } from '@ngrx/signals';
import { UserApiClient } from '$core/api-clients/user-api-client';
import { CompendiumApiClient } from '$core/api-clients/compendium-api-client';
import { WorkoutSection, UserExerciseScheme } from '$generated/user-models';

export interface ExerciseDisplayInfo {
  name: string;
  summary: string;
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
    const compendiumApi = inject(CompendiumApiClient);

    return {
      async loadExerciseDisplay(sections: WorkoutSection[]) {
        patchState(store, { isLoadingDisplay: true });

        const display: Record<number, ExerciseDisplayInfo> = {};

        // 1. Fetch all schemes in parallel
        const schemeResults = await Promise.all(
          sections.flatMap((s) =>
            (s.exercises ?? []).map((ex) =>
              userApi
                .fetchExerciseScheme(ex.userExerciseSchemeId)
                .then((scheme) => ({ schemeId: ex.userExerciseSchemeId, scheme }))
                .catch(() => null),
            ),
          ),
        );
        const schemes = schemeResults.filter(
          (r): r is { schemeId: number; scheme: UserExerciseScheme } => r !== null,
        );

        // 2. Fetch unique user exercises to get compendium IDs
        const uniqueUserExerciseIds = [...new Set(schemes.map((s) => s.scheme.userExerciseId))];
        const userExerciseResults = await Promise.all(
          uniqueUserExerciseIds.map((id) => userApi.fetchUserExercise(id).catch(() => null)),
        );

        // 3. Fetch exercise versions for names in parallel
        const exerciseNames: Record<number, string> = {};
        await Promise.all(
          userExerciseResults
            .filter((ue) => ue !== null)
            .map(async (ue) => {
              try {
                const version = await compendiumApi.fetchExerciseVersion(
                  ue.compendiumExerciseId,
                  ue.compendiumVersion,
                );
                exerciseNames[ue.id] = version.snapshot?.name ?? `Exercise #${ue.id}`;
              } catch {
                exerciseNames[ue.id] = `Exercise #${ue.id}`;
              }
            }),
        );

        // 4. Build display map
        for (const item of schemes) {
          display[item.schemeId] = {
            name:
              exerciseNames[item.scheme.userExerciseId] ??
              `Exercise #${item.scheme.userExerciseId}`,
            summary: formatSchemeSummary(item.scheme),
          };
        }

        patchState(store, { exerciseDisplay: display, isLoadingDisplay: false });
      },
    };
  }),
);
