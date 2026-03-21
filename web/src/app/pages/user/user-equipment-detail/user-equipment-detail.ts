import { Component, inject, computed, signal } from '@angular/core';
import { ActivatedRoute, Router } from '@angular/router';
import { toSignal } from '@angular/core/rxjs-interop';
import {
  injectQuery,
  injectMutation,
  injectQueryClient,
} from '@tanstack/angular-query-experimental';
import { UserApiClient } from '$core/api-clients/user-api-client';
import { CompendiumApiClient } from '$core/api-clients/compendium-api-client';
import { userEquipmentKeys } from '$core/query-keys';
import { PageLayout } from '../../../layout/page-layout';
import { ConfirmDialog } from '$ui/confirm-dialog/confirm-dialog';
import { Equipment } from '$generated/models';

@Component({
  selector: 'app-user-equipment-detail',
  imports: [PageLayout, ConfirmDialog],
  template: `
    <app-page-layout
      [header]="equipmentName()"
      [isPending]="detailQuery.isPending()"
      [errorMessage]="detailQuery.isError() ? detailQuery.error().message : undefined"
    >
      <div actions class="flex gap-2">
        <button
          type="button"
          (click)="showDeleteDialog.set(true)"
          class="rounded-md bg-red-600 px-4 py-2 text-sm font-medium text-white hover:bg-red-700"
        >
          Delete
        </button>
      </div>
      <app-confirm-dialog
        [open]="showDeleteDialog()"
        title="Remove Equipment"
        [message]="deleteMessage()"
        [isPending]="deleteMutation.isPending()"
        (confirmed)="deleteMutation.mutate()"
        (cancelled)="showDeleteDialog.set(false)"
      />
      @if (snapshot(); as equipment) {
        <dl class="grid grid-cols-1 gap-4 sm:grid-cols-2">
          <div>
            <dt class="text-sm font-medium text-gray-500 dark:text-gray-400">Display Name</dt>
            <dd class="text-sm text-gray-900 dark:text-gray-100">{{ equipment.displayName }}</dd>
          </div>
          <div>
            <dt class="text-sm font-medium text-gray-500 dark:text-gray-400">Category</dt>
            <dd class="text-sm text-gray-900 dark:text-gray-100">{{ equipment.category }}</dd>
          </div>
          <div>
            <dt class="text-sm font-medium text-gray-500 dark:text-gray-400">Name</dt>
            <dd class="text-sm text-gray-900 dark:text-gray-100">{{ equipment.name }}</dd>
          </div>
          <div>
            <dt class="text-sm font-medium text-gray-500 dark:text-gray-400">Version</dt>
            <dd class="text-sm text-gray-900 dark:text-gray-100">
              v{{ detailQuery.data()?.userEquipment?.compendiumVersion }}
            </dd>
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
export class UserEquipmentDetail {
  private userApi = inject(UserApiClient);
  private compendiumApi = inject(CompendiumApiClient);
  private router = inject(Router);
  private queryClient = injectQueryClient();
  private params = toSignal(inject(ActivatedRoute).paramMap);

  private id = computed(() => Number(this.params()?.get('id')));

  showDeleteDialog = signal(false);

  detailQuery = injectQuery(() => ({
    queryKey: userEquipmentKeys.detail(this.id()),
    queryFn: async () => {
      const userEquipment = await this.userApi.fetchUserEquipmentItem(this.id());
      const version = await this.compendiumApi.fetchEquipmentVersion(
        userEquipment.compendiumEquipmentId,
        userEquipment.compendiumVersion,
      );
      return { userEquipment, version };
    },
    enabled: !!this.id(),
  }));

  snapshot = computed(() => this.detailQuery.data()?.version.snapshot as Equipment | undefined);

  equipmentName = computed(() => this.snapshot()?.displayName ?? 'Equipment');

  deleteMessage = computed(
    () => `Are you sure you want to remove '${this.equipmentName()}' from your equipment?`,
  );

  deleteMutation = injectMutation(() => ({
    mutationFn: () => this.userApi.deleteUserEquipment(this.id()),
    onSuccess: () => {
      this.queryClient.invalidateQueries({ queryKey: userEquipmentKeys.all() });
      this.router.navigate(['/user/equipment']);
    },
  }));
}
