import { inject } from '@angular/core';
import { signalStore, withState, withMethods, patchState } from '@ngrx/signals';
import { UserApiClient } from '$core/api-clients/user-api-client';
import { CompendiumApiClient } from '$core/api-clients/compendium-api-client';
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
    const compendiumApi = inject(CompendiumApiClient);

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

        // Fetch unique user exercises
        const uniqueUserExerciseIds = [
          ...new Set(validSchemes.map((s) => s.scheme.userExerciseId)),
        ];
        const userExercises = await Promise.all(
          uniqueUserExerciseIds.map((id) => userApi.fetchUserExercise(id).catch(() => null)),
        );

        // Fetch exercise names from compendium
        const exerciseNameMap: Record<number, string> = {};
        await Promise.all(
          userExercises
            .filter((ue) => ue !== null)
            .map(async (ue) => {
              try {
                const version = await compendiumApi.fetchExerciseVersion(
                  ue.compendiumExerciseId,
                  ue.compendiumVersion,
                );
                exerciseNameMap[ue.id] = version.snapshot?.name ?? `Exercise #${ue.id}`;
              } catch {
                exerciseNameMap[ue.id] = `Exercise #${ue.id}`;
              }
            }),
        );

        // Map scheme IDs to exercise names
        for (const item of validSchemes) {
          names[item.schemeId] =
            exerciseNameMap[item.scheme.userExerciseId] ??
            `Exercise #${item.scheme.userExerciseId}`;
        }

        patchState(store, { exerciseNames: names, isLoading: false });
      },
    };
  }),
);
