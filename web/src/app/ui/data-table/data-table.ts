import {
  Component,
  ElementRef,
  HostListener,
  Injector,
  afterNextRender,
  inject,
  input,
  signal,
  viewChild,
} from '@angular/core';
import { FormsModule } from '@angular/forms';
import { NgClass } from '@angular/common';
import { ActivatedRoute, Router } from '@angular/router';
import { toSignal, takeUntilDestroyed } from '@angular/core/rxjs-interop';
import { Subject, debounceTime } from 'rxjs';

export interface DataTableColumn {
  label: string;
  options?: string[];
  filterParam?: string;
  searchParam?: string;
}

@Component({
  selector: 'app-data-table',
  imports: [FormsModule, NgClass],
  template: `
    <div
      class="rounded-lg border border-gray-200 dark:border-gray-800"
      [class.overflow-x-auto]="!activeFilter()"
      [class.overflow-visible]="!!activeFilter()"
    >
      <table class="min-w-full divide-y divide-gray-200 dark:divide-gray-800">
        <thead class="bg-gray-50 dark:bg-gray-900">
          <tr>
            @for (col of columns(); track col.label) {
              <th
                class="px-4 py-2 text-left text-xs font-medium tracking-wider text-gray-500 uppercase dark:text-gray-400"
                [class.relative]="!!col.filterParam || !!col.searchParam"
              >
                @if (col.filterParam) {
                  @if (activeFilter() === col.label) {
                    <!-- eslint-disable-next-line @angular-eslint/template/click-events-have-key-events, @angular-eslint/template/interactive-supports-focus -->
                    <div class="relative" (click)="$event.stopPropagation()">
                      <input
                        type="text"
                        [placeholder]="'Filter ' + col.label.toLowerCase() + '...'"
                        [ngModel]="searchTerm()"
                        (ngModelChange)="searchTerm.set($event)"
                        (keydown.escape)="activeFilter.set(null)"
                        class="w-full min-w-32 rounded border border-blue-400 bg-white px-2 py-0.5 text-xs font-normal normal-case focus:outline-none dark:border-blue-600 dark:bg-gray-800 dark:text-gray-100"
                        #filterInput
                      />
                      <div
                        class="absolute left-0 z-50 mt-1 max-h-48 min-w-full overflow-y-auto rounded-md border border-gray-200 bg-white shadow-lg dark:border-gray-700 dark:bg-gray-800"
                      >
                        <button
                          type="button"
                          (click)="selectOption(col, '')"
                          class="w-full whitespace-nowrap px-3 py-1.5 text-left text-xs font-normal normal-case text-gray-500 hover:bg-gray-100 dark:text-gray-400 dark:hover:bg-gray-700"
                        >
                          All
                        </button>
                        @for (opt of getFilteredOptions(col); track opt) {
                          <button
                            type="button"
                            (click)="selectOption(col, opt)"
                            class="w-full whitespace-nowrap px-3 py-1.5 text-left text-xs font-normal normal-case hover:bg-gray-100 dark:text-gray-100 dark:hover:bg-gray-700"
                            [ngClass]="{
                              'bg-blue-50 dark:bg-blue-900': getParamValue(col.filterParam) === opt,
                            }"
                          >
                            {{ opt }}
                          </button>
                        }
                      </div>
                    </div>
                  } @else {
                    <button
                      type="button"
                      (click)="openFilter(col, $event)"
                      class="flex cursor-pointer items-center gap-1 hover:text-gray-700 dark:hover:text-gray-200"
                    >
                      {{ col.label }}
                      @if (getParamValue(col.filterParam); as val) {
                        <span
                          class="rounded bg-blue-100 px-1.5 py-0.5 text-[10px] font-semibold text-blue-700 dark:bg-blue-900 dark:text-blue-300"
                        >
                          {{ val }}
                        </span>
                      }
                      <svg
                        class="h-3 w-3"
                        [class.opacity-40]="!getParamValue(col.filterParam)"
                        viewBox="0 0 20 20"
                        fill="currentColor"
                      >
                        <path
                          fill-rule="evenodd"
                          d="M2.628 1.601C5.028 1.206 7.49 1 10 1s4.973.206 7.372.601a.75.75 0 01.628.74v2.288a2.25 2.25 0 01-.659 1.59l-4.682 4.683a2.25 2.25 0 00-.659 1.59v3.037c0 .684-.31 1.33-.844 1.757l-1.937 1.55A.75.75 0 018 18.25v-5.757a2.25 2.25 0 00-.659-1.591L2.659 6.22A2.25 2.25 0 012 4.629V2.34a.75.75 0 01.628-.74z"
                          clip-rule="evenodd"
                        />
                      </svg>
                    </button>
                  }
                } @else if (col.searchParam) {
                  @if (activeFilter() === col.label) {
                    <!-- eslint-disable-next-line @angular-eslint/template/click-events-have-key-events, @angular-eslint/template/interactive-supports-focus -->
                    <div (click)="$event.stopPropagation()">
                      <input
                        type="text"
                        [placeholder]="'Search ' + col.label.toLowerCase() + '...'"
                        [ngModel]="searchTerm()"
                        (ngModelChange)="onSearchInput(col, $event)"
                        (keydown.escape)="activeFilter.set(null)"
                        class="w-full min-w-32 rounded border border-blue-400 bg-white px-2 py-0.5 text-xs font-normal normal-case focus:outline-none dark:border-blue-600 dark:bg-gray-800 dark:text-gray-100"
                        #filterInput
                      />
                    </div>
                  } @else {
                    <button
                      type="button"
                      (click)="openFilter(col, $event)"
                      class="flex cursor-pointer items-center gap-1 hover:text-gray-700 dark:hover:text-gray-200"
                    >
                      {{ col.label }}
                      @if (getParamValue(col.searchParam); as val) {
                        <span
                          class="rounded bg-blue-100 px-1.5 py-0.5 text-[10px] font-semibold text-blue-700 dark:bg-blue-900 dark:text-blue-300"
                        >
                          {{ val }}
                        </span>
                      }
                      <svg
                        class="h-3 w-3"
                        [class.opacity-40]="!getParamValue(col.searchParam)"
                        viewBox="0 0 20 20"
                        fill="currentColor"
                      >
                        <path
                          fill-rule="evenodd"
                          d="M9 3.5a5.5 5.5 0 100 11 5.5 5.5 0 000-11zM2 9a7 7 0 1112.452 4.391l3.328 3.329a.75.75 0 11-1.06 1.06l-3.329-3.328A7 7 0 012 9z"
                          clip-rule="evenodd"
                        />
                      </svg>
                    </button>
                  }
                } @else {
                  {{ col.label }}
                }
              </th>
            }
          </tr>
        </thead>
        <tbody
          class="divide-y divide-gray-200 bg-white transition-opacity duration-200 dark:divide-gray-800 dark:bg-gray-950"
          [class.opacity-50]="stale()"
        >
          <ng-content />
        </tbody>
      </table>
    </div>
  `,
})
export class DataTable {
  private injector = inject(Injector);
  private router = inject(Router);
  private route = inject(ActivatedRoute);
  private queryParams = toSignal(this.route.queryParamMap);

  columns = input.required<DataTableColumn[]>();
  stale = input(false);
  activeFilter = signal<string | null>(null);
  searchTerm = signal('');

  filterInput = viewChild<ElementRef<HTMLInputElement>>('filterInput');

  private search$ = new Subject<{ param: string; value: string }>();

  constructor() {
    this.search$.pipe(debounceTime(300), takeUntilDestroyed()).subscribe(({ param, value }) => {
      this.router.navigate([], {
        relativeTo: this.route,
        queryParams: { [param]: value || null, offset: null },
        queryParamsHandling: 'merge',
        replaceUrl: true,
      });
    });
  }

  getParamValue(param: string): string {
    return this.queryParams()?.get(param) ?? '';
  }

  openFilter(col: DataTableColumn, event: Event) {
    event.stopPropagation();
    this.searchTerm.set(col.searchParam ? this.getParamValue(col.searchParam) : '');
    this.activeFilter.set(col.label);
    afterNextRender(
      () => {
        this.filterInput()?.nativeElement.focus();
      },
      { injector: this.injector },
    );
  }

  onSearchInput(col: DataTableColumn, value: string) {
    this.searchTerm.set(value);
    this.search$.next({ param: col.searchParam!, value });
  }

  selectOption(col: DataTableColumn, value: string) {
    this.router.navigate([], {
      relativeTo: this.route,
      queryParams: { [col.filterParam!]: value || null, offset: null },
      queryParamsHandling: 'merge',
    });
    this.activeFilter.set(null);
  }

  getFilteredOptions(col: DataTableColumn): string[] {
    const term = this.searchTerm().toLowerCase();
    if (!term) return col.options ?? [];
    return (col.options ?? []).filter((o) => o.toLowerCase().includes(term));
  }

  @HostListener('document:click')
  onDocumentClick() {
    if (this.activeFilter()) {
      this.activeFilter.set(null);
    }
  }
}
