import { Component } from '@angular/core';

import { ConversationList } from './features/conversation-list/conversation-list';
import { ConversationDetail } from './features/conversation-detail/conversation-detail';

@Component({
  selector: 'app-root',
  imports: [ConversationList, ConversationDetail],
  templateUrl: './app.html',
  styleUrl: './app.css',
})
export class App {}
