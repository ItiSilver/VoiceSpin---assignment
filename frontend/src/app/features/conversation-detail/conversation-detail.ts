import { Component, inject } from '@angular/core';
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
  selector: 'app-conversation-detail',
  imports: [DatePipe, FormsModule],
  templateUrl: './conversation-detail.html',
  styleUrl: './conversation-detail.css',
})
export class ConversationDetail {
  private readonly service = inject(ConversationService);

  protected readonly conversation = this.service.selected;
  protected readonly statuses = STATUSES;
  protected readonly priorities = PRIORITIES;

  protected changeStatus(status: Status): void {
    const current = this.conversation();
    if (current && current.status !== status) {
      this.service.update(current.id, { status });
    }
  }

  protected changePriority(priority: Priority): void {
    const current = this.conversation();
    if (current && current.priority !== priority) {
      this.service.update(current.id, { priority });
    }
  }
}
