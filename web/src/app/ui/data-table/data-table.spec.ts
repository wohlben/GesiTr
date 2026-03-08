import { render, screen } from '@testing-library/angular';
import userEvent from '@testing-library/user-event';
import { TestBed } from '@angular/core/testing';
import { provideRouter, Router } from '@angular/router';
import { provideLocationMocks } from '@angular/common/testing';
import { DataTable, DataTableColumn } from './data-table';

describe('DataTable', () => {
  const staticColumns: DataTableColumn[] = [
    { label: 'Name' },
    { label: 'Value' },
  ];

  const filterableColumns: DataTableColumn[] = [
    { label: 'Name', searchParam: 'q' },
    { label: 'Type', filterParam: 'type', options: ['Alpha', 'Beta', 'Gamma'] },
    { label: 'Status' },
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

  it('renders static columns without buttons', async () => {
    await renderTable(staticColumns);
    expect(screen.queryByRole('button')).toBeNull();
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
});
