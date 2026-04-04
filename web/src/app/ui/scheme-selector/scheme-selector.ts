import { Component, inject, input, output, computed } from '@angular/core';
import { injectQuery } from '@tanstack/angular-query-experimental';
import { UserApiClient } from '$core/api-clients/user-api-client';
import { exerciseSchemeKeys } from '$core/query-keys';
import { formatSchemeSummary } from '$core/scheme-utils';
import { TranslocoDirective } from '@jsverse/transloco';

@Component({
  selector: 'app-scheme-selector',
  imports: [TranslocoDirective],
  template: `
    <div *transloco="let t" class="mt-3">
      <span class="mb-1 block text-xs font-medium text-gray-700 dark:text-gray-300">{{
        t('user.workouts.schemeLabel')
      }}</span>
      @if (exerciseId()) {
        <div class="flex flex-wrap gap-1.5">
          <button
            type="button"
            (click)="createRequested.emit()"
            class="rounded-md border border-dashed border-gray-400 px-2.5 py-1 text-xs font-medium text-gray-600 transition-colors hover:border-gray-600 hover:text-gray-800 dark:border-gray-500 dark:text-gray-400 dark:hover:border-gray-300 dark:hover:text-gray-200"
          >
            + {{ t('common.new') }}
          </button>
          @for (scheme of schemes(); track scheme.id) {
            <button
              type="button"
              (click)="schemeSelected.emit(scheme.id)"
              class="rounded-md border px-2.5 py-1 text-xs font-medium transition-colors"
              [class]="
                scheme.id === selectedSchemeId()
                  ? 'border-blue-600 bg-blue-600 text-white'
                  : 'border-gray-300 bg-white text-gray-700 hover:bg-gray-50 dark:border-gray-600 dark:bg-gray-800 dark:text-gray-300 dark:hover:bg-gray-700'
              "
            >
              {{ summary(scheme) }}
            </button>
          }
        </div>
        @if (schemesQuery.isSuccess() && schemes().length === 0 && !selectedSchemeId()) {
          <p class="mt-1.5 text-xs text-amber-600 dark:text-amber-400">
            {{ t('user.workouts.noSchemeWarning') }}
          </p>
        }
      }
    </div>
  `,
})
export class SchemeSelector {
  private userApi = inject(UserApiClient);

  exerciseId = input<number | null>(null);
  selectedSchemeId = input<number | null>(null);

  schemeSelected = output<number | null>();
  createRequested = output<void>();

  schemesQuery = injectQuery(() => ({
    queryKey: exerciseSchemeKeys.list({ exerciseId: this.exerciseId() ?? undefined }),
    queryFn: () => this.userApi.fetchExerciseSchemes({ exerciseId: this.exerciseId()! }),
    enabled: this.exerciseId() != null,
  }));

  schemes = computed(() => this.schemesQuery.data() ?? []);

  summary = formatSchemeSummary;
}
