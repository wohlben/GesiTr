import { Component, signal } from '@angular/core';
import { TestBed } from '@angular/core/testing';
import { form } from '@angular/forms/signals';
import { provideTranslocoForTest } from '$core/testing/transloco-testing';
import { WorkoutStartExerciseItem } from './workout-start-exercise-item';
import { ExerciseDisplayInfo } from './workout-start.store';
import { StartModel } from './workout-start.models';

const REP_BASED_DISPLAY: ExerciseDisplayInfo = {
  name: 'Bench Press',
  summary: '3x10 @ 60kg',
  measurementType: 'REP_BASED',
  sets: [],
};

function makeModel(): StartModel {
  return {
    name: 'Test',
    notes: '',
    sections: [
      {
        id: 1,
        type: 'main',
        label: '',
        exercises: [
          {
            id: 1,
            sourceExerciseSchemeId: 10,
            breakAfterSeconds: 60,
            sets: [
              {
                id: 1,
                targetReps: 10,
                targetWeight: 60,
                targetDuration: null,
                targetDistance: null,
                targetTime: null,
                restAfterSeconds: null,
              },
            ],
          },
        ],
        pendingGroups: [],
      },
    ],
  };
}

@Component({
  selector: 'app-test-host',
  template: `
    <app-workout-start-exercise-item
      [exercise]="startForm.sections[0].exercises[0]"
      [displayInfo]="displayInfo"
      [isLast]="isLast"
      [readonly]="isReadonly"
      (removed)="removedCalled = true"
      (setChanged)="lastSetChanged = $event"
    />
  `,
  imports: [WorkoutStartExerciseItem],
})
class TestHost {
  model = signal(makeModel());
  startForm = form(this.model);
  displayInfo: ExerciseDisplayInfo | undefined = REP_BASED_DISPLAY;
  isLast = false;
  isReadonly = false;
  removedCalled = false;
  lastSetChanged: { setIndex: number } | null = null;
}

describe('WorkoutStartExerciseItem', () => {
  beforeEach(() => {
    TestBed.configureTestingModule({
      imports: [TestHost],
      providers: [provideTranslocoForTest()],
    });
  });

  function setup(
    overrides: { readonly?: boolean; isLast?: boolean; displayInfo?: ExerciseDisplayInfo } = {},
  ) {
    const fixture = TestBed.createComponent(TestHost);
    if (overrides.readonly != null) fixture.componentInstance.isReadonly = overrides.readonly;
    if (overrides.isLast != null) fixture.componentInstance.isLast = overrides.isLast;
    if (overrides.displayInfo !== undefined)
      fixture.componentInstance.displayInfo = overrides.displayInfo;
    fixture.detectChanges();
    return fixture;
  }

  it('shows drag handle and remove button in editable mode', () => {
    const fixture = setup({ readonly: false });
    const el = fixture.nativeElement as HTMLElement;

    expect(el.querySelector('[cdkDragHandle]')).toBeTruthy();
    expect(el.textContent).toContain('common.remove');
  });

  it('hides drag handle and remove button in readonly mode', () => {
    const fixture = setup({ readonly: true });
    const el = fixture.nativeElement as HTMLElement;

    expect(el.querySelector('[cdkDragHandle]')).toBeNull();
    expect(el.textContent).not.toContain('common.remove');
  });

  it('shows exercise name from display info', () => {
    const fixture = setup();
    const el = fixture.nativeElement as HTMLElement;

    expect(el.textContent).toContain('Bench Press');
    expect(el.textContent).toContain('3x10 @ 60kg');
  });

  it('shows set inputs in editable mode', () => {
    const fixture = setup({ readonly: false });
    const el = fixture.nativeElement as HTMLElement;

    const inputs = el.querySelectorAll(
      'input[data-field="targetReps"], input[data-field="targetWeight"]',
    );
    expect(inputs.length).toBe(2);
  });

  it('keeps set inputs in readonly mode', () => {
    const fixture = setup({ readonly: true });
    const el = fixture.nativeElement as HTMLElement;

    const inputs = el.querySelectorAll(
      'input[data-field="targetReps"], input[data-field="targetWeight"]',
    );
    expect(inputs.length).toBe(2);
  });

  it('shows break-after-exercise when not last', () => {
    const fixture = setup({ isLast: false });
    const el = fixture.nativeElement as HTMLElement;

    expect(el.querySelector('input[data-field="breakAfterSeconds"]')).toBeTruthy();
  });

  it('hides break-after-exercise when last', () => {
    const fixture = setup({ isLast: true });
    const el = fixture.nativeElement as HTMLElement;

    expect(el.querySelector('input[data-field="breakAfterSeconds"]')).toBeNull();
  });

  it('emits removed event when remove button clicked', () => {
    const fixture = setup({ readonly: false });
    const el = fixture.nativeElement as HTMLElement;

    const removeBtn = Array.from(el.querySelectorAll('button')).find(
      (b) => b.textContent?.trim() === 'common.remove',
    );
    removeBtn?.click();
    fixture.detectChanges();

    expect(fixture.componentInstance.removedCalled).toBe(true);
  });

  it('emits setChanged when a set input changes', () => {
    const fixture = setup({ readonly: false });
    const el = fixture.nativeElement as HTMLElement;

    const repsInput = el.querySelector('input[data-field="targetReps"]') as HTMLInputElement;
    repsInput.value = '12';
    repsInput.dispatchEvent(new Event('change'));
    fixture.detectChanges();

    expect(fixture.componentInstance.lastSetChanged).toEqual({ setIndex: 0 });
  });
});
