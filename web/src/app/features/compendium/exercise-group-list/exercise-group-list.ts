import { Component, inject, signal, computed, effect, untracked } from '@angular/core';
import { injectQuery } from '@tanstack/angular-query-experimental';
import { CompendiumApiClient } from '$core/api-clients/compendium-api-client';
import { ExerciseGroupListItem } from '$ui/compendium/exercise-group-list-item/exercise-group-list-item';
import { SearchInput } from '$ui/inputs/search-input/search-input';
import { DataTable } from '$ui/data-table/data-table';
import { Pagination } from '$ui/pagination/pagination';
import { PageLayout } from '../../../layout/page-layout';

@Component({
  selector: 'app-exercise-group-list',
  imports: [SearchInput, ExerciseGroupListItem, DataTable, Pagination, PageLayout],
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
        <app-pagination [page]="page" [(offset)]="offset" emptyLabel="No groups found" />
      }
    </app-page-layout>
  `,
})
export class ExerciseGroupList {
  private api = inject(CompendiumApiClient);

  q = signal('');
  offset = signal(0);

  constructor() {
    effect(() => {
      this.q();
      untracked(() => this.offset.set(0));
    });
  }

  filters = computed(() => ({
    q: this.q() || undefined,
    offset: this.offset() || undefined,
  }));

  groupsQuery = injectQuery(() => ({
    queryKey: ['exercise-groups', this.filters()],
    queryFn: () => this.api.fetchExerciseGroups(this.filters()),
  }));
}
