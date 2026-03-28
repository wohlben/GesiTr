import { Component, computed, effect, inject, signal } from '@angular/core';
import { ActivatedRoute, Router } from '@angular/router';
import { toSignal } from '@angular/core/rxjs-interop';
import { form, FormField, min, max } from '@angular/forms/signals';
import { injectQuery, injectMutation, QueryClient } from '@tanstack/angular-query-experimental';
import { TranslocoDirective } from '@jsverse/transloco';
import { HlmToggleGroupImports } from '@spartan-ng/helm/toggle-group';
import { HlmTooltip } from '@spartan-ng/helm/tooltip';
import { HlmDatePicker } from '@spartan-ng/helm/date-picker';
import { provideNativeDateAdapter } from '@spartan-ng/brain/date-time';
import { PeriodDayPicker } from '$ui/period-day-picker/period-day-picker';
import { HlmInput } from '@spartan-ng/helm/input';
import { HlmAlertImports } from '@spartan-ng/helm/alert';
import { NgIcon, provideIcons } from '@ng-icons/core';
import { lucideInfo, lucideTriangleAlert } from '@ng-icons/lucide';
import { UserApiClient } from '$core/api-clients/user-api-client';
import { workoutScheduleKeys, schedulePeriodKeys } from '$core/query-keys';
import { PageLayout } from '../../../layout/page-layout';

@Component({
  selector: 'app-workout-schedule-edit',
  imports: [
    PageLayout,
    FormField,
    TranslocoDirective,
    HlmToggleGroupImports,
    PeriodDayPicker,
    HlmDatePicker,
    HlmInput,
    HlmAlertImports,
    HlmTooltip,
    NgIcon,
  ],
  providers: [provideNativeDateAdapter(), provideIcons({ lucideInfo, lucideTriangleAlert })],
  template: `
    <ng-container *transloco="let t">
      <app-page-layout
        [header]="
          isCreateMode() ? t('user.schedules.newSchedule') : t('user.schedules.editSchedule')
        "
        [isPending]="!isCreateMode() && scheduleQuery.isPending()"
        [errorMessage]="scheduleQuery.isError() ? scheduleQuery.error().message : undefined"
      >
        <form (submit)="onSubmit(); $event.preventDefault()" class="mx-auto max-w-2xl space-y-8">
          <!-- Initial Status -->
          <section>
            <h3
              class="mb-3 text-sm font-semibold tracking-wider text-gray-500 uppercase dark:text-gray-400"
            >
              {{ t('user.schedules.initialStatus') }}
            </h3>
            <div class="flex items-start gap-3">
              <div
                hlmToggleGroup
                type="single"
                [formField]="scheduleForm.initialStatus"
                variant="outline"
              >
                <button hlmToggleGroupItem value="committed">
                  {{ t('enums.workoutLogStatus.committed') }}
                </button>
                <button hlmToggleGroupItem value="proposed">
                  {{ t('enums.workoutLogStatus.proposed') }}
                </button>
              </div>
              <ng-icon
                name="lucideInfo"
                class="mt-2 shrink-0 cursor-help text-gray-400 dark:text-gray-500"
                [hlmTooltip]="
                  model().initialStatus === 'committed'
                    ? t('user.schedules.committedHelp')
                    : t('user.schedules.proposedHelp')
                "
              />
            </div>
          </section>

          <!-- Schedule Start Date -->
          <section>
            <h3
              class="mb-3 text-sm font-semibold tracking-wider text-gray-500 uppercase dark:text-gray-400"
            >
              {{ t('user.schedules.activeRange') }}
            </h3>
            <div class="flex items-center gap-3">
              <label
                for="startDate"
                class="shrink-0 text-sm font-medium text-gray-700 dark:text-gray-300"
              >
                {{ t('common.from') }}
              </label>
              <input
                id="startDate"
                type="date"
                [formField]="scheduleForm.startDate"
                hlmInput
                class="w-48"
              />
            </div>
          </section>

          <!-- Next Period -->
          <section>
            <h3
              class="mb-3 text-sm font-semibold tracking-wider text-gray-500 uppercase dark:text-gray-400"
            >
              {{ t('user.schedules.nextPeriod') }}
            </h3>
            <div class="grid grid-cols-[auto_auto_1fr] items-end gap-x-6 gap-y-1">
              <!-- Row 1: labels + button group -->
              <label for="periodStart" class="text-sm font-medium text-gray-700 dark:text-gray-300">
                {{ t('common.from') }}
              </label>
              <label for="periodEnd" class="text-sm font-medium text-gray-700 dark:text-gray-300">
                {{ t('common.to') }}
              </label>
              <div></div>

              <!-- Row 2: pickers + mode toggle -->
              <hlm-date-picker [formField]="periodForm.startDate" [min]="scheduleStartDate()" />
              <hlm-date-picker [formField]="periodForm.endDate" [min]="tomorrow" />
              @if (hasPeriod()) {
                <div class="flex items-center gap-2">
                  <div hlmToggleGroup type="single" [formField]="periodForm.mode" variant="outline">
                    <button hlmToggleGroupItem value="normal">
                      {{ periodDayCount() }} {{ periodDayCount() === 1 ? 'day' : 'days' }}
                    </button>
                    @if (looksMonthly()) {
                      <button hlmToggleGroupItem value="monthly">
                        {{ t('user.schedules.modeMonthly') }}
                      </button>
                    }
                  </div>
                  <ng-icon
                    name="lucideInfo"
                    class="shrink-0 cursor-help text-gray-400 dark:text-gray-500"
                    [hlmTooltip]="
                      periodModel().mode === 'monthly'
                        ? t('user.schedules.modeMonthlyHelp')
                        : t('user.schedules.modeNormalHelp')
                    "
                  />
                </div>
              } @else {
                <div></div>
              }
            </div>
          </section>

          <!-- Period Type -->
          <section>
            <h3
              class="mb-3 text-sm font-semibold tracking-wider text-gray-500 uppercase dark:text-gray-400"
            >
              {{ t('user.schedules.type') }}
            </h3>
            <div class="flex items-start gap-3">
              <div hlmToggleGroup type="single" [formField]="periodForm.type" variant="outline">
                <button hlmToggleGroupItem value="fixed_date">
                  {{ t('enums.scheduleType.fixed_date') }}
                </button>
                <button hlmToggleGroupItem value="frequency">
                  {{ t('enums.scheduleType.frequency') }}
                </button>
              </div>
              <ng-icon
                name="lucideInfo"
                class="mt-2 shrink-0 cursor-help text-gray-400 dark:text-gray-500"
                [hlmTooltip]="
                  periodModel().type === 'fixed_date'
                    ? t('user.schedules.fixedDateHelp')
                    : t('user.schedules.frequencyHelp')
                "
              />
            </div>
          </section>

          @if (hasPeriod()) {
            <!-- Fixed Date: pick specific days -->
            @if (periodModel().type === 'fixed_date') {
              <section>
                <h3
                  class="mb-3 text-sm font-semibold tracking-wider text-gray-500 uppercase dark:text-gray-400"
                >
                  {{ t('user.schedules.selectDays') }}
                </h3>
                <p class="mb-4 text-xs text-gray-500 dark:text-gray-400">
                  {{ t('user.schedules.selectDaysHelp') }}
                </p>
                <app-period-day-picker
                  [periodStart]="periodStartDate()!"
                  [periodEnd]="periodEndDate()!"
                  [formField]="periodForm.selectedDates"
                />
              </section>
            }

            <!-- Frequency: number input -->
            @if (periodModel().type === 'frequency') {
              <section>
                <div class="flex items-center gap-3">
                  <label
                    for="frequencyCount"
                    class="shrink-0 text-sm font-medium text-gray-700 dark:text-gray-300"
                  >
                    {{ t('user.schedules.frequencyCount') }}
                  </label>
                  <input
                    id="frequencyCount"
                    type="number"
                    [formField]="periodForm.frequencyCount"
                    hlmInput
                    class="w-24"
                  />
                </div>
              </section>
            }
          }

          <!-- Warning: period starts today or earlier -->
          @if (periodIsImmediate()) {
            <div hlmAlert variant="destructive">
              <ng-icon hlmAlertIcon name="lucideTriangleAlert" />
              <h4 hlmAlertTitle>{{ t('user.schedules.immediatePeriodTitle') }}</h4>
              <p hlmAlertDescription>{{ t('user.schedules.immediatePeriodWarning') }}</p>
            </div>
          }

          <!-- Submit -->
          <div class="flex gap-3">
            <button
              type="submit"
              [disabled]="
                createMutation.isPending() ||
                updateMutation.isPending() ||
                (isCreateMode() && !canCreate())
              "
              class="rounded-md bg-blue-600 px-4 py-2 text-sm font-medium text-white hover:bg-blue-700 disabled:opacity-50"
            >
              @if (isCreateMode()) {
                {{ t('common.create') }}
              } @else {
                {{ t('common.save') }}
              }
            </button>
          </div>
        </form>
      </app-page-layout>
    </ng-container>
  `,
})
export class WorkoutScheduleEdit {
  private route = inject(ActivatedRoute);
  private router = inject(Router);
  private userApi = inject(UserApiClient);
  private queryClient = inject(QueryClient);
  private params = toSignal(this.route.paramMap);

  private workoutId = computed(() => Number(this.params()?.get('id')));
  private scheduleId = computed(() => Number(this.params()?.get('scheduleId')));
  isCreateMode = computed(() => !this.params()?.get('scheduleId'));

  model = signal({
    startDate: '',
    initialStatus: 'committed',
  });

  scheduleForm = form(this.model);

  // Period form (separate from schedule — period is its own entity)
  periodModel = signal({
    startDate: null as Date | null,
    endDate: null as Date | null,
    type: 'fixed_date',
    mode: 'normal',
    frequencyCount: 3,
    selectedDates: [] as Date[],
  });

  periodForm = form(this.periodModel, (f) => {
    min(f.frequencyCount, 1);
    max(f.frequencyCount, () => this.periodDayCount());
  });

  tomorrow = (() => {
    const d = new Date();
    d.setDate(d.getDate() + 1);
    d.setHours(0, 0, 0, 0);
    return d;
  })();

  scheduleStartDate = computed(() => {
    const s = this.model().startDate;
    return s ? new Date(s) : undefined;
  });

  periodStartDate = computed(() => this.periodModel().startDate ?? undefined);

  periodEndDate = computed(() => this.periodModel().endDate ?? undefined);

  hasPeriod = computed(() => !!this.periodModel().startDate && !!this.periodModel().endDate);

  periodDayCount = computed(() => {
    const { startDate, endDate } = this.periodModel();
    if (!startDate || !endDate) return 1;
    const diff = endDate.getTime() - startDate.getTime();
    return Math.max(1, Math.ceil(diff / (1000 * 60 * 60 * 24)) + 1);
  });

  looksMonthly = computed(() => {
    const { startDate, endDate } = this.periodModel();
    if (!startDate || !endDate) return false;
    const oneMonthLater = new Date(startDate);
    oneMonthLater.setMonth(oneMonthLater.getMonth() + 1);
    oneMonthLater.setDate(oneMonthLater.getDate() - 1);
    return (
      endDate.getFullYear() === oneMonthLater.getFullYear() &&
      endDate.getMonth() === oneMonthLater.getMonth() &&
      endDate.getDate() === oneMonthLater.getDate()
    );
  });

  canCreate = computed(() => {
    const p = this.periodModel();
    if (!p.startDate || !p.endDate) return false;
    if (p.type === 'fixed_date' && p.selectedDates.length === 0) return false;
    if (p.type === 'frequency' && p.frequencyCount < 1) return false;
    return true;
  });

  periodIsImmediate = computed(() => {
    const s = this.periodModel().startDate;
    if (!s) return false;
    const today = new Date();
    today.setHours(0, 0, 0, 0);
    return s <= today;
  });

  scheduleQuery = injectQuery(() => ({
    queryKey: workoutScheduleKeys.detail(this.scheduleId()),
    queryFn: () => this.userApi.fetchWorkoutSchedule(this.scheduleId()),
    enabled: !!this.scheduleId() && !this.isCreateMode(),
  }));

  periodsQuery = injectQuery(() => ({
    queryKey: schedulePeriodKeys.list(this.scheduleId()),
    queryFn: () => this.userApi.fetchSchedulePeriods({ scheduleId: this.scheduleId() }),
    enabled: !!this.scheduleId() && !this.isCreateMode(),
  }));

  lastPeriod = computed(() => {
    const periods = this.periodsQuery.data();
    if (!periods || periods.length === 0) return undefined;
    return periods[periods.length - 1];
  });

  commitmentsQuery = injectQuery(() => {
    const period = this.lastPeriod();
    return {
      queryKey: ['schedule-commitments', 'list', period?.id],
      queryFn: () => this.userApi.fetchScheduleCommitments({ periodId: period!.id }),
      enabled: !!period,
    };
  });

  createMutation = injectMutation(() => ({
    mutationFn: async (data: {
      schedule: Record<string, unknown>;
      period?: { start: Date; end: Date; type: string; mode: string };
      commitmentDates?: Date[];
      frequencyCount?: number;
    }) => {
      const schedule = await this.userApi.createWorkoutSchedule(data.schedule);
      if (data.period) {
        const period = await this.userApi.createSchedulePeriod({
          scheduleId: schedule.id,
          periodStart: data.period.start.toISOString(),
          periodEnd: data.period.end.toISOString(),
          type: data.period.type,
          mode: data.period.mode,
        });
        // fixed_date: create commitments with specific dates
        if (data.commitmentDates) {
          for (const date of data.commitmentDates) {
            await this.userApi.createScheduleCommitment({
              periodId: period.id,
              date: date.toISOString(),
            });
          }
        }
        // frequency: create N commitments without dates
        if (data.frequencyCount) {
          for (let i = 0; i < data.frequencyCount; i++) {
            await this.userApi.createScheduleCommitment({
              periodId: period.id,
            });
          }
        }
      }
      return schedule;
    },
    onSuccess: () => {
      this.queryClient.invalidateQueries({ queryKey: workoutScheduleKeys.all() });
      this.router.navigate(['..'], { relativeTo: this.route });
    },
  }));

  updateMutation = injectMutation(() => ({
    mutationFn: (data: Record<string, unknown>) =>
      this.userApi.updateWorkoutSchedule(this.scheduleId(), data),
    onSuccess: () => {
      this.queryClient.invalidateQueries({ queryKey: workoutScheduleKeys.all() });
      this.router.navigate(['../..'], { relativeTo: this.route });
    },
  }));

  constructor() {
    const tomorrow = new Date();
    tomorrow.setDate(tomorrow.getDate() + 1);
    const tomorrowStr = tomorrow.toISOString().substring(0, 10);
    this.model.set({
      startDate: tomorrowStr,
      initialStatus: 'committed',
    });
    const defaultStart = new Date(tomorrow);
    defaultStart.setHours(0, 0, 0, 0);
    const defaultEnd = new Date(defaultStart);
    defaultEnd.setDate(defaultEnd.getDate() + 6);
    this.periodModel.set({
      startDate: defaultStart,
      endDate: defaultEnd,
      type: 'fixed_date',
      mode: 'normal',
      frequencyCount: 3,
      selectedDates: [],
    });

    effect(() => {
      const data = this.scheduleQuery.data();
      if (!data) return;
      this.model.set({
        startDate: data.startDate.substring(0, 10),
        initialStatus: data.initialStatus,
      });
    });

    // Populate period form from last period + commitments in edit mode
    effect(() => {
      const period = this.lastPeriod();
      const commitments = this.commitmentsQuery.data();
      if (!period) return;
      this.periodModel.set({
        startDate: new Date(period.periodStart),
        endDate: new Date(period.periodEnd),
        type: period.type,
        mode: period.mode ?? 'normal',
        frequencyCount: commitments?.length ?? 3,
        selectedDates: (commitments ?? []).filter((c) => c.date).map((c) => new Date(c.date!)),
      });
    });
  }

  onSubmit() {
    const val = this.model();
    if (this.isCreateMode()) {
      const schedule: Record<string, unknown> = {
        workoutId: this.workoutId(),
        initialStatus: val.initialStatus,
      };
      if (val.startDate) {
        schedule['startDate'] = new Date(val.startDate).toISOString();
      }

      const pVal = this.periodModel();
      const period =
        pVal.startDate && pVal.endDate
          ? {
              start: pVal.startDate,
              end: pVal.endDate,
              type: pVal.type,
              mode: this.looksMonthly() ? pVal.mode : 'normal',
            }
          : undefined;

      const commitmentDates =
        pVal.type === 'fixed_date' && pVal.selectedDates.length > 0
          ? pVal.selectedDates
          : undefined;

      const frequencyCount = pVal.type === 'frequency' ? pVal.frequencyCount : undefined;

      this.createMutation.mutate({ schedule, period, commitmentDates, frequencyCount });
    } else {
      this.updateMutation.mutate({});
    }
  }
}
