import { Component, inject, computed, signal, viewChild } from '@angular/core';
import { ActivatedRoute, RouterLink, Router } from '@angular/router';
import { toSignal } from '@angular/core/rxjs-interop';
import { TranslocoDirective, TranslocoService } from '@jsverse/transloco';
import { UserApiClient } from '$core/api-clients/user-api-client';
import { ExerciseConfig } from '$ui/exercise-config/exercise-config';
import { ExerciseRunner } from '$ui/exercise-runner/exercise-runner';
import { HlmButton } from '@spartan-ng/helm/button';
import { HlmSeparator } from '@spartan-ng/helm/separator';
import { PageLayout } from '../../../layout/page-layout';

@Component({
  selector: 'app-exercise-track',
  imports: [
    PageLayout,
    RouterLink,
    ExerciseConfig,
    ExerciseRunner,
    HlmButton,
    HlmSeparator,
    TranslocoDirective,
  ],
  template: `
    <ng-container *transloco="let t">
      <app-page-layout [header]="t('user.exerciseTrack.title')" [isPending]="false">
        <a
          actions
          routerLink="/user/exercises"
          class="rounded-md border border-gray-300 px-4 py-2 text-sm font-medium text-gray-700 hover:bg-gray-50 dark:border-gray-600 dark:text-gray-300 dark:hover:bg-gray-800"
        >
          {{ t('common.back') }}
        </a>

        <!-- Phase 1: Exercise configuration (exercise preselected) -->
        <app-exercise-config #exerciseConfig [preselectedExerciseId]="exerciseId()" />

        <hlm-separator class="my-4" />

        <!-- Phase 2: Set preview -->
        @if (exerciseConfig.model().userExerciseId && exerciseConfig.model().sets) {
          <app-exercise-runner
            [exerciseName]="exerciseConfig.selectedExerciseName()"
            [measurementType]="exerciseConfig.model().measurementType"
            [setCount]="exerciseConfig.model().sets!"
            [defaultReps]="exerciseConfig.model().reps"
            [defaultWeight]="exerciseConfig.model().weight"
            [defaultDuration]="exerciseConfig.model().duration"
            [defaultDistance]="exerciseConfig.model().distance"
            [defaultRest]="exerciseConfig.model().restBetweenSets"
          />
        } @else {
          <div class="py-4 text-center text-sm text-gray-400 dark:text-gray-500">
            {{ t('user.exerciseTrack.configureHint') }}
          </div>
        }

        @if (errorMessage()) {
          <div
            class="mt-2 rounded-md border border-red-300 bg-red-50 px-3 py-2 text-sm text-red-700 dark:border-red-700 dark:bg-red-900/20 dark:text-red-400"
          >
            {{ errorMessage() }}
          </div>
        }

        <div class="mt-6">
          <button
            hlmBtn
            (click)="startWorkout()"
            [disabled]="!exerciseConfig.canConfirm() || isStarting()"
            class="w-full"
          >
            @if (isStarting()) {
              {{ t('common.starting') }}
            } @else {
              {{ t('user.exerciseTrack.startWorkout') }}
            }
          </button>
        </div>
      </app-page-layout>
    </ng-container>
  `,
})
export class ExerciseTrack {
  private userApi = inject(UserApiClient);
  private router = inject(Router);
  private route = inject(ActivatedRoute);
  private transloco = inject(TranslocoService);
  private params = toSignal(this.route.paramMap);

  exerciseConfig = viewChild.required<ExerciseConfig>('exerciseConfig');

  exerciseId = computed(() => Number(this.params()?.get('id')));

  isStarting = signal(false);
  errorMessage = signal('');

  async startWorkout() {
    this.isStarting.set(true);
    this.errorMessage.set('');
    try {
      const exerciseName = this.exerciseConfig().selectedExerciseName();

      // 1. Create scheme from config
      const scheme = await this.exerciseConfig().confirm();

      // 2. Create workout log in planning status
      const log = await this.userApi.createWorkoutLog({
        name: exerciseName || this.transloco.translate('user.exerciseTrack.title'),
      });

      // 3. Create section
      const section = await this.userApi.createWorkoutLogSection({
        workoutLogId: log.id,
        type: 'main',
        label: exerciseName || this.transloco.translate('user.exerciseTrack.title'),
        position: 0,
      });

      // 4. Create exercise (backend auto-creates sets from scheme)
      await this.userApi.createWorkoutLogExercise({
        workoutLogSectionId: section.id,
        sourceExerciseSchemeId: scheme.id,
        position: 0,
      });

      // 5. Start the workout (transitions everything to in_progress)
      await this.userApi.startWorkoutLog(log.id);

      // 6. Navigate to workout log detail for tracking
      this.router.navigate(['/user/workout-logs', log.id]);
    } catch (err) {
      console.error('Failed to start quick track workout:', err);
      this.errorMessage.set(err instanceof Error ? err.message : 'Failed to start workout');
    } finally {
      this.isStarting.set(false);
    }
  }
}
