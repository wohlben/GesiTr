import { Component, inject, computed } from '@angular/core';
import { ActivatedRoute, RouterLink } from '@angular/router';
import { toSignal } from '@angular/core/rxjs-interop';
import { injectQuery, keepPreviousData } from '@tanstack/angular-query-experimental';
import { CompendiumApiClient } from '$core/api-clients/compendium-api-client';
import { equipmentKeys } from '$core/query-keys';
import { EquipmentListItem } from '$ui/compendium/equipment-list-item/equipment-list-item';
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
  imports: [EquipmentListItem, DataTable, Pagination, PageLayout, RouterLink],
  template: `
    <app-page-layout
      header="Equipment"
      [isPending]="equipmentQuery.isPending()"
      [errorMessage]="equipmentQuery.isError() ? equipmentQuery.error().message : undefined"
    >
      <a
        actions
        routerLink="./new"
        class="rounded-md bg-blue-600 px-4 py-2 text-sm font-medium text-white hover:bg-blue-700"
        >New</a
      >
      @if (equipmentQuery.data(); as page) {
        <app-data-table
          [columns]="equipmentColumns"
          [stale]="equipmentQuery.isPlaceholderData()"
          [initialHiddenColumns]="savedHiddenColumns"
          (hiddenColumnsChange)="onHiddenColumnsChange($event)"
        >
          @for (item of page.items; track item.id) {
            <tr app-equipment-list-item [equipment]="item"></tr>
          }
        </app-data-table>
        <app-pagination [page]="page" emptyLabel="No equipment found" />
      }
    </app-page-layout>
  `,
})
export class EquipmentList {
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

  equipmentQuery = injectQuery(() => ({
    queryKey: equipmentKeys.list(this.filters()),
    queryFn: () => this.api.fetchEquipment(this.filters()),
    placeholderData: keepPreviousData,
  }));

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
    { label: 'Name', searchParam: 'q' },
    {
      label: 'Category',
      filterParam: 'category',
      options: [
        EquipmentCategoryFreeWeights,
        EquipmentCategoryAccessories,
        EquipmentCategoryBenches,
        EquipmentCategoryMachines,
        EquipmentCategoryFunctional,
        EquipmentCategoryOther,
      ],
    },
    { label: 'Description' },
    { label: 'Internal name', defaultHidden: true },
    { label: 'Version', defaultHidden: true },
    { label: 'Created by', defaultHidden: true },
    { label: 'Created at', defaultHidden: true },
    { label: 'Updated at', defaultHidden: true },
  ];
}
