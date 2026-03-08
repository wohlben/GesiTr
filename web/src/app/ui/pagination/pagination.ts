import { Component, computed, inject, input } from '@angular/core';
import { ActivatedRoute, Router } from '@angular/router';
import { PaginatedResponse } from '$core/api-clients/paginated-response';

@Component({
  selector: 'app-pagination',
  template: `
    @if (page(); as p) {
      <div class="flex items-center justify-between text-sm text-gray-500 dark:text-gray-400">
        <p>{{ p.total === 0 ? emptyLabel() : 'Showing ' + (p.offset + 1) + '–' + (p.offset + p.items.length) + ' of ' + p.total }}</p>
        @if (totalPages() > 1) {
          <div class="flex gap-2">
            <button
              class="rounded border border-gray-300 px-3 py-1 disabled:opacity-50 dark:border-gray-600"
              [disabled]="!hasPrev()"
              (click)="prev()">Previous</button>
            <button
              class="rounded border border-gray-300 px-3 py-1 disabled:opacity-50 dark:border-gray-600"
              [disabled]="!hasNext()"
              (click)="next()">Next</button>
          </div>
        }
      </div>
    }
  `,
})
export class Pagination {
  private router = inject(Router);
  private route = inject(ActivatedRoute);

  page = input.required<PaginatedResponse<unknown>>();
  emptyLabel = input('No results found');

  totalPages = computed(() => {
    const p = this.page();
    return Math.ceil(p.total / p.limit);
  });

  hasPrev = computed(() => this.page().offset > 0);
  hasNext = computed(() => {
    const p = this.page();
    return p.offset + p.limit < p.total;
  });

  prev() {
    const p = this.page();
    const newOffset = Math.max(0, p.offset - p.limit);
    this.navigate(newOffset);
  }

  next() {
    const p = this.page();
    this.navigate(p.offset + p.limit);
  }

  private navigate(offset: number) {
    this.router.navigate([], {
      relativeTo: this.route,
      queryParams: { offset: offset || null },
      queryParamsHandling: 'merge',
    });
  }
}
