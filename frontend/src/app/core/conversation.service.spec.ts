import { TestBed } from '@angular/core/testing';
import { provideHttpClient } from '@angular/common/http';
import {
  HttpTestingController,
  provideHttpClientTesting,
} from '@angular/common/http/testing';

import { ConversationService } from './conversation.service';
import { API_BASE } from '../api';
import { Conversation } from '../models/conversation';

const sample: Conversation = {
  id: 'c1',
  customerName: 'John Carter',
  customerEmail: 'john.carter@example.com',
  subject: 'Cannot log into my account',
  status: 'OPEN',
  priority: 'HIGH',
  createdAt: '2026-08-31T09:00:00Z',
};

describe('ConversationService', () => {
  let service: ConversationService;
  let http: HttpTestingController;

  beforeEach(() => {
    TestBed.configureTestingModule({
      providers: [provideHttpClient(), provideHttpClientTesting()],
    });
    service = TestBed.inject(ConversationService);
    http = TestBed.inject(HttpTestingController);
  });

  afterEach(() => http.verify());

  it('load() sends only the filters that are set, and stores the result', () => {
    service.load({ status: 'OPEN', search: 'john' });

    const req = http.expectOne((r) => r.url === `${API_BASE}/api/conversations`);
    expect(req.request.method).toBe('GET');
    expect(req.request.params.get('status')).toBe('OPEN');
    expect(req.request.params.get('search')).toBe('john');
    expect(req.request.params.has('priority')).toBe(false);

    req.flush([sample]);

    expect(service.loading()).toBe(false);
    expect(service.conversations()).toEqual([sample]);
  });

  it('update() PATCHes the change and updates the row in place', () => {
    service.load();
    http.expectOne(`${API_BASE}/api/conversations`).flush([sample]);

    service.update('c1', { status: 'RESOLVED' });

    const req = http.expectOne(`${API_BASE}/api/conversations/c1`);
    expect(req.request.method).toBe('PATCH');
    expect(req.request.body).toEqual({ status: 'RESOLVED' });

    req.flush({ ...sample, status: 'RESOLVED' });

    expect(service.conversations()[0].status).toBe('RESOLVED');
  });

  it('load() surfaces an error message when the request fails', () => {
    service.load();
    http
      .expectOne(`${API_BASE}/api/conversations`)
      .flush('boom', { status: 500, statusText: 'Server Error' });

    expect(service.error()).toBeTruthy();
    expect(service.loading()).toBe(false);
  });
});
