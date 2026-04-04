import { computed, inject, Signal } from '@angular/core';
import { signalStore, withComputed, withMethods, patchState } from '@ngrx/signals';
import { withEntities, setEntity, setAllEntities } from '@ngrx/signals/entities';
import { QueryClient } from '@tanstack/angular-query-experimental';
import { CompendiumApiClient } from '$core/api-clients/compendium-api-client';
import { exerciseKeys } from '$core/query-keys';
import type { Exercise } from '$generated/models';

export enum LoadingState {
  LOADING = 'LOADING',
  ERROR = 'ERROR',
  READY = 'READY',
}

export type ExerciseEntry =
  | { id: number; loading: LoadingState.LOADING }
  | { id: number; loading: LoadingState.ERROR; error: string }
  | ({ loading: LoadingState.READY } & Exercise);

export const ExerciseStore = signalStore(
  { providedIn: 'root' },
  withEntities<ExerciseEntry>(),
  withComputed((store) => ({
    readyExercises: computed(() =>
      store
        .entities()
        .filter(
          (e): e is { loading: LoadingState.READY } & Exercise => e.loading === LoadingState.READY,
        ),
    ),
  })),
  withMethods((store) => {
    const queryClient = inject(QueryClient);
    const api = inject(CompendiumApiClient);

    return {
      exercise(id: number): Signal<ExerciseEntry | undefined> {
        return computed(() => store.entityMap()[id]);
      },

      setAllFromQuery(exercises: Exercise[]): void {
        const entries: ExerciseEntry[] = exercises.map((exercise) => ({
          ...exercise,
          loading: LoadingState.READY as const,
        }));
        patchState(store, setAllEntities<ExerciseEntry>(entries));
      },

      async loadExercises(ids: number[]): Promise<void> {
        const uniqueIds = [...new Set(ids)];
        const currentMap = store.entityMap();
        const idsToFetch = uniqueIds.filter((id) => {
          const existing = currentMap[id];
          return !existing || existing.loading === LoadingState.ERROR;
        });

        if (idsToFetch.length === 0) return;

        for (const id of idsToFetch) {
          patchState(store, setEntity<ExerciseEntry>({ id, loading: LoadingState.LOADING }));
        }

        await Promise.all(
          idsToFetch.map(async (id) => {
            try {
              const exercise = await queryClient.fetchQuery({
                queryKey: exerciseKeys.detail(id),
                queryFn: () => api.fetchExercise(id),
              });
              patchState(
                store,
                setEntity<ExerciseEntry>({ ...exercise, loading: LoadingState.READY }),
              );
            } catch (err) {
              patchState(
                store,
                setEntity<ExerciseEntry>({
                  id,
                  loading: LoadingState.ERROR,
                  error: err instanceof Error ? err.message : 'Unknown error',
                }),
              );
            }
          }),
        );
      },
    };
  }),
);
