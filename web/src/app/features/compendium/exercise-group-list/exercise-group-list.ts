import { Component, inject, signal, computed, effect, untracked } from '@angular/core';
import { injectQuery } from '@tanstack/angular-query-experimental';
import { CompendiumApiClient } from '$core/api-clients/compendium-api-client';
import { ExerciseGroupListItem } from '$ui/compendium/exercise-group-list-item/exercise-group-list-item';
import { SearchInput } from '$ui/inputs/search-input/search-input';
import { DataTable } from '$ui/data-table/data-table';
import { PageLayout } from '../../../layout/page-layout';

@Component({
  selector: 'app-exercise-group-list',
  imports: [SearchInput, ExerciseGroupListItem, DataTable, PageLayout],
  template: `
    <app-page-layout
      header="Exercise Groups"
      [isPending]="groupsQuery.isPending()"
      [errorMessage]="groupsQuery.isError() ? groupsQuery.error().message : undefined">

      <div filters class="flex flex-wrap gap-3">
        <app-search-input placeholder="Search exercise groups..." [(value)]="q" />
      </div>

      @if (groupsQuery.data(); as page) {
        <app-data-table [columns]="['Name', 'Description']">
          @for (group of page.items; track group.id) {
            <tr app-exercise-group-list-item [group]="group"></tr>
          }
        </app-data-table>
        <div class="flex items-center justify-between text-sm text-gray-500 dark:text-gray-400">
          <p>{{ page.total === 0 ? 'No groups found' : 'Showing ' + (page.offset + 1) + '–' + (page.offset + page.items.length) + ' of ' + page.total + ' groups' }}</p>
          <div class="flex gap-2">
            <button
              class="rounded border border-gray-300 px-3 py-1 disabled:opacity-50 dark:border-gray-600"
              [disabled]="page.offset === 0"
              (click)="prevPage()">Previous</button>
            <button
              class="rounded border border-gray-300 px-3 py-1 disabled:opacity-50 dark:border-gray-600"
              [disabled]="page.offset + page.limit >= page.total"
              (click)="nextPage()">Next</button>
          </div>
        </div>
      }
    </app-page-layout>
  `,
})
export class ExerciseGroupList {
  private api = inject(CompendiumApiClient);

  q = signal('');
  offset = signal(0);

  private resetOffset = effect(() => {
    this.q();
    untracked(() => this.offset.set(0));
  });

  filters = computed(() => ({
    q: this.q() || undefined,
    offset: this.offset() || undefined,
  }));

  groupsQuery = injectQuery(() => ({
    queryKey: ['exercise-groups', this.filters()],
    queryFn: () => this.api.fetchExerciseGroups(this.filters()),
  }));

  prevPage() {
    const page = this.groupsQuery.data();
    if (page) this.offset.set(Math.max(0, page.offset - page.limit));
  }

  nextPage() {
    const page = this.groupsQuery.data();
    if (page) this.offset.set(page.offset + page.limit);
  }
}
