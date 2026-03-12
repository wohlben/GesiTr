import { Component, inject, computed } from '@angular/core';
import { ActivatedRoute } from '@angular/router';
import { toSignal } from '@angular/core/rxjs-interop';
import { JsonPipe } from '@angular/common';
import { injectQuery } from '@tanstack/angular-query-experimental';
import { UserApiClient } from '$core/api-clients/user-api-client';
import { workoutLogKeys } from '$core/query-keys';
import { PageLayout } from '../../../layout/page-layout';

@Component({
  selector: 'app-workout-log-detail',
  imports: [PageLayout, JsonPipe],
  template: `
    <app-page-layout
      header="Workout Log"
      [isPending]="logQuery.isPending()"
      [errorMessage]="logQuery.isError() ? logQuery.error().message : undefined"
    >
      @if (logQuery.data(); as log) {
        <pre
          class="overflow-x-auto rounded-lg bg-gray-100 p-4 text-sm text-gray-800 dark:bg-gray-800 dark:text-gray-200"
          >{{ log | json }}</pre
        >
      }
    </app-page-layout>
  `,
})
export class WorkoutLogDetail {
  private userApi = inject(UserApiClient);
  private route = inject(ActivatedRoute);
  private params = toSignal(this.route.paramMap);

  private id = computed(() => Number(this.params()?.get('id')));

  logQuery = injectQuery(() => ({
    queryKey: workoutLogKeys.detail(this.id()),
    queryFn: () => this.userApi.fetchWorkoutLog(this.id()),
    enabled: !!this.id(),
  }));
}
