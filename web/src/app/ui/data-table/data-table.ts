import {
  Component,
  DestroyRef,
  ElementRef,
  HostListener,
  Injector,
  afterNextRender,
  computed,
  effect,
  inject,
  input,
  output,
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
  hideable?: boolean;
  defaultHidden?: boolean;
}

@Component({
  selector: 'app-data-table',
  imports: [FormsModule, NgClass],
  template: `
    <div
      [id]="tableId"
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
            @if (hideableColumns().length > 0) {
              <th class="w-8 bg-gray-50 px-2 py-2 dark:bg-gray-900">
                <button
                  type="button"
                  (click)="showColumnSettings.set(true); $event.stopPropagation()"
                  class="text-gray-400 hover:text-gray-600 dark:hover:text-gray-200"
                  aria-label="Column settings"
                >
                  <svg class="h-4 w-4" viewBox="0 0 20 20" fill="currentColor">
                    <path
                      fill-rule="evenodd"
                      d="M8.34 1.804A1 1 0 0 1 9.32 1h1.36a1 1 0 0 1 .98.804l.295 1.473c.497.144.971.342 1.416.587l1.25-.834a1 1 0 0 1 1.262.125l.962.962a1 1 0 0 1 .125 1.262l-.834 1.25c.245.445.443.919.587 1.416l1.473.294a1 1 0 0 1 .804.98v1.361a1 1 0 0 1-.804.98l-1.473.295a6.95 6.95 0 0 1-.587 1.416l.834 1.25a1 1 0 0 1-.125 1.262l-.962.962a1 1 0 0 1-1.262.125l-1.25-.834a6.953 6.953 0 0 1-1.416.587l-.294 1.473a1 1 0 0 1-.98.804H9.32a1 1 0 0 1-.98-.804l-.295-1.473a6.957 6.957 0 0 1-1.416-.587l-1.25.834a1 1 0 0 1-1.262-.125l-.962-.962a1 1 0 0 1-.125-1.262l.834-1.25a6.957 6.957 0 0 1-.587-1.416l-1.473-.294A1 1 0 0 1 1 10.68V9.32a1 1 0 0 1 .804-.98l1.473-.295c.144-.497.342-.971.587-1.416l-.834-1.25a1 1 0 0 1 .125-1.262l.962-.962A1 1 0 0 1 5.38 3.03l1.25.834a6.957 6.957 0 0 1 1.416-.587l.294-1.473ZM13 10a3 3 0 1 1-6 0 3 3 0 0 1 6 0Z"
                      clip-rule="evenodd"
                    />
                  </svg>
                </button>
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
    @if (showColumnSettings()) {
      <!-- eslint-disable-next-line @angular-eslint/template/click-events-have-key-events, @angular-eslint/template/interactive-supports-focus -->
      <div
        class="fixed inset-0 z-50 flex items-center justify-center bg-black/50"
        (click)="showColumnSettings.set(false)"
      >
        <!-- eslint-disable-next-line @angular-eslint/template/click-events-have-key-events, @angular-eslint/template/interactive-supports-focus -->
        <div
          class="min-w-48 rounded-lg bg-white p-4 shadow-xl dark:bg-gray-800"
          (click)="$event.stopPropagation()"
          role="dialog"
          aria-label="Column visibility settings"
        >
          <div class="mb-3 flex items-center justify-between gap-4">
            <h3 class="text-sm font-medium text-gray-900 dark:text-gray-100">Columns</h3>
            <button
              type="button"
              (click)="showColumnSettings.set(false)"
              class="text-gray-400 hover:text-gray-600 dark:hover:text-gray-200"
              aria-label="Close"
            >
              <svg class="h-4 w-4" viewBox="0 0 20 20" fill="currentColor">
                <path
                  d="M6.28 5.22a.75.75 0 0 0-1.06 1.06L8.94 10l-3.72 3.72a.75.75 0 1 0 1.06 1.06L10 11.06l3.72 3.72a.75.75 0 1 0 1.06-1.06L11.06 10l3.72-3.72a.75.75 0 0 0-1.06-1.06L10 8.94 6.28 5.22Z"
                />
              </svg>
            </button>
          </div>
          @for (col of hideableColumns(); track col.label) {
            <label class="flex items-center gap-2 py-1 text-sm text-gray-700 dark:text-gray-300">
              <input
                type="checkbox"
                [checked]="!hiddenColumns().has(col.label)"
                (change)="toggleColumn(col.label)"
                class="rounded"
              />
              {{ col.label }}
            </label>
          }
        </div>
      </div>
    }
  `,
})
export class DataTable {
  private injector = inject(Injector);
  private destroyRef = inject(DestroyRef);
  private router = inject(Router);
  private route = inject(ActivatedRoute);
  private queryParams = toSignal(this.route.queryParamMap);

  columns = input.required<DataTableColumn[]>();
  stale = input(false);
  initialHiddenColumns = input<string[]>();

  hiddenColumnsChange = output<string[]>();

  activeFilter = signal<string | null>(null);
  searchTerm = signal('');
  hiddenColumns = signal<Set<string>>(new Set());
  showColumnSettings = signal(false);

  filterInput = viewChild<ElementRef<HTMLInputElement>>('filterInput');

  tableId = 'dt-' + Math.random().toString(36).slice(2, 9);
  private styleEl: HTMLStyleElement | null = null;

  hideableColumns = computed(() => this.columns().filter((col) => col.hideable !== false));

  columnHideStyles = computed(() => {
    const hidden = this.hiddenColumns();
    if (hidden.size === 0) return '';
    const rules: string[] = [];
    this.columns().forEach((col, i) => {
      if (hidden.has(col.label)) {
        const nth = i + 1;
        rules.push(
          `#${this.tableId} th:nth-child(${nth}), #${this.tableId} td:nth-child(${nth}) { display: none; }`,
        );
      }
    });
    return rules.join('\n');
  });

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

    // Keep style element in sync with column hide rules
    effect(() => {
      const styles = this.columnHideStyles();
      if (this.styleEl) {
        this.styleEl.textContent = styles;
      }
    });

    afterNextRender(
      () => {
        // Create style element for column hiding CSS
        this.styleEl = document.createElement('style');
        this.styleEl.textContent = this.columnHideStyles();
        document.head.appendChild(this.styleEl);
        this.destroyRef.onDestroy(() => this.styleEl?.remove());

        // Seed hidden columns from input or defaultHidden
        const initial = this.initialHiddenColumns();
        if (initial !== undefined) {
          this.hiddenColumns.set(new Set(initial));
        } else {
          const defaults = this.columns()
            .filter((col) => col.defaultHidden)
            .map((col) => col.label);
          if (defaults.length > 0) {
            this.hiddenColumns.set(new Set(defaults));
          }
        }
      },
      { injector: this.injector },
    );
  }

  toggleColumn(label: string) {
    this.hiddenColumns.update((set) => {
      const next = new Set(set);
      next.has(label) ? next.delete(label) : next.add(label);
      return next;
    });
    this.hiddenColumnsChange.emit([...this.hiddenColumns()]);
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

  @HostListener('document:keydown.escape')
  onEscapeKey() {
    if (this.showColumnSettings()) {
      this.showColumnSettings.set(false);
    }
  }

  @HostListener('document:click')
  onDocumentClick() {
    if (this.activeFilter()) {
      this.activeFilter.set(null);
    }
  }
}
