import { Component, inject, input, output, signal } from '@angular/core';
import { form, FormField } from '@angular/forms/signals';
import { injectQuery, injectMutation, QueryClient } from '@tanstack/angular-query-experimental';
import { CompendiumApiClient } from '$core/api-clients/compendium-api-client';
import { ownershipGroupKeys } from '$core/query-keys';
import { HlmDialogImports } from '@spartan-ng/helm/dialog';
import { HlmButton } from '@spartan-ng/helm/button';

@Component({
  selector: 'app-ownership-group-panel',
  imports: [HlmDialogImports, HlmButton, FormField],
  template: `
    <hlm-dialog [state]="open() ? 'open' : 'closed'" (closed)="onClose(); closed.emit()">
      <ng-template hlmDialogPortal>
        <hlm-dialog-content [showCloseButton]="false" class="max-h-[90dvh] overflow-y-auto">
          <hlm-dialog-header>
            <h3 hlmDialogTitle>Shared with</h3>
          </hlm-dialog-header>

          @if (membershipsQuery.data(); as memberships) {
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
                @for (m of memberships; track m.id) {
                  <tr>
                    <td class="px-4 py-2 text-sm text-gray-900 dark:text-gray-100">
                      {{ m.userId }}
                    </td>
                    <td class="px-4 py-2 text-sm">
                      <span
                        class="inline-flex items-center rounded-full px-2.5 py-0.5 text-xs font-medium"
                        [class]="
                          m.role === 'owner'
                            ? 'bg-green-100 text-green-800 dark:bg-green-900 dark:text-green-200'
                            : m.role === 'admin'
                              ? 'bg-purple-100 text-purple-800 dark:bg-purple-900 dark:text-purple-200'
                              : 'bg-blue-100 text-blue-800 dark:bg-blue-900 dark:text-blue-200'
                        "
                      >
                        {{ m.role }}
                      </span>
                    </td>
                    <td class="px-4 py-2 text-right">
                      @if (m.role !== 'owner') {
                        <button
                          type="button"
                          (click)="removeMember(m.id)"
                          class="text-sm text-red-600 hover:underline dark:text-red-400"
                        >
                          Remove
                        </button>
                      }
                    </td>
                  </tr>
                }
              </tbody>
            </table>
          }

          @if (membershipsQuery.isPending()) {
            <p class="mb-4 text-sm text-gray-500 dark:text-gray-400">Loading...</p>
          }

          <!-- Add member form -->
          <form (submit)="addMember(); $event.preventDefault()" class="flex gap-2">
            <input
              [formField]="addMemberForm.userId"
              class="block w-full rounded-md border border-gray-300 px-3 py-2 text-sm focus:border-blue-500 focus:ring-blue-500 dark:border-gray-600 dark:bg-gray-800 dark:text-gray-100"
              placeholder="User ID"
            />
            <button
              type="submit"
              [disabled]="addMemberMutation.isPending() || !model().userId"
              class="shrink-0 rounded-md bg-blue-600 px-4 py-2 text-sm font-medium text-white hover:bg-blue-700 disabled:opacity-50"
            >
              Add
            </button>
          </form>

          <hlm-dialog-footer>
            <button hlmBtn variant="outline" hlmDialogClose>Close</button>
          </hlm-dialog-footer>
        </hlm-dialog-content>
      </ng-template>
    </hlm-dialog>
  `,
})
export class OwnershipGroupPanel {
  ownershipGroupId = input.required<number>();
  open = input(false);
  closed = output();

  private api = inject(CompendiumApiClient);
  private queryClient = inject(QueryClient);

  model = signal({ userId: '' });
  addMemberForm = form(this.model);

  membershipsQuery = injectQuery(() => ({
    queryKey: ownershipGroupKeys.memberships(this.ownershipGroupId()),
    queryFn: () => this.api.fetchOwnershipGroupMemberships(this.ownershipGroupId()),
    enabled: !!this.ownershipGroupId() && this.open(),
  }));

  addMemberMutation = injectMutation(() => ({
    mutationFn: (data: { userId: string }) =>
      this.api.createOwnershipGroupMembership(this.ownershipGroupId(), data),
    onSuccess: () => {
      this.queryClient.invalidateQueries({
        queryKey: ownershipGroupKeys.memberships(this.ownershipGroupId()),
      });
      this.model.set({ userId: '' });
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

  onClose() {
    this.model.set({ userId: '' });
  }

  addMember() {
    const userId = this.model().userId.trim();
    if (userId) {
      this.addMemberMutation.mutate({ userId });
    }
  }

  removeMember(id: number) {
    this.removeMemberMutation.mutate(id);
  }
}
