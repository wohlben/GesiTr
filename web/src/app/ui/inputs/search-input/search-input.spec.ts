import { render, screen } from '@testing-library/angular';
import userEvent from '@testing-library/user-event';
import { SearchInput } from './search-input';

describe('SearchInput', () => {
  it('renders with default placeholder', async () => {
    await render(SearchInput);
    expect(screen.getByPlaceholderText('Search...')).toBeTruthy();
  });

  it('renders with custom placeholder', async () => {
    await render(SearchInput, {
      inputs: { placeholder: 'Find something...' },
    });
    expect(screen.getByPlaceholderText('Find something...')).toBeTruthy();
  });

  it('displays the bound value', async () => {
    await render(SearchInput, {
      inputs: { value: 'hello' },
    });
    expect(screen.getByDisplayValue('hello')).toBeTruthy();
  });

  it('emits value changes on user input', async () => {
    const user = userEvent.setup();
    const { fixture } = await render(SearchInput);

    const input = screen.getByRole('textbox');
    await user.type(input, 'test');

    expect(fixture.componentInstance.value()).toBe('test');
  });
});
