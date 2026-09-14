export type Status = 'OPEN' | 'IN_PROGRESS' | 'RESOLVED';
export type Priority = 'LOW' | 'MEDIUM' | 'HIGH';

export const STATUSES: readonly Status[] = ['OPEN', 'IN_PROGRESS', 'RESOLVED'];
export const PRIORITIES: readonly Priority[] = ['LOW', 'MEDIUM', 'HIGH'];

export type Conversation = {
  id: string;
  customerName: string;
  customerEmail: string;
  subject: string;
  status: Status;
  priority: Priority;
  createdAt: string;
}

export type ConversationFilters = {
  status?: Status | '';
  priority?: Priority | '';
  search?: string;
}

export type ConversationPatch = {
  status?: Status;
  priority?: Priority;
}
