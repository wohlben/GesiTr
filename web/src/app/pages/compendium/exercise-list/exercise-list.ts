import { Component, inject, computed } from '@angular/core';
import { ActivatedRoute, Router, RouterLink } from '@angular/router';
import { toSignal } from '@angular/core/rxjs-interop';
import {
  injectQuery,
  injectMutation,
  keepPreviousData,
  QueryClient,
} from '@tanstack/angular-query-experimental';
import { CompendiumApiClient } from '$core/api-clients/compendium-api-client';
import { UserApiClient } from '$core/api-clients/user-api-client';
import { exerciseKeys, masteryKeys, namePreferenceKeys } from '$core/query-keys';
import { Exercise } from '$generated/models';
import { ExerciseMastery } from '$generated/user-mastery';
import { TranslocoDirective } from '@jsverse/transloco';
import { ExerciseListItem } from '$ui/compendium/exercise-list-item/exercise-list-item';
import { LocalityToggleGroup } from '$ui/compendium/locality-toggle-group/locality-toggle-group';
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
  imports: [
    ExerciseListItem,
    LocalityToggleGroup,
    DataTable,
    Pagination,
    PageLayout,
    RouterLink,
    TranslocoDirective,
  ],
  template: `
    <ng-container *transloco="let t">
      <app-page-layout
        [header]="t('compendium.exercises.title')"
        [isPending]="exercisesQuery.isPending()"
        [errorMessage]="exercisesQuery.isError() ? exercisesQuery.error().message : undefined"
      >
        <div actions class="flex items-center gap-2">
          <app-locality-toggle-group (selectedChange)="onLocalityChange($event)" />
          <div class="flex-grow"></div>
          <a
            routerLink="./new"
            class="rounded-md bg-blue-600 px-4 py-2 text-sm font-medium text-white hover:bg-blue-700"
            >{{ t('common.new') }}</a
          >
        </div>
        @if (exercisesQuery.data(); as page) {
          <app-data-table
            [columns]="exerciseColumns"
            [stale]="exercisesQuery.isPlaceholderData()"
            [initialHiddenColumns]="savedHiddenColumns"
            (hiddenColumnsChange)="onHiddenColumnsChange($event)"
          >
            @for (ex of page.items; track ex.id) {
              <tr
                app-exercise-list-item
                [exercise]="ex"
                [displayName]="getDisplayName(ex)"
                [matchingNames]="getMatchingNames(ex)"
                [mastery]="masteryMap().get(ex.id)"
                (nameClicked)="onNameClicked(ex.id, $event)"
              ></tr>
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
  private userApi = inject(UserApiClient);
  private queryClient = inject(QueryClient);
  private router = inject(Router);
  private route = inject(ActivatedRoute);
  private queryParams = toSignal(this.route.queryParamMap);

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

  masteryQuery = injectQuery(() => ({
    queryKey: masteryKeys.list(),
    queryFn: () => this.userApi.fetchMasteryList(),
  }));

  masteryMap = computed(() => {
    const map = new Map<number, ExerciseMastery>();
    for (const m of this.masteryQuery.data() ?? []) {
      map.set(m.exerciseId, m);
    }
    return map;
  });

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

  namePreferenceQuery = injectQuery(() => ({
    queryKey: namePreferenceKeys.list(),
    queryFn: () => this.userApi.fetchExerciseNamePreferences(),
  }));

  // Maps exerciseId → preferred exercise_name.id
  preferenceMap = computed(() => {
    const map = new Map<number, number>();
    for (const p of this.namePreferenceQuery.data() ?? []) {
      map.set(p.exerciseId, p.exerciseNameId);
    }
    return map;
  });

  searchQuery = computed(() => this.filters()['q'] ?? '');

  savePreferenceMutation = injectMutation(() => ({
    mutationFn: (vars: { exerciseId: number; exerciseNameId: number }) =>
      this.userApi.setExerciseNamePreference(vars.exerciseId, vars.exerciseNameId),
    onSuccess: () => this.queryClient.invalidateQueries({ queryKey: namePreferenceKeys.all() }),
  }));

  getDisplayName(ex: Exercise): string {
    const q = this.searchQuery().toLowerCase();
    if (q) {
      const names = ex.names ?? [];
      const startsWithMatch = names.find((n) => n.name.toLowerCase().startsWith(q));
      if (startsWithMatch) return startsWithMatch.name;
      const containsMatch = names.find((n) => n.name.toLowerCase().includes(q));
      if (containsMatch) return containsMatch.name;
    }
    const prefId = this.preferenceMap().get(ex.id);
    if (prefId) {
      const preferred = ex.names?.find((n) => n.id === prefId);
      if (preferred) return preferred.name;
    }
    return ex.names?.[0]?.name ?? '';
  }

  getMatchingNames(ex: Exercise): string[] {
    const q = this.searchQuery().toLowerCase();
    if (!q) return [];
    const primary = this.getDisplayName(ex);
    return (ex.names ?? [])
      .filter((n) => n.name !== primary && n.name.toLowerCase().includes(q))
      .map((n) => n.name);
  }

  onNameClicked(exerciseId: number, name: string) {
    const ex = this.exercisesQuery.data()?.items.find((e) => e.id === exerciseId);
    const nameEntry = ex?.names?.find((n) => n.name === name);
    if (nameEntry) {
      this.savePreferenceMutation.mutate({ exerciseId, exerciseNameId: nameEntry.id });
    }
  }

  onLocalityChange(localityId: number | null) {
    this.router.navigate([], {
      relativeTo: this.route,
      queryParams: { localityId: localityId ?? null, offset: null },
      queryParamsHandling: 'merge',
    });
  }

  exerciseColumns: DataTableColumn[] = [
    { label: 'Name', labelKey: 'fields.name', searchParam: 'q' },
    { label: 'Mastery', labelKey: 'fields.mastery' },
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
    { label: 'Body weight scaling', labelKey: 'fields.bodyWeightScaling', defaultHidden: true },
    {
      label: 'Measurement paradigms',
      labelKey: 'fields.measurementParadigms',
      defaultHidden: true,
    },
    { label: 'Description', labelKey: 'fields.description', defaultHidden: true },
    { label: 'Names', labelKey: 'fields.names', defaultHidden: true },
    { label: 'Author', labelKey: 'fields.author', defaultHidden: true },
    { label: 'Version', labelKey: 'fields.version', defaultHidden: true },
    { label: 'Owner', labelKey: 'fields.owner', defaultHidden: true },
    { label: 'Created at', labelKey: 'fields.createdAt', defaultHidden: true },
    { label: 'Updated at', labelKey: 'fields.updatedAt', defaultHidden: true },
  ];
}
