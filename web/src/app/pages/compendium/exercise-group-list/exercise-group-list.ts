import { Component, inject, computed } from '@angular/core';
import { ActivatedRoute, RouterLink } from '@angular/router';
import { toSignal } from '@angular/core/rxjs-interop';
import { injectQuery, keepPreviousData } from '@tanstack/angular-query-experimental';
import { CompendiumApiClient } from '$core/api-clients/compendium-api-client';
import { exerciseGroupKeys } from '$core/query-keys';
import { TranslocoDirective } from '@jsverse/transloco';
import { ExerciseGroupListItem } from '$ui/compendium/exercise-group-list-item/exercise-group-list-item';
import { DataTable, DataTableColumn } from '$ui/data-table/data-table';
import { Pagination } from '$ui/pagination/pagination';
import { PageLayout } from '../../../layout/page-layout';

@Component({
  selector: 'app-exercise-group-list',
  imports: [
    ExerciseGroupListItem,
    DataTable,
    Pagination,
    PageLayout,
    RouterLink,
    TranslocoDirective,
  ],
  template: `
    <ng-container *transloco="let t">
      <app-page-layout
        [header]="t('compendium.exerciseGroups.title')"
        [isPending]="groupsQuery.isPending()"
        [errorMessage]="groupsQuery.isError() ? groupsQuery.error().message : undefined"
      >
        <a
          actions
          routerLink="./new"
          class="rounded-md bg-blue-600 px-4 py-2 text-sm font-medium text-white hover:bg-blue-700"
          >{{ t('common.new') }}</a
        >
        @if (groupsQuery.data(); as page) {
          <app-data-table
            [columns]="groupColumns"
            [stale]="groupsQuery.isPlaceholderData()"
            [initialHiddenColumns]="savedHiddenColumns"
            (hiddenColumnsChange)="onHiddenColumnsChange($event)"
          >
            @for (group of page.items; track group.id) {
              <tr app-exercise-group-list-item [group]="group"></tr>
            }
          </app-data-table>
          <app-pagination [page]="page" [emptyLabel]="t('compendium.exerciseGroups.noResults')" />
        }
      </app-page-layout>
    </ng-container>
  `,
})
export class ExerciseGroupList {
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

  groupsQuery = injectQuery(() => ({
    queryKey: exerciseGroupKeys.list(this.filters()),
    queryFn: () => this.api.fetchExerciseGroups(this.filters()),
    placeholderData: keepPreviousData,
  }));

  private static readonly STORAGE_KEY = 'dt-columns-exercise-groups';

  savedHiddenColumns = ExerciseGroupList.loadHiddenColumns();

  onHiddenColumnsChange(labels: string[]) {
    localStorage.setItem(ExerciseGroupList.STORAGE_KEY, JSON.stringify(labels));
  }

  private static loadHiddenColumns(): string[] | undefined {
    try {
      const stored = localStorage.getItem(ExerciseGroupList.STORAGE_KEY);
      return stored ? JSON.parse(stored) : undefined;
    } catch {
      return undefined;
    }
  }

  groupColumns: DataTableColumn[] = [
    { label: 'Name', labelKey: 'fields.name', searchParam: 'q' },
    { label: 'Description', labelKey: 'fields.description' },
    { label: 'Created by', labelKey: 'fields.createdBy', defaultHidden: true },
    { label: 'Created at', labelKey: 'fields.createdAt', defaultHidden: true },
    { label: 'Updated at', labelKey: 'fields.updatedAt', defaultHidden: true },
  ];
}
