import { Component, inject, signal, computed } from '@angular/core';
import { injectQuery } from '@tanstack/angular-query-experimental';
import { CompendiumApiClient } from '$core/api-clients/compendium-api-client';
import { EquipmentListItem } from '$ui/compendium/equipment-list-item/equipment-list-item';
import { SearchInput } from '$ui/inputs/search-input/search-input';
import { FilterSelect } from '$ui/inputs/filter-select/filter-select';
import { DataTable } from '$ui/data-table/data-table';
import { PageLayout } from '../../../layout/page-layout';
import {
  EquipmentCategory,
  EquipmentCategoryFreeWeights,
  EquipmentCategoryAccessories,
  EquipmentCategoryBenches,
  EquipmentCategoryMachines,
  EquipmentCategoryFunctional,
  EquipmentCategoryOther,
} from '$generated/models';

@Component({
  selector: 'app-equipment-list',
  imports: [SearchInput, FilterSelect, EquipmentListItem, DataTable, PageLayout],
  template: `
    <app-page-layout
      header="Equipment"
      [isPending]="equipmentQuery.isPending()"
      [errorMessage]="equipmentQuery.isError() ? equipmentQuery.error().message : undefined">

      <div filters class="flex flex-wrap gap-3">
        <app-search-input placeholder="Search equipment..." [(value)]="q" />
        <app-filter-select allLabel="All categories" [options]="categoryOptions" [(value)]="category" />
      </div>

      @if (equipmentQuery.data(); as items) {
        <app-data-table [columns]="['Name', 'Category', 'Description']">
          @for (item of items; track item.id) {
            <tr app-equipment-list-item [equipment]="item"></tr>
          }
        </app-data-table>
        <p class="text-sm text-gray-500 dark:text-gray-400">{{ items.length }} items</p>
      }
    </app-page-layout>
  `,
})
export class EquipmentList {
  private api = inject(CompendiumApiClient);

  q = signal('');
  category = signal<EquipmentCategory | ''>('');

  filters = computed(() => ({
    q: this.q() || undefined,
    category: this.category() || undefined,
  }));

  equipmentQuery = injectQuery(() => ({
    queryKey: ['equipment', this.filters()],
    queryFn: () => this.api.fetchEquipment(this.filters()),
  }));

  categoryOptions: EquipmentCategory[] = [
    EquipmentCategoryFreeWeights,
    EquipmentCategoryAccessories,
    EquipmentCategoryBenches,
    EquipmentCategoryMachines,
    EquipmentCategoryFunctional,
    EquipmentCategoryOther,
  ];
}
