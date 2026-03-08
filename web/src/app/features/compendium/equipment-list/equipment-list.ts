import { Component, inject, computed } from '@angular/core';
import { ActivatedRoute } from '@angular/router';
import { toSignal } from '@angular/core/rxjs-interop';
import { injectQuery, keepPreviousData } from '@tanstack/angular-query-experimental';
import { CompendiumApiClient } from '$core/api-clients/compendium-api-client';
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
  imports: [EquipmentListItem, DataTable, Pagination, PageLayout],
  template: `
    <app-page-layout
      header="Equipment"
      [isPending]="equipmentQuery.isPending()"
      [errorMessage]="equipmentQuery.isError() ? equipmentQuery.error().message : undefined"
    >
      @if (equipmentQuery.data(); as page) {
        <app-data-table [columns]="equipmentColumns" [stale]="equipmentQuery.isPlaceholderData()">
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
    queryKey: ['equipment', this.filters()],
    queryFn: () => this.api.fetchEquipment(this.filters()),
    placeholderData: keepPreviousData,
  }));

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
  ];
}
