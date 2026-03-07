import { render, screen } from '@testing-library/angular';
import userEvent from '@testing-library/user-event';
import { FilterSelect } from './filter-select';

describe('FilterSelect', () => {
  it('renders with default "All" label', async () => {
    await render(FilterSelect);
    expect(screen.getByRole('option', { name: 'All' })).toBeTruthy();
  });

  it('renders with custom allLabel', async () => {
    await render(FilterSelect, {
      inputs: { allLabel: 'All types' },
    });
    expect(screen.getByRole('option', { name: 'All types' })).toBeTruthy();
  });

  it('renders provided options', async () => {
    await render(FilterSelect, {
      inputs: { options: ['Alpha', 'Beta', 'Gamma'] },
    });
    expect(screen.getByRole('option', { name: 'Alpha' })).toBeTruthy();
    expect(screen.getByRole('option', { name: 'Beta' })).toBeTruthy();
    expect(screen.getByRole('option', { name: 'Gamma' })).toBeTruthy();
  });

  it('selects the empty value by default', async () => {
    await render(FilterSelect, {
      inputs: { options: ['A', 'B'], allLabel: 'Pick one' },
    });
    const select = screen.getByRole('combobox') as HTMLSelectElement;
    expect(select.value).toBe('');
  });

  it('updates value when user selects an option', async () => {
    const user = userEvent.setup();
    const { fixture } = await render(FilterSelect, {
      inputs: { options: ['X', 'Y'], allLabel: 'All' },
    });

    await user.selectOptions(screen.getByRole('combobox'), 'Y');
    expect(fixture.componentInstance.value()).toBe('Y');
  });
});
