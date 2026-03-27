import { Component, input, output } from '@angular/core';
import { FormField } from '@angular/forms/signals';
import { CdkDragDrop, CdkDrag, CdkDropList, CdkDragHandle } from '@angular/cdk/drag-drop';
import { TranslocoDirective } from '@jsverse/transloco';
import { BrnSelectImports } from '@spartan-ng/brain/select';
import { HlmSelectImports } from '@spartan-ng/helm/select';
import { HlmInput } from '@spartan-ng/helm/input';
import { WorkoutSectionTypeMain, WorkoutSectionTypeSupplementary } from '$generated/user-models';
import { ExerciseDisplayInfo } from './workout-start.store';
import { WorkoutStartExerciseItem } from './workout-start-exercise-item';

@Component({
  selector: 'app-workout-start-section',
  imports: [
    FormField,
    CdkDropList,
    CdkDrag,
    CdkDragHandle,
    BrnSelectImports,
    HlmSelectImports,
    HlmInput,
    TranslocoDirective,
    WorkoutStartExerciseItem,
  ],
  template: `
    <ng-container *transloco="let t">
      <div class="rounded-lg border border-gray-200 p-4 dark:border-gray-700">
        <div class="mb-3 flex items-center justify-between">
          <div class="flex items-center gap-2">
            @if (!readonly()) {
              <!-- Section drag handle -->
              <div
                cdkDragHandle
                class="flex cursor-grab flex-col gap-0.5 px-1 py-1 text-gray-400 active:cursor-grabbing dark:text-gray-500"
              >
                <div class="flex gap-0.5">
                  <div class="h-1 w-1 rounded-full bg-current"></div>
                  <div class="h-1 w-1 rounded-full bg-current"></div>
                </div>
                <div class="flex gap-0.5">
                  <div class="h-1 w-1 rounded-full bg-current"></div>
                  <div class="h-1 w-1 rounded-full bg-current"></div>
                </div>
                <div class="flex gap-0.5">
                  <div class="h-1 w-1 rounded-full bg-current"></div>
                  <div class="h-1 w-1 rounded-full bg-current"></div>
                </div>
              </div>
            }
            <h3 class="text-sm font-semibold text-gray-900 dark:text-gray-100">
              {{ t('user.workouts.sectionLabel', { n: sectionIndex() + 1 }) }}
            </h3>
          </div>
          @if (!readonly()) {
            <button
              type="button"
              (click)="removed.emit()"
              class="text-sm text-red-600 hover:text-red-800 dark:text-red-400 dark:hover:text-red-300"
            >
              {{ t('common.remove') }}
            </button>
          }
        </div>

        <div class="mb-3 grid grid-cols-1 gap-3 sm:grid-cols-2">
          <div>
            <span class="block text-xs font-medium text-gray-700 dark:text-gray-300">{{
              t('fields.type')
            }}</span>
            @if (readonly()) {
              <span class="mt-1 block text-sm text-gray-900 dark:text-gray-100">{{
                section().type().value() === SECTION_TYPE_MAIN
                  ? t('enums.workoutSectionType.main')
                  : t('enums.workoutSectionType.supplementary')
              }}</span>
            } @else {
              <brn-select
                [formField]="section().type"
                (valueChange)="sectionChanged.emit()"
                class="mt-1"
                hlm
              >
                <hlm-select-trigger class="w-full">
                  <hlm-select-value />
                </hlm-select-trigger>
                <hlm-select-content>
                  <hlm-option [value]="SECTION_TYPE_MAIN">{{
                    t('enums.workoutSectionType.main')
                  }}</hlm-option>
                  <hlm-option [value]="SECTION_TYPE_SUPPLEMENTARY">{{
                    t('enums.workoutSectionType.supplementary')
                  }}</hlm-option>
                </hlm-select-content>
              </brn-select>
            }
          </div>
          <div>
            <span class="block text-xs font-medium text-gray-700 dark:text-gray-300">{{
              t('fields.label')
            }}</span>
            @if (readonly()) {
              <span class="mt-1 block text-sm text-gray-900 dark:text-gray-100">{{
                section().label().value() || '—'
              }}</span>
            } @else {
              <input
                hlmInput
                [formField]="section().label"
                (change)="sectionChanged.emit()"
                class="mt-1"
              />
            }
          </div>
        </div>

        <!-- Exercise cards -->
        <div
          cdkDropList
          [cdkDropListData]="section().exercises"
          [cdkDropListDisabled]="readonly()"
          (cdkDropListDropped)="onExerciseDrop($event)"
        >
          @for (
            exercise of section().exercises;
            track $index;
            let ei = $index;
            let lastEx = $last
          ) {
            <app-workout-start-exercise-item
              cdkDrag
              [exercise]="exercise"
              [displayInfo]="exerciseDisplayMap()[exercise.id().value()!]"
              [isLast]="lastEx"
              [readonly]="readonly()"
              (removed)="exerciseRemoved.emit({ exerciseIndex: ei })"
              (exerciseChanged)="exerciseChanged.emit({ exerciseIndex: ei })"
              (setChanged)="setChanged.emit({ exerciseIndex: ei, setIndex: $event.setIndex })"
            />
          }
        </div>

        <!-- Pending exercise groups -->
        @for (group of section().pendingGroups; track $index; let gi = $index) {
          <div
            class="mt-2 rounded-md border-2 border-dashed border-amber-300 bg-amber-50 p-3 dark:border-amber-600 dark:bg-amber-950/30"
            data-testid="pending-group"
          >
            <div class="mb-2 text-sm font-semibold text-amber-800 dark:text-amber-300">
              {{
                group.groupName().value() ||
                  t('common.unnamedGroup', { id: group.groupId().value() })
              }}
            </div>
            <div class="mb-1 text-xs text-gray-600 dark:text-gray-400">
              {{ t('user.workoutStart.pickExercise') }}
            </div>
            <select hlmInput class="w-full" (change)="onGroupExercisePicked(gi, $event)">
              <option value="">{{ t('common.select') }}</option>
              @for (member of group.members().value(); track member.id) {
                <option [value]="member.id">{{ member.name }}</option>
              }
            </select>
          </div>
        }

        <!-- Add Exercise button -->
        @if (!readonly()) {
          <button
            type="button"
            (click)="addExerciseRequested.emit()"
            class="mt-2 text-sm text-blue-500/70 hover:text-blue-600 dark:text-blue-400/70 dark:hover:text-blue-300"
          >
            {{ t('user.workouts.addExercise') }}
          </button>
        }
      </div>
    </ng-container>
  `,
})
export class WorkoutStartSection {
  readonly SECTION_TYPE_MAIN = WorkoutSectionTypeMain;
  readonly SECTION_TYPE_SUPPLEMENTARY = WorkoutSectionTypeSupplementary;

  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  section = input.required<any>();
  sectionIndex = input.required<number>();
  exerciseDisplayMap = input.required<Record<number, ExerciseDisplayInfo>>();
  readonly = input(false);

  removed = output<void>();
  sectionChanged = output<void>();
  exerciseRemoved = output<{ exerciseIndex: number }>();
  exerciseChanged = output<{ exerciseIndex: number }>();
  setChanged = output<{ exerciseIndex: number; setIndex: number }>();
  exerciseDropped = output<{ previousIndex: number; currentIndex: number }>();
  addExerciseRequested = output<void>();
  groupExercisePicked = output<{ groupIndex: number; exerciseId: number }>();

  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  onExerciseDrop(event: CdkDragDrop<any>) {
    if (event.previousIndex === event.currentIndex) return;
    this.exerciseDropped.emit({
      previousIndex: event.previousIndex,
      currentIndex: event.currentIndex,
    });
  }

  onGroupExercisePicked(groupIndex: number, event: Event) {
    const exerciseId = Number((event.target as HTMLSelectElement).value);
    if (!exerciseId || isNaN(exerciseId)) return;
    this.groupExercisePicked.emit({ groupIndex, exerciseId });
  }
}
