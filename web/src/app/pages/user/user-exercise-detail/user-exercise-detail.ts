import { Component, inject, computed, signal } from '@angular/core';
import { ActivatedRoute, Router } from '@angular/router';
import { toSignal } from '@angular/core/rxjs-interop';
import { injectQuery, injectMutation, QueryClient } from '@tanstack/angular-query-experimental';
import { TranslocoDirective, TranslocoService } from '@jsverse/transloco';
import { UserApiClient } from '$core/api-clients/user-api-client';
import { userExerciseKeys } from '$core/query-keys';
import { PageLayout } from '../../../layout/page-layout';
import { ConfirmDialog } from '$ui/confirm-dialog/confirm-dialog';

@Component({
  selector: 'app-user-exercise-detail',
  imports: [PageLayout, ConfirmDialog, TranslocoDirective],
  template: `
    <ng-container *transloco="let t">
      <app-page-layout
        [header]="exerciseName()"
        [isPending]="detailQuery.isPending()"
        [errorMessage]="detailQuery.isError() ? detailQuery.error().message : undefined"
      >
        <div actions class="flex gap-2">
          <button
            type="button"
            (click)="showDeleteDialog.set(true)"
            class="rounded-md bg-red-600 px-4 py-2 text-sm font-medium text-white hover:bg-red-700"
          >
            {{ t('common.delete') }}
          </button>
        </div>
        <app-confirm-dialog
          [open]="showDeleteDialog()"
          [title]="t('user.exercises.deleteTitle')"
          [message]="t('user.exercises.deleteMessage', { name: exerciseName() })"
          [isPending]="deleteMutation.isPending()"
          (confirmed)="deleteMutation.mutate()"
          (cancelled)="showDeleteDialog.set(false)"
        />
        @if (detailQuery.data(); as exercise) {
          <dl class="grid grid-cols-1 gap-4 sm:grid-cols-2">
            <div>
              <dt class="text-sm font-medium text-gray-500 dark:text-gray-400">
                {{ t('fields.type') }}
              </dt>
              <dd class="text-sm text-gray-900 dark:text-gray-100">{{ exercise.type }}</dd>
            </div>
            <div>
              <dt class="text-sm font-medium text-gray-500 dark:text-gray-400">
                {{ t('fields.difficulty') }}
              </dt>
              <dd class="text-sm text-gray-900 dark:text-gray-100">
                {{ exercise.technicalDifficulty }}
              </dd>
            </div>
            <div>
              <dt class="text-sm font-medium text-gray-500 dark:text-gray-400">
                {{ t('fields.primaryMuscles') }}
              </dt>
              <dd class="text-sm text-gray-900 dark:text-gray-100">
                {{ exercise.primaryMuscles?.join(', ') }}
              </dd>
            </div>
            <div>
              <dt class="text-sm font-medium text-gray-500 dark:text-gray-400">
                {{ t('fields.secondaryMuscles') }}
              </dt>
              <dd class="text-sm text-gray-900 dark:text-gray-100">
                {{ exercise.secondaryMuscles?.join(', ') }}
              </dd>
            </div>
            <div>
              <dt class="text-sm font-medium text-gray-500 dark:text-gray-400">
                {{ t('fields.force') }}
              </dt>
              <dd class="text-sm text-gray-900 dark:text-gray-100">
                {{ exercise.force?.join(', ') }}
              </dd>
            </div>
            <div>
              <dt class="text-sm font-medium text-gray-500 dark:text-gray-400">
                {{ t('fields.version') }}
              </dt>
              <dd class="text-sm text-gray-900 dark:text-gray-100">v{{ exercise.version }}</dd>
            </div>
            <div class="sm:col-span-2">
              <dt class="text-sm font-medium text-gray-500 dark:text-gray-400">
                {{ t('fields.description') }}
              </dt>
              <dd class="text-sm text-gray-900 dark:text-gray-100">{{ exercise.description }}</dd>
            </div>
            @if (exercise.instructions?.length) {
              <div class="sm:col-span-2">
                <dt class="text-sm font-medium text-gray-500 dark:text-gray-400">
                  {{ t('fields.instructions') }}
                </dt>
                <dd class="text-sm text-gray-900 dark:text-gray-100">
                  <ol class="list-inside list-decimal space-y-1">
                    @for (step of exercise.instructions; track $index) {
                      <li>{{ step }}</li>
                    }
                  </ol>
                </dd>
              </div>
            }
          </dl>
        }
      </app-page-layout>
    </ng-container>
  `,
})
export class UserExerciseDetail {
  private userApi = inject(UserApiClient);
  private router = inject(Router);
  private queryClient = inject(QueryClient);
  private transloco = inject(TranslocoService);
  private params = toSignal(inject(ActivatedRoute).paramMap);

  private id = computed(() => Number(this.params()?.get('id')));

  showDeleteDialog = signal(false);

  detailQuery = injectQuery(() => ({
    queryKey: userExerciseKeys.detail(this.id()),
    queryFn: () => this.userApi.fetchUserExercise(this.id()),
    enabled: !!this.id(),
  }));

  exerciseName = computed(
    () => this.detailQuery.data()?.name ?? this.transloco.translate('common.loading'),
  );

  deleteMutation = injectMutation(() => ({
    mutationFn: () => this.userApi.deleteUserExercise(this.id()),
    onSuccess: () => {
      this.queryClient.invalidateQueries({ queryKey: userExerciseKeys.all() });
      this.router.navigate(['/user/exercises']);
    },
  }));
}
