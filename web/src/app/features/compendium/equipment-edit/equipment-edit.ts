import { Component, inject, computed, effect } from '@angular/core';
import { ActivatedRoute, Router, RouterLink } from '@angular/router';
import { toSignal } from '@angular/core/rxjs-interop';
import { ReactiveFormsModule, FormGroup, FormControl, Validators } from '@angular/forms';
import {
  injectQuery,
  injectMutation,
  injectQueryClient,
} from '@tanstack/angular-query-experimental';
import { CompendiumApiClient } from '$core/api-clients/compendium-api-client';
import { equipmentKeys } from '$core/query-keys';
import { SlugifyPipe } from '$ui/pipes/slugify';
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
  imports: [PageLayout, ReactiveFormsModule, RouterLink],
  template: `
    <app-page-layout
      [header]="isCreateMode() ? 'New Equipment' : 'Edit Equipment'"
      [isPending]="!isCreateMode() && equipmentQuery.isPending()"
      [errorMessage]="
        !isCreateMode() && equipmentQuery.isError() ? equipmentQuery.error().message : undefined
      "
    >
      @if (isCreateMode() || equipmentQuery.data()) {
        <form [formGroup]="form" (ngSubmit)="onSubmit()" class="space-y-4">
          <div>
            <label for="name" class="block text-sm font-medium text-gray-700 dark:text-gray-300"
              >Name *</label
            >
            <input
              id="name"
              formControlName="name"
              class="mt-1 block w-full rounded-md border border-gray-300 px-3 py-2 shadow-sm focus:border-blue-500 focus:outline-none focus:ring-1 focus:ring-blue-500 dark:border-gray-600 dark:bg-gray-800 dark:text-gray-100"
            />
          </div>

          <div>
            <label
              for="displayName"
              class="block text-sm font-medium text-gray-700 dark:text-gray-300"
              >Display Name *</label
            >
            <input
              id="displayName"
              formControlName="displayName"
              class="mt-1 block w-full rounded-md border border-gray-300 px-3 py-2 shadow-sm focus:border-blue-500 focus:outline-none focus:ring-1 focus:ring-blue-500 dark:border-gray-600 dark:bg-gray-800 dark:text-gray-100"
            />
          </div>

          <div>
            <label
              for="description"
              class="block text-sm font-medium text-gray-700 dark:text-gray-300"
              >Description</label
            >
            <textarea
              id="description"
              formControlName="description"
              rows="4"
              class="mt-1 block w-full rounded-md border border-gray-300 px-3 py-2 shadow-sm focus:border-blue-500 focus:outline-none focus:ring-1 focus:ring-blue-500 dark:border-gray-600 dark:bg-gray-800 dark:text-gray-100"
            ></textarea>
          </div>

          <div>
            <label for="category" class="block text-sm font-medium text-gray-700 dark:text-gray-300"
              >Category *</label
            >
            <select
              id="category"
              formControlName="category"
              class="mt-1 block w-full rounded-md border border-gray-300 px-3 py-2 shadow-sm focus:border-blue-500 focus:outline-none focus:ring-1 focus:ring-blue-500 dark:border-gray-600 dark:bg-gray-800 dark:text-gray-100"
            >
              @for (cat of categories; track cat) {
                <option [value]="cat">{{ cat }}</option>
              }
            </select>
          </div>

          <div>
            <label for="imageUrl" class="block text-sm font-medium text-gray-700 dark:text-gray-300"
              >Image URL</label
            >
            <input
              id="imageUrl"
              formControlName="imageUrl"
              class="mt-1 block w-full rounded-md border border-gray-300 px-3 py-2 shadow-sm focus:border-blue-500 focus:outline-none focus:ring-1 focus:ring-blue-500 dark:border-gray-600 dark:bg-gray-800 dark:text-gray-100"
            />
          </div>

          <div class="flex gap-2">
            <button
              type="submit"
              [disabled]="form.invalid || mutation.isPending() || createMutation.isPending()"
              class="rounded-md bg-blue-600 px-4 py-2 text-sm font-medium text-white hover:bg-blue-700 disabled:opacity-50"
            >
              Save
            </button>
            <a
              [routerLink]="isCreateMode() ? ['/compendium/equipment'] : ['..']"
              class="rounded-md border border-gray-300 px-4 py-2 text-sm font-medium text-gray-700 hover:bg-gray-50 dark:border-gray-600 dark:text-gray-300 dark:hover:bg-gray-800"
            >
              Cancel
            </a>
          </div>
        </form>
      }
    </app-page-layout>
  `,
})
export class EquipmentEdit {
  private api = inject(CompendiumApiClient);
  private router = inject(Router);
  private route = inject(ActivatedRoute);
  private queryClient = injectQueryClient();
  private slugify = new SlugifyPipe();
  private params = toSignal(this.route.paramMap);

  private id = computed(() => Number(this.params()?.get('id')));
  isCreateMode = computed(() => !this.params()?.get('id'));

  categories = EQUIPMENT_CATEGORIES;

  form = new FormGroup({
    name: new FormControl('', { nonNullable: true, validators: [Validators.required] }),
    displayName: new FormControl('', { nonNullable: true, validators: [Validators.required] }),
    description: new FormControl('', { nonNullable: true }),
    category: new FormControl<EquipmentCategory>(EquipmentCategoryFreeWeights, {
      nonNullable: true,
      validators: [Validators.required],
    }),
    imageUrl: new FormControl('', { nonNullable: true }),
  });

  equipmentQuery = injectQuery(() => ({
    queryKey: equipmentKeys.detail(this.id()),
    queryFn: () => this.api.fetchEquipmentItem(this.id()),
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
      const data = this.equipmentQuery.data();
      if (data) {
        this.form.patchValue({
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
    if (this.form.valid) {
      const val = this.form.getRawValue();
      const payload = {
        ...(this.isCreateMode() ? {} : this.equipmentQuery.data()!),
        name: val.name,

        displayName: val.displayName,
        description: val.description,
        category: val.category,
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
