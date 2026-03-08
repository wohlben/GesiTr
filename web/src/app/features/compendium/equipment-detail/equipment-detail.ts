import { Component, inject, computed } from '@angular/core';
import { ActivatedRoute, RouterLink } from '@angular/router';
import { toSignal } from '@angular/core/rxjs-interop';
import { injectQuery } from '@tanstack/angular-query-experimental';
import { CompendiumApiClient } from '$core/api-clients/compendium-api-client';
import { equipmentKeys } from '$core/query-keys';
import { PageLayout } from '../../../layout/page-layout';

@Component({
  selector: 'app-equipment-detail',
  imports: [PageLayout, RouterLink],
  template: `
    <app-page-layout
      [header]="equipmentQuery.data()?.displayName ?? 'Equipment'"
      [isPending]="equipmentQuery.isPending()"
      [errorMessage]="equipmentQuery.isError() ? equipmentQuery.error().message : undefined"
    >
      <a
        actions
        routerLink="./edit"
        class="rounded-md bg-blue-600 px-4 py-2 text-sm font-medium text-white hover:bg-blue-700"
        >Edit</a
      >
      @if (equipmentQuery.data(); as equipment) {
        <dl class="grid grid-cols-1 gap-4 sm:grid-cols-2">
          <div>
            <dt class="text-sm font-medium text-gray-500 dark:text-gray-400">Category</dt>
            <dd class="text-sm text-gray-900 dark:text-gray-100">{{ equipment.category }}</dd>
          </div>
          <div>
            <dt class="text-sm font-medium text-gray-500 dark:text-gray-400">Name</dt>
            <dd class="text-sm text-gray-900 dark:text-gray-100">{{ equipment.name }}</dd>
          </div>
          <div class="sm:col-span-2">
            <dt class="text-sm font-medium text-gray-500 dark:text-gray-400">Description</dt>
            <dd class="text-sm text-gray-900 dark:text-gray-100">{{ equipment.description }}</dd>
          </div>
        </dl>
      }
    </app-page-layout>
  `,
})
export class EquipmentDetail {
  private api = inject(CompendiumApiClient);
  private params = toSignal(inject(ActivatedRoute).paramMap);

  private id = computed(() => Number(this.params()?.get('id')));

  equipmentQuery = injectQuery(() => ({
    queryKey: equipmentKeys.detail(this.id()),
    queryFn: () => this.api.fetchEquipmentItem(this.id()),
    enabled: !!this.id(),
  }));
}
