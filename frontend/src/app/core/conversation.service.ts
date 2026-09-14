import { Injectable, inject, signal } from '@angular/core';
import { HttpClient, HttpParams } from '@angular/common/http';
import { API_BASE } from '../api';
import {
  Conversation,
  ConversationFilters,
  ConversationPatch,
} from '../models/conversation';

/**
 * ConversationService owns all HTTP calls AND the UI state for this feature.
 *
 * State is held in signals (think React `useState`, but framework-native):
 * components read `conversations()`, `selected()`, `loading()`, `error()`
 * and re-render automatically when they change. There's no NgRx here — the
 * app is small enough that a single service is the simplest thing that works.
 */
@Injectable({ providedIn: 'root' })
export class ConversationService {
  private readonly http = inject(HttpClient);
  private readonly baseUrl = `${API_BASE}/api/conversations`;

  readonly conversations = signal<Conversation[]>([]);
  readonly selected = signal<Conversation | null>(null);
  readonly loading = signal(false);
  readonly error = signal<string | null>(null);

  /** GET /api/conversations, applying whichever filters are set. */
  load(filters: ConversationFilters = {}): void {
    let params = new HttpParams();
    if (filters.status) params = params.set('status', filters.status);
    if (filters.priority) params = params.set('priority', filters.priority);
    if (filters.search?.trim()) params = params.set('search', filters.search.trim());

    this.loading.set(true);
    this.error.set(null);

    this.http.get<Conversation[]>(this.baseUrl, { params }).subscribe({
      next: (list) => {
        this.conversations.set(list);
        this.loading.set(false);

        // If the currently-selected conversation fell outside the new
        // filter, drop the selection so the detail pane stays consistent.
        const current = this.selected();
        if (current && !list.some((c) => c.id === current.id)) {
          this.selected.set(null);
        }
      },
      error: () => {
        this.error.set('Could not load conversations. Is the API running on :8080?');
        this.loading.set(false);
      },
    });
  }

  /** GET /api/conversations/:id and show it in the detail pane. */
  select(id: string): void {
    // Show what we already have from the list immediately, then refresh.
    this.selected.set(this.conversations().find((c) => c.id === id) ?? null);

    this.http.get<Conversation>(`${this.baseUrl}/${id}`).subscribe({
      next: (conversation) => this.selected.set(conversation),
      error: () => this.error.set('Could not load that conversation.'),
    });
  }

  clearSelection(): void {
    this.selected.set(null);
  }

  /** PATCH /api/conversations/:id, then update the row + detail in place. */
  update(id: string, patch: ConversationPatch): void {
    this.http.patch<Conversation>(`${this.baseUrl}/${id}`, patch).subscribe({
      next: (updated) => {
        this.conversations.update((list) =>
          list.map((c) => (c.id === updated.id ? updated : c)),
        );
        if (this.selected()?.id === updated.id) {
          this.selected.set(updated);
        }
      },
      error: () => this.error.set('Could not save your change. Please try again.'),
    });
  }
}
