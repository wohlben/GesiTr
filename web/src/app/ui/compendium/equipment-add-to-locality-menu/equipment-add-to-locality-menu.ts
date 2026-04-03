import { Component, inject, input, signal, computed } from '@angular/core';
import { Router } from '@angular/router';
import { injectQuery, injectMutation, QueryClient } from '@tanstack/angular-query-experimental';
import { CompendiumApiClient } from '$core/api-clients/compendium-api-client';
import { localityKeys, localityAvailabilityKeys } from '$core/query-keys';
import { TranslocoDirective } from '@jsverse/transloco';
import { NgIcon, provideIcons } from '@ng-icons/core';
import { lucideMapPin, lucideHome, lucideX } from '@ng-icons/lucide';
import { HlmIconImports } from '@spartan-ng/helm/icon';
import { HlmPopoverImports } from '@spartan-ng/helm/popover';

@Component({
  selector: 'app-equipment-add-to-locality-menu',
  imports: [TranslocoDirective, NgIcon, HlmIconImports, HlmPopoverImports],
  providers: [provideIcons({ lucideMapPin, lucideHome, lucideX })],
  template: `
    <ng-container *transloco="let t">
      <div hlmPopover [state]="open() ? 'open' : 'closed'" (closed)="open.set(false)" align="end">
        <button
          hlmPopoverTrigger
          (click)="open.set(!open())"
          [disabled]="addMutation.isPending() || removeMutation.isPending()"
          class="rounded-md border border-gray-300 px-4 py-2 text-sm font-medium text-gray-700 hover:bg-gray-50 disabled:opacity-50 dark:border-gray-600 dark:text-gray-300 dark:hover:bg-gray-800"
        >
          <span class="flex items-center gap-1.5">
            <ng-icon hlm name="lucideMapPin" size="sm" />
            {{ t('compendium.equipment.addTo') }}
          </span>
        </button>
        <ng-template hlmPopoverPortal>
          <div hlmPopoverContent class="w-56 p-2">
            @if (localitiesQuery.isPending()) {
              <div class="px-3 py-2 text-sm text-gray-500">
                {{ t('common.loading') }}
              </div>
            } @else if (localitiesQuery.data(); as page) {
              @if (page.items.length === 0) {
                <button
                  (click)="goToLocalities()"
                  class="flex w-full items-center gap-2 rounded-md px-3 py-2 text-left text-sm transition-colors hover:bg-gray-100 dark:hover:bg-gray-800"
                >
                  <ng-icon hlm name="lucideHome" size="sm" />
                  <div>
                    <span class="font-medium">{{ t('compendium.localities.home') }}</span>
                    <p class="text-xs text-gray-500 dark:text-gray-400">
                      {{ t('compendium.localities.homeDescription') }}
                    </p>
                  </div>
                </button>
              } @else {
                @for (locality of page.items; track locality.id) {
                  @if (availabilityMap().get(locality.id); as avail) {
                    <div
                      class="flex w-full items-center justify-between rounded-md px-3 py-1.5 text-sm"
                    >
                      <span class="text-gray-900 dark:text-gray-100">{{ locality.name }}</span>
                      <button
                        (click)="removeFromLocality(avail.id)"
                        [disabled]="removeMutation.isPending()"
                        class="rounded p-0.5 text-red-500 transition-colors hover:bg-red-100 disabled:opacity-50 dark:text-red-400 dark:hover:bg-red-900/30"
                      >
                        <ng-icon hlm name="lucideX" size="sm" />
                      </button>
                    </div>
                  } @else {
                    <button
                      (click)="addToLocality(locality.id)"
                      [disabled]="addMutation.isPending()"
                      class="w-full rounded-md px-3 py-1.5 text-left text-sm transition-colors hover:bg-gray-100 disabled:opacity-50 dark:hover:bg-gray-800"
                    >
                      {{ locality.name }}
                    </button>
                  }
                }
              }
            }
          </div>
        </ng-template>
      </div>
    </ng-container>
  `,
})
export class EquipmentAddToLocalityMenu {
  equipmentId = input.required<number>();

  private api = inject(CompendiumApiClient);
  private router = inject(Router);
  private queryClient = inject(QueryClient);

  open = signal(false);

  localitiesQuery = injectQuery(() => ({
    queryKey: localityKeys.list({ owner: 'me', limit: 100 }),
    queryFn: () => this.api.fetchLocalities({ owner: 'me', limit: 100 }),
  }));

  availabilitiesQuery = injectQuery(() => ({
    queryKey: localityAvailabilityKeys.list({ equipmentId: this.equipmentId() }),
    queryFn: () => this.api.fetchLocalityAvailabilities({ equipmentId: this.equipmentId() }),
    enabled: !!this.equipmentId(),
  }));

  availabilityMap = computed(() => {
    const map = new Map<number, { id: number; localityId: number }>();
    for (const a of this.availabilitiesQuery.data() ?? []) {
      map.set(a.localityId, { id: a.id, localityId: a.localityId });
    }
    return map;
  });

  addMutation = injectMutation(() => ({
    mutationFn: (localityId: number) =>
      this.api.createLocalityAvailability({
        localityId,
        equipmentId: this.equipmentId(),
      }),
    onSuccess: () => {
      this.queryClient.invalidateQueries({ queryKey: localityAvailabilityKeys.all() });
    },
  }));

  removeMutation = injectMutation(() => ({
    mutationFn: (availabilityId: number) => this.api.deleteLocalityAvailability(availabilityId),
    onSuccess: () => {
      this.queryClient.invalidateQueries({ queryKey: localityAvailabilityKeys.all() });
    },
  }));

  addToLocality(localityId: number) {
    this.addMutation.mutate(localityId);
  }

  removeFromLocality(availabilityId: number) {
    this.removeMutation.mutate(availabilityId);
  }

  goToLocalities() {
    this.open.set(false);
    this.router.navigate(['/compendium/localities']);
  }
}
