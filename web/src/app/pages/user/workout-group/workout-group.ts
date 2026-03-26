import { Component, inject, computed, signal } from '@angular/core';
import { ActivatedRoute, RouterLink } from '@angular/router';
import { injectQuery, injectMutation, QueryClient } from '@tanstack/angular-query-experimental';
import { TranslocoDirective } from '@jsverse/transloco';
import { UserApiClient } from '$core/api-clients/user-api-client';
import { workoutKeys, workoutGroupKeys } from '$core/query-keys';
import { WorkoutGroupRoleMember } from '$generated/user-workoutgroup';
import { PageLayout } from '../../../layout/page-layout';
import { ConfirmDialog } from '$ui/confirm-dialog/confirm-dialog';

@Component({
  selector: 'app-workout-group',
  imports: [PageLayout, RouterLink, TranslocoDirective, ConfirmDialog],
  template: `
    <ng-container *transloco="let t">
      <app-page-layout
        [header]="headerText()"
        [isPending]="workoutQuery.isPending()"
        [errorMessage]="workoutQuery.isError() ? workoutQuery.error().message : undefined"
      >
        <a
          actions
          routerLink="/user/workouts"
          class="text-sm text-blue-600 hover:underline dark:text-blue-400"
        >
          &larr; {{ t('common.back') }}
        </a>

        @if (workoutQuery.data(); as workout) {
          <!-- No group yet: create form -->
          @if (!group()) {
            <div class="rounded-lg border border-gray-200 p-6 dark:border-gray-700">
              <label
                for="create-group-name"
                class="mb-2 block text-sm font-medium text-gray-700 dark:text-gray-300"
              >
                {{ t('user.workouts.groupName') }}
              </label>
              <input
                id="create-group-name"
                type="text"
                [value]="groupNameInput()"
                (input)="groupNameInput.set($any($event.target).value)"
                class="mb-4 block w-full rounded-md border border-gray-300 px-3 py-2 text-sm focus:border-blue-500 focus:ring-blue-500 dark:border-gray-600 dark:bg-gray-800 dark:text-gray-100"
                [placeholder]="t('user.workouts.groupName')"
              />
              <button
                type="button"
                (click)="createGroup()"
                [disabled]="createGroupMutation.isPending() || !groupNameInput()"
                class="rounded-md bg-blue-600 px-4 py-2 text-sm font-medium text-white hover:bg-blue-700 disabled:opacity-50"
              >
                {{ t('user.workouts.createGroup') }}
              </button>
            </div>
          }

          <!-- Group exists -->
          @if (group(); as g) {
            <!-- Group name -->
            <div class="rounded-lg border border-gray-200 p-6 dark:border-gray-700">
              <label
                for="edit-group-name"
                class="mb-2 block text-sm font-medium text-gray-700 dark:text-gray-300"
              >
                {{ t('user.workouts.groupName') }}
              </label>
              <div class="flex gap-2">
                <input
                  id="edit-group-name"
                  type="text"
                  [value]="groupNameInput()"
                  (input)="groupNameInput.set($any($event.target).value)"
                  class="block w-full rounded-md border border-gray-300 px-3 py-2 text-sm focus:border-blue-500 focus:ring-blue-500 dark:border-gray-600 dark:bg-gray-800 dark:text-gray-100"
                />
                <button
                  type="button"
                  (click)="updateGroupName()"
                  [disabled]="updateGroupMutation.isPending() || groupNameInput() === g.name"
                  class="shrink-0 rounded-md bg-blue-600 px-4 py-2 text-sm font-medium text-white hover:bg-blue-700 disabled:opacity-50"
                >
                  {{ t('common.save') }}
                </button>
              </div>
              <button
                type="button"
                (click)="showDeleteDialog.set(true)"
                class="mt-3 text-sm text-red-600 hover:underline dark:text-red-400"
              >
                {{ t('user.workouts.deleteGroup') }}
              </button>
            </div>

            <!-- Members -->
            <div class="mt-6 rounded-lg border border-gray-200 p-6 dark:border-gray-700">
              <h2 class="mb-4 text-lg font-semibold text-gray-900 dark:text-gray-100">
                {{ t('user.workouts.groupMembers') }}
              </h2>

              @if (membershipsQuery.data(); as memberships) {
                @if (memberships.length === 0) {
                  <p class="mb-4 text-sm text-gray-500 dark:text-gray-400">
                    {{ t('user.workouts.noMembers') }}
                  </p>
                } @else {
                  <table class="mb-4 min-w-full divide-y divide-gray-200 dark:divide-gray-700">
                    <thead>
                      <tr>
                        <th
                          class="px-4 py-2 text-left text-xs font-medium tracking-wider text-gray-500 uppercase dark:text-gray-400"
                        >
                          {{ t('user.workouts.userId') }}
                        </th>
                        <th
                          class="px-4 py-2 text-left text-xs font-medium tracking-wider text-gray-500 uppercase dark:text-gray-400"
                        >
                          {{ t('user.workouts.role') }}
                        </th>
                        <th class="px-4 py-2"></th>
                      </tr>
                    </thead>
                    <tbody class="divide-y divide-gray-200 dark:divide-gray-700">
                      @for (m of memberships; track m.id) {
                        <tr>
                          <td class="px-4 py-2 text-sm text-gray-900 dark:text-gray-100">
                            {{ m.userId }}
                          </td>
                          <td class="px-4 py-2 text-sm">
                            <select
                              [value]="m.role"
                              (change)="updateMemberRole(m.id, $any($event.target).value)"
                              class="rounded-md border border-gray-300 px-2 py-1 text-sm dark:border-gray-600 dark:bg-gray-800 dark:text-gray-100"
                            >
                              <option value="member">{{ t('user.workouts.roleMember') }}</option>
                              <option value="admin">{{ t('user.workouts.roleAdmin') }}</option>
                            </select>
                          </td>
                          <td class="px-4 py-2 text-right">
                            <button
                              type="button"
                              (click)="removeMember(m.id)"
                              class="text-sm text-red-600 hover:underline dark:text-red-400"
                            >
                              {{ t('common.remove') }}
                            </button>
                          </td>
                        </tr>
                      }
                    </tbody>
                  </table>
                }
              }

              <!-- Add member form -->
              <div class="flex gap-2">
                <input
                  type="text"
                  [value]="newMemberUserId()"
                  (input)="newMemberUserId.set($any($event.target).value)"
                  class="block w-full rounded-md border border-gray-300 px-3 py-2 text-sm focus:border-blue-500 focus:ring-blue-500 dark:border-gray-600 dark:bg-gray-800 dark:text-gray-100"
                  [placeholder]="t('user.workouts.userId')"
                />
                <select
                  [value]="newMemberRole()"
                  (change)="newMemberRole.set($any($event.target).value)"
                  class="shrink-0 rounded-md border border-gray-300 px-2 py-2 text-sm dark:border-gray-600 dark:bg-gray-800 dark:text-gray-100"
                >
                  <option value="member">{{ t('user.workouts.roleMember') }}</option>
                  <option value="admin">{{ t('user.workouts.roleAdmin') }}</option>
                </select>
                <button
                  type="button"
                  (click)="addMember()"
                  [disabled]="addMemberMutation.isPending() || !newMemberUserId()"
                  class="shrink-0 rounded-md bg-blue-600 px-4 py-2 text-sm font-medium text-white hover:bg-blue-700 disabled:opacity-50"
                >
                  {{ t('user.workouts.addMember') }}
                </button>
              </div>
            </div>
          }
        }
      </app-page-layout>

      <app-confirm-dialog
        [open]="showDeleteDialog()"
        [title]="t('user.workouts.deleteGroup')"
        [message]="t('user.workouts.deleteGroupMessage')"
        [isPending]="deleteGroupMutation.isPending()"
        (confirmed)="deleteGroup()"
        (cancelled)="showDeleteDialog.set(false)"
      />
    </ng-container>
  `,
})
export class WorkoutGroup {
  private route = inject(ActivatedRoute);
  private userApi = inject(UserApiClient);
  private queryClient = inject(QueryClient);

  private workoutId = Number(this.route.snapshot.paramMap.get('id'));

  groupNameInput = signal('');
  newMemberUserId = signal('');
  newMemberRole = signal<string>(WorkoutGroupRoleMember);
  showDeleteDialog = signal(false);

  workoutQuery = injectQuery(() => ({
    queryKey: workoutKeys.detail(this.workoutId),
    queryFn: () => this.userApi.fetchWorkout(this.workoutId),
  }));

  groupsQuery = injectQuery(() => ({
    queryKey: workoutGroupKeys.list(),
    queryFn: () => this.userApi.fetchWorkoutGroups(),
  }));

  group = computed(() => {
    const groups = this.groupsQuery.data();
    if (!groups) return undefined;
    const found = groups.find((g) => g.workoutId === this.workoutId);
    if (found && !this.groupNameInput()) {
      this.groupNameInput.set(found.name);
    }
    return found;
  });

  membershipsQuery = injectQuery(() => {
    const g = this.group();
    return {
      queryKey: workoutGroupKeys.memberships(g?.id ?? 0),
      queryFn: () => this.userApi.fetchWorkoutGroupMemberships({ groupId: g!.id }),
      enabled: !!g,
    };
  });

  headerText = computed(() => {
    const workout = this.workoutQuery.data();
    return workout ? `Group — ${workout.name}` : 'Group';
  });

  createGroupMutation = injectMutation(() => ({
    mutationFn: (data: { name: string; workoutId: number }) =>
      this.userApi.createWorkoutGroup(data),
    onSuccess: () => {
      this.queryClient.invalidateQueries({ queryKey: workoutGroupKeys.all() });
    },
  }));

  updateGroupMutation = injectMutation(() => ({
    mutationFn: (data: { id: number; name: string }) =>
      this.userApi.updateWorkoutGroup(data.id, { name: data.name }),
    onSuccess: () => {
      this.queryClient.invalidateQueries({ queryKey: workoutGroupKeys.all() });
    },
  }));

  deleteGroupMutation = injectMutation(() => ({
    mutationFn: (id: number) => this.userApi.deleteWorkoutGroup(id),
    onSuccess: () => {
      this.showDeleteDialog.set(false);
      this.groupNameInput.set('');
      this.queryClient.invalidateQueries({ queryKey: workoutGroupKeys.all() });
    },
  }));

  addMemberMutation = injectMutation(() => ({
    mutationFn: (data: { groupId: number; userId: string; role: string }) =>
      this.userApi.createWorkoutGroupMembership(data),
    onSuccess: () => {
      this.newMemberUserId.set('');
      const g = this.group();
      if (g) this.queryClient.invalidateQueries({ queryKey: workoutGroupKeys.memberships(g.id) });
    },
  }));

  updateMemberRoleMutation = injectMutation(() => ({
    mutationFn: (data: { id: number; role: string }) =>
      this.userApi.updateWorkoutGroupMembership(data.id, { role: data.role }),
    onSuccess: () => {
      const g = this.group();
      if (g) this.queryClient.invalidateQueries({ queryKey: workoutGroupKeys.memberships(g.id) });
    },
  }));

  removeMemberMutation = injectMutation(() => ({
    mutationFn: (id: number) => this.userApi.deleteWorkoutGroupMembership(id),
    onSuccess: () => {
      const g = this.group();
      if (g) this.queryClient.invalidateQueries({ queryKey: workoutGroupKeys.memberships(g.id) });
    },
  }));

  createGroup() {
    this.createGroupMutation.mutate({
      name: this.groupNameInput(),
      workoutId: this.workoutId,
    });
  }

  updateGroupName() {
    const g = this.group();
    if (g) {
      this.updateGroupMutation.mutate({ id: g.id, name: this.groupNameInput() });
    }
  }

  deleteGroup() {
    const g = this.group();
    if (g) {
      this.deleteGroupMutation.mutate(g.id);
    }
  }

  addMember() {
    const g = this.group();
    if (g) {
      this.addMemberMutation.mutate({
        groupId: g.id,
        userId: this.newMemberUserId(),
        role: this.newMemberRole(),
      });
    }
  }

  updateMemberRole(membershipId: number, role: string) {
    this.updateMemberRoleMutation.mutate({ id: membershipId, role });
  }

  removeMember(membershipId: number) {
    this.removeMemberMutation.mutate(membershipId);
  }
}
