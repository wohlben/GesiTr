import { Component, inject, signal, computed, effect, untracked } from '@angular/core';
import { injectQuery } from '@tanstack/angular-query-experimental';
import { CompendiumApiClient } from '$core/api-clients/compendium-api-client';
import { ExerciseListItem } from '$ui/compendium/exercise-list-item/exercise-list-item';
import { SearchInput } from '$ui/inputs/search-input/search-input';
import { FilterSelect } from '$ui/inputs/filter-select/filter-select';
import { DataTable } from '$ui/data-table/data-table';
import { PageLayout } from '../../../layout/page-layout';
import {
  ExerciseType,
  ExerciseTypeStrength,
  ExerciseTypeCardio,
  ExerciseTypeStretching,
  ExerciseTypePlyometric,
  ExerciseTypeStrongman,
  ExerciseTypeOlympicWeightlifting,
  ExerciseTypePowerlifting,
  TechnicalDifficulty,
  DifficultyBeginner,
  DifficultyIntermediate,
  DifficultyAdvanced,
  DifficultyExpert,
  Force,
  ForcePull,
  ForcePush,
  ForceStatic,
  ForceDynamic,
  ForceHinge,
  ForceRotation,
  Muscle,
  MuscleAbs,
  MuscleAbductors,
  MuscleAdductors,
  MuscleBiceps,
  MuscleCalves,
  MuscleChest,
  MuscleForearms,
  MuscleGlutes,
  MuscleHamstrings,
  MuscleHipFlexors,
  MuscleLats,
  MuscleLowerBack,
  MuscleMiddleBack,
  MuscleNeck,
  MuscleObliques,
  MuscleQuads,
  MuscleShoulders,
  MuscleTraps,
  MuscleTriceps,
  MuscleFrontDelts,
  MuscleRearDelts,
  MuscleRhomboids,
  MuscleSideDelts,
} from '$generated/models';

@Component({
  selector: 'app-exercise-list',
  imports: [SearchInput, FilterSelect, ExerciseListItem, DataTable, PageLayout],
  template: `
    <app-page-layout
      header="Exercises"
      [isPending]="exercisesQuery.isPending()"
      [errorMessage]="exercisesQuery.isError() ? exercisesQuery.error().message : undefined">

      <div filters class="flex flex-wrap gap-3">
        <app-search-input placeholder="Search exercises..." [(value)]="q" />
        <app-filter-select allLabel="All types" [options]="typeOptions" [(value)]="type" />
        <app-filter-select allLabel="All difficulties" [options]="difficultyOptions" [(value)]="difficulty" />
        <app-filter-select allLabel="All forces" [options]="forceOptions" [(value)]="force" />
        <app-filter-select allLabel="All muscles" [options]="muscleOptions" [(value)]="muscle" />
      </div>

      @if (exercisesQuery.data(); as page) {
        <app-data-table [columns]="['Name', 'Type', 'Difficulty', 'Force', 'Primary muscles']">
          @for (ex of page.items; track ex.id) {
            <tr app-exercise-list-item [exercise]="ex"></tr>
          }
        </app-data-table>
        <div class="flex items-center justify-between text-sm text-gray-500 dark:text-gray-400">
          <p>{{ page.total === 0 ? 'No exercises found' : 'Showing ' + (page.offset + 1) + '–' + (page.offset + page.items.length) + ' of ' + page.total + ' exercises' }}</p>
          <div class="flex gap-2">
            <button
              class="rounded border border-gray-300 px-3 py-1 disabled:opacity-50 dark:border-gray-600"
              [disabled]="page.offset === 0"
              (click)="prevPage()">Previous</button>
            <button
              class="rounded border border-gray-300 px-3 py-1 disabled:opacity-50 dark:border-gray-600"
              [disabled]="page.offset + page.limit >= page.total"
              (click)="nextPage()">Next</button>
          </div>
        </div>
      }
    </app-page-layout>
  `,
})
export class ExerciseList {
  private api = inject(CompendiumApiClient);

  q = signal('');
  type = signal<ExerciseType | ''>('');
  difficulty = signal<TechnicalDifficulty | ''>('');
  force = signal<Force | ''>('');
  muscle = signal<Muscle | ''>('');
  offset = signal(0);

  private resetOffset = effect(() => {
    this.q(); this.type(); this.difficulty(); this.force(); this.muscle();
    untracked(() => this.offset.set(0));
  });

  filters = computed(() => ({
    q: this.q() || undefined,
    type: this.type() || undefined,
    difficulty: this.difficulty() || undefined,
    force: this.force() || undefined,
    muscle: this.muscle() || undefined,
    offset: this.offset() || undefined,
  }));

  exercisesQuery = injectQuery(() => ({
    queryKey: ['exercises', this.filters()],
    queryFn: () => this.api.fetchExercises(this.filters()),
  }));

  prevPage() {
    const page = this.exercisesQuery.data();
    if (page) this.offset.set(Math.max(0, page.offset - page.limit));
  }

  nextPage() {
    const page = this.exercisesQuery.data();
    if (page) this.offset.set(page.offset + page.limit);
  }

  typeOptions: ExerciseType[] = [
    ExerciseTypeStrength,
    ExerciseTypeCardio,
    ExerciseTypeStretching,
    ExerciseTypePlyometric,
    ExerciseTypeStrongman,
    ExerciseTypeOlympicWeightlifting,
    ExerciseTypePowerlifting,
  ];

  difficultyOptions: TechnicalDifficulty[] = [
    DifficultyBeginner,
    DifficultyIntermediate,
    DifficultyAdvanced,
    DifficultyExpert,
  ];

  forceOptions: Force[] = [
    ForcePull,
    ForcePush,
    ForceStatic,
    ForceDynamic,
    ForceHinge,
    ForceRotation,
  ];

  muscleOptions: Muscle[] = [
    MuscleAbs,
    MuscleAbductors,
    MuscleAdductors,
    MuscleBiceps,
    MuscleCalves,
    MuscleChest,
    MuscleForearms,
    MuscleGlutes,
    MuscleHamstrings,
    MuscleHipFlexors,
    MuscleLats,
    MuscleLowerBack,
    MuscleMiddleBack,
    MuscleNeck,
    MuscleObliques,
    MuscleQuads,
    MuscleShoulders,
    MuscleTraps,
    MuscleTriceps,
    MuscleFrontDelts,
    MuscleRearDelts,
    MuscleRhomboids,
    MuscleSideDelts,
  ];
}
