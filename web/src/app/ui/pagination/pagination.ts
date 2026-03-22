import { Component, effect, inject, input, linkedSignal } from '@angular/core';
import { ActivatedRoute, Router } from '@angular/router';
import { PaginatedResponse } from '$core/api-clients/paginated-response';
import { HlmNumberedPagination } from '@spartan-ng/helm/pagination';
import { TranslocoDirective } from '@jsverse/transloco';

@Component({
  selector: 'app-pagination',
  imports: [HlmNumberedPagination, TranslocoDirective],
  template: `
    <ng-container *transloco="let t">
      @if (page(); as p) {
        @if (p.total === 0) {
          <p class="text-sm text-gray-500 dark:text-gray-400">{{ t(emptyLabel()) }}</p>
        } @else {
          <hlm-numbered-pagination
            [(currentPage)]="currentPage"
            [(itemsPerPage)]="pageSize"
            [totalItems]="p.total"
          />
        }
      }
    </ng-container>
  `,
})
export class Pagination {
  private router = inject(Router);
  private route = inject(ActivatedRoute);

  page = input.required<PaginatedResponse<unknown>>();
  emptyLabel = input('ui.pagination.noResults');

  currentPage = linkedSignal(() => {
    const p = this.page();
    return p.limit > 0 ? Math.floor(p.offset / p.limit) + 1 : 1;
  });

  pageSize = linkedSignal(() => this.page().limit);

  constructor() {
    effect(() => {
      const newPage = this.currentPage();
      const newSize = this.pageSize();
      const p = this.page();
      const expectedPage = p.limit > 0 ? Math.floor(p.offset / p.limit) + 1 : 1;

      if (newPage !== expectedPage || newSize !== p.limit) {
        const offset = (newPage - 1) * newSize;
        this.router.navigate([], {
          relativeTo: this.route,
          queryParams: {
            offset: offset || null,
            limit: newSize !== 50 ? newSize : null,
          },
          queryParamsHandling: 'merge',
        });
      }
    });
  }
}
