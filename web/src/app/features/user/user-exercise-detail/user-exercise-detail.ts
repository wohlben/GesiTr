import { Component, inject, computed, signal } from '@angular/core';
import { ActivatedRoute, Router } from '@angular/router';
import { toSignal } from '@angular/core/rxjs-interop';
import {
  injectQuery,
  injectMutation,
  injectQueryClient,
} from '@tanstack/angular-query-experimental';
import { UserApiClient } from '$core/api-clients/user-api-client';
import { CompendiumApiClient } from '$core/api-clients/compendium-api-client';
import { userExerciseKeys } from '$core/query-keys';
import { PageLayout } from '../../../layout/page-layout';
import { ConfirmDialog } from '$ui/confirm-dialog/confirm-dialog';
import { Exercise } from '$generated/models';

@Component({
  selector: 'app-user-exercise-detail',
  imports: [PageLayout, ConfirmDialog],
  template: `
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
          Delete
        </button>
      </div>
      <app-confirm-dialog
        [open]="showDeleteDialog()"
        title="Remove Exercise"
        [message]="deleteMessage()"
        [isPending]="deleteMutation.isPending()"
        (confirmed)="deleteMutation.mutate()"
        (cancelled)="showDeleteDialog.set(false)"
      />
      @if (snapshot(); as exercise) {
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
          <div>
            <dt class="text-sm font-medium text-gray-500 dark:text-gray-400">Force</dt>
            <dd class="text-sm text-gray-900 dark:text-gray-100">
              {{ exercise.force?.join(', ') }}
            </dd>
          </div>
          <div>
            <dt class="text-sm font-medium text-gray-500 dark:text-gray-400">Version</dt>
            <dd class="text-sm text-gray-900 dark:text-gray-100">
              v{{ detailQuery.data()?.userExercise?.compendiumVersion }}
            </dd>
          </div>
          <div class="sm:col-span-2">
            <dt class="text-sm font-medium text-gray-500 dark:text-gray-400">Description</dt>
            <dd class="text-sm text-gray-900 dark:text-gray-100">{{ exercise.description }}</dd>
          </div>
          @if (exercise.instructions?.length) {
            <div class="sm:col-span-2">
              <dt class="text-sm font-medium text-gray-500 dark:text-gray-400">Instructions</dt>
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
  `,
})
export class UserExerciseDetail {
  private userApi = inject(UserApiClient);
  private compendiumApi = inject(CompendiumApiClient);
  private router = inject(Router);
  private queryClient = injectQueryClient();
  private params = toSignal(inject(ActivatedRoute).paramMap);

  private id = computed(() => Number(this.params()?.get('id')));

  showDeleteDialog = signal(false);

  detailQuery = injectQuery(() => ({
    queryKey: userExerciseKeys.detail(this.id()),
    queryFn: async () => {
      const userExercise = await this.userApi.fetchUserExercise(this.id());
      const version = await this.compendiumApi.fetchExerciseVersion(
        userExercise.compendiumExerciseId,
        userExercise.compendiumVersion,
      );
      return { userExercise, version };
    },
    enabled: !!this.id(),
  }));

  snapshot = computed(() => this.detailQuery.data()?.version.snapshot as Exercise | undefined);

  exerciseName = computed(() => this.snapshot()?.name ?? 'Exercise');

  deleteMessage = computed(
    () => `Are you sure you want to remove '${this.exerciseName()}' from your exercises?`,
  );

  deleteMutation = injectMutation(() => ({
    mutationFn: () => this.userApi.deleteUserExercise(this.id()),
    onSuccess: () => {
      this.queryClient.invalidateQueries({ queryKey: userExerciseKeys.all() });
      this.router.navigate(['/user/exercises']);
    },
  }));
}
