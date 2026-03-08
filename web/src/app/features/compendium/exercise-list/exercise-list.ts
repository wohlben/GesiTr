import { Component, inject, computed } from '@angular/core';
import { ActivatedRoute, RouterLink } from '@angular/router';
import { toSignal } from '@angular/core/rxjs-interop';
import { injectQuery, keepPreviousData } from '@tanstack/angular-query-experimental';
import { CompendiumApiClient } from '$core/api-clients/compendium-api-client';
import { exerciseKeys } from '$core/query-keys';
import { ExerciseListItem } from '$ui/compendium/exercise-list-item/exercise-list-item';
import { DataTable, DataTableColumn } from '$ui/data-table/data-table';
import { Pagination } from '$ui/pagination/pagination';
import { PageLayout } from '../../../layout/page-layout';
import {
  ExerciseTypeStrength,
  ExerciseTypeCardio,
  ExerciseTypeStretching,
  ExerciseTypeStrongman,
  DifficultyBeginner,
  DifficultyIntermediate,
  DifficultyAdvanced,
  ForcePull,
  ForcePush,
  ForceStatic,
  ForceDynamic,
  ForceHinge,
  ForceRotation,
  MuscleAbs,
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
  MuscleNeck,
  MuscleObliques,
  MuscleQuads,
  MuscleTraps,
  MuscleTriceps,
  MuscleFrontDelts,
  MuscleRearDelts,
  MuscleRhomboids,
  MuscleSideDelts,
} from '$generated/models';

@Component({
  selector: 'app-exercise-list',
  imports: [ExerciseListItem, DataTable, Pagination, PageLayout, RouterLink],
  template: `
    <app-page-layout
      header="Exercises"
      [isPending]="exercisesQuery.isPending()"
      [errorMessage]="exercisesQuery.isError() ? exercisesQuery.error().message : undefined"
    >
      <a
        actions
        routerLink="./new"
        class="rounded-md bg-blue-600 px-4 py-2 text-sm font-medium text-white hover:bg-blue-700"
        >New</a
      >
      @if (exercisesQuery.data(); as page) {
        <app-data-table
          [columns]="exerciseColumns"
          [stale]="exercisesQuery.isPlaceholderData()"
          [initialHiddenColumns]="savedHiddenColumns"
          (hiddenColumnsChange)="onHiddenColumnsChange($event)"
        >
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
    queryKey: exerciseKeys.list(this.filters()),
    queryFn: () => this.api.fetchExercises(this.filters()),
    placeholderData: keepPreviousData,
  }));

  private static readonly STORAGE_KEY = 'dt-columns-exercises';

  savedHiddenColumns = ExerciseList.loadHiddenColumns();

  onHiddenColumnsChange(labels: string[]) {
    localStorage.setItem(ExerciseList.STORAGE_KEY, JSON.stringify(labels));
  }

  private static loadHiddenColumns(): string[] | undefined {
    try {
      const stored = localStorage.getItem(ExerciseList.STORAGE_KEY);
      return stored ? JSON.parse(stored) : undefined;
    } catch {
      return undefined;
    }
  }

  exerciseColumns: DataTableColumn[] = [
    { label: 'Name', searchParam: 'q' },
    {
      label: 'Type',
      filterParam: 'type',
      options: [
        ExerciseTypeStrength,
        ExerciseTypeCardio,
        ExerciseTypeStretching,
        ExerciseTypeStrongman,
      ],
    },
    {
      label: 'Difficulty',
      filterParam: 'difficulty',
      options: [DifficultyBeginner, DifficultyIntermediate, DifficultyAdvanced],
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
        MuscleNeck,
        MuscleObliques,
        MuscleQuads,
        MuscleTraps,
        MuscleTriceps,
        MuscleFrontDelts,
        MuscleRearDelts,
        MuscleRhomboids,
        MuscleSideDelts,
      ],
    },
    { label: 'Secondary muscles', defaultHidden: true },
    { label: 'Slug', defaultHidden: true },
    { label: 'Body weight scaling', defaultHidden: true },
    { label: 'Measurement paradigms', defaultHidden: true },
    { label: 'Description', defaultHidden: true },
    { label: 'Alternative names', defaultHidden: true },
    { label: 'Author', defaultHidden: true },
    { label: 'Version', defaultHidden: true },
    { label: 'Created by', defaultHidden: true },
    { label: 'Created at', defaultHidden: true },
    { label: 'Updated at', defaultHidden: true },
  ];
}
