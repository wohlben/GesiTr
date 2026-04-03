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
        <!-- Home locality entry (always visible) -->
        <button
          type="button"
          (click)="goHome()"
          [disabled]="homeMutation.isPending()"
          class="mb-6 flex w-full items-center gap-3 rounded-lg border-2 border-dashed border-gray-300 p-4 text-left transition-colors hover:border-blue-400 hover:bg-blue-50 disabled:opacity-50 dark:border-gray-600 dark:hover:border-blue-500 dark:hover:bg-blue-900/20"
        >
          <span
            class="flex h-10 w-10 shrink-0 items-center justify-center rounded-full bg-blue-100 text-blue-600 dark:bg-blue-900/30 dark:text-blue-400"
          >
            <svg class="h-5 w-5" viewBox="0 0 20 20" fill="currentColor">
              <path
                d="M10.707 2.293a1 1 0 00-1.414 0l-7 7a1 1 0 001.414 1.414L4 10.414V17a1 1 0 001 1h2a1 1 0 001-1v-2a1 1 0 011-1h2a1 1 0 011 1v2a1 1 0 001 1h2a1 1 0 001-1v-6.586l.293.293a1 1 0 001.414-1.414l-7-7z"
              />
            </svg>
          </span>
          <div>
            <span class="text-sm font-medium text-gray-900 dark:text-gray-100">
              @if (homeMutation.isPending()) {
                {{ t('common.starting') }}
              } @else {
                {{ t('compendium.localities.home') }}
              }
            </span>
            <p class="text-xs text-gray-500 dark:text-gray-400">
              {{ t('compendium.localities.homeDescription') }}
            </p>
          </div>
        </button>
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
  private router = inject(Router);
  private queryClient = inject(QueryClient);
  private slugify = new SlugifyPipe();
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

  private homeLocalityQuery = injectQuery(() => ({
    queryKey: localityKeys.list({ owner: 'me', public: 'false', limit: 1 }),
    queryFn: () => this.api.fetchLocalities({ owner: 'me', public: 'false', limit: 1 }),
  }));

  homeMutation = injectMutation(() => ({
    mutationFn: () => this.api.createLocality({ name: 'Home', public: false }),
    onSuccess: (result) => {
      this.queryClient.invalidateQueries({ queryKey: localityKeys.all() });
      this.router.navigate([
        '/compendium/localities',
        result.id,
        this.slugify.transform(result.name),
      ]);
    },
  }));

  goHome() {
    const home = this.homeLocalityQuery.data()?.items?.[0];
    if (home) {
      this.router.navigate(['/compendium/localities', home.id, this.slugify.transform(home.name)]);
    } else {
      this.homeMutation.mutate();
    }
  }

  localityColumns: DataTableColumn[] = [
    { label: 'Name', labelKey: 'fields.name', searchParam: 'q' },
    { label: 'Owner', labelKey: 'fields.owner', defaultHidden: true },
    { label: 'Created at', labelKey: 'fields.createdAt', defaultHidden: true },
  ];
}
