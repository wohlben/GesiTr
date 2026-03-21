import { Component, inject, computed, signal } from '@angular/core';
import { ActivatedRoute, Router, RouterLink } from '@angular/router';
import { toSignal } from '@angular/core/rxjs-interop';
import {
  injectQuery,
  injectMutation,
  injectQueryClient,
} from '@tanstack/angular-query-experimental';
import { CompendiumApiClient } from '$core/api-clients/compendium-api-client';
import { UserApiClient } from '$core/api-clients/user-api-client';
import { equipmentKeys, userEquipmentKeys } from '$core/query-keys';
import { PageLayout } from '../../../layout/page-layout';
import { ConfirmDialog } from '$ui/confirm-dialog/confirm-dialog';

@Component({
  selector: 'app-equipment-detail',
  imports: [PageLayout, RouterLink, ConfirmDialog],
  template: `
    <app-page-layout
      [header]="equipmentQuery.data()?.displayName ?? 'Equipment'"
      [isPending]="equipmentQuery.isPending()"
      [errorMessage]="equipmentQuery.isError() ? equipmentQuery.error().message : undefined"
    >
      <div actions class="flex gap-2">
        <button
          type="button"
          (click)="addMutation.mutate()"
          [disabled]="addMutation.isPending()"
          class="rounded-md bg-green-600 px-4 py-2 text-sm font-medium text-white hover:bg-green-700 disabled:opacity-50"
        >
          {{ addMutation.isPending() ? 'Adding...' : 'Add to My Equipment' }}
        </button>
        @if (hasHistory()) {
          <a
            routerLink="./history"
            class="rounded-md bg-gray-600 px-4 py-2 text-sm font-medium text-white hover:bg-gray-700"
            >History</a
          >
        }
        <a
          routerLink="./edit"
          class="rounded-md bg-blue-600 px-4 py-2 text-sm font-medium text-white hover:bg-blue-700"
          >Edit</a
        >
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
        title="Delete Equipment"
        [message]="deleteMessage()"
        [isPending]="deleteMutation.isPending()"
        (confirmed)="deleteMutation.mutate()"
        (cancelled)="showDeleteDialog.set(false)"
      />
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
  private userApi = inject(UserApiClient);
  private router = inject(Router);
  private queryClient = injectQueryClient();
  private params = toSignal(inject(ActivatedRoute).paramMap);

  private id = computed(() => Number(this.params()?.get('id')));

  showDeleteDialog = signal(false);
  deleteMessage = computed(
    () => `Are you sure you want to delete '${this.equipmentQuery.data()?.displayName ?? ''}'?`,
  );

  equipmentQuery = injectQuery(() => ({
    queryKey: equipmentKeys.detail(this.id()),
    queryFn: () => this.api.fetchEquipmentItem(this.id()),
    enabled: !!this.id(),
  }));

  versionsQuery = injectQuery(() => ({
    queryKey: equipmentKeys.versions(this.id()),
    queryFn: () => this.api.fetchEquipmentVersions(this.id()),
    enabled: !!this.id(),
  }));

  hasHistory = computed(() => (this.versionsQuery.data()?.length ?? 0) > 1);

  deleteMutation = injectMutation(() => ({
    mutationFn: () => this.api.deleteEquipment(this.id()),
    onSuccess: () => {
      this.queryClient.invalidateQueries({ queryKey: equipmentKeys.all() });
      this.router.navigate(['/compendium/equipment']);
    },
  }));

  addMutation = injectMutation(() => ({
    mutationFn: () => {
      const equipment = this.equipmentQuery.data()!;
      return this.userApi.createUserEquipment({
        compendiumEquipmentId: equipment.templateId,
        compendiumVersion: equipment.version,
      });
    },
    onSuccess: (created) => {
      this.queryClient.invalidateQueries({ queryKey: userEquipmentKeys.all() });
      this.router.navigate(['/user/equipment', created.id]);
    },
  }));
}
