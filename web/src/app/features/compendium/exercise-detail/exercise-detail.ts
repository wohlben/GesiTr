import { Component, inject, computed } from '@angular/core';
import { ActivatedRoute, RouterLink } from '@angular/router';
import { toSignal } from '@angular/core/rxjs-interop';
import { injectQuery } from '@tanstack/angular-query-experimental';
import { CompendiumApiClient } from '$core/api-clients/compendium-api-client';
import { exerciseKeys } from '$core/query-keys';
import { PageLayout } from '../../../layout/page-layout';

@Component({
  selector: 'app-exercise-detail',
  imports: [PageLayout, RouterLink],
  template: `
    <app-page-layout
      [header]="exerciseQuery.data()?.name ?? 'Exercise'"
      [isPending]="exerciseQuery.isPending()"
      [errorMessage]="exerciseQuery.isError() ? exerciseQuery.error().message : undefined"
    >
      <div actions class="flex gap-2">
        @if (hasHistory()) {
          <a
            routerLink="./history"
            class="rounded-md bg-gray-600 px-4 py-2 text-sm font-medium text-white hover:bg-gray-700"
            >History</a
          >
        }
        <a
          routerLink="./edit"
          class="rounded-md bg-blue-600 px-4 py-2 text-sm font-medium text-white hover:bg-blue-700"
          >Edit</a
        >
      </div>
      @if (exerciseQuery.data(); as exercise) {
        <dl class="grid grid-cols-1 gap-4 sm:grid-cols-2">
          <div>
            <dt class="text-sm font-medium text-gray-500 dark:text-gray-400">Type</dt>
            <dd class="text-sm text-gray-900 dark:text-gray-100">{{ exercise.type }}</dd>
          </div>
          <div>
            <dt class="text-sm font-medium text-gray-500 dark:text-gray-400">Difficulty</dt>
            <dd class="text-sm text-gray-900 dark:text-gray-100">
              {{ exercise.technicalDifficulty }}
            </dd>
          </div>
          <div>
            <dt class="text-sm font-medium text-gray-500 dark:text-gray-400">Primary Muscles</dt>
            <dd class="text-sm text-gray-900 dark:text-gray-100">
              {{ exercise.primaryMuscles?.join(', ') }}
            </dd>
          </div>
          <div>
            <dt class="text-sm font-medium text-gray-500 dark:text-gray-400">Secondary Muscles</dt>
            <dd class="text-sm text-gray-900 dark:text-gray-100">
              {{ exercise.secondaryMuscles?.join(', ') }}
            </dd>
          </div>
          <div class="sm:col-span-2">
            <dt class="text-sm font-medium text-gray-500 dark:text-gray-400">Description</dt>
            <dd class="text-sm text-gray-900 dark:text-gray-100">{{ exercise.description }}</dd>
          </div>
        </dl>
      }
    </app-page-layout>
  `,
})
export class ExerciseDetail {
  private api = inject(CompendiumApiClient);
  private params = toSignal(inject(ActivatedRoute).paramMap);

  private id = computed(() => Number(this.params()?.get('id')));

  exerciseQuery = injectQuery(() => ({
    queryKey: exerciseKeys.detail(this.id()),
    queryFn: () => this.api.fetchExercise(this.id()),
    enabled: !!this.id(),
  }));

  versionsQuery = injectQuery(() => ({
    queryKey: exerciseKeys.versions(this.id()),
    queryFn: () => this.api.fetchExerciseVersions(this.id()),
    enabled: !!this.id(),
  }));

  hasHistory = computed(() => (this.versionsQuery.data()?.length ?? 0) > 1);
}
