import { Component, inject, computed } from '@angular/core';
import { RouterLink } from '@angular/router';
import { injectQuery } from '@tanstack/angular-query-experimental';
import { injectQueries } from '@tanstack/angular-query-experimental/inject-queries-experimental';
import { UserApiClient } from '$core/api-clients/user-api-client';
import { CompendiumApiClient } from '$core/api-clients/compendium-api-client';
import { userExerciseKeys, exerciseKeys } from '$core/query-keys';
import { PageLayout } from '../../../layout/page-layout';

@Component({
  selector: 'app-user-exercise-list',
  imports: [PageLayout, RouterLink],
  template: `
    <app-page-layout
      header="My Exercises"
      [isPending]="userExercisesQuery.isPending()"
      [errorMessage]="userExercisesQuery.isError() ? userExercisesQuery.error().message : undefined"
    >
      @if (enrichedExercises(); as exercises) {
        @if (exercises.length === 0) {
          <p class="text-sm text-gray-500 dark:text-gray-400">No exercises imported yet.</p>
        } @else {
          <div class="overflow-x-auto">
            <table class="min-w-full divide-y divide-gray-200 dark:divide-gray-700">
              <thead>
                <tr>
                  <th
                    class="px-4 py-3 text-left text-xs font-medium tracking-wider text-gray-500 uppercase dark:text-gray-400"
                  >
                    Name
                  </th>
                  <th
                    class="px-4 py-3 text-left text-xs font-medium tracking-wider text-gray-500 uppercase dark:text-gray-400"
                  >
                    Type
                  </th>
                  <th
                    class="px-4 py-3 text-left text-xs font-medium tracking-wider text-gray-500 uppercase dark:text-gray-400"
                  >
                    Version
                  </th>
                </tr>
              </thead>
              <tbody class="divide-y divide-gray-200 dark:divide-gray-700">
                @for (item of exercises; track item.id) {
                  <tr
                    class="hover:bg-gray-50 dark:hover:bg-gray-800/50"
                    [routerLink]="['./', item.id]"
                    class="cursor-pointer"
                  >
                    <td class="px-4 py-3 text-sm text-gray-900 dark:text-gray-100">
                      {{ item.name }}
                    </td>
                    <td class="px-4 py-3 text-sm text-gray-600 dark:text-gray-400">
                      {{ item.type }}
                    </td>
                    <td class="px-4 py-3 text-sm text-gray-600 dark:text-gray-400">
                      v{{ item.compendiumVersion }}
                    </td>
                  </tr>
                }
              </tbody>
            </table>
          </div>
        }
      }
    </app-page-layout>
  `,
})
export class UserExerciseList {
  private userApi = inject(UserApiClient);
  private compendiumApi = inject(CompendiumApiClient);

  userExercisesQuery = injectQuery(() => ({
    queryKey: userExerciseKeys.list(),
    queryFn: () => this.userApi.fetchUserExercises(),
  }));

  private snapshotQueries = injectQueries(() => ({
    queries: (this.userExercisesQuery.data() ?? []).map((ue) => ({
      queryKey: exerciseKeys.version(ue.compendiumExerciseId, ue.compendiumVersion),
      queryFn: () =>
        this.compendiumApi.fetchExerciseVersion(ue.compendiumExerciseId, ue.compendiumVersion),
      staleTime: Infinity,
    })),
  }));

  enrichedExercises = computed(() => {
    const userExercises = this.userExercisesQuery.data();
    if (!userExercises) return undefined;

    const snapshots = this.snapshotQueries();

    return userExercises.map((ue, i) => {
      const versionEntry = snapshots[i]?.data();
      const exercise = versionEntry?.snapshot;
      return {
        id: ue.id,
        compendiumVersion: ue.compendiumVersion,
        name: exercise?.name ?? '...',
        type: exercise?.type ?? '',
      };
    });
  });
}
