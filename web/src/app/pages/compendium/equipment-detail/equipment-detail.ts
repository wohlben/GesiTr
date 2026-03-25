import { Component, inject, computed, signal } from '@angular/core';
import { ActivatedRoute, Router, RouterLink } from '@angular/router';
import { toSignal } from '@angular/core/rxjs-interop';
import { injectQuery, injectMutation, QueryClient } from '@tanstack/angular-query-experimental';
import { CompendiumApiClient } from '$core/api-clients/compendium-api-client';
import { UserApiClient } from '$core/api-clients/user-api-client';
import { equipmentKeys, equipmentRelationshipKeys, userEquipmentKeys } from '$core/query-keys';
import { TranslocoDirective } from '@jsverse/transloco';
import { PageLayout } from '../../../layout/page-layout';
import { ConfirmDialog } from '$ui/confirm-dialog/confirm-dialog';

@Component({
  selector: 'app-equipment-detail',
  imports: [PageLayout, RouterLink, ConfirmDialog, TranslocoDirective],
  template: `
    <ng-container *transloco="let t">
      <app-page-layout
        [header]="equipmentQuery.data()?.displayName ?? 'Equipment'"
        [isPending]="equipmentQuery.isPending()"
        [errorMessage]="equipmentQuery.isError() ? equipmentQuery.error().message : undefined"
      >
        <div actions class="flex gap-2">
          @if (alreadyAdded(); as existing) {
            <a
              [routerLink]="['/user/equipment', existing.id]"
              class="rounded-md bg-gray-500 px-4 py-2 text-sm font-medium text-white hover:bg-gray-600"
            >
              {{ t('compendium.equipment.alreadyAdded') }}
            </a>
          } @else {
            <button
              type="button"
              (click)="addMutation.mutate()"
              [disabled]="addMutation.isPending()"
              class="rounded-md bg-green-600 px-4 py-2 text-sm font-medium text-white hover:bg-green-700 disabled:opacity-50"
            >
              {{
                addMutation.isPending() ? t('common.adding') : t('compendium.equipment.addToMine')
              }}
            </button>
          }
          @if (hasHistory()) {
            <a
              routerLink="./history"
              class="rounded-md bg-gray-600 px-4 py-2 text-sm font-medium text-white hover:bg-gray-700"
              >{{ t('common.history') }}</a
            >
          }
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
          [title]="t('compendium.equipment.deleteTitle')"
          [message]="
            t('compendium.equipment.deleteMessage', {
              name: equipmentQuery.data()?.displayName ?? '',
            })
          "
          [isPending]="deleteMutation.isPending()"
          (confirmed)="deleteMutation.mutate()"
          (cancelled)="showDeleteDialog.set(false)"
        />
        @if (equipmentQuery.data(); as equipment) {
          <dl class="grid grid-cols-1 gap-4 sm:grid-cols-2">
            <div>
              <dt class="text-sm font-medium text-gray-500 dark:text-gray-400">
                {{ t('fields.category') }}
              </dt>
              <dd class="text-sm text-gray-900 dark:text-gray-100">
                {{ t('enums.equipmentCategory.' + equipment.category) }}
              </dd>
            </div>
            <div>
              <dt class="text-sm font-medium text-gray-500 dark:text-gray-400">
                {{ t('fields.name') }}
              </dt>
              <dd class="text-sm text-gray-900 dark:text-gray-100">{{ equipment.name }}</dd>
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
export class EquipmentDetail {
  private api = inject(CompendiumApiClient);
  private userApi = inject(UserApiClient);
  private router = inject(Router);
  private queryClient = inject(QueryClient);
  private params = toSignal(inject(ActivatedRoute).paramMap);

  private id = computed(() => Number(this.params()?.get('id')));

  showDeleteDialog = signal(false);

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

  forkedRelationshipsQuery = injectQuery(() => ({
    queryKey: equipmentRelationshipKeys.list({
      toEquipmentId: this.id(),
      relationshipType: 'forked',
      owner: 'me',
    }),
    queryFn: () =>
      this.api.fetchEquipmentRelationships({
        toEquipmentId: this.id(),
        relationshipType: 'forked',
        owner: 'me',
      }),
    enabled: !!this.id(),
  }));

  permissionsQuery = injectQuery(() => ({
    queryKey: equipmentKeys.permissions(this.id()),
    queryFn: () => this.api.fetchEquipmentPermissions(this.id()),
    enabled: !!this.id(),
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

  hasHistory = computed(() => (this.versionsQuery.data()?.length ?? 0) > 1);

  alreadyAdded = computed(() => {
    const rels = this.forkedRelationshipsQuery.data();
    if (!rels || rels.length === 0) return undefined;
    return { id: rels[0].fromEquipmentId };
  });

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
        name: equipment.name,
        displayName: equipment.displayName,
        description: equipment.description,
        category: equipment.category,
        imageUrl: equipment.imageUrl,
        public: false,
        sourceEquipmentId: equipment.id,
      } as Record<string, unknown>);
    },
    onSuccess: (created) => {
      this.queryClient.invalidateQueries({ queryKey: userEquipmentKeys.all() });
      this.queryClient.invalidateQueries({
        queryKey: equipmentRelationshipKeys.all(),
      });
      this.router.navigate(['/user/equipment', created.id]);
    },
  }));
}
