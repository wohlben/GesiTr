import { Component, inject, signal, computed, effect, untracked } from '@angular/core';
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

      @if (equipmentQuery.data(); as page) {
        <app-data-table [columns]="['Name', 'Category', 'Description']">
          @for (item of page.items; track item.id) {
            <tr app-equipment-list-item [equipment]="item"></tr>
          }
        </app-data-table>
        <div class="flex items-center justify-between text-sm text-gray-500 dark:text-gray-400">
          <p>{{ page.total === 0 ? 'No equipment found' : 'Showing ' + (page.offset + 1) + '–' + (page.offset + page.items.length) + ' of ' + page.total + ' items' }}</p>
          <div class="flex gap-2">
            <button
              class="rounded border border-gray-300 px-3 py-1 disabled:opacity-50 dark:border-gray-600"
              [disabled]="page.offset === 0"
              (click)="prevPage()">Previous</button>
            <button
              class="rounded border border-gray-300 px-3 py-1 disabled:opacity-50 dark:border-gray-600"
              [disabled]="page.offset + page.limit >= page.total"
              (click)="nextPage()">Next</button>
          </div>
        </div>
      }
    </app-page-layout>
  `,
})
export class EquipmentList {
  private api = inject(CompendiumApiClient);

  q = signal('');
  category = signal<EquipmentCategory | ''>('');
  offset = signal(0);

  private resetOffset = effect(() => {
    this.q(); this.category();
    untracked(() => this.offset.set(0));
  });

  filters = computed(() => ({
    q: this.q() || undefined,
    category: this.category() || undefined,
    offset: this.offset() || undefined,
  }));

  equipmentQuery = injectQuery(() => ({
    queryKey: ['equipment', this.filters()],
    queryFn: () => this.api.fetchEquipment(this.filters()),
  }));

  prevPage() {
    const page = this.equipmentQuery.data();
    if (page) this.offset.set(Math.max(0, page.offset - page.limit));
  }

  nextPage() {
    const page = this.equipmentQuery.data();
    if (page) this.offset.set(page.offset + page.limit);
  }

  categoryOptions: EquipmentCategory[] = [
    EquipmentCategoryFreeWeights,
    EquipmentCategoryAccessories,
    EquipmentCategoryBenches,
    EquipmentCategoryMachines,
    EquipmentCategoryFunctional,
    EquipmentCategoryOther,
  ];
}
