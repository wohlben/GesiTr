import { Component, inject, computed } from '@angular/core';
import { ActivatedRoute, RouterLink } from '@angular/router';
import { toSignal } from '@angular/core/rxjs-interop';
import { injectQuery, keepPreviousData } from '@tanstack/angular-query-experimental';
import { CompendiumApiClient } from '$core/api-clients/compendium-api-client';
import { localityKeys } from '$core/query-keys';
import { TranslocoDirective } from '@jsverse/transloco';
import { DataTable, DataTableColumn } from '$ui/data-table/data-table';
import { Pagination } from '$ui/pagination/pagination';
import { PageLayout } from '../../../layout/page-layout';
import { SlugifyPipe } from '$ui/pipes/slugify';
import { DatePipe } from '@angular/common';

@Component({
  selector: 'app-locality-list',
  imports: [
    DataTable,
    Pagination,
    PageLayout,
    RouterLink,
    TranslocoDirective,
    SlugifyPipe,
    DatePipe,
  ],
  template: `
    <ng-container *transloco="let t">
      <app-page-layout
        [header]="t('compendium.localities.title')"
        [isPending]="localityQuery.isPending()"
        [errorMessage]="localityQuery.isError() ? localityQuery.error().message : undefined"
      >
        <a
          actions
          routerLink="./new"
          class="rounded-md bg-blue-600 px-4 py-2 text-sm font-medium text-white hover:bg-blue-700"
          >{{ t('common.new') }}</a
        >
        @if (localityQuery.data(); as page) {
          <app-data-table [columns]="localityColumns" [stale]="localityQuery.isPlaceholderData()">
            @for (locality of page.items; track locality.id) {
              <tr
                [routerLink]="['./', locality.id, locality.name | slugify]"
                class="cursor-pointer hover:bg-gray-50 dark:hover:bg-gray-800"
              >
                <td class="px-4 py-3 text-sm text-gray-900 dark:text-gray-100">
                  {{ locality.name }}
                </td>
                <td class="px-4 py-3 text-sm text-gray-500 dark:text-gray-400">
                  {{ locality.owner }}
                </td>
                <td class="px-4 py-3 text-sm text-gray-500 dark:text-gray-400">
                  {{ locality.createdAt | date: 'mediumDate' }}
                </td>
              </tr>
            }
          </app-data-table>
          <app-pagination [page]="page" [emptyLabel]="t('compendium.localities.noResults')" />
        }
      </app-page-layout>
    </ng-container>
  `,
})
export class LocalityList {
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

  localityQuery = injectQuery(() => ({
    queryKey: localityKeys.list(this.filters()),
    queryFn: () => this.api.fetchLocalities(this.filters()),
    placeholderData: keepPreviousData,
  }));

  localityColumns: DataTableColumn[] = [
    { label: 'Name', labelKey: 'fields.name', searchParam: 'q' },
    { label: 'Owner', labelKey: 'fields.owner', defaultHidden: true },
    { label: 'Created at', labelKey: 'fields.createdAt', defaultHidden: true },
  ];
}
