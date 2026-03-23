import { Component, inject, computed, effect, signal } from '@angular/core';
import { ActivatedRoute, Router, RouterLink } from '@angular/router';
import { toSignal } from '@angular/core/rxjs-interop';
import { form, required, FormField } from '@angular/forms/signals';
import { injectQuery, injectMutation, QueryClient } from '@tanstack/angular-query-experimental';
import { CompendiumApiClient } from '$core/api-clients/compendium-api-client';
import { equipmentKeys } from '$core/query-keys';
import { TranslocoDirective } from '@jsverse/transloco';
import { SlugifyPipe } from '$ui/pipes/slugify';
import { PageLayout } from '../../../layout/page-layout';
import { BrnSelectImports } from '@spartan-ng/brain/select';
import { HlmSelectImports } from '@spartan-ng/helm/select';
import { HlmInput } from '@spartan-ng/helm/input';
import { HlmTextarea } from '@spartan-ng/helm/textarea';
import {
  EquipmentCategory,
  EquipmentCategoryFreeWeights,
  EquipmentCategoryAccessories,
  EquipmentCategoryBenches,
  EquipmentCategoryMachines,
  EquipmentCategoryFunctional,
  EquipmentCategoryOther,
} from '$generated/models';

const EQUIPMENT_CATEGORIES: EquipmentCategory[] = [
  EquipmentCategoryFreeWeights,
  EquipmentCategoryAccessories,
  EquipmentCategoryBenches,
  EquipmentCategoryMachines,
  EquipmentCategoryFunctional,
  EquipmentCategoryOther,
];

@Component({
  selector: 'app-equipment-edit',
  imports: [
    PageLayout,
    FormField,
    RouterLink,
    BrnSelectImports,
    HlmSelectImports,
    HlmInput,
    HlmTextarea,
    TranslocoDirective,
  ],
  template: `
    <ng-container *transloco="let t">
      <app-page-layout
        [header]="
          isCreateMode() ? t('compendium.equipment.newTitle') : t('compendium.equipment.editTitle')
        "
        [isPending]="!isCreateMode() && equipmentQuery.isPending()"
        [errorMessage]="
          !isCreateMode() && equipmentQuery.isError() ? equipmentQuery.error().message : undefined
        "
      >
        @if (isCreateMode() || equipmentQuery.data()) {
          <form (submit)="onSubmit(); $event.preventDefault()" class="space-y-4">
            <div>
              <label for="name" class="block text-sm font-medium text-gray-700 dark:text-gray-300"
                >{{ t('fields.name') }} *</label
              >
              <input id="name" [formField]="equipmentForm.name" hlmInput class="mt-1" />
            </div>

            <div>
              <label
                for="displayName"
                class="block text-sm font-medium text-gray-700 dark:text-gray-300"
                >{{ t('fields.displayName') }} *</label
              >
              <input
                id="displayName"
                [formField]="equipmentForm.displayName"
                hlmInput
                class="mt-1"
              />
            </div>

            <div>
              <label
                for="description"
                class="block text-sm font-medium text-gray-700 dark:text-gray-300"
                >{{ t('fields.description') }}</label
              >
              <textarea
                id="description"
                [formField]="equipmentForm.description"
                rows="4"
                hlmTextarea
                class="mt-1"
              ></textarea>
            </div>

            <div>
              <label
                for="category"
                class="block text-sm font-medium text-gray-700 dark:text-gray-300"
                >{{ t('fields.category') }} *</label
              >
              <brn-select [formField]="equipmentForm.category" class="mt-1" hlm>
                <hlm-select-trigger class="w-full">
                  <hlm-select-value />
                </hlm-select-trigger>
                <hlm-select-content>
                  @for (cat of categories; track cat) {
                    <hlm-option [value]="cat">{{ t('enums.equipmentCategory.' + cat) }}</hlm-option>
                  }
                </hlm-select-content>
              </brn-select>
            </div>

            <div>
              <label
                for="imageUrl"
                class="block text-sm font-medium text-gray-700 dark:text-gray-300"
                >{{ t('fields.imageUrl') }}</label
              >
              <input id="imageUrl" [formField]="equipmentForm.imageUrl" hlmInput class="mt-1" />
            </div>

            <div class="flex gap-2">
              <button
                type="submit"
                [disabled]="
                  !equipmentForm().valid() || mutation.isPending() || createMutation.isPending()
                "
                class="rounded-md bg-blue-600 px-4 py-2 text-sm font-medium text-white hover:bg-blue-700 disabled:opacity-50"
              >
                {{ t('common.save') }}
              </button>
              <a
                [routerLink]="isCreateMode() ? ['/compendium/equipment'] : ['..']"
                class="rounded-md border border-gray-300 px-4 py-2 text-sm font-medium text-gray-700 hover:bg-gray-50 dark:border-gray-600 dark:text-gray-300 dark:hover:bg-gray-800"
              >
                {{ t('common.cancel') }}
              </a>
            </div>
          </form>
        }
      </app-page-layout>
    </ng-container>
  `,
})
export class EquipmentEdit {
  private api = inject(CompendiumApiClient);
  private router = inject(Router);
  private route = inject(ActivatedRoute);
  private queryClient = inject(QueryClient);
  private slugify = new SlugifyPipe();
  private params = toSignal(this.route.paramMap);

  private id = computed(() => Number(this.params()?.get('id')));
  isCreateMode = computed(() => !this.params()?.get('id'));

  categories = EQUIPMENT_CATEGORIES;

  model = signal({
    name: '',
    displayName: '',
    description: '',
    category: EquipmentCategoryFreeWeights as string,
    imageUrl: '',
  });
  equipmentForm = form(this.model, (f) => {
    required(f.name);
    required(f.displayName);
    required(f.category);
  });

  equipmentQuery = injectQuery(() => ({
    queryKey: equipmentKeys.detail(this.id()),
    queryFn: () => this.api.fetchEquipmentItem(this.id()),
    enabled: !!this.id() && !this.isCreateMode(),
  }));

  permissionsQuery = injectQuery(() => ({
    queryKey: equipmentKeys.permissions(this.id()),
    queryFn: () => this.api.fetchEquipmentPermissions(this.id()),
    enabled: !!this.id() && !this.isCreateMode(),
  }));

  mutation = injectMutation(() => ({
    mutationFn: (data: Parameters<typeof this.api.updateEquipment>[1]) =>
      this.api.updateEquipment(this.id(), data),
    onSuccess: () => {
      this.queryClient.invalidateQueries({ queryKey: equipmentKeys.all() });
      this.router.navigate(['..'], { relativeTo: this.route });
    },
  }));

  createMutation = injectMutation(() => ({
    mutationFn: (data: Parameters<typeof this.api.createEquipment>[0]) =>
      this.api.createEquipment(data),
    onSuccess: (result) => {
      this.queryClient.invalidateQueries({ queryKey: equipmentKeys.all() });
      this.router.navigate([
        '/compendium/equipment',
        result.id,
        this.slugify.transform(result.displayName),
      ]);
    },
  }));

  constructor() {
    effect(() => {
      const perms = this.permissionsQuery.data();
      if (perms && !perms.permissions.includes('MODIFY')) {
        this.router.navigate(['..'], { relativeTo: this.route });
      }
    });

    effect(() => {
      const data = this.equipmentQuery.data();
      if (data) {
        this.model.set({
          name: data.name,
          displayName: data.displayName,
          description: data.description,
          category: data.category,
          imageUrl: data.imageUrl ?? '',
        });
      }
    });
  }

  onSubmit() {
    if (this.equipmentForm().valid()) {
      const val = this.model();
      const payload = {
        ...(this.isCreateMode() ? {} : this.equipmentQuery.data()!),
        name: val.name,
        displayName: val.displayName,
        description: val.description,
        category: val.category as EquipmentCategory,
        imageUrl: val.imageUrl || undefined,
      };
      if (this.isCreateMode()) {
        this.createMutation.mutate(payload);
      } else {
        this.mutation.mutate(payload);
      }
    }
  }
}
