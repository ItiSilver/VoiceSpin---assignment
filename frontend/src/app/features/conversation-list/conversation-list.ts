import { Component, computed, inject } from '@angular/core';
import { DatePipe } from '@angular/common';
import { FormsModule } from '@angular/forms';

import { ConversationService } from '../../core/conversation.service';
import {
  Priority,
  PRIORITIES,
  Status,
  STATUSES,
} from '../../models/conversation';

@Component({
  selector: 'app-conversation-list',
  imports: [FormsModule, DatePipe],
  templateUrl: './conversation-list.html',
  styleUrl: './conversation-list.css',
})
export class ConversationList {
  private readonly service = inject(ConversationService);

  protected readonly conversations = this.service.conversations;
  protected readonly loading = this.service.loading;
  protected readonly error = this.service.error;
  protected readonly selectedId = computed(() => this.service.selected()?.id ?? null);

  protected readonly statuses = STATUSES;
  protected readonly priorities = PRIORITIES;

  protected search = '';
  protected status: Status | '' = '';
  protected priority: Priority | '' = '';

  private debounceId?: ReturnType<typeof setTimeout>;

  constructor() {
    this.reload();
  }

  protected onSearchInput(value: string): void {
    this.search = value;
    clearTimeout(this.debounceId);
    this.debounceId = setTimeout(() => this.reload(), 300);
  }

  protected onFilterChange(): void {
    this.reload();
  }

  protected select(id: string): void {
    this.service.select(id);
  }

  protected reload(): void {
    this.service.load({
      search: this.search,
      status: this.status || undefined,
      priority: this.priority || undefined,
    });
  }
}
