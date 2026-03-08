import { Component, inject, computed } from '@angular/core';
import { ActivatedRoute } from '@angular/router';
import { toSignal } from '@angular/core/rxjs-interop';
import { injectQuery, keepPreviousData } from '@tanstack/angular-query-experimental';
import { CompendiumApiClient } from '$core/api-clients/compendium-api-client';
import { ExerciseListItem } from '$ui/compendium/exercise-list-item/exercise-list-item';
import { DataTable, DataTableColumn } from '$ui/data-table/data-table';
import { Pagination } from '$ui/pagination/pagination';
import { PageLayout } from '../../../layout/page-layout';
import {
  ExerciseTypeStrength,
  ExerciseTypeCardio,
  ExerciseTypeStretching,
  ExerciseTypePlyometric,
  ExerciseTypeStrongman,
  ExerciseTypeOlympicWeightlifting,
  ExerciseTypePowerlifting,
  DifficultyBeginner,
  DifficultyIntermediate,
  DifficultyAdvanced,
  DifficultyExpert,
  ForcePull,
  ForcePush,
  ForceStatic,
  ForceDynamic,
  ForceHinge,
  ForceRotation,
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
  imports: [ExerciseListItem, DataTable, Pagination, PageLayout],
  template: `
    <app-page-layout
      header="Exercises"
      [isPending]="exercisesQuery.isPending()"
      [errorMessage]="exercisesQuery.isError() ? exercisesQuery.error().message : undefined">

      @if (exercisesQuery.data(); as page) {
        <app-data-table [columns]="exerciseColumns" [stale]="exercisesQuery.isPlaceholderData()">
          @for (ex of page.items; track ex.id) {
            <tr app-exercise-list-item [exercise]="ex"></tr>
          }
        </app-data-table>
        <app-pagination [page]="page" emptyLabel="No exercises found" />
      }
    </app-page-layout>
  `,
})
export class ExerciseList {
  private api = inject(CompendiumApiClient);
  private queryParams = toSignal(inject(ActivatedRoute).queryParamMap);

  filters = computed(() => {
    const params: Record<string, string> = {};
    const qp = this.queryParams();
    if (qp) {
      for (const key of qp.keys) {
        const val = qp.get(key);
        if (val) params[key] = val;
      }
    }
    return params;
  });

  exercisesQuery = injectQuery(() => ({
    queryKey: ['exercises', this.filters()],
    queryFn: () => this.api.fetchExercises(this.filters()),
    placeholderData: keepPreviousData,
  }));

  exerciseColumns: DataTableColumn[] = [
    { label: 'Name', searchParam: 'q' },
    {
      label: 'Type',
      filterParam: 'type',
      options: [
        ExerciseTypeStrength,
        ExerciseTypeCardio,
        ExerciseTypeStretching,
        ExerciseTypePlyometric,
        ExerciseTypeStrongman,
        ExerciseTypeOlympicWeightlifting,
        ExerciseTypePowerlifting,
      ],
    },
    {
      label: 'Difficulty',
      filterParam: 'difficulty',
      options: [DifficultyBeginner, DifficultyIntermediate, DifficultyAdvanced, DifficultyExpert],
    },
    {
      label: 'Force',
      filterParam: 'force',
      options: [ForcePull, ForcePush, ForceStatic, ForceDynamic, ForceHinge, ForceRotation],
    },
    {
      label: 'Primary muscles',
      filterParam: 'muscle',
      options: [
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
      ],
    },
  ];
}
