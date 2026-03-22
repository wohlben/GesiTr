import { Component, inject, computed } from '@angular/core';
import { ActivatedRoute, RouterLink } from '@angular/router';
import { toSignal } from '@angular/core/rxjs-interop';
import { injectQuery, keepPreviousData } from '@tanstack/angular-query-experimental';
import { CompendiumApiClient } from '$core/api-clients/compendium-api-client';
import { exerciseKeys } from '$core/query-keys';
import { TranslocoDirective } from '@jsverse/transloco';
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
  imports: [ExerciseListItem, DataTable, Pagination, PageLayout, RouterLink, TranslocoDirective],
  template: `
    <ng-container *transloco="let t">
      <app-page-layout
        [header]="t('compendium.exercises.title')"
        [isPending]="exercisesQuery.isPending()"
        [errorMessage]="exercisesQuery.isError() ? exercisesQuery.error().message : undefined"
      >
        <a
          actions
          routerLink="./new"
          class="rounded-md bg-blue-600 px-4 py-2 text-sm font-medium text-white hover:bg-blue-700"
          >{{ t('common.new') }}</a
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
          <app-pagination [page]="page" [emptyLabel]="t('compendium.exercises.noResults')" />
        }
      </app-page-layout>
    </ng-container>
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
    { label: 'Name', labelKey: 'fields.name', searchParam: 'q' },
    {
      label: 'Type',
      labelKey: 'fields.type',
      filterParam: 'type',
      optionKeyPrefix: 'enums.exerciseType',
      options: [
        ExerciseTypeStrength,
        ExerciseTypeCardio,
        ExerciseTypeStretching,
        ExerciseTypeStrongman,
      ],
    },
    {
      label: 'Difficulty',
      labelKey: 'fields.difficulty',
      filterParam: 'difficulty',
      optionKeyPrefix: 'enums.difficulty',
      options: [DifficultyBeginner, DifficultyIntermediate, DifficultyAdvanced],
    },
    {
      label: 'Force',
      labelKey: 'fields.force',
      filterParam: 'force',
      optionKeyPrefix: 'enums.force',
      options: [ForcePull, ForcePush, ForceStatic, ForceDynamic, ForceHinge, ForceRotation],
    },
    {
      label: 'Primary muscles',
      labelKey: 'fields.primaryMuscles',
      filterParam: 'muscle',
      optionKeyPrefix: 'enums.muscle',
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
    { label: 'Secondary muscles', labelKey: 'fields.secondaryMuscles', defaultHidden: true },
    { label: 'Slug', labelKey: 'fields.slug', defaultHidden: true },
    { label: 'Body weight scaling', labelKey: 'fields.bodyWeightScaling', defaultHidden: true },
    {
      label: 'Measurement paradigms',
      labelKey: 'fields.measurementParadigms',
      defaultHidden: true,
    },
    { label: 'Description', labelKey: 'fields.description', defaultHidden: true },
    { label: 'Alternative names', labelKey: 'fields.alternativeNames', defaultHidden: true },
    { label: 'Author', labelKey: 'fields.author', defaultHidden: true },
    { label: 'Version', labelKey: 'fields.version', defaultHidden: true },
    { label: 'Created by', labelKey: 'fields.createdBy', defaultHidden: true },
    { label: 'Created at', labelKey: 'fields.createdAt', defaultHidden: true },
    { label: 'Updated at', labelKey: 'fields.updatedAt', defaultHidden: true },
  ];
}
