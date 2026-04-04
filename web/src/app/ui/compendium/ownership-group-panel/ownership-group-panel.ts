import { Component, inject, input, signal } from '@angular/core';
import { injectQuery, injectMutation, QueryClient } from '@tanstack/angular-query-experimental';
import { CompendiumApiClient } from '$core/api-clients/compendium-api-client';
import { ownershipGroupKeys } from '$core/query-keys';
import { OwnershipGroupMembership } from '$generated/ownershipgroup';

@Component({
  selector: 'app-ownership-group-panel',
  template: `
    @if (ownershipGroupId()) {
      <div class="mt-6 rounded-lg border border-gray-200 p-6 dark:border-gray-700">
        <h2 class="mb-4 text-lg font-semibold text-gray-900 dark:text-gray-100">Shared with</h2>

        @if (membershipsQuery.data(); as memberships) {
          @if (nonOwnerMembers(memberships).length === 0) {
            <p class="mb-4 text-sm text-gray-500 dark:text-gray-400">Not shared with anyone yet.</p>
          } @else {
            <table class="mb-4 min-w-full divide-y divide-gray-200 dark:divide-gray-700">
              <thead>
                <tr>
                  <th
                    class="px-4 py-2 text-left text-xs font-medium tracking-wider text-gray-500 uppercase dark:text-gray-400"
                  >
                    User
                  </th>
                  <th
                    class="px-4 py-2 text-left text-xs font-medium tracking-wider text-gray-500 uppercase dark:text-gray-400"
                  >
                    Role
                  </th>
                  <th class="px-4 py-2"></th>
                </tr>
              </thead>
              <tbody class="divide-y divide-gray-200 dark:divide-gray-700">
                @for (m of nonOwnerMembers(memberships); track m.id) {
                  <tr>
                    <td class="px-4 py-2 text-sm text-gray-900 dark:text-gray-100">
                      {{ m.userId }}
                    </td>
                    <td class="px-4 py-2 text-sm">
                      <span
                        class="inline-flex items-center rounded-full px-2.5 py-0.5 text-xs font-medium"
                        [class]="
                          m.role === 'admin'
                            ? 'bg-purple-100 text-purple-800 dark:bg-purple-900 dark:text-purple-200'
                            : 'bg-blue-100 text-blue-800 dark:bg-blue-900 dark:text-blue-200'
                        "
                      >
                        {{ m.role }}
                      </span>
                    </td>
                    <td class="px-4 py-2 text-right">
                      <button
                        type="button"
                        (click)="removeMember(m.id)"
                        class="text-sm text-red-600 hover:underline dark:text-red-400"
                      >
                        Remove
                      </button>
                    </td>
                  </tr>
                }
              </tbody>
            </table>
          }
        }

        @if (membershipsQuery.isPending()) {
          <p class="mb-4 text-sm text-gray-500 dark:text-gray-400">Loading...</p>
        }

        <!-- Add member form -->
        <div class="flex gap-2">
          <input
            type="text"
            [value]="newMemberUserId()"
            (input)="newMemberUserId.set($any($event.target).value)"
            class="block w-full rounded-md border border-gray-300 px-3 py-2 text-sm focus:border-blue-500 focus:ring-blue-500 dark:border-gray-600 dark:bg-gray-800 dark:text-gray-100"
            placeholder="User ID"
          />
          <button
            type="button"
            (click)="addMember()"
            [disabled]="addMemberMutation.isPending() || !newMemberUserId()"
            class="shrink-0 rounded-md bg-blue-600 px-4 py-2 text-sm font-medium text-white hover:bg-blue-700 disabled:opacity-50"
          >
            Add
          </button>
        </div>
      </div>
    }
  `,
})
export class OwnershipGroupPanel {
  ownershipGroupId = input.required<number>();

  private api = inject(CompendiumApiClient);
  private queryClient = inject(QueryClient);

  newMemberUserId = signal('');

  membershipsQuery = injectQuery(() => ({
    queryKey: ownershipGroupKeys.memberships(this.ownershipGroupId()),
    queryFn: () => this.api.fetchOwnershipGroupMemberships(this.ownershipGroupId()),
    enabled: !!this.ownershipGroupId(),
  }));

  addMemberMutation = injectMutation(() => ({
    mutationFn: (data: { userId: string }) =>
      this.api.createOwnershipGroupMembership(this.ownershipGroupId(), data),
    onSuccess: () => {
      this.queryClient.invalidateQueries({
        queryKey: ownershipGroupKeys.memberships(this.ownershipGroupId()),
      });
      this.newMemberUserId.set('');
    },
  }));

  removeMemberMutation = injectMutation(() => ({
    mutationFn: (id: number) => this.api.deleteOwnershipGroupMembership(id),
    onSuccess: () => {
      this.queryClient.invalidateQueries({
        queryKey: ownershipGroupKeys.memberships(this.ownershipGroupId()),
      });
    },
  }));

  nonOwnerMembers(memberships: OwnershipGroupMembership[]) {
    return memberships.filter((m) => m.role !== 'owner');
  }

  addMember() {
    const userId = this.newMemberUserId().trim();
    if (userId) {
      this.addMemberMutation.mutate({ userId });
    }
  }

  removeMember(id: number) {
    this.removeMemberMutation.mutate(id);
  }
}
