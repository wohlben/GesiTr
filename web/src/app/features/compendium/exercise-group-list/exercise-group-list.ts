import { Component, inject, signal, computed } from '@angular/core';
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

      @if (groupsQuery.data(); as groups) {
        <app-data-table [columns]="['Name', 'Description']">
          @for (group of groups; track group.id) {
            <tr app-exercise-group-list-item [group]="group"></tr>
          }
        </app-data-table>
        <p class="text-sm text-gray-500 dark:text-gray-400">{{ groups.length }} groups</p>
      }
    </app-page-layout>
  `,
})
export class ExerciseGroupList {
  private api = inject(CompendiumApiClient);

  q = signal('');

  filters = computed(() => ({
    q: this.q() || undefined,
  }));

  groupsQuery = injectQuery(() => ({
    queryKey: ['exercise-groups', this.filters()],
    queryFn: () => this.api.fetchExerciseGroups(this.filters()),
  }));
}
