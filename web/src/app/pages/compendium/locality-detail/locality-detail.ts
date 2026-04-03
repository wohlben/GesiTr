import { Component, inject, computed, signal } from '@angular/core';
import { ActivatedRoute, Router, RouterLink } from '@angular/router';
import { toSignal } from '@angular/core/rxjs-interop';
import { injectQuery, injectMutation, QueryClient } from '@tanstack/angular-query-experimental';
import { CompendiumApiClient } from '$core/api-clients/compendium-api-client';
import { localityKeys, localityAvailabilityKeys, equipmentKeys } from '$core/query-keys';
import { TranslocoDirective } from '@jsverse/transloco';
import { PageLayout } from '../../../layout/page-layout';
import { ConfirmDialog } from '$ui/confirm-dialog/confirm-dialog';
import { Equipment, LocalityAvailability } from '$generated/models';
import { HlmInput } from '@spartan-ng/helm/input';
import { FormsModule } from '@angular/forms';

@Component({
  selector: 'app-locality-detail',
  imports: [PageLayout, RouterLink, ConfirmDialog, TranslocoDirective, HlmInput, FormsModule],
  template: `
    <ng-container *transloco="let t">
      <app-page-layout
        [header]="localityQuery.data()?.name ?? t('compendium.localities.title')"
        [isPending]="localityQuery.isPending()"
        [errorMessage]="localityQuery.isError() ? localityQuery.error().message : undefined"
      >
        <div actions class="flex gap-2">
          @if (canModify()) {
            <a
              routerLink="./edit"
              class="rounded-md bg-blue-600 px-4 py-2 text-sm font-medium text-white hover:bg-blue-700"
              >{{ t('common.edit') }}</a
            >
          }
          @if (canDelete()) {
            <button
              type="button"
              (click)="showDeleteDialog.set(true)"
              class="rounded-md bg-red-600 px-4 py-2 text-sm font-medium text-white hover:bg-red-700"
            >
              {{ t('common.delete') }}
            </button>
          }
        </div>
        <app-confirm-dialog
          [open]="showDeleteDialog()"
          [title]="t('compendium.localities.deleteTitle')"
          [message]="
            t('compendium.localities.deleteMessage', {
              name: localityQuery.data()?.name ?? '',
            })
          "
          [isPending]="deleteMutation.isPending()"
          (confirmed)="deleteMutation.mutate()"
          (cancelled)="showDeleteDialog.set(false)"
        />
        @if (localityQuery.data(); as locality) {
          <section class="space-y-4">
            <div class="flex items-center justify-between">
              <h2 class="text-lg font-semibold text-gray-900 dark:text-gray-100">
                {{ t('compendium.localities.equipmentAtLocality') }}
              </h2>
              @if (canModify()) {
                <button
                  type="button"
                  (click)="showAddSection.set(!showAddSection())"
                  class="rounded-md bg-blue-600 px-3 py-1.5 text-sm font-medium text-white hover:bg-blue-700"
                >
                  {{
                    showAddSection() ? t('common.cancel') : t('compendium.localities.addEquipment')
                  }}
                </button>
              }
            </div>

            @if (showAddSection() && canModify()) {
              <div class="rounded-md border border-gray-200 p-4 dark:border-gray-700">
                <label
                  for="equipmentSearch"
                  class="block text-sm font-medium text-gray-700 dark:text-gray-300"
                  >{{ t('compendium.localities.searchEquipment') }}</label
                >
                <input
                  id="equipmentSearch"
                  hlmInput
                  class="mt-1"
                  [(ngModel)]="equipmentSearchTerm"
                  (ngModelChange)="onEquipmentSearch($event)"
                  [placeholder]="t('compendium.localities.searchPlaceholder')"
                />
                @if (equipmentSearchQuery.data(); as page) {
                  <ul class="mt-2 max-h-48 overflow-y-auto">
                    @for (equipment of page.items; track equipment.id) {
                      <li>
                        <button
                          type="button"
                          (click)="addEquipment(equipment)"
                          [disabled]="isAlreadyAdded(equipment.id)"
                          class="w-full px-3 py-2 text-left text-sm hover:bg-gray-100 disabled:opacity-50 dark:hover:bg-gray-800"
                        >
                          {{ equipment.displayName }}
                          @if (isAlreadyAdded(equipment.id)) {
                            <span class="text-xs text-gray-400"
                              >({{ t('compendium.localities.alreadyAdded') }})</span
                            >
                          }
                        </button>
                      </li>
                    } @empty {
                      <li class="px-3 py-2 text-sm text-gray-500">{{ t('common.noResults') }}</li>
                    }
                  </ul>
                }
              </div>
            }

            @if (enrichedAvailabilities(); as items) {
              @if (items.length === 0) {
                <p class="text-sm text-gray-500 dark:text-gray-400">
                  {{ t('compendium.localities.noEquipment') }}
                </p>
              } @else {
                <ul
                  class="divide-y divide-gray-200 rounded-md border border-gray-200 dark:divide-gray-700 dark:border-gray-700"
                >
                  @for (item of items; track item.availability.id) {
                    <li class="flex items-center justify-between px-4 py-3">
                      <div class="flex items-center gap-3">
                        <button
                          type="button"
                          (click)="toggleAvailability(item.availability)"
                          [title]="t('compendium.localities.toggleAvailability')"
                          class="flex h-5 w-9 shrink-0 items-center rounded-full transition-colors"
                          [class]="
                            item.availability.available
                              ? 'bg-green-500'
                              : 'bg-gray-300 dark:bg-gray-600'
                          "
                          [disabled]="!canModify()"
                        >
                          <span
                            class="inline-block h-4 w-4 rounded-full bg-white shadow transition-transform"
                            [class]="
                              item.availability.available ? 'translate-x-4' : 'translate-x-0.5'
                            "
                          ></span>
                        </button>
                        <span
                          class="text-sm"
                          [class]="
                            item.availability.available
                              ? 'text-gray-900 dark:text-gray-100'
                              : 'text-gray-400 line-through dark:text-gray-500'
                          "
                        >
                          {{ item.equipmentName }}
                        </span>
                      </div>
                      @if (canModify()) {
                        <button
                          type="button"
                          (click)="removeAvailability(item.availability.id)"
                          class="text-sm text-red-600 hover:text-red-800 dark:text-red-400 dark:hover:text-red-300"
                        >
                          {{ t('common.remove') }}
                        </button>
                      }
                    </li>
                  }
                </ul>
              }
            }
          </section>
        }
      </app-page-layout>
    </ng-container>
  `,
})
export class LocalityDetail {
  private api = inject(CompendiumApiClient);
  private router = inject(Router);
  private queryClient = inject(QueryClient);
  private params = toSignal(inject(ActivatedRoute).paramMap);

  private id = computed(() => Number(this.params()?.get('id')));

  showDeleteDialog = signal(false);
  showAddSection = signal(false);
  equipmentSearchTerm = '';
  equipmentSearch = signal('');

  localityQuery = injectQuery(() => ({
    queryKey: localityKeys.detail(this.id()),
    queryFn: () => this.api.fetchLocality(this.id()),
    enabled: !!this.id(),
  }));

  permissionsQuery = injectQuery(() => ({
    queryKey: localityKeys.permissions(this.id()),
    queryFn: () => this.api.fetchLocalityPermissions(this.id()),
    enabled: !!this.id(),
  }));

  availabilityQuery = injectQuery(() => ({
    queryKey: localityAvailabilityKeys.list({ localityId: this.id() }),
    queryFn: () => this.api.fetchLocalityAvailabilities({ localityId: this.id() }),
    enabled: !!this.id(),
  }));

  equipmentQuery = injectQuery(() => ({
    queryKey: equipmentKeys.list({ owner: 'me' }),
    queryFn: () => this.api.fetchEquipment({ owner: 'me', limit: 1000 }),
    enabled: !!this.id(),
  }));

  equipmentSearchQuery = injectQuery(() => ({
    queryKey: equipmentKeys.list({ q: this.equipmentSearch(), limit: 20 }),
    queryFn: () => this.api.fetchEquipment({ q: this.equipmentSearch(), limit: 20 }),
    enabled: this.showAddSection(),
  }));

  canModify = computed(
    () =>
      this.permissionsQuery.isSuccess() &&
      (this.permissionsQuery.data()?.permissions?.includes('MODIFY') ?? false),
  );
  canDelete = computed(
    () =>
      this.permissionsQuery.isSuccess() &&
      (this.permissionsQuery.data()?.permissions?.includes('DELETE') ?? false),
  );

  private equipmentMap = computed(() => {
    const map = new Map<number, Equipment>();
    for (const e of this.equipmentQuery.data()?.items ?? []) {
      map.set(e.id, e);
    }
    return map;
  });

  private addedEquipmentIds = computed(() => {
    const set = new Set<number>();
    for (const a of this.availabilityQuery.data() ?? []) {
      set.add(a.equipmentId);
    }
    return set;
  });

  enrichedAvailabilities = computed(() => {
    const avails = this.availabilityQuery.data();
    if (!avails) return undefined;
    const eqMap = this.equipmentMap();
    return avails.map((a) => ({
      availability: a,
      equipmentName: eqMap.get(a.equipmentId)?.displayName ?? `Equipment #${a.equipmentId}`,
    }));
  });

  isAlreadyAdded(equipmentId: number): boolean {
    return this.addedEquipmentIds().has(equipmentId);
  }

  onEquipmentSearch(term: string) {
    this.equipmentSearch.set(term);
  }

  deleteMutation = injectMutation(() => ({
    mutationFn: () => this.api.deleteLocality(this.id()),
    onSuccess: () => {
      this.queryClient.invalidateQueries({ queryKey: localityKeys.all() });
      this.router.navigate(['/compendium/localities']);
    },
  }));

  addEquipmentMutation = injectMutation(() => ({
    mutationFn: (equipmentId: number) =>
      this.api.createLocalityAvailability({ localityId: this.id(), equipmentId }),
    onSuccess: () => {
      this.queryClient.invalidateQueries({
        queryKey: localityAvailabilityKeys.list({ localityId: this.id() }),
      });
    },
  }));

  toggleMutation = injectMutation(() => ({
    mutationFn: (vars: { id: number; available: boolean }) =>
      this.api.updateLocalityAvailability(vars.id, { available: vars.available }),
    onSuccess: () => {
      this.queryClient.invalidateQueries({
        queryKey: localityAvailabilityKeys.list({ localityId: this.id() }),
      });
    },
  }));

  removeMutation = injectMutation(() => ({
    mutationFn: (availabilityId: number) => this.api.deleteLocalityAvailability(availabilityId),
    onSuccess: () => {
      this.queryClient.invalidateQueries({
        queryKey: localityAvailabilityKeys.list({ localityId: this.id() }),
      });
    },
  }));

  addEquipment(equipment: Equipment) {
    this.addEquipmentMutation.mutate(equipment.id);
  }

  toggleAvailability(availability: LocalityAvailability) {
    this.toggleMutation.mutate({ id: availability.id, available: !availability.available });
  }

  removeAvailability(id: number) {
    this.removeMutation.mutate(id);
  }
}
