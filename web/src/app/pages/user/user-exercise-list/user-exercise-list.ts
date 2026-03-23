import { Component, inject } from '@angular/core';
import { RouterLink } from '@angular/router';
import { injectQuery } from '@tanstack/angular-query-experimental';
import { TranslocoDirective } from '@jsverse/transloco';
import { UserApiClient } from '$core/api-clients/user-api-client';
import { userExerciseKeys } from '$core/query-keys';
import { PageLayout } from '../../../layout/page-layout';
import { NgIcon, provideIcons } from '@ng-icons/core';
import { lucideListCheck } from '@ng-icons/lucide';

@Component({
  selector: 'app-user-exercise-list',
  imports: [PageLayout, RouterLink, NgIcon, TranslocoDirective],
  providers: [provideIcons({ lucideListCheck })],
  template: `
    <ng-container *transloco="let t">
      <app-page-layout
        [header]="t('user.exercises.title')"
        [isPending]="userExercisesQuery.isPending()"
        [errorMessage]="
          userExercisesQuery.isError() ? userExercisesQuery.error().message : undefined
        "
      >
        @if (userExercisesQuery.data(); as exercises) {
          @if (exercises.length === 0) {
            <p class="text-sm text-gray-500 dark:text-gray-400">
              {{ t('user.exercises.noResults') }}
            </p>
          } @else {
            <div class="overflow-x-auto">
              <table class="min-w-full divide-y divide-gray-200 dark:divide-gray-700">
                <thead>
                  <tr>
                    <th
                      class="px-4 py-3 text-left text-xs font-medium tracking-wider text-gray-500 uppercase dark:text-gray-400"
                    >
                      {{ t('fields.name') }}
                    </th>
                    <th
                      class="px-4 py-3 text-left text-xs font-medium tracking-wider text-gray-500 uppercase dark:text-gray-400"
                    >
                      {{ t('fields.type') }}
                    </th>
                    <th
                      class="px-4 py-3 text-left text-xs font-medium tracking-wider text-gray-500 uppercase dark:text-gray-400"
                    >
                      {{ t('fields.version') }}
                    </th>
                    <th class="px-4 py-3"></th>
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
                        v{{ item.version }}
                      </td>
                      <td class="px-4 py-3 text-right">
                        <a
                          [routerLink]="['./', item.id, 'track']"
                          (click)="$event.stopPropagation()"
                          class="inline-flex items-center rounded-md p-1.5 text-green-600 hover:bg-green-50 dark:text-green-400 dark:hover:bg-green-900/30"
                          [title]="t('user.exercises.quickTrack')"
                        >
                          <ng-icon name="lucideListCheck" class="text-lg" />
                        </a>
                      </td>
                    </tr>
                  }
                </tbody>
              </table>
            </div>
          }
        }
      </app-page-layout>
    </ng-container>
  `,
})
export class UserExerciseList {
  private userApi = inject(UserApiClient);

  userExercisesQuery = injectQuery(() => ({
    queryKey: userExerciseKeys.list(),
    queryFn: () => this.userApi.fetchUserExercises(),
  }));
}
