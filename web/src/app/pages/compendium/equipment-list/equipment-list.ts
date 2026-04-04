import { Component, inject, computed, signal } from '@angular/core';
import { ActivatedRoute, RouterLink } from '@angular/router';
import { toSignal } from '@angular/core/rxjs-interop';
import { injectQuery, keepPreviousData } from '@tanstack/angular-query-experimental';
import { CompendiumApiClient } from '$core/api-clients/compendium-api-client';
import { UserApiClient } from '$core/api-clients/user-api-client';
import { equipmentKeys, equipmentMasteryKeys, localityAvailabilityKeys } from '$core/query-keys';
import { EquipmentMastery } from '$generated/user-mastery';
import { TranslocoDirective } from '@jsverse/transloco';
import { EquipmentListItem } from '$ui/compendium/equipment-list-item/equipment-list-item';
import { LocalityToggleGroup } from '$ui/compendium/locality-toggle-group/locality-toggle-group';
import { DataTable, DataTableColumn } from '$ui/data-table/data-table';
import { Pagination } from '$ui/pagination/pagination';
import { PageLayout } from '../../../layout/page-layout';
import {
  EquipmentCategoryFreeWeights,
  EquipmentCategoryAccessories,
  EquipmentCategoryBenches,
  EquipmentCategoryMachines,
  EquipmentCategoryFunctional,
  EquipmentCategoryOther,
} from '$generated/models';

@Component({
  selector: 'app-equipment-list',
  imports: [
    EquipmentListItem,
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
        [header]="t('compendium.equipment.title')"
        [isPending]="equipmentQuery.isPending()"
        [errorMessage]="equipmentQuery.isError() ? equipmentQuery.error().message : undefined"
      >
        <div actions class="flex items-center gap-2">
          <app-locality-toggle-group (selectedChange)="selectedLocalityId.set($event)" />
          <div class="flex-grow"></div>
          <a
            routerLink="/compendium/localities"
            class="rounded-md border border-gray-300 px-4 py-2 text-sm font-medium text-gray-700 hover:bg-gray-50 dark:border-gray-600 dark:text-gray-300 dark:hover:bg-gray-800"
            >{{ t('nav.localities') }}</a
          >
          <a
            routerLink="./new"
            class="rounded-md bg-blue-600 px-4 py-2 text-sm font-medium text-white hover:bg-blue-700"
            >{{ t('common.new') }}</a
          >
        </div>
        @if (filteredPage(); as page) {
          <app-data-table
            [columns]="equipmentColumns"
            [stale]="equipmentQuery.isPlaceholderData()"
            [initialHiddenColumns]="savedHiddenColumns"
            (hiddenColumnsChange)="onHiddenColumnsChange($event)"
          >
            @for (item of page.items; track item.id) {
              <tr
                app-equipment-list-item
                [equipment]="item"
                [mastery]="masteryMap().get(item.id)"
              ></tr>
            }
          </app-data-table>
          <app-pagination [page]="page" [emptyLabel]="t('compendium.equipment.noResults')" />
        }
      </app-page-layout>
    </ng-container>
  `,
})
export class EquipmentList {
  private api = inject(CompendiumApiClient);
  private userApi = inject(UserApiClient);
  private queryParams = toSignal(inject(ActivatedRoute).queryParamMap);

  selectedLocalityId = signal<number | null>(null);

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

  equipmentQuery = injectQuery(() => ({
    queryKey: equipmentKeys.list(this.filters()),
    queryFn: () => this.api.fetchEquipment(this.filters()),
    placeholderData: keepPreviousData,
  }));

  masteryQuery = injectQuery(() => ({
    queryKey: equipmentMasteryKeys.list(),
    queryFn: () => this.userApi.fetchEquipmentMasteryList(),
  }));

  availabilitiesQuery = injectQuery(() => ({
    queryKey: localityAvailabilityKeys.list({ localityId: this.selectedLocalityId()! }),
    queryFn: () => this.api.fetchLocalityAvailabilities({ localityId: this.selectedLocalityId()! }),
    enabled: this.selectedLocalityId() !== null,
  }));

  private availableEquipmentIds = computed(() => {
    const ids = new Set<number>();
    for (const a of this.availabilitiesQuery.data() ?? []) {
      ids.add(a.equipmentId);
    }
    return ids;
  });

  filteredPage = computed(() => {
    const page = this.equipmentQuery.data();
    if (!page) return undefined;
    if (this.selectedLocalityId() === null) return page;
    const ids = this.availableEquipmentIds();
    const items = page.items.filter((e) => ids.has(e.id));
    return { ...page, items, total: items.length };
  });

  masteryMap = computed(() => {
    const map = new Map<number, EquipmentMastery>();
    for (const m of this.masteryQuery.data() ?? []) {
      map.set(m.equipmentId, m);
    }
    return map;
  });

  private static readonly STORAGE_KEY = 'dt-columns-equipment';

  savedHiddenColumns = EquipmentList.loadHiddenColumns();

  onHiddenColumnsChange(labels: string[]) {
    localStorage.setItem(EquipmentList.STORAGE_KEY, JSON.stringify(labels));
  }

  private static loadHiddenColumns(): string[] | undefined {
    try {
      const stored = localStorage.getItem(EquipmentList.STORAGE_KEY);
      return stored ? JSON.parse(stored) : undefined;
    } catch {
      return undefined;
    }
  }

  equipmentColumns: DataTableColumn[] = [
    { label: 'Name', labelKey: 'fields.name', searchParam: 'q' },
    { label: 'Mastery', labelKey: 'fields.mastery' },
    {
      label: 'Category',
      labelKey: 'fields.category',
      filterParam: 'category',
      optionKeyPrefix: 'enums.equipmentCategory',
      options: [
        EquipmentCategoryFreeWeights,
        EquipmentCategoryAccessories,
        EquipmentCategoryBenches,
        EquipmentCategoryMachines,
        EquipmentCategoryFunctional,
        EquipmentCategoryOther,
      ],
    },
    { label: 'Description', labelKey: 'fields.description' },
    { label: 'Internal name', labelKey: 'fields.internalName', defaultHidden: true },
    { label: 'Version', labelKey: 'fields.version', defaultHidden: true },
    { label: 'Created at', labelKey: 'fields.createdAt', defaultHidden: true },
    { label: 'Updated at', labelKey: 'fields.updatedAt', defaultHidden: true },
  ];
}
