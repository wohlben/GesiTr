import { TestBed } from '@angular/core/testing';
import { ExerciseGroupConfig, GroupConfigValue } from './exercise-group-config';
import { provideTranslocoForTest } from '$core/testing/transloco-testing';
import { ExerciseGroup } from '$generated/models';

// brn-select uses ResizeObserver which isn't available in the unit test environment
beforeAll(() => {
  globalThis.ResizeObserver ??= class {
    // eslint-disable-next-line @typescript-eslint/no-empty-function
    observe() {}
    // eslint-disable-next-line @typescript-eslint/no-empty-function
    unobserve() {}
    // eslint-disable-next-line @typescript-eslint/no-empty-function
    disconnect() {}
  } as unknown as typeof ResizeObserver;
});

const EXERCISES = [
  { id: 1, name: 'Bench Press' },
  { id: 2, name: 'Dumbbell Press' },
];

const GROUP: ExerciseGroup = {
  id: 5,
  createdAt: '',
  updatedAt: '',
  name: 'Push Group',
  owner: 'anon',
};

describe('ExerciseGroupConfig', () => {
  beforeEach(() => {
    TestBed.configureTestingModule({
      imports: [ExerciseGroupConfig],
      providers: [provideTranslocoForTest()],
    });
  });

  function setup(initial: GroupConfigValue, existingGroups: ExerciseGroup[] = []) {
    const fixture = TestBed.createComponent(ExerciseGroupConfig);
    fixture.componentRef.setInput('existingGroups', existingGroups);
    fixture.componentRef.setInput('exercises', EXERCISES);
    fixture.componentInstance.value.set(initial);
    fixture.detectChanges();
    return { fixture };
  }

  it('preserves exerciseGroupId when groups load later', async () => {
    const { fixture } = setup({ exerciseGroupId: 5, name: 'Push Group', members: [1, 2] }, []);

    expect(fixture.componentInstance.value().exerciseGroupId).toBe(5);

    fixture.componentRef.setInput('existingGroups', [GROUP]);
    fixture.detectChanges();
    await fixture.whenStable();

    expect(fixture.componentInstance.value().exerciseGroupId).toBe(5);
  });

  it('writes name changes to the model signal', () => {
    const { fixture } = setup({ exerciseGroupId: null, name: 'Initial', members: [] });

    const nameInput = fixture.nativeElement.querySelector('input') as HTMLInputElement;
    expect(nameInput.value).toBe('Initial');

    nameInput.value = 'Updated';
    nameInput.dispatchEvent(new Event('input'));
    fixture.detectChanges();

    expect(fixture.componentInstance.value().name).toBe('Updated');
  });

  it('populates group name from existingGroups on selection change', () => {
    const { fixture } = setup({ exerciseGroupId: null, name: 'Old Name', members: [1] }, [GROUP]);

    fixture.componentInstance.onGroupSelect(5);
    fixture.detectChanges();

    const v = fixture.componentInstance.value();
    expect(v.exerciseGroupId).toBe(5);
    expect(v.name).toBe('Push Group');
    expect(v.members).toEqual([]);
  });

  it('clears name when switching to New Group', () => {
    const { fixture } = setup({ exerciseGroupId: 5, name: 'Push Group', members: [1] }, [GROUP]);

    fixture.componentInstance.onGroupSelect(null);
    fixture.detectChanges();

    const v = fixture.componentInstance.value();
    expect(v.exerciseGroupId).toBeNull();
    expect(v.name).toBe('');
    expect(v.members).toEqual([]);
  });

  it('adds members to the value', () => {
    const { fixture } = setup({ exerciseGroupId: null, name: '', members: [1] });

    const memberSelect = fixture.nativeElement.querySelector('select') as HTMLSelectElement;
    memberSelect.value = '2';
    memberSelect.dispatchEvent(new Event('change'));

    expect(fixture.componentInstance.value().members).toEqual([1, 2]);
  });

  it('removes members from the value', () => {
    const { fixture } = setup({ exerciseGroupId: null, name: '', members: [1, 2] });

    fixture.componentInstance.onRemoveMember(1);

    expect(fixture.componentInstance.value().members).toEqual([2]);
  });
});
