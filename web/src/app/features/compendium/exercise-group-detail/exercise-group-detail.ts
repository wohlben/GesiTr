import { Component, inject, computed } from '@angular/core';
import { ActivatedRoute } from '@angular/router';
import { toSignal } from '@angular/core/rxjs-interop';
import { injectQuery } from '@tanstack/angular-query-experimental';
import { CompendiumApiClient } from '$core/api-clients/compendium-api-client';
import { PageLayout } from '../../../layout/page-layout';

@Component({
  selector: 'app-exercise-group-detail',
  imports: [PageLayout],
  template: `
    <app-page-layout
      [header]="groupQuery.data()?.name ?? 'Exercise Group'"
      [isPending]="groupQuery.isPending()"
      [errorMessage]="groupQuery.isError() ? groupQuery.error().message : undefined"
    >
      @if (groupQuery.data(); as group) {
        <dl class="grid grid-cols-1 gap-4 sm:grid-cols-2">
          <div class="sm:col-span-2">
            <dt class="text-sm font-medium text-gray-500 dark:text-gray-400">Description</dt>
            <dd class="text-sm text-gray-900 dark:text-gray-100">{{ group.description }}</dd>
          </div>
          <div>
            <dt class="text-sm font-medium text-gray-500 dark:text-gray-400">Created By</dt>
            <dd class="text-sm text-gray-900 dark:text-gray-100">{{ group.createdBy }}</dd>
          </div>
        </dl>
      }
    </app-page-layout>
  `,
})
export class ExerciseGroupDetail {
  private api = inject(CompendiumApiClient);
  private params = toSignal(inject(ActivatedRoute).paramMap);

  private id = computed(() => Number(this.params()?.get('id')));

  groupQuery = injectQuery(() => ({
    queryKey: ['exercise-group', this.id()],
    queryFn: () => this.api.fetchExerciseGroup(this.id()),
    enabled: !!this.id(),
  }));
}
