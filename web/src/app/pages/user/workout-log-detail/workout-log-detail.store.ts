import { inject } from '@angular/core';
import { signalStore, withState, withMethods, patchState } from '@ngrx/signals';
import { UserApiClient } from '$core/api-clients/user-api-client';
import { WorkoutLogSection } from '$generated/user-models';

interface WorkoutLogDetailState {
  exerciseNames: Record<number, string>;
  isLoading: boolean;
}

const initialState: WorkoutLogDetailState = {
  exerciseNames: {},
  isLoading: false,
};

export const WorkoutLogDetailStore = signalStore(
  withState(initialState),
  withMethods((store) => {
    const userApi = inject(UserApiClient);

    return {
      async loadExerciseNames(sections: WorkoutLogSection[]) {
        patchState(store, { isLoading: true });

        const names: Record<number, string> = {};

        // Collect unique scheme IDs from all exercises
        const schemeIds = [
          ...new Set(
            sections.flatMap((s) => (s.exercises ?? []).map((ex) => ex.sourceExerciseSchemeId)),
          ),
        ];

        // Fetch all schemes in parallel
        const schemes = await Promise.all(
          schemeIds.map((id) =>
            userApi
              .fetchExerciseScheme(id)
              .then((scheme) => ({ schemeId: id, scheme }))
              .catch(() => null),
          ),
        );
        const validSchemes = schemes.filter((s): s is NonNullable<typeof s> => s !== null);

        // Fetch unique exercises to get names directly
        const uniqueExerciseIds = [...new Set(validSchemes.map((s) => s.scheme.exerciseId))];
        const exerciseResults = await Promise.all(
          uniqueExerciseIds.map((id) => userApi.fetchUserExercise(id).catch(() => null)),
        );

        // Build exercise name map
        const exerciseNameMap: Record<number, string> = {};
        for (const exercise of exerciseResults) {
          if (exercise) {
            exerciseNameMap[exercise.id] = exercise.name;
          }
        }

        // Map scheme IDs to exercise names
        for (const item of validSchemes) {
          names[item.schemeId] =
            exerciseNameMap[item.scheme.exerciseId] ?? `Exercise #${item.scheme.exerciseId}`;
        }

        patchState(store, { exerciseNames: names, isLoading: false });
      },
    };
  }),
);
