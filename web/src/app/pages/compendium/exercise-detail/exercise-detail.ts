import { Component, inject, computed, signal } from '@angular/core';
import { ActivatedRoute, Router, RouterLink } from '@angular/router';
import { toSignal } from '@angular/core/rxjs-interop';
import { injectQuery, injectMutation, QueryClient } from '@tanstack/angular-query-experimental';
import { CompendiumApiClient } from '$core/api-clients/compendium-api-client';
import { UserApiClient } from '$core/api-clients/user-api-client';
import {
  exerciseKeys,
  exerciseRelationshipKeys,
  equipmentKeys,
  userExerciseKeys,
  masteryKeys,
} from '$core/query-keys';
import { TranslocoDirective } from '@jsverse/transloco';
import { PageLayout } from '../../../layout/page-layout';
import { ConfirmDialog } from '$ui/confirm-dialog/confirm-dialog';
import { DecimalPipe } from '@angular/common';

@Component({
  selector: 'app-exercise-detail',
  imports: [PageLayout, RouterLink, ConfirmDialog, TranslocoDirective, DecimalPipe],
  template: `
    <ng-container *transloco="let t">
      <app-page-layout
        [header]="exerciseQuery.data()?.name ?? 'Exercise'"
        [isPending]="exerciseQuery.isPending()"
        [errorMessage]="exerciseQuery.isError() ? exerciseQuery.error().message : undefined"
      >
        <div actions class="flex gap-2">
          @if (alreadyAdded(); as existing) {
            <a
              [routerLink]="['/user/exercises', existing.id]"
              class="rounded-md bg-gray-500 px-4 py-2 text-sm font-medium text-white hover:bg-gray-600"
            >
              {{ t('compendium.exercises.alreadyAdded') }}
            </a>
          } @else {
            <button
              type="button"
              (click)="addMutation.mutate()"
              [disabled]="addMutation.isPending()"
              class="rounded-md bg-green-600 px-4 py-2 text-sm font-medium text-white hover:bg-green-700 disabled:opacity-50"
            >
              {{
                addMutation.isPending() ? t('common.adding') : t('compendium.exercises.addToMine')
              }}
            </button>
          }
          @if (hasHistory()) {
            <a
              routerLink="./history"
              class="rounded-md bg-gray-600 px-4 py-2 text-sm font-medium text-white hover:bg-gray-700"
              >{{ t('common.history') }}</a
            >
          }
          @if (canModify()) {
            <a
              routerLink="./edit"
              class="rounded-md bg-blue-600 px-4 py-2 text-sm font-medium text-white hover:bg-blue-700"
              >{{ t('common.edit') }}</a
            >
          }
          @if (canDelete()) {
            <button
              type="button"
              (click)="showDeleteDialog.set(true)"
              class="rounded-md bg-red-600 px-4 py-2 text-sm font-medium text-white hover:bg-red-700"
            >
              {{ t('common.delete') }}
            </button>
          }
        </div>
        <app-confirm-dialog
          [open]="showDeleteDialog()"
          [title]="t('compendium.exercises.deleteTitle')"
          [message]="
            t('compendium.exercises.deleteMessage', { name: exerciseQuery.data()?.name ?? '' })
          "
          [isPending]="deleteMutation.isPending()"
          (confirmed)="deleteMutation.mutate()"
          (cancelled)="showDeleteDialog.set(false)"
        />
        @if (exerciseQuery.data(); as exercise) {
          @if (masteryQuery.data(); as mastery) {
            @if (mastery.totalXp > 0) {
              <div
                class="mb-6 rounded-lg border border-amber-200 bg-amber-50 p-4 dark:border-amber-800 dark:bg-amber-900/20"
              >
                <h3 class="text-sm font-medium text-amber-900 dark:text-amber-200">Your Mastery</h3>
                <div class="mt-2 grid grid-cols-3 gap-4 text-sm">
                  <div>
                    <span class="text-amber-600 dark:text-amber-400">Level</span>
                    <span class="block font-semibold text-amber-900 dark:text-amber-100">{{
                      mastery.level
                    }}</span>
                  </div>
                  <div>
                    <span class="text-amber-600 dark:text-amber-400">Tier</span>
                    <span
                      class="block font-semibold capitalize text-amber-900 dark:text-amber-100"
                      >{{ mastery.tier }}</span
                    >
                  </div>
                  <div>
                    <span class="text-amber-600 dark:text-amber-400">XP</span>
                    <span class="block font-semibold text-amber-900 dark:text-amber-100">{{
                      mastery.effectiveXp | number: '1.0-0'
                    }}</span>
                  </div>
                </div>
                <div class="mt-3 h-2 w-full rounded-full bg-amber-200 dark:bg-amber-800">
                  <div
                    class="h-2 rounded-full bg-amber-500"
                    [style.width.%]="mastery.progress * 100"
                  ></div>
                </div>
              </div>
            }
          }
          <dl class="grid grid-cols-1 gap-4 sm:grid-cols-2">
            <div>
              <dt class="text-sm font-medium text-gray-500 dark:text-gray-400">
                {{ t('fields.type') }}
              </dt>
              <dd class="text-sm text-gray-900 dark:text-gray-100">
                {{ t('enums.exerciseType.' + exercise.type) }}
              </dd>
            </div>
            <div>
              <dt class="text-sm font-medium text-gray-500 dark:text-gray-400">
                {{ t('fields.difficulty') }}
              </dt>
              <dd class="text-sm text-gray-900 dark:text-gray-100">
                {{ t('enums.difficulty.' + exercise.technicalDifficulty) }}
              </dd>
            </div>
            @if (exercise.force?.length) {
              <div>
                <dt class="text-sm font-medium text-gray-500 dark:text-gray-400">
                  {{ t('fields.force') }}
                </dt>
                <dd class="text-sm text-gray-900 dark:text-gray-100">
                  @for (f of exercise.force; track f; let last = $last) {
                    {{ t('enums.force.' + f) }}{{ last ? '' : ', ' }}
                  }
                </dd>
              </div>
            }
            <div>
              <dt class="text-sm font-medium text-gray-500 dark:text-gray-400">
                {{ t('fields.primaryMuscles') }}
              </dt>
              <dd class="text-sm text-gray-900 dark:text-gray-100">
                @for (m of exercise.primaryMuscles; track m; let last = $last) {
                  {{ t('enums.muscle.' + m) }}{{ last ? '' : ', ' }}
                }
              </dd>
            </div>
            <div>
              <dt class="text-sm font-medium text-gray-500 dark:text-gray-400">
                {{ t('fields.secondaryMuscles') }}
              </dt>
              <dd class="text-sm text-gray-900 dark:text-gray-100">
                @for (m of exercise.secondaryMuscles; track m; let last = $last) {
                  {{ t('enums.muscle.' + m) }}{{ last ? '' : ', ' }}
                }
              </dd>
            </div>
            @if (exercise.suggestedMeasurementParadigms?.length) {
              <div>
                <dt class="text-sm font-medium text-gray-500 dark:text-gray-400">
                  {{ t('fields.suggestedMeasurementParadigms') }}
                </dt>
                <dd class="text-sm text-gray-900 dark:text-gray-100">
                  @for (p of exercise.suggestedMeasurementParadigms; track p; let last = $last) {
                    {{ t('enums.measurementType.' + p) }}{{ last ? '' : ', ' }}
                  }
                </dd>
              </div>
            }
            @if (equipmentNames().length) {
              <div>
                <dt class="text-sm font-medium text-gray-500 dark:text-gray-400">
                  {{ t('fields.equipmentIds') }}
                </dt>
                <dd class="text-sm text-gray-900 dark:text-gray-100">
                  {{ equipmentNames().join(', ') }}
                </dd>
              </div>
            }
            @if (exercise.alternativeNames?.length) {
              <div>
                <dt class="text-sm font-medium text-gray-500 dark:text-gray-400">
                  {{ t('fields.alternativeNames') }}
                </dt>
                <dd class="text-sm text-gray-900 dark:text-gray-100">
                  {{ exercise.alternativeNames.join(', ') }}
                </dd>
              </div>
            }
            @if (exercise.authorName) {
              <div>
                <dt class="text-sm font-medium text-gray-500 dark:text-gray-400">
                  {{ t('fields.authorName') }}
                </dt>
                <dd class="text-sm text-gray-900 dark:text-gray-100">
                  @if (exercise.authorUrl) {
                    <a
                      [href]="exercise.authorUrl"
                      target="_blank"
                      rel="noopener"
                      class="text-blue-600 hover:underline dark:text-blue-400"
                      >{{ exercise.authorName }}</a
                    >
                  } @else {
                    {{ exercise.authorName }}
                  }
                </dd>
              </div>
            }
            @if (exercise.description) {
              <div class="sm:col-span-2">
                <dt class="text-sm font-medium text-gray-500 dark:text-gray-400">
                  {{ t('fields.description') }}
                </dt>
                <dd class="whitespace-pre-line text-sm text-gray-900 dark:text-gray-100">
                  {{ exercise.description }}
                </dd>
              </div>
            }
          </dl>

          @if (exercise.instructions?.length) {
            <div class="mt-6">
              <h3 class="text-sm font-medium text-gray-500 dark:text-gray-400">
                {{ t('fields.instructions') }}
              </h3>
              <ol
                class="mt-2 list-inside list-decimal space-y-1 text-sm text-gray-900 dark:text-gray-100"
              >
                @for (step of exercise.instructions; track $index) {
                  <li>{{ step }}</li>
                }
              </ol>
            </div>
          }

          @if (allRelationships().length) {
            <div class="mt-6">
              <h3 class="text-sm font-medium text-gray-500 dark:text-gray-400">
                {{ t('compendium.exercises.relationships') }}
              </h3>
              <ul class="mt-2 space-y-1 text-sm">
                @for (rel of allRelationships(); track rel.id) {
                  <li class="flex items-center gap-2">
                    <span
                      class="rounded bg-gray-100 px-1.5 py-0.5 text-xs text-gray-600 dark:bg-gray-800 dark:text-gray-400"
                    >
                      {{ t('enums.exerciseRelationshipType.' + rel.relationshipType) }}
                    </span>
                    <a
                      [routerLink]="['/compendium/exercises', rel.linkedId, rel.linkedSlug]"
                      class="text-blue-600 hover:underline dark:text-blue-400"
                    >
                      {{ rel.linkedName }}
                    </a>
                  </li>
                }
              </ul>
            </div>
          }
        }
      </app-page-layout>
    </ng-container>
  `,
})
export class ExerciseDetail {
  private api = inject(CompendiumApiClient);
  private userApi = inject(UserApiClient);
  private router = inject(Router);
  private queryClient = inject(QueryClient);
  private params = toSignal(inject(ActivatedRoute).paramMap);

  private id = computed(() => Number(this.params()?.get('id')));

  showDeleteDialog = signal(false);

  exerciseQuery = injectQuery(() => ({
    queryKey: exerciseKeys.detail(this.id()),
    queryFn: () => this.api.fetchExercise(this.id()),
    enabled: !!this.id(),
  }));

  versionsQuery = injectQuery(() => ({
    queryKey: exerciseKeys.versions(this.id()),
    queryFn: () => this.api.fetchExerciseVersions(this.id()),
    enabled: !!this.id(),
  }));

  forkedRelationshipsQuery = injectQuery(() => ({
    queryKey: exerciseRelationshipKeys.list({
      toExerciseId: this.id(),
      relationshipType: 'forked',
      owner: 'me',
    }),
    queryFn: () =>
      this.api.fetchExerciseRelationships({
        toExerciseId: this.id(),
        relationshipType: 'forked',
        owner: 'me',
      }),
    enabled: !!this.id(),
  }));

  permissionsQuery = injectQuery(() => ({
    queryKey: exerciseKeys.permissions(this.id()),
    queryFn: () => this.api.fetchExercisePermissions(this.id()),
    enabled: !!this.id(),
  }));

  masteryQuery = injectQuery(() => ({
    queryKey: masteryKeys.detail(this.id()),
    queryFn: () => this.userApi.fetchMastery(this.id()),
    enabled: !!this.id(),
  }));

  canModify = computed(
    () =>
      this.permissionsQuery.isSuccess() &&
      (this.permissionsQuery.data()?.permissions?.includes('MODIFY') ?? false),
  );
  canDelete = computed(
    () =>
      this.permissionsQuery.isSuccess() &&
      (this.permissionsQuery.data()?.permissions?.includes('DELETE') ?? false),
  );

  // All relationships (both directions) for this exercise
  private fromRelationshipsQuery = injectQuery(() => ({
    queryKey: exerciseRelationshipKeys.list({ fromExerciseId: this.id() }),
    queryFn: () => this.api.fetchExerciseRelationships({ fromExerciseId: this.id() }),
    enabled: !!this.id(),
  }));

  private toRelationshipsQuery = injectQuery(() => ({
    queryKey: exerciseRelationshipKeys.list({ toExerciseId: this.id() }),
    queryFn: () => this.api.fetchExerciseRelationships({ toExerciseId: this.id() }),
    enabled: !!this.id(),
  }));

  // Equipment list (to resolve names from IDs)
  private equipmentQuery = injectQuery(() => ({
    queryKey: equipmentKeys.list({ limit: 200 }),
    queryFn: () => this.api.fetchEquipment({ limit: 200 }),
    enabled: !!this.exerciseQuery.data()?.equipmentIds?.length,
  }));

  equipmentNames = computed(() => {
    const ids = new Set(this.exerciseQuery.data()?.equipmentIds ?? []);
    if (!ids.size) return [];
    const items = this.equipmentQuery.data()?.items ?? [];
    return items.filter((e) => ids.has(e.id)).map((e) => e.displayName ?? e.name);
  });

  allRelationships = computed(() => {
    const id = this.id();
    const from = (this.fromRelationshipsQuery.data() ?? []).map((r) => ({
      ...r,
      linkedId: r.toExerciseId,
      linkedName: this.exerciseNameMap().get(r.toExerciseId) ?? `Exercise #${r.toExerciseId}`,
      linkedSlug: this.exerciseSlugMap().get(r.toExerciseId) ?? '',
    }));
    const to = (this.toRelationshipsQuery.data() ?? [])
      .filter((r) => r.fromExerciseId !== id)
      .map((r) => ({
        ...r,
        linkedId: r.fromExerciseId,
        linkedName: this.exerciseNameMap().get(r.fromExerciseId) ?? `Exercise #${r.fromExerciseId}`,
        linkedSlug: this.exerciseSlugMap().get(r.fromExerciseId) ?? '',
      }));
    return [...from, ...to];
  });

  // Exercise names for relationship links — fetch all exercises involved
  private relatedExerciseIds = computed(() => {
    const ids = new Set<number>();
    for (const r of this.fromRelationshipsQuery.data() ?? []) ids.add(r.toExerciseId);
    for (const r of this.toRelationshipsQuery.data() ?? []) ids.add(r.fromExerciseId);
    ids.delete(this.id());
    return ids;
  });

  private relatedExercisesQuery = injectQuery(() => ({
    queryKey: exerciseKeys.list({ limit: 200, relatedTo: this.id() }),
    queryFn: () => this.api.fetchExercises({ limit: 200 }),
    enabled: this.relatedExerciseIds().size > 0,
  }));

  private exerciseNameMap = computed(() => {
    const map = new Map<number, string>();
    for (const e of this.relatedExercisesQuery.data()?.items ?? []) {
      map.set(e.id, e.name);
    }
    return map;
  });

  private exerciseSlugMap = computed(() => {
    const map = new Map<number, string>();
    for (const e of this.relatedExercisesQuery.data()?.items ?? []) {
      map.set(
        e.id,
        e.name
          .toLowerCase()
          .replace(/[^a-z0-9]+/g, '-')
          .replace(/^-|-$/g, ''),
      );
    }
    return map;
  });

  hasHistory = computed(() => (this.versionsQuery.data()?.length ?? 0) > 1);

  alreadyAdded = computed(() => {
    const rels = this.forkedRelationshipsQuery.data();
    if (!rels || rels.length === 0) return undefined;
    return { id: rels[0].fromExerciseId };
  });

  deleteMutation = injectMutation(() => ({
    mutationFn: () => this.api.deleteExercise(this.id()),
    onSuccess: () => {
      this.queryClient.invalidateQueries({ queryKey: exerciseKeys.all() });
      this.router.navigate(['/compendium/exercises']);
    },
  }));

  addMutation = injectMutation(() => ({
    mutationFn: () => {
      const exercise = this.exerciseQuery.data()!;
      return this.userApi.createUserExercise({
        name: exercise.name,
        type: exercise.type,
        force: exercise.force,
        primaryMuscles: exercise.primaryMuscles,
        secondaryMuscles: exercise.secondaryMuscles,
        technicalDifficulty: exercise.technicalDifficulty,
        bodyWeightScaling: exercise.bodyWeightScaling,
        suggestedMeasurementParadigms: exercise.suggestedMeasurementParadigms,
        description: exercise.description,
        instructions: exercise.instructions,
        images: exercise.images,
        alternativeNames: exercise.alternativeNames,
        authorName: exercise.authorName,
        authorUrl: exercise.authorUrl,
        parentExerciseId: exercise.parentExerciseId,
        equipmentIds: exercise.equipmentIds,
        public: false,
        sourceExerciseId: exercise.id,
      } as Record<string, unknown>);
    },
    onSuccess: (created) => {
      this.queryClient.invalidateQueries({ queryKey: userExerciseKeys.all() });
      this.queryClient.invalidateQueries({
        queryKey: exerciseRelationshipKeys.all(),
      });
      this.router.navigate(['/user/exercises', created.id]);
    },
  }));
}
