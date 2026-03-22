import { Component, inject, computed } from '@angular/core';
import { RouterLink } from '@angular/router';
import { injectQuery } from '@tanstack/angular-query-experimental';
import { injectQueries } from '@tanstack/angular-query-experimental/inject-queries-experimental';
import { TranslocoDirective } from '@jsverse/transloco';
import { UserApiClient } from '$core/api-clients/user-api-client';
import { CompendiumApiClient } from '$core/api-clients/compendium-api-client';
import { userEquipmentKeys, equipmentKeys } from '$core/query-keys';
import { PageLayout } from '../../../layout/page-layout';

@Component({
  selector: 'app-user-equipment-list',
  imports: [PageLayout, RouterLink, TranslocoDirective],
  template: `
    <ng-container *transloco="let t">
      <app-page-layout
        [header]="t('user.equipment.title')"
        [isPending]="userEquipmentQuery.isPending()"
        [errorMessage]="
          userEquipmentQuery.isError() ? userEquipmentQuery.error().message : undefined
        "
      >
        @if (enrichedEquipment(); as equipment) {
          @if (equipment.length === 0) {
            <p class="text-sm text-gray-500 dark:text-gray-400">
              {{ t('user.equipment.noResults') }}
            </p>
          } @else {
            <div class="overflow-x-auto">
              <table class="min-w-full divide-y divide-gray-200 dark:divide-gray-700">
                <thead>
                  <tr>
                    <th
                      class="px-4 py-3 text-left text-xs font-medium tracking-wider text-gray-500 uppercase dark:text-gray-400"
                    >
                      {{ t('fields.name') }}
                    </th>
                    <th
                      class="px-4 py-3 text-left text-xs font-medium tracking-wider text-gray-500 uppercase dark:text-gray-400"
                    >
                      {{ t('fields.category') }}
                    </th>
                    <th
                      class="px-4 py-3 text-left text-xs font-medium tracking-wider text-gray-500 uppercase dark:text-gray-400"
                    >
                      {{ t('fields.version') }}
                    </th>
                  </tr>
                </thead>
                <tbody class="divide-y divide-gray-200 dark:divide-gray-700">
                  @for (item of equipment; track item.id) {
                    <tr
                      class="hover:bg-gray-50 dark:hover:bg-gray-800/50"
                      [routerLink]="['./', item.id]"
                      class="cursor-pointer"
                    >
                      <td class="px-4 py-3 text-sm text-gray-900 dark:text-gray-100">
                        {{ item.displayName }}
                      </td>
                      <td class="px-4 py-3 text-sm text-gray-600 dark:text-gray-400">
                        {{ item.category }}
                      </td>
                      <td class="px-4 py-3 text-sm text-gray-600 dark:text-gray-400">
                        v{{ item.compendiumVersion }}
                      </td>
                    </tr>
                  }
                </tbody>
              </table>
            </div>
          }
        }
      </app-page-layout>
    </ng-container>
  `,
})
export class UserEquipmentList {
  private userApi = inject(UserApiClient);
  private compendiumApi = inject(CompendiumApiClient);

  userEquipmentQuery = injectQuery(() => ({
    queryKey: userEquipmentKeys.list(),
    queryFn: () => this.userApi.fetchUserEquipment(),
  }));

  private snapshotQueries = injectQueries(() => ({
    queries: (this.userEquipmentQuery.data() ?? []).map((ue) => ({
      queryKey: equipmentKeys.version(ue.compendiumEquipmentId, ue.compendiumVersion),
      queryFn: () =>
        this.compendiumApi.fetchEquipmentVersion(ue.compendiumEquipmentId, ue.compendiumVersion),
      staleTime: Infinity,
    })),
  }));

  enrichedEquipment = computed(() => {
    const userEquipment = this.userEquipmentQuery.data();
    if (!userEquipment) return undefined;

    const snapshots = this.snapshotQueries();

    return userEquipment.map((ue, i) => {
      const versionEntry = snapshots[i]?.data();
      const equipment = versionEntry?.snapshot;
      return {
        id: ue.id,
        compendiumVersion: ue.compendiumVersion,
        displayName: equipment?.displayName ?? '...',
        category: equipment?.category ?? '',
      };
    });
  });
}
