import { Component, inject, computed, signal } from '@angular/core';
import { ActivatedRoute, Router, RouterLink } from '@angular/router';
import { toSignal } from '@angular/core/rxjs-interop';
import { injectQuery, injectMutation, QueryClient } from '@tanstack/angular-query-experimental';
import { CompendiumApiClient } from '$core/api-clients/compendium-api-client';
import { UserApiClient } from '$core/api-clients/user-api-client';
import { exerciseKeys, userExerciseKeys } from '$core/query-keys';
import { TranslocoDirective } from '@jsverse/transloco';
import { PageLayout } from '../../../layout/page-layout';
import { ConfirmDialog } from '$ui/confirm-dialog/confirm-dialog';

@Component({
  selector: 'app-exercise-detail',
  imports: [PageLayout, RouterLink, ConfirmDialog, TranslocoDirective],
  template: `
    <ng-container *transloco="let t">
      <app-page-layout
        [header]="exerciseQuery.data()?.name ?? 'Exercise'"
        [isPending]="exerciseQuery.isPending()"
        [errorMessage]="exerciseQuery.isError() ? exerciseQuery.error().message : undefined"
      >
        <div actions class="flex gap-2">
          @if (alreadyAdded(); as existing) {
            <a
              [routerLink]="['/user/exercises', existing.id]"
              class="rounded-md bg-gray-500 px-4 py-2 text-sm font-medium text-white hover:bg-gray-600"
            >
              {{ t('compendium.exercises.alreadyAdded') }}
            </a>
          } @else {
            <button
              type="button"
              (click)="addMutation.mutate()"
              [disabled]="addMutation.isPending()"
              class="rounded-md bg-green-600 px-4 py-2 text-sm font-medium text-white hover:bg-green-700 disabled:opacity-50"
            >
              {{
                addMutation.isPending() ? t('common.adding') : t('compendium.exercises.addToMine')
              }}
            </button>
          }
          @if (hasHistory()) {
            <a
              routerLink="./history"
              class="rounded-md bg-gray-600 px-4 py-2 text-sm font-medium text-white hover:bg-gray-700"
              >{{ t('common.history') }}</a
            >
          }
          <a
            routerLink="./edit"
            class="rounded-md bg-blue-600 px-4 py-2 text-sm font-medium text-white hover:bg-blue-700"
            >{{ t('common.edit') }}</a
          >
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
          [title]="t('compendium.exercises.deleteTitle')"
          [message]="
            t('compendium.exercises.deleteMessage', { name: exerciseQuery.data()?.name ?? '' })
          "
          [isPending]="deleteMutation.isPending()"
          (confirmed)="deleteMutation.mutate()"
          (cancelled)="showDeleteDialog.set(false)"
        />
        @if (exerciseQuery.data(); as exercise) {
          <dl class="grid grid-cols-1 gap-4 sm:grid-cols-2">
            <div>
              <dt class="text-sm font-medium text-gray-500 dark:text-gray-400">
                {{ t('fields.type') }}
              </dt>
              <dd class="text-sm text-gray-900 dark:text-gray-100">
                {{ t('enums.exerciseType.' + exercise.type) }}
              </dd>
            </div>
            <div>
              <dt class="text-sm font-medium text-gray-500 dark:text-gray-400">
                {{ t('fields.difficulty') }}
              </dt>
              <dd class="text-sm text-gray-900 dark:text-gray-100">
                {{ t('enums.difficulty.' + exercise.technicalDifficulty) }}
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
            <div class="sm:col-span-2">
              <dt class="text-sm font-medium text-gray-500 dark:text-gray-400">
                {{ t('fields.description') }}
              </dt>
              <dd class="text-sm text-gray-900 dark:text-gray-100">{{ exercise.description }}</dd>
            </div>
          </dl>
        }
      </app-page-layout>
    </ng-container>
  `,
})
export class ExerciseDetail {
  private api = inject(CompendiumApiClient);
  private userApi = inject(UserApiClient);
  private router = inject(Router);
  private queryClient = inject(QueryClient);
  private params = toSignal(inject(ActivatedRoute).paramMap);

  private id = computed(() => Number(this.params()?.get('id')));

  showDeleteDialog = signal(false);

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

  userExercisesQuery = injectQuery(() => ({
    queryKey: userExerciseKeys.list(),
    queryFn: () => this.userApi.fetchUserExercises(),
  }));

  hasHistory = computed(() => (this.versionsQuery.data()?.length ?? 0) > 1);

  alreadyAdded = computed(() => {
    const templateId = this.exerciseQuery.data()?.templateId;
    const userExercises = this.userExercisesQuery.data();
    if (!templateId || !userExercises) return undefined;
    return userExercises.find((ue) => ue.templateId === templateId);
  });

  deleteMutation = injectMutation(() => ({
    mutationFn: () => this.api.deleteExercise(this.id()),
    onSuccess: () => {
      this.queryClient.invalidateQueries({ queryKey: exerciseKeys.all() });
      this.router.navigate(['/compendium/exercises']);
    },
  }));

  addMutation = injectMutation(() => ({
    mutationFn: () => {
      const exercise = this.exerciseQuery.data()!;
      return this.userApi.createUserExercise({
        name: exercise.name,
        type: exercise.type,
        force: exercise.force,
        primaryMuscles: exercise.primaryMuscles,
        secondaryMuscles: exercise.secondaryMuscles,
        technicalDifficulty: exercise.technicalDifficulty,
        bodyWeightScaling: exercise.bodyWeightScaling,
        suggestedMeasurementParadigms: exercise.suggestedMeasurementParadigms,
        description: exercise.description,
        instructions: exercise.instructions,
        images: exercise.images,
        alternativeNames: exercise.alternativeNames,
        authorName: exercise.authorName,
        authorUrl: exercise.authorUrl,
        version: exercise.version,
        parentExerciseId: exercise.parentExerciseId,
        templateId: exercise.templateId,
        equipmentIds: exercise.equipmentIds,
        public: false,
      });
    },
    onSuccess: (created) => {
      this.queryClient.invalidateQueries({ queryKey: userExerciseKeys.all() });
      this.router.navigate(['/user/exercises', created.id]);
    },
  }));
}
