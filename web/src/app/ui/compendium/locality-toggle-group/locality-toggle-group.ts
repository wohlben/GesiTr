import { Component, inject, output, signal } from '@angular/core';
import { injectQuery } from '@tanstack/angular-query-experimental';
import { CompendiumApiClient } from '$core/api-clients/compendium-api-client';
import { localityKeys } from '$core/query-keys';
import { TranslocoDirective } from '@jsverse/transloco';

@Component({
  selector: 'app-locality-toggle-group',
  imports: [TranslocoDirective],
  template: `
    <ng-container *transloco="let t">
      @if (localitiesQuery.data(); as page) {
        <div class="flex overflow-hidden rounded-md border border-gray-300 dark:border-gray-600">
          <button
            type="button"
            (click)="select(null)"
            class="px-3 py-2 text-sm font-medium transition-colors"
            [class]="
              selected() === null
                ? 'bg-blue-600 text-white'
                : 'bg-white text-gray-700 hover:bg-gray-50 dark:bg-gray-800 dark:text-gray-300 dark:hover:bg-gray-700'
            "
          >
            {{ t('common.all') }}
          </button>
          @for (locality of page.items; track locality.id) {
            <button
              type="button"
              (click)="select(locality.id)"
              class="border-l border-gray-300 px-3 py-2 text-sm font-medium transition-colors dark:border-gray-600"
              [class]="
                selected() === locality.id
                  ? 'bg-blue-600 text-white'
                  : 'bg-white text-gray-700 hover:bg-gray-50 dark:bg-gray-800 dark:text-gray-300 dark:hover:bg-gray-700'
              "
            >
              {{ locality.name }}
            </button>
          }
        </div>
      }
    </ng-container>
  `,
})
export class LocalityToggleGroup {
  private api = inject(CompendiumApiClient);

  selected = signal<number | null>(null);
  selectedChange = output<number | null>();

  localitiesQuery = injectQuery(() => ({
    queryKey: localityKeys.list({ owner: 'me', limit: 100 }),
    queryFn: () => this.api.fetchLocalities({ owner: 'me', limit: 100 }),
  }));

  select(id: number | null) {
    this.selected.set(id);
    this.selectedChange.emit(id);
  }
}
