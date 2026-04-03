import { Component, inject, computed } from '@angular/core';
import { ActivatedRoute, RouterLink } from '@angular/router';
import { DatePipe, JsonPipe } from '@angular/common';
import { toSignal } from '@angular/core/rxjs-interop';
import { injectQuery } from '@tanstack/angular-query-experimental';
import { TranslocoDirective } from '@jsverse/transloco';
import { CompendiumApiClient } from '$core/api-clients/compendium-api-client';
import { exerciseKeys } from '$core/query-keys';
import { PageLayout } from '../../../layout/page-layout';

@Component({
  selector: 'app-exercise-history',
  imports: [PageLayout, RouterLink, DatePipe, JsonPipe, TranslocoDirective],
  template: `
    <ng-container *transloco="let t">
      <app-page-layout
        [header]="
          (exerciseQuery.data()?.names?.[0]?.name ?? 'Exercise') +
          ' — ' +
          t('compendium.exercises.historyTitle')
        "
        [isPending]="exerciseQuery.isPending() || versionsQuery.isPending()"
        [errorMessage]="
          exerciseQuery.isError()
            ? exerciseQuery.error().message
            : versionsQuery.isError()
              ? versionsQuery.error().message
              : undefined
        "
      >
        <a
          actions
          routerLink=".."
          class="rounded-md bg-gray-600 px-4 py-2 text-sm font-medium text-white hover:bg-gray-700"
          >{{ t('common.back') }}</a
        >
        @if (versionsQuery.data(); as versions) {
          <div class="space-y-4">
            @for (entry of versions; track entry.version) {
              <div
                class="rounded-lg border border-gray-200 bg-white p-4 dark:border-gray-700 dark:bg-gray-800"
              >
                <div class="mb-2 flex items-center justify-between">
                  <span class="text-sm font-semibold text-gray-900 dark:text-gray-100">{{
                    t('compendium.exercises.versionLabel', { version: entry.version })
                  }}</span>
                  <span class="text-xs text-gray-500 dark:text-gray-400">
                    {{ entry.changedAt | date: 'medium'
                    }}{{ t('compendium.exercises.byLabel', { changedBy: entry.changedBy }) }}
                  </span>
                </div>
                <pre
                  class="overflow-x-auto rounded bg-gray-50 p-2 text-xs text-gray-700 dark:bg-gray-900 dark:text-gray-300"
                  >{{ entry.snapshot | json }}</pre
                >
              </div>
            }
          </div>
        }
      </app-page-layout>
    </ng-container>
  `,
})
export class ExerciseHistory {
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
}
