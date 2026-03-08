import { render, screen } from '@testing-library/angular';
import userEvent from '@testing-library/user-event';
import { TestBed } from '@angular/core/testing';
import { provideRouter, Router } from '@angular/router';
import { provideLocationMocks } from '@angular/common/testing';
import { DataTable, DataTableColumn } from './data-table';

describe('DataTable', () => {
  const staticColumns: DataTableColumn[] = [
    { label: 'Name', hideable: false },
    { label: 'Value', hideable: false },
  ];

  const filterableColumns: DataTableColumn[] = [
    { label: 'Name', searchParam: 'q' },
    { label: 'Type', filterParam: 'type', options: ['Alpha', 'Beta', 'Gamma'] },
    { label: 'Status' },
  ];

  const hideableColumns: DataTableColumn[] = [
    { label: 'Name' },
    { label: 'Value' },
    { label: 'Fixed', hideable: false },
  ];

  const template = `
    <app-data-table [columns]="columns" [stale]="stale">
      <tr><td>row</td></tr>
    </app-data-table>
  `;

  const renderTable = (columns: DataTableColumn[], extras: Record<string, unknown> = {}) =>
    render(template, {
      imports: [DataTable],
      providers: [provideRouter([]), provideLocationMocks()],
      componentProperties: { columns, stale: false, ...extras },
    });

  it('renders column headers', async () => {
    await renderTable(staticColumns);
    expect(screen.getByText('Name')).toBeTruthy();
    expect(screen.getByText('Value')).toBeTruthy();
  });

  it('renders static columns without filter buttons', async () => {
    await renderTable(staticColumns);
    expect(screen.queryByRole('button', { name: /name/i })).toBeNull();
    expect(screen.queryByRole('button', { name: /value/i })).toBeNull();
  });

  it('shows clickable button for filterable columns', async () => {
    await renderTable(filterableColumns);
    expect(screen.getByRole('button', { name: /type/i })).toBeTruthy();
  });

  it('shows clickable button for searchable columns', async () => {
    await renderTable(filterableColumns);
    expect(screen.getByRole('button', { name: /name/i })).toBeTruthy();
  });

  it('does not show button for static columns', async () => {
    await renderTable(filterableColumns);
    expect(screen.queryByRole('button', { name: /status/i })).toBeNull();
  });

  it('opens dropdown when clicking a filterable header', async () => {
    const user = userEvent.setup();
    await renderTable(filterableColumns);

    await user.click(screen.getByRole('button', { name: /type/i }));

    expect(screen.getByText('All')).toBeTruthy();
    expect(screen.getByText('Alpha')).toBeTruthy();
    expect(screen.getByText('Beta')).toBeTruthy();
    expect(screen.getByText('Gamma')).toBeTruthy();
  });

  it('filters dropdown options when typing', async () => {
    const user = userEvent.setup();
    await renderTable(filterableColumns);

    await user.click(screen.getByRole('button', { name: /type/i }));
    await user.type(screen.getByPlaceholderText('Filter type...'), 'alp');

    expect(screen.getByText('Alpha')).toBeTruthy();
    expect(screen.queryByText('Beta')).toBeNull();
    expect(screen.queryByText('Gamma')).toBeNull();
  });

  it('navigates with filter param when selecting an option', async () => {
    const user = userEvent.setup();
    await renderTable(filterableColumns);

    const router = TestBed.inject(Router);
    const navigateSpy = vi.spyOn(router, 'navigate');

    await user.click(screen.getByRole('button', { name: /type/i }));
    await user.click(screen.getByText('Beta'));

    expect(navigateSpy).toHaveBeenCalledWith(
      [],
      expect.objectContaining({
        queryParams: { type: 'Beta', offset: null },
        queryParamsHandling: 'merge',
      }),
    );
  });

  it('clears filter when selecting "All"', async () => {
    const user = userEvent.setup();
    await renderTable(filterableColumns);

    const router = TestBed.inject(Router);
    const navigateSpy = vi.spyOn(router, 'navigate');

    await user.click(screen.getByRole('button', { name: /type/i }));
    await user.click(screen.getByText('All'));

    expect(navigateSpy).toHaveBeenCalledWith(
      [],
      expect.objectContaining({
        queryParams: { type: null, offset: null },
      }),
    );
  });

  it('opens search input when clicking a searchable header', async () => {
    const user = userEvent.setup();
    await renderTable(filterableColumns);

    await user.click(screen.getByRole('button', { name: /name/i }));

    expect(screen.getByPlaceholderText('Search name...')).toBeTruthy();
  });

  it('closes dropdown on escape', async () => {
    const user = userEvent.setup();
    await renderTable(filterableColumns);

    await user.click(screen.getByRole('button', { name: /type/i }));
    expect(screen.getByText('All')).toBeTruthy();

    await user.keyboard('{Escape}');
    expect(screen.queryByText('All')).toBeNull();
  });

  it('applies stale opacity when stale is true', async () => {
    const { fixture } = await renderTable(staticColumns, { stale: true });
    const tbody = fixture.nativeElement.querySelector('tbody');
    expect(tbody.classList.contains('opacity-50')).toBe(true);
  });

  it('does not apply stale opacity when stale is false', async () => {
    const { fixture } = await renderTable(staticColumns);
    const tbody = fixture.nativeElement.querySelector('tbody');
    expect(tbody.classList.contains('opacity-50')).toBe(false);
  });

  describe('column visibility', () => {
    it('renders settings gear button when columns are hideable', async () => {
      await renderTable(hideableColumns);
      expect(screen.getByRole('button', { name: /column settings/i })).toBeTruthy();
    });

    it('does not render gear button when all columns have hideable: false', async () => {
      await renderTable(staticColumns);
      expect(screen.queryByRole('button', { name: /column settings/i })).toBeNull();
    });

    it('opens modal on gear click', async () => {
      const user = userEvent.setup();
      await renderTable(hideableColumns);

      await user.click(screen.getByRole('button', { name: /column settings/i }));

      expect(screen.getByRole('dialog')).toBeTruthy();
      expect(screen.getByText('Columns')).toBeTruthy();
    });

    it('shows checkboxes only for hideable columns', async () => {
      const user = userEvent.setup();
      await renderTable(hideableColumns);

      await user.click(screen.getByRole('button', { name: /column settings/i }));

      const checkboxes = screen.getAllByRole('checkbox');
      expect(checkboxes).toHaveLength(2);
      expect(screen.getByLabelText('Name')).toBeTruthy();
      expect(screen.getByLabelText('Value')).toBeTruthy();
    });

    it('generates CSS hide rule when unchecking a column', async () => {
      const user = userEvent.setup();
      const { fixture } = await renderTable(hideableColumns);

      await user.click(screen.getByRole('button', { name: /column settings/i }));
      await user.click(screen.getByLabelText('Name'));

      fixture.detectChanges();

      const tableEl = fixture.nativeElement.querySelector('table')?.parentElement;
      const tableId = tableEl?.id;
      const styleEl = document.head.querySelector(`style`);
      expect(styleEl?.textContent).toContain(`#${tableId} th:nth-child(1)`);
      expect(styleEl?.textContent).toContain(`#${tableId} td:nth-child(1)`);
      expect(styleEl?.textContent).toContain('display: none');
    });

    it('removes CSS hide rule when re-checking a column', async () => {
      const user = userEvent.setup();
      const { fixture } = await renderTable(hideableColumns);

      await user.click(screen.getByRole('button', { name: /column settings/i }));
      await user.click(screen.getByLabelText('Name'));
      fixture.detectChanges();

      await user.click(screen.getByLabelText('Name'));
      fixture.detectChanges();

      const styleEl = document.head.querySelector('style');
      expect(styleEl?.textContent).toBe('');
    });

    it('closes modal on backdrop click', async () => {
      const user = userEvent.setup();
      await renderTable(hideableColumns);

      await user.click(screen.getByRole('button', { name: /column settings/i }));
      expect(screen.getByRole('dialog')).toBeTruthy();

      // Click the backdrop (the fixed overlay parent of the dialog)
      const backdrop = screen.getByRole('dialog').parentElement!;
      await user.click(backdrop);

      expect(screen.queryByRole('dialog')).toBeNull();
    });

    it('closes modal on close button click', async () => {
      const user = userEvent.setup();
      await renderTable(hideableColumns);

      await user.click(screen.getByRole('button', { name: /column settings/i }));
      expect(screen.getByRole('dialog')).toBeTruthy();

      await user.click(screen.getByRole('button', { name: /close/i }));

      expect(screen.queryByRole('dialog')).toBeNull();
    });

    it('closes modal on Escape key', async () => {
      const user = userEvent.setup();
      await renderTable(hideableColumns);

      await user.click(screen.getByRole('button', { name: /column settings/i }));
      expect(screen.getByRole('dialog')).toBeTruthy();

      await user.keyboard('{Escape}');

      expect(screen.queryByRole('dialog')).toBeNull();
    });

    it('hides defaultHidden columns initially', async () => {
      const columnsWithDefaults: DataTableColumn[] = [
        { label: 'Name' },
        { label: 'Extra', defaultHidden: true },
        { label: 'Fixed', hideable: false },
      ];
      const user = userEvent.setup();
      const { fixture } = await renderTable(columnsWithDefaults);

      await fixture.whenStable();
      fixture.detectChanges();

      await user.click(screen.getByRole('button', { name: /column settings/i }));

      const nameCheckbox = screen.getByLabelText('Name') as HTMLInputElement;
      expect(nameCheckbox.checked).toBe(true);

      const extraCheckbox = screen.getByLabelText('Extra') as HTMLInputElement;
      expect(extraCheckbox.checked).toBe(false);
    });

    it('emits hiddenColumnsChange on toggle', async () => {
      const user = userEvent.setup();
      const changeSpy = vi.fn();
      await render(
        `<app-data-table [columns]="columns" (hiddenColumnsChange)="onChange($event)">
          <tr><td>row</td></tr>
        </app-data-table>`,
        {
          imports: [DataTable],
          providers: [provideRouter([]), provideLocationMocks()],
          componentProperties: { columns: hideableColumns, onChange: changeSpy },
        },
      );

      await user.click(screen.getByRole('button', { name: /column settings/i }));
      await user.click(screen.getByLabelText('Value'));

      expect(changeSpy).toHaveBeenCalledWith(['Value']);
    });

    it('uses initialHiddenColumns when provided', async () => {
      const user = userEvent.setup();
      const { fixture } = await render(
        `<app-data-table [columns]="columns" [initialHiddenColumns]="initialHidden">
          <tr><td>row</td></tr>
        </app-data-table>`,
        {
          imports: [DataTable],
          providers: [provideRouter([]), provideLocationMocks()],
          componentProperties: { columns: hideableColumns, initialHidden: ['Name'] },
        },
      );

      await fixture.whenStable();
      fixture.detectChanges();

      await user.click(screen.getByRole('button', { name: /column settings/i }));

      const nameCheckbox = screen.getByLabelText('Name') as HTMLInputElement;
      expect(nameCheckbox.checked).toBe(false);

      const valueCheckbox = screen.getByLabelText('Value') as HTMLInputElement;
      expect(valueCheckbox.checked).toBe(true);
    });
  });
});
