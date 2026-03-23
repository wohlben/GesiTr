import { Component, inject, computed, signal } from '@angular/core';
import { ActivatedRoute, Router } from '@angular/router';
import { toSignal } from '@angular/core/rxjs-interop';
import { injectQuery, injectMutation, QueryClient } from '@tanstack/angular-query-experimental';
import { TranslocoDirective, TranslocoService } from '@jsverse/transloco';
import { UserApiClient } from '$core/api-clients/user-api-client';
import { userEquipmentKeys } from '$core/query-keys';
import { PageLayout } from '../../../layout/page-layout';
import { ConfirmDialog } from '$ui/confirm-dialog/confirm-dialog';

@Component({
  selector: 'app-user-equipment-detail',
  imports: [PageLayout, ConfirmDialog, TranslocoDirective],
  template: `
    <ng-container *transloco="let t">
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
            {{ t('common.delete') }}
          </button>
        </div>
        <app-confirm-dialog
          [open]="showDeleteDialog()"
          [title]="t('user.equipment.deleteTitle')"
          [message]="t('user.equipment.deleteMessage', { name: equipmentName() })"
          [isPending]="deleteMutation.isPending()"
          (confirmed)="deleteMutation.mutate()"
          (cancelled)="showDeleteDialog.set(false)"
        />
        @if (detailQuery.data(); as equipment) {
          <dl class="grid grid-cols-1 gap-4 sm:grid-cols-2">
            <div>
              <dt class="text-sm font-medium text-gray-500 dark:text-gray-400">
                {{ t('fields.displayName') }}
              </dt>
              <dd class="text-sm text-gray-900 dark:text-gray-100">{{ equipment.displayName }}</dd>
            </div>
            <div>
              <dt class="text-sm font-medium text-gray-500 dark:text-gray-400">
                {{ t('fields.category') }}
              </dt>
              <dd class="text-sm text-gray-900 dark:text-gray-100">{{ equipment.category }}</dd>
            </div>
            <div>
              <dt class="text-sm font-medium text-gray-500 dark:text-gray-400">
                {{ t('fields.name') }}
              </dt>
              <dd class="text-sm text-gray-900 dark:text-gray-100">{{ equipment.name }}</dd>
            </div>
            <div>
              <dt class="text-sm font-medium text-gray-500 dark:text-gray-400">
                {{ t('fields.version') }}
              </dt>
              <dd class="text-sm text-gray-900 dark:text-gray-100">v{{ equipment.version }}</dd>
            </div>
            <div class="sm:col-span-2">
              <dt class="text-sm font-medium text-gray-500 dark:text-gray-400">
                {{ t('fields.description') }}
              </dt>
              <dd class="text-sm text-gray-900 dark:text-gray-100">{{ equipment.description }}</dd>
            </div>
          </dl>
        }
      </app-page-layout>
    </ng-container>
  `,
})
export class UserEquipmentDetail {
  private userApi = inject(UserApiClient);
  private router = inject(Router);
  private queryClient = inject(QueryClient);
  private transloco = inject(TranslocoService);
  private params = toSignal(inject(ActivatedRoute).paramMap);

  private id = computed(() => Number(this.params()?.get('id')));

  showDeleteDialog = signal(false);

  detailQuery = injectQuery(() => ({
    queryKey: userEquipmentKeys.detail(this.id()),
    queryFn: () => this.userApi.fetchUserEquipmentItem(this.id()),
    enabled: !!this.id(),
  }));

  equipmentName = computed(
    () => this.detailQuery.data()?.displayName ?? this.transloco.translate('common.loading'),
  );

  deleteMutation = injectMutation(() => ({
    mutationFn: () => this.userApi.deleteUserEquipment(this.id()),
    onSuccess: () => {
      this.queryClient.invalidateQueries({ queryKey: userEquipmentKeys.all() });
      this.router.navigate(['/user/equipment']);
    },
  }));
}
